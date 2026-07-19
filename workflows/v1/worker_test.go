//go:build unix

package workflows

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	examplesv1 "github.com/tilebox/tilebox-go/protogen/examples/v1"
	tileboxv1 "github.com/tilebox/tilebox-go/protogen/tilebox/v1"
	workflowsv1 "github.com/tilebox/tilebox-go/protogen/workflows/v1"
	"github.com/tilebox/tilebox-go/protogen/workflows/v1/workflowsv1connect"
	"go.opentelemetry.io/otel/trace/noop"
	"golang.org/x/net/http2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

type workerRegistrationTask struct{}

func (*workerRegistrationTask) Execute(context.Context) error { return nil }

type workerInvalidIdentifierTask struct{}

func (*workerInvalidIdentifierTask) Execute(context.Context) error { return nil }
func (*workerInvalidIdentifierTask) Identifier() TaskIdentifier {
	return NewTaskIdentifier("", "not-a-version")
}

type workerNilIdentifierTask struct{}

func (*workerNilIdentifierTask) Execute(context.Context) error { return nil }
func (*workerNilIdentifierTask) Identifier() TaskIdentifier    { return nil }

type workerValueTask struct{}

func (workerValueTask) Execute(context.Context) error { return nil }

var workerJSONValues chan string

type workerJSONTask struct {
	Value string `json:"value"`
}

func (*workerJSONTask) Identifier() TaskIdentifier {
	return NewTaskIdentifier("tilebox.com/tests/json", "v1.0")
}

func (task *workerJSONTask) Execute(context.Context) error {
	workerJSONValues <- task.Value
	return nil
}

var workerProtoValues chan int64

type workerProtoTask struct {
	examplesv1.SpawnWorkflowTreeTask
}

func (*workerProtoTask) Identifier() TaskIdentifier {
	return NewTaskIdentifier("tilebox.com/tests/protobuf", "v1.0")
}

func (task *workerProtoTask) Execute(context.Context) error {
	workerProtoValues <- task.GetDepth()
	return nil
}

type workerDispatchV1Task struct{}

func (*workerDispatchV1Task) Identifier() TaskIdentifier {
	return NewTaskIdentifier("tilebox.com/tests/dispatch", "v1.0")
}
func (*workerDispatchV1Task) Execute(context.Context) error { return nil }

type workerDispatchV2Task struct{}

func (*workerDispatchV2Task) Identifier() TaskIdentifier {
	return NewTaskIdentifier("tilebox.com/tests/dispatch", "v2.0")
}
func (*workerDispatchV2Task) Execute(context.Context) error { return nil }

type workerChildTask struct {
	Value string `json:"value"`
}

func (*workerChildTask) Identifier() TaskIdentifier {
	return NewTaskIdentifier("tilebox.com/tests/child", "v1.0")
}

type workerTaskContextResult struct {
	cluster   string
	hasClient bool
	err       error
}

var workerTaskContextResults chan workerTaskContextResult

type workerContextTask struct{}

func (*workerContextTask) Identifier() TaskIdentifier {
	return NewTaskIdentifier("tilebox.com/tests/context", "v1.0")
}

func (*workerContextTask) Execute(ctx context.Context) error {
	client, clientErr := GetClient(ctx)
	cluster, clusterErr := GetCurrentCluster(ctx)
	if clientErr == nil && clusterErr == nil {
		progress := Progress("work")
		if err := progress.Add(ctx, 3); err != nil {
			clientErr = err
		} else if err := progress.Done(ctx, 2); err != nil {
			clientErr = err
		} else if _, err := SubmitSubtask(ctx, &workerChildTask{Value: "child"}); err != nil {
			clientErr = err
		} else if err := SetTaskDisplay(ctx, "context ready"); err != nil {
			clientErr = err
		}
	}
	workerTaskContextResults <- workerTaskContextResult{
		cluster:   cluster,
		hasClient: client != nil,
		err:       errors.Join(clientErr, clusterErr),
	}
	return errors.Join(clientErr, clusterErr)
}

type workerUserErrorTask struct{}

func (*workerUserErrorTask) Identifier() TaskIdentifier {
	return NewTaskIdentifier("tilebox.com/tests/user-error", "v1.0")
}

func (*workerUserErrorTask) Execute(ctx context.Context) error {
	if err := Progress("before-error").Add(ctx, 5); err != nil {
		return err
	}
	return errors.New("user task failed")
}

