package workflows

import (
	"connectrpc.com/connect"
	"context"
	"encoding/json"
	"fmt"
	"github.com/avast/retry-go/v4"
	"github.com/google/uuid"
	"github.com/tilebox/tilebox-go/observability"
	workflowsv1 "github.com/tilebox/tilebox-go/protogen/go/workflows/v1"
	"github.com/tilebox/tilebox-go/protogen/go/workflows/v1/workflowsv1connect"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"
	"log/slog"
	"math/rand/v2"
	"os/signal"
	"reflect"
	"syscall"
	"time"
)

type ContextKeyTaskExecutionType string

const ContextKeyTaskExecution ContextKeyTaskExecutionType = "x-tilebox-task-execution-object"

const DefaultClusterSlug = "testing-4qgCk4qHH85qR7"

//const DefaultClusterSlug = "workflow-dev-EifhUozDpwAJDL"

const pollingInterval = 5 * time.Second
const jitterInterval = 5 * time.Second

type TaskRunner struct {
	Client          workflowsv1connect.TaskServiceClient
	taskDefinitions map[TaskIdentifier]ExecutableTask
	tracer          trace.Tracer
}

func NewTaskRunner(client workflowsv1connect.TaskServiceClient) *TaskRunner {
	return &TaskRunner{
		Client:          client,
		taskDefinitions: make(map[TaskIdentifier]ExecutableTask),
		tracer:          otel.Tracer("tilebox.com/observability"),
	}
}

func (t *TaskRunner) RegisterTask(task ExecutableTask) error {
	identifier := identifierFromTask(task)
	err := ValidateIdentifier(identifier)
	if err != nil {
		return err
	}
	t.taskDefinitions[identifier] = task
	return nil
}

func (t *TaskRunner) RegisterTasks(task ...ExecutableTask) error {
	for _, task := range task {
		err := t.RegisterTask(task)
		if err != nil {
			return err
		}
	}
	return nil
}

func protobufToUuid(id *workflowsv1.UUID) (uuid.UUID, error) {
	if id == nil || len(id.Uuid) == 0 {
		return uuid.Nil, nil
	}

	bytes, err := uuid.FromBytes(id.Uuid)
	if err != nil {
		return uuid.Nil, err
	}

	return bytes, nil
}

