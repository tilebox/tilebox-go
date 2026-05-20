package workflows

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tileboxv1 "github.com/tilebox/tilebox-go/protogen/tilebox/v1"
	workflowsv1 "github.com/tilebox/tilebox-go/protogen/workflows/v1"
	"github.com/tilebox/tilebox-go/query"
	"github.com/tilebox/tilebox-go/workflows/v1/job"
	commonv1 "go.opentelemetry.io/proto/otlp/common/v1"
	logsv1 "go.opentelemetry.io/proto/otlp/logs/v1"
	resourcev1 "go.opentelemetry.io/proto/otlp/resource/v1"
	tracev1 "go.opentelemetry.io/proto/otlp/trace/v1"
)

type mockJobService struct {
	JobService

	reqs          []*workflowsv1.SubmitJobRequest
	queryPages    []*workflowsv1.QueryJobsResponse
	queryFilters  []*workflowsv1.QueryFilters
	queryPageReqs []*tileboxv1.Pagination
	querySortReqs []tileboxv1.SortDirection
}

type mockTelemetryService struct {
	TelemetryService

	logPages       []*workflowsv1.PaginatedLogsData
	spanPages      []*workflowsv1.PaginatedSpansData
	logPageReqs    []*tileboxv1.Pagination
	spanPageReqs   []*tileboxv1.Pagination
	logSortReqs    []*tileboxv1.SortDirection
	spanSortReqs   []*tileboxv1.SortDirection
	queryLogJobID  uuid.UUID
	querySpanJobID uuid.UUID
}

func (m *mockJobService) SubmitJob(_ context.Context, req *workflowsv1.SubmitJobRequest) (*workflowsv1.Job, error) {
	m.reqs = append(m.reqs, req)

	return workflowsv1.Job_builder{
		Name: req.GetJobName(),
	}.Build(), nil
}

func (m *mockJobService) QueryJobs(_ context.Context, filters *workflowsv1.QueryFilters, page *tileboxv1.Pagination, sortDirection tileboxv1.SortDirection) (*workflowsv1.QueryJobsResponse, error) {
	m.queryFilters = append(m.queryFilters, filters)
	m.queryPageReqs = append(m.queryPageReqs, page)
	m.querySortReqs = append(m.querySortReqs, sortDirection)

	pageIndex := len(m.queryPageReqs) - 1
	if pageIndex >= len(m.queryPages) {
		return workflowsv1.QueryJobsResponse_builder{}.Build(), nil
	}
	return m.queryPages[pageIndex], nil
}

func (m *mockTelemetryService) QueryJobLogs(_ context.Context, jobID uuid.UUID, page *tileboxv1.Pagination, sortDirection *tileboxv1.SortDirection) (*workflowsv1.PaginatedLogsData, error) {
	m.queryLogJobID = jobID
	m.logPageReqs = append(m.logPageReqs, page)
	m.logSortReqs = append(m.logSortReqs, sortDirection)

	pageIndex := len(m.logPageReqs) - 1
	return m.logPages[pageIndex], nil
}

func (m *mockTelemetryService) QueryJobSpans(_ context.Context, jobID uuid.UUID, page *tileboxv1.Pagination, sortDirection *tileboxv1.SortDirection) (*workflowsv1.PaginatedSpansData, error) {
	m.querySpanJobID = jobID
	m.spanPageReqs = append(m.spanPageReqs, page)
	m.spanSortReqs = append(m.spanSortReqs, sortDirection)

	pageIndex := len(m.spanPageReqs) - 1
	return m.spanPages[pageIndex], nil
}