type workerPanicTask struct{}

func (*workerPanicTask) Identifier() TaskIdentifier {
	return NewTaskIdentifier("tilebox.com/tests/panic", "v1.0")
}

func (*workerPanicTask) Execute(context.Context) error { panic("user panic") }

var (
	workerCancellationStarted  chan struct{}
	workerCancellationObserved chan error
)

type workerCancellationTask struct{}

func (*workerCancellationTask) Identifier() TaskIdentifier {
	return NewTaskIdentifier("tilebox.com/tests/cancellation", "v1.0")
}

func (*workerCancellationTask) Execute(ctx context.Context) error {
	close(workerCancellationStarted)
	<-ctx.Done()
	workerCancellationObserved <- ctx.Err()
	return ctx.Err()
}

var (
	workerGracefulStarted chan struct{}
	workerGracefulRelease chan struct{}
)

type workerGracefulTask struct{}

func (*workerGracefulTask) Identifier() TaskIdentifier {
	return NewTaskIdentifier("tilebox.com/tests/graceful", "v1.0")
}

func (*workerGracefulTask) Execute(context.Context) error {
	close(workerGracefulStarted)
	<-workerGracefulRelease
	return nil
}

var (
	workerQueuedStarted       chan struct{}
	workerQueuedCanceled      chan struct{}
	workerQueuedSecondRan     atomic.Bool
	workerQueuedSecondStarted chan struct{}
)

type workerQueuedFirstTask struct{}

func (*workerQueuedFirstTask) Identifier() TaskIdentifier {
	return NewTaskIdentifier("tilebox.com/tests/queued-first", "v1.0")
}

func (*workerQueuedFirstTask) Execute(ctx context.Context) error {
	close(workerQueuedStarted)
	<-ctx.Done()
	close(workerQueuedCanceled)
	return ctx.Err()
}

type workerQueuedSecondTask struct{}

func (*workerQueuedSecondTask) Identifier() TaskIdentifier {
	return NewTaskIdentifier("tilebox.com/tests/queued-second", "v1.0")
}

func (*workerQueuedSecondTask) Execute(context.Context) error {
	workerQueuedSecondRan.Store(true)
	return nil
}

var workerSensitiveValue string

type workerSensitiveTask struct{}

func (*workerSensitiveTask) Identifier() TaskIdentifier {
	return NewTaskIdentifier("tilebox.com/tests/sensitive", "v1.0")
}

func (*workerSensitiveTask) Execute(ctx context.Context) error {
	if err := SetTaskDisplay(ctx, "credential "+workerSensitiveValue); err != nil {
		return err
	}
	return errors.New("credential failure: " + workerSensitiveValue)
}

type testWorkerServer struct {
	client     workflowsv1.WorkerServiceClient
	connection *grpc.ClientConn
	socketPath string
	cancel     context.CancelFunc
	done       chan struct{}

	errMu sync.Mutex
	err   error
}

func startTestWorker(t *testing.T, worker *Worker) *testWorkerServer {
	t.Helper()
	return startTestWorkerAt(t, worker, filepath.Join(t.TempDir(), "worker.sock"))
}

func startTestWorkerAt(t *testing.T, worker *Worker, socketPath string) *testWorkerServer {
	t.Helper()
	t.Setenv(workerAddressEnvironmentVariable, "unix://"+socketPath)

	ctx, cancel := context.WithCancel(context.Background())
	server := &testWorkerServer{socketPath: socketPath, cancel: cancel, done: make(chan struct{})}
	go func() {
		err := worker.Serve(ctx)
		server.errMu.Lock()
		server.err = err
		server.errMu.Unlock()
		close(server.done)
	}()

	connection, err := grpc.NewClient(
		"passthrough:///tilebox-worker",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
			return (&net.Dialer{}).DialContext(ctx, "unix", socketPath)
		}),
	)
	require.NoError(t, err)
	server.connection = connection
	server.client = workflowsv1.NewWorkerServiceClient(connection)

	deadline := time.Now().Add(3 * time.Second)
	for {
		ctx, cancelList := context.WithTimeout(context.Background(), 100*time.Millisecond)
		_, listErr := server.client.ListRegisteredTasks(ctx, &emptypb.Empty{})
		cancelList()
		if listErr == nil {
			break
		}
		if time.Now().After(deadline) {
			require.NoError(t, listErr, "worker never became ready")
		}
		time.Sleep(5 * time.Millisecond)
	}

	info, err := os.Lstat(socketPath)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o600), info.Mode().Perm())

	t.Cleanup(func() {
		_ = connection.Close()
		cancel()
		select {
		case <-server.done:
		case <-time.After(3 * time.Second):
			t.Errorf("worker did not stop during test cleanup")
		}
	})
	return server
}

