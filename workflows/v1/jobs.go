package workflows

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"iter"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	tileboxv1 "github.com/tilebox/tilebox-go/protogen/tilebox/v1"
	workflowsv1 "github.com/tilebox/tilebox-go/protogen/workflows/v1"
	"github.com/tilebox/tilebox-go/workflows/v1/job"
	commonv1 "go.opentelemetry.io/proto/otlp/common/v1"
	logsv1 "go.opentelemetry.io/proto/otlp/logs/v1"
	tracev1 "go.opentelemetry.io/proto/otlp/trace/v1"
	"google.golang.org/protobuf/proto"
)

// Job represents a Tilebox Workflows job.
//
// Documentation: https://docs.tilebox.com/workflows/concepts/jobs
type Job struct {
	// ID is the unique identifier of the job.
	ID uuid.UUID
	// Name is the name of the job.
	Name string
	// State is the current state of the job.
	State job.State
	// SubmittedAt is the time the job was submitted.
	SubmittedAt time.Time
	// AutomationID is the ID of the automation that submitted the job.
	AutomationID uuid.UUID
	// Progress is a list of progress indicators for the job.
	Progress []*ProgressIndicator
	// ExecutionStats contains execution statistics for the job.
	ExecutionStats *ExecutionStats
}

// JobTask is a task in a job's user-visible task tree.
type JobTask struct {
	// ID is the unique identifier of the task.
	ID uuid.UUID
	// ParentID is the closest user-visible parent task. Nil indicates a root task.
	ParentID *uuid.UUID
	// Display is a human-readable representation of the task.
	Display string
	// State is the current state of the task.
	State TaskState
	// SubmittedAt is the time the task was submitted.
	SubmittedAt time.Time
	// StartedAt is the time the task started running. It is zero if the task has not started.
	StartedAt time.Time
	// StoppedAt is the time the task stopped running. It is zero if the task has not stopped.
	StoppedAt time.Time
	// Input is the serialized task input.
	Input []byte
	// RetryCount is the number of times the task has been retried.
	RetryCount int64
	// MaxRetries is the maximum number of times the task may be retried.
	MaxRetries int64
	// HasChildren reports whether the task has children in the user-visible tree.
	HasChildren bool
	// ClusterSlug is the cluster on which the task is scheduled to run.
	ClusterSlug string
	// Optional reports whether failure of the task is optional.
	Optional bool
}

// TaskState is the state of a Task.
type TaskState int32

type ProgressIndicator struct {
	Label string
	Total uint64
	Done  uint64
}

// ExecutionStats contains execution statistics for a job.
type ExecutionStats struct {
	// FirstTaskStartedAt is the time the first task of the job was started.
	FirstTaskStartedAt time.Time
	// LastTaskStoppedAt is the time the last task of the job was stopped.
	LastTaskStoppedAt time.Time
	// ComputeTime is the total compute time of the job, as sum of all task compute times.
	ComputeTime time.Duration
	// ElapsedTime is the elapsed time of the job (wall time), which is the time from first task started to last task stopped.
	ElapsedTime time.Duration
	// Parallelism is the parallelism factor of the job, which is the average number of tasks running at any given time.
	Parallelism float64
	// TotalTasks is the total number of tasks of the job.
	TotalTasks uint64
	// TasksByState is the number of tasks by their state.
	TasksByState []*TaskStateCount
}

// TaskStateCount represents the count of tasks in a particular state.
type TaskStateCount struct {
	// State is the task state.
	State TaskState
	// Count is the number of tasks in this state.
	Count uint64
}

// LogRecord is a log record emitted while running a job.
type LogRecord struct {
	// Time is the time the log record occurred.
	Time time.Time
	// SeverityNumber is the OpenTelemetry severity number.
	SeverityNumber int
	// SeverityText is the textual severity, such as INFO or ERROR.
	SeverityText string
	// Body is the log message body.
	Body string
	// TraceID is the hex-encoded OpenTelemetry trace ID. Empty if the log is not associated with a trace.
	TraceID string
	// SpanID is the hex-encoded OpenTelemetry span ID. Empty if the log is not associated with a span.
	SpanID string
	// Attributes are the log record attributes.
	Attributes map[string]any
	// RunnerAttributes are resource attributes describing the task runner that emitted the log.
	RunnerAttributes map[string]any
}

// Span is a trace span emitted while running a job.
type Span struct {
	// StartTime is the time the span started.
	StartTime time.Time
	// EndTime is the time the span ended.
	EndTime time.Time
	// TraceID is the hex-encoded OpenTelemetry trace ID.
	TraceID string
	// SpanID is the hex-encoded OpenTelemetry span ID.
	SpanID string
	// ParentSpanID is the hex-encoded parent span ID. Empty for root spans.
	ParentSpanID string
	// Name is the span name.
	Name string
	// StatusCode is the OpenTelemetry status code name.
	StatusCode string
	// StatusMessage is the OpenTelemetry status message.
	StatusMessage string
	// Attributes are the span attributes.
	Attributes map[string]any
	// RunnerAttributes are resource attributes describing the task runner that emitted the span.
	RunnerAttributes map[string]any
	// Events are the events attached to the span.
	Events []*SpanEvent
}

