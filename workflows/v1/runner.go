package workflows // import "github.com/tilebox/tilebox-go/workflows/v1"

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"os/signal"
	"reflect"
	"strings"
	"syscall"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/google/uuid"
	"github.com/tilebox/tilebox-go/observability"
	workflowsv1 "github.com/tilebox/tilebox-go/protogen/go/workflows/v1"
	"github.com/tilebox/tilebox-go/workflows/v1/runner"
	"github.com/tilebox/tilebox-go/workflows/v1/subtask"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/proto"
)

type contextKeyTaskExecutionType string

const contextKeyTaskExecution contextKeyTaskExecutionType = "x-tilebox-task-execution-object"

const (
	pollingInterval = 5 * time.Second
	jitterInterval  = 5 * time.Second
)

// TaskRunner executes tasks.
type TaskRunner struct {
	service         TaskService
	taskDefinitions map[taskIdentifier]ExecutableTask

	cluster            string
	tracer             trace.Tracer
	logger             *slog.Logger
	taskDurationMetric metric.Float64Histogram
}

func newTaskRunner(service TaskService, tracer trace.Tracer, cluster *Cluster, options ...runner.Option) (*TaskRunner, error) {
	opts := &runner.Options{
		Logger:        slog.Default(),
		MeterProvider: otel.GetMeterProvider(),
	}
	for _, option := range options {
		option(opts)
	}

	meter := opts.MeterProvider.Meter(otelMeterName)

	taskDurationMetric, err := meter.Float64Histogram(
		"task.execution.duration",
		metric.WithDescription("Task execution duration"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create task execution duration metric: %w", err)
	}

	return &TaskRunner{
		service:         service,
		taskDefinitions: make(map[taskIdentifier]ExecutableTask),

		cluster:            cluster.Slug,
		tracer:             tracer,
		logger:             opts.Logger,
		taskDurationMetric: taskDurationMetric,
	}, nil
}

// registerTask makes the task runner aware of a task.
func (t *TaskRunner) registerTask(task ExecutableTask) error {
	identifier := identifierFromTask(task)
	err := ValidateIdentifier(identifier)
	if err != nil {
		return err
	}
	t.taskDefinitions[taskIdentifier{name: identifier.Name(), version: identifier.Version()}] = task
	return nil
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

func isEmpty(id *workflowsv1.UUID) bool {
	taskID, err := protoToUUID(id)
	if err != nil {
		return false
	}
	return taskID == uuid.Nil
}

// RunForever runs the task runner forever, looking for new tasks to run and polling for new tasks when idle.
func (t *TaskRunner) RunForever(ctx context.Context) {
	t.run(ctx, false)
}

// RunAll run the task runner and execute all tasks until there are no more tasks available.
func (t *TaskRunner) RunAll(ctx context.Context) {
	t.run(ctx, true)
}

func (t *TaskRunner) run(ctx context.Context, stopWhenIdling bool) {
	// Catch signals to gracefully shutdown
	ctxSignal, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGTSTP, syscall.SIGQUIT)
	defer stop()

	identifiers := make([]*workflowsv1.TaskIdentifier, 0)

	for _, task := range t.taskDefinitions {
		identifier := identifierFromTask(task)
		identifiers = append(identifiers, &workflowsv1.TaskIdentifier{
			Name:    identifier.Name(),
			Version: identifier.Version(),
		})
	}

	var task *workflowsv1.Task

	for {
		if task == nil { // if we don't have a task, let's try work-stealing one
			taskResponse, err := t.service.NextTask(ctx,
				nil,
				&workflowsv1.NextTaskToRun{ClusterSlug: t.cluster, Identifiers: identifiers},
			)
			if err != nil {
				t.logger.ErrorContext(ctx, "failed to work-steal a task", slog.Any("error", err))
				// return  // should we even try again, or just stop here?
			} else {
				task = taskResponse.GetNextTask()
			}
		}

		if task != nil { // we have a task to execute
			if isEmpty(task.GetId()) {
				t.logger.ErrorContext(ctx, "got a task without an ID - skipping to the next task")
				task = nil
				continue
			}
			executionContext, errTask := t.executeTask(ctx, task)
			stopExecution := false
			if errTask == nil { // in case we got no error, let's mark the task as computed and get the next one
				computedTask := &workflowsv1.ComputedTask{
					Id:       task.GetId(),
					Display:  task.GetDisplay(),
					SubTasks: nil,
				}
				if executionContext != nil && len(executionContext.Subtasks) > 0 {
					computedTask.SubTasks = executionContext.Subtasks
				}
				nextTaskToRun := &workflowsv1.NextTaskToRun{ClusterSlug: t.cluster, Identifiers: identifiers}
				select {
				case <-ctxSignal.Done():
					// if we got a context cancellation, don't request a new task
					nextTaskToRun = nil
					stopExecution = true
				case <-ctx.Done():
					// if we got a context cancellation, don't request a new task
					nextTaskToRun = nil
					stopExecution = true
				default:
				}

				var err error
				task, err = retry.DoWithData(
					func() (*workflowsv1.Task, error) {
						taskResponse, err := t.service.NextTask(ctx, computedTask, nextTaskToRun)
						if err != nil {
							t.logger.ErrorContext(ctx, "failed to mark task as computed, retrying", slog.Any("error", err))
							return nil, err
						}
						return taskResponse.GetNextTask(), nil
					}, retry.Context(ctxSignal), retry.DelayType(retry.CombineDelay(retry.BackOffDelay, retry.RandomDelay)),
				)
				if err != nil {
					if !errors.Is(err, context.Canceled) {
						t.logger.ErrorContext(ctx, "failed to retry NextTask", slog.Any("error", err))
					}
					return // we got a cancellation signal, so let's just stop here
				}
			} else { // errTask != nil
				taskID, err := protoToUUID(task.GetId())
				if err != nil {
					t.logger.ErrorContext(ctx, "failed to convert task id to uuid", slog.Any("error", err))
					return
				}

				// the error itself is already logged in executeTask, so we just need to report the task as failed
				err = t.taskFailed(ctx, ctxSignal, taskID, errTask, task.GetDisplay())
				if err != nil {
					if !errors.Is(err, context.Canceled) {
						t.logger.ErrorContext(ctx, "failed to retry TaskFailed", slog.Any("error", err))
					}
					return // we got a cancellation signal, so let's just stop here
				}
				task = nil // reported a task failure, let's work-steal again
			}
			if stopExecution {
				return
			}
		} else {
			// if we didn't get a task, let's wait for a bit and try work-stealing again
			t.logger.DebugContext(ctx, "no task to run")

			if stopWhenIdling {
				return
			}

			// instead of time.Sleep we set a timer and select on it, so we still can catch signals like SIGINT
			timer := time.NewTimer(pollingInterval + rand.N(jitterInterval))
			select {
			case <-ctxSignal.Done():
				timer.Stop() // stop the timer before returning, avoids a memory leak
				return       // if we got a context cancellation, let's just stop here
			case <-ctx.Done():
				timer.Stop() // stop the timer before returning, avoids a memory leak
				return       // if we got a context cancellation, let's just stop here
			case <-timer.C: // the timer expired, let's try to work-steal a task again
			}
		}
	}
}

func (t *TaskRunner) taskFailed(ctx, ctxSignal context.Context, taskID uuid.UUID, taskError error, taskDisplay string) error {
	// job output is limited to 1KB, so truncate the error message if necessary
	errorMessage := taskError.Error()
	if len(errorMessage) > 1024 {
		errorMessage = errorMessage[:1024]
	}

	display := taskDisplay
	if errorMessage != "" {
		display = strings.Join([]string{taskDisplay, errorMessage}, "\n")
	}

	return retry.Do(
		func() error {
			_, err := t.service.TaskFailed(ctx, taskID, display, true)
			if err != nil {
				t.logger.ErrorContext(ctx, "failed to report task failure", slog.Any("error", err))
				return err
			}
			return nil
		}, retry.Context(ctxSignal), retry.DelayType(retry.CombineDelay(retry.BackOffDelay, retry.RandomDelay)),
	)
}

func (t *TaskRunner) executeTask(ctx context.Context, task *workflowsv1.Task) (*taskExecutionContext, error) {
	taskID, err := protoToUUID(task.GetId())
	if err != nil {
		return nil, fmt.Errorf("failed to convert task id to uuid: %w", err)
	}

	// start a goroutine to extend the lease of the task continuously until the task execution is finished
	leaseCtx, stopLeaseExtensions := context.WithCancel(ctx)

	go t.extendTaskLease(
		leaseCtx,
		t.service,
		taskID,
		task.GetLease().GetLease().AsDuration(),
		task.GetLease().GetRecommendedWaitUntilNextExtension().AsDuration(),
	)
	defer stopLeaseExtensions()

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

	jobID, err := protoToUUID(task.GetJob().GetId())
	if err != nil {
		return nil, fmt.Errorf("failed to convert job id to uuid: %w", err)
	}

	return observability.StartJobSpan(ctx, t.tracer, fmt.Sprintf("task/%s", identifier.Name()), task.GetJob(), func(ctx context.Context) (taskExecutionContext *taskExecutionContext, err error) { //nolint:nonamedreturns // needed to return a value in case of panic
		t.logger.DebugContext(ctx, "executing task", slog.String("task", identifier.Name()), slog.String("version", identifier.Version()))
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
			slog.String("job_id", jobID.String()),
			slog.String("task", identifier.Name()),
			slog.String("version", identifier.Version()),
			slog.Time("start_time", beforeTime),
		)

		defer func() {
			if r := recover(); r != nil {
				// recover from panics during task executions, so we can still report the error to the server and continue
				// with other tasks
				log.ErrorContext(ctx, "task execution failed", slog.String("error", "panic"), slog.Int64("retry_attempt", task.GetRetryCount()))
				taskExecutionContext = nil
				err = fmt.Errorf("task panicked: %v", r)
			}
		}()

		executionContext := t.withTaskExecutionContext(ctx, task)
		err = taskStruct.Execute(executionContext)

		executionTime := time.Since(beforeTime)

		log = log.With(
			slog.Duration("execution_time", executionTime),
			slog.String("execution_time_human", roundDuration(executionTime, 2).String()),
		)

		if err != nil {
			log.ErrorContext(ctx, "task execution failed", slog.Any("error", err), slog.Int64("retry_attempt", task.GetRetryCount()))
			return nil, fmt.Errorf("failed to execute task: %w", err)
		}

		// record the time it took to run a successful task
		t.taskDurationMetric.Record(ctx, executionTime.Seconds())

		return getTaskExecutionContext(executionContext), nil
	})
}

