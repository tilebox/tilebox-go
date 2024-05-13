package workflows

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"os/signal"
	"reflect"
	"syscall"
	"time"

	"go.opentelemetry.io/otel"

	"connectrpc.com/connect"
	"github.com/avast/retry-go/v4"
	"github.com/google/uuid"
	"github.com/tilebox/tilebox-go/observability"
	workflowsv1 "github.com/tilebox/tilebox-go/protogen/go/workflows/v1"
	"github.com/tilebox/tilebox-go/protogen/go/workflows/v1/workflowsv1connect"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"
)

type ContextKeyTaskExecutionType string

const ContextKeyTaskExecution ContextKeyTaskExecutionType = "x-tilebox-task-execution-object"

const pollingInterval = 5 * time.Second
const jitterInterval = 5 * time.Second

type taskRunnerConfig struct {
	clusterSlug    string
	tracerProvider trace.TracerProvider
	tracerName     string
	logger         *slog.Logger
}

type TaskRunnerOption func(*taskRunnerConfig)

func WithCluster(clusterSlug string) TaskRunnerOption {
	return func(cfg *taskRunnerConfig) {
		cfg.clusterSlug = clusterSlug
	}
}

func WithRunnerTracerProvider(tracerProvider trace.TracerProvider) TaskRunnerOption {
	return func(cfg *taskRunnerConfig) {
		cfg.tracerProvider = tracerProvider
	}
}

func WithRunnerTracerName(tracerName string) TaskRunnerOption {
	return func(cfg *taskRunnerConfig) {
		cfg.tracerName = tracerName
	}
}

func WithRunnerLogger(logger *slog.Logger) TaskRunnerOption {
	return func(cfg *taskRunnerConfig) {
		cfg.logger = logger
	}
}

func newTaskRunnerConfig(options []TaskRunnerOption) (*taskRunnerConfig, error) {
	cfg := &taskRunnerConfig{
		tracerProvider: otel.GetTracerProvider(),    // use the global tracer provider by default
		tracerName:     "tilebox.com/observability", // the default tracer name we use
		logger:         slog.Default(),
	}
	for _, option := range options {
		option(cfg)
	}

	if cfg.clusterSlug == "" {
		return nil, errors.New("cluster slug is required")
	}

	return cfg, nil
}

type TaskRunner struct {
	client          workflowsv1connect.TaskServiceClient
	taskDefinitions map[taskIdentifier]ExecutableTask

	cluster string
	tracer  trace.Tracer
	logger  *slog.Logger
}

func NewTaskRunner(client workflowsv1connect.TaskServiceClient, options ...TaskRunnerOption) (*TaskRunner, error) {
	cfg, err := newTaskRunnerConfig(options)
	if err != nil {
		return nil, err
	}
	return &TaskRunner{
		client:          client,
		taskDefinitions: make(map[taskIdentifier]ExecutableTask),

		cluster: cfg.clusterSlug,
		tracer:  cfg.tracerProvider.Tracer(cfg.tracerName),
		logger:  cfg.logger,
	}, nil
}

func (t *TaskRunner) RegisterTask(task ExecutableTask) error {
	identifier := identifierFromTask(task)
	err := ValidateIdentifier(identifier)
	if err != nil {
		return err
	}
	t.taskDefinitions[taskIdentifier{name: identifier.Name(), version: identifier.Version()}] = task
	return nil
}

func (t *TaskRunner) GetRegisteredTask(identifier TaskIdentifier) (ExecutableTask, bool) {
	registeredTask, found := t.taskDefinitions[taskIdentifier{name: identifier.Name(), version: identifier.Version()}]
	return registeredTask, found
}

func (t *TaskRunner) RegisterTasks(tasks ...ExecutableTask) error {
	for _, task := range tasks {
		err := t.RegisterTask(task)
		if err != nil {
			return err
		}
	}
	return nil
}

func protobufToUUID(id *workflowsv1.UUID) (uuid.UUID, error) {
	if id == nil || len(id.GetUuid()) == 0 {
		return uuid.Nil, nil
	}

	bytes, err := uuid.FromBytes(id.GetUuid())
	if err != nil {
		return uuid.Nil, err
	}

	return bytes, nil
}

func isEmpty(id *workflowsv1.UUID) bool {
	taskID, err := protobufToUUID(id)
	if err != nil {
		return false
	}
	return taskID == uuid.Nil
}