func TestTaskStateString(t *testing.T) {
	tests := []struct {
		name  string
		state TaskState
		want  string
	}{
		{name: "queued", state: TaskQueued, want: "queued"},
		{name: "running", state: TaskRunning, want: "running"},
		{name: "computed", state: TaskComputed, want: "computed"},
		{name: "failed", state: TaskFailed, want: "failed"},
		{name: "skipped", state: TaskSkipped, want: "skipped"},
		{name: "failed optional", state: TaskFailedOptional, want: "failed_optional"},
		{name: "unspecified", state: 0, want: "unspecified"},
		{name: "unknown", state: 99, want: "unspecified"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, test.state.String())

			data, err := json.Marshal(test.state)
			require.NoError(t, err)
			assert.JSONEq(t, `"`+test.want+`"`, string(data))
		})
	}
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
							OptionalValues:      []bool{false},
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

			// verify the request was sent
			assert.Len(t, service.reqs, 1)

			assert.Equal(t, tt.wantReq.GetJobName(), service.reqs[0].GetJobName())
			// we don't compare traceparent because it's generated randomly

			taskSubmissions := service.reqs[0].GetTasks()
			assert.Equal(t, tt.wantReq.GetTasks().GetTaskGroups(), taskSubmissions.GetTaskGroups())
			assert.Equal(t, tt.wantReq.GetTasks().GetClusterSlugLookup(), taskSubmissions.GetClusterSlugLookup())
			assert.Equal(t, tt.wantReq.GetTasks().GetIdentifierLookup(), taskSubmissions.GetIdentifierLookup())
			assert.Equal(t, tt.wantReq.GetTasks().GetDisplayLookup(), taskSubmissions.GetDisplayLookup())
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

func Test_jobClient_Query_WithSortDirection(t *testing.T) {
	service := &mockJobService{}
	client := jobClient{service: service}
	timeInterval := query.NewTimeInterval(
		time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2025, time.February, 1, 0, 0, 0, 0, time.UTC),
	)

	jobs, err := Collect(client.Query(context.Background(), job.WithTemporalExtent(timeInterval), job.WithSortDirection(job.Ascending)))
	require.NoError(t, err)
	assert.Empty(t, jobs)

	require.Len(t, service.querySortReqs, 1)
	assert.Equal(t, tileboxv1.SortDirection_SORT_DIRECTION_ASCENDING, service.querySortReqs[0])
}

func Test_jobClient_Query_WithoutTemporalExtent(t *testing.T) {
	service := &mockJobService{}
	client := jobClient{service: service}

	jobs, err := Collect(client.Query(context.Background()))
	require.NoError(t, err)
	assert.Empty(t, jobs)

	require.Len(t, service.queryFilters, 1)
	assert.Nil(t, service.queryFilters[0].GetTimeInterval())
	assert.Nil(t, service.queryFilters[0].GetIdInterval())
}

func Test_jobClient_Query_WithTaskStates(t *testing.T) {
	service := &mockJobService{}
	client := jobClient{service: service}
	timeInterval := query.NewTimeInterval(
		time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2025, time.February, 1, 0, 0, 0, 0, time.UTC),
	)

	jobs, err := Collect(client.Query(context.Background(), job.WithTemporalExtent(timeInterval), job.WithTaskStates(job.TaskRunning, job.TaskFailed)))
	require.NoError(t, err)
	assert.Empty(t, jobs)

	require.Len(t, service.queryFilters, 1)
	assert.Equal(t, []workflowsv1.TaskState{workflowsv1.TaskState_TASK_STATE_RUNNING, workflowsv1.TaskState_TASK_STATE_FAILED}, service.queryFilters[0].GetTaskStates())
}

