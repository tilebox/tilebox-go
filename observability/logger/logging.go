package logger // import "github.com/tilebox/tilebox-go/observability/logger"

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
	slogmulti "github.com/samber/slog-multi"
	"github.com/tilebox/tilebox-go/observability"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/sdk/log"
	oteltrace "go.opentelemetry.io/otel/trace"
)

const (
	axiomLogsEndpoint = "https://api.axiom.co/v1/logs"
	otelLoggerName    = "tilebox.com/observability"
)

// NewOtelLogProcessor creates a new OpenTelemetry HTTP log processor.
func NewOtelLogProcessor(ctx context.Context, options ...Option) (log.Processor, error) {
	opts := &Options{
		ExportInterval: 1 * time.Second,
	}
	for _, option := range options {
		option(opts)
	}

	var exporterOptions []otlploghttp.Option
	if opts.Endpoint != "" {
		exporterOptions = append(exporterOptions, otlploghttp.WithEndpointURL(opts.Endpoint))
	}
	if len(opts.Headers) > 0 {
		exporterOptions = append(exporterOptions, otlploghttp.WithHeaders(opts.Headers))
	}

	exporter, err := otlploghttp.New(ctx, exporterOptions...)
	if err != nil {
		return nil, err
	}

	if opts.ExportInterval > 0 {
		return log.NewBatchProcessor(exporter, log.WithExportInterval(opts.ExportInterval)), nil
	}

	return log.NewSimpleProcessor(exporter), nil
}

// NewLoggingProviderWithProcessors creates a new OpenTelemetry logger provider with the given processors.
func NewLoggingProviderWithProcessors(otelService *observability.Service, processors ...log.Processor) *log.LoggerProvider {
	opts := make([]log.LoggerProviderOption, 0, len(processors)+1)
	opts = append(opts, log.WithResource(observability.NewResource(otelService)))

	for _, processor := range processors {
		opts = append(opts, log.WithProcessor(processor))
	}

	return log.NewLoggerProvider(opts...)
}

// tracingAndlevelFilterHandler is a slog.Handler that filters log messages based on the configured level
// and adds tracing information to each log record.
type tracingAndlevelFilterHandler struct {
	level        slog.Level
	next         slog.Handler
	addSpanEvent bool
}

var _ slog.Handler = &tracingAndlevelFilterHandler{}

type contextSlogAttributesKeyType string

const contextSlogAttributesKey contextSlogAttributesKeyType = "x-tilebox-ctx-slog-attributes"

// ContextWithSlogAttributes returns a context whose log records include attrs.
func ContextWithSlogAttributes(ctx context.Context, attributes ...slog.Attr) context.Context {
	if len(attributes) == 0 {
		return ctx
	}

	existing := contextAttributes(ctx)
	contextAttributes := make([]slog.Attr, 0, len(existing)+len(attributes))
	contextAttributes = append(contextAttributes, existing...)
	contextAttributes = append(contextAttributes, attributes...)
	return context.WithValue(ctx, contextSlogAttributesKey, contextAttributes)
}

func contextAttributes(ctx context.Context) []slog.Attr {
	attrs, _ := ctx.Value(contextSlogAttributesKey).([]slog.Attr)
	return attrs
}

// Enabled checks if a log message with the given level should be logged.
func (h *tracingAndlevelFilterHandler) Enabled(ctx context.Context, level slog.Level) bool {
	if level < h.level { // if we are below the configured level, we for sure don't want to log this
		return false
	}
	// otherwise, we want to log if the underlying handler is enabled
	return h.next.Enabled(ctx, level)
}

