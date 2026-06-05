package workflows

import (
	"context"
	"errors"
	"log/slog"
	"math/rand/v2"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	workflowsv1 "github.com/tilebox/tilebox-go/protogen/workflows/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TaskExecutor is responsible for actually executing a task and keeping track of which tasks can be executed.
type TaskExecutor interface {
	// TaskIdentifiers returns the current task capabilities of this executor.
	// This method must be cheap and non-blocking. It should only read current in-memory state and must not poll remote
	// APIs, download artifacts, or start runtimes.
	TaskIdentifiers() []*workflowsv1.TaskIdentifier

	ExecuteTask(ctx context.Context, task *workflowsv1.Task) (*workflowsv1.ExecuteTaskResponse, error)
}

type TaskExecutionError struct {
	Err             error
	ProgressUpdates []*workflowsv1.Progress
}

func (e *TaskExecutionError) Error() string {
	if e == nil || e.Err == nil {
		return "task execution failed"
	}
	return e.Err.Error()
}

func (e *TaskExecutionError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

type PollingTaskRunner struct {
	service TaskService

	clusterSlug string
	logger      *slog.Logger
	executor    TaskExecutor

	requestNewTasks atomic.Bool

	activeMu          sync.Mutex
	activeTask        *workflowsv1.Task
	cancelActiveLease context.CancelFunc
}

func NewPollingTaskRunner(service TaskService, clusterSlug string, executor TaskExecutor, logger *slog.Logger) *PollingTaskRunner {
	if logger == nil {
		logger = slog.Default()
	}
	runner := &PollingTaskRunner{
		service:     service,
		clusterSlug: clusterSlug,
		logger:      logger,
		executor:    executor,
	}
	runner.requestNewTasks.Store(true)
	return runner
}

func (r *PollingTaskRunner) RunForever(ctx context.Context) error {
	return r.run(ctx, false)
}

func (r *PollingTaskRunner) RunAll(ctx context.Context) error {
	return r.run(ctx, true)
}

func (r *PollingTaskRunner) StopRequestingNewTasks() {
	r.requestNewTasks.Store(false)
}

func (r *PollingTaskRunner) IsRequestingTasks() bool {
	return r.requestNewTasks.Load()
}

func (r *PollingTaskRunner) HasActiveTask() bool {
	r.activeMu.Lock()
	defer r.activeMu.Unlock()
	return r.activeTask != nil
}

func (r *PollingTaskRunner) InterruptActiveTask(ctx context.Context) error {
	r.activeMu.Lock()
	activeTask := r.activeTask
	cancelLease := r.cancelActiveLease
	if cancelLease != nil {
		cancelLease()
	}
	r.activeMu.Unlock()

	if activeTask == nil || activeTask.GetId() == nil || activeTask.GetId().AsUUID() == uuid.Nil {
		return nil
	}
	return r.reportTaskFailed(ctx, activeTask.GetId().AsUUID(), activeTask.GetDisplay(), false, nil)
}

func (r *PollingTaskRunner) run(ctx context.Context, stopWhenIdling bool) error {
	var work *workflowsv1.NextTaskResponse
	var lastExecutionResult *workflowsv1.ExecuteTaskResponse

	for {
		if err := ctx.Err(); err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return nil
			}
			return err
		}

		// check if we need to report a task failure:
		if lastExecutionResult != nil && lastExecutionResult.GetFailedTask() != nil {
			// the last task execution resulted in a task failure, so let's report that back to the orchestrator
			pendingFailed := lastExecutionResult.GetFailedTask()
			_, err := r.service.TaskFailed(ctx, pendingFailed.GetTaskId().AsUUID(), pendingFailed.GetDisplay(), pendingFailed.GetWasWorkflowError(), pendingFailed.GetProgressUpdates())
			if err != nil {
				r.logError(ctx, err, "failed to report task failure", slog.String("task_id", pendingFailed.GetTaskId().AsUUID().String()), slog.String("task_display", pendingFailed.GetDisplay()))
				// a failure to report a task as failed typically means something is wrong with our API, so let's do one idle time before continuing with the normal runner loop
				if r.idle(ctx, randomFallbackIdleDuration()) {
					return nil
				}
				// we did not clear lastExecutionResult yet, so the next loop will attempt to report the task as failed again
				continue
			}
			lastExecutionResult = nil
			work = nil
		}

		// we don't have a current task to work on, so let's request one
		if work == nil || work.GetNextTask() == nil {
			// while requesting a new task, in the same request we also report the result of the previous task if there is one
			var lastComputedTask *workflowsv1.ComputedTask
			if lastExecutionResult.GetComputedTask() != nil {
				lastComputedTask = lastExecutionResult.GetComputedTask()
			}

			// check whether the runner should still request a new task
			requestingNewTasks := r.requestNewTasks.Load()
			// and get a list of all known task identifiers we are requesting work for
			taskIdentifiersCapableOfRunning := r.executor.TaskIdentifiers()

			if lastComputedTask == nil && len(taskIdentifiersCapableOfRunning) == 0 {
				// no task result to report back, and no tasks we can run, so no need to even send a request, let's just idle locally
				if stopWhenIdling {
					return nil
				}
				idleDuration := randomFallbackIdleDuration()
				if requestingNewTasks {
					r.logger.DebugContext(ctx, "not requesting any work, idling", slog.Duration("duration", idleDuration))
				} else {
					r.logger.DebugContext(ctx, "not requesting any work, runner is about to shut down", slog.Duration("duration", idleDuration))
				}

				if r.idle(ctx, idleDuration) {
					return nil
				}
				continue
			}

			var nextTaskToRun *workflowsv1.NextTaskToRun
			if requestingNewTasks && len(taskIdentifiersCapableOfRunning) > 0 {
				nextTaskToRun = workflowsv1.NextTaskToRun_builder{ClusterSlug: r.clusterSlug, Identifiers: taskIdentifiersCapableOfRunning}.Build()
			}

			// now let's log a debug message of what exactly we're doing with the NextTask request, depending on whether we have a lastComputedTask, and if we're a requesting a next task or not

			switch {
			case nextTaskToRun != nil && lastComputedTask != nil:
				// we report a task result, and immediately request a new task
				r.logger.DebugContext(ctx, "reporting task result and requesting new task", slog.String("computed_task_id", lastComputedTask.GetId().AsUUID().String()), slog.String("computed_task_display", lastComputedTask.GetDisplay()), slog.Int("known_task_identifiers", len(taskIdentifiersCapableOfRunning)))
			case nextTaskToRun != nil && lastComputedTask == nil:
				r.logger.DebugContext(ctx, "requesting a task to run", slog.Int("known_task_identifiers", len(taskIdentifiersCapableOfRunning)))
			case nextTaskToRun == nil && lastComputedTask != nil:
				// no next task to run, but a last computed task means we just report task results
				r.logger.DebugContext(ctx, "reporting task computed result", slog.String("task_id", lastComputedTask.GetId().AsUUID().String()), slog.String("task_display", lastComputedTask.GetDisplay()))
			default:
				r.logger.WarnContext(ctx, "unexpected NextTask request, when both computed_task and next_task is not set")
			}

			taskResponse, err := r.service.NextTask(ctx, lastComputedTask, nextTaskToRun)
			if err != nil {
				r.logError(ctx, err, "failed to request next task")
				if r.idle(ctx, randomFallbackIdleDuration()) {
					return nil
				}
				continue
			}
			lastExecutionResult = nil
			work = taskResponse
		}

		if work != nil && work.GetNextTask() != nil {
			task := work.GetNextTask()
			work = nil
			if isEmpty(task.GetId()) {
				r.logError(ctx, nil, "got a task without an ID - skipping to the next task")
				continue
			}
			if task.GetRetryCount() > 0 {
				r.logger.DebugContext(ctx, "retrying task", slog.String("task_id", task.GetId().AsUUID().String()), slog.String("task_display", task.GetDisplay()), slog.Int64("retry_count", task.GetRetryCount()))
			} else {
				r.logger.DebugContext(ctx, "executing task", slog.String("task_id", task.GetId().AsUUID().String()), slog.String("task_display", task.GetDisplay()))
			}

			response, err := r.executeTask(ctx, task)
			if err != nil {
				lastExecutionResult = workflowsv1.ExecuteTaskResponse_builder{FailedTask: r.failedTaskFromError(task, err)}.Build()
				continue
			}
			if response.GetFailedTask() == nil && response.GetComputedTask() == nil {
				lastExecutionResult = workflowsv1.ExecuteTaskResponse_builder{FailedTask: r.failedTaskFromError(task, errors.New("executor returned neither computed nor failed task"))}.Build()
				continue
			}
			lastExecutionResult = response
			continue // don't idle, immediately return the result
		}

		idleDuration := randomFallbackIdleDuration()
		if work != nil && work.GetIdling().GetSuggestedIdlingDuration() != nil {
			idleDuration = lo.Clamp(work.GetIdling().GetSuggestedIdlingDuration().AsDuration(), minIdlingDuration, maxIdlingDuration)
		}

		r.logger.DebugContext(ctx, "no task to run, idling", slog.Duration("duration", idleDuration))
		if stopWhenIdling {
			return nil
		}

		if r.idle(ctx, idleDuration) {
			return nil
		}
	}
}