func Test_jobClient_Query_WithLimit(t *testing.T) {
	nextPage := tileboxv1.Pagination_builder{StartingAfter: tileboxv1.NewUUID(uuid.MustParse("01994da4-255e-740d-9df7-b8c1aa41c75b"))}.Build()
	service := &mockJobService{
		queryPages: []*workflowsv1.QueryJobsResponse{
			workflowsv1.QueryJobsResponse_builder{
				Jobs: []*workflowsv1.Job{
					workflowsv1.Job_builder{Name: "first", Id: tileboxv1.NewUUID(uuid.New())}.Build(),
				},
				NextPage: nextPage,
			}.Build(),
			workflowsv1.QueryJobsResponse_builder{
				Jobs: []*workflowsv1.Job{
					workflowsv1.Job_builder{Name: "second", Id: tileboxv1.NewUUID(uuid.New())}.Build(),
					workflowsv1.Job_builder{Name: "third", Id: tileboxv1.NewUUID(uuid.New())}.Build(),
				},
			}.Build(),
		},
	}
	client := jobClient{service: service}
	timeInterval := query.NewTimeInterval(
		time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2025, time.February, 1, 0, 0, 0, 0, time.UTC),
	)

	jobs, err := Collect(client.Query(context.Background(), job.WithTemporalExtent(timeInterval), job.WithLimit(2)))
	require.NoError(t, err)
	require.Len(t, jobs, 2)

	require.Len(t, service.queryPageReqs, 2)
	assert.Equal(t, int64(2), service.queryPageReqs[0].GetLimit())
	assert.Nil(t, service.queryPageReqs[0].GetStartingAfter())
	assert.Equal(t, int64(1), service.queryPageReqs[1].GetLimit())
	assert.Equal(t, nextPage.GetStartingAfter(), service.queryPageReqs[1].GetStartingAfter())
	assert.Equal(t, "first", jobs[0].Name)
	assert.Equal(t, "second", jobs[1].Name)
}

func Test_jobClient_QueryPage(t *testing.T) {
	cursorID := uuid.MustParse("01994da4-255e-740d-9df7-b8c1aa41c75b")
	nextCursorID := uuid.MustParse("01994da4-255e-740d-9df7-b8c1aa41c75c")
	service := &mockJobService{
		queryPages: []*workflowsv1.QueryJobsResponse{
			workflowsv1.QueryJobsResponse_builder{
				Jobs: []*workflowsv1.Job{
					workflowsv1.Job_builder{Name: "first", Id: tileboxv1.NewUUID(uuid.New())}.Build(),
				},
				NextPage: tileboxv1.Pagination_builder{StartingAfter: tileboxv1.NewUUID(nextCursorID)}.Build(),
			}.Build(),
		},
	}
	client := jobClient{service: service}
	timeInterval := query.NewTimeInterval(
		time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2025, time.February, 1, 0, 0, 0, 0, time.UTC),
	)

	page, err := client.QueryPage(context.Background(), job.WithTemporalExtent(timeInterval), job.WithCursor(job.NewCursor(cursorID)), job.WithLimit(10), job.WithSortDirection(job.Descending))
	require.NoError(t, err)
	require.Len(t, page.Jobs, 1)
	require.NotNil(t, page.NextCursor)

	assert.Equal(t, "first", page.Jobs[0].Name)
	assert.Equal(t, nextCursorID, page.NextCursor.StartingAfter())
	require.Len(t, service.queryPageReqs, 1)
	assert.Equal(t, int64(10), service.queryPageReqs[0].GetLimit())
	assert.Equal(t, cursorID, service.queryPageReqs[0].GetStartingAfter().AsUUID())
	require.Len(t, service.querySortReqs, 1)
	assert.Equal(t, tileboxv1.SortDirection_SORT_DIRECTION_DESCENDING, service.querySortReqs[0])
}

func Test_jobClient_QueryLogs(t *testing.T) {
	jobID := uuid.MustParse("01994da3-779b-7d01-a570-2a04d0ac163b")
	nextPage := tileboxv1.Pagination_builder{StartingAfter: tileboxv1.NewUUID(uuid.MustParse("01994da4-255e-740d-9df7-b8c1aa41c75b"))}.Build()
	service := &mockTelemetryService{
		logPages: []*workflowsv1.PaginatedLogsData{
			workflowsv1.PaginatedLogsData_builder{
				ResourceLogs: []*logsv1.ResourceLogs{testResourceLogs("first log", 10)},
				NextPage:     nextPage,
			}.Build(),
			workflowsv1.PaginatedLogsData_builder{
				ResourceLogs: []*logsv1.ResourceLogs{testResourceLogs("second log", 20)},
			}.Build(),
		},
	}
	client := jobClient{telemetryService: service}

	logs, err := Collect(client.QueryLogs(context.Background(), jobID))
	require.NoError(t, err)
	require.Len(t, logs, 2)

	assert.Equal(t, jobID, service.queryLogJobID)
	require.Len(t, service.logPageReqs, 2)
	assert.Nil(t, service.logPageReqs[0])
	assert.Equal(t, nextPage, service.logPageReqs[1])
	assert.Equal(t, []*tileboxv1.SortDirection{nil, nil}, service.logSortReqs)
	assert.Equal(t, "first log", logs[0].Body)
	assert.Equal(t, 10, logs[0].SeverityNumber)
	assert.Equal(t, "INFO", logs[0].SeverityText)
	assert.Equal(t, "010203", logs[0].TraceID)
	assert.Equal(t, map[string]any{"attempt": int64(3)}, logs[0].Attributes)
	assert.Equal(t, map[string]any{"runner": "test-runner"}, logs[0].RunnerAttributes)
}

