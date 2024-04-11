package observability

import (
	"context"
	"crypto/rand"
	"fmt"
	"log/slog"

	adapter "github.com/axiomhq/axiom-go/adapters/slog"
	"github.com/axiomhq/axiom-go/axiom"
	axiotel "github.com/axiomhq/axiom-go/axiom/otel"
	workflowsv1 "github.com/tilebox/tilebox-go/protogen/go/workflows/v1"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

var propagator = propagation.TraceContext{}

// AxiomLogHandler returns an Axiom handler for slog.
func AxiomLogHandler(dataset, token string, level slog.Level) (*adapter.Handler, error) {
	client, err := axiom.NewClient(axiom.SetToken(token))
	if err != nil {
		return nil, err
	}

	return adapter.New(
		adapter.SetDataset(dataset),
		adapter.SetClient(client),
		adapter.SetLevel(level),
	)
}

// AxiomTraceExporter returns an Axiom OpenTelemetry trace exporter.
func AxiomTraceExporter(ctx context.Context, dataset, token string) (trace.SpanExporter, error) {
	return axiotel.TraceExporter(ctx, dataset, axiotel.SetToken(token))
}

func SetupOtelTracing(serviceName, serviceVersion string, exporters ...trace.SpanExporter) func(ctx context.Context) {
	tp := tracerProvider(serviceName, serviceVersion, exporters)
	otel.SetTracerProvider(tp)

	shutDownFunc := func(ctx context.Context) {
		_ = tp.Shutdown(ctx)
	}

	return shutDownFunc
}

// tracerProvider configures and returns a new OpenTelemetry tracer provider.
func tracerProvider(serviceName, serviceVersion string, exporters []trace.SpanExporter) *trace.TracerProvider {
	rs := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(serviceName),
		semconv.ServiceVersionKey.String(serviceVersion),
	)

	opts := []trace.TracerProviderOption{
		trace.WithResource(rs),
	}

	for _, exporter := range exporters {
		opts = append(opts, trace.WithBatcher(exporter, trace.WithMaxQueueSize(10*1024)))
	}

	return trace.NewTracerProvider(opts...)
}

// generateTraceParent generates a random traceparent.
// ex: 00-51651128eabd3c6f7695b9e6ee0e337d-03e60705fdf5efe0-01
func generateTraceParent() (string, error) {
	traceID := oteltrace.TraceID{}
	_, err := rand.Read(traceID[:])
	if err != nil {
		return "", err
	}

	spanID := oteltrace.SpanID{}
	_, err = rand.Read(spanID[:])
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("00-%s-%s-01", traceID, spanID), nil
}

func GetTraceParentOfCurrentSpan(ctx context.Context) string {
	carrier := propagation.MapCarrier{}

	propagator.Inject(ctx, carrier)

	traceparent := carrier.Get("traceparent")
	if traceparent == "" {
		// If the tracer is from a NoopTracerProvider, it will not generate a traceparent.
		// In this case, we generate a random one.
		traceparent, err := generateTraceParent()
		if err != nil {
			slog.ErrorContext(ctx, "failed to generate traceparent", "error", err)
			return ""
		}
		return traceparent
	}
	return traceparent
}

func StartJobSpan[Result any](ctx context.Context, tracer oteltrace.Tracer, spanName string, job *workflowsv1.Job, f func(ctx context.Context) (Result, error)) (Result, error) {
	carrier := propagation.MapCarrier{"traceparent": job.GetTraceParent()}
	ctx = propagator.Extract(ctx, carrier)

	return WithSpanResult(ctx, tracer, spanName, f)
}

// WithSpanResult wraps a function call that returns a result and an error with a tracing span of the given name
func WithSpanResult[Result any](ctx context.Context, tracer oteltrace.Tracer, name string, f func(ctx context.Context) (Result, error)) (Result, error) {
	ctx, span := tracer.Start(ctx, name)
	defer span.End()

	result, err := f(ctx)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	return result, err
}
