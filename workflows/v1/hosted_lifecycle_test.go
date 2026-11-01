package workflows

import (
	"context"
	"errors"
	"log/slog"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tileboxv1 "github.com/tilebox/tilebox-go/protogen/tilebox/v1"
	workflowsv1 "github.com/tilebox/tilebox-go/protogen/workflows/v1"
	"google.golang.org/protobuf/types/known/durationpb"
)

type hostedTestService struct {
	next   func(context.Context, *workflowsv1.ComputedTask, *workflowsv1.NextTaskToRun) (*workflowsv1.NextTaskResponse, error)
	failed func() error
	lease  func() (*workflowsv1.TaskLease, error)
}

func (s *hostedTestService) NextTask(ctx context.Context, computed *workflowsv1.ComputedTask, next *workflowsv1.NextTaskToRun) (*workflowsv1.NextTaskResponse, error) {
	return s.next(ctx, computed, next)
}

func (s *hostedTestService) TaskFailed(context.Context, uuid.UUID, string, bool, []*workflowsv1.Progress) (*workflowsv1.TaskStateResponse, error) {
	return &workflowsv1.TaskStateResponse{}, s.failed()
}

func (s *hostedTestService) ExtendTaskLease(context.Context, uuid.UUID, time.Duration) (*workflowsv1.TaskLease, error) {
	if s.lease != nil {
		return s.lease()
	}
	return workflowsv1.TaskLease_builder{Lease: durationpb.New(time.Hour), RecommendedWaitUntilNextExtension: durationpb.New(time.Hour)}.Build(), nil
}

type hostedTestExecutor struct {
	execute func(context.Context, *workflowsv1.Task) (*workflowsv1.ExecuteTaskResponse, error)
}

func (*hostedTestExecutor) TaskIdentifiers() []*workflowsv1.TaskIdentifier {
	return []*workflowsv1.TaskIdentifier{workflowsv1.TaskIdentifier_builder{Name: "test", Version: "v1"}.Build()}
}

func (e *hostedTestExecutor) ExecuteTask(ctx context.Context, task *workflowsv1.Task) (*workflowsv1.ExecuteTaskResponse, error) {
	return e.execute(ctx, task)
}

func hostedTestTask() *workflowsv1.Task {
	return workflowsv1.Task_builder{Id: tileboxv1.NewUUID(uuid.New()), Lease: workflowsv1.TaskLease_builder{
		Lease: durationpb.New(time.Hour), RecommendedWaitUntilNextExtension: durationpb.New(time.Second),
	}.Build()}.Build()
}