func (s *testWorkerServer) wait(t *testing.T) error {
	t.Helper()
	select {
	case <-s.done:
	case <-time.After(3 * time.Second):
		t.Fatal("worker Serve did not return")
	}
	s.errMu.Lock()
	defer s.errMu.Unlock()
	return s.err
}

func testWorkerInitializationRequest() *workflowsv1.InitializeRunnerRequest {
	apiURL := filepath.Join(os.TempDir(), "tilebox-worker-test-api.sock")
	apiToken := ""
	traceParent := "00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01"
	return workflowsv1.InitializeRunnerRequest_builder{
		RunnerId:    tileboxv1.NewUUID(uuid.New()),
		TraceParent: &traceParent,
		Cluster: workflowsv1.Cluster_builder{
			Slug:        "test-cluster",
			DisplayName: "Test Cluster",
		}.Build(),
		Workflow: workflowsv1.Workflow_builder{
			Slug: "test-workflow",
			Name: "Test Workflow",
			Releases: []*workflowsv1.WorkflowRelease{
				workflowsv1.WorkflowRelease_builder{Id: tileboxv1.NewUUID(uuid.New())}.Build(),
			},
		}.Build(),
		ApiConnection: workflowsv1.TileboxAPIConnection_builder{
			Url:   &apiURL,
			Token: &apiToken,
		}.Build(),
	}.Build()
}

func testWorkerTask(identifier TaskIdentifier, input []byte) *workflowsv1.Task {
	return workflowsv1.Task_builder{
		Id: tileboxv1.NewUUID(uuid.New()),
		Identifier: workflowsv1.TaskIdentifier_builder{
			Name:    identifier.Name(),
			Version: identifier.Version(),
		}.Build(),
		State:   workflowsv1.TaskState_TASK_STATE_RUNNING,
		Input:   input,
		Display: proto.String(identifier.Display()),
		Job: workflowsv1.Job_builder{
			Id:          tileboxv1.NewUUID(uuid.New()),
			Name:        "worker test job",
			TraceParent: "00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01",
		}.Build(),
	}.Build()
}

func initializeTestWorker(t *testing.T, server *testWorkerServer) {
	t.Helper()
	request := testWorkerInitializationRequest()
	_, err := server.client.InitializeWorker(context.Background(), request)
	require.NoError(t, err)
}

func TestWorkerRegistrationAndListBeforeInitialization(t *testing.T) {
	t.Setenv("TILEBOX_API_URL", "https://network-must-not-be-used.invalid")
	t.Setenv("TILEBOX_API_KEY", "construction-must-not-read-this-credential")

	worker := NewWorker()
	require.NoError(t, worker.RegisterTasks(&workerRegistrationTask{}, &workerJSONTask{}))
	require.ErrorContains(t, worker.RegisterTasks(&workerRegistrationTask{}), "duplicate task identifier")
	require.ErrorContains(t, worker.RegisterTasks(&workerInvalidIdentifierTask{}), "task name is empty")
	require.ErrorContains(t, worker.RegisterTasks(&workerNilIdentifierTask{}), "task identifier is nil")
	require.ErrorContains(t, worker.RegisterTasks(workerValueTask{}), "must be a non-nil pointer to a struct")
	var nilTask *workerRegistrationTask
	require.ErrorContains(t, worker.RegisterTasks(nilTask), "must be a non-nil pointer to a struct")

	server := startTestWorker(t, worker)
	response, err := server.client.ListRegisteredTasks(context.Background(), &emptypb.Empty{})
	require.NoError(t, err)
	require.Len(t, response.GetIdentifiers(), 2)
	assert.Equal(t, "workerRegistrationTask", response.GetIdentifiers()[0].GetName())
	assert.Equal(t, "v0.0", response.GetIdentifiers()[0].GetVersion())
	assert.Equal(t, "tilebox.com/tests/json", response.GetIdentifiers()[1].GetName())
	require.ErrorContains(t, worker.RegisterTasks(&workerPanicTask{}), "after worker serving has started")
}

