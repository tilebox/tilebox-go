package job // import "github.com/tilebox/tilebox-go/workflows/v1/job"

import (
	"github.com/google/uuid"
	"github.com/tilebox/tilebox-go/query"
)

// SubmitOptions contains the configuration for a Submit request.
type SubmitOptions struct {
	MaxRetries int64
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

// QueryOptions contains the configuration for a Query request.
type QueryOptions struct {
	TemporalExtent query.TemporalExtent
	AutomationID   uuid.UUID
}

type QueryOption func(*QueryOptions)

// WithTemporalExtent specifies the time interval for which jobs should be queried.
// Right now, a temporal extent is required for every query.
func WithTemporalExtent(temporalExtent query.TemporalExtent) QueryOption {
	return func(cfg *QueryOptions) {
		cfg.TemporalExtent = temporalExtent
	}
}

// WithAutomationID specifies the automation ID to filter jobs by.
func WithAutomationID(automationID uuid.UUID) QueryOption {
	return func(cfg *QueryOptions) {
		cfg.AutomationID = automationID
	}
}
