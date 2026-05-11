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
var service = &observability.Service{Name: "axiom-example", Version: "dev"}

type MyAxiomTask struct{}

func (t *MyAxiomTask) Execute(ctx context.Context) error {
	// Start an openTelemetry tracing span that will be exported to Axiom
	result, err := workflows.WithSpanResult(ctx, "compute", func(ctx context.Context) (int, error) {
		return 6 * 7, nil
	})
	if err != nil {
		return fmt.Errorf("failed to compute: %w", err)
	}

	// Log the result to Axiom
	slog.InfoContext(ctx, "Computed result", slog.Int("result", result))
	return nil
}

func main() {
	ctx := context.Background()

	// Setup OpenTelemetry logging and slog
	axiomHandler, shutdownLogger, err := logger.NewAxiomHandlerFromEnv(ctx, service,
		logger.WithLevel(slog.LevelInfo), // export logs at info level and above as OTEL logs
	)
	defer shutdownLogger(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "failed to set up axiom log handler", slog.Any("error", err))
		return
	}
	tileboxLogger := logger.New( // initialize a slog.Logger
		axiomHandler, // export logs to Axiom
	)
	slog.SetDefault(tileboxLogger) // all future slog calls will be forwarded to the tilebox logger
	workflows.ConfigureConsoleLogging(slog.LevelWarn)

	// Setup an OpenTelemetry trace span processor, exporting traces and spans to Axiom
	tileboxTracerProvider, shutdown, err := tracer.NewAxiomProviderFromEnv(ctx, service)
	defer shutdown(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "failed to set up axiom tracer provider", slog.Any("error", err))
		return
	}
	otel.SetTracerProvider(tileboxTracerProvider) // set the tilebox tracer provider as the global OTEL tracer provider

	client := workflows.NewClient()

	job, err := client.Jobs.Submit(ctx, "Axiom Observability",
		[]workflows.Task{&MyAxiomTask{}},
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

	if err := taskRunner.RegisterTasks(&MyAxiomTask{}); err != nil {
		slog.ErrorContext(ctx, "failed to register tasks", slog.Any("error", err))
		return
	}

	taskRunner.RunForever(ctx)
}
