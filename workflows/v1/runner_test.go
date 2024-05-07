package workflows

import (
	"context"
	"reflect"
	"testing"

	workflowsv1 "github.com/tilebox/tilebox-go/protogen/go/workflows/v1"
	"github.com/tilebox/tilebox-go/protogen/go/workflows/v1/workflowsv1connect"
)

type mockClient struct {
	workflowsv1connect.TaskServiceClient
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

func TestTaskRunner_RegisterTask(t *testing.T) {
	type args struct {
		task ExecutableTask
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Register Task",
			args: args{
				task: &testTask1{},
			},
		},
		{
			name: "Register Task bad identifier",
			args: args{
				task: &badIdentifierTask{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t1 *testing.T) {
			runner := NewTaskRunner(mockClient{})

			err := runner.RegisterTask(tt.args.task)
			if (err != nil) != tt.wantErr {
				t1.Errorf("RegisterTask() error = %v, wantErr %v", err, tt.wantErr)
			}

			identifier := identifierFromTask(tt.args.task)
			_, ok := runner.GetRegisteredTask(identifier)
			if ok && tt.wantErr {
				t1.Errorf("RegisterTask() task found in taskDefinitions")
			}
			if !ok && !tt.wantErr {
				t1.Errorf("RegisterTask() task not found in taskDefinitions")
			}
		})
	}
}

func TestTaskRunner_RegisterTasks(t *testing.T) {
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
		t.Run(tt.name, func(t1 *testing.T) {
			runner := NewTaskRunner(mockClient{})

			err := runner.RegisterTasks(tt.args.tasks...)
			if (err != nil) != tt.wantErr {
				t1.Errorf("RegisterTasks() error = %v, wantErr %v", err, tt.wantErr)
			}

			if len(runner.taskDefinitions) != len(tt.wantTasks) {
				t1.Errorf("RegisterTasks() taskDefinitions length = %v, want %v", len(runner.taskDefinitions), len(tt.wantTasks))
			}
			for _, task := range tt.wantTasks {
				identifier := identifierFromTask(task)
				_, ok := runner.GetRegisteredTask(identifier)
				if !ok {
					t1.Errorf("RegisterTasks() task not found in taskDefinitions")
				}
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
			if got := isEmpty(tt.args.id); got != tt.want {
				t.Errorf("isEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_withTaskExecutionContextRoundtrip(t *testing.T) {
	type args struct {
		ctx    context.Context
		client workflowsv1connect.TaskServiceClient
		task   *workflowsv1.Task
	}
	tests := []struct {
		name string
		args args
		want *taskExecutionContext
	}{
		{
			name: "withTaskExecutionContext",
			args: args{
				ctx:    context.Background(),
				client: mockClient{},
				task:   nil,
			},
			want: &taskExecutionContext{
				CurrentTask: nil,
				Client:      mockClient{},
				Subtasks:    []*workflowsv1.TaskSubmission{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updatedCtx := withTaskExecutionContext(tt.args.ctx, tt.args.client, tt.args.task)
			got := getTaskExecutionContext(updatedCtx)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("withTaskExecutionContext() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSubmitSubtasks(t *testing.T) {
	type args struct {
		tasks []Task
	}
	tests := []struct {
		name         string
		args         args
		wantErr      bool
		wantSubtasks []*workflowsv1.TaskSubmission
	}{
		{
			name: "Submit no Subtasks",
			args: args{
				tasks: []Task{},
			},
			wantSubtasks: []*workflowsv1.TaskSubmission{},
		},
		{
			name: "Submit one Subtasks",
			args: args{
				tasks: []Task{
					&testTask1{},
				},
			},
			wantSubtasks: []*workflowsv1.TaskSubmission{
				{
					ClusterSlug:  DefaultClusterSlug,
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
				tasks: []Task{
					&testTask1{},
					&badIdentifierTask{},
					&testTask2{},
				},
			},
			wantErr: true,
			wantSubtasks: []*workflowsv1.TaskSubmission{
				{
					ClusterSlug:  DefaultClusterSlug,
					Identifier:   &workflowsv1.TaskIdentifier{Name: "testTask1", Version: "v0.0"},
					Input:        []byte("{\"ExecutableTask\":null}"),
					Display:      "testTask1",
					Dependencies: nil,
					MaxRetries:   0,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := withTaskExecutionContext(context.Background(), nil, nil)

			err := SubmitSubtasks(ctx, tt.args.tasks...)
			if (err != nil) != tt.wantErr {
				t.Errorf("SubmitSubtasks() error = %v, wantErr %v", err, tt.wantErr)
			}

			te := getTaskExecutionContext(ctx)
			if len(te.Subtasks) != len(tt.wantSubtasks) {
				t.Errorf("SubmitSubtasks() Subtasks length = %v, want %v", len(te.Subtasks), len(tt.wantSubtasks))
			}

			for i, task := range tt.wantSubtasks {
				if !reflect.DeepEqual(te.Subtasks[i], task) {
					t.Errorf("SubmitSubtasks() Subtask %v = %v, want %v", i, te.Subtasks[i], task)
				}
			}
		})
	}
}
