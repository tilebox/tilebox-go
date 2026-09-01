package workflows

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"os/signal"
	"reflect"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/tilebox/tilebox-go/internal/span"
	obslogger "github.com/tilebox/tilebox-go/observability/logger"
	tileboxv1 "github.com/tilebox/tilebox-go/protogen/tilebox/v1"
	workflowsv1 "github.com/tilebox/tilebox-go/protogen/workflows/v1"
	"github.com/tilebox/tilebox-go/workflows/v1/runner"
	"github.com/tilebox/tilebox-go/workflows/v1/subtask"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/proto"
)

type contextKeyTaskExecutionType string

const contextKeyTaskExecution contextKeyTaskExecutionType = "x-tilebox-task-execution-object"

const (
	// A maximum idling duration, as a safeguard to avoid way too long sleep times in case the suggested idling duration is
	// ever too long. 5 minutes should be plenty of time to wait.
	maxIdlingDuration = 5 * time.Minute
	// A minimum idling duration, as a safeguard to avoid too short sleep times in case the suggested idling duration is ever too short.
	minIdlingDuration = 1 * time.Millisecond

	// Fallback polling interval and jitter in case the workflows API fails to respond with a suggested idling duration
	fallbackPollingInterval = 10 * time.Second
	fallbackJitterInterval  = 5 * time.Second
)

const (
	UnitSeconds       = "s"
	UnitDimensionless = "1"
	UnitBytes         = "By"
)

type taskRunnerMetrics struct {
	tasksExecutedMetric metric.Int64Counter
	tasksComputedMetric metric.Int64Counter
	tasksFailedMetric   metric.Int64Counter

	taskInputSizeMetric         metric.Int64Histogram
	taskExecutionDurationMetric metric.Float64Histogram
}

func newTaskRunnerMetrics(meter metric.Meter) (*taskRunnerMetrics, error) {
	tasksExecutedMetric, err := meter.Int64Counter(
		"task.executed.count",
		metric.WithDescription("Number of tasks executed"),
		metric.WithUnit(UnitDimensionless),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create task count metric: %w", err)
	}

	tasksComputedMetric, err := meter.Int64Counter(
		"task.computed.count",
		metric.WithDescription("Number of tasks computed"),
		metric.WithUnit(UnitDimensionless),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create task computed metric: %w", err)
	}

	tasksFailedMetric, err := meter.Int64Counter(
		"task.failed.count",
		metric.WithDescription("Number of tasks failed"),
		metric.WithUnit(UnitDimensionless),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create task failed metric: %w", err)
	}

	taskArgsSizeMetric, err := meter.Int64Histogram(
		"task.input.size",
		metric.WithDescription("Task arguments size"),
		metric.WithUnit(UnitBytes),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create task input size metric: %w", err)
	}

	taskExecutionDurationMetric, err := meter.Float64Histogram(
		"task.execution.duration",
		metric.WithDescription("Task execution duration"),
		metric.WithUnit(UnitSeconds),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create task duration metric: %w", err)
	}

	return &taskRunnerMetrics{
		tasksExecutedMetric:         tasksExecutedMetric,
		tasksComputedMetric:         tasksComputedMetric,
		tasksFailedMetric:           tasksFailedMetric,
		taskInputSizeMetric:         taskArgsSizeMetric,
		taskExecutionDurationMetric: taskExecutionDurationMetric,
	}, nil
}

// TaskRunner executes tasks.
//
// Documentation: https://docs.tilebox.com/workflows/concepts/task-runners
type TaskRunner struct {
	pollingRunner   *PollingTaskRunner
	service         TaskService
	taskDefinitions map[taskIdentifier]ExecutableTask

	cluster string
	tracer  trace.Tracer
	logger  *slog.Logger
	metrics *taskRunnerMetrics
}

func newTaskRunner(ctx context.Context, service TaskService, clusterClient ClusterClient, tracer trace.Tracer, options ...runner.Option) (*TaskRunner, error) {
	opts := &runner.Options{
		ClusterSlug:   "",
		Logger:        slog.Default(),
		MeterProvider: otel.GetMeterProvider(),
	}
	for _, option := range options {
		option(opts)
	}

	cluster, err := clusterClient.Get(ctx, opts.ClusterSlug)
	if err != nil {
		return nil, fmt.Errorf("failed to get cluster: %w", err)
	}

	metrics, err := newTaskRunnerMetrics(opts.MeterProvider.Meter(otelMeterName))
	if err != nil {
		return nil, fmt.Errorf("failed to create task runner metrics: %w", err)
	}

	taskRunner := &TaskRunner{
		service:         service,
		taskDefinitions: make(map[taskIdentifier]ExecutableTask),

		cluster: cluster.Slug,
		tracer:  tracer,
		logger:  opts.Logger,
		metrics: metrics,
	}
	taskRunner.pollingRunner = NewPollingTaskRunner(service, cluster.Slug, taskRunner, opts.Logger)
	return taskRunner, nil
}