// Duration returns the elapsed time of the span.
func (s *Span) Duration() time.Duration {
	return s.EndTime.Sub(s.StartTime)
}

// SpanEvent is an event attached to a trace span.
type SpanEvent struct {
	// Time is the time the event occurred.
	Time time.Time
	// Name is the event name.
	Name string
	// Attributes are the event attributes.
	Attributes map[string]any
}

// JobPage is a single page of jobs returned by a manual page query.
type JobPage struct {
	// Jobs is the list of jobs in this page.
	Jobs []*Job
	// NextCursor can be used to request the next page. Nil means there are no more pages.
	NextCursor *job.Cursor
}

// JobTaskPage is one page of root tasks or direct children of one task.
type JobTaskPage struct {
	// ParentTaskID is the parent of this sibling collection. Nil indicates root tasks.
	ParentTaskID *uuid.UUID
	// Tasks is the list of tasks in this page.
	Tasks []*JobTask
	// NextCursor can be used with the same parent task to request the next page.
	NextCursor *job.Cursor
	// PrefetchedChildPages contains the first child page for tasks in this page when prefetching was requested.
	PrefetchedChildPages []*JobTaskPage
}

// LogPage is a single page of log records returned by a manual page query.
type LogPage struct {
	// Logs is the list of log records in this page.
	Logs []*LogRecord
	// NextCursor can be used to request the next page. Nil means there are no more pages.
	NextCursor *job.Cursor
}

// SpanPage is a single page of spans returned by a manual page query.
type SpanPage struct {
	// Spans is the list of spans in this page.
	Spans []*Span
	// NextCursor can be used to request the next page. Nil means there are no more pages.
	NextCursor *job.Cursor
}

// TaskState values.
const (
	_                  TaskState = iota
	TaskQueued                   // The task is queued and waiting to be run.
	TaskRunning                  // The task is currently running on some task runner.
	TaskComputed                 // The task has been computed and the output is available.
	TaskFailed                   // The task has failed.
	TaskSkipped                  // The task has been skipped.
	TaskFailedOptional           // The task has failed, but was marked as optional.
)

func (state TaskState) String() string {
	switch state {
	case TaskQueued:
		return "queued"
	case TaskRunning:
		return "running"
	case TaskComputed:
		return "computed"
	case TaskFailed:
		return "failed"
	case TaskSkipped:
		return "skipped"
	case TaskFailedOptional:
		return "failed_optional"
	default:
		return "unspecified"
	}
}

func (state TaskState) MarshalJSON() ([]byte, error) {
	return json.Marshal(state.String())
}