// Run runs the task runner forever, looking for new tasks to run and polling for new tasks when idle.
func (t *TaskRunner) Run(ctx context.Context) {
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
			taskResponse, err := t.client.NextTask(ctx, connect.NewRequest(&workflowsv1.NextTaskRequest{
				NextTaskToRun: &workflowsv1.NextTaskToRun{ClusterSlug: t.cluster, Identifiers: identifiers},
			}))
			if err != nil {
				t.logger.ErrorContext(ctx, "failed to work-steal a task", "error", err)
				// return  // should we even try again, or just stop here?
			} else {
				task = taskResponse.Msg.GetNextTask()
			}
		}

		if task != nil { // we have a task to execute
			if isEmpty(task.GetId()) {
				t.logger.ErrorContext(ctx, "got a task without an ID - skipping to the next task")
				task = nil
				continue
			}
			executionContext, err := t.executeTask(ctx, task)
			stopExecution := false
			if err == nil { // in case we got no error, let's mark the task as computed and get the next one
				computedTask := &workflowsv1.ComputedTask{
					Id:       task.GetId(),
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
				default:
				}

				task, err = retry.DoWithData(
					func() (*workflowsv1.Task, error) {
						taskResponse, err := t.client.NextTask(ctx, connect.NewRequest(&workflowsv1.NextTaskRequest{
							ComputedTask: computedTask, NextTaskToRun: nextTaskToRun,
						}))
						if err != nil {
							t.logger.ErrorContext(ctx, "failed to mark task as computed, retrying", "error", err)
							return nil, err
						}
						return taskResponse.Msg.GetNextTask(), nil
					}, retry.Context(ctxSignal), retry.DelayType(retry.CombineDelay(retry.BackOffDelay, retry.RandomDelay)),
				)
				if err != nil {
					t.logger.ErrorContext(ctx, "failed to retry NextTask", "error", err)
					return // we got a cancellation signal, so let's just stop here
				}
			} else { // err != nil
				t.logger.ErrorContext(ctx, "task execution failed", "error", err)
				err = retry.Do(
					func() error {
						_, err := t.client.TaskFailed(ctx, connect.NewRequest(&workflowsv1.TaskFailedRequest{
							TaskId:    task.GetId(),
							CancelJob: true,
						}))
						if err != nil {
							t.logger.ErrorContext(ctx, "failed to report task failure", "error", err)
							return err
						}
						return nil
					}, retry.Context(ctxSignal), retry.DelayType(retry.CombineDelay(retry.BackOffDelay, retry.RandomDelay)),
				)
				if err != nil {
					t.logger.ErrorContext(ctx, "failed to retry TaskFailed", "error", err)
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

			// instead of time.Sleep we set a timer and select on it, so we still can catch signals like SIGINT
			timer := time.NewTimer(pollingInterval + rand.N(jitterInterval))
			select {
			case <-ctxSignal.Done():
				timer.Stop() // stop the timer before returning, avoids a memory leak
				return       // if we got a context cancellation, let's just stop here
			case <-timer.C: // the timer expired, let's try to work-steal a task again
			}
		}
	}
}

func (t *TaskRunner) executeTask(ctx context.Context, task *workflowsv1.Task) (*taskExecutionContext, error) {
	// start a goroutine to extend the lease of the task continuously until the task execution is finished
	leaseCtx, stopLeaseExtensions := context.WithCancel(ctx)
	go t.extendTaskLease(leaseCtx, t.client, task.GetId(), task.GetLease().GetLease().AsDuration(), task.GetLease().GetRecommendedWaitUntilNextExtension().AsDuration())
	defer stopLeaseExtensions()

	// actually execute the task
	if task.GetIdentifier() == nil {
		return nil, errors.New("task has no identifier")
	}
	identifier := NewTaskIdentifier(task.GetIdentifier().GetName(), task.GetIdentifier().GetVersion())
	taskPrototype, found := t.GetRegisteredTask(identifier)
	if !found {
		return nil, fmt.Errorf("task %s is not registered on this runner", task.GetIdentifier().GetName())
	}

	return observability.StartJobSpan(ctx, t.tracer, fmt.Sprintf("task/%s", identifier.Name()), task.GetJob(), func(ctx context.Context) (*taskExecutionContext, error) {
		t.logger.DebugContext(ctx, "executing task", "task", identifier.Name, "version", identifier.Version)
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

		executionContext := t.withTaskExecutionContext(ctx, task)
		err := taskStruct.Execute(executionContext)
		if r := recover(); r != nil {
			// recover from panics during task executions, so we can still report the error to the server and continue
			// with other tasks
			return nil, fmt.Errorf("task panicked: %v", r)
		}
		if err != nil {
			return nil, fmt.Errorf("failed to execute task: %w", err)
		}

		return getTaskExecutionContext(executionContext), nil
	})
}

