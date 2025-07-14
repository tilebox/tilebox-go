package workflows

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tileboxv1 "github.com/tilebox/tilebox-go/protogen/go/tilebox/v1"
	workflowsv1 "github.com/tilebox/tilebox-go/protogen/go/workflows/v1"
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
		wantErr   bool
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
			name: "Register Tasks duplicated task",
			args: args{
				tasks: []ExecutableTask{
					&testTask1{},
					&testTask1{},
				},
			},
			wantTasks: []ExecutableTask{
				&testTask1{},
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
			wantErr: true,
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
			if (err != nil) != tt.wantErr {
				t.Errorf("RegisterTasks() error = %v, wantErr %v", err, tt.wantErr)
			}

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

func (m *mockMinimalTaskService) TaskFailed(context.Context, uuid.UUID, string, bool) (*workflowsv1.TaskStateResponse, error) {
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
			assert.Empty(t, got.Subtasks)
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
		name           string
		setup          []Task
		args           args
		wantErr        string
		wantSubmission *workflowsv1.TaskSubmission
		wantFutureTask subtask.FutureTask
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
			wantSubmission: workflowsv1.TaskSubmission_builder{
				ClusterSlug:  "default-5GGzdzBZEA3oeD",
				Identifier:   workflowsv1.TaskIdentifier_builder{Name: "testTask1", Version: "v0.0"}.Build(),
				Input:        []byte("{\"ExecutableTask\":null}"),
				Display:      "testTask1",
				Dependencies: nil,
				MaxRetries:   0,
			}.Build(),
			wantFutureTask: subtask.FutureTask(0),
		},
		{
			name: "Submit bad identifier Subtasks",
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
			wantSubmission: workflowsv1.TaskSubmission_builder{
				ClusterSlug:  "my-production-cluster",
				Identifier:   workflowsv1.TaskIdentifier_builder{Name: "testTask1", Version: "v0.0"}.Build(),
				Input:        []byte("{\"ExecutableTask\":null}"),
				Display:      "testTask1",
				Dependencies: nil,
				MaxRetries:   0,
			}.Build(),
			wantFutureTask: subtask.FutureTask(0),
		},
		{
			name: "Submit subtask WithMaxRetries",
			args: args{
				task: &testTask1{},
				options: []subtask.SubmitOption{
					subtask.WithMaxRetries(12),
				},
			},
			wantSubmission: workflowsv1.TaskSubmission_builder{
				ClusterSlug:  "default-5GGzdzBZEA3oeD",
				Identifier:   workflowsv1.TaskIdentifier_builder{Name: "testTask1", Version: "v0.0"}.Build(),
				Input:        []byte("{\"ExecutableTask\":null}"),
				Display:      "testTask1",
				Dependencies: nil,
				MaxRetries:   12,
			}.Build(),
			wantFutureTask: subtask.FutureTask(0),
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
			wantSubmission: workflowsv1.TaskSubmission_builder{
				ClusterSlug:  "default-5GGzdzBZEA3oeD",
				Identifier:   workflowsv1.TaskIdentifier_builder{Name: "testTask1", Version: "v0.0"}.Build(),
				Input:        []byte("{\"ExecutableTask\":null}"),
				Display:      "testTask1",
				Dependencies: []int64{2, 0},
				MaxRetries:   0,
			}.Build(),
			wantFutureTask: subtask.FutureTask(3),
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
			assert.Equal(t, tt.wantSubmission, te.Subtasks[len(te.Subtasks)-1]) // last submission
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
		name            string
		args            args
		wantErr         bool
		wantSubmissions []*workflowsv1.TaskSubmission
		wantFutureTasks []subtask.FutureTask
	}{
		{
			name: "Submit no Subtasks",
			args: args{
				tasks: []Task{},
			},
			wantSubmissions: []*workflowsv1.TaskSubmission{},
			wantFutureTasks: []subtask.FutureTask{},
		},
		{
			name: "Submit one Subtasks",
			args: args{
				tasks: []Task{
					&testTask1{},
				},
			},
			wantSubmissions: []*workflowsv1.TaskSubmission{
				workflowsv1.TaskSubmission_builder{
					ClusterSlug:  "default-5GGzdzBZEA3oeD",
					Identifier:   workflowsv1.TaskIdentifier_builder{Name: "testTask1", Version: "v0.0"}.Build(),
					Input:        []byte("{\"ExecutableTask\":null}"),
					Display:      "testTask1",
					Dependencies: nil,
					MaxRetries:   0,
				}.Build(),
			},
			wantFutureTasks: []subtask.FutureTask{0},
		},
		{
			name: "Submit bad identifier Subtasks",
			args: args{
				tasks: []Task{
					&testTask1{},
					&badIdentifierTask{},
					&testTask2{},
				},
			},
			wantErr: true,
			wantSubmissions: []*workflowsv1.TaskSubmission{
				workflowsv1.TaskSubmission_builder{
					ClusterSlug:  "default-5GGzdzBZEA3oeD",
					Identifier:   workflowsv1.TaskIdentifier_builder{Name: "testTask1", Version: "v0.0"}.Build(),
					Input:        []byte("{\"ExecutableTask\":null}"),
					Display:      "testTask1",
					Dependencies: nil,
					MaxRetries:   0,
				}.Build(),
			},
			wantFutureTasks: []subtask.FutureTask{0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := runner.withTaskExecutionContext(context.Background(), workflowsv1.Task_builder{
				Id: tileboxv1.ID_builder{Uuid: currentTaskID[:]}.Build(),
			}.Build())

			futureTasks, err := SubmitSubtasks(ctx, tt.args.tasks)
			if (err != nil) != tt.wantErr {
				t.Errorf("SubmitSubtasks() error = %v, wantErr %v", err, tt.wantErr)
			}

			te := getTaskExecutionContext(ctx)

			assert.Len(t, te.Subtasks, len(tt.wantSubmissions))
			for i, task := range tt.wantSubmissions {
				assert.Equal(t, task, te.Subtasks[i])
			}

			assert.Len(t, futureTasks, len(tt.wantFutureTasks))
			for i, task := range tt.wantFutureTasks {
				assert.Equal(t, task, futureTasks[i])
			}
		})
	}
}