type JobClient interface {
	// Submit submits a job to a cluster.
	//
	// Options:
	//   - job.WithMaxRetries: sets the maximum number of times a job can be automatically retried. Defaults to 0.
	//   - job.WithClusterSlug: sets the cluster slug of the cluster where the job will be executed. Defaults to the default cluster.
	//
	// Documentation: https://docs.tilebox.com/workflows/concepts/jobs#submission
	Submit(ctx context.Context, jobName string, tasks []Task, options ...job.SubmitOption) (*Job, error)

	// Get returns a job by its ID.
	Get(ctx context.Context, jobID uuid.UUID) (*Job, error)

	// ListTasks returns root tasks or the direct children of one task.
	//
	// Options:
	//   - job.WithParentTaskID: lists direct children of the specified task instead of root tasks. (Optional)
	//   - job.WithCursor: starts after a cursor returned for the same sibling collection. (Optional)
	//   - job.WithLimit: limits the total number of tasks returned. Defaults to unlimited.
	//
	// Pagination is handled automatically for this sibling collection. Child collections are not traversed.
	ListTasks(ctx context.Context, jobID uuid.UUID, options ...job.TaskListOption) iter.Seq2[*JobTask, error]

	// ListTasksPage returns one page of root tasks or direct children of one task.
	//
	// Options:
	//   - job.WithParentTaskID: lists direct children of the specified task instead of root tasks. (Optional)
	//   - job.WithCursor: starts after a cursor returned for the same sibling collection. (Optional)
	//   - job.WithLimit: limits the number of tasks returned in this page. Defaults to the server default.
	//   - job.WithPrefetchChildren: includes the first child page for each task returned in this page. (Optional)
	ListTasksPage(ctx context.Context, jobID uuid.UUID, options ...job.TaskPageOption) (*JobTaskPage, error)

	// Retry retries a job.
	//
	// Returns the number of rescheduled tasks.
	//
	// Documentation: https://docs.tilebox.com/workflows/concepts/jobs#retry-handling
	Retry(ctx context.Context, jobID uuid.UUID) (int64, error)

	// Cancel cancels a job.
	//
	// Documentation: https://docs.tilebox.com/workflows/concepts/jobs#cancellation
	Cancel(ctx context.Context, jobID uuid.UUID) error

	// Query returns a list of jobs matching the provided options.
	//
	// Options:
	//   - job.WithTemporalExtent: specifies the time or ID interval for which jobs should be queried. (Optional)
	//   - job.WithAutomationIDs: specifies the automation IDs to filter jobs by. Only jobs submitted by the specified
	//  automation will be returned. (Optional)
	//   - job.WithClusterSlugs: specifies the cluster slugs to filter jobs by. Only jobs submitted to the specified
	//     clusters will be returned. (Optional)
	//   - job.WithJobStates: specifies the job state to filter jobs by. Only jobs with the specified state will be returned. (Optional)
	//   - job.WithTaskStates: specifies task states to filter jobs by. Only jobs with at least one task in one of the
	//     specified states will be returned. (Optional)
	//   - job.WithJobName: specifies the job name to filter jobs by. Only jobs with the specified name will be returned. (Optional)
	//   - job.WithCursor: starts the query after a cursor returned by a previous page. (Optional)
	//   - job.WithLimit: limits the total number of jobs returned. Defaults to unlimited.
	//   - job.WithSortDirection: specifies the direction to sort jobs by submission date. (Optional)
	//
	// Pagination is handled automatically under the hood. The jobs are lazily loaded and returned as a sequence.
	// The output sequence can be transformed into a slice using Collect.
	Query(ctx context.Context, options ...job.QueryOption) iter.Seq2[*Job, error]

	// QueryPage returns a single page of jobs matching the provided options.
	//
	// Options:
	//   - job.WithTemporalExtent: specifies the time or ID interval for which jobs should be queried. (Optional)
	//   - job.WithAutomationIDs: specifies the automation IDs to filter jobs by. Only jobs submitted by the specified
	//  automation will be returned. (Optional)
	//   - job.WithClusterSlugs: specifies the cluster slugs to filter jobs by. Only jobs submitted to the specified
	//     clusters will be returned. (Optional)
	//   - job.WithJobStates: specifies the job state to filter jobs by. Only jobs with the specified state will be returned. (Optional)
	//   - job.WithTaskStates: specifies task states to filter jobs by. Only jobs with at least one task in one of the
	//     specified states will be returned. (Optional)
	//   - job.WithJobName: specifies the job name to filter jobs by. Only jobs with the specified name will be returned. (Optional)
	//   - job.WithCursor: starts the query after a cursor returned by a previous page. (Optional)
	//   - job.WithLimit: limits the number of jobs returned in this page. Defaults to the server default.
	//   - job.WithSortDirection: specifies the direction to sort jobs by submission date. (Optional)
	QueryPage(ctx context.Context, options ...job.QueryOption) (*JobPage, error)

	// QueryLogs returns the logs emitted while running a job.
	//
	// Options:
	//   - job.WithTaskID: filters logs by task ID. Defaults to all tasks in the job.
	//   - job.WithCursor: starts the query after a cursor returned by a previous page. (Optional)
	//   - job.WithLimit: limits the total number of records returned. Defaults to unlimited.
	//   - job.WithSortDirection: sorts records ascending or descending. Defaults to the server default.
	//   - job.WithLogSeverityGroups: filters logs by severity group. Defaults to all severities.
	//
	// Pagination is handled automatically under the hood. The logs are lazily loaded and returned as a sequence.
	// The output sequence can be transformed into a slice using Collect.
	QueryLogs(ctx context.Context, jobID uuid.UUID, options ...job.LogQueryOption) iter.Seq2[*LogRecord, error]

	// QueryLogsPage returns a single page of logs emitted while running a job.
	//
	// Options:
	//   - job.WithTaskID: filters logs by task ID. Defaults to all tasks in the job.
	//   - job.WithCursor: starts the query after a cursor returned by a previous page. (Optional)
	//   - job.WithLimit: limits the number of records returned in this page. Defaults to the server default.
	//   - job.WithSortDirection: sorts records ascending or descending. Defaults to the server default.
	//   - job.WithLogSeverityGroups: filters logs by severity group. Defaults to all severities.
	QueryLogsPage(ctx context.Context, jobID uuid.UUID, options ...job.LogQueryOption) (*LogPage, error)

	// QuerySpans returns the spans emitted while running a job.
	//
	// Options:
	//   - job.WithTaskID: filters spans by task ID. Defaults to all tasks in the job.
	//   - job.WithCursor: starts the query after a cursor returned by a previous page. (Optional)
	//   - job.WithLimit: limits the total number of records returned. Defaults to unlimited.
	//   - job.WithSortDirection: sorts records ascending or descending. Defaults to the server default.
	//
	// Pagination is handled automatically under the hood. The spans are lazily loaded and returned as a sequence.
	// The output sequence can be transformed into a slice using Collect.
	QuerySpans(ctx context.Context, jobID uuid.UUID, options ...job.TelemetryQueryOption) iter.Seq2[*Span, error]

	// QuerySpansPage returns a single page of spans emitted while running a job.
	//
	// Options:
	//   - job.WithTaskID: filters spans by task ID. Defaults to all tasks in the job.
	//   - job.WithCursor: starts the query after a cursor returned by a previous page. (Optional)
	//   - job.WithLimit: limits the number of records returned in this page. Defaults to the server default.
	//   - job.WithSortDirection: sorts records ascending or descending. Defaults to the server default.
	QuerySpansPage(ctx context.Context, jobID uuid.UUID, options ...job.TelemetryQueryOption) (*SpanPage, error)
}

