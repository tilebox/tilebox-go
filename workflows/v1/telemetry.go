package workflows

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"path"
	"reflect"
	"strings"
	"sync"

	"github.com/tilebox/tilebox-go/observability"
	"github.com/tilebox/tilebox-go/observability/logger"
	obstracer "github.com/tilebox/tilebox-go/observability/tracer"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

const (
	otlpLogsPath   = "/v1/logs"
	otlpTracesPath = "/v1/traces"
)

var (
	tileboxTelemetryService = &observability.Service{Name: "tilebox-go"}
	tileboxLogHandlerOnce   sync.Once
	consoleLogHandlerOnce   sync.Once
	tileboxTraceProviders   sync.Map
	slogHandlersMu          sync.Mutex
	slogHandlers            []slog.Handler
	slogBaseLogger          *slog.Logger
)

// ConfigureConsoleLogging configures the default slog logger to write logs to stdout at the given level.
//
// Console logging is composed with the Tilebox log exporter configured by NewClient. Calling
// ConfigureConsoleLogging before or after NewClient causes future slog records to be sent to both
// destinations. If NewClient is not called, logs are only written to the console.
func ConfigureConsoleLogging(level slog.Level) {
	consoleLogHandlerOnce.Do(func() {
		configureSlogHandler(logger.NewConsoleHandler(logger.WithLevel(level)))
	})
}

// WithSpanResult wraps f with a tracing span using the current task runner's tracer.
// If ctx is not a task execution context or the task runner has no tracer, f is called without creating a span.
func WithSpanResult[Result any](ctx context.Context, name string, f func(ctx context.Context) (Result, error)) (Result, error) {
	executionContext := getTaskExecutionContext(ctx)
	if executionContext == nil || executionContext.executor.tracer == nil {
		return f(ctx)
	}

	return observability.WithSpanResult(ctx, executionContext.executor.tracer, name, f)
}

// WithSpan wraps f with a tracing span using the current task runner's tracer.
// If ctx is not a task execution context or the task runner has no tracer, f is called without creating a span.
func WithSpan(ctx context.Context, name string, f func(ctx context.Context) error) error {
	executionContext := getTaskExecutionContext(ctx)
	if executionContext == nil || executionContext.executor.tracer == nil {
		return f(ctx)
	}

	return observability.WithSpan(ctx, executionContext.executor.tracer, name, f)
}

func configureTileboxTelemetry(ctx context.Context, cfg *clientConfig) {
	if cfg.apiKey == "" || !isHTTPURL(cfg.url) {
		return
	}

	if !cfg.disableTracing {
		cfg.tracerProvider = configureTileboxTracing(ctx, cfg)
	}

	if !cfg.disableLogging {
		configureTileboxLogging(ctx, cfg)
	}
}

func configureTileboxTracing(ctx context.Context, cfg *clientConfig) trace.TracerProvider {
	endpoint, ok := otlpEndpointURL(cfg.url, otlpTracesPath)
	if !ok {
		return cfg.tracerProvider
	}

	globalTracerProvider := otel.GetTracerProvider()
	if tracerProvider, ok := globalTracerProvider.(*sdktrace.TracerProvider); ok {
		if _, loaded := tileboxTraceProviders.Load(tracerProvider); loaded {
			return tracerProvider
		}
	}

	spanProcessor, err := obstracer.NewOtelSpanProcessor(ctx,
		obstracer.WithEndpointURL(endpoint),
		obstracer.WithHeaders(tileboxTelemetryHeaders(cfg.apiKey)),
	)
	if err != nil {
		otel.Handle(fmt.Errorf("failed to configure Tilebox trace exporter: %w", err))
		return cfg.tracerProvider
	}

	if tracerProvider, ok := globalTracerProvider.(*sdktrace.TracerProvider); ok {
		tracerProvider.RegisterSpanProcessor(spanProcessor)
		tileboxTraceProviders.Store(tracerProvider, endpoint)
		return tracerProvider
	}

	if isDefaultTracerProvider(globalTracerProvider) {
		tracerProvider := obstracer.NewOtelProviderWithProcessors(tileboxTelemetryService, spanProcessor)
		otel.SetTracerProvider(tracerProvider)
		tileboxTraceProviders.Store(tracerProvider, endpoint)
		return tracerProvider
	}

	return cfg.tracerProvider
}

func configureTileboxLogging(ctx context.Context, cfg *clientConfig) {
	tileboxLogHandlerOnce.Do(func() {
		otelHandler, _, err := NewTileboxLogHandler(ctx, cfg.url, cfg.apiKey, slog.LevelInfo)
		if err != nil {
			otel.Handle(fmt.Errorf("failed to configure Tilebox log exporter: %w", err))
			return
		}

		configureSlogHandler(otelHandler)
	})
}