func TestWorkerConnectGRPCInteroperability(t *testing.T) {
	workerJSONValues = make(chan string, 1)
	worker := NewWorker()
	require.NoError(t, worker.RegisterTasks(&workerJSONTask{}))
	server := startTestWorker(t, worker)

	transport := &http2.Transport{
		AllowHTTP: true,
		DialTLSContext: func(ctx context.Context, _, _ string, _ *tls.Config) (net.Conn, error) {
			return (&net.Dialer{}).DialContext(ctx, "unix", server.socketPath)
		},
	}
	t.Cleanup(transport.CloseIdleConnections)
	client := workflowsv1connect.NewWorkerServiceClient(
		&http.Client{Transport: transport},
		"http://tilebox-worker",
		connect.WithGRPC(),
	)

	registered, err := client.ListRegisteredTasks(context.Background(), connect.NewRequest(&emptypb.Empty{}))
	require.NoError(t, err)
	require.Len(t, registered.Msg.GetIdentifiers(), 1)
	assert.Equal(t, (&workerJSONTask{}).Identifier().Name(), registered.Msg.GetIdentifiers()[0].GetName())

	_, err = client.InitializeWorker(context.Background(), connect.NewRequest(testWorkerInitializationRequest()))
	require.NoError(t, err)
	input, err := json.Marshal(&workerJSONTask{Value: "connect gRPC decoded"})
	require.NoError(t, err)
	executed, err := client.ExecuteTask(context.Background(), connect.NewRequest(testWorkerTask((&workerJSONTask{}).Identifier(), input)))
	require.NoError(t, err)
	require.NotNil(t, executed.Msg.GetComputedTask())
	assert.Equal(t, "connect gRPC decoded", <-workerJSONValues)

	_, err = client.ShutdownWorker(context.Background(), connect.NewRequest(&emptypb.Empty{}))
	require.NoError(t, err)
	require.NoError(t, server.wait(t))
}

func TestWorkerConcurrentRegistrationAndServeLifecycle(t *testing.T) {
	worker := NewWorker()
	socketPath := filepath.Join(t.TempDir(), "worker.sock")
	t.Setenv(workerAddressEnvironmentVariable, "unix://"+socketPath)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	start := make(chan struct{})
	errorsCh := make(chan error, 33)
	for range 32 {
		go func() {
			<-start
			errorsCh <- worker.RegisterTasks()
		}()
	}
	serveDone := make(chan error, 1)
	go func() {
		<-start
		serveDone <- worker.Serve(ctx)
	}()
	close(start)

	deadline := time.Now().Add(3 * time.Second)
	for {
		if _, err := os.Lstat(socketPath); err == nil {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("worker did not bind socket")
		}
		time.Sleep(time.Millisecond)
	}

	require.ErrorContains(t, worker.Serve(context.Background()), "only be served once")
	for range 32 {
		err := <-errorsCh
		if err != nil {
			require.ErrorContains(t, err, "after worker serving has started")
		}
	}
	cancel()
	require.NoError(t, <-serveDone)
	_, err := os.Lstat(socketPath)
	require.ErrorIs(t, err, os.ErrNotExist)
	require.ErrorContains(t, worker.RegisterTasks(), "after worker serving has started")
}