var _ JobClient = &jobClient{}

type jobClient struct {
	service          JobService
	telemetryService TelemetryService
}

func (c jobClient) Submit(ctx context.Context, jobName string, tasks []Task, options ...job.SubmitOption) (*Job, error) {
	opts := &job.SubmitOptions{}
	for _, option := range options {
		option(opts)
	}

	jobRequest, err := validateJob(jobName, opts.ClusterSlug, opts.MaxRetries, tasks...)
	if err != nil {
		return nil, err
	}

	response, err := c.service.SubmitJob(ctx, jobRequest)
	if err != nil {
		return nil, err
	}

	return protoToJob(response), nil
}

func (c jobClient) Get(ctx context.Context, jobID uuid.UUID) (*Job, error) {
	response, err := c.service.GetJob(ctx, jobID)
	if err != nil {
		return nil, err
	}

	return protoToJob(response), nil
}

func (c jobClient) ListTasks(ctx context.Context, jobID uuid.UUID, options ...job.TaskListOption) iter.Seq2[*JobTask, error] {
	opts := &job.TaskListOptions{}
	for _, option := range options {
		option.ApplyTaskListOption(opts)
	}

	return func(yield func(*JobTask, error) bool) {
		cursor := opts.Cursor
		remaining := opts.Limit

		for {
			pageOpts := &job.TaskPageOptions{TaskListOptions: *opts}
			pageOpts.Cursor = cursor
			if opts.Limit > 0 {
				if remaining == 0 {
					break
				}
				pageOpts.Limit = remaining
			}

			taskPage, err := c.listTasksPage(ctx, jobID, pageOpts)
			if err != nil {
				yield(nil, err)
				return
			}

			for _, task := range taskPage.Tasks {
				if opts.Limit > 0 && remaining == 0 {
					return
				}
				if !yield(task, nil) {
					return
				}
				if opts.Limit > 0 {
					remaining--
				}
			}

			cursor = taskPage.NextCursor
			if cursor == nil || (opts.Limit > 0 && remaining == 0) {
				break
			}
		}
	}
}

func (c jobClient) ListTasksPage(ctx context.Context, jobID uuid.UUID, options ...job.TaskPageOption) (*JobTaskPage, error) {
	opts := &job.TaskPageOptions{}
	for _, option := range options {
		option.ApplyTaskPageOption(opts)
	}
	return c.listTasksPage(ctx, jobID, opts)
}

func (c jobClient) Retry(ctx context.Context, jobID uuid.UUID) (int64, error) {
	response, err := c.service.RetryJob(ctx, jobID)
	if err != nil {
		return 0, err
	}

	return response.GetNumTasksRescheduled(), nil
}

func (c jobClient) Cancel(ctx context.Context, jobID uuid.UUID) error {
	return c.service.CancelJob(ctx, jobID)
}

func (c jobClient) Query(ctx context.Context, options ...job.QueryOption) iter.Seq2[*Job, error] {
	opts := &job.QueryOptions{}
	for _, option := range options {
		option.ApplyQueryOption(opts)
	}

	return func(yield func(*Job, error) bool) {
		cursor := opts.Cursor
		remaining := opts.Limit

		for {
			pageOpts := *opts
			pageOpts.Cursor = cursor
			if opts.Limit > 0 {
				if remaining == 0 {
					break
				}
				pageOpts.Limit = remaining
			}

			jobsPage, err := c.queryPage(ctx, &pageOpts)
			if err != nil {
				yield(nil, err)
				return
			}

			for _, result := range jobsPage.Jobs {
				if opts.Limit > 0 && remaining == 0 {
					return
				}

				if !yield(result, nil) {
					return
				}
				if opts.Limit > 0 {
					remaining--
				}
			}

			cursor = jobsPage.NextCursor
			if cursor == nil || (opts.Limit > 0 && remaining == 0) {
				break
			}
		}
	}
}

func (c jobClient) QueryPage(ctx context.Context, options ...job.QueryOption) (*JobPage, error) {
	opts := &job.QueryOptions{}
	for _, option := range options {
		option.ApplyQueryOption(opts)
	}
	return c.queryPage(ctx, opts)
}

