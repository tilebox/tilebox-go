package job

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/tilebox/tilebox-go/query"
)

// SubmitOptions contains the configuration for a Submit request.
type SubmitOptions struct {
	MaxRetries  int64
	ClusterSlug string
}

type SubmitOption func(*SubmitOptions)

// WithMaxRetries sets the maximum number of times a job can be automatically retried.
//
// Defaults to 0.
func WithMaxRetries(maxRetries int64) SubmitOption {
	return func(cfg *SubmitOptions) {
		cfg.MaxRetries = maxRetries
	}
}

// WithClusterSlug sets the cluster slug of the cluster where the job will be executed.
//
// Defaults to the default cluster.
func WithClusterSlug(clusterSlug string) SubmitOption {
	return func(cfg *SubmitOptions) {
		cfg.ClusterSlug = clusterSlug
	}
}

// State is an alias to the workflows State type for use in query options.
type State int32

// State values for use in query filters.
const (
	// Submitted means the job has been submitted, queued and waiting for its first task to be run.
	Submitted State = 1
	// Running means the job is running, i.e. at least one task is running.
	Running State = 2
	// Started means the job has started running, i.e. at least one task has been computed, but currently no tasks are running.
	Started State = 3
	// Completed means the job has completed successfully.
	Completed State = 4
	// Failed means the job has failed.
	Failed State = 5
	// Canceled means the job has been canceled on user request.
	Canceled State = 6
)

func (state State) String() string {
	switch state {
	case Submitted:
		return "submitted"
	case Running:
		return "running"
	case Started:
		return "started"
	case Completed:
		return "completed"
	case Failed:
		return "failed"
	case Canceled:
		return "canceled"
	default:
		return "unspecified"
	}
}

func (state State) MarshalJSON() ([]byte, error) {
	return json.Marshal(state.String())
}

// TaskState is a task state value for use in query filters.
type TaskState int32

