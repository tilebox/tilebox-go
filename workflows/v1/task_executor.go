package workflows

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"sync"
	"time"

	"github.com/tilebox/tilebox-go/internal/span"
	obslogger "github.com/tilebox/tilebox-go/observability/logger"
	workflowsv1 "github.com/tilebox/tilebox-go/protogen/workflows/v1"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/proto"
)

type taskRegistry struct {
	mu          sync.RWMutex
	definitions map[taskIdentifier]ExecutableTask
	order       []taskIdentifier
}

func newTaskRegistry() *taskRegistry {
	return &taskRegistry{definitions: make(map[taskIdentifier]ExecutableTask)}
}

func (r *taskRegistry) registerTasks(tasks ...ExecutableTask) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, task := range tasks {
		if err := validateExecutableTask(task); err != nil {
			return err
		}

		identifier := identifierFromTask(task)
		if err := ValidateIdentifier(identifier); err != nil {
			return err
		}

		key := taskIdentifier{name: identifier.Name(), version: identifier.Version()}
		if _, found := r.definitions[key]; found {
			return fmt.Errorf(
				"duplicate task identifier: a task '%s' with version '%s' is already registered",
				identifier.Name(),
				identifier.Version(),
			)
		}
		r.definitions[key] = task
		r.order = append(r.order, key)
	}
	return nil
}

func validateExecutableTask(task ExecutableTask) error {
	if task == nil {
		return errors.New("cannot register nil task")
	}
	taskType := reflect.TypeOf(task)
	taskValue := reflect.ValueOf(task)
	if taskType.Kind() != reflect.Pointer || taskType.Elem().Kind() != reflect.Struct || taskValue.IsNil() {
		return fmt.Errorf("task %T must be a non-nil pointer to a struct", task)
	}
	return nil
}

func (r *taskRegistry) get(identifier TaskIdentifier) (ExecutableTask, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	registeredTask, found := r.definitions[taskIdentifier{name: identifier.Name(), version: identifier.Version()}]
	return registeredTask, found
}

func (r *taskRegistry) identifiers() []*workflowsv1.TaskIdentifier {
	r.mu.RLock()
	defer r.mu.RUnlock()

	identifiers := make([]*workflowsv1.TaskIdentifier, 0, len(r.order))
	for _, identifier := range r.order {
		identifiers = append(identifiers, workflowsv1.TaskIdentifier_builder{
			Name:    identifier.name,
			Version: identifier.version,
		}.Build())
	}
	return identifiers
}

type classifiedTaskError struct {
	err              error
	wasWorkflowError bool
}

func (e *classifiedTaskError) Error() string {
	return e.err.Error()
}

func (e *classifiedTaskError) Unwrap() error {
	return e.err
}

type taskExecutor struct {
	registry *taskRegistry

	cluster string
	client  *Client
	tracer  trace.Tracer
	logger  *slog.Logger
	metrics *taskRunnerMetrics

	redactError  func(error) error
	redactString func(string) string
}

func newTaskExecutor(
	registry *taskRegistry,
	cluster string,
	client *Client,
	tracer trace.Tracer,
	logger *slog.Logger,
	meterProvider metric.MeterProvider,
) (*taskExecutor, error) {
	metrics, err := newTaskRunnerMetrics(meterProvider.Meter(otelMeterName))
	if err != nil {
		return nil, fmt.Errorf("failed to create task runner metrics: %w", err)
	}
	return &taskExecutor{
		registry:     registry,
		cluster:      cluster,
		client:       client,
		tracer:       tracer,
		logger:       logger,
		metrics:      metrics,
		redactError:  func(err error) error { return err },
		redactString: func(value string) string { return value },
	}, nil
}

func (e *taskExecutor) ExecuteTask(ctx context.Context, task *workflowsv1.Task) (*workflowsv1.ExecuteTaskResponse, error) {
	executionContext, err := e.executeTask(ctx, task)
	if err != nil {
		wasWorkflowError := false
		var classifiedError *classifiedTaskError
		if errors.As(err, &classifiedError) {
			wasWorkflowError = classifiedError.wasWorkflowError
		}
		err = e.redactError(err)

		var progressUpdates []*workflowsv1.Progress
		if executionContext != nil {
			progressUpdates = executionContext.getProgressUpdates()
		}
		failedTask := workflowsv1.TaskFailedRequest_builder{
			TaskId:           task.GetId(),
			Display:          failedTaskDisplay(e.redactString(task.GetDisplay()), err),
			WasWorkflowError: wasWorkflowError,
			ProgressUpdates:  progressUpdates,
		}.Build()
		return workflowsv1.ExecuteTaskResponse_builder{FailedTask: failedTask}.Build(), nil
	}

	computedTask := workflowsv1.ComputedTask_builder{
		Id:              task.GetId(),
		Display:         e.redactString(task.GetDisplay()),
		SubTasks:        executionContext.getSubTasks(),
		ProgressUpdates: executionContext.getProgressUpdates(),
	}.Build()
	return workflowsv1.ExecuteTaskResponse_builder{ComputedTask: computedTask}.Build(), nil
}