func (c jobClient) QueryLogs(ctx context.Context, jobID uuid.UUID, options ...job.LogQueryOption) iter.Seq2[*LogRecord, error] {
	opts := &job.LogQueryOptions{}
	for _, option := range options {
		option.ApplyLogQueryOption(opts)
	}

	return func(yield func(*LogRecord, error) bool) {
		cursor := opts.Cursor
		remaining := opts.Limit

		for {
			pageOpts := *opts
			pageOpts.Cursor = cursor
			if opts.Limit > 0 {
				if remaining == 0 {
					break
				}
				pageOpts.Limit = remaining
			}

			logsPage, err := c.queryLogsPage(ctx, jobID, &pageOpts)
			if err != nil {
				yield(nil, err)
				return
			}

			for _, logRecord := range logsPage.Logs {
				if opts.Limit > 0 && remaining == 0 {
					return
				}

				if !yield(logRecord, nil) {
					return
				}
				if opts.Limit > 0 {
					remaining--
				}
			}

			cursor = logsPage.NextCursor
			if cursor == nil || (opts.Limit > 0 && remaining == 0) {
				break
			}
		}
	}
}

func (c jobClient) QueryLogsPage(ctx context.Context, jobID uuid.UUID, options ...job.LogQueryOption) (*LogPage, error) {
	opts := &job.LogQueryOptions{}
	for _, option := range options {
		option.ApplyLogQueryOption(opts)
	}
	return c.queryLogsPage(ctx, jobID, opts)
}

func (c jobClient) QuerySpans(ctx context.Context, jobID uuid.UUID, options ...job.TelemetryQueryOption) iter.Seq2[*Span, error] {
	opts := &job.TelemetryQueryOptions{}
	for _, option := range options {
		option.ApplyTelemetryQueryOption(opts)
	}

	return func(yield func(*Span, error) bool) {
		cursor := opts.Cursor
		remaining := opts.Limit

		for {
			pageOpts := *opts
			pageOpts.Cursor = cursor
			if opts.Limit > 0 {
				if remaining == 0 {
					break
				}
				pageOpts.Limit = remaining
			}

			spansPage, err := c.querySpansPage(ctx, jobID, &pageOpts)
			if err != nil {
				yield(nil, err)
				return
			}

			for _, span := range spansPage.Spans {
				if opts.Limit > 0 && remaining == 0 {
					return
				}

				if !yield(span, nil) {
					return
				}
				if opts.Limit > 0 {
					remaining--
				}
			}

			cursor = spansPage.NextCursor
			if cursor == nil || (opts.Limit > 0 && remaining == 0) {
				break
			}
		}
	}
}

func (c jobClient) QuerySpansPage(ctx context.Context, jobID uuid.UUID, options ...job.TelemetryQueryOption) (*SpanPage, error) {
	opts := &job.TelemetryQueryOptions{}
	for _, option := range options {
		option.ApplyTelemetryQueryOption(opts)
	}
	return c.querySpansPage(ctx, jobID, opts)
}

func (c jobClient) querySpansPage(ctx context.Context, jobID uuid.UUID, opts *job.TelemetryQueryOptions) (*SpanPage, error) {
	spansMessage, err := c.telemetryService.QueryJobSpans(ctx, jobID, taskIDToProto(opts.TaskID), paginationFromOptions(opts.Limit, opts.Cursor), sortDirectionToProto(opts.SortDirection))
	if err != nil {
		return nil, err
	}

	spans := make([]*Span, 0, len(spansMessage.GetResourceSpans()))
	for _, resourceSpans := range spansMessage.GetResourceSpans() {
		span, err := protoToSpan(resourceSpans)
		if err != nil {
			return nil, err
		}
		spans = append(spans, span)
	}

	return &SpanPage{
		Spans:      spans,
		NextCursor: cursorFromPagination(spansMessage.GetNextPage()),
	}, nil
}

func (c jobClient) listTasksPage(ctx context.Context, jobID uuid.UUID, opts *job.TaskPageOptions) (*JobTaskPage, error) {
	var prefetchChildren *workflowsv1.JobTaskChildrenPrefetch
	if opts.PrefetchChildrenLimit > 0 {
		prefetchChildren = workflowsv1.JobTaskChildrenPrefetch_builder{
			Limit: opts.PrefetchChildrenLimit,
		}.Build()
	}

	response, err := c.service.ListJobTasks(
		ctx,
		jobID,
		taskIDToProto(opts.ParentTaskID),
		paginationFromOptions(opts.Limit, opts.Cursor),
		prefetchChildren,
	)
	if err != nil {
		return nil, err
	}

	page := protoToJobTaskPage(response.GetPage())
	for _, childPage := range response.GetPrefetchedChildPages() {
		page.PrefetchedChildPages = append(page.PrefetchedChildPages, protoToJobTaskPage(childPage))
	}
	return page, nil
}

