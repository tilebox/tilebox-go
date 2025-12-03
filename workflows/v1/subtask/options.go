package subtask // import "github.com/tilebox/tilebox-go/workflows/v1/subtask"

import "fmt"

type FutureTask uint32

func (ft FutureTask) String() string {
	return fmt.Sprintf("FutureTask(%d)", ft)
}

// SubmitOptions contains the configuration for a SubmitSubtask request.
type SubmitOptions struct {
	Dependencies []FutureTask
	ClusterSlug  string
	MaxRetries   int64
}

// SubmitOption is an interface for configuring a SubmitSubtask request.
type SubmitOption func(*SubmitOptions)

// WithDependencies sets the dependencies of the task.
func WithDependencies(dependencies ...FutureTask) SubmitOption {
	return func(cfg *SubmitOptions) {
		cfg.Dependencies = append(cfg.Dependencies, dependencies...)
	}
}

// WithClusterSlug sets the cluster slug of the cluster where the task will be executed.
//
// Defaults to the cluster of the task runner.
func WithClusterSlug(clusterSlug string) SubmitOption {
	return func(cfg *SubmitOptions) {
		cfg.ClusterSlug = clusterSlug
	}
}

// WithMaxRetries sets the maximum number of times a task can be automatically retried.
//
// Defaults to 0.
func WithMaxRetries(maxRetries int64) SubmitOption {
	return func(cfg *SubmitOptions) {
		cfg.MaxRetries = maxRetries
	}
}
