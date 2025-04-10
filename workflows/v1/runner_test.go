package workflows

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	workflowsv1 "github.com/tilebox/tilebox-go/protogen/go/workflows/v1"
	"go.opentelemetry.io/otel/trace/noop"
)

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
			runner, err := newTaskRunner(service, tracer, WithCluster("testing-4qgCk4qHH85qR7"))
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
		id *workflowsv1.UUID
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
				id: &workflowsv1.UUID{Uuid: []byte{1, 2, 3}},
			},
			want: false,
		},
		{
			name: "isEmpty not nil but valid id",
			args: args{
				id: &workflowsv1.UUID{Uuid: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6}},
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

func Test_withTaskExecutionContextRoundtrip(t *testing.T) {
	tracer := noop.NewTracerProvider().Tracer("")
	service := mockTaskService{}
	runner, err := newTaskRunner(service, tracer, WithCluster("testing-4qgCk4qHH85qR7"))
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
				task: &workflowsv1.Task{
					Id: &workflowsv1.UUID{Uuid: taskID[:]},
				},
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

func TestSubmitSubtask(t *testing.T) {
	tracer := noop.NewTracerProvider().Tracer("")
	service := mockTaskService{}
	cluster := "testing-4qgCk4qHH85qR7"
	runner, err := newTaskRunner(service, tracer, WithCluster(cluster))
	require.NoError(t, err)

	currentTaskID := uuid.New()

	type args struct {
		task    Task
		options []SubmitSubtaskOption
	}
	tests := []struct {
		name        string
		args        args
		wantErr     string
		wantSubtask *FutureTask
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
			wantSubtask: &FutureTask{
				submission: &workflowsv1.TaskSubmission{
					ClusterSlug:  cluster,
					Identifier:   &workflowsv1.TaskIdentifier{Name: "testTask1", Version: "v0.0"},
					Input:        []byte("{\"ExecutableTask\":null}"),
					Display:      "testTask1",
					Dependencies: nil,
					MaxRetries:   0,
				},
			},
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
				options: []SubmitSubtaskOption{
					WithClusterSlug("my-production-cluster"),
				},
			},
			wantSubtask: &FutureTask{
				submission: &workflowsv1.TaskSubmission{
					ClusterSlug:  "my-production-cluster",
					Identifier:   &workflowsv1.TaskIdentifier{Name: "testTask1", Version: "v0.0"},
					Input:        []byte("{\"ExecutableTask\":null}"),
					Display:      "testTask1",
					Dependencies: nil,
					MaxRetries:   0,
				},
			},
		},
		{
			name: "Submit subtask WithMaxRetries",
			args: args{
				task: &testTask1{},
				options: []SubmitSubtaskOption{
					WithMaxRetries(12),
				},
			},
			wantSubtask: &FutureTask{
				submission: &workflowsv1.TaskSubmission{
					ClusterSlug:  cluster,
					Identifier:   &workflowsv1.TaskIdentifier{Name: "testTask1", Version: "v0.0"},
					Input:        []byte("{\"ExecutableTask\":null}"),
					Display:      "testTask1",
					Dependencies: nil,
					MaxRetries:   12,
				},
			},
		},
		{
			name: "Submit subtask WithDependencies",
			args: args{
				task: &testTask1{},
				options: []SubmitSubtaskOption{
					WithDependencies(
						&FutureTask{index: 3},
						&FutureTask{index: 1},
					),
				},
			},
			wantSubtask: &FutureTask{
				submission: &workflowsv1.TaskSubmission{
					ClusterSlug:  cluster,
					Identifier:   &workflowsv1.TaskIdentifier{Name: "testTask1", Version: "v0.0"},
					Input:        []byte("{\"ExecutableTask\":null}"),
					Display:      "testTask1",
					Dependencies: []int64{3, 1},
					MaxRetries:   0,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := runner.withTaskExecutionContext(context.Background(), &workflowsv1.Task{
				Id: &workflowsv1.UUID{Uuid: currentTaskID[:]},
			})

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

			assert.Equal(t, tt.wantSubtask, futureTask)
			assert.Equal(t, tt.wantSubtask, te.Subtasks[0])
		})
	}
}

func TestSubmitSubtasks(t *testing.T) {
	tracer := noop.NewTracerProvider().Tracer("")
	service := mockTaskService{}
	cluster := "testing-4qgCk4qHH85qR7"
	runner, err := newTaskRunner(service, tracer, WithCluster(cluster))
	require.NoError(t, err)

	currentTaskID := uuid.New()

	type args struct {
		tasks []Task
	}
	tests := []struct {
		name         string
		args         args
		wantErr      bool
		wantSubtasks []*FutureTask
	}{
		{
			name: "Submit no Subtasks",
			args: args{
				tasks: []Task{},
			},
			wantSubtasks: []*FutureTask{},
		},
		{
			name: "Submit one Subtasks",
			args: args{
				tasks: []Task{
					&testTask1{},
				},
			},
			wantSubtasks: []*FutureTask{
				{
					submission: &workflowsv1.TaskSubmission{
						ClusterSlug:  cluster,
						Identifier:   &workflowsv1.TaskIdentifier{Name: "testTask1", Version: "v0.0"},
						Input:        []byte("{\"ExecutableTask\":null}"),
						Display:      "testTask1",
						Dependencies: nil,
						MaxRetries:   0,
					},
				},
			},
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
			wantSubtasks: []*FutureTask{
				{
					submission: &workflowsv1.TaskSubmission{
						ClusterSlug:  cluster,
						Identifier:   &workflowsv1.TaskIdentifier{Name: "testTask1", Version: "v0.0"},
						Input:        []byte("{\"ExecutableTask\":null}"),
						Display:      "testTask1",
						Dependencies: nil,
						MaxRetries:   0,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := runner.withTaskExecutionContext(context.Background(), &workflowsv1.Task{
				Id: &workflowsv1.UUID{Uuid: currentTaskID[:]},
			})

			futureTasks, err := SubmitSubtasks(ctx, tt.args.tasks)
			if (err != nil) != tt.wantErr {
				t.Errorf("SubmitSubtasks() error = %v, wantErr %v", err, tt.wantErr)
			}

			te := getTaskExecutionContext(ctx)
			assert.Len(t, futureTasks, len(tt.wantSubtasks))
			assert.Len(t, te.Subtasks, len(tt.wantSubtasks))

			for i, task := range tt.wantSubtasks {
				assert.Equal(t, task, futureTasks[i])
				assert.Equal(t, task, te.Subtasks[i])
			}
		})
	}
}
