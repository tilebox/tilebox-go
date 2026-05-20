package job

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tilebox/tilebox-go/query"
)

func TestStateStringAndMarshalJSON(t *testing.T) {
	tests := []struct {
		name  string
		state State
		want  string
	}{
		{name: "submitted", state: Submitted, want: "submitted"},
		{name: "running", state: Running, want: "running"},
		{name: "started", state: Started, want: "started"},
		{name: "completed", state: Completed, want: "completed"},
		{name: "failed", state: Failed, want: "failed"},
		{name: "canceled", state: Canceled, want: "canceled"},
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

func TestTaskStateStringAndMarshalJSON(t *testing.T) {
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

func Test_SubmitJobOptions(t *testing.T) {
	tests := []struct {
		name    string
		options []SubmitOption
		want    SubmitOptions
	}{
		{
			name: "with max retries",
			options: []SubmitOption{
				WithMaxRetries(7),
			},
			want: SubmitOptions{
				MaxRetries: 7,
			},
		},
		{
			name: "with cluster slug",
			options: []SubmitOption{
				WithClusterSlug("slug"),
			},
			want: SubmitOptions{
				ClusterSlug: "slug",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var options SubmitOptions
			for _, option := range tt.options {
				option(&options)
			}
			assert.Equal(t, tt.want, options)
		})
	}
}

func Test_QueryOptions(t *testing.T) {
	now := time.Now()
	automationID := uuid.New()
	anotherAutomationID := uuid.New()
	cursor := NewCursor(uuid.New())

	tests := []struct {
		name    string
		options []QueryOption
		want    QueryOptions
	}{
		{
			name: "with temporal extent",
			options: []QueryOption{
				WithTemporalExtent(query.NewTimeInterval(
					now,
					now.Add(time.Hour),
				)),
			},
			want: QueryOptions{
				TemporalExtent: query.NewTimeInterval(
					now,
					now.Add(time.Hour),
				),
			},
		},
		{
			name: "with automation id",
			options: []QueryOption{
				WithAutomationIDs(automationID),
				WithAutomationIDs(anotherAutomationID),
			},
			want: QueryOptions{
				AutomationIDs: []uuid.UUID{automationID, anotherAutomationID},
			},
		},
		{
			name: "with job states",
			options: []QueryOption{
				WithJobStates(Submitted, Running),
				WithJobStates(Started),
			},
			want: QueryOptions{
				States: []State{Submitted, Running, Started},
			},
		},
		{
			name: "with task states",
			options: []QueryOption{
				WithTaskStates(TaskQueued, TaskRunning),
				WithTaskStates(TaskComputed),
			},
			want: QueryOptions{
				TaskStates: []TaskState{TaskQueued, TaskRunning, TaskComputed},
			},
		},
		{
			name: "with job name",
			options: []QueryOption{
				WithJobName("my-job"),
			},
			want: QueryOptions{
				Name: "my-job",
			},
		},
		{
			name: "with cursor",
			options: []QueryOption{
				WithCursor(cursor),
			},
			want: QueryOptions{
				Cursor: cursor,
			},
		},
		{
			name: "with limit",
			options: []QueryOption{
				WithLimit(10),
			},
			want: QueryOptions{
				Limit: 10,
			},
		},
		{
			name: "with sort direction",
			options: []QueryOption{
				WithSortDirection(Descending),
			},
			want: QueryOptions{
				SortDirection: Descending,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var options QueryOptions
			for _, option := range tt.options {
				option.ApplyQueryOption(&options)
			}
			assert.Equal(t, tt.want, options)
		})
	}
}

func Test_TelemetryQueryOptions(t *testing.T) {
	cursor := NewCursor(uuid.New())

	tests := []struct {
		name    string
		options []TelemetryQueryOption
		want    TelemetryQueryOptions
	}{
		{
			name: "with cursor",
			options: []TelemetryQueryOption{
				WithCursor(cursor),
			},
			want: TelemetryQueryOptions{
				Cursor: cursor,
			},
		},
		{
			name: "with limit",
			options: []TelemetryQueryOption{
				WithLimit(10),
			},
			want: TelemetryQueryOptions{
				Limit: 10,
			},
		},
		{
			name: "with sort direction",
			options: []TelemetryQueryOption{
				WithSortDirection(Ascending),
			},
			want: TelemetryQueryOptions{
				SortDirection: Ascending,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var options TelemetryQueryOptions
			for _, option := range tt.options {
				option.ApplyTelemetryQueryOption(&options)
			}
			assert.Equal(t, tt.want, options)
		})
	}
}

func TestCursor(t *testing.T) {
	id := uuid.New()
	cursor := NewCursor(id)

	assert.Equal(t, id, cursor.StartingAfter())
	assert.Equal(t, id.String(), cursor.String())

	parsed, err := ParseCursor(cursor.String())
	require.NoError(t, err)
	assert.Equal(t, id, parsed.StartingAfter())
}