func TestWorkerProtocolInitializationAndExecution(t *testing.T) {
	workerJSONValues = make(chan string, 1)
	workerProtoValues = make(chan int64, 1)
	workerTaskContextResults = make(chan workerTaskContextResult, 1)

	worker := NewWorker()
	require.NoError(t, worker.RegisterTasks(
		&workerJSONTask{},
		&workerProtoTask{},
		&workerDispatchV1Task{},
		&workerDispatchV2Task{},
		&workerContextTask{},
	))
	server := startTestWorker(t, worker)

	preInitializeTask := testWorkerTask((&workerJSONTask{}).Identifier(), []byte(`{"value":"before"}`))
	response, err := server.client.ExecuteTask(context.Background(), preInitializeTask)
	require.NoError(t, err)
	require.NotNil(t, response.GetFailedTask())
	assert.False(t, response.GetFailedTask().GetWasWorkflowError())
	assert.Contains(t, response.GetFailedTask().GetDisplay(), "worker is not initialized")

	initialization := testWorkerInitializationRequest()
	const concurrentInitializations = 12
	errorsCh := make(chan error, concurrentInitializations)
	for range concurrentInitializations {
		go func() {
			_, err := server.client.InitializeWorker(context.Background(), proto.CloneOf(initialization))
			errorsCh <- err
		}()
	}
	for range concurrentInitializations {
		require.NoError(t, <-errorsCh)
	}
	retriedTraceParent := proto.CloneOf(initialization)
	retriedTraceParent.SetTraceParent("00-11111111111111111111111111111111-2222222222222222-01")
	_, err = server.client.InitializeWorker(context.Background(), retriedTraceParent)
	require.NoError(t, err)

	conflicting := proto.CloneOf(initialization)
	conflicting.SetCluster(workflowsv1.Cluster_builder{Slug: "other-cluster"}.Build())
	_, err = server.client.InitializeWorker(context.Background(), conflicting)
	require.Error(t, err)
	assert.Equal(t, codes.FailedPrecondition, status.Code(err))

	jsonInput, err := json.Marshal(&workerJSONTask{Value: "json decoded"})
	require.NoError(t, err)
	response, err = server.client.ExecuteTask(context.Background(), testWorkerTask((&workerJSONTask{}).Identifier(), jsonInput))
	require.NoError(t, err)
	require.NotNil(t, response.GetComputedTask())
	assert.Equal(t, "json decoded", <-workerJSONValues)

	protoInput, err := proto.Marshal(&workerProtoTask{SpawnWorkflowTreeTask: *examplesv1.SpawnWorkflowTreeTask_builder{
		Depth: proto.Int64(7),
	}.Build()})
	require.NoError(t, err)
	response, err = server.client.ExecuteTask(context.Background(), testWorkerTask((&workerProtoTask{}).Identifier(), protoInput))
	require.NoError(t, err)
	require.NotNil(t, response.GetComputedTask())
	assert.Equal(t, int64(7), <-workerProtoValues)

	for _, identifier := range []TaskIdentifier{(&workerDispatchV1Task{}).Identifier(), (&workerDispatchV2Task{}).Identifier()} {
		response, err = server.client.ExecuteTask(context.Background(), testWorkerTask(identifier, []byte(`{}`)))
		require.NoError(t, err)
		require.NotNil(t, response.GetComputedTask())
	}

	response, err = server.client.ExecuteTask(context.Background(), testWorkerTask((&workerContextTask{}).Identifier(), []byte(`{}`)))
	require.NoError(t, err)
	require.NotNil(t, response.GetComputedTask())
	contextResult := <-workerTaskContextResults
	require.NoError(t, contextResult.err)
	assert.True(t, contextResult.hasClient)
	assert.Equal(t, "test-cluster", contextResult.cluster)
	assert.Equal(t, "context ready", response.GetComputedTask().GetDisplay())
	require.Len(t, response.GetComputedTask().GetProgressUpdates(), 1)
	assert.Equal(t, "work", response.GetComputedTask().GetProgressUpdates()[0].GetLabel())
	assert.Equal(t, uint64(3), response.GetComputedTask().GetProgressUpdates()[0].GetTotal())
	assert.Equal(t, uint64(2), response.GetComputedTask().GetProgressUpdates()[0].GetDone())
	require.NotNil(t, response.GetComputedTask().GetSubTasks())
	require.Len(t, response.GetComputedTask().GetSubTasks().GetTaskGroups(), 1)
}

func TestWorkerExecutionFailureClassification(t *testing.T) {
	worker := NewWorker()
	require.NoError(t, worker.RegisterTasks(&workerJSONTask{}, &workerUserErrorTask{}, &workerPanicTask{}))
	server := startTestWorker(t, worker)
	initializeTestWorker(t, server)

	tests := []struct {
		name             string
		task             *workflowsv1.Task
		wasWorkflowError bool
		displayContains  string
		progressUpdates  int
	}{
		{
			name:             "unknown exact version",
			task:             testWorkerTask(NewTaskIdentifier("tilebox.com/tests/json", "v9.0"), []byte(`{}`)),
			wasWorkflowError: false,
			displayContains:  "is not registered",
		},
		{
			name:             "JSON decode",
			task:             testWorkerTask((&workerJSONTask{}).Identifier(), []byte(`{"value":`)),
			wasWorkflowError: true,
			displayContains:  "failed to unmarshal json task",
		},
		{
			name:             "user error",
			task:             testWorkerTask((&workerUserErrorTask{}).Identifier(), []byte(`{}`)),
			wasWorkflowError: true,
			displayContains:  "user task failed",
			progressUpdates:  1,
		},
		{
			name:             "panic recovery",
			task:             testWorkerTask((&workerPanicTask{}).Identifier(), []byte(`{}`)),
			wasWorkflowError: true,
			displayContains:  "task panicked: user panic",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response, err := server.client.ExecuteTask(context.Background(), test.task)
			require.NoError(t, err)
			failedTask := response.GetFailedTask()
			require.NotNil(t, failedTask)
			assert.Equal(t, test.wasWorkflowError, failedTask.GetWasWorkflowError())
			assert.Contains(t, failedTask.GetDisplay(), test.displayContains)
			assert.Len(t, failedTask.GetProgressUpdates(), test.progressUpdates)
		})
	}
}