// GetRegisteredTask returns the task with the given identifier.
func (t *TaskRunner) GetRegisteredTask(identifier TaskIdentifier) (ExecutableTask, bool) {
	registeredTask, found := t.taskDefinitions[taskIdentifier{name: identifier.Name(), version: identifier.Version()}]
	return registeredTask, found
}

// RegisterTasks makes the task runner aware of multiple tasks.
func (t *TaskRunner) RegisterTasks(tasks ...ExecutableTask) error {
	for _, task := range tasks {
		err := t.registerTask(task)
		if err != nil {
			return err
		}
	}
	return nil
}

func isEmpty(id *tileboxv1.ID) bool {
	return id.AsUUID() == uuid.Nil
}

// RunForever runs the task runner forever, looking for new tasks to run and polling for new tasks when idle.
func (t *TaskRunner) RunForever(ctx context.Context) error {
	ctxSignal, stop := signal.NotifyContext(ctx, runnerShutdownSignals()...)
	defer stop()
	t.pollingRunner.service = t.service
	return t.pollingRunner.RunForever(ctxSignal)
}

// RunAll run the task runner and execute all tasks until there are no more tasks available.
func (t *TaskRunner) RunAll(ctx context.Context) error {
	t.pollingRunner.service = t.service
	return t.pollingRunner.RunAll(ctx)
}

func (t *TaskRunner) TaskIdentifiers() []*workflowsv1.TaskIdentifier {
	identifiers := make([]*workflowsv1.TaskIdentifier, 0, len(t.taskDefinitions))
	for _, task := range t.taskDefinitions {
		identifier := identifierFromTask(task)
		identifiers = append(identifiers, workflowsv1.TaskIdentifier_builder{
			Name:    identifier.Name(),
			Version: identifier.Version(),
		}.Build())
	}
	return identifiers
}

func (t *TaskRunner) ExecuteTask(ctx context.Context, task *workflowsv1.Task) (*workflowsv1.ExecuteTaskResponse, error) {
	executionContext, err := t.executeTask(ctx, task)
	if err != nil {
		var progressUpdates []*workflowsv1.Progress
		if executionContext != nil {
			progressUpdates = executionContext.getProgressUpdates()
		}
		failedTask := workflowsv1.TaskFailedRequest_builder{
			TaskId:           task.GetId(),
			Display:          failedTaskDisplay(task.GetDisplay(), err),
			WasWorkflowError: true,
			ProgressUpdates:  progressUpdates,
		}.Build()
		return workflowsv1.ExecuteTaskResponse_builder{FailedTask: failedTask}.Build(), nil
	}

	computedTask := workflowsv1.ComputedTask_builder{
		Id:              task.GetId(),
		Display:         task.GetDisplay(),
		SubTasks:        executionContext.getSubTasks(),
		ProgressUpdates: executionContext.getProgressUpdates(),
	}.Build()
	return workflowsv1.ExecuteTaskResponse_builder{ComputedTask: computedTask}.Build(), nil
}

