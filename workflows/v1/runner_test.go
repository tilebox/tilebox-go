package workflows

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tileboxv1 "github.com/tilebox/tilebox-go/protogen/tilebox/v1"
	workflowsv1 "github.com/tilebox/tilebox-go/protogen/workflows/v1"
	"github.com/tilebox/tilebox-go/workflows/v1/subtask"
	"go.opentelemetry.io/otel/trace/noop"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"
)

type mockClusterService struct {
	WorkflowService
}

func (m *mockClusterService) GetCluster(ctx context.Context, slug string) (*workflowsv1.Cluster, error) {
	if slug != "" {
		panic("mock only support empty slug")
	}

	return workflowsv1.Cluster_builder{
		Slug:        "default-5GGzdzBZEA3oeD",
		DisplayName: "Default",
		Deletable:   false,
	}.Build(), nil
}

type mockTaskService struct {
	TaskService
}

type testTask1 struct {
	ExecutableTask
}

type testTask2 struct {
	ExecutableTask
}

type badIdentifierTask struct {
	ExecutableTask
}

func (t *badIdentifierTask) Identifier() TaskIdentifier {
	return NewTaskIdentifier("", "")
}

type SampleTaskV1 struct {
	N int
}

func (t *SampleTaskV1) Identifier() TaskIdentifier {
	return NewTaskIdentifier("tilebox.com/sample/SampleTask", "v1.0")
}

func (t *SampleTaskV1) Execute(ctx context.Context) error {
	return nil
}

type SampleTaskV2 struct {
	N int
}

func (t *SampleTaskV2) Identifier() TaskIdentifier {
	return NewTaskIdentifier("tilebox.com/sample/SampleTask", "v2.0")
}

func (t *SampleTaskV2) Execute(ctx context.Context) error {
	return nil
}

func TestTaskRunner_RegisterTasks(t *testing.T) {
	tracer := noop.NewTracerProvider().Tracer("")
	service := mockTaskService{}
	mockClusterClient := clusterClient{service: &mockClusterService{}}

	type args struct {
		tasks []ExecutableTask
	}
	tests := []struct {
		name      string
		args      args
		wantErr   string
		wantTasks []ExecutableTask
	}{
		{
			name: "Register Tasks no tasks",
			args: args{
				tasks: []ExecutableTask{},
			},
			wantTasks: []ExecutableTask{},
		},
		{
			name: "Register Tasks duplicated implicitly generated task identifiers",
			args: args{
				tasks: []ExecutableTask{
					&testTask1{},
					&testTask1{},
				},
			},
			wantErr: "duplicate task identifier: a task 'testTask1' with version 'v0.0' is already registered",
		},
		{
			name: "Register Tasks duplicated explicitly defined task identifiers",
			args: args{
				tasks: []ExecutableTask{
					&SampleTaskV1{},
					&SampleTaskV1{},
				},
			},
			wantErr: "duplicate task identifier: a task 'tilebox.com/sample/SampleTask' with version 'v1.0' is already registered",
		},
		{
			name: "Register Tasks duplicated identifier names with different versions",
			args: args{
				tasks: []ExecutableTask{
					&SampleTaskV1{},
					&SampleTaskV2{},
				},
			},
			wantTasks: []ExecutableTask{
				&SampleTaskV1{},
				&SampleTaskV2{},
			},
		},
		{
			name: "Register Tasks bad identifier",
			args: args{
				tasks: []ExecutableTask{
					&testTask1{},
					&badIdentifierTask{},
					&testTask2{},
				},
			},
			wantErr: "task name is empty",
			wantTasks: []ExecutableTask{
				&testTask1{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runner, err := newTaskRunner(context.Background(), service, mockClusterClient, tracer)
			require.NoError(t, err)

			err = runner.RegisterTasks(tt.args.tasks...)
			if tt.wantErr != "" {
				require.ErrorContains(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)

			assert.Len(t, runner.taskDefinitions, len(tt.wantTasks))
			for _, task := range tt.wantTasks {
				identifier := identifierFromTask(task)
				_, ok := runner.GetRegisteredTask(identifier)
				assert.True(t, ok, "task not found in taskDefinitions")
			}
		})
	}
}

func Test_isEmpty(t *testing.T) {
	type args struct {
		id *tileboxv1.ID
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "isEmpty nil",
			args: args{
				id: nil,
			},
			want: true,
		},
		{
			name: "isEmpty not nil but invalid id",
			args: args{
				id: tileboxv1.ID_builder{Uuid: []byte{1, 2, 3}}.Build(),
			},
			want: true,
		},
		{
			name: "isEmpty not nil but valid id",
			args: args{
				id: tileboxv1.ID_builder{Uuid: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6}}.Build(),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isEmpty(tt.args.id)
			assert.Equal(t, tt.want, got)
		})
	}
}