func isEmpty(id *workflowsv1.UUID) bool {
	taskId, err := protobufToUuid(id)
	if err != nil {
		return false
	}
	return taskId == uuid.Nil
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
			Name:    identifier.Name,
			Version: identifier.Version,
		})
	}

	var task *workflowsv1.Task

	for {
		if task == nil { // if we don't have a task, let's try work-stealing one
			taskResponse, err := t.Client.NextTask(ctx, connect.NewRequest(&workflowsv1.NextTaskRequest{
				NextTaskToRun: &workflowsv1.NextTaskToRun{ClusterSlug: DefaultClusterSlug, Identifiers: identifiers},
			}))
			if err != nil {
				slog.ErrorContext(ctx, "failed to work-steal a task", "error", err)
				// return  // should we even try again, or just stop here?
			} else {
				task = taskResponse.Msg.NextTask
			}
		}

		if task != nil { // we have a task to execute
			if isEmpty(task.Id) {
				slog.ErrorContext(ctx, "got a task without an ID - skipping to the next task")
				task = nil
				continue
			}
			executionContext, err := t.executeTask(ctx, task)
			stopExecution := false
			if err == nil { // in case we got no error, let's mark the task as computed and get the next one
				computedTask := &workflowsv1.ComputedTask{
					Id:       task.Id,
					SubTasks: nil,
				}
				if executionContext != nil && len(executionContext.Subtasks) > 0 {
					computedTask.SubTasks = executionContext.Subtasks
				}
				nextTaskToRun := &workflowsv1.NextTaskToRun{ClusterSlug: DefaultClusterSlug, Identifiers: identifiers}
				select {
				case <-ctxSignal.Done():
					// if we got a context cancellation, don't request a new task
					nextTaskToRun = nil
					stopExecution = true
				default:
				}

				task, err = retry.DoWithData(
					func() (*workflowsv1.Task, error) {
						taskResponse, err := t.Client.NextTask(ctx, connect.NewRequest(&workflowsv1.NextTaskRequest{
							ComputedTask: computedTask, NextTaskToRun: nextTaskToRun,
						}))
						if err != nil {
							slog.ErrorContext(ctx, "failed to mark task as computed, retrying", "error", err)
							return nil, err
						}
						return taskResponse.Msg.NextTask, nil
					}, retry.Context(ctxSignal), retry.DelayType(retry.CombineDelay(retry.BackOffDelay, retry.RandomDelay)),
				)
				if err != nil {
					slog.ErrorContext(ctx, "failed to retry NextTask", "error", err)
					return // we got a cancellation signal, so let's just stop here
				}
			} else { // err != nil
				slog.ErrorContext(ctx, "task execution failed", "error", err)
				err = retry.Do(
					func() error {
						_, err := t.Client.TaskFailed(ctx, connect.NewRequest(&workflowsv1.TaskFailedRequest{
							TaskId:    task.Id,
							CancelJob: true,
						}))
						if err != nil {
							slog.ErrorContext(ctx, "failed to report task failure", "error", err)
							return err
						}
						return nil
					}, retry.Context(ctxSignal), retry.DelayType(retry.CombineDelay(retry.BackOffDelay, retry.RandomDelay)),
				)
				if err != nil {
					slog.ErrorContext(ctx, "failed to retry TaskFailed", "error", err)
					return // we got a cancellation signal, so let's just stop here
				}
				task = nil // reported a task failure, let's work-steal again
			}
			if stopExecution {
				return
			}
		} else {
			// if we didn't get a task, let's wait for a bit and try work-stealing again
			slog.DebugContext(ctx, "no task to run")

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
	go extendTaskLease(leaseCtx, t.Client, task.Id, task.Lease.Lease.AsDuration(), task.Lease.RecommendedWaitUntilNextExtension.AsDuration())
	defer stopLeaseExtensions()

	// actually execute the task
	if task.Identifier == nil {
		return nil, fmt.Errorf("task has no identifier")
	}
	identifier := TaskIdentifier{Name: task.Identifier.Name, Version: task.Identifier.Version}
	taskPrototype, found := t.taskDefinitions[identifier]
	if !found {
		return nil, fmt.Errorf("task %s is not registered on this runner", task.Identifier.Name)
	}

	return observability.StartJobSpan(t.tracer, ctx, fmt.Sprintf("task/%s", identifier.Name), task.GetJob(), func(ctx context.Context) (*taskExecutionContext, error) {
		slog.DebugContext(ctx, "executing task", "task", identifier.Name, "version", identifier.Version)
		taskStruct := reflect.New(reflect.ValueOf(taskPrototype).Elem().Type()).Interface().(ExecutableTask)

		_, isProtobuf := taskStruct.(proto.Message)
		if isProtobuf {
			err := proto.Unmarshal(task.Input, taskStruct.(proto.Message))
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal protobuf task: %w", err)
			}
		} else {
			err := json.Unmarshal(task.Input, taskStruct)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal json task: %w", err)
			}
		}

		executionContext := withTaskExecutionContext(ctx, t.Client, task)
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
func extendTaskLease(ctx context.Context, client workflowsv1connect.TaskServiceClient, taskId *workflowsv1.UUID, initialLease, initialWait time.Duration) {
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
		slog.DebugContext(ctx, "extending task lease", "task_id", uuid.Must(uuid.FromBytes(taskId.Uuid)), "lease", lease, "wait", wait)
		req := &workflowsv1.TaskLeaseRequest{
			TaskId:         taskId,
			RequestedLease: durationpb.New(2 * lease), // double the current lease duration for the next extension
		}
		extension, err := client.ExtendTaskLease(ctx, connect.NewRequest(req))
		if err != nil {
			slog.ErrorContext(ctx, "failed to extend task lease", "error", err, "task_id", uuid.Must(uuid.FromBytes(taskId.Uuid)))
			// The server probably has an internal error, but there is no point in trying to extend the lease again
			// because it will be expired then, so let's just return
			return
		}
		if extension.Msg.Lease == nil {
			// the server did not return a lease extension, it means that there is no need in trying to extend the lease
			slog.DebugContext(ctx, "task lease extension not granted", "task_id", uuid.Must(uuid.FromBytes(taskId.Uuid)))
			return
		}
		// will probably be double the previous lease (since we requested that) or capped by the server at maxLeaseDuration
		lease = extension.Msg.Lease.AsDuration()
		wait = extension.Msg.RecommendedWaitUntilNextExtension.AsDuration()
	}
}

type taskExecutionContext struct {
	CurrentTask *workflowsv1.Task
	Client      workflowsv1connect.TaskServiceClient
	Subtasks    []*workflowsv1.TaskSubmission
}

func withTaskExecutionContext(ctx context.Context, client workflowsv1connect.TaskServiceClient, task *workflowsv1.Task) context.Context {
	return context.WithValue(ctx, ContextKeyTaskExecution, &taskExecutionContext{
		CurrentTask: task,
		Client:      client,
		Subtasks:    make([]*workflowsv1.TaskSubmission, 0),
	})
}

func getTaskExecutionContext(ctx context.Context) *taskExecutionContext {
	return ctx.Value(ContextKeyTaskExecution).(*taskExecutionContext)
}

func SubmitSubtasks(ctx context.Context, tasks ...Task) error {
	executionContext := getTaskExecutionContext(ctx)
	if executionContext == nil {
		return fmt.Errorf("cannot submit subtask without task execution context")
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
			ClusterSlug: DefaultClusterSlug,
			Identifier: &workflowsv1.TaskIdentifier{
				Name:    identifier.Name,
				Version: identifier.Version,
			},
			Input:        subtaskInput,
			Dependencies: nil,
			Display:      identifier.Name,
		})
	}

	return nil
}
