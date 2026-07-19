package workflows

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/tilebox/tilebox-go/internal/span"
	obstracer "github.com/tilebox/tilebox-go/observability/tracer"
	workflowsv1 "github.com/tilebox/tilebox-go/protogen/workflows/v1"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	workerAddressEnvironmentVariable = "TILEBOX_WORKER_ADDRESS"
	workerShutdownGracePeriod        = 5 * time.Second
	workerCleanupTimeout             = 5 * time.Second
)

type workerState uint8

const (
	workerStateNew workerState = iota
	workerStateStarting
	workerStateServing
	workerStateInitialized
	workerStateStopping
	workerStateStopped
)

type workerRuntime struct {
	executor *taskExecutor
	cleanup  func(context.Context)
}

var errInvalidWorkerConfiguration = errors.New("invalid worker configuration")

// Worker serves registered Go tasks as an execution-only child process managed by tilebox runner start.
//
// A Worker never polls for tasks, acquires or extends leases, or reports results to the Tilebox API. Its generated
// WorkerService is served on the private Unix socket supplied in TILEBOX_WORKER_ADDRESS. The managing runner remains
// responsible for polling, routing, lease management, and result reporting.
//
// Shutdown is graceful for up to five seconds. After that deadline, the active task context is canceled and the RPC
// server is stopped. Task implementations must honor context cancellation; the managing runner may terminate a child
// process that does not exit.
type Worker struct {
	workflowsv1.UnimplementedWorkerServiceServer

	mu       sync.Mutex
	state    workerState
	registry *taskRegistry
	server   *grpc.Server
	runtime  *workerRuntime
	init     *workflowsv1.InitializeRunnerRequest

	stopRequested chan struct{}
	executionSlot chan struct{}
	activeCancel  context.CancelFunc
	activeDone    chan struct{}

	shutdownGracePeriod time.Duration
	cleanupTimeout      time.Duration
}

var _ workflowsv1.WorkerServiceServer = (*Worker)(nil)

// NewWorker creates an execution-only workflow worker.
//
// Construction only allocates in-memory lifecycle and task registration state. It does not create a client, perform
// cluster lookup, access credentials, poll a queue, or make any network request.
func NewWorker() *Worker {
	executionSlot := make(chan struct{}, 1)
	executionSlot <- struct{}{}
	return &Worker{
		state:               workerStateNew,
		registry:            newTaskRegistry(),
		stopRequested:       make(chan struct{}),
		executionSlot:       executionSlot,
		shutdownGracePeriod: workerShutdownGracePeriod,
		cleanupTimeout:      workerCleanupTimeout,
	}
}

// RegisterTasks makes the worker aware of tasks it can execute.
//
// Registration is concurrency-safe and uses the same identifier validation and duplicate checks as TaskRunner.
// Registration is permanently frozen when Serve starts, including when startup subsequently fails.
func (w *Worker) RegisterTasks(tasks ...ExecutableTask) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.state != workerStateNew {
		return errors.New("cannot register tasks after worker serving has started")
	}
	return w.registry.registerTasks(tasks...)
}

// Serve binds TILEBOX_WORKER_ADDRESS, serves the generated WorkerService, and blocks until shutdown.
//
// Serve can be called once. It returns nil after context cancellation, a supported process signal, or a
// ShutdownWorker RPC. The socket is removed before Serve returns. Concurrent or repeated calls return an error.
func (w *Worker) Serve(ctx context.Context) error {
	if ctx == nil {
		return errors.New("worker serve context must not be nil")
	}

	w.mu.Lock()
	if w.state != workerStateNew {
		w.mu.Unlock()
		return errors.New("worker can only be served once")
	}
	w.state = workerStateStarting
	w.mu.Unlock()

	listener, cleanupSocket, err := listenWorker(os.Getenv(workerAddressEnvironmentVariable))
	if err != nil {
		w.finishServe()
		return fmt.Errorf("failed to start worker service: %w", err)
	}

	server := grpc.NewServer()
	workflowsv1.RegisterWorkerServiceServer(server, w)

	w.mu.Lock()
	w.server = server
	w.state = workerStateServing
	w.mu.Unlock()

	ctxSignal, stopSignals := signal.NotifyContext(ctx, runnerShutdownSignals()...)
	defer stopSignals()

	serveErrors := make(chan error, 1)
	go func() {
		serveErrors <- server.Serve(listener)
	}()

	select {
	case err = <-serveErrors:
		w.markStopping()
		w.cancelActiveTask()
		w.stopServer(server, 0)
	case <-ctxSignal.Done():
		w.markStopping()
		w.cancelActiveTask()
		w.stopServer(server, 0)
		err = <-serveErrors
	}

	w.finishServe()
	cleanupErr := cleanupSocket()
	if err == nil || errors.Is(err, grpc.ErrServerStopped) || ctxSignal.Err() != nil {
		if cleanupErr != nil {
			return fmt.Errorf("failed to clean up worker socket: %w", cleanupErr)
		}
		return nil
	}
	return errors.Join(fmt.Errorf("worker service stopped unexpectedly: %w", err), cleanupErr)
}