func TestHostedDrainRacesClaimExecutionAndReport(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		lifecycle := NewHostedLifecycle()
		claimStarted, releaseClaim := make(chan struct{}), make(chan struct{})
		taskStarted, releaseTask := make(chan struct{}), make(chan struct{})
		reportStarted, releaseReport := make(chan struct{}), make(chan struct{})
		task := hostedTestTask()
		service := &hostedTestService{
			next: func(_ context.Context, computed *workflowsv1.ComputedTask, next *workflowsv1.NextTaskToRun) (*workflowsv1.NextTaskResponse, error) {
				if computed != nil {
					assert.Nil(t, next, "drain must report without another claim")
					close(reportStarted)
					<-releaseReport
					return &workflowsv1.NextTaskResponse{}, nil
				}
				close(claimStarted)
				<-releaseClaim
				return workflowsv1.NextTaskResponse_builder{NextTask: task}.Build(), nil
			},
			lease: func() (*workflowsv1.TaskLease, error) { return nil, errors.New("lease backend unavailable") },
		}
		executor := &hostedTestExecutor{execute: func(ctx context.Context, task *workflowsv1.Task) (*workflowsv1.ExecuteTaskResponse, error) {
			close(taskStarted)
			<-releaseTask
			assert.NoError(t, ctx.Err(), "drain and lease failure must not cancel execution") //nolint:testifylint // Executor runs outside the test goroutine; require would abort it.
			return workflowsv1.ExecuteTaskResponse_builder{ComputedTask: workflowsv1.ComputedTask_builder{Id: task.GetId()}.Build()}.Build(), nil
		}}
		runner := NewPollingTaskRunner(service, "default", executor, slog.New(slog.DiscardHandler), WithHostedLifecycle(lifecycle))
		done := make(chan error, 1)
		go func() { done <- runner.RunForever(t.Context()) }()
		synctest.Wait()
		assert.False(t, runner.IsRequestingTasks())
		select {
		case <-claimStarted:
			t.Fatal("polled before StartIntake")
		default:
		}
		assert.False(t, lifecycle.StartIntake(), "initialization is required")
		lifecycle.MarkReady()
		require.True(t, lifecycle.StartIntake())
		<-claimStarted
		require.EqualValues(t, 1, lifecycle.Snapshot().GetActiveClaims())
		lifecycle.Drain()
		assert.False(t, lifecycle.StartIntake(), "drain is irreversible")
		assert.False(t, lifecycle.Snapshot().GetSafeToStop())
		close(releaseClaim)
		<-taskStarted
		time.Sleep(6 * time.Hour)
		assert.True(t, lifecycle.Snapshot().GetTaskLeaseProblem())
		assert.EqualValues(t, 1, lifecycle.Snapshot().GetActiveTasks())
		assert.False(t, lifecycle.Snapshot().GetIdle())
		close(releaseTask)
		<-reportStarted
		assert.False(t, runner.HasActiveTask())
		assert.EqualValues(t, 1, lifecycle.Snapshot().GetPendingReports())
		assert.False(t, lifecycle.Snapshot().GetSafeToStop())
		close(releaseReport)
		require.NoError(t, <-done)
		assert.True(t, lifecycle.Snapshot().GetSafeToStop())
	})
}

func TestHostedLostClaimRemainsUncertainAfterSuccessfulPoll(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		lifecycle := NewHostedLifecycle()
		calls := 0
		service := &hostedTestService{next: func(context.Context, *workflowsv1.ComputedTask, *workflowsv1.NextTaskToRun) (*workflowsv1.NextTaskResponse, error) {
			calls++
			if calls == 1 {
				return nil, errors.New("claim committed but reply lost")
			}
			lifecycle.Drain()
			return &workflowsv1.NextTaskResponse{}, nil
		}}
		runner := NewPollingTaskRunner(service, "default", &hostedTestExecutor{}, slog.New(slog.DiscardHandler), WithHostedLifecycle(lifecycle))
		lifecycle.MarkReady()
		require.True(t, lifecycle.StartIntake())
		require.NoError(t, runner.RunAll(t.Context()))
		assert.Equal(t, 2, calls)
		assert.True(t, lifecycle.Snapshot().GetClaimUncertain())
		assert.False(t, lifecycle.Snapshot().GetIdle())
		assert.False(t, lifecycle.Snapshot().GetSafeToStop())
		assert.Equal(t, workflowsv1.HostedRunnerBackendState_HOSTED_RUNNER_BACKEND_STATE_REACHABLE, lifecycle.Snapshot().GetBackend())
	})
}

