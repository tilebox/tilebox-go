package workflows

import (
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tileboxv1 "github.com/tilebox/tilebox-go/protogen/tilebox/v1"
	workflowsv1 "github.com/tilebox/tilebox-go/protogen/workflows/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"
)

func TestPollingTaskRunnerReportsExecutorReturnedFailedTaskAsIs(t *testing.T) {
	taskID := uuid.New()
	taskDisplay := "Python task"
	mockNextTask := workflowsv1.Task_builder{
		Id: tileboxv1.NewUUID(taskID),
		Identifier: workflowsv1.TaskIdentifier_builder{
			Name:    "python.Task",
			Version: "v1.0",
		}.Build(),
		State:   workflowsv1.TaskState_TASK_STATE_RUNNING,
		Display: &taskDisplay,
		Lease: workflowsv1.TaskLease_builder{
			Lease:                             durationpb.New(5 * time.Minute),
			RecommendedWaitUntilNextExtension: durationpb.New(5 * time.Minute),
		}.Build(),
	}.Build()
	progressUpdates := []*workflowsv1.Progress{
		workflowsv1.Progress_builder{Label: "work", Total: 10, Done: 4}.Build(),
	}
	service := &mockMinimalTaskService{nextTask: mockNextTask}
	executor := &failedResponseExecutor{
		response: workflowsv1.ExecuteTaskResponse_builder{
			FailedTask: workflowsv1.TaskFailedRequest_builder{
				TaskId:           tileboxv1.NewUUID(taskID),
				Display:          "Python classified failure",
				WasWorkflowError: true,
				ProgressUpdates:  progressUpdates,
			}.Build(),
		}.Build(),
	}
	runner := NewPollingTaskRunner(service, "default", executor, slog.Default())

	require.NoError(t, runner.RunAll(context.Background()))

	require.True(t, service.failed)
	require.Equal(t, "Python classified failure", service.failedDisplay)
	require.True(t, service.failedWasWorkflowError)
	require.Equal(t, progressUpdates, service.failedProgressUpdates)
}

func TestPollingTaskRunnerReportsComputedTaskRequestErrorAsWorkflowFailure(t *testing.T) {
	taskID := uuid.New()
	taskDisplay := "Python task"
	mockNextTask := workflowsv1.Task_builder{
		Id: tileboxv1.NewUUID(taskID),
		Identifier: workflowsv1.TaskIdentifier_builder{
			Name:    "python.Task",
			Version: "v1.0",
		}.Build(),
		State:   workflowsv1.TaskState_TASK_STATE_RUNNING,
		Display: &taskDisplay,
		Lease: workflowsv1.TaskLease_builder{
			Lease:                             durationpb.New(5 * time.Minute),
			RecommendedWaitUntilNextExtension: durationpb.New(5 * time.Minute),
		}.Build(),
	}.Build()
	service := &requestErrorTaskService{
		nextTask:                mockNextTask,
		nextTaskComputedTaskErr: connect.NewError(connect.CodeInvalidArgument, errors.New("invalid computed task")),
	}
	executor := &failedResponseExecutor{
		response: workflowsv1.ExecuteTaskResponse_builder{
			ComputedTask: workflowsv1.ComputedTask_builder{
				Id:      tileboxv1.NewUUID(taskID),
				Display: taskDisplay,
			}.Build(),
		}.Build(),
	}
	runner := NewPollingTaskRunner(service, "default", executor, slog.Default())

	require.NoError(t, runner.RunAll(context.Background()))

	require.Len(t, service.taskFailedRequests, 1)
	failedTask := service.taskFailedRequests[0]
	assert.Equal(t, taskID, failedTask.GetTaskId().AsUUID())
	assert.True(t, failedTask.GetWasWorkflowError())
	assert.Contains(t, failedTask.GetDisplay(), taskDisplay)
	assert.Contains(t, failedTask.GetDisplay(), "invalid computed task")
	assert.Empty(t, failedTask.GetProgressUpdates())
}

func TestPollingTaskRunnerDoesNotFailComputedTaskWhenCombinedNextTaskRequestIsInvalid(t *testing.T) {
	taskID := uuid.New()
	taskDisplay := "Python task"
	mockNextTask := workflowsv1.Task_builder{
		Id: tileboxv1.NewUUID(taskID),
		Identifier: workflowsv1.TaskIdentifier_builder{
			Name:    "python.Task",
			Version: "v1.0",
		}.Build(),
		State:   workflowsv1.TaskState_TASK_STATE_RUNNING,
		Display: &taskDisplay,
		Lease: workflowsv1.TaskLease_builder{
			Lease:                             durationpb.New(5 * time.Minute),
			RecommendedWaitUntilNextExtension: durationpb.New(5 * time.Minute),
		}.Build(),
	}.Build()
	service := &requestErrorTaskService{
		nextTask:                      mockNextTask,
		nextTaskComputedAndRequestErr: connect.NewError(connect.CodeInvalidArgument, errors.New("invalid work request")),
	}
	executor := &failedResponseExecutor{
		response: workflowsv1.ExecuteTaskResponse_builder{
			ComputedTask: workflowsv1.ComputedTask_builder{
				Id:      tileboxv1.NewUUID(taskID),
				Display: taskDisplay,
			}.Build(),
		}.Build(),
	}
	runner := NewPollingTaskRunner(service, "default", executor, slog.Default())

	require.NoError(t, runner.RunAll(context.Background()))

	assert.Empty(t, service.taskFailedRequests)
	require.Len(t, service.computedTasks, 1)
	assert.Equal(t, taskID, service.computedTasks[0].GetId().AsUUID())
}

