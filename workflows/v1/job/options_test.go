package job

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/tilebox/tilebox-go/query"
)

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
			},
			want: QueryOptions{
				AutomationIDs: []uuid.UUID{automationID},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var options QueryOptions
			for _, option := range tt.options {
				option(&options)
			}
			assert.Equal(t, tt.want, options)
		})
	}
}
