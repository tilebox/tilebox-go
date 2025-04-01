package workflows

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	workflowsv1 "github.com/tilebox/tilebox-go/protogen/go/workflows/v1"
)

type mockJobService struct {
	JobService
	reqs []*workflowsv1.SubmitJobRequest
}

func (m *mockJobService) SubmitJob(_ context.Context, req *workflowsv1.SubmitJobRequest) (*workflowsv1.Job, error) {
	m.reqs = append(m.reqs, req)

	return &workflowsv1.Job{
		Name: req.GetJobName(),
	}, nil
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
		want    *Job
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
			want: &Job{
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
			service := &mockJobService{}
			js := jobClient{service: service}
			got, err := js.Submit(ctx, tt.args.jobName, tt.args.clusterSlug, 0, tt.args.tasks...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Submit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.want.Name, got.Name)

			// Verify the submitted request
			assert.Len(t, service.reqs, 1)
			assert.Equal(t, tt.wantReq.GetTasks(), service.reqs[0].GetTasks())
			assert.Equal(t, tt.wantReq.GetJobName(), service.reqs[0].GetJobName())
			// We don't compare traceparent because it's generated randomly
		})
	}
}

func Test_jobClient_Get(t *testing.T) {
	ctx := context.Background()
	client := NewReplayClient(t, "Job")

	job, err := client.Jobs.Get(ctx, uuid.MustParse("0194ad17-bdaf-ff8e-983b-d1299fd2d235"))
	require.NoError(t, err)

	assert.Equal(t, "my-windows-job", job.Name)
	assert.Equal(t, "0194ad17-bdaf-ff8e-983b-d1299fd2d235", job.ID.String())
	assert.Equal(t, JobCompleted, job.State)
}

func Test_jobClient_List(t *testing.T) {
	ctx := context.Background()
	client := NewReplayClient(t, "Jobs")

	jobs, err := Collect(client.Jobs.List(ctx, &workflowsv1.IDInterval{
		StartId: "00000000-0000-0000-0000-000000000000",
		EndId:   "ffffffff-ffff-ffff-ffff-ffffffffffff",
	}))
	require.NoError(t, err)

	job := jobs[0]
	assert.Equal(t, "my-windows-job", job.Name)
	assert.Equal(t, "0194ad17-bdaf-ff8e-983b-d1299fd2d235", job.ID.String())
	assert.Equal(t, JobCompleted, job.State)
}
