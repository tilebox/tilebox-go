package job // import "github.com/tilebox/tilebox-go/workflows/v1/job"

// SubmitJobOptions contains the configuration for Tilebox Workflows Task Runner.
type SubmitJobOptions struct {
	MaxRetries int64
}

type SubmitOption func(*SubmitJobOptions)

// WithMaxRetries sets the maximum number of times a job can be automatically retried.
//
// Defaults to 0.
func WithMaxRetries(maxRetries int64) SubmitOption {
	return func(cfg *SubmitJobOptions) {
		cfg.MaxRetries = maxRetries
	}
}