// ListRegisteredTasks implements workflows.v1.WorkerService and is available before initialization.
func (w *Worker) ListRegisteredTasks(context.Context, *emptypb.Empty) (*workflowsv1.TaskIdentifiers, error) {
	return workflowsv1.TaskIdentifiers_builder{Identifiers: w.registry.identifiers()}.Build(), nil
}

// InitializeWorker implements workflows.v1.WorkerService initialization.
//
// A semantically identical repeated request is idempotent. Any conflicting request is rejected, and a failed
// initialization may be retried. Request API connection fields are authoritative; initialization does not fall back
// to TILEBOX_API_URL or TILEBOX_API_KEY.
func (w *Worker) InitializeWorker(ctx context.Context, request *workflowsv1.InitializeRunnerRequest) (*workflowsv1.InitializeRunnerResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "worker initialization request is required")
	}
	if err := request.GetRunnerId().CheckValid(); err != nil {
		return nil, status.Error(codes.InvalidArgument, "worker initialization runner ID is invalid")
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	switch w.state {
	case workerStateInitialized:
		if equivalentWorkerInitialization(w.init, request) {
			return &workflowsv1.InitializeRunnerResponse{}, nil
		}
		return nil, status.Error(codes.FailedPrecondition, "worker is already initialized with different configuration")
	case workerStateServing:
		// Continue below.
	case workerStateStopping, workerStateStopped:
		return nil, status.Error(codes.FailedPrecondition, "worker is shutting down")
	case workerStateNew, workerStateStarting:
		return nil, status.Error(codes.FailedPrecondition, "worker is not serving")
	}

	runtime, err := newWorkerRuntime(ctx, w.registry, request, otel.GetMeterProvider())
	if err != nil {
		if !errors.Is(err, errInvalidWorkerConfiguration) {
			return nil, status.Error(codes.Internal, "worker initialization failed")
		}
		return nil, status.Error(codes.InvalidArgument, "worker initialization configuration is invalid")
	}
	w.runtime = runtime
	w.init = proto.CloneOf(request)
	w.state = workerStateInitialized
	return &workflowsv1.InitializeRunnerResponse{}, nil
}

// ExecuteTask implements workflows.v1.WorkerService task execution.
//
// Calls are serialized because a release worker is managed by one routing runner. RPC cancellation is propagated to
// the task context. Calls made before successful initialization or while shutting down return an infrastructure-class
// failed task response rather than a transport error.
func (w *Worker) ExecuteTask(ctx context.Context, task *workflowsv1.Task) (*workflowsv1.ExecuteTaskResponse, error) {
	select {
	case <-ctx.Done():
		return nil, status.FromContextError(ctx.Err()).Err()
	case <-w.stopRequested:
		return w.failedResponse(task, errors.New("worker is shutting down")), nil
	case <-w.executionSlot:
	}
	defer func() { w.executionSlot <- struct{}{} }()

	w.mu.Lock()
	if w.state != workerStateInitialized || w.runtime == nil {
		stopping := w.state == workerStateStopping || w.state == workerStateStopped
		w.mu.Unlock()
		if stopping {
			return w.failedResponse(task, errors.New("worker is shutting down")), nil
		}
		return w.failedResponse(task, errors.New("worker is not initialized")), nil
	}

	runtime := w.runtime
	if task.GetJob() == nil {
		w.mu.Unlock()
		return w.failedResponse(task, errors.New("task has no job")), nil
	}
	executionContext, cancel := context.WithCancel(ctx)
	done := make(chan struct{})
	w.activeCancel = cancel
	w.activeDone = done
	w.mu.Unlock()

	response, err := runtime.executor.ExecuteTask(executionContext, task)
	cancel()
	close(done)

	w.mu.Lock()
	if w.activeDone == done {
		w.activeCancel = nil
		w.activeDone = nil
	}
	w.mu.Unlock()
	return response, err
}