func (c jobClient) queryPage(ctx context.Context, opts *job.QueryOptions) (*JobPage, error) {
	sortDirection := jobQuerySortDirectionToProto(opts.SortDirection)

	var timeInterval *tileboxv1.TimeInterval
	var idInterval *tileboxv1.IDInterval
	if opts.TemporalExtent != nil {
		timeInterval = opts.TemporalExtent.ToProtoTimeInterval()
		idInterval = opts.TemporalExtent.ToProtoIDInterval()

		if timeInterval == nil && idInterval == nil {
			return nil, errors.New("invalid temporal extent")
		}
	}

	automationIDs := lo.Map(opts.AutomationIDs, func(id uuid.UUID, _ int) *tileboxv1.ID {
		return tileboxv1.NewUUID(id)
	})

	states := lo.Map(opts.States, func(state job.State, _ int) workflowsv1.JobState {
		return workflowsv1.JobState(state)
	})
	taskStates := lo.Map(opts.TaskStates, func(state job.TaskState, _ int) workflowsv1.TaskState {
		return workflowsv1.TaskState(state)
	})

	filters := workflowsv1.QueryFilters_builder{
		TimeInterval:  timeInterval,
		IdInterval:    idInterval,
		AutomationIds: automationIDs,
		ClusterSlugs:  opts.ClusterSlugs,
		States:        states,
		TaskStates:    taskStates,
		Name:          opts.Name,
	}.Build()

	jobsMessage, err := c.service.QueryJobs(ctx, filters, paginationFromOptions(opts.Limit, opts.Cursor), sortDirection)
	if err != nil {
		return nil, err
	}

	jobs := make([]*Job, 0, len(jobsMessage.GetJobs()))
	for _, jobMessage := range jobsMessage.GetJobs() {
		jobs = append(jobs, protoToJob(jobMessage))
	}

	return &JobPage{
		Jobs:       jobs,
		NextCursor: cursorFromPagination(jobsMessage.GetNextPage()),
	}, nil
}

func (c jobClient) queryLogsPage(ctx context.Context, jobID uuid.UUID, opts *job.LogQueryOptions) (*LogPage, error) {
	var filters *workflowsv1.LogQueryFilters
	if len(opts.SeverityGroups) > 0 {
		severityLevels := make([]workflowsv1.LogSeverityGroup, len(opts.SeverityGroups))
		for index, group := range opts.SeverityGroups {
			severityLevels[index] = workflowsv1.LogSeverityGroup(group)
		}
		filters = workflowsv1.LogQueryFilters_builder{SeverityLevels: severityLevels}.Build()
	}

	logsMessage, err := c.telemetryService.QueryJobLogs(ctx, jobID, taskIDToProto(opts.TaskID), paginationFromOptions(opts.Limit, opts.Cursor), sortDirectionToProto(opts.SortDirection), filters)
	if err != nil {
		return nil, err
	}

	logs := make([]*LogRecord, 0, len(logsMessage.GetResourceLogs()))
	for _, resourceLogs := range logsMessage.GetResourceLogs() {
		logRecord, err := protoToLogRecord(resourceLogs)
		if err != nil {
			return nil, err
		}
		logs = append(logs, logRecord)
	}

	return &LogPage{
		Logs:       logs,
		NextCursor: cursorFromPagination(logsMessage.GetNextPage()),
	}, nil
}

func validateJob(jobName string, clusterSlug string, maxRetries int64, tasks ...Task) (*workflowsv1.SubmitJobRequest, error) {
	if len(tasks) == 0 {
		return nil, errors.New("no tasks to submit")
	}

	submissions := make([]*futureTask, 0, len(tasks))

	for _, task := range tasks {
		var input []byte
		var err error

		identifier := identifierFromTask(task)
		err = ValidateIdentifier(identifier)
		if err != nil {
			return nil, fmt.Errorf("task has invalid task identifier: %w", err)
		}

		taskProto, isProtobuf := task.(proto.Message)
		if isProtobuf {
			input, err = proto.Marshal(taskProto)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal protobuf task: %w", err)
			}
		} else {
			input, err = json.Marshal(task)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal task: %w", err)
			}
		}

		submissions = append(submissions, &futureTask{
			clusterSlug: clusterSlug,
			identifier:  identifier,
			input:       input,
			maxRetries:  maxRetries,
		})
	}

	return workflowsv1.SubmitJobRequest_builder{
		Tasks:   mergeFutureTasksToSubmissions(submissions),
		JobName: jobName,
	}.Build(), nil
}

func protoToJob(jobMessage *workflowsv1.Job) *Job {
	progressIndicators := make([]*ProgressIndicator, 0, len(jobMessage.GetProgress()))
	for _, progress := range jobMessage.GetProgress() {
		progressIndicators = append(progressIndicators, protoToProgressIndicator(progress))
	}

	return &Job{
		ID:             jobMessage.GetId().AsUUID(),
		Name:           jobMessage.GetName(),
		State:          job.State(jobMessage.GetState()),
		SubmittedAt:    jobMessage.GetSubmittedAt().AsTime(),
		AutomationID:   jobMessage.GetAutomationId().AsUUID(),
		Progress:       progressIndicators,
		ExecutionStats: protoToExecutionStats(jobMessage.GetExecutionStats()),
	}
}