func (t *TaskRunner) executeTask(ctx context.Context, task *workflowsv1.Task) (*taskExecutionContext, error) {
	// actually execute the task
	beforeTime := time.Now().UTC()

	if task.GetIdentifier() == nil {
		return nil, errors.New("task has no identifier")
	}
	identifier := NewTaskIdentifier(task.GetIdentifier().GetName(), task.GetIdentifier().GetVersion())
	taskPrototype, found := t.GetRegisteredTask(identifier)
	if !found {
		return nil, fmt.Errorf("task %s is not registered on this runner", task.GetIdentifier().GetName())
	}

	return span.StartJobSpan(ctx, t.tracer, fmt.Sprintf("task/%s", identifier.Name()), task.GetJob(), func(ctx context.Context) (taskExecutionContext *taskExecutionContext, err error) { //nolint:nonamedreturns // needed to return a value in case of panic
		t.logger.DebugContext(ctx, "executing task",
			slog.String("task_id", task.GetId().AsUUID().String()),
			slog.String("task", identifier.Name()),
			slog.String("version", identifier.Version()),
		)
		taskStruct := reflect.New(reflect.ValueOf(taskPrototype).Elem().Type()).Interface().(ExecutableTask)

		_, isProtobuf := taskStruct.(proto.Message)
		if isProtobuf {
			err := proto.Unmarshal(task.GetInput(), taskStruct.(proto.Message))
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal protobuf task: %w", err)
			}
		} else {
			err := json.Unmarshal(task.GetInput(), taskStruct)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal json task: %w", err)
			}
		}

		log := t.logger.With(
			slog.String("job_id", task.GetJob().GetId().AsUUID().String()),
			slog.String("task", identifier.Name()),
			slog.String("version", identifier.Version()),
			slog.Time("start_time", beforeTime),
		)

		taskMetricAttributes := metric.WithAttributes(attribute.String("task_identifier", identifier.Name()), attribute.String("task_version", identifier.Version()))

		defer func() {
			if r := recover(); r != nil {
				// recover from panics during task executions, so we can still report the error to the server and continue
				// with other tasks
				log.ErrorContext(ctx, "task execution failed", slog.String("error", "panic"), slog.Int64("retry_attempt", task.GetRetryCount()))
				taskExecutionContext = nil
				err = fmt.Errorf("task panicked: %v", r)

				// also instrument the panic as a failed task execution
				t.metrics.tasksFailedMetric.Add(ctx, 1, taskMetricAttributes)
				t.metrics.taskExecutionDurationMetric.Record(ctx, time.Since(beforeTime).Seconds(), taskMetricAttributes, metric.WithAttributes(attribute.String("state", "failed")))
			}
		}()

		t.metrics.taskInputSizeMetric.Record(ctx, int64(len(task.GetInput())), taskMetricAttributes)
		t.metrics.tasksExecutedMetric.Add(ctx, 1, taskMetricAttributes)

		executionContext := t.withTaskExecutionContext(ctx, task)
		executionContext = obslogger.ContextWithSlogAttributes(executionContext, slog.String("task_id", task.GetId().AsUUID().String()))
		err = taskStruct.Execute(executionContext)

		executionTime := time.Since(beforeTime)
		log = log.With(
			slog.Duration("execution_time", executionTime),
			slog.String("execution_time_human", roundDuration(executionTime, 2).String()),
		)

		if err != nil {
			t.metrics.tasksFailedMetric.Add(ctx, 1, taskMetricAttributes)
			t.metrics.taskExecutionDurationMetric.Record(ctx, executionTime.Seconds(), taskMetricAttributes, metric.WithAttributes(attribute.String("state", "failed")))

			log.ErrorContext(executionContext, "task execution failed", slog.Any("error", err), slog.Int64("retry_attempt", task.GetRetryCount()))
			return getTaskExecutionContext(executionContext), fmt.Errorf("failed to execute task: %w", err)
		}

		t.metrics.tasksComputedMetric.Add(ctx, 1, taskMetricAttributes)
		t.metrics.taskExecutionDurationMetric.Record(ctx, executionTime.Seconds(), taskMetricAttributes, metric.WithAttributes(attribute.String("state", "computed")))

		return getTaskExecutionContext(executionContext), nil
	})
}

// extendTaskLease is a function designed to be run as a goroutine, extending the lease of a task continuously until the
// context is cancelled, which indicates that the execution of the task is finished.
// registerTask makes the task runner aware of a task.
func (t *TaskRunner) registerTask(task ExecutableTask) error {
	identifier := identifierFromTask(task)
	err := ValidateIdentifier(identifier)
	if err != nil {
		return err
	}
	key := taskIdentifier{name: identifier.Name(), version: identifier.Version()}
	_, found := t.taskDefinitions[key]
	if found {
		return fmt.Errorf(
			"duplicate task identifier: a task '%s' with version '%s' is already registered",
			identifier.Name(),
			identifier.Version(),
		)
	}
	t.taskDefinitions[key] = task
	return nil
}

type taskExecutionContext struct {
	CurrentTask *workflowsv1.Task
	runner      *TaskRunner

	subtasksMutex sync.Mutex
	subtasks      []*futureTask

	progressMutex      sync.Mutex
	progressIndicators map[string]*taskProgressIndicator
}