func Test_jobClient_QueryLogs_WithOptions(t *testing.T) {
	jobID := uuid.MustParse("01994da3-779b-7d01-a570-2a04d0ac163b")
	nextPage := tileboxv1.Pagination_builder{StartingAfter: tileboxv1.NewUUID(uuid.MustParse("01994da4-255e-740d-9df7-b8c1aa41c75b"))}.Build()
	service := &mockTelemetryService{
		logPages: []*workflowsv1.PaginatedLogsData{
			workflowsv1.PaginatedLogsData_builder{
				ResourceLogs: []*logsv1.ResourceLogs{testResourceLogs("first log", 10)},
				NextPage:     nextPage,
			}.Build(),
			workflowsv1.PaginatedLogsData_builder{
				ResourceLogs: []*logsv1.ResourceLogs{testResourceLogs("second log", 20), testResourceLogs("third log", 30)},
			}.Build(),
		},
	}
	client := jobClient{telemetryService: service}

	logs, err := Collect(client.QueryLogs(context.Background(), jobID, job.WithLimit(2), job.WithSortDirection(job.Descending)))
	require.NoError(t, err)
	require.Len(t, logs, 2)

	require.Len(t, service.logPageReqs, 2)
	assert.Equal(t, int64(2), service.logPageReqs[0].GetLimit())
	assert.Nil(t, service.logPageReqs[0].GetStartingAfter())
	assert.Equal(t, int64(1), service.logPageReqs[1].GetLimit())
	assert.Equal(t, nextPage.GetStartingAfter(), service.logPageReqs[1].GetStartingAfter())
	require.Len(t, service.logSortReqs, 2)
	assert.Equal(t, tileboxv1.SortDirection_SORT_DIRECTION_DESCENDING, *service.logSortReqs[0])
	assert.Equal(t, tileboxv1.SortDirection_SORT_DIRECTION_DESCENDING, *service.logSortReqs[1])
	assert.Equal(t, "first log", logs[0].Body)
	assert.Equal(t, "second log", logs[1].Body)
}

func Test_jobClient_QueryLogsPage(t *testing.T) {
	jobID := uuid.MustParse("01994da3-779b-7d01-a570-2a04d0ac163b")
	cursorID := uuid.MustParse("01994da4-255e-740d-9df7-b8c1aa41c75b")
	nextCursorID := uuid.MustParse("01994da4-255e-740d-9df7-b8c1aa41c75c")
	service := &mockTelemetryService{
		logPages: []*workflowsv1.PaginatedLogsData{
			workflowsv1.PaginatedLogsData_builder{
				ResourceLogs: []*logsv1.ResourceLogs{testResourceLogs("first log", 10)},
				NextPage:     tileboxv1.Pagination_builder{StartingAfter: tileboxv1.NewUUID(nextCursorID)}.Build(),
			}.Build(),
		},
	}
	client := jobClient{telemetryService: service}

	page, err := client.QueryLogsPage(context.Background(), jobID, job.WithCursor(job.NewCursor(cursorID)), job.WithLimit(10), job.WithSortDirection(job.Descending))
	require.NoError(t, err)
	require.Len(t, page.Logs, 1)
	require.NotNil(t, page.NextCursor)

	assert.Equal(t, "first log", page.Logs[0].Body)
	assert.Equal(t, nextCursorID, page.NextCursor.StartingAfter())
	require.Len(t, service.logPageReqs, 1)
	assert.Equal(t, int64(10), service.logPageReqs[0].GetLimit())
	assert.Equal(t, cursorID, service.logPageReqs[0].GetStartingAfter().AsUUID())
	require.Len(t, service.logSortReqs, 1)
	assert.Equal(t, tileboxv1.SortDirection_SORT_DIRECTION_DESCENDING, *service.logSortReqs[0])
}

