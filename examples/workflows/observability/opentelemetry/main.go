package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/tilebox/tilebox-go/observability"
	"github.com/tilebox/tilebox-go/observability/logger"
	"github.com/tilebox/tilebox-go/observability/tracer"
	"github.com/tilebox/tilebox-go/workflows/v1"
	"go.opentelemetry.io/otel"
)

// specify a service name and version to identify the instrumenting application in traces and logs
var service = &observability.Service{Name: "opentelemetry-example", Version: "dev"}

type MyOpenTelemetryTask struct{}

func (t *MyOpenTelemetryTask) Execute(ctx context.Context) error {
	// Start an openTelemetry tracing span that will be exported to the OpenTelemetry endpoint
	result, err := workflows.WithSpanResult(ctx, "compute", func(ctx context.Context) (int, error) {
		return 6 * 7, nil
	})
	if err != nil {
		return fmt.Errorf("failed to compute: %w", err)
	}

	// Log the result to the OpenTelemetry endpoint
	slog.InfoContext(ctx, "Computed result", slog.Int("result", result))
	return nil
}

func main() {
	ctx := context.Background()

	endpoint := "http://localhost:4318"
	headers := map[string]string{
		"Authorization": "Bearer <ENDPOINT_AUTH>",
	}

	// Setup an OpenTelemetry log handler, exporting logs to an OTEL compatible log endpoint
	otelHandler, shutdownLogger, err := logger.NewOtelHandler(ctx, service,
		logger.WithEndpointURL(endpoint),
		logger.WithHeaders(headers),
		logger.WithLevel(slog.LevelInfo), // export logs at info level and above as OTEL logs
	)
	defer shutdownLogger(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "failed to set up otel log handler", slog.Any("error", err))
		return
	}
	tileboxLogger := logger.New( // initialize a slog.Logger
		otelHandler, // export logs to the OTEL handler
	)
	slog.SetDefault(tileboxLogger) // all future slog calls will be forwarded to the tilebox logger
	workflows.ConfigureConsoleLogging(slog.LevelWarn)

	// Setup an OpenTelemetry trace span processor, exporting traces and spans to an OTEL compatible trace endpoint
	tileboxTracerProvider, shutdown, err := tracer.NewOtelProvider(ctx, service,
		tracer.WithEndpointURL(endpoint),
		tracer.WithHeaders(headers),
	)
	defer shutdown(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "failed to set up otel span processor", slog.Any("error", err))
		return
	}
	otel.SetTracerProvider(tileboxTracerProvider) // set the tilebox tracer provider as the global OTEL tracer provider

	client := workflows.NewClient()

	job, err := client.Jobs.Submit(ctx, "OpenTelemetry Observability",
		[]workflows.Task{&MyOpenTelemetryTask{}},
	)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to submit job", "error", err)
		return
	}

	slog.InfoContext(ctx, "Job submitted", slog.String("job_id", job.ID.String()))

	taskRunner, err := client.NewTaskRunner(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "failed to create task runner", slog.Any("error", err))
		return
	}

	if err := taskRunner.RegisterTasks(&MyOpenTelemetryTask{}); err != nil {
		slog.ErrorContext(ctx, "failed to register tasks", slog.Any("error", err))
		return
	}

	taskRunner.RunForever(ctx)
}
