package tracer // import "github.com/tilebox/tilebox-go/observability/tracer"

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/tilebox/tilebox-go/observability"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.30.0"
)

const axiomTracesEndpoint = "https://api.axiom.co/v1/traces"

// NewOtelSpanProcessor creates a new OpenTelemetry HTTP span processor.
func NewOtelSpanProcessor(ctx context.Context, options ...Option) (trace.SpanProcessor, error) {
	opts := &Options{
		ExportInterval: trace.DefaultScheduleDelay * time.Millisecond,
	}
	for _, option := range options {
		option(opts)
	}

	exporterOptions := make([]otlptracehttp.Option, 0, 2)
	if opts.Endpoint != "" {
		exporterOptions = append(exporterOptions, otlptracehttp.WithEndpointURL(opts.Endpoint))
	}
	if len(opts.Headers) > 0 {
		exporterOptions = append(exporterOptions, otlptracehttp.WithHeaders(opts.Headers))
	}

	exporter, err := otlptracehttp.New(ctx, exporterOptions...)
	if err != nil {
		return nil, err
	}

	if opts.ExportInterval > 0 {
		return trace.NewBatchSpanProcessor(exporter,
			trace.WithMaxQueueSize(10*1024),
			trace.WithBatchTimeout(opts.ExportInterval),
		), nil
	}

	return trace.NewSimpleSpanProcessor(exporter), nil
}

// NewOtelProviderWithProcessors creates a new OpenTelemetry tracer provider with the given processors.
func NewOtelProviderWithProcessors(otelService *observability.Service, processors ...trace.SpanProcessor) *trace.TracerProvider {
	rs := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(otelService.Name),
		semconv.ServiceVersionKey.String(otelService.Version),
	)

	opts := make([]trace.TracerProviderOption, 0, len(processors)+2)
	opts = append(opts, trace.WithResource(rs))
	// Ensure all traces are sampled. It a good default when generating little data (< 100 traces per second).
	// https://opentelemetry.io/docs/concepts/sampling/#when-not-to-sample
	opts = append(opts, trace.WithSampler(trace.AlwaysSample()))

	for _, processor := range processors {
		opts = append(opts, trace.WithSpanProcessor(processor))
	}

	return trace.NewTracerProvider(opts...)
}

func noShutdown(context.Context) {}

// NewOtelProvider creates a new OpenTelemetry tracer provider.
func NewOtelProvider(ctx context.Context, otelService *observability.Service, options ...Option) (*trace.TracerProvider, func(ctx context.Context), error) {
	processor, err := NewOtelSpanProcessor(ctx, options...)
	if err != nil {
		return nil, noShutdown, err
	}

	tracerProvider := NewOtelProviderWithProcessors(otelService, processor)

	return tracerProvider, func(ctx context.Context) {
		_ = tracerProvider.Shutdown(ctx)
	}, nil
}

// NewAxiomProviderFromEnv creates a new OpenTelemetry tracer provider that sends traces to Axiom.
// It reads the following environment variables:
//
//   - AXIOM_API_KEY: The Axiom API key.
//   - AXIOM_TRACES_DATASET: The dataset where traces should be stored.
func NewAxiomProviderFromEnv(ctx context.Context, otelService *observability.Service, options ...Option) (*trace.TracerProvider, func(ctx context.Context), error) {
	apiKey := os.Getenv("AXIOM_API_KEY")
	if apiKey == "" {
		return nil, noShutdown, errors.New("AXIOM_API_KEY environment variable is not set")
	}

	dataset := os.Getenv("AXIOM_TRACES_DATASET")
	if dataset == "" {
		return nil, noShutdown, errors.New("AXIOM_TRACES_DATASET environment variable is not set")
	}

	return NewAxiomProvider(ctx, otelService, dataset, apiKey, options...)
}

// NewAxiomProvider creates a new OpenTelemetry tracer provider that sends traces to Axiom to the given dataset and authenticating with the given API key.
func NewAxiomProvider(ctx context.Context, otelService *observability.Service, dataset, apiKey string, options ...Option) (*trace.TracerProvider, func(ctx context.Context), error) {
	headers := map[string]string{
		"Authorization":   fmt.Sprintf("Bearer %s", apiKey),
		"X-Axiom-Dataset": dataset,
	}

	opts := make([]Option, 0, len(options)+2)
	opts = append(opts, WithEndpointURL(axiomTracesEndpoint), WithHeaders(headers))
	opts = append(opts, options...)

	return NewOtelProvider(ctx, otelService, opts...)
}
