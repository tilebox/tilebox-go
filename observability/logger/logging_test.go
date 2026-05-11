package logger

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

type recordingHandler struct {
	records []slog.Record
}

func (h *recordingHandler) Enabled(context.Context, slog.Level) bool {
	return true
}

func (h *recordingHandler) Handle(_ context.Context, record slog.Record) error {
	h.records = append(h.records, record.Clone())
	return nil
}

func (h *recordingHandler) WithAttrs([]slog.Attr) slog.Handler {
	return h
}

func (h *recordingHandler) WithGroup(string) slog.Handler {
	return h
}

func TestTracingAndLevelFilterHandlerAddsTraceAndContextAttrsAndSpanEvent(t *testing.T) {
	spanRecorder := tracetest.NewSpanRecorder()
	tracerProvider := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(spanRecorder))
	tracer := tracerProvider.Tracer("test")

	ctx, span := tracer.Start(context.Background(), "test-span")
	ctx = ContextWithSlogAttributes(ctx, slog.String("task_id", "019e070c-63b8-ce26-7a15-cec85a85de2b"))

	sink := &recordingHandler{}
	logHandler := &tracingAndlevelFilterHandler{level: slog.LevelInfo, next: sink, addSpanEvent: true}
	slog.New(logHandler).InfoContext(ctx, "hello", slog.String("custom", "value"))
	span.End()

	require.Len(t, sink.records, 1)
	logAttrs := slogRecordAttrs(sink.records[0])
	spanContext := span.SpanContext()
	assert.Equal(t, spanContext.TraceID().String(), logAttrs["trace_id"])
	assert.Equal(t, spanContext.SpanID().String(), logAttrs["span_id"])
	assert.Equal(t, "019e070c-63b8-ce26-7a15-cec85a85de2b", logAttrs["task_id"])
	assert.Equal(t, "value", logAttrs["custom"])

	endedSpans := spanRecorder.Ended()
	require.Len(t, endedSpans, 1)
	events := endedSpans[0].Events()
	require.Len(t, events, 1)
	assert.Equal(t, "log.message", events[0].Name)
	eventAttrs := otelAttrs(events[0].Attributes)
	assert.Equal(t, "hello", eventAttrs["body"])
	assert.Equal(t, "INFO", eventAttrs[slog.LevelKey])
	assert.Equal(t, spanContext.TraceID().String(), eventAttrs["trace_id"])
	assert.Equal(t, spanContext.SpanID().String(), eventAttrs["span_id"])
	assert.Equal(t, "019e070c-63b8-ce26-7a15-cec85a85de2b", eventAttrs["task_id"])
	assert.Equal(t, "value", eventAttrs["custom"])
}

func slogRecordAttrs(record slog.Record) map[string]string {
	attrs := make(map[string]string, record.NumAttrs())
	record.Attrs(func(attr slog.Attr) bool {
		attrs[attr.Key] = attr.Value.Resolve().String()
		return true
	})
	return attrs
}

func otelAttrs(kvs []attribute.KeyValue) map[string]string {
	attrs := make(map[string]string, len(kvs))
	for _, kv := range kvs {
		attrs[string(kv.Key)] = kv.Value.AsString()
	}
	return attrs
}
