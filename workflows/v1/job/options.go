package job // import "github.com/tilebox/tilebox-go/workflows/v1/job"

import (
	"encoding/json"
	"fmt"

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

// TelemetryQueryOptions contains options for querying job observability.
type TelemetryQueryOptions struct {
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

type sharedQueryOption interface {
	QueryOption
	TelemetryQueryOption
}

type SortDirectionOption = sharedQueryOption

type LimitOption = sharedQueryOption

type CursorOption = sharedQueryOption

// Cursor identifies where to continue a paginated query.
//
// Cursors are returned as NextCursor values from page query methods and should only be reused with the same endpoint,
// filters and sort direction that produced them.
type Cursor struct {
	startingAfter uuid.UUID
}

// NewCursor creates a cursor that starts after the entry with the given ID.
func NewCursor(startingAfter uuid.UUID) *Cursor {
	return &Cursor{startingAfter: startingAfter}
}

// ParseCursor parses a cursor string returned by Cursor.String.
func ParseCursor(value string) (*Cursor, error) {
	startingAfter, err := uuid.Parse(value)
	if err != nil {
		return nil, fmt.Errorf("invalid cursor: %w", err)
	}
	return NewCursor(startingAfter), nil
}

// String returns the cursor as a string that can be passed to ParseCursor.
func (c *Cursor) String() string {
	if c == nil {
		return ""
	}
	return c.startingAfter.String()
}

// StartingAfter returns the ID after which the next page should start.
func (c *Cursor) StartingAfter() uuid.UUID {
	if c == nil {
		return uuid.Nil
	}
	return c.startingAfter
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
