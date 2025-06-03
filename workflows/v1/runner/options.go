package runner // import "github.com/tilebox/tilebox-go/workflows/v1/runner"

import (
	"log/slog"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
)

// Options contains the configuration for Tilebox Workflows Task Runner.
type Options struct {
	ClusterSlug   string
	Logger        *slog.Logger
	MeterProvider metric.MeterProvider
}

type Option func(*Options)

// WithClusterSlug sets the cluster to listen tasks to.
//
// Defaults to the default cluster.
func WithClusterSlug(clusterSlug string) Option {
	return func(cfg *Options) {
		cfg.ClusterSlug = clusterSlug
	}
}

// WithRunnerLogger sets the logger to use for the task runner.
//
// Defaults to slog.Default().
func WithRunnerLogger(logger *slog.Logger) Option {
	return func(cfg *Options) {
		cfg.Logger = logger
	}
}

// WithDisableMetrics disables OpenTelemetry metrics for the task runner.
func WithDisableMetrics() Option {
	return func(cfg *Options) {
		cfg.MeterProvider = noop.NewMeterProvider()
	}
}