func TestWorkerExecuteCancellationReachesTaskContext(t *testing.T) {
	workerCancellationStarted = make(chan struct{})
	workerCancellationObserved = make(chan error, 1)
	worker := NewWorker()
	require.NoError(t, worker.RegisterTasks(&workerCancellationTask{}))
	server := startTestWorker(t, worker)
	initializeTestWorker(t, server)

	ctx, cancel := context.WithCancel(context.Background())
	executeDone := make(chan error, 1)
	go func() {
		_, err := server.client.ExecuteTask(ctx, testWorkerTask((&workerCancellationTask{}).Identifier(), []byte(`{}`)))
		executeDone <- err
	}()
	select {
	case <-workerCancellationStarted:
	case <-time.After(time.Second):
		t.Fatal("task did not start")
	}
	cancel()

	err := <-executeDone
	require.Error(t, err)
	assert.Equal(t, codes.Canceled, status.Code(err))
	select {
	case observed := <-workerCancellationObserved:
		require.ErrorIs(t, observed, context.Canceled)
	case <-time.After(time.Second):
		t.Fatal("task did not observe RPC cancellation")
	}
}

func TestWorkerShutdownIdleCleansUp(t *testing.T) {
	worker := NewWorker()
	server := startTestWorker(t, worker)
	initializeTestWorker(t, server)

	var cleanupCalls atomic.Int64
	worker.mu.Lock()
	originalCleanup := worker.runtime.cleanup
	worker.runtime.cleanup = func(ctx context.Context) {
		cleanupCalls.Add(1)
		originalCleanup(ctx)
	}
	worker.mu.Unlock()

	_, err := server.client.ShutdownWorker(context.Background(), &emptypb.Empty{})
	require.NoError(t, err)
	require.NoError(t, server.wait(t))
	assert.Equal(t, int64(1), cleanupCalls.Load())
	_, err = os.Lstat(server.socketPath)
	assert.ErrorIs(t, err, os.ErrNotExist)
}

func TestWorkerShutdownActiveIsBoundedAndCancelsTask(t *testing.T) {
	workerCancellationStarted = make(chan struct{})
	workerCancellationObserved = make(chan error, 1)
	worker := NewWorker()
	worker.shutdownGracePeriod = 25 * time.Millisecond
	require.NoError(t, worker.RegisterTasks(&workerCancellationTask{}))
	server := startTestWorker(t, worker)
	initializeTestWorker(t, server)

	executeDone := make(chan error, 1)
	go func() {
		_, err := server.client.ExecuteTask(context.Background(), testWorkerTask((&workerCancellationTask{}).Identifier(), []byte(`{}`)))
		executeDone <- err
	}()
	select {
	case <-workerCancellationStarted:
	case <-time.After(time.Second):
		t.Fatal("task did not start")
	}

	startedShutdown := time.Now()
	_, err := server.client.ShutdownWorker(context.Background(), &emptypb.Empty{})
	require.NoError(t, err)
	require.NoError(t, server.wait(t))
	assert.Less(t, time.Since(startedShutdown), time.Second)
	select {
	case observed := <-workerCancellationObserved:
		require.ErrorIs(t, observed, context.Canceled)
	case <-time.After(time.Second):
		t.Fatal("active task did not observe shutdown cancellation")
	}
	require.Error(t, <-executeDone)
	_, err = os.Lstat(server.socketPath)
	assert.ErrorIs(t, err, os.ErrNotExist)
}