func TestHostedExhaustedReportSurvivesResetAndLaterAcknowledgment(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		lifecycle := NewHostedLifecycle()
		claims, reports := 0, 0
		service := &hostedTestService{
			next: func(context.Context, *workflowsv1.ComputedTask, *workflowsv1.NextTaskToRun) (*workflowsv1.NextTaskResponse, error) {
				claims++
				return workflowsv1.NextTaskResponse_builder{NextTask: hostedTestTask()}.Build(), nil
			},
			failed: func() error {
				reports++
				if reports <= maxTaskFailedRetries {
					return errors.New("report acknowledgment lost")
				}
				lifecycle.Drain()
				return nil
			},
		}
		executor := &hostedTestExecutor{execute: func(context.Context, *workflowsv1.Task) (*workflowsv1.ExecuteTaskResponse, error) {
			return nil, errors.New("runtime failed")
		}}
		runner := NewPollingTaskRunner(service, "default", executor, slog.New(slog.DiscardHandler), WithHostedLifecycle(lifecycle))
		lifecycle.MarkReady()
		require.True(t, lifecycle.StartIntake())
		require.NoError(t, runner.RunForever(t.Context()))
		assert.Equal(t, 2, claims)
		assert.Equal(t, 4, reports)
		assert.True(t, lifecycle.Snapshot().GetReportUncertain())
		assert.EqualValues(t, 1, lifecycle.Snapshot().GetPendingReports())
		assert.False(t, lifecycle.Snapshot().GetSafeToStop())
	})
}

func TestHostedReportOnlyLostAcknowledgmentClearsOnlyOnAcknowledgment(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		lifecycle := NewHostedLifecycle()
		reports := 0
		service := &hostedTestService{next: func(_ context.Context, computed *workflowsv1.ComputedTask, next *workflowsv1.NextTaskToRun) (*workflowsv1.NextTaskResponse, error) {
			if computed == nil {
				return workflowsv1.NextTaskResponse_builder{NextTask: hostedTestTask()}.Build(), nil
			}
			assert.Nil(t, next)
			reports++
			if reports == 1 {
				return nil, errors.New("lost report ack")
			}
			assert.True(t, lifecycle.Snapshot().GetReportUncertain())
			assert.False(t, lifecycle.Snapshot().GetSafeToStop())
			return &workflowsv1.NextTaskResponse{}, nil
		}}
		executor := &hostedTestExecutor{execute: func(_ context.Context, task *workflowsv1.Task) (*workflowsv1.ExecuteTaskResponse, error) {
			lifecycle.Drain()
			return workflowsv1.ExecuteTaskResponse_builder{ComputedTask: workflowsv1.ComputedTask_builder{Id: task.GetId()}.Build()}.Build(), nil
		}}
		runner := NewPollingTaskRunner(service, "default", executor, slog.New(slog.DiscardHandler), WithHostedLifecycle(lifecycle))
		lifecycle.MarkReady()
		require.True(t, lifecycle.StartIntake())
		require.NoError(t, runner.RunForever(t.Context()))
		assert.Equal(t, 2, reports)
		assert.True(t, lifecycle.Snapshot().GetSafeToStop())
		assert.False(t, lifecycle.Snapshot().GetClaimUncertain())
	})
}

func TestHostedAggregateClaimsRaceDrain(t *testing.T) {
	for range 50 {
		lifecycle := NewHostedLifecycle()
		var calls atomic.Int32
		service := &hostedTestService{next: func(context.Context, *workflowsv1.ComputedTask, *workflowsv1.NextTaskToRun) (*workflowsv1.NextTaskResponse, error) {
			calls.Add(1)
			return &workflowsv1.NextTaskResponse{}, nil
		}}
		runners := make([]*PollingTaskRunner, 4)
		for i := range runners {
			runners[i] = NewPollingTaskRunner(service, "default", &hostedTestExecutor{}, nil, WithHostedLifecycle(lifecycle))
		}
		lifecycle.MarkReady()
		require.True(t, lifecycle.StartIntake())
		done := make(chan error, len(runners))
		for _, runner := range runners {
			go func() { done <- runner.RunAll(t.Context()) }()
		}
		lifecycle.Drain()
		for range runners {
			require.NoError(t, <-done)
		}
		assert.True(t, lifecycle.Snapshot().GetSafeToStop())
		assert.Zero(t, lifecycle.Snapshot().GetActivityGeneration(), "empty polls must not reset idle generation")
		assert.LessOrEqual(t, calls.Load(), int32(4))
	}
}
