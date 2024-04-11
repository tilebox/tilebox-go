package workflows

import (
	"context"
	"reflect"
	"testing"

	"connectrpc.com/connect"
	workflowsv1 "github.com/tilebox/tilebox-go/protogen/go/workflows/v1"
	"github.com/tilebox/tilebox-go/protogen/go/workflows/v1/workflowsv1connect"
)

type mockJobServiceClient struct {
	workflowsv1connect.JobServiceClient
	reqs []*workflowsv1.SubmitJobRequest
}

func (m *mockJobServiceClient) SubmitJob(_ context.Context, req *connect.Request[workflowsv1.SubmitJobRequest]) (*connect.Response[workflowsv1.Job], error) {
	m.reqs = append(m.reqs, req.Msg)

	return connect.NewResponse(&workflowsv1.Job{
		Name: req.Msg.GetJobName(),
	}), nil
}

func TestJobService_Submit(t *testing.T) {
	ctx := context.Background()

	type args struct {
		jobName     string
		clusterSlug string
		tasks       []Task
	}
	tests := []struct {
		name    string
		args    args
		want    *workflowsv1.Job
		wantReq *workflowsv1.SubmitJobRequest
		wantErr bool
	}{
		{
			name: "Submit Job",
			args: args{
				jobName:     "test-job",
				clusterSlug: "test-cluster",
				tasks:       []Task{&testTask1{}},
			},
			want: &workflowsv1.Job{
				Name: "test-job",
			},
			wantReq: &workflowsv1.SubmitJobRequest{
				Tasks: []*workflowsv1.TaskSubmission{
					{
						ClusterSlug: "test-cluster",
						Identifier: &workflowsv1.TaskIdentifier{
							Name:    "testTask1",
							Version: "v0.0",
						},
						Input:   []byte("{\"ExecutableTask\":null}"),
						Display: "testTask1",
					},
				},
				JobName: "test-job",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := mockJobServiceClient{}
			js := NewJobService(&client)
			got, err := js.Submit(ctx, tt.args.jobName, tt.args.clusterSlug, tt.args.tasks...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Submit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Submit() got = %v, want %v", got, tt.want)
			}

			// Verify the submitted request
			if len(client.reqs) != 1 {
				t.Fatalf("Submit() expected 1 request, got %d", len(client.reqs))
			}
			if !reflect.DeepEqual(client.reqs[0].GetTasks(), tt.wantReq.GetTasks()) {
				t.Errorf("Submit() request = %v, want %v", client.reqs[0], tt.wantReq)
			}
			if !reflect.DeepEqual(client.reqs[0].GetJobName(), tt.wantReq.GetJobName()) {
				t.Errorf("Submit() request = %v, want %v", client.reqs[0], tt.wantReq)
			}
			// We don't compare traceparent because it's generated randomly
		})
	}
}