func TestWorkerShutdownGracefullyDeliversActiveResult(t *testing.T) {
	workerGracefulStarted = make(chan struct{})
	workerGracefulRelease = make(chan struct{})
	worker := NewWorker()
	worker.shutdownGracePeriod = time.Second
	require.NoError(t, worker.RegisterTasks(&workerGracefulTask{}))
	server := startTestWorker(t, worker)
	initializeTestWorker(t, server)

	type executeResult struct {
		response *workflowsv1.ExecuteTaskResponse
		err      error
	}
	executeDone := make(chan executeResult, 1)
	go func() {
		response, err := server.client.ExecuteTask(context.Background(), testWorkerTask((&workerGracefulTask{}).Identifier(), []byte(`{}`)))
		executeDone <- executeResult{response: response, err: err}
	}()
	select {
	case <-workerGracefulStarted:
	case <-time.After(time.Second):
		t.Fatal("task did not start")
	}

	_, err := server.client.ShutdownWorker(context.Background(), &emptypb.Empty{})
	require.NoError(t, err)
	close(workerGracefulRelease)
	result := <-executeDone
	require.NoError(t, result.err)
	require.NotNil(t, result.response.GetComputedTask())
	require.NoError(t, server.wait(t))
}

func TestWorkerShutdownRejectsQueuedExecution(t *testing.T) {
	workerQueuedStarted = make(chan struct{})
	workerQueuedCanceled = make(chan struct{})
	workerQueuedSecondStarted = make(chan struct{})
	workerQueuedSecondRan.Store(false)
	worker := NewWorker()
	worker.shutdownGracePeriod = 25 * time.Millisecond
	require.NoError(t, worker.RegisterTasks(&workerQueuedFirstTask{}, &workerQueuedSecondTask{}))
	server := startTestWorker(t, worker)
	initializeTestWorker(t, server)

	firstDone := make(chan error, 1)
	go func() {
		_, err := server.client.ExecuteTask(context.Background(), testWorkerTask((&workerQueuedFirstTask{}).Identifier(), []byte(`{}`)))
		firstDone <- err
	}()
	select {
	case <-workerQueuedStarted:
	case <-time.After(time.Second):
		t.Fatal("first task did not start")
	}

	type executeResult struct {
		response *workflowsv1.ExecuteTaskResponse
		err      error
	}
	secondDone := make(chan executeResult, 1)
	go func() {
		close(workerQueuedSecondStarted)
		response, err := server.client.ExecuteTask(context.Background(), testWorkerTask((&workerQueuedSecondTask{}).Identifier(), []byte(`{}`)))
		secondDone <- executeResult{response: response, err: err}
	}()
	<-workerQueuedSecondStarted
	time.Sleep(5 * time.Millisecond)

	_, err := server.client.ShutdownWorker(context.Background(), &emptypb.Empty{})
	require.NoError(t, err)
	second := <-secondDone
	require.NoError(t, second.err)
	require.NotNil(t, second.response.GetFailedTask())
	assert.False(t, second.response.GetFailedTask().GetWasWorkflowError())
	assert.Contains(t, second.response.GetFailedTask().GetDisplay(), "worker is shutting down")
	assert.False(t, workerQueuedSecondRan.Load())
	select {
	case <-workerQueuedCanceled:
	case <-time.After(time.Second):
		t.Fatal("first task did not observe shutdown cancellation")
	}
	require.Error(t, <-firstDone)
	require.NoError(t, server.wait(t))
}

func TestWorkerInitializationCanRetryAfterInvalidConfiguration(t *testing.T) {
	worker := NewWorker()
	server := startTestWorker(t, worker)
	invalid := testWorkerInitializationRequest()
	invalidURL := "http://"
	token := "not-exported"
	invalid.SetApiConnection(workflowsv1.TileboxAPIConnection_builder{Url: &invalidURL, Token: &token}.Build())

	_, err := server.client.InitializeWorker(context.Background(), invalid)
	require.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
	assert.NotContains(t, err.Error(), token)

	valid := testWorkerInitializationRequest()
	valid.SetRunnerId(invalid.GetRunnerId())
	_, err = server.client.InitializeWorker(context.Background(), valid)
	require.NoError(t, err)
}

