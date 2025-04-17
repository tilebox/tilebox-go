package observability // import "github.com/tilebox/tilebox-go/observability"

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
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

// NewAxiomSpanProcessor creates a new span processor that sends traces to Axiom.
func NewAxiomSpanProcessor(ctx context.Context, dataset string, apiKey string) (trace.SpanProcessor, error) {
	headers := map[string]string{
		"Authorization":   fmt.Sprintf("Bearer %s", apiKey),
		"X-Axiom-Dataset": dataset,
	}

	return NewOtelSpanProcessor(ctx, WithEndpointURL(axiomTracesEndpoint), WithHeaders(headers))
}

// NewTracerProvider creates a new OpenTelemetry tracer provider with the given processors.
func NewTracerProvider(otelService *Service, processors ...trace.SpanProcessor) *trace.TracerProvider {
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

// InitializeTracing initializes the OpenTelemetry tracer with the given processors.
// It sets the global tracer provider and returns a shutdown function.
func InitializeTracing(otelService *Service, processors ...trace.SpanProcessor) func(ctx context.Context) {
	tracerProvider := NewTracerProvider(otelService, processors...)
	otel.SetTracerProvider(tracerProvider)

	return func(ctx context.Context) {
		_ = tracerProvider.Shutdown(ctx)
	}
}

// InitializePropagator configures the global OpenTelemetry propagator to support W3C Trace Context and Baggage.
func InitializePropagator() {
	prop := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)

	otel.SetTextMapPropagator(prop)
}