// extendTaskLease is a function designed to be run as a goroutine, extending the lease of a task continuously until the
// context is cancelled, which indicates that the execution of the task is finished.
func (t *TaskRunner) extendTaskLease(ctx context.Context, client workflowsv1connect.TaskServiceClient, taskID *workflowsv1.UUID, initialLease, initialWait time.Duration) {
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
		t.logger.DebugContext(ctx, "extending task lease", "task_id", uuid.Must(uuid.FromBytes(taskID.GetUuid())), "lease", lease, "wait", wait)
		req := &workflowsv1.TaskLeaseRequest{
			TaskId:         taskID,
			RequestedLease: durationpb.New(2 * lease), // double the current lease duration for the next extension
		}
		extension, err := client.ExtendTaskLease(ctx, connect.NewRequest(req))
		if err != nil {
			t.logger.ErrorContext(ctx, "failed to extend task lease", "error", err, "task_id", uuid.Must(uuid.FromBytes(taskID.GetUuid())))
			// The server probably has an internal error, but there is no point in trying to extend the lease again
			// because it will be expired then, so let's just return
			return
		}
		if extension.Msg.GetLease() == nil {
			// the server did not return a lease extension, it means that there is no need in trying to extend the lease
			t.logger.DebugContext(ctx, "task lease extension not granted", "task_id", uuid.Must(uuid.FromBytes(taskID.GetUuid())))
			return
		}
		// will probably be double the previous lease (since we requested that) or capped by the server at maxLeaseDuration
		lease = extension.Msg.GetLease().AsDuration()
		wait = extension.Msg.GetRecommendedWaitUntilNextExtension().AsDuration()
	}
}

type taskExecutionContext struct {
	CurrentTask *workflowsv1.Task
	runner      *TaskRunner
	Subtasks    []*workflowsv1.TaskSubmission
}

func (t *TaskRunner) withTaskExecutionContext(ctx context.Context, task *workflowsv1.Task) context.Context {
	return context.WithValue(ctx, ContextKeyTaskExecution, &taskExecutionContext{
		CurrentTask: task,
		runner:      t,
		Subtasks:    make([]*workflowsv1.TaskSubmission, 0),
	})
}

func getTaskExecutionContext(ctx context.Context) *taskExecutionContext {
	return ctx.Value(ContextKeyTaskExecution).(*taskExecutionContext)
}

func SubmitSubtasks(ctx context.Context, tasks ...Task) error {
	executionContext := getTaskExecutionContext(ctx)
	if executionContext == nil {
		return errors.New("cannot submit subtask without task execution context")
	}

	for _, task := range tasks {
		var subtaskInput []byte
		var err error

		taskProto, isProtobuf := task.(proto.Message)
		if isProtobuf {
			subtaskInput, err = proto.Marshal(taskProto)
			if err != nil {
				return fmt.Errorf("failed to marshal protobuf task: %w", err)
			}
		} else {
			subtaskInput, err = json.Marshal(task)
			if err != nil {
				return fmt.Errorf("failed to marshal task: %w", err)
			}
		}

		identifier := identifierFromTask(task)
		err = ValidateIdentifier(identifier)
		if err != nil {
			return fmt.Errorf("subtask has invalid task identifier: %w", err)
		}

		executionContext.Subtasks = append(executionContext.Subtasks, &workflowsv1.TaskSubmission{
			ClusterSlug: executionContext.runner.cluster,
			Identifier: &workflowsv1.TaskIdentifier{
				Name:    identifier.Name(),
				Version: identifier.Version(),
			},
			Input:        subtaskInput,
			Dependencies: nil,
			Display:      identifier.Display(),
		})
	}

	return nil
}