func (e *taskExecutor) executeTask(ctx context.Context, task *workflowsv1.Task) (*taskExecutionContext, error) {
	beforeTime := time.Now().UTC()

	if task.GetIdentifier() == nil {
		return nil, &classifiedTaskError{err: errors.New("task has no identifier")}
	}
	identifier := NewTaskIdentifier(task.GetIdentifier().GetName(), task.GetIdentifier().GetVersion())
	taskPrototype, found := e.registry.get(identifier)
	if !found {
		return nil, &classifiedTaskError{
			err: fmt.Errorf("task %s@%s is not registered on this runner", identifier.Name(), identifier.Version()),
		}
	}

	return span.StartJobSpan(ctx, e.tracer, fmt.Sprintf("task/%s", identifier.Name()), task.GetJob(), func(ctx context.Context) (taskExecutionContext *taskExecutionContext, err error) { //nolint:nonamedreturns // needed to return state after a recovered panic
		e.logger.DebugContext(ctx, "executing task",
			slog.String("task_id", task.GetId().AsUUID().String()),
			slog.String("task", identifier.Name()),
			slog.String("version", identifier.Version()),
		)

		log := e.logger.With(
			slog.String("job_id", task.GetJob().GetId().AsUUID().String()),
			slog.String("task", identifier.Name()),
			slog.String("version", identifier.Version()),
			slog.Time("start_time", beforeTime),
		)
		taskMetricAttributes := metric.WithAttributes(
			attribute.String("task_identifier", identifier.Name()),
			attribute.String("task_version", identifier.Version()),
		)

		defer func() {
			if recovered := recover(); recovered != nil {
				panicError := e.redactError(fmt.Errorf("task panicked: %v", recovered))
				log.ErrorContext(ctx, "task execution failed", slog.String("error", panicError.Error()), slog.Int64("retry_attempt", task.GetRetryCount()))
				err = &classifiedTaskError{err: panicError, wasWorkflowError: true}

				e.metrics.tasksFailedMetric.Add(ctx, 1, taskMetricAttributes)
				e.metrics.taskExecutionDurationMetric.Record(ctx, time.Since(beforeTime).Seconds(), taskMetricAttributes, metric.WithAttributes(attribute.String("state", "failed")))
			}
		}()

		taskStruct := reflect.New(reflect.TypeOf(taskPrototype).Elem()).Interface().(ExecutableTask)
		if taskMessage, isProtobuf := taskStruct.(proto.Message); isProtobuf {
			if err := proto.Unmarshal(task.GetInput(), taskMessage); err != nil {
				return nil, &classifiedTaskError{err: fmt.Errorf("failed to unmarshal protobuf task: %w", err), wasWorkflowError: true}
			}
		} else if err := json.Unmarshal(task.GetInput(), taskStruct); err != nil {
			return nil, &classifiedTaskError{err: fmt.Errorf("failed to unmarshal json task: %w", err), wasWorkflowError: true}
		}

		e.metrics.taskInputSizeMetric.Record(ctx, int64(len(task.GetInput())), taskMetricAttributes)
		e.metrics.tasksExecutedMetric.Add(ctx, 1, taskMetricAttributes)

		executionContext := e.withTaskExecutionContext(ctx, task)
		executionContext = obslogger.ContextWithSlogAttributes(executionContext, slog.String("task_id", task.GetId().AsUUID().String()))
		err = taskStruct.Execute(executionContext)

		executionTime := time.Since(beforeTime)
		log = log.With(
			slog.Duration("execution_time", executionTime),
			slog.String("execution_time_human", roundDuration(executionTime, 2).String()),
		)

		if err != nil {
			err = e.redactError(err)
			e.metrics.tasksFailedMetric.Add(ctx, 1, taskMetricAttributes)
			e.metrics.taskExecutionDurationMetric.Record(ctx, executionTime.Seconds(), taskMetricAttributes, metric.WithAttributes(attribute.String("state", "failed")))
			log.ErrorContext(executionContext, "task execution failed", slog.String("error", err.Error()), slog.Int64("retry_attempt", task.GetRetryCount()))
			return getTaskExecutionContext(executionContext), &classifiedTaskError{
				err:              fmt.Errorf("failed to execute task: %w", err),
				wasWorkflowError: true,
			}
		}

		e.metrics.tasksComputedMetric.Add(ctx, 1, taskMetricAttributes)
		e.metrics.taskExecutionDurationMetric.Record(ctx, executionTime.Seconds(), taskMetricAttributes, metric.WithAttributes(attribute.String("state", "computed")))

		return getTaskExecutionContext(executionContext), nil
	})
}

func (e *taskExecutor) withTaskExecutionContext(ctx context.Context, task *workflowsv1.Task) context.Context {
	return context.WithValue(ctx, contextKeyTaskExecution, &taskExecutionContext{
		CurrentTask:        task,
		executor:           e,
		subtasks:           make([]*futureTask, 0),
		progressIndicators: make(map[string]*taskProgressIndicator),
	})
}