// getProgressUpdates converts the internal progress indicators into the protobuf representation we need for
// reporting progress to the workflows API
func (e *taskExecutionContext) getProgressUpdates() []*workflowsv1.Progress {
	e.progressMutex.Lock()
	defer e.progressMutex.Unlock()
	progressUpdates := make([]*workflowsv1.Progress, 0, len(e.progressIndicators))
	for _, progress := range e.progressIndicators {
		progressUpdates = append(progressUpdates, workflowsv1.Progress_builder{
			Label: progress.label,
			Total: progress.total,
			Done:  progress.done,
		}.Build())
	}
	return progressUpdates
}

// getSubTasks converts the internal subtask submissions into the protobuf TaskSubmissions format.
func (e *taskExecutionContext) getSubTasks() *workflowsv1.TaskSubmissions {
	e.subtasksMutex.Lock()
	defer e.subtasksMutex.Unlock()
	return mergeFutureTasksToSubmissions(e.subtasks)
}

func (t *TaskRunner) withTaskExecutionContext(ctx context.Context, task *workflowsv1.Task) context.Context {
	return context.WithValue(ctx, contextKeyTaskExecution, &taskExecutionContext{
		CurrentTask:        task,
		runner:             t,
		subtasksMutex:      sync.Mutex{},
		subtasks:           make([]*futureTask, 0),
		progressIndicators: make(map[string]*taskProgressIndicator),
		progressMutex:      sync.Mutex{},
	})
}

func getTaskExecutionContext(ctx context.Context) *taskExecutionContext {
	executionContext := ctx.Value(contextKeyTaskExecution)
	if executionContext == nil {
		return nil
	}
	return executionContext.(*taskExecutionContext)
}

// GetCurrentCluster returns the current cluster slug.
//
// This function is intended to be used in tasks to get the current cluster slug.
func GetCurrentCluster(ctx context.Context) (string, error) {
	executionContext := getTaskExecutionContext(ctx)
	if executionContext == nil {
		return "", errors.New("cannot get current cluster without task execution context")
	}
	return executionContext.runner.cluster, nil
}

// SetTaskDisplay sets the label name of the current task.
func SetTaskDisplay(ctx context.Context, display string) error {
	executionContext := getTaskExecutionContext(ctx)
	if executionContext == nil {
		return errors.New("cannot set task label name without task execution context")
	}
	executionContext.CurrentTask.SetDisplay(display)
	return nil
}

// SubmitSubtask submits a task to the task runner as a subtask of the current task.
//
// Options:
//   - subtask.WithDependencies: sets the dependencies of the task.
//   - subtask.WithClusterSlug: sets the cluster slug of the cluster where the task will be executed. Defaults to the cluster of the task runner.
//   - subtask.WithMaxRetries: sets the maximum number of times a task can be automatically retried. Defaults to 0.
func SubmitSubtask(ctx context.Context, task Task, options ...subtask.SubmitOption) (subtask.FutureTask, error) {
	executionContext := getTaskExecutionContext(ctx)
	if executionContext == nil {
		return 0, errors.New("cannot submit subtask without task execution context")
	}

	opts := &subtask.SubmitOptions{
		Dependencies: nil,
		ClusterSlug:  executionContext.runner.cluster,
		MaxRetries:   0,
		Optional:     false,
	}
	for _, option := range options {
		option(opts)
	}

	var subtaskInput []byte
	var err error

	if task == nil {
		return 0, errors.New("cannot submit nil task")
	}
	taskProto, isProtobuf := task.(proto.Message)
	if isProtobuf {
		subtaskInput, err = proto.Marshal(taskProto)
		if err != nil {
			return 0, fmt.Errorf("failed to marshal protobuf task: %w", err)
		}
	} else {
		subtaskInput, err = json.Marshal(task)
		if err != nil {
			return 0, fmt.Errorf("failed to marshal task: %w", err)
		}
	}

	identifier := identifierFromTask(task)
	err = ValidateIdentifier(identifier)
	if err != nil {
		return 0, fmt.Errorf("subtask has invalid task identifier: %w", err)
	}

	executionContext.subtasksMutex.Lock()
	defer executionContext.subtasksMutex.Unlock()
	if len(executionContext.subtasks) >= math.MaxUint32 {
		return 0, errors.New("too many subtasks")
	}
	newTaskIndex := uint32(len(executionContext.subtasks)) //nolint:gosec // we checked that we don't overflow

	var dependencies []uint32
	if len(opts.Dependencies) >= 1 {
		dependencies = make([]uint32, 0, len(opts.Dependencies))
		for _, futureTask := range opts.Dependencies {
			if uint32(futureTask) >= newTaskIndex { // the new index is the last task, so no larger indices than that can exist
				return 0, fmt.Errorf("invalid dependency: future task %d doesn't exist", futureTask)
			}
			dependencies = append(dependencies, uint32(futureTask))
		}
	}

	sub := &futureTask{
		clusterSlug:  opts.ClusterSlug,
		identifier:   identifier,
		input:        subtaskInput,
		dependencies: dependencies,
		maxRetries:   opts.MaxRetries,
		optional:     opts.Optional,
	}

	executionContext.subtasks = append(executionContext.subtasks, sub)
	return subtask.FutureTask(newTaskIndex), nil
}