// ShutdownWorker implements workflows.v1.WorkerService graceful shutdown.
//
// Server shutdown runs asynchronously so this RPC can return before GracefulStop waits for its own handler.
func (w *Worker) ShutdownWorker(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	w.mu.Lock()
	switch w.state {
	case workerStateStopping, workerStateStopped:
		w.mu.Unlock()
		return &emptypb.Empty{}, nil
	case workerStateServing, workerStateInitialized:
		w.state = workerStateStopping
		close(w.stopRequested)
		server := w.server
		gracePeriod := w.shutdownGracePeriod
		w.mu.Unlock()
		go w.stopServer(server, gracePeriod)
		return &emptypb.Empty{}, nil
	case workerStateNew, workerStateStarting:
		w.mu.Unlock()
		return nil, status.Error(codes.FailedPrecondition, "worker is not serving")
	}
	w.mu.Unlock()
	return nil, status.Error(codes.Internal, "worker has invalid lifecycle state")
}

func (w *Worker) markStopping() {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.state != workerStateStopping && w.state != workerStateStopped {
		w.state = workerStateStopping
		close(w.stopRequested)
	}
}

func (w *Worker) cancelActiveTask() {
	w.mu.Lock()
	cancel := w.activeCancel
	w.mu.Unlock()
	if cancel != nil {
		cancel()
	}
}

func (w *Worker) stopServer(server *grpc.Server, gracePeriod time.Duration) {
	if server == nil {
		return
	}

	if gracePeriod <= 0 {
		w.cancelActiveTask()
		server.Stop()
		return
	}

	gracefulStopDone := make(chan struct{})
	go func() {
		server.GracefulStop()
		close(gracefulStopDone)
	}()

	timer := time.NewTimer(gracePeriod)
	defer timer.Stop()
	select {
	case <-gracefulStopDone:
		return
	case <-timer.C:
		w.cancelActiveTask()
		server.Stop()
	}
}

func (w *Worker) finishServe() {
	w.mu.Lock()
	runtime := w.runtime
	w.runtime = nil
	w.init = nil
	w.server = nil
	w.state = workerStateStopped
	cleanupTimeout := w.cleanupTimeout
	w.mu.Unlock()

	if runtime != nil {
		cleanupContext, cancel := context.WithTimeout(context.Background(), cleanupTimeout)
		runtime.cleanup(cleanupContext)
		cancel()
	}
}