// Handle processes a log record.
func (h *tracingAndlevelFilterHandler) Handle(ctx context.Context, record slog.Record) error {
	attrs := make([]slog.Attr, 0, 2)
	attrs = append(attrs, contextAttributes(ctx)...)

	spanContext := oteltrace.SpanContextFromContext(ctx)
	if spanContext.IsValid() {
		attrs = append(attrs,
			slog.String("trace_id", spanContext.TraceID().String()),
			slog.String("span_id", spanContext.SpanID().String()),
		)
	}

	if len(attrs) > 0 {
		record = record.Clone()
		record.AddAttrs(attrs...)
	}

	span := oteltrace.SpanFromContext(ctx)
	if h.addSpanEvent && span.IsRecording() {
		eventAttrs := make([]attribute.KeyValue, 0, record.NumAttrs())
		eventAttrs = append(eventAttrs, attribute.String("body", record.Message))
		eventAttrs = append(eventAttrs, attribute.String(slog.LevelKey, record.Level.String()))
		eventAttrs = append(eventAttrs, attribute.String(slog.TimeKey, record.Time.Format(time.RFC3339Nano)))
		record.Attrs(func(attr slog.Attr) bool {
			eventAttrs = append(eventAttrs, attribute.KeyValue{
				Key:   attribute.Key(attr.Key),
				Value: attribute.StringValue(fmt.Sprintf("%+v", attr.Value.Resolve().Any())),
			})

			return true
		})

		span.AddEvent("log.message", oteltrace.WithAttributes(eventAttrs...))
	}

	return h.next.Handle(ctx, record)
}

// WithGroup creates a new handler with the given group name.
func (h *tracingAndlevelFilterHandler) WithGroup(name string) slog.Handler {
	return &tracingAndlevelFilterHandler{level: h.level, next: h.next.WithGroup(name), addSpanEvent: h.addSpanEvent}
}

// WithAttrs creates a new handler with the given attributes.
func (h *tracingAndlevelFilterHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &tracingAndlevelFilterHandler{level: h.level, next: h.next.WithAttrs(attrs), addSpanEvent: h.addSpanEvent}
}

func noShutdown(context.Context) {}

// NewOtelHandler creates a new OpenTelemetry log handler.
func NewOtelHandler(ctx context.Context, otelService *observability.Service, options ...Option) (slog.Handler, func(context.Context), error) {
	opts := &Options{}
	for _, option := range options {
		option(opts)
	}

	processor, err := NewOtelLogProcessor(ctx, options...)
	if err != nil {
		return nil, noShutdown, err
	}

	loggerProvider := NewLoggingProviderWithProcessors(otelService, processor)

	return &tracingAndlevelFilterHandler{
			level:        opts.Level,
			next:         otelslog.NewHandler(otelLoggerName, otelslog.WithLoggerProvider(loggerProvider)),
			addSpanEvent: true,
		},
		func(ctx context.Context) {
			_ = loggerProvider.Shutdown(ctx)
		}, nil
}

// NewAxiomHandlerFromEnv creates a new Axiom log handler. It reads the following environment variables:
//
//   - AXIOM_API_KEY: The Axiom API key.
//   - AXIOM_LOGS_DATASET: The dataset where logs should be stored.
func NewAxiomHandlerFromEnv(ctx context.Context, otelService *observability.Service, options ...Option) (slog.Handler, func(context.Context), error) {
	apiKey := os.Getenv("AXIOM_API_KEY")
	if apiKey == "" {
		return nil, noShutdown, errors.New("AXIOM_API_KEY environment variable is not set")
	}

	dataset := os.Getenv("AXIOM_LOGS_DATASET")
	if dataset == "" {
		return nil, noShutdown, errors.New("AXIOM_LOGS_DATASET environment variable is not set")
	}

	return NewAxiomHandler(ctx, otelService, dataset, apiKey, options...)
}

// NewAxiomHandler creates a new Axiom log handler, emitting logs to the given dataset and authenticating with the given API key.
func NewAxiomHandler(ctx context.Context, otelService *observability.Service, dataset string, apiKey string, options ...Option) (slog.Handler, func(context.Context), error) {
	headers := map[string]string{
		"Authorization":   fmt.Sprintf("Bearer %s", apiKey),
		"X-Axiom-Dataset": dataset,
	}

	opts := make([]Option, 0, len(options)+2)
	opts = append(opts, WithEndpointURL(axiomLogsEndpoint), WithHeaders(headers))
	opts = append(opts, options...)

	return NewOtelHandler(ctx, otelService, opts...)
}

// NewConsoleHandler creates a new console log handler.
func NewConsoleHandler(options ...Option) slog.Handler {
	opts := &Options{}
	for _, option := range options {
		option(opts)
	}

	return &tracingAndlevelFilterHandler{
		level: opts.Level,
		next:  tint.NewHandler(os.Stdout, &tint.Options{Level: opts.Level}),
	}
}

// New creates a new slog.Logger with the given handlers.
func New(handlers ...slog.Handler) *slog.Logger {
	return slog.New(slogmulti.Fanout(handlers...))
}
