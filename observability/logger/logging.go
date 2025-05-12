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
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.30.0"
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
	rs := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(otelService.Name),
		semconv.ServiceVersionKey.String(otelService.Version),
	)

	opts := []log.LoggerProviderOption{
		log.WithResource(rs),
	}

	for _, processor := range processors {
		opts = append(opts, log.WithProcessor(processor))
	}

	return log.NewLoggerProvider(opts...)
}

// tracingAndlevelFilterHandler is a slog.Handler that filters log messages based on the configured level
// and adds tracing information to each log record.
type tracingAndlevelFilterHandler struct {
	level slog.Level
	next  slog.Handler
}

var _ slog.Handler = &tracingAndlevelFilterHandler{}

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
	// add log info to span event
	span := oteltrace.SpanFromContext(ctx)
	if span != nil && span.IsRecording() {
		eventAttrs := make([]attribute.KeyValue, 0, record.NumAttrs())
		eventAttrs = append(eventAttrs, attribute.String(slog.MessageKey, record.Message))
		eventAttrs = append(eventAttrs, attribute.String(slog.LevelKey, record.Level.String()))
		eventAttrs = append(eventAttrs, attribute.String(slog.TimeKey, record.Time.Format(time.RFC3339Nano)))
		record.Attrs(func(attr slog.Attr) bool {
			eventAttrs = append(eventAttrs, attribute.KeyValue{
				Key:   attribute.Key(attr.Key),
				Value: attribute.StringValue(fmt.Sprintf("%+v", attr.Value.Resolve().Any())),
			})

			return true
		})

		span.AddEvent("log_record", oteltrace.WithAttributes(eventAttrs...))
	}
	return h.next.Handle(ctx, record)
}

// WithGroup creates a new handler with the given group name.
func (h *tracingAndlevelFilterHandler) WithGroup(name string) slog.Handler {
	return &tracingAndlevelFilterHandler{level: h.level, next: h.next.WithGroup(name)}
}

// WithAttrs creates a new handler with the given attributes.
func (h *tracingAndlevelFilterHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &tracingAndlevelFilterHandler{level: h.level, next: h.next.WithAttrs(attrs)}
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
			level: opts.Level,
			next:  otelslog.NewHandler(otelLoggerName, otelslog.WithLoggerProvider(loggerProvider)),
		},
		func(ctx context.Context) {
			_ = loggerProvider.Shutdown(ctx)
		}, nil
}

// NewAxiomHandler creates a new Axiom log handler. It reads the following environment variables:
//
//   - AXIOM_API_KEY: The Axiom API key.
//   - AXIOM_LOGS_DATASET: The dataset where logs should be stored.
func NewAxiomHandler(ctx context.Context, otelService *observability.Service, options ...Option) (slog.Handler, func(context.Context), error) {
	apiKey := os.Getenv("AXIOM_API_KEY")
	if apiKey == "" {
		return nil, noShutdown, errors.New("AXIOM_API_KEY environment variable is not set")
	}

	dataset := os.Getenv("AXIOM_LOGS_DATASET")
	if dataset == "" {
		return nil, noShutdown, errors.New("AXIOM_LOGS_DATASET environment variable is not set")
	}

	headers := map[string]string{
		"Authorization":   fmt.Sprintf("Bearer %s", apiKey),
		"X-Axiom-Dataset": dataset,
	}

	opts := []Option{WithEndpointURL(axiomLogsEndpoint), WithHeaders(headers)}
	opts = append(opts, options...)

	return NewOtelHandler(ctx, otelService, opts...)
}

// NewConsoleHandler creates a new console log handler.
func NewConsoleHandler(options ...Option) slog.Handler {
	opts := &Options{}
	for _, option := range options {
		option(opts)
	}

	return tint.NewHandler(os.Stdout, &tint.Options{Level: opts.Level})
}

// New creates a new slog.Logger with the given handlers.
func New(handlers ...slog.Handler) *slog.Logger {
	return slog.New(slogmulti.Fanout(handlers...))
}