func TestPollingTaskRunnerRetriesTaskFailedRequestErrorOnceAsWorkflowFailure(t *testing.T) {
	taskID := uuid.New()
	taskDisplay := "Python task"
	mockNextTask := workflowsv1.Task_builder{
		Id: tileboxv1.NewUUID(taskID),
		Identifier: workflowsv1.TaskIdentifier_builder{
			Name:    "python.Task",
			Version: "v1.0",
		}.Build(),
		State:   workflowsv1.TaskState_TASK_STATE_RUNNING,
		Display: &taskDisplay,
		Lease: workflowsv1.TaskLease_builder{
			Lease:                             durationpb.New(5 * time.Minute),
			RecommendedWaitUntilNextExtension: durationpb.New(5 * time.Minute),
		}.Build(),
	}.Build()
	service := &requestErrorTaskService{
		nextTask: mockNextTask,
		taskFailedErrors: []error{
			connect.NewError(connect.CodeInvalidArgument, errors.New("invalid failed task request")),
		},
	}
	progressUpdates := []*workflowsv1.Progress{
		workflowsv1.Progress_builder{Label: "work", Total: 10, Done: 4}.Build(),
	}
	executor := &failedResponseExecutor{
		response: workflowsv1.ExecuteTaskResponse_builder{
			FailedTask: workflowsv1.TaskFailedRequest_builder{
				TaskId:           tileboxv1.NewUUID(taskID),
				Display:          "user defined display",
				WasWorkflowError: false,
				ProgressUpdates:  progressUpdates,
			}.Build(),
		}.Build(),
	}
	runner := NewPollingTaskRunner(service, "default", executor, slog.Default())

	require.NoError(t, runner.RunAll(context.Background()))

	require.Len(t, service.taskFailedRequests, 2)
	assert.Equal(t, "user defined display", service.taskFailedRequests[0].GetDisplay())
	assert.False(t, service.taskFailedRequests[0].GetWasWorkflowError())
	assert.Equal(t, progressUpdates, service.taskFailedRequests[0].GetProgressUpdates())
	assert.NotContains(t, service.taskFailedRequests[1].GetDisplay(), "user defined display")
	assert.Contains(t, service.taskFailedRequests[1].GetDisplay(), "invalid failed task request")
	assert.True(t, service.taskFailedRequests[1].GetWasWorkflowError())
	assert.Empty(t, service.taskFailedRequests[1].GetProgressUpdates())
}

func TestPollingTaskRunnerStopsRetryingTaskFailedAfterSecondRequestError(t *testing.T) {
	taskID := uuid.New()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	mockNextTask := workflowsv1.Task_builder{
		Id: tileboxv1.NewUUID(taskID),
		Identifier: workflowsv1.TaskIdentifier_builder{
			Name:    "python.Task",
			Version: "v1.0",
		}.Build(),
		State: workflowsv1.TaskState_TASK_STATE_RUNNING,
		Lease: workflowsv1.TaskLease_builder{
			Lease:                             durationpb.New(5 * time.Minute),
			RecommendedWaitUntilNextExtension: durationpb.New(5 * time.Minute),
		}.Build(),
	}.Build()
	service := &requestErrorTaskService{
		nextTask: mockNextTask,
		taskFailedErrors: []error{
			connect.NewError(connect.CodeInvalidArgument, errors.New("invalid failed task request")),
			connect.NewError(connect.CodeInvalidArgument, errors.New("still invalid")),
		},
		taskFailedHook: func(requests int) {
			if requests == 2 {
				cancel()
			}
		},
	}
	executor := &failedResponseExecutor{
		response: workflowsv1.ExecuteTaskResponse_builder{
			FailedTask: workflowsv1.TaskFailedRequest_builder{
				TaskId:  tileboxv1.NewUUID(taskID),
				Display: "user defined display",
			}.Build(),
		}.Build(),
	}
	runner := NewPollingTaskRunner(service, "default", executor, slog.Default())

	require.NoError(t, runner.RunAll(ctx))

	require.Len(t, service.taskFailedRequests, 2)
	assert.Equal(t, "user defined display", service.taskFailedRequests[0].GetDisplay())
	assert.NotContains(t, service.taskFailedRequests[1].GetDisplay(), "user defined display")
}