func TestWorkerParentCancellationAndListenFailures(t *testing.T) {
	t.Run("parent cancellation", func(t *testing.T) {
		worker := NewWorker()
		server := startTestWorker(t, worker)
		server.cancel()
		require.NoError(t, server.wait(t))
		_, err := os.Lstat(server.socketPath)
		assert.ErrorIs(t, err, os.ErrNotExist)
	})

	t.Run("missing address", func(t *testing.T) {
		t.Setenv(workerAddressEnvironmentVariable, "")
		worker := NewWorker()
		err := worker.Serve(context.Background())
		require.ErrorContains(t, err, workerAddressEnvironmentVariable+" is not set")
		require.ErrorContains(t, worker.RegisterTasks(), "after worker serving has started")
	})

	t.Run("existing regular file is preserved", func(t *testing.T) {
		socketPath := filepath.Join(t.TempDir(), "worker.sock")
		require.NoError(t, os.WriteFile(socketPath, []byte("not a socket"), 0o600))
		t.Setenv(workerAddressEnvironmentVariable, "unix://"+socketPath)
		err := NewWorker().Serve(context.Background())
		require.ErrorContains(t, err, "is not a socket")
		contents, readErr := os.ReadFile(socketPath)
		require.NoError(t, readErr)
		assert.Equal(t, "not a socket", string(contents))
	})

	t.Run("stale socket is safely replaced and removed", func(t *testing.T) {
		socketPath := filepath.Join(t.TempDir(), "worker.sock")
		listener, err := (&net.ListenConfig{}).Listen(context.Background(), "unix", socketPath)
		require.NoError(t, err)
		require.NoError(t, listener.Close())

		worker := NewWorker()
		server := startTestWorkerAt(t, worker, socketPath)
		server.cancel()
		require.NoError(t, server.wait(t))
		_, err = os.Lstat(socketPath)
		assert.ErrorIs(t, err, os.ErrNotExist)
	})
}

func TestWorkerCredentialsAreRedacted(t *testing.T) {
	var logs bytes.Buffer
	previousLogger := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&logs, nil)))
	t.Cleanup(func() { slog.SetDefault(previousLogger) })

	requestToken := "worker-request-secret-token"
	environmentToken := "worker-environment-secret-token"
	t.Setenv("TILEBOX_API_KEY", environmentToken)
	workerSensitiveValue = requestToken + " " + environmentToken
	worker := NewWorker()
	require.NoError(t, worker.RegisterTasks(&workerSensitiveTask{}))
	server := startTestWorker(t, worker)
	request := testWorkerInitializationRequest()
	apiURL := filepath.Join(t.TempDir(), "unused-api.sock")
	request.SetApiConnection(workflowsv1.TileboxAPIConnection_builder{
		Url:   &apiURL,
		Token: &requestToken,
	}.Build())
	_, err := server.client.InitializeWorker(context.Background(), request)
	require.NoError(t, err)

	task := testWorkerTask((&workerSensitiveTask{}).Identifier(), []byte(`{}`))
	task.SetDisplay("input display " + workerSensitiveValue)
	response, err := server.client.ExecuteTask(context.Background(), task)
	require.NoError(t, err)
	require.NotNil(t, response.GetFailedTask())
	assert.NotContains(t, response.GetFailedTask().GetDisplay(), requestToken)
	assert.NotContains(t, response.GetFailedTask().GetDisplay(), environmentToken)
	assert.Contains(t, response.GetFailedTask().GetDisplay(), "[REDACTED]")
	assert.NotContains(t, logs.String(), requestToken)
	assert.NotContains(t, logs.String(), environmentToken)
}

func TestClientConfigDefaultsURLFromEnvironment(t *testing.T) {
	t.Setenv("TILEBOX_API_URL", "https://workflows.example.test")
	t.Setenv("TILEBOX_API_KEY", "")
	cfg := newClientConfig([]ClientOption{WithHTTPClient(&http.Client{})})
	assert.Equal(t, "https://workflows.example.test", cfg.url)

	cfg = newClientConfig([]ClientOption{
		WithHTTPClient(&http.Client{}),
		WithURL("https://explicit.example.test"),
	})
	assert.Equal(t, "https://explicit.example.test", cfg.url)
}

func TestTaskRunnerKeepsNativeFailureClassification(t *testing.T) {
	runner, err := newTaskRunner(
		context.Background(),
		mockTaskService{},
		clusterClient{service: &mockClusterService{}},
		noop.NewTracerProvider().Tracer("test"),
	)
	require.NoError(t, err)
	response, err := runner.ExecuteTask(context.Background(), testWorkerTask(NewTaskIdentifier("unknown", "v1.0"), []byte(`{}`)))
	require.NoError(t, err)
	require.NotNil(t, response.GetFailedTask())
	assert.True(t, response.GetFailedTask().GetWasWorkflowError())
}