func newWorkerRuntime(ctx context.Context, registry *taskRegistry, request *workflowsv1.InitializeRunnerRequest, meterProvider metric.MeterProvider) (*workerRuntime, error) {
	apiURL := defaultTileboxAPIURL
	apiToken := ""
	if connection := request.GetApiConnection(); connection != nil {
		if strings.TrimSpace(connection.GetUrl()) != "" {
			apiURL = connection.GetUrl()
		}
		apiToken = connection.GetToken()
	}

	cfg := newClientConfig([]ClientOption{WithURL(apiURL), WithAPIKey(apiToken)})
	cleanupFunctions := make([]func(context.Context), 0, 3)

	if apiToken != "" && isHTTPURL(apiURL) {
		traceEndpoint, ok := otlpEndpointURL(apiURL, otlpTracesPath)
		if !ok {
			return nil, fmt.Errorf("%w: API URL", errInvalidWorkerConfiguration)
		}
		tracerProvider, cleanupTracing, err := obstracer.NewOtelProvider(ctx, tileboxTelemetryService,
			obstracer.WithEndpointURL(traceEndpoint),
			obstracer.WithHeaders(tileboxTelemetryHeaders(apiToken)),
		)
		if err != nil {
			return nil, errors.New("failed to configure worker tracing")
		}
		cfg.tracerProvider = tracerProvider
		cleanupFunctions = append(cleanupFunctions, cleanupTracing)

		logHandler, cleanupLogging, err := NewTileboxLogHandler(ctx, apiURL, apiToken, slog.LevelInfo)
		if err != nil {
			cleanupTracing(ctx)
			return nil, errors.New("failed to configure worker logging")
		}
		removeLogHandler := configureRemovableSlogHandler(logHandler)
		cleanupFunctions = append(cleanupFunctions, cleanupLogging)
		cleanupFunctions = append(cleanupFunctions, func(context.Context) { removeLogHandler() })
	}

	if closeIdleConnections, ok := cfg.httpClient.(interface{ CloseIdleConnections() }); ok {
		cleanupFunctions = append(cleanupFunctions, func(context.Context) { closeIdleConnections.CloseIdleConnections() })
	}

	client := newClient(cfg)
	logger := workerLogger(request)
	executor, err := newTaskExecutor(registry, request.GetCluster().GetSlug(), client, cfg.tracerProvider.Tracer(otelTracerName), logger, meterProvider)
	if err != nil {
		for _, cleanup := range slices.Backward(cleanupFunctions) {
			cleanup(ctx)
		}
		return nil, err
	}
	executor.redactError, executor.redactString = secretRedactors(apiToken, os.Getenv("TILEBOX_API_KEY"))

	logger.DebugContext(span.ContextWithTraceParent(ctx, request.GetTraceParent()), "worker initialized")
	return &workerRuntime{
		executor: executor,
		cleanup: func(ctx context.Context) {
			for _, cleanup := range slices.Backward(cleanupFunctions) {
				cleanup(ctx)
			}
		},
	}, nil
}

func workerLogger(request *workflowsv1.InitializeRunnerRequest) *slog.Logger {
	attributes := []any{
		slog.String("runner_id", request.GetRunnerId().AsUUID().String()),
		slog.String("cluster", request.GetCluster().GetSlug()),
		slog.String("workflow", request.GetWorkflow().GetSlug()),
	}
	if releases := request.GetWorkflow().GetReleases(); len(releases) > 0 {
		attributes = append(attributes, slog.String("release_id", releases[0].GetId().AsUUID().String()))
	}
	return slog.Default().With(attributes...)
}

func secretRedactors(secrets ...string) (func(error) error, func(string) string) {
	redactString := func(message string) string {
		for _, secret := range secrets {
			if secret != "" {
				message = strings.ReplaceAll(message, secret, "[REDACTED]")
			}
		}
		return message
	}
	redactError := func(err error) error {
		if err == nil {
			return nil
		}
		return errors.New(redactString(err.Error()))
	}
	return redactError, redactString
}

func equivalentWorkerInitialization(first, second *workflowsv1.InitializeRunnerRequest) bool {
	first = proto.CloneOf(first)
	second = proto.CloneOf(second)
	first.ClearTraceParent()
	second.ClearTraceParent()
	return proto.Equal(first, second)
}

func (w *Worker) failedResponse(task *workflowsv1.Task, err error) *workflowsv1.ExecuteTaskResponse {
	w.mu.Lock()
	runtime := w.runtime
	w.mu.Unlock()

	redactError, redactString := secretRedactors(os.Getenv("TILEBOX_API_KEY"))
	if runtime != nil {
		redactError = runtime.executor.redactError
		redactString = runtime.executor.redactString
	}
	err = redactError(err)
	return workflowsv1.ExecuteTaskResponse_builder{FailedTask: workflowsv1.TaskFailedRequest_builder{
		TaskId:           task.GetId(),
		Display:          failedTaskDisplay(redactString(task.GetDisplay()), err),
		WasWorkflowError: false,
	}.Build()}.Build()
}