func (r *PollingTaskRunner) executeTask(ctx context.Context, task *workflowsv1.Task) (*workflowsv1.ExecuteTaskResponse, error) {
	leaseCtx, stopLeaseExtensions := context.WithCancel(ctx)
	r.activeMu.Lock()
	r.activeTask = task
	r.cancelActiveLease = stopLeaseExtensions
	r.activeMu.Unlock()

	go r.extendTaskLease(
		leaseCtx,
		task.GetId().AsUUID(),
		task.GetLease().GetLease().AsDuration(),
		task.GetLease().GetRecommendedWaitUntilNextExtension().AsDuration(),
	)

	defer func() {
		stopLeaseExtensions()
		r.activeMu.Lock()
		if r.activeTask == task {
			r.activeTask = nil
			r.cancelActiveLease = nil
		}
		r.activeMu.Unlock()
	}()

	return r.executor.ExecuteTask(ctx, task)
}

func (r *PollingTaskRunner) failedTaskFromError(task *workflowsv1.Task, taskError error) *workflowsv1.TaskFailedRequest {
	var executionError *TaskExecutionError
	progressUpdates := []*workflowsv1.Progress(nil)
	if errors.As(taskError, &executionError) {
		progressUpdates = executionError.ProgressUpdates
	}
	return workflowsv1.TaskFailedRequest_builder{
		TaskId:           task.GetId(),
		Display:          failedTaskDisplay(task.GetDisplay(), taskError),
		WasWorkflowError: false,
		ProgressUpdates:  progressUpdates,
	}.Build()
}

