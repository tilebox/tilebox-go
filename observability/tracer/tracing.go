package tracer // import "github.com/tilebox/tilebox-go/observability/tracer"

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.30.0"
)

const axiomTracesEndpoint = "https://api.axiom.co/v1/traces"

// Service represents the service being monitored.
type Service struct {
	Name    string // Name represents the logical name of the service.
	Version string // Version represents the version of the service, e.g. "2.0.0".
}

// NewOtelSpanProcessor creates a new OpenTelemetry HTTP span processor.
func NewOtelSpanProcessor(ctx context.Context, options ...Option) (trace.SpanProcessor, error) {
	opts := &Options{
		ExportInterval: trace.DefaultScheduleDelay,
	}
	for _, option := range options {
		option(opts)
	}

	var exporterOptions []otlptracehttp.Option
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
func NewOtelProviderWithProcessors(otelService *Service, processors ...trace.SpanProcessor) *trace.TracerProvider {
	rs := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(otelService.Name),
		semconv.ServiceVersionKey.String(otelService.Version),
	)

	opts := []trace.TracerProviderOption{
		trace.WithResource(rs),
		// Ensure all traces are sampled. It a good default when generating little data (< 100 traces per second).
		// https://opentelemetry.io/docs/concepts/sampling/#when-not-to-sample
		trace.WithSampler(trace.AlwaysSample()),
	}

	for _, processor := range processors {
		opts = append(opts, trace.WithSpanProcessor(processor))
	}

	return trace.NewTracerProvider(opts...)
}

func noShutdown(context.Context) {}

// NewOtelProvider creates a new OpenTelemetry tracer provider.
func NewOtelProvider(ctx context.Context, otelService *Service, options ...Option) (*trace.TracerProvider, func(ctx context.Context), error) {
	processor, err := NewOtelSpanProcessor(ctx, options...)
	if err != nil {
		return nil, noShutdown, err
	}

	tracerProvider := NewOtelProviderWithProcessors(otelService, processor)

	return tracerProvider, func(ctx context.Context) {
		_ = tracerProvider.Shutdown(ctx)
	}, nil
}

// NewAxiomProvider creates a new OpenTelemetry tracer provider that sends traces to Axiom.
func NewAxiomProvider(ctx context.Context, otelService *Service, dataset string, apiKey string) (*trace.TracerProvider, func(ctx context.Context), error) {
	headers := map[string]string{
		"Authorization":   fmt.Sprintf("Bearer %s", apiKey),
		"X-Axiom-Dataset": dataset,
	}

	return NewOtelProvider(ctx, otelService, WithEndpointURL(axiomTracesEndpoint), WithHeaders(headers))
}
