package observability

import (
	"context"
	"crypto/rand"
	"fmt"
	"log/slog"

	adapter "github.com/axiomhq/axiom-go/adapters/slog"
	"github.com/axiomhq/axiom-go/axiom"
	axiotel "github.com/axiomhq/axiom-go/axiom/otel"
	slogotel "github.com/remychantenay/slog-otel"
	workflowsv1 "github.com/tilebox/tilebox-go/protogen/go/workflows/v1"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

var propagator = propagation.TraceContext{}

// NewAxiomLogger returns a slog.Logger that logs to Axiom.
// It also returns a shutdown function that should be called when the logger is no longer needed, to ensure
// all logs are flushed.
func NewAxiomLogger(dataset, token string, level slog.Level, addTracing bool) (slog.Handler, func(), error) {
	noShutdown := func() {}
	client, err := axiom.NewClient(axiom.SetToken(token))
	if err != nil {
		return nil, noShutdown, err
	}

	axiomHandler, err := adapter.New(
		adapter.SetDataset(dataset),
		adapter.SetClient(client),
		adapter.SetLevel(level),
	)
	if err != nil {
		return nil, noShutdown, err
	}

	if !addTracing {
		return axiomHandler, axiomHandler.Close, nil
	}

	return slogotel.OtelHandler{Next: axiomHandler}, axiomHandler.Close, nil
}

func NewAxiomTracerProvider(ctx context.Context, dataset, token, serviceName, serviceVersion string) (oteltrace.TracerProvider, func(), error) {
	noShutdown := func() {}
	exporter, err := axiotel.TraceExporter(ctx, dataset, axiotel.SetToken(token))
	if err != nil {
		return nil, noShutdown, err
	}

	traceResource := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(serviceName),
		semconv.ServiceVersionKey.String(serviceVersion),
	)

	provider := trace.NewTracerProvider(
		trace.WithResource(traceResource),
		trace.WithBatcher(exporter, trace.WithMaxQueueSize(10*1024)),
	)
	shutdown := func() {
		_ = provider.Shutdown(ctx)
	}
	return provider, shutdown, nil
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
			slog.ErrorContext(ctx, "failed to generate traceparent", slog.Any("error", err))
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

// WithSpan wraps a function call with a tracing span of the given name
func WithSpan(ctx context.Context, tracer oteltrace.Tracer, name string, f func(ctx context.Context) error) error {
	ctx, span := tracer.Start(ctx, name)
	defer span.End()

	err := f(ctx)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	return err
}