// TaskState values for use in query filters.
const (
	// TaskQueued means the task is queued and waiting to be run.
	TaskQueued TaskState = 1
	// TaskRunning means the task is currently running on some task runner.
	TaskRunning TaskState = 2
	// TaskComputed means the task has been computed and the output is available.
	TaskComputed TaskState = 3
	// TaskFailed means the task has failed.
	TaskFailed TaskState = 4
	// TaskSkipped means the task has been skipped.
	TaskSkipped TaskState = 5
	// TaskFailedOptional means the task has failed, but was marked as optional.
	TaskFailedOptional TaskState = 6
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

// QueryOptions contains the configuration for a Query request.
type QueryOptions struct {
	// TemporalExtent is the time or ID interval for which jobs should be queried.
	// Leave unset to query jobs without a temporal extent filter.
	TemporalExtent query.TemporalExtent
	// Cursor starts the query after a cursor returned by a previous page.
	Cursor *Cursor
	// AutomationIDs filters jobs by the automations that submitted them.
	AutomationIDs []uuid.UUID
	// ClusterSlugs filters jobs by the clusters they were submitted to.
	ClusterSlugs []string
	// States filters jobs by their states.
	States []State
	// TaskStates filters jobs by the states of their tasks.
	TaskStates []TaskState
	// Name filters jobs by their name.
	Name string
	// Limit is the maximum number of jobs to return.
	// Leave unset or set to 0 to paginate through and return all jobs.
	Limit int64
	// SortDirection is the direction in which jobs should be sorted by submission date.
	// Leave unset to let the server choose its default sort direction.
	SortDirection SortDirection
}

type QueryOption interface {
	ApplyQueryOption(*QueryOptions)
}

type queryOptionFunc func(*QueryOptions)

func (f queryOptionFunc) ApplyQueryOption(cfg *QueryOptions) {
	f(cfg)
}

// TaskListOptions contains options for listing one sibling collection in a job's task tree.
type TaskListOptions struct {
	// ParentTaskID selects the task whose children to list. Leave unset to list root tasks.
	ParentTaskID *uuid.UUID
	// Cursor starts the listing after a cursor returned for the same sibling collection.
	Cursor *Cursor
	// Limit is the maximum number of tasks to return.
	// Leave unset or set to 0 to paginate through and return all tasks.
	Limit int64
}

type TaskListOption interface {
	ApplyTaskListOption(*TaskListOptions)
}

// TaskPageOptions contains options for listing one page in a job's task tree.
type TaskPageOptions struct {
	TaskListOptions

	// PrefetchChildrenLimit is the maximum number of children to prefetch for each task in the page.
	// Leave unset or set to 0 to disable child prefetching.
	PrefetchChildrenLimit int64
}

type TaskPageOption interface {
	ApplyTaskPageOption(*TaskPageOptions)
}

// TelemetryQueryOptions contains options for querying job observability.
type TelemetryQueryOptions struct {
	// TaskID filters observability records by task. Leave unset to return records for all tasks in the job.
	TaskID *uuid.UUID
	// Cursor starts the query after a cursor returned by a previous page.
	Cursor *Cursor
	// Limit is the maximum number of observability records to return.
	// Leave unset or set to 0 to paginate through and return all records.
	Limit int64
	// SortDirection is the direction in which observability records should be sorted.
	// Leave unset to let the server choose its default sort direction.
	SortDirection SortDirection
}

type TelemetryQueryOption interface {
	ApplyTelemetryQueryOption(*TelemetryQueryOptions)
}

// LogQueryOptions contains options for querying job logs.
type LogQueryOptions struct {
	TelemetryQueryOptions

	// SeverityGroups filters logs by severity group. Leave unset to return logs of all severities.
	SeverityGroups []LogSeverityGroup
}

type LogQueryOption interface {
	ApplyLogQueryOption(*LogQueryOptions)
}

type sharedQueryOption interface {
	QueryOption
	TelemetryQueryOption
	LogQueryOption
}

type SortDirectionOption = sharedQueryOption

type paginationOption interface {
	sharedQueryOption
	TaskListOption
	TaskPageOption
}

type LimitOption = paginationOption

type CursorOption = paginationOption

type taskCollectionOption interface {
	TaskListOption
	TaskPageOption
}

type ParentTaskIDOption = taskCollectionOption

type telemetryQueryOption interface {
	TelemetryQueryOption
	LogQueryOption
}

type TaskIDOption = telemetryQueryOption

// Cursor identifies where to continue a paginated query.
//
// Cursors are returned as NextCursor values from page query methods and should only be reused with the same endpoint,
// filters and sort direction that produced them.
type Cursor = query.Cursor

// NewCursor creates a cursor that starts after the entry with the given ID.
func NewCursor(startingAfter uuid.UUID) *Cursor {
	return query.NewCursor(startingAfter)
}

// ParseCursor parses a cursor string returned by Cursor.String.
func ParseCursor(value string) (*Cursor, error) {
	return query.ParseCursor(value)
}

type cursorOption struct {
	cursor *Cursor
}

func (o cursorOption) ApplyQueryOption(cfg *QueryOptions) {
	cfg.Cursor = o.cursor
}

func (o cursorOption) ApplyTelemetryQueryOption(cfg *TelemetryQueryOptions) {
	cfg.Cursor = o.cursor
}

func (o cursorOption) ApplyLogQueryOption(cfg *LogQueryOptions) {
	cfg.Cursor = o.cursor
}

func (o cursorOption) ApplyTaskListOption(cfg *TaskListOptions) {
	cfg.Cursor = o.cursor
}

func (o cursorOption) ApplyTaskPageOption(cfg *TaskPageOptions) {
	cfg.Cursor = o.cursor
}

type limitOption struct {
	limit int64
}

func (o limitOption) ApplyQueryOption(cfg *QueryOptions) {
	if o.limit > 0 {
		cfg.Limit = o.limit
	}
}

func (o limitOption) ApplyTelemetryQueryOption(cfg *TelemetryQueryOptions) {
	if o.limit > 0 {
		cfg.Limit = o.limit
	}
}

func (o limitOption) ApplyLogQueryOption(cfg *LogQueryOptions) {
	if o.limit > 0 {
		cfg.Limit = o.limit
	}
}

func (o limitOption) ApplyTaskListOption(cfg *TaskListOptions) {
	if o.limit > 0 {
		cfg.Limit = o.limit
	}
}

func (o limitOption) ApplyTaskPageOption(cfg *TaskPageOptions) {
	if o.limit > 0 {
		cfg.Limit = o.limit
	}
}

type sortDirectionOption struct {
	direction SortDirection
}

func (o sortDirectionOption) ApplyQueryOption(cfg *QueryOptions) {
	if o.direction != 0 {
		cfg.SortDirection = o.direction
	}
}

func (o sortDirectionOption) ApplyTelemetryQueryOption(cfg *TelemetryQueryOptions) {
	if o.direction != 0 {
		cfg.SortDirection = o.direction
	}
}

func (o sortDirectionOption) ApplyLogQueryOption(cfg *LogQueryOptions) {
	if o.direction != 0 {
		cfg.SortDirection = o.direction
	}
}

type taskIDOption struct {
	taskID uuid.UUID
}

func (o taskIDOption) ApplyTelemetryQueryOption(cfg *TelemetryQueryOptions) {
	cfg.TaskID = &o.taskID
}

func (o taskIDOption) ApplyLogQueryOption(cfg *LogQueryOptions) {
	cfg.TaskID = &o.taskID
}

type parentTaskIDOption struct {
	parentTaskID uuid.UUID
}

func (o parentTaskIDOption) ApplyTaskListOption(cfg *TaskListOptions) {
	cfg.ParentTaskID = &o.parentTaskID
}

func (o parentTaskIDOption) ApplyTaskPageOption(cfg *TaskPageOptions) {
	cfg.ParentTaskID = &o.parentTaskID
}

type prefetchChildrenOption struct {
	limit int64
}

func (o prefetchChildrenOption) ApplyTaskPageOption(cfg *TaskPageOptions) {
	if o.limit > 0 {
		cfg.PrefetchChildrenLimit = o.limit
	}
}

// LogSeverityGroup is a log severity group for use in query filters.
type LogSeverityGroup int32

// Log severity values for use in query filters.
const (
	_ LogSeverityGroup = iota
	// LogSeverityTrace includes OpenTelemetry trace severities.
	LogSeverityTrace
	// LogSeverityDebug includes OpenTelemetry debug severities.
	LogSeverityDebug
	// LogSeverityInfo includes OpenTelemetry info severities.
	LogSeverityInfo
	// LogSeverityWarning includes OpenTelemetry warning severities.
	LogSeverityWarning
	// LogSeverityError includes OpenTelemetry error and fatal severities.
	LogSeverityError
)

type logSeverityGroupsOption struct {
	groups []LogSeverityGroup
}

func (o logSeverityGroupsOption) ApplyLogQueryOption(cfg *LogQueryOptions) {
	cfg.SeverityGroups = append(cfg.SeverityGroups, o.groups...)
}

// SortDirection specifies the sort direction for job and observability query results.
type SortDirection int32

// SortDirection values.
const (
	_ SortDirection = iota
	// Ascending sorts query results oldest first.
	Ascending
	// Descending sorts query results newest first.
	Descending
)

// WithTemporalExtent specifies the time or ID interval for which jobs should be queried.
func WithTemporalExtent(temporalExtent query.TemporalExtent) QueryOption {
	return queryOptionFunc(func(cfg *QueryOptions) {
		cfg.TemporalExtent = temporalExtent
	})
}

// WithAutomationIDs specifies multiple automation IDs to filter jobs by.
// Only jobs submitted by any of the specified automations will be returned.
func WithAutomationIDs(automationIDs ...uuid.UUID) QueryOption {
	return queryOptionFunc(func(cfg *QueryOptions) {
		cfg.AutomationIDs = append(cfg.AutomationIDs, automationIDs...)
	})
}

// WithClusterSlugs specifies multiple cluster slugs to filter jobs by.
// Only jobs that have tasks on any of the specified clusters will be returned.
func WithClusterSlugs(clusterSlugs ...string) QueryOption {
	return queryOptionFunc(func(cfg *QueryOptions) {
		cfg.ClusterSlugs = append(cfg.ClusterSlugs, clusterSlugs...)
	})
}

// WithJobStates filters jobs by their state.
// Only jobs in any of the specified states will be returned.
func WithJobStates(states ...State) QueryOption {
	return queryOptionFunc(func(cfg *QueryOptions) {
		cfg.States = append(cfg.States, states...)
	})
}

// WithTaskStates filters jobs by task state.
// Only jobs with at least one task in any of the specified states will be returned.
func WithTaskStates(states ...TaskState) QueryOption {
	return queryOptionFunc(func(cfg *QueryOptions) {
		cfg.TaskStates = append(cfg.TaskStates, states...)
	})
}

// WithJobName filters jobs by name.
func WithJobName(name string) QueryOption {
	return queryOptionFunc(func(cfg *QueryOptions) {
		cfg.Name = name
	})
}

// WithCursor starts the query after a cursor returned by a previous page.
func WithCursor(cursor *Cursor) CursorOption {
	return cursorOption{cursor: cursor}
}

// WithLimit limits the number of query results returned.
//
// For auto-paginated query methods, the limit applies to the total number of results yielded. For page query methods,
// the limit applies to the single page returned.
//
// Defaults to unlimited.
func WithLimit(limit int64) LimitOption {
	return limitOption{limit: limit}
}

// WithSortDirection sets the sort direction for job and observability query results.
//
// Defaults to the server default.
func WithSortDirection(direction SortDirection) SortDirectionOption {
	return sortDirectionOption{direction: direction}
}

// WithTaskID filters job logs or spans by task ID.
func WithTaskID(taskID uuid.UUID) TaskIDOption {
	return taskIDOption{taskID: taskID}
}

// WithParentTaskID lists the direct children of the specified task instead of root tasks.
func WithParentTaskID(parentTaskID uuid.UUID) ParentTaskIDOption {
	return parentTaskIDOption{parentTaskID: parentTaskID}
}

// WithPrefetchChildren includes the first page of children for each task in a manually requested task page.
// The limit is applied separately to each child page.
func WithPrefetchChildren(limit int64) TaskPageOption {
	return prefetchChildrenOption{limit: limit}
}

// WithLogSeverityGroups filters logs by severity group. Logs matching any of the specified groups are returned.
// If no groups are specified, logs of all severities are returned.
func WithLogSeverityGroups(groups ...LogSeverityGroup) LogQueryOption {
	return logSeverityGroupsOption{groups: groups}
}