// mockMinimalTaskService is a mock task service that returns a single next task response. after that it returns an empty response every time
type mockMinimalTaskService struct {
	computedTasks []*workflowsv1.ComputedTask
	nextTask      *workflowsv1.Task
	failed        bool
}

var _ TaskService = &mockMinimalTaskService{}

func (m *mockMinimalTaskService) NextTask(_ context.Context, computedTask *workflowsv1.ComputedTask, _ *workflowsv1.NextTaskToRun) (*workflowsv1.NextTaskResponse, error) {
	if computedTask != nil {
		m.computedTasks = append(m.computedTasks, computedTask)
	}

	if m.nextTask != nil {
		resp := workflowsv1.NextTaskResponse_builder{NextTask: proto.CloneOf(m.nextTask)}.Build()
		m.nextTask = nil
		return resp, nil
	}
	// subsequent calls => no next task
	return workflowsv1.NextTaskResponse_builder{
		NextTask: nil,
	}.Build(), nil
}

func (m *mockMinimalTaskService) TaskFailed(context.Context, uuid.UUID, string, bool, []*workflowsv1.Progress) (*workflowsv1.TaskStateResponse, error) {
	m.failed = true
	return workflowsv1.TaskStateResponse_builder{
		State: workflowsv1.TaskState_TASK_STATE_FAILED,
	}.Build(), nil
}

func (m *mockMinimalTaskService) ExtendTaskLease(context.Context, uuid.UUID, time.Duration) (*workflowsv1.TaskLease, error) {
	return workflowsv1.TaskLease_builder{
		Lease:                             durationpb.New(5 * time.Minute),
		RecommendedWaitUntilNextExtension: durationpb.New(5 * time.Minute),
	}.Build(), nil
}

type testTaskExecute struct{}

func (testTaskExecute) Execute(context.Context) error {
	return nil
}

func (testTaskExecute) Identifier() TaskIdentifier {
	return NewTaskIdentifier("testTaskExecute", "v0.1")
}

// simulate an actual runner with one task that it should receive
func TestTaskRunner_RunForever(t *testing.T) {
	task := &testTaskExecute{}

	taskInput, err := json.Marshal(task)
	require.NoError(t, err)

	// we simulate a response from NextTask to let our runner start our task instead of calling Execute() ourselves
	mockNextTask := workflowsv1.Task_builder{
		Id: tileboxv1.NewUUID(uuid.New()),
		Identifier: workflowsv1.TaskIdentifier_builder{
			Name:    task.Identifier().Name(),
			Version: task.Identifier().Version(),
		}.Build(),
		State: workflowsv1.TaskState_TASK_STATE_RUNNING,
		Input: taskInput,
		Job: workflowsv1.Job_builder{
			Id:          tileboxv1.NewUUID(uuid.New()),
			Name:        "tilebox-test-job",
			TraceParent: "00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01",
		}.Build(),
	}.Build()
	mockTaskClient := &mockMinimalTaskService{nextTask: mockNextTask}

	tracer := noop.NewTracerProvider().Tracer("")
	service := mockTaskService{}
	mockClusterClient := clusterClient{service: &mockClusterService{}}
	runner, err := newTaskRunner(context.Background(), service, mockClusterClient, tracer)
	require.NoError(t, err, "failed to create task runner")
	runner.service = mockTaskClient

	err = runner.RegisterTasks(task)
	require.NoError(t, err, "failed to register task")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond) // run our one task, then stop
	defer cancel()
	runner.RunForever(ctx)

	assert.False(t, mockTaskClient.failed, "task should not have failed")
	require.Len(t, mockTaskClient.computedTasks, 1)
	computedTask := mockTaskClient.computedTasks[0]
	assert.Equal(t, mockNextTask.GetId(), computedTask.GetId())
}

