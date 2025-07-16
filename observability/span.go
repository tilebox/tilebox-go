package observability // import "github.com/tilebox/tilebox-go/observability"

import (
	"context"
	"crypto/rand"
	"fmt"
	"log/slog"

	workflowsv1 "github.com/tilebox/tilebox-go/protogen/workflows/v1"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// Service represents the service being monitored.
type Service struct {
	Name    string // Name represents the logical name of the service.
	Version string // Version represents the version of the service, e.g. "2.0.0".
}

var propagator = propagation.TraceContext{}

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
