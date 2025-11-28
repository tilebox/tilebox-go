package workflows

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	workflowsv1 "github.com/tilebox/tilebox-go/protogen/workflows/v1"
	"github.com/tilebox/tilebox-go/query"
	"github.com/tilebox/tilebox-go/workflows/v1/job"
)

type mockJobService struct {
	JobService
	reqs []*workflowsv1.SubmitJobRequest
}

func (m *mockJobService) SubmitJob(_ context.Context, req *workflowsv1.SubmitJobRequest) (*workflowsv1.Job, error) {
	m.reqs = append(m.reqs, req)

	return workflowsv1.Job_builder{
		Name: req.GetJobName(),
	}.Build(), nil
}

func TestJobService_Submit(t *testing.T) {
	type args struct {
		jobName string
		tasks   []Task
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
				jobName: "test-job",
				tasks:   []Task{&testTask1{}},
			},
			want: &Job{
				Name: "test-job",
			},
			wantReq: workflowsv1.SubmitJobRequest_builder{
				Tasks: workflowsv1.TaskSubmissions_builder{
					TaskGroups: []*workflowsv1.TaskSubmissionGroup{
						workflowsv1.TaskSubmissionGroup_builder{
							Inputs:              [][]byte{[]byte("{\"ExecutableTask\":null}")},
							IdentifierPointers:  []uint64{0},
							ClusterSlugPointers: []uint64{0},
							DisplayPointers:     []uint64{0},
							MaxRetriesValues:    []int64{0},
						}.Build(),
					},
					ClusterSlugLookup: []string{""},
					IdentifierLookup: []*workflowsv1.TaskIdentifier{
						workflowsv1.TaskIdentifier_builder{
							Name:    "testTask1",
							Version: "v0.0",
						}.Build(),
					},
					DisplayLookup: []string{"testTask1"},
				}.Build(),
				JobName: "test-job",
			}.Build(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &mockJobService{}
			js := jobClient{service: service}
			got, err := js.Submit(context.Background(), tt.args.jobName, tt.args.tasks)
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

	foundJob, err := client.Jobs.Get(ctx, uuid.MustParse("0199614b-6102-1498-815b-3cb23f72b489"))
	require.NoError(t, err)

	assert.Equal(t, "progress-example", foundJob.Name)
	assert.Equal(t, "0199614b-6102-1498-815b-3cb23f72b489", foundJob.ID.String())
	assert.Equal(t, job.Completed, foundJob.State)
	assert.Len(t, foundJob.TaskSummaries, 6)
	require.Len(t, foundJob.Progress, 1)
	assert.Empty(t, foundJob.Progress[0].Label)
	assert.Equal(t, uint64(5), foundJob.Progress[0].Total)
	assert.Equal(t, uint64(5), foundJob.Progress[0].Done)
}

func Test_jobClient_Query(t *testing.T) {
	ctx := context.Background()
	client := NewReplayClient(t, "Jobs")

	startDate := time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, time.February, 1, 0, 0, 0, 0, time.UTC)
	timeInterval := query.NewTimeInterval(startDate, endDate)
	jobs, err := Collect(client.Jobs.Query(ctx, job.WithTemporalExtent(timeInterval)))
	require.NoError(t, err)

	firstJob := jobs[0]
	assert.Equal(t, "my-windows-job", firstJob.Name)
	assert.Equal(t, "0194ad17-bdaf-ff8e-983b-d1299fd2d235", firstJob.ID.String())
	assert.Equal(t, job.Completed, firstJob.State)
}