func TestTaskRunner_RunAll(t *testing.T) {
	task := &testTaskExecute{}

	taskInput, err := json.Marshal(task)
	require.NoError(t, err)

	// we simulate a response from NextTask to let our runner start our task instead of calling Execute() ourselves
	mockNextTask := workflowsv1.Task_builder{
		Id: tileboxv1.NewUUID(uuid.New()),
		Identifier: workflowsv1.TaskIdentifier_builder{
			Name:    task.Identifier().Name(),
			Version: task.Identifier().Version(),
		}.Build(),
		State: workflowsv1.TaskState_TASK_STATE_RUNNING,
		Input: taskInput,
		Job: workflowsv1.Job_builder{
			Id:          tileboxv1.NewUUID(uuid.New()),
			Name:        "tilebox-test-job",
			TraceParent: "00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01",
		}.Build(),
	}.Build()
	mockTaskClient := &mockMinimalTaskService{nextTask: mockNextTask}

	tracer := noop.NewTracerProvider().Tracer("")
	service := mockTaskService{}
	mockClusterClient := clusterClient{service: &mockClusterService{}}
	runner, err := newTaskRunner(context.Background(), service, mockClusterClient, tracer)
	require.NoError(t, err, "failed to create task runner")
	runner.service = mockTaskClient

	err = runner.RegisterTasks(task)
	require.NoError(t, err, "failed to register task")

	// run our one task, then stop
	runner.RunAll(context.Background())

	assert.False(t, mockTaskClient.failed, "task should not have failed")
	require.Len(t, mockTaskClient.computedTasks, 1)
	computedTask := mockTaskClient.computedTasks[0]
	assert.Equal(t, mockNextTask.GetId(), computedTask.GetId())
}

func Test_withTaskExecutionContextRoundtrip(t *testing.T) {
	tracer := noop.NewTracerProvider().Tracer("")
	service := mockTaskService{}
	mockClusterClient := clusterClient{service: &mockClusterService{}}
	runner, err := newTaskRunner(context.Background(), service, mockClusterClient, tracer)
	require.NoError(t, err)

	taskID := uuid.New()

	type args struct {
		ctx  context.Context
		task *workflowsv1.Task
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "withTaskExecutionContext",
			args: args{
				ctx: context.Background(),
				task: workflowsv1.Task_builder{
					Id: tileboxv1.ID_builder{Uuid: taskID[:]}.Build(),
				}.Build(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updatedCtx := runner.withTaskExecutionContext(tt.args.ctx, tt.args.task)
			got := getTaskExecutionContext(updatedCtx)
			assert.Empty(t, got.subtasks)
			assert.Equal(t, tt.args.task.GetId(), got.CurrentTask.GetId())
		})
	}
}

func TestGetCurrentCluster(t *testing.T) {
	tracer := noop.NewTracerProvider().Tracer("")
	service := mockTaskService{}
	mockClusterClient := clusterClient{service: &mockClusterService{}}
	runner, err := newTaskRunner(context.Background(), service, mockClusterClient, tracer)
	require.NoError(t, err)

	currentTaskID := uuid.New()

	type args struct{}
	tests := []struct {
		name        string
		args        args
		wantErr     bool
		wantCluster string
	}{
		{
			name:        "GetCurrentCluster",
			args:        args{},
			wantCluster: "default-5GGzdzBZEA3oeD",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := runner.withTaskExecutionContext(context.Background(), workflowsv1.Task_builder{
				Id: tileboxv1.ID_builder{Uuid: currentTaskID[:]}.Build(),
			}.Build())

			gotCluster, err := GetCurrentCluster(ctx)
			if tt.wantErr {
				require.Error(t, err, "expected an error, got none")
				return
			}
			require.NoError(t, err, "got an unexpected error")

			assert.Equal(t, tt.wantCluster, gotCluster)
		})
	}
}