func protoToExecutionStats(es *workflowsv1.ExecutionStats) *ExecutionStats {
	if es == nil {
		return nil
	}
	tasksByState := make([]*TaskStateCount, 0, len(es.GetTasksByState()))
	for _, tsc := range es.GetTasksByState() {
		tasksByState = append(tasksByState, &TaskStateCount{
			State: TaskState(tsc.GetState()),
			Count: tsc.GetCount(),
		})
	}
	return &ExecutionStats{
		FirstTaskStartedAt: es.GetFirstTaskStartedAt().AsTime(),
		LastTaskStoppedAt:  es.GetLastTaskStoppedAt().AsTime(),
		ComputeTime:        es.GetComputeTime().AsDuration(),
		ElapsedTime:        es.GetElapsedTime().AsDuration(),
		Parallelism:        es.GetParallelism(),
		TotalTasks:         es.GetTotalTasks(),
		TasksByState:       tasksByState,
	}
}

func protoToJobTaskPage(page *workflowsv1.JobTaskPage) *JobTaskPage {
	tasks := make([]*JobTask, 0, len(page.GetTasks()))
	for _, task := range page.GetTasks() {
		tasks = append(tasks, protoToJobTask(task))
	}

	return &JobTaskPage{
		ParentTaskID: uuidPointerFromProto(page.GetParentTaskId()),
		Tasks:        tasks,
		NextCursor:   cursorFromPagination(page.GetNextPage()),
	}
}

func protoToJobTask(task *workflowsv1.TaskSummary) *JobTask {
	var submittedAt time.Time
	if timestamp := task.GetSubmittedAt(); timestamp != nil {
		submittedAt = timestamp.AsTime()
	}
	var startedAt time.Time
	if timestamp := task.GetStartedAt(); timestamp != nil {
		startedAt = timestamp.AsTime()
	}
	var stoppedAt time.Time
	if timestamp := task.GetStoppedAt(); timestamp != nil {
		stoppedAt = timestamp.AsTime()
	}

	return &JobTask{
		ID:          task.GetId().AsUUID(),
		ParentID:    uuidPointerFromProto(task.GetParentId()),
		Display:     task.GetDisplay(),
		State:       TaskState(task.GetState()),
		SubmittedAt: submittedAt,
		StartedAt:   startedAt,
		StoppedAt:   stoppedAt,
		Input:       task.GetInput(),
		RetryCount:  task.GetRetryCount(),
		MaxRetries:  task.GetMaxRetries(),
		HasChildren: task.GetHasChildren(),
		ClusterSlug: task.GetClusterSlug(),
		Optional:    task.GetOptional(),
	}
}

func protoToProgressIndicator(p *workflowsv1.Progress) *ProgressIndicator {
	return &ProgressIndicator{
		Label: p.GetLabel(),
		Total: p.GetTotal(),
		Done:  p.GetDone(),
	}
}

func jobQuerySortDirectionToProto(direction job.SortDirection) tileboxv1.SortDirection {
	protoDirection := sortDirectionToProto(direction)
	if protoDirection == nil {
		return tileboxv1.SortDirection_SORT_DIRECTION_UNSPECIFIED
	}
	return *protoDirection
}

func sortDirectionToProto(direction job.SortDirection) *tileboxv1.SortDirection {
	var protoDirection tileboxv1.SortDirection
	switch direction {
	case job.Ascending:
		protoDirection = tileboxv1.SortDirection_SORT_DIRECTION_ASCENDING
	case job.Descending:
		protoDirection = tileboxv1.SortDirection_SORT_DIRECTION_DESCENDING
	default:
		return nil
	}
	return &protoDirection
}

func paginationFromOptions(limit int64, cursor *job.Cursor) *tileboxv1.Pagination {
	if limit <= 0 && cursor == nil {
		return nil
	}

	var startingAfter *tileboxv1.ID
	if cursor != nil {
		startingAfter = tileboxv1.NewUUID(cursor.StartingAfter())
	}

	if limit <= 0 {
		return tileboxv1.Pagination_builder{
			StartingAfter: startingAfter,
		}.Build()
	}

	return tileboxv1.Pagination_builder{
		Limit:         &limit,
		StartingAfter: startingAfter,
	}.Build()
}

func taskIDToProto(taskID *uuid.UUID) *tileboxv1.ID {
	if taskID == nil {
		return nil
	}
	return tileboxv1.NewUUID(*taskID)
}

func uuidPointerFromProto(id *tileboxv1.ID) *uuid.UUID {
	if id == nil {
		return nil
	}
	value := id.AsUUID()
	return &value
}

