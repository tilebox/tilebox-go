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

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/tilebox/tilebox-go/internal/span"
	obslogger "github.com/tilebox/tilebox-go/observability/logger"
	workflowsv1 "github.com/tilebox/tilebox-go/protogen/workflows/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const maxTaskFailedRetries = 3

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
	_, err := r.service.TaskFailed(ctx, activeTask.GetId().AsUUID(), activeTask.GetDisplay(), false, nil)
	return err
}

func (r *PollingTaskRunner) run(ctx context.Context, stopWhenIdling bool) error {
	var work *workflowsv1.NextTaskResponse
	pending := pendingReport{}
	// currentTaskLogContext carries task_id and trace context for logs related to the task
	// that is currently executing or whose result is pending. When no task result is pending,
	// it is reset to the main runner context so callers never have to handle a nil context.
	currentTaskLogContext := ctx

	for {
		if err := ctx.Err(); err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return nil
			}
			return err
		}

		if pending.hasFailureToReport() {
			var outcome reportOutcome
			pending, outcome = r.reportPendingFailure(ctx, currentTaskLogContext, pending)
			switch outcome {
			case reportSucceeded:
				currentTaskLogContext = ctx //nolint:fatcontext // Reset task-scoped log state to the root runner context.
			case reportRetryNow:
				continue
			case reportRetryLater:
				if r.idle(ctx, randomFallbackIdleDuration()) {
					return nil
				}
				continue
			case reportResetRunner:
				work = nil
				pending = pendingReport{}
				currentTaskLogContext = ctx
				if r.idle(ctx, randomFallbackIdleDuration()) {
					return nil
				}
				continue
			}
		}

		// we don't have a current task to work on, so let's request one
		if work == nil || work.GetNextTask() == nil {
			// while requesting a new task, in the same request we also report the result of the previous task if there is one
			var lastComputedTask *workflowsv1.ComputedTask
			if pending.hasComputedTaskToReport() {
				lastComputedTask = pending.result.GetComputedTask()
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

			if lastComputedTask != nil {
				var outcome reportOutcome
				work, pending, outcome = r.reportPendingComputed(ctx, currentTaskLogContext, pending, nextTaskToRun)
				switch outcome {
				case reportSucceeded:
					currentTaskLogContext = ctx
				case reportRetryNow:
					continue
				case reportRetryLater:
					if r.idle(ctx, randomFallbackIdleDuration()) {
						return nil
					}
					continue
				case reportResetRunner:
					// during a next task request we should never get that, but just to be defensive in case of future
					// changes
					work = nil
					pending = pendingReport{}
					currentTaskLogContext = ctx
					if r.idle(ctx, randomFallbackIdleDuration()) {
						return nil
					}
					continue
				}
			} else {
				// work stealing:
				taskResponse, err := r.service.NextTask(ctx, nil, nextTaskToRun)
				if err != nil {
					// easy retry case, since we don't have a task result to report anyway
					logError(r.logger, ctx, err, "failed to request next task, will retry")
					if r.idle(ctx, randomFallbackIdleDuration()) {
						return nil
					}
					continue
				}
				work = taskResponse
			}
		}

		if work != nil && work.GetNextTask() != nil {
			task := work.GetNextTask()
			work = nil
			if isEmpty(task.GetId()) {
				logError(r.logger, ctx, nil, "got a task without an ID - skipping to the next task")
				continue
			}
			currentTaskLogContext = taskContextForLogs(ctx, task)
			if task.GetRetryCount() > 0 {
				r.logger.DebugContext(currentTaskLogContext, "retrying task", slog.String("task_display", task.GetDisplay()), slog.Int64("retry_count", task.GetRetryCount()))
			} else {
				r.logger.DebugContext(currentTaskLogContext, "executing task", slog.String("task_display", task.GetDisplay()))
			}

			response, err := r.executeTask(ctx, task)
			if err != nil {
				pending = pendingReport{
					result: workflowsv1.ExecuteTaskResponse_builder{FailedTask: r.failedTaskFromError(task, err)}.Build(),
				}
				continue
			}
			if response.GetFailedTask() == nil && response.GetComputedTask() == nil {
				pending = pendingReport{
					result: workflowsv1.ExecuteTaskResponse_builder{FailedTask: r.failedTaskFromError(task, errors.New("executor returned neither computed nor failed task"))}.Build(),
				}
				continue
			}
			pending = pendingReport{
				result: response,
			}
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
	progressUpdates := []*workflowsv1.Progress(nil)
	if executionError, ok := errors.AsType[*TaskExecutionError](taskError); ok {
		progressUpdates = executionError.ProgressUpdates
	}
	return workflowsv1.TaskFailedRequest_builder{
		TaskId:           task.GetId(),
		Display:          failedTaskDisplay(task.GetDisplay(), taskError),
		WasWorkflowError: false,
		ProgressUpdates:  progressUpdates,
	}.Build()
}

type reportOutcome int

const (
	reportSucceeded reportOutcome = iota
	reportRetryNow
	reportRetryLater
	reportResetRunner
)

type pendingReport struct {
	result *workflowsv1.ExecuteTaskResponse

	taskFailedAttempts   int
	taskFailedSimplified bool
}

func (p pendingReport) hasFailureToReport() bool {
	return p.result != nil && p.result.HasFailedTask()
}

func (p pendingReport) hasComputedTaskToReport() bool {
	return p.result != nil && p.result.HasComputedTask()
}

// reportPendingFailure owns the retry policy for TaskFailed RPCs. The first request-shaped error
// returns a pending report with a simplified TaskFailedRequest, so the runner never resends the same
// rejected user display/progress payload. After simplification, another request-shaped error or three
// total failed attempts resets the runner state; retryable errors below that cap are retried after idling.
func (r *PollingTaskRunner) reportPendingFailure(ctx, logCtx context.Context, pending pendingReport) (pendingReport, reportOutcome) {
	failedTask := pending.result.GetFailedTask()
	_, err := r.service.TaskFailed(ctx, failedTask.GetTaskId().AsUUID(), failedTask.GetDisplay(), failedTask.GetWasWorkflowError(), failedTask.GetProgressUpdates())
	if err == nil {
		return pendingReport{}, reportSucceeded
	}

	retryAttempt := pending.taskFailedAttempts + 1
	resetRunnerState := retryAttempt >= maxTaskFailedRetries
	simplifiedRequest := false
	failedTaskToRetry := failedTask
	taskFailedSimplified := pending.taskFailedSimplified
	if !shouldRetryRPCAfterTimeout(err) {
		if pending.taskFailedSimplified {
			resetRunnerState = true
		} else {
			failedTaskToRetry = simplifiedTaskFailedRequest(failedTask, err)
			taskFailedSimplified = true
			simplifiedRequest = true
		}
	}

	// we want to avoid modifying the input argument silently, so we construct a new result to return
	pendingRetry := pendingReport{
		result:               workflowsv1.ExecuteTaskResponse_builder{FailedTask: failedTaskToRetry}.Build(),
		taskFailedAttempts:   retryAttempt,
		taskFailedSimplified: taskFailedSimplified,
	}

	logError(r.logger, logCtx, err, "failed to report task failure back to Tilebox", slog.Int("retry_count", retryAttempt), slog.Int("max_retries", maxTaskFailedRetries), slog.Bool("request_simplified", taskFailedSimplified), slog.Bool("resetting_runner", resetRunnerState))
	if simplifiedRequest && !resetRunnerState {
		return pendingRetry, reportRetryNow
	}
	if resetRunnerState {
		return pendingRetry, reportResetRunner
	}
	return pendingRetry, reportRetryLater
}

// reportPendingComputed reports a computed task through NextTask. If the combined
// computed-result-plus-new-work request fails with a request-shaped error, it retries once without
// the new-work half to avoid blaming task output for cluster/identifier/auth request failures. Only
// payload-shape errors from an isolated computed-task report are converted into workflow failures.
func (r *PollingTaskRunner) reportPendingComputed(ctx, logCtx context.Context, pending pendingReport, nextTaskToRun *workflowsv1.NextTaskToRun) (*workflowsv1.NextTaskResponse, pendingReport, reportOutcome) {
	computedTask := pending.result.GetComputedTask()
	taskResponse, err := r.service.NextTask(ctx, computedTask, nextTaskToRun)
	if err == nil {
		return taskResponse, pendingReport{}, reportSucceeded
	}
	if shouldRetryRPCAfterTimeout(err) {
		logError(r.logger, logCtx, err, "failed to report computed task, will retry", slog.String("task_id", computedTask.GetId().AsUUID().String()), slog.String("task_display", computedTask.GetDisplay()))
		return nil, pending, reportRetryLater
	}

	if nextTaskToRun != nil {
		logError(r.logger, logCtx, err, "failed to report computed task and request next task due to request error, retrying without new work request", slog.String("task_id", computedTask.GetId().AsUUID().String()), slog.String("task_display", computedTask.GetDisplay()))
		taskResponse, err = r.service.NextTask(ctx, computedTask, nil)
		if err == nil {
			return taskResponse, pendingReport{}, reportSucceeded
		}
		if shouldRetryRPCAfterTimeout(err) || !shouldFailComputedTaskForRPC(err) {
			logError(r.logger, logCtx, err, "failed to report computed task without requesting next task, will retry", slog.String("task_id", computedTask.GetId().AsUUID().String()), slog.String("task_display", computedTask.GetDisplay()))
			return nil, pending, reportRetryLater
		}
	}

	if !shouldFailComputedTaskForRPC(err) {
		logError(r.logger, logCtx, err, "failed to report computed task due to request error, will retry", slog.String("task_id", computedTask.GetId().AsUUID().String()), slog.String("task_display", computedTask.GetDisplay()))
		return nil, pending, reportRetryLater
	}

	logError(r.logger, logCtx, err, "failed to report computed task due to invalid payload, will report task as failed", slog.String("task_id", computedTask.GetId().AsUUID().String()), slog.String("task_display", computedTask.GetDisplay()))

	return nil, pendingReport{
		result:               workflowsv1.ExecuteTaskResponse_builder{FailedTask: failedTaskFromComputedTaskRequestError(computedTask, err)}.Build(),
		taskFailedAttempts:   0,
		taskFailedSimplified: false,
	}, reportRetryNow
}

func simplifiedTaskFailedRequest(failedTask *workflowsv1.TaskFailedRequest, err error) *workflowsv1.TaskFailedRequest {
	return workflowsv1.TaskFailedRequest_builder{
		TaskId:           failedTask.GetTaskId(),
		Display:          failedTaskDisplay("", err),
		WasWorkflowError: true,
	}.Build()
}

func taskContextForLogs(ctx context.Context, task *workflowsv1.Task) context.Context {
	if task == nil {
		return ctx
	}

	if !isEmpty(task.GetId()) {
		ctx = obslogger.ContextWithSlogAttributes(ctx, slog.String("task_id", task.GetId().AsUUID().String()))
	}
	return span.ContextWithTraceParent(ctx, task.GetJob().GetTraceParent())
}

func failedTaskFromComputedTaskRequestError(computedTask *workflowsv1.ComputedTask, err error) *workflowsv1.TaskFailedRequest {
	return workflowsv1.TaskFailedRequest_builder{
		TaskId:           computedTask.GetId(),
		Display:          failedTaskDisplay(computedTask.GetDisplay(), err),
		WasWorkflowError: true,
	}.Build()
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
			logError(r.logger, ctx, err, "failed to extend task lease", slog.String("task_id", taskID.String()))
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

func shouldRetryRPCAfterTimeout(err error) bool {
	switch connect.CodeOf(err) { //nolint:exhaustive // Only request-error codes are special-cased; all other codes retry.
	case connect.CodeInvalidArgument,
		connect.CodeNotFound,
		connect.CodeAlreadyExists,
		connect.CodeFailedPrecondition,
		connect.CodeOutOfRange,
		connect.CodeUnimplemented:
		return false
	}

	switch status.Code(err) { //nolint:exhaustive // Only request-error codes are special-cased; all other codes retry.
	case codes.InvalidArgument,
		codes.NotFound,
		codes.AlreadyExists,
		codes.FailedPrecondition,
		codes.OutOfRange,
		codes.Unimplemented:
		return false
	}
	return true
}

func shouldFailComputedTaskForRPC(err error) bool {
	switch connect.CodeOf(err) { //nolint:exhaustive // Only payload-shape errors fail the computed task.
	case connect.CodeInvalidArgument, connect.CodeOutOfRange:
		return true
	}

	switch status.Code(err) { //nolint:exhaustive // Only payload-shape errors fail the computed task.
	case codes.InvalidArgument, codes.OutOfRange:
		return true
	default:
		return false
	}
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

func logError(logger *slog.Logger, ctx context.Context, err error, msg string, args ...any) {
	switch {
	case errors.Is(err, context.Canceled):
		return
	case connect.CodeOf(err) == connect.CodeCanceled:
		return
	case status.Code(err) == codes.Canceled:
		return
	}

	fields := make([]any, 0, len(args)+1)
	fields = append(fields, slog.Any("error", err))
	fields = append(fields, args...)
	logger.ErrorContext(ctx, msg, fields...)
}