// NewTileboxLogHandler creates a slog handler that exports log records to the Tilebox API.
func NewTileboxLogHandler(ctx context.Context, apiURL string, apiKey string, level slog.Level) (slog.Handler, func(context.Context), error) {
	cleanup := func(context.Context) {}
	if strings.TrimSpace(apiKey) == "" {
		return nil, cleanup, errors.New("tilebox API key is required for log export")
	}
	if !isHTTPURL(apiURL) {
		return nil, cleanup, fmt.Errorf("tilebox API URL must be an HTTP(S) URL for log export: %s", apiURL)
	}

	endpoint, ok := otlpEndpointURL(apiURL, otlpLogsPath)
	if !ok {
		return nil, cleanup, fmt.Errorf("invalid Tilebox API URL for log export: %s", apiURL)
	}

	handler, cleanup, err := logger.NewOtelHandler(ctx, tileboxTelemetryService,
		logger.WithEndpointURL(endpoint),
		logger.WithHeaders(tileboxTelemetryHeaders(apiKey)),
		logger.WithLevel(level),
	)
	if cleanup == nil {
		cleanup = func(context.Context) {}
	}
	return handler, cleanup, err
}

func configureSlogHandler(handler slog.Handler) {
	slogHandlersMu.Lock()
	defer slogHandlersMu.Unlock()
	configureSlogHandlerLocked(handler)
}

func configureSlogHandlerLocked(handler slog.Handler) {
	if len(slogHandlers) == 0 {
		slogBaseLogger = slog.Default()
		existingHandler := slogBaseLogger.Handler()
		if !isDefaultSlogHandler(existingHandler) {
			slogHandlers = append(slogHandlers, existingHandler)
		}
	}

	slogHandlers = append(slogHandlers, handler)
	slog.SetDefault(logger.New(slogHandlers...))
}

type removableSlogHandler struct {
	slog.Handler
}

func configureRemovableSlogHandler(handler slog.Handler) func() {
	registeredHandler := &removableSlogHandler{Handler: handler}
	slogHandlersMu.Lock()
	configureSlogHandlerLocked(registeredHandler)
	slogHandlersMu.Unlock()

	var removeOnce sync.Once
	return func() {
		removeOnce.Do(func() {
			slogHandlersMu.Lock()
			defer slogHandlersMu.Unlock()

			for i, configuredHandler := range slogHandlers {
				if configuredHandler == registeredHandler {
					slogHandlers = append(slogHandlers[:i], slogHandlers[i+1:]...)
					break
				}
			}

			if len(slogHandlers) > 0 {
				slog.SetDefault(logger.New(slogHandlers...))
				return
			}
			if slogBaseLogger != nil {
				slog.SetDefault(slogBaseLogger)
				slogBaseLogger = nil
			}
		})
	}
}

func isDefaultSlogHandler(handler slog.Handler) bool {
	handlerType := reflect.TypeOf(handler)
	if handlerType == nil || handlerType.Kind() != reflect.Pointer {
		return false
	}

	handlerType = handlerType.Elem()
	return handlerType.PkgPath() == "log/slog" && handlerType.Name() == "defaultHandler"
}

func isDefaultTracerProvider(tracerProvider trace.TracerProvider) bool {
	tracerProviderType := reflect.TypeOf(tracerProvider)
	if tracerProviderType == nil || tracerProviderType.Kind() != reflect.Pointer {
		return false
	}

	tracerProviderType = tracerProviderType.Elem()
	return tracerProviderType.PkgPath() == "go.opentelemetry.io/otel/internal/global" && tracerProviderType.Name() == "tracerProvider"
}

func tileboxTelemetryHeaders(apiKey string) map[string]string {
	return map[string]string{"Authorization": "Bearer " + apiKey}
}

func isHTTPURL(rawURL string) bool {
	return strings.HasPrefix(rawURL, "https://") || strings.HasPrefix(rawURL, "http://")
}

func otlpEndpointURL(baseURL, signalPath string) (string, bool) {
	parsedURL, err := url.Parse(baseURL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return "", false
	}

	parsedURL.RawQuery = ""
	parsedURL.Fragment = ""
	parsedURL.Path = path.Join(parsedURL.Path, signalPath)
	if !strings.HasPrefix(parsedURL.Path, "/") {
		parsedURL.Path = "/" + parsedURL.Path
	}

	return parsedURL.String(), true
}