// extendTaskLease is a function designed to be run as a goroutine, extending the lease of a task continuously until the
// context is cancelled, which indicates that the execution of the task is finished.
func (t *TaskRunner) extendTaskLease(ctx context.Context, service TaskService, taskID uuid.UUID, initialLease, initialWait time.Duration) {
	wait := initialWait
	lease := initialLease
	for {
		timer := time.NewTimer(wait)
		select {
		case <-ctx.Done(): // if the context is cancelled, let's stop trying to extend the lease -> the task finished
			timer.Stop() // stop the timer before returning, avoids a memory leak
			return
		case <-timer.C: // the timer expired, let's try to extend the lease
		}
		t.logger.DebugContext(ctx, "extending task lease", slog.String("task_id", taskID.String()), slog.Duration("lease", lease), slog.Duration("wait", wait))
		// double the current lease duration for the next extension
		extension, err := service.ExtendTaskLease(ctx, taskID, 2*lease)
		if err != nil {
			t.logger.ErrorContext(ctx, "failed to extend task lease", slog.Any("error", err), slog.String("task_id", taskID.String()))
			// The server probably has an internal error, but there is no point in trying to extend the lease again
			// because it will be expired then, so let's just return
			return
		}
		if extension.GetLease() == nil {
			// the server did not return a lease extension, it means that there is no need in trying to extend the lease
			t.logger.DebugContext(ctx, "task lease extension not granted", slog.String("task_id", taskID.String()))
			return
		}
		// will probably be double the previous lease (since we requested that) or capped by the server at maxLeaseDuration
		lease = extension.GetLease().AsDuration()
		wait = extension.GetRecommendedWaitUntilNextExtension().AsDuration()
	}
}