func cursorFromPagination(page *tileboxv1.Pagination) *job.Cursor {
	if page == nil || page.GetStartingAfter() == nil {
		return nil
	}
	return job.NewCursor(page.GetStartingAfter().AsUUID())
}

func protoToLogRecord(resourceLogs *logsv1.ResourceLogs) (*LogRecord, error) {
	if len(resourceLogs.GetScopeLogs()) == 0 || len(resourceLogs.GetScopeLogs()[0].GetLogRecords()) == 0 {
		return nil, errors.New("resource logs does not contain a log record")
	}

	logRecord := resourceLogs.GetScopeLogs()[0].GetLogRecords()[0]
	bodyString := ""
	if logRecord.GetBody() != nil {
		body := anyValueToGo(logRecord.GetBody())
		var ok bool
		bodyString, ok = body.(string)
		if !ok {
			bodyString = fmt.Sprint(body)
		}
	}

	return &LogRecord{
		Time:             timeFromUnixNanos(logRecord.GetTimeUnixNano()),
		SeverityNumber:   int(logRecord.GetSeverityNumber()),
		SeverityText:     logRecord.GetSeverityText(),
		Body:             bodyString,
		TraceID:          hex.EncodeToString(logRecord.GetTraceId()),
		SpanID:           hex.EncodeToString(logRecord.GetSpanId()),
		Attributes:       keyValuesToMap(logRecord.GetAttributes()),
		RunnerAttributes: keyValuesToMap(resourceLogs.GetResource().GetAttributes()),
	}, nil
}

func protoToSpan(resourceSpans *tracev1.ResourceSpans) (*Span, error) {
	if len(resourceSpans.GetScopeSpans()) == 0 || len(resourceSpans.GetScopeSpans()[0].GetSpans()) == 0 {
		return nil, errors.New("resource spans does not contain a span")
	}

	span := resourceSpans.GetScopeSpans()[0].GetSpans()[0]
	events := make([]*SpanEvent, 0, len(span.GetEvents()))
	for _, event := range span.GetEvents() {
		events = append(events, &SpanEvent{
			Time:       timeFromUnixNanos(event.GetTimeUnixNano()),
			Name:       event.GetName(),
			Attributes: keyValuesToMap(event.GetAttributes()),
		})
	}

	return &Span{
		StartTime:        timeFromUnixNanos(span.GetStartTimeUnixNano()),
		EndTime:          timeFromUnixNanos(span.GetEndTimeUnixNano()),
		TraceID:          hex.EncodeToString(span.GetTraceId()),
		SpanID:           hex.EncodeToString(span.GetSpanId()),
		ParentSpanID:     hex.EncodeToString(span.GetParentSpanId()),
		Name:             span.GetName(),
		StatusCode:       span.GetStatus().GetCode().String(),
		StatusMessage:    span.GetStatus().GetMessage(),
		Attributes:       keyValuesToMap(span.GetAttributes()),
		RunnerAttributes: keyValuesToMap(resourceSpans.GetResource().GetAttributes()),
		Events:           events,
	}, nil
}

func keyValuesToMap(keyValues []*commonv1.KeyValue) map[string]any {
	attributes := make(map[string]any, len(keyValues))
	for _, keyValue := range keyValues {
		attributes[keyValue.GetKey()] = anyValueToGo(keyValue.GetValue())
	}
	return attributes
}

func anyValueToGo(value *commonv1.AnyValue) any {
	if value == nil {
		return nil
	}

	switch value := value.GetValue().(type) {
	case *commonv1.AnyValue_StringValue:
		return value.StringValue
	case *commonv1.AnyValue_BoolValue:
		return value.BoolValue
	case *commonv1.AnyValue_IntValue:
		return value.IntValue
	case *commonv1.AnyValue_DoubleValue:
		return value.DoubleValue
	case *commonv1.AnyValue_BytesValue:
		return value.BytesValue
	case *commonv1.AnyValue_ArrayValue:
		values := value.ArrayValue.GetValues()
		items := make([]any, 0, len(values))
		for _, item := range values {
			items = append(items, anyValueToGo(item))
		}
		return items
	case *commonv1.AnyValue_KvlistValue:
		return keyValuesToMap(value.KvlistValue.GetValues())
	default:
		return nil
	}
}

func timeFromUnixNanos(unixNanos uint64) time.Time {
	const maxInt64 = uint64(1<<63 - 1)
	if unixNanos > maxInt64 {
		return time.Unix(0, int64(maxInt64)).UTC()
	}

	return time.Unix(0, int64(unixNanos)).UTC()
}

// Collect converts any sequence into a slice.
//
// It returns an error if any of the elements in the sequence has a non-nil error.
func Collect[K any](seq iter.Seq2[K, error]) ([]K, error) {
	s := make([]K, 0)

	for k, err := range seq {
		if err != nil {
			return nil, err
		}
		s = append(s, k)
	}
	return s, nil
}