func (r *PollingTaskRunner) reportTaskFailed(ctx context.Context, taskID uuid.UUID, display string, wasWorkflowError bool, progressUpdates []*workflowsv1.Progress) error {
	_, err := r.service.TaskFailed(ctx, taskID, display, wasWorkflowError, progressUpdates)
	return err
}

func (r *PollingTaskRunner) extendTaskLease(ctx context.Context, taskID uuid.UUID, initialLease, initialWait time.Duration) {
	wait := initialWait
	lease := initialLease
	for {
		timer := time.NewTimer(wait)
		select {
		case <-ctx.Done():
			timer.Stop()
			return
		case <-timer.C:
		}
		r.logger.DebugContext(ctx, "extending task lease", slog.String("task_id", taskID.String()), slog.Duration("lease", lease), slog.Duration("wait", wait))
		extension, err := r.service.ExtendTaskLease(ctx, taskID, 2*lease)
		if err != nil {
			r.logError(ctx, err, "failed to extend task lease", slog.String("task_id", taskID.String()))
			return
		}
		if extension.GetLease() == nil {
			r.logger.DebugContext(ctx, "task lease extension not granted", slog.String("task_id", taskID.String()))
			return
		}
		lease = extension.GetLease().AsDuration()
		wait = extension.GetRecommendedWaitUntilNextExtension().AsDuration()
	}
}

func failedTaskDisplay(taskDisplay string, taskError error) string {
	errorMessage := taskError.Error()
	if len(errorMessage) > 1024 {
		errorMessage = errorMessage[:1024]
	}
	if errorMessage == "" {
		return taskDisplay
	}
	return strings.Join([]string{taskDisplay, errorMessage}, "\n")
}

func randomFallbackIdleDuration() time.Duration {
	return fallbackPollingInterval + rand.N(fallbackJitterInterval)
}

// idle for the given duration, and return true if idling was interrupted due to a context error (e.g. cancellation)
func (r *PollingTaskRunner) idle(ctx context.Context, duration time.Duration) bool {
	timer := time.NewTimer(duration)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err() != nil
	case <-timer.C:
		return false
	}
}

func (r *PollingTaskRunner) logError(ctx context.Context, err error, msg string, args ...any) {
	switch {
	case errors.Is(err, context.Canceled):
		return
	case status.Code(err) == codes.Canceled:
		return
	}

	fields := make([]any, 0, len(args)+1)
	fields = append(fields, slog.Any("error", err))
	fields = append(fields, args...)
	r.logger.ErrorContext(ctx, msg, fields...)
}
