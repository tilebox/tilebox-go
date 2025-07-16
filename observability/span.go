package observability // import "github.com/tilebox/tilebox-go/observability"

import (
	"context"

	"go.opentelemetry.io/otel/codes"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// Service represents the service being monitored.
type Service struct {
	Name    string // Name represents the logical name of the service.
	Version string // Version represents the version of the service, e.g. "2.0.0".
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