// SubmitSubtasks submits multiple tasks to the task runner as subtask of the current task.
// It is similar to SubmitSubtask, but it takes a slice of tasks instead of a single task.
func SubmitSubtasks(ctx context.Context, tasks []Task, options ...subtask.SubmitOption) ([]subtask.FutureTask, error) {
	futureTasks := make([]subtask.FutureTask, 0, len(tasks))
	for _, task := range tasks {
		futureTask, err := SubmitSubtask(ctx, task, options...)
		if err != nil {
			return futureTasks, err
		}
		futureTasks = append(futureTasks, futureTask)
	}
	return futureTasks, nil
}

// taskProgressIndicator is the internal struct that keeps track of the progress of a single progress indicator within
// a task.
type taskProgressIndicator struct {
	label string
	total uint64
	done  uint64
	mutex sync.Mutex
}

// ProgressTracker is an interface for updating the total work and completed work units for a named progress
// indicator for a job.
type ProgressTracker interface {
	// Add a given amount of total work to be done to the progress indicator.
	Add(context.Context, uint64) error

	// Done marks a given amount of work as done.
	Done(context.Context, uint64) error
}

// progressTracker is an intermediate facade struct which allows us to delay eventual errors, e.g. due to a context
// without a taskExecutionContext to the actual calls to Add and Done, allowing for a more convenient API.
type progressTracker struct {
	label string
}

func (p *progressTracker) Add(ctx context.Context, n uint64) error {
	progress, err := p.getProgressUpdate(ctx)
	if err != nil {
		return err
	}
	progress.mutex.Lock()
	defer progress.mutex.Unlock()
	progress.total += n
	return nil
}

func (p *progressTracker) Done(ctx context.Context, n uint64) error {
	progress, err := p.getProgressUpdate(ctx)
	if err != nil {
		return err
	}
	progress.mutex.Lock()
	defer progress.mutex.Unlock()
	progress.done += n
	return nil
}

func (p *progressTracker) getProgressUpdate(ctx context.Context) (*taskProgressIndicator, error) {
	executionContext := getTaskExecutionContext(ctx)
	if executionContext == nil {
		return nil, errors.New("cannot track progress without a task execution context")
	}

	// avoid concurrent access to the progress indicators in case multiple goroutines are used within a task
	executionContext.progressMutex.Lock()
	defer executionContext.progressMutex.Unlock()
	progress, found := executionContext.progressIndicators[p.label]
	if !found {
		progress = &taskProgressIndicator{label: p.label, mutex: sync.Mutex{}}
		executionContext.progressIndicators[p.label] = progress
	}
	return progress, nil
}

// DefaultProgress returns the default, unnamed progress indicator instance for tracking job progress.
func DefaultProgress() ProgressTracker {
	return &progressTracker{}
}

// Progress returns a named progress indicator instance for tracking job progress.
func Progress(label string) ProgressTracker {
	return &progressTracker{label: label}
}

// WithTaskSpanResult is a helper function that wraps a function with a tracing span.
// It returns the result of the function and an error if any.
func WithTaskSpanResult[Result any](ctx context.Context, name string, f func(ctx context.Context) (Result, error)) (Result, error) {
	return WithSpanResult(ctx, name, f)
}

// WithTaskSpan is a helper function that wraps a function with a tracing span.
func WithTaskSpan(ctx context.Context, name string, f func(ctx context.Context) error) error {
	return WithSpan(ctx, name, f)
}

var divs = []time.Duration{
	time.Duration(1), time.Duration(10), time.Duration(100), time.Duration(1000),
}

// human-readable, rounded duration, taken from
// https://stackoverflow.com/questions/58414820/limiting-significant-digits-in-formatted-durations
func roundDuration(d time.Duration, digits int) time.Duration {
	switch {
	case d > time.Second:
		d = d.Round(time.Second / divs[digits])
	case d > time.Millisecond:
		d = d.Round(time.Millisecond / divs[digits])
	case d > time.Microsecond:
		d = d.Round(time.Microsecond / divs[digits])
	}
	return d
}
