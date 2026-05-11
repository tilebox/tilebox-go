package workflows

import (
	"context"
	"log/slog"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	oteltrace "go.opentelemetry.io/otel/trace"
)

var initialDefaultSlogLogger = slog.Default()

type recordingSlogHandler struct {
	name    string
	records *[]string
}

func (h *recordingSlogHandler) Enabled(context.Context, slog.Level) bool {
	return true
}

func (h *recordingSlogHandler) Handle(_ context.Context, record slog.Record) error {
	*h.records = append(*h.records, h.name+":"+record.Message)
	return nil
}

func (h *recordingSlogHandler) WithAttrs([]slog.Attr) slog.Handler {
	return h
}

func (h *recordingSlogHandler) WithGroup(string) slog.Handler {
	return h
}

func TestConfigureSlogHandlerFansOutToConfiguredHandlers(t *testing.T) {
	resetSlogConfigurationForTest(t)

	records := make([]string, 0)
	configureSlogHandler(&recordingSlogHandler{name: "first", records: &records})
	slog.Info("one")

	configureSlogHandler(&recordingSlogHandler{name: "second", records: &records})
	slog.Info("two")

	assert.Equal(t, []string{"first:one", "first:two", "second:two"}, records)
}

func TestConfigureSlogHandlerPreservesExistingCustomHandler(t *testing.T) {
	resetSlogConfigurationForTest(t)

	records := make([]string, 0)
	slog.SetDefault(slog.New(&recordingSlogHandler{name: "existing", records: &records}))

	configureSlogHandler(&recordingSlogHandler{name: "configured", records: &records})
	slog.Info("hello")

	assert.Equal(t, []string{"existing:hello", "configured:hello"}, records)
}

func TestConfigureConsoleLoggingIsIdempotent(t *testing.T) {
	resetSlogConfigurationForTest(t)

	ConfigureConsoleLogging(slog.LevelInfo)
	ConfigureConsoleLogging(slog.LevelDebug)

	slogHandlersMu.Lock()
	defer slogHandlersMu.Unlock()
	assert.Len(t, slogHandlers, 1)
}

func TestWithSpanResultUsesTaskRunnerTracer(t *testing.T) {
	spanRecorder := tracetest.NewSpanRecorder()
	tracerProvider := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(spanRecorder))
	ctx := context.WithValue(context.Background(), contextKeyTaskExecution, &taskExecutionContext{
		runner: &TaskRunner{tracer: tracerProvider.Tracer("test")},
	})

	result, err := WithSpanResult(ctx, "expensive-compute", func(ctx context.Context) (string, error) {
		assert.True(t, oteltrace.SpanContextFromContext(ctx).IsValid())
		return "ok", nil
	})

	require.NoError(t, err)
	assert.Equal(t, "ok", result)
	endedSpans := spanRecorder.Ended()
	require.Len(t, endedSpans, 1)
	assert.Equal(t, "expensive-compute", endedSpans[0].Name())
}

func TestWithSpanRunsWithoutTracerOutsideTaskContext(t *testing.T) {
	called := false
	err := WithSpan(context.Background(), "without-task-context", func(ctx context.Context) error {
		called = true
		assert.False(t, oteltrace.SpanContextFromContext(ctx).IsValid())
		return nil
	})

	require.NoError(t, err)
	assert.True(t, called)
}

func resetSlogConfigurationForTest(t *testing.T) {
	t.Helper()

	defaultLogger := slog.Default()
	slog.SetDefault(initialDefaultSlogLogger)
	slogHandlersMu.Lock()
	previousHandlers := slogHandlers
	slogHandlers = nil
	consoleLogHandlerOnce = sync.Once{}
	slogHandlersMu.Unlock()

	t.Cleanup(func() {
		slog.SetDefault(defaultLogger)

		slogHandlersMu.Lock()
		defer slogHandlersMu.Unlock()
		slogHandlers = previousHandlers
		consoleLogHandlerOnce = sync.Once{}
	})
}
