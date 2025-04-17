package observability // import "github.com/tilebox/tilebox-go/observability"

import (
	"log/slog"
	"time"
)

// Options contain the configuration for observability.
type Options struct {
	Level          slog.Level
	Endpoint       string
	Headers        map[string]string
	ExportInterval time.Duration
}

type Option func(*Options)

// WithLevel sets the log level.
//
// Defaults to slog.LevelInfo.
func WithLevel(level slog.Level) Option {
	return func(cfg *Options) {
		cfg.Level = level
	}
}

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
// If set to 0, records will not be batched.
//
// Defaults to 5s for tracing and 1s for logs.
func WithExportInterval(exportInterval time.Duration) Option {
	return func(cfg *Options) {
		cfg.ExportInterval = exportInterval
	}
}