func Test_jobClient_QuerySpans(t *testing.T) {
	jobID := uuid.MustParse("01994da3-779b-7d01-a570-2a04d0ac163b")
	service := &mockTelemetryService{
		spanPages: []*workflowsv1.PaginatedSpansData{
			workflowsv1.PaginatedSpansData_builder{
				ResourceSpans: []*tracev1.ResourceSpans{testResourceSpans()},
			}.Build(),
		},
	}
	client := jobClient{telemetryService: service}

	spans, err := Collect(client.QuerySpans(context.Background(), jobID))
	require.NoError(t, err)
	require.Len(t, spans, 1)

	assert.Equal(t, jobID, service.querySpanJobID)
	require.Len(t, service.spanPageReqs, 1)
	assert.Nil(t, service.spanPageReqs[0])
	assert.Equal(t, []*tileboxv1.SortDirection{nil}, service.spanSortReqs)
	assert.Equal(t, "task/run", spans[0].Name)
	assert.Equal(t, "010203", spans[0].TraceID)
	assert.Equal(t, "040506", spans[0].SpanID)
	assert.Equal(t, "070809", spans[0].ParentSpanID)
	assert.Equal(t, "STATUS_CODE_ERROR", spans[0].StatusCode)
	assert.Equal(t, "failed", spans[0].StatusMessage)
	assert.Equal(t, 10*time.Millisecond, spans[0].Duration())
	assert.Equal(t, map[string]any{"task": "example"}, spans[0].Attributes)
	assert.Equal(t, map[string]any{"runner": "test-runner"}, spans[0].RunnerAttributes)
	require.Len(t, spans[0].Events, 1)
	assert.Equal(t, "retry", spans[0].Events[0].Name)
	assert.Equal(t, map[string]any{"count": int64(1)}, spans[0].Events[0].Attributes)
}

func Test_jobClient_QuerySpans_WithOptions(t *testing.T) {
	jobID := uuid.MustParse("01994da3-779b-7d01-a570-2a04d0ac163b")
	service := &mockTelemetryService{
		spanPages: []*workflowsv1.PaginatedSpansData{
			workflowsv1.PaginatedSpansData_builder{
				ResourceSpans: []*tracev1.ResourceSpans{testResourceSpans(), testResourceSpans()},
			}.Build(),
		},
	}
	client := jobClient{telemetryService: service}

	spans, err := Collect(client.QuerySpans(context.Background(), jobID, job.WithLimit(1), job.WithSortDirection(job.Ascending)))
	require.NoError(t, err)
	require.Len(t, spans, 1)

	require.Len(t, service.spanPageReqs, 1)
	assert.Equal(t, int64(1), service.spanPageReqs[0].GetLimit())
	require.Len(t, service.spanSortReqs, 1)
	assert.Equal(t, tileboxv1.SortDirection_SORT_DIRECTION_ASCENDING, *service.spanSortReqs[0])
}

