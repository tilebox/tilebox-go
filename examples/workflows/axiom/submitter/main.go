package main

import (
	"context"
	"log/slog"

	"github.com/tilebox/tilebox-go/examples/workflows/axiom"
	"github.com/tilebox/tilebox-go/observability"
	"github.com/tilebox/tilebox-go/observability/logger"
	"github.com/tilebox/tilebox-go/observability/tracer"
	"github.com/tilebox/tilebox-go/workflows/v1"
	"go.opentelemetry.io/otel"
)

// specify a service name and version to identify the instrumenting application in traces and logs
var service = &observability.Service{Name: "submitter", Version: "dev"}

func main() {
	ctx := context.Background()

	// Setup OpenTelemetry logging and slog
	axiomHandler, shutdownLogger, err := logger.NewAxiomHandlerFromEnv(ctx, service,
		logger.WithLevel(slog.LevelInfo), // export logs at info level and above as OTEL logs
	)
	defer shutdownLogger(ctx)
	if err != nil {
		slog.Error("failed to set up axiom log handler", slog.Any("error", err))
		return
	}
	tileboxLogger := logger.New( // initialize a slog.Logger
		axiomHandler, // export logs to Axiom
		logger.NewConsoleHandler(logger.WithLevel(slog.LevelWarn)), // and additionally, export WARN and ERROR logs to stdout
	)
	slog.SetDefault(tileboxLogger) // all future slog calls will be forwarded to the tilebox logger

	// Setup an OpenTelemetry trace span processor, exporting traces and spans to Axiom
	tileboxTracerProvider, shutdown, err := tracer.NewAxiomProviderFromEnv(ctx, service)
	defer shutdown(ctx)
	if err != nil {
		slog.Error("failed to set up axiom tracer provider", slog.Any("error", err))
		return
	}
	otel.SetTracerProvider(tileboxTracerProvider) // set the tilebox tracer provider as the global OTEL tracer provider

	client := workflows.NewClient()
	job, err := client.Jobs.Submit(ctx, "Axiom Observability",
		[]workflows.Task{&axiom.MyAxiomTask{}},
	)
	if err != nil {
		slog.Error("Failed to submit job", "error", err)
		return
	}

	slog.Info("Job submitted", slog.String("job_id", job.ID.String()))
}