func TestPollingTaskRunnerResetsTaskFailedAfterMaxRetryableErrors(t *testing.T) {
	taskID := uuid.New()
	service := &requestErrorTaskService{
		taskFailedErrors: []error{
			connect.NewError(connect.CodeUnavailable, errors.New("unavailable 1")),
			connect.NewError(connect.CodeUnavailable, errors.New("unavailable 2")),
			connect.NewError(connect.CodeUnavailable, errors.New("unavailable 3")),
		},
	}
	runner := NewPollingTaskRunner(service, "default", &failedResponseExecutor{}, slog.Default())
	pending := pendingReport{
		result: workflowsv1.ExecuteTaskResponse_builder{
			FailedTask: workflowsv1.TaskFailedRequest_builder{
				TaskId:  tileboxv1.NewUUID(taskID),
				Display: "task failed",
			}.Build(),
		}.Build(),
	}

	var outcome reportOutcome
	pending, outcome = runner.reportPendingFailure(context.Background(), context.Background(), pending)
	assert.Equal(t, reportRetryLater, outcome)
	pending, outcome = runner.reportPendingFailure(context.Background(), context.Background(), pending)
	assert.Equal(t, reportRetryLater, outcome)
	_, outcome = runner.reportPendingFailure(context.Background(), context.Background(), pending)
	assert.Equal(t, reportResetRunner, outcome)
	require.Len(t, service.taskFailedRequests, maxTaskFailedRetries)
}

type failedResponseExecutor struct {
	response *workflowsv1.ExecuteTaskResponse
}

func (e *failedResponseExecutor) TaskIdentifiers() []*workflowsv1.TaskIdentifier {
	return []*workflowsv1.TaskIdentifier{
		workflowsv1.TaskIdentifier_builder{Name: "python.Task", Version: "v1.0"}.Build(),
	}
}

func (e *failedResponseExecutor) ExecuteTask(context.Context, *workflowsv1.Task) (*workflowsv1.ExecuteTaskResponse, error) {
	return e.response, nil
}

type requestErrorTaskService struct {
	computedTasks                 []*workflowsv1.ComputedTask
	nextTask                      *workflowsv1.Task
	nextTaskComputedAndRequestErr error
	nextTaskComputedTaskErr       error
	taskFailedErrors              []error
	taskFailedRequests            []*workflowsv1.TaskFailedRequest
	taskFailedHook                func(requests int)
}

var _ TaskService = &requestErrorTaskService{}

func (s *requestErrorTaskService) NextTask(_ context.Context, computedTask *workflowsv1.ComputedTask, nextTaskToRun *workflowsv1.NextTaskToRun) (*workflowsv1.NextTaskResponse, error) {
	if computedTask != nil && nextTaskToRun != nil && s.nextTaskComputedAndRequestErr != nil {
		return nil, s.nextTaskComputedAndRequestErr
	}
	if computedTask != nil && s.nextTaskComputedTaskErr != nil {
		return nil, s.nextTaskComputedTaskErr
	}
	if computedTask != nil {
		s.computedTasks = append(s.computedTasks, computedTask)
	}
	if s.nextTask != nil {
		response := workflowsv1.NextTaskResponse_builder{NextTask: proto.CloneOf(s.nextTask)}.Build()
		s.nextTask = nil
		return response, nil
	}
	return workflowsv1.NextTaskResponse_builder{}.Build(), nil
}

func (s *requestErrorTaskService) TaskFailed(_ context.Context, taskID uuid.UUID, display string, wasWorkflowError bool, progressUpdates []*workflowsv1.Progress) (*workflowsv1.TaskStateResponse, error) {
	s.taskFailedRequests = append(s.taskFailedRequests, workflowsv1.TaskFailedRequest_builder{
		TaskId:           tileboxv1.NewUUID(taskID),
		Display:          display,
		WasWorkflowError: wasWorkflowError,
		ProgressUpdates:  progressUpdates,
	}.Build())
	if s.taskFailedHook != nil {
		s.taskFailedHook(len(s.taskFailedRequests))
	}
	if len(s.taskFailedErrors) > 0 {
		err := s.taskFailedErrors[0]
		s.taskFailedErrors = s.taskFailedErrors[1:]
		return nil, err
	}
	return workflowsv1.TaskStateResponse_builder{State: workflowsv1.TaskState_TASK_STATE_FAILED}.Build(), nil
}

func (s *requestErrorTaskService) ExtendTaskLease(context.Context, uuid.UUID, time.Duration) (*workflowsv1.TaskLease, error) {
	return workflowsv1.TaskLease_builder{
		Lease:                             durationpb.New(5 * time.Minute),
		RecommendedWaitUntilNextExtension: durationpb.New(5 * time.Minute),
	}.Build(), nil
}