type taskExecutionContext struct {
	CurrentTask *workflowsv1.Task
	runner      *TaskRunner
	Subtasks    []*workflowsv1.TaskSubmission
}

func (t *TaskRunner) withTaskExecutionContext(ctx context.Context, task *workflowsv1.Task) context.Context {
	return context.WithValue(ctx, contextKeyTaskExecution, &taskExecutionContext{
		CurrentTask: task,
		runner:      t,
		Subtasks:    make([]*workflowsv1.TaskSubmission, 0),
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

// SetTaskDisplay sets the display name of the current task.
func SetTaskDisplay(ctx context.Context, display string) error {
	executionContext := getTaskExecutionContext(ctx)
	if executionContext == nil {
		return errors.New("cannot set task display name without task execution context")
	}
	executionContext.CurrentTask.Display = &display
	return nil
}

// SubmitSubtask submits a task to the task runner as subtask of the current task.
func SubmitSubtask(ctx context.Context, task Task, options ...subtask.SubmitOption) (subtask.FutureTask, error) {
	executionContext := getTaskExecutionContext(ctx)
	if executionContext == nil {
		return 0, errors.New("cannot submit subtask without task execution context")
	}

	opts := &subtask.SubmitOptions{
		Dependencies: nil,
		ClusterSlug:  executionContext.runner.cluster,
		MaxRetries:   0,
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

	var dependencies []int64 //nolint:prealloc // we want to keep it nil if there are no dependencies
	taskIndex := int64(len(executionContext.Subtasks))

	for _, ft := range opts.Dependencies {
		if ft < 0 || int64(ft) >= taskIndex {
			return 0, fmt.Errorf("invalid dependency: future task %d doesn't exist", ft)
		}
		dependencies = append(dependencies, int64(ft))
	}

	subtaskSubmission := &workflowsv1.TaskSubmission{
		ClusterSlug: opts.ClusterSlug,
		Identifier: &workflowsv1.TaskIdentifier{
			Name:    identifier.Name(),
			Version: identifier.Version(),
		},
		Input:        subtaskInput,
		Dependencies: dependencies,
		MaxRetries:   opts.MaxRetries,
		Display:      identifier.Display(),
	}

	executionContext.Subtasks = append(executionContext.Subtasks, subtaskSubmission)
	return subtask.FutureTask(taskIndex), nil
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

// WithTaskSpanResult is a helper function that wraps a function with a tracing span.
// It returns the result of the function and an error if any.
func WithTaskSpanResult[Result any](ctx context.Context, name string, f func(ctx context.Context) (Result, error)) (Result, error) {
	executionContext := getTaskExecutionContext(ctx)
	if executionContext == nil || executionContext.runner.tracer == nil {
		// if we don't have a task execution context or the tracer is not configured, just run the function without creating a span
		return f(ctx)
	}
	return observability.WithSpanResult(ctx, executionContext.runner.tracer, name, f)
}

// WithTaskSpan is a helper function that wraps a function with a tracing span.
func WithTaskSpan(ctx context.Context, name string, f func(ctx context.Context) error) error {
	executionContext := getTaskExecutionContext(ctx)
	if executionContext == nil || executionContext.runner.tracer == nil {
		// if we don't have a task execution context or the tracer is not configured, just run the function without creating a span
		return f(ctx)
	}
	return observability.WithSpan(ctx, executionContext.runner.tracer, name, f)
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