func TestSetTaskDisplay(t *testing.T) {
	tracer := noop.NewTracerProvider().Tracer("")
	service := mockTaskService{}
	mockClusterClient := clusterClient{service: &mockClusterService{}}
	runner, err := newTaskRunner(context.Background(), service, mockClusterClient, tracer)
	require.NoError(t, err)

	currentTaskID := uuid.New()

	type args struct {
		display string
	}
	tests := []struct {
		name        string
		args        args
		wantErr     bool
		wantDisplay string
	}{
		{
			name: "SetTaskDisplay",
			args: args{
				display: "my task",
			},
			wantDisplay: "my task",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := runner.withTaskExecutionContext(context.Background(), workflowsv1.Task_builder{
				Id: tileboxv1.ID_builder{Uuid: currentTaskID[:]}.Build(),
			}.Build())

			err := SetTaskDisplay(ctx, tt.args.display)
			if tt.wantErr {
				require.Error(t, err, "expected an error, got none")
				return
			}
			require.NoError(t, err, "got an unexpected error")

			te := getTaskExecutionContext(ctx)

			assert.Equal(t, tt.wantDisplay, te.CurrentTask.GetDisplay())
		})
	}
}

func TestSubmitSubtask(t *testing.T) {
	tracer := noop.NewTracerProvider().Tracer("")
	service := mockTaskService{}
	mockClusterClient := clusterClient{service: &mockClusterService{}}
	runner, err := newTaskRunner(context.Background(), service, mockClusterClient, tracer)
	require.NoError(t, err)

	currentTaskID := uuid.New()

	type args struct {
		task    Task
		options []subtask.SubmitOption
	}
	tests := []struct {
		name             string
		setup            []Task
		args             args
		wantErr          string
		wantClusterSlug  string
		wantIdentifier   TaskIdentifier
		wantInput        []byte
		wantDependencies []int64
		wantMaxRetries   int64
		wantFutureTask   subtask.FutureTask
	}{
		{
			name: "Submit nil task",
			args: args{
				task: nil,
			},
			wantErr: "cannot submit nil task",
		},
		{
			name: "Submit subtask",
			args: args{
				task: &testTask1{},
			},
			wantClusterSlug:  "default-5GGzdzBZEA3oeD",
			wantIdentifier:   NewTaskIdentifier("testTask1", "v0.0"),
			wantInput:        []byte("{\"ExecutableTask\":null}"),
			wantDependencies: nil,
			wantMaxRetries:   0,
			wantFutureTask:   subtask.FutureTask(0),
		},
		{
			name: "Submit bad identifier subtasks",
			args: args{
				task: &badIdentifierTask{},
			},
			wantErr: "task name is empty",
		},
		{
			name: "Submit subtask WithClusterSlug",
			args: args{
				task: &testTask1{},
				options: []subtask.SubmitOption{
					subtask.WithClusterSlug("my-production-cluster"),
				},
			},
			wantClusterSlug:  "my-production-cluster",
			wantIdentifier:   NewTaskIdentifier("testTask1", "v0.0"),
			wantInput:        []byte("{\"ExecutableTask\":null}"),
			wantDependencies: nil,
			wantMaxRetries:   0,
			wantFutureTask:   subtask.FutureTask(0),
		},
		{
			name: "Submit subtask WithMaxRetries",
			args: args{
				task: &testTask1{},
				options: []subtask.SubmitOption{
					subtask.WithMaxRetries(12),
				},
			},
			wantClusterSlug:  "default-5GGzdzBZEA3oeD",
			wantIdentifier:   NewTaskIdentifier("testTask1", "v0.0"),
			wantInput:        []byte("{\"ExecutableTask\":null}"),
			wantDependencies: nil,
			wantMaxRetries:   12,
			wantFutureTask:   subtask.FutureTask(0),
		},
		{
			name: "Submit subtask WithDependencies doesn't exist",
			args: args{
				task: &testTask1{},
				options: []subtask.SubmitOption{
					subtask.WithDependencies(subtask.FutureTask(3)),
				},
			},
			wantErr: "future task 3 doesn't exist",
		},
		{
			name: "Submit subtask WithDependencies",
			setup: []Task{
				&testTask2{}, // future task 0
				&testTask2{}, // future task 1
				&testTask2{}, // future task 2
			},
			args: args{
				task: &testTask1{},
				options: []subtask.SubmitOption{
					subtask.WithDependencies(
						subtask.FutureTask(2),
						subtask.FutureTask(0),
					),
				},
			},
			wantClusterSlug:  "default-5GGzdzBZEA3oeD",
			wantIdentifier:   NewTaskIdentifier("testTask1", "v0.0"),
			wantInput:        []byte("{\"ExecutableTask\":null}"),
			wantDependencies: []int64{2, 0},
			wantMaxRetries:   0,
			wantFutureTask:   subtask.FutureTask(3),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := runner.withTaskExecutionContext(context.Background(), workflowsv1.Task_builder{
				Id: tileboxv1.ID_builder{Uuid: currentTaskID[:]}.Build(),
			}.Build())

			for _, task := range tt.setup {
				_, err := SubmitSubtask(ctx, task)
				require.NoError(t, err, "setup failed")
			}

			futureTask, err := SubmitSubtask(ctx, tt.args.task, tt.args.options...)
			if tt.wantErr != "" {
				// we wanted an error, let's check if we got one
				require.Error(t, err, "expected an error, got none")
				assert.Contains(t, err.Error(), tt.wantErr, "error didn't contain expected message: '%s', got error '%s' instead.", tt.wantErr, err.Error())
				return
			}
			// we didn't want an error:
			require.NoError(t, err, "got an unexpected error")

			te := getTaskExecutionContext(ctx)

			assert.Equal(t, tt.wantFutureTask, futureTask)
			lastSubmission := te.subtasks[len(te.subtasks)-1]
			assert.Equal(t, tt.wantClusterSlug, lastSubmission.clusterSlug)
			assert.Equal(t, tt.wantIdentifier.Name(), lastSubmission.identifier.Name())
			assert.Equal(t, tt.wantIdentifier.Version(), lastSubmission.identifier.Version())
			assert.Equal(t, tt.wantInput, lastSubmission.input)
			assert.Equal(t, tt.wantDependencies, lastSubmission.dependencies)
			assert.Equal(t, tt.wantMaxRetries, lastSubmission.maxRetries)
		})
	}
}

