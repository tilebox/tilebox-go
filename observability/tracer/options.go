package tracer // import "github.com/tilebox/tilebox-go/observability/tracer"

import (
	"time"
)

// Options contain the configuration for tracing.
type Options struct {
	Endpoint       string
	Headers        map[string]string
	ExportInterval time.Duration
}

type Option func(*Options)

// WithEndpointURL sets the target endpoint URL (scheme, host, port, path) the Exporter will connect to.
func WithEndpointURL(endpoint string) Option {
	return func(cfg *Options) {
		cfg.Endpoint = endpoint
	}
}

// WithHeaders sets headers to send on each HTTP request.
func WithHeaders(headers map[string]string) Option {
	return func(cfg *Options) {
		cfg.Headers = headers
	}
}

// WithExportInterval sets the maximum duration between batched exports.
// If set to 0, traces will not be batched.
//
// Defaults to 5s.
func WithExportInterval(exportInterval time.Duration) Option {
	return func(cfg *Options) {
		cfg.ExportInterval = exportInterval
	}
}
