package job // import "github.com/tilebox/tilebox-go/workflows/v1/job"

import (
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

// QueryOptions contains the configuration for a Query request.
type QueryOptions struct {
	TemporalExtent query.TemporalExtent
	// AutomationID is kept for backward compatibility. Use AutomationIDs instead.
	AutomationID  uuid.UUID
	AutomationIDs []uuid.UUID
	States        []State
	Name          string
}

type QueryOption func(*QueryOptions)

// WithTemporalExtent specifies the time interval for which jobs should be queried.
// Right now, a temporal extent is required for every query.
func WithTemporalExtent(temporalExtent query.TemporalExtent) QueryOption {
	return func(cfg *QueryOptions) {
		cfg.TemporalExtent = temporalExtent
	}
}

// WithAutomationIDs specifies multiple automation IDs to filter jobs by.
// Only jobs submitted by any of the specified automations will be returned.
func WithAutomationIDs(automationIDs ...uuid.UUID) QueryOption {
	return func(cfg *QueryOptions) {
		cfg.AutomationIDs = append(cfg.AutomationIDs, automationIDs...)
	}
}

// WithJobStates filters jobs by their state.
// Only jobs in any of the specified states will be returned.
func WithJobStates(states ...State) QueryOption {
	return func(cfg *QueryOptions) {
		cfg.States = append(cfg.States, states...)
	}
}

// WithJobName filters jobs by name.
func WithJobName(name string) QueryOption {
	return func(cfg *QueryOptions) {
		cfg.Name = name
	}
}