func Test_jobClient_QuerySpansPage(t *testing.T) {
	jobID := uuid.MustParse("01994da3-779b-7d01-a570-2a04d0ac163b")
	cursorID := uuid.MustParse("01994da4-255e-740d-9df7-b8c1aa41c75b")
	nextCursorID := uuid.MustParse("01994da4-255e-740d-9df7-b8c1aa41c75c")
	service := &mockTelemetryService{
		spanPages: []*workflowsv1.PaginatedSpansData{
			workflowsv1.PaginatedSpansData_builder{
				ResourceSpans: []*tracev1.ResourceSpans{testResourceSpans()},
				NextPage:      tileboxv1.Pagination_builder{StartingAfter: tileboxv1.NewUUID(nextCursorID)}.Build(),
			}.Build(),
		},
	}
	client := jobClient{telemetryService: service}

	page, err := client.QuerySpansPage(context.Background(), jobID, job.WithCursor(job.NewCursor(cursorID)), job.WithLimit(10), job.WithSortDirection(job.Ascending))
	require.NoError(t, err)
	require.Len(t, page.Spans, 1)
	require.NotNil(t, page.NextCursor)

	assert.Equal(t, "task/run", page.Spans[0].Name)
	assert.Equal(t, nextCursorID, page.NextCursor.StartingAfter())
	require.Len(t, service.spanPageReqs, 1)
	assert.Equal(t, int64(10), service.spanPageReqs[0].GetLimit())
	assert.Equal(t, cursorID, service.spanPageReqs[0].GetStartingAfter().AsUUID())
	require.Len(t, service.spanSortReqs, 1)
	assert.Equal(t, tileboxv1.SortDirection_SORT_DIRECTION_ASCENDING, *service.spanSortReqs[0])
}

func testResourceLogs(body string, severity logsv1.SeverityNumber) *logsv1.ResourceLogs {
	return &logsv1.ResourceLogs{
		Resource: &resourcev1.Resource{Attributes: []*commonv1.KeyValue{stringKeyValue("runner", "test-runner")}},
		ScopeLogs: []*logsv1.ScopeLogs{
			{
				LogRecords: []*logsv1.LogRecord{
					{
						TimeUnixNano:   uint64(time.Unix(1, 2).UnixNano()),
						SeverityNumber: severity,
						SeverityText:   "INFO",
						Body:           stringValue(body),
						TraceId:        []byte{1, 2, 3},
						SpanId:         []byte{4, 5, 6},
						Attributes:     []*commonv1.KeyValue{intKeyValue("attempt", 3)},
					},
				},
			},
		},
	}
}

func testResourceSpans() *tracev1.ResourceSpans {
	return &tracev1.ResourceSpans{
		Resource: &resourcev1.Resource{Attributes: []*commonv1.KeyValue{stringKeyValue("runner", "test-runner")}},
		ScopeSpans: []*tracev1.ScopeSpans{
			{
				Spans: []*tracev1.Span{
					{
						TraceId:           []byte{1, 2, 3},
						SpanId:            []byte{4, 5, 6},
						ParentSpanId:      []byte{7, 8, 9},
						Name:              "task/run",
						StartTimeUnixNano: uint64(time.Unix(1, 0).UnixNano()),
						EndTimeUnixNano:   uint64(time.Unix(1, int64(10*time.Millisecond)).UnixNano()),
						Attributes:        []*commonv1.KeyValue{stringKeyValue("task", "example")},
						Status:            &tracev1.Status{Code: tracev1.Status_STATUS_CODE_ERROR, Message: "failed"},
						Events: []*tracev1.Span_Event{
							{
								TimeUnixNano: uint64(time.Unix(1, int64(5*time.Millisecond)).UnixNano()),
								Name:         "retry",
								Attributes:   []*commonv1.KeyValue{intKeyValue("count", 1)},
							},
						},
					},
				},
			},
		},
	}
}

func stringKeyValue(key string, value string) *commonv1.KeyValue {
	return &commonv1.KeyValue{Key: key, Value: stringValue(value)}
}

func intKeyValue(key string, value int64) *commonv1.KeyValue {
	return &commonv1.KeyValue{Key: key, Value: &commonv1.AnyValue{Value: &commonv1.AnyValue_IntValue{IntValue: value}}}
}

func stringValue(value string) *commonv1.AnyValue {
	return &commonv1.AnyValue{Value: &commonv1.AnyValue_StringValue{StringValue: value}}
}