func TestSubmitSubtasks(t *testing.T) {
	tracer := noop.NewTracerProvider().Tracer("")
	service := mockTaskService{}
	mockClusterClient := clusterClient{service: &mockClusterService{}}
	runner, err := newTaskRunner(context.Background(), service, mockClusterClient, tracer)
	require.NoError(t, err)

	currentTaskID := uuid.New()

	type args struct {
		tasks []Task
	}
	tests := []struct {
		name             string
		args             args
		wantErr          bool
		wantSubtaskCount int
		wantFutureTasks  []subtask.FutureTask
	}{
		{
			name: "Submit no subtasks",
			args: args{
				tasks: []Task{},
			},
			wantSubtaskCount: 0,
			wantFutureTasks:  []subtask.FutureTask{},
		},
		{
			name: "Submit one subtasks",
			args: args{
				tasks: []Task{
					&testTask1{},
				},
			},
			wantSubtaskCount: 1,
			wantFutureTasks:  []subtask.FutureTask{0},
		},
		{
			name: "Submit bad identifier subtasks",
			args: args{
				tasks: []Task{
					&testTask1{},
					&badIdentifierTask{},
					&testTask2{},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := runner.withTaskExecutionContext(context.Background(), workflowsv1.Task_builder{
				Id: tileboxv1.ID_builder{Uuid: currentTaskID[:]}.Build(),
			}.Build())

			futureTasks, err := SubmitSubtasks(ctx, tt.args.tasks)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			te := getTaskExecutionContext(ctx)

			assert.Len(t, te.subtasks, tt.wantSubtaskCount)

			assert.Len(t, futureTasks, len(tt.wantFutureTasks))
			for i, task := range tt.wantFutureTasks {
				assert.Equal(t, task, futureTasks[i])
			}
		})
	}
}

func TestTrackProgress(t *testing.T) {
	tracer := noop.NewTracerProvider().Tracer("")
	service := mockTaskService{}
	mockClusterClient := clusterClient{service: &mockClusterService{}}
	runner, err := newTaskRunner(context.Background(), service, mockClusterClient, tracer)
	require.NoError(t, err)

	currentTaskID := uuid.New()

	type args struct {
		label string
		total uint64
		done  uint64
	}
	tests := []struct {
		name         string
		args         []args
		wantErr      bool
		wantProgress map[string]*workflowsv1.Progress
	}{
		{
			name:         "no progress updates",
			args:         []args{},
			wantProgress: map[string]*workflowsv1.Progress{},
		},
		{
			name: "single progress indicator with single updates",
			args: []args{
				{label: "my progress", total: 10, done: 0},
			},
			wantProgress: map[string]*workflowsv1.Progress{
				"my progress": workflowsv1.Progress_builder{Label: "my progress", Total: 10, Done: 0}.Build(),
			},
		},
		{
			name: "single progress indicator with multiple updates",
			args: []args{
				{label: "my progress", total: 10, done: 0},
				{label: "my progress", total: 0, done: 2},
				{label: "my progress", total: 0, done: 3},
			},
			wantProgress: map[string]*workflowsv1.Progress{
				"my progress": workflowsv1.Progress_builder{Label: "my progress", Total: 10, Done: 5}.Build(),
			},
		},
		{
			name: "default progress indicator with multiple updates",
			args: []args{
				{total: 10, done: 0},
				{total: 0, done: 2},
				{total: 0, done: 3},
			},
			wantProgress: map[string]*workflowsv1.Progress{
				"": workflowsv1.Progress_builder{Label: "", Total: 10, Done: 5}.Build(),
			},
		},
		{
			name: "multiple progress indicator with multiple updates",
			args: []args{
				{label: "my progress", total: 10, done: 5},
				{total: 4, done: 2}, // default progress indicator
				{label: "other progress", total: 2, done: 2},
				{label: "my progress", total: 5, done: 2},
			},
			wantProgress: map[string]*workflowsv1.Progress{
				"":               workflowsv1.Progress_builder{Label: "", Total: 4, Done: 2}.Build(),
				"my progress":    workflowsv1.Progress_builder{Label: "my progress", Total: 15, Done: 7}.Build(),
				"other progress": workflowsv1.Progress_builder{Label: "other progress", Total: 2, Done: 2}.Build(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := runner.withTaskExecutionContext(context.Background(), workflowsv1.Task_builder{
				Id: tileboxv1.ID_builder{Uuid: currentTaskID[:]}.Build(),
			}.Build())

			for _, progressUpdate := range tt.args {
				var tracker ProgressTracker
				if progressUpdate.label == "" {
					tracker = DefaultProgress()
				} else {
					tracker = Progress(progressUpdate.label)
				}

				if progressUpdate.total > 0 {
					err := tracker.Add(ctx, progressUpdate.total)
					require.NoError(t, err, "failed to add to total progress")
				}

				// we process all the done updates one by one, and we do all of that in parallel to test goroutine safety
				var wg sync.WaitGroup
				errorChannel := make(chan error, 100)
				for i := range progressUpdate.done {
					wg.Go(func() {
						err := tracker.Done(ctx, 1)
						if err != nil {
							errorChannel <- fmt.Errorf("ProgressIndicator.Done() error for update %d: %w", i, err)
						}
					})
				}
				wg.Wait()
				errorChannel <- nil // marker for end of channel

				errors := make([]error, 0)
				for err := range errorChannel {
					if err == nil { // end of channel marker received
						break
					}
					errors = append(errors, err)
				}
				if len(errors) > 0 {
					t.Errorf("received %d errors in parallel progress tracking", len(errors))
					for _, err := range errors {
						t.Error(err.Error())
					}
					return
				}
			}

			progressIndicators := getTaskExecutionContext(ctx).getProgressUpdates()
			require.Len(t, progressIndicators, len(tt.wantProgress))
			for _, progress := range progressIndicators {
				wantProgress, found := tt.wantProgress[progress.GetLabel()]
				if !found {
					t.Errorf("got unexpected progress indicator: %s", progress.GetLabel())
					continue
				}
				assert.Equal(t, wantProgress.GetLabel(), progress.GetLabel())
				assert.Equal(t, wantProgress.GetTotal(), progress.GetTotal())
				assert.Equal(t, wantProgress.GetDone(), progress.GetDone())
			}
		})
	}
}

type cacheKey string

const fibonacciCacheCtxKey cacheKey = "x-fibonacci-test-workflow-cache"

type FibonacciTask struct {
	N int
}

func (t *FibonacciTask) Identifier() TaskIdentifier {
	return NewTaskIdentifier("tilebox.com/fibonacci/FibonacciTask", "v1.0")
}

func (t *FibonacciTask) Execute(ctx context.Context) error {
	cache := ctx.Value(fibonacciCacheCtxKey).(map[string]int)

	key := fmt.Sprintf("fib_%d", t.N)
	_, ok := cache[key]
	if ok {
		// result already cached and computed, so nothing to do
		return nil
	}
	if t.N <= 2 { // base cases: fib(1) = 1, fib(2) = 1
		cache[key] = 1
	} else {
		fibMinus1And2, err := SubmitSubtasks(ctx, []Task{&FibonacciTask{N: t.N - 1}, &FibonacciTask{N: t.N - 2}})
		if err != nil {
			return err
		}

		_, err = SubmitSubtask(ctx, &SumResultTask{N: t.N}, subtask.WithDependencies(fibMinus1And2...))
		if err != nil {
			return err
		}
	}

	return nil
}

type SumResultTask struct {
	N int
}

func (t *SumResultTask) Identifier() TaskIdentifier {
	return NewTaskIdentifier("tilebox.com/fibonacci/SumResultTask", "v1.0")
}

func (t *SumResultTask) Execute(ctx context.Context) error {
	cache := ctx.Value(fibonacciCacheCtxKey).(map[string]int)
	key := fmt.Sprintf("fib_%d", t.N)
	// calculate and store the result
	cache[key] = cache[fmt.Sprintf("fib_%d", t.N-1)] + cache[fmt.Sprintf("fib_%d", t.N-2)]
	return nil
}

func TestRunnerFibonacciWorkflow(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})))
	client := NewReplayClient(t, "fibonacci_workflow.json")

	ctx := context.Background()
	cache := make(map[string]int)
	ctx = context.WithValue(ctx, fibonacciCacheCtxKey, cache)

	runner, err := client.NewTaskRunner(ctx)
	require.NoError(t, err, "failed to create task runner")

	err = runner.RegisterTasks(&FibonacciTask{}, &SumResultTask{})
	require.NoError(t, err, "failed to register task")

	_, err = client.Jobs.Submit(ctx, "fibonacci", []Task{&FibonacciTask{N: 7}})
	require.NoError(t, err, "failed to submit job")

	runner.RunAll(ctx)

	// make sure all results are computed
	assert.Equal(t, 13, cache["fib_7"])
	assert.Equal(t, 8, cache["fib_6"])
	assert.Equal(t, 5, cache["fib_5"])
	assert.Equal(t, 3, cache["fib_4"])
	assert.Equal(t, 2, cache["fib_3"])
	assert.Equal(t, 1, cache["fib_2"])
	assert.Equal(t, 1, cache["fib_1"])
}
