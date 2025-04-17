package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/tilebox/tilebox-go/examples/workflows/axiom"
	"github.com/tilebox/tilebox-go/observability"
	"github.com/tilebox/tilebox-go/workflows/v1"
)

var (
	serviceName    = "task-runner"
	serviceVersion = "dev"
)

func main() {
	ctx := context.Background()

	tileboxAPIKey := os.Getenv("TILEBOX_API_KEY")
	endpoint := "https://api.axiom.co/v1/logs"
	headers := map[string]string{
		"Authorization": "Bearer <apikey>",
	}

	// Setup OpenTelemetry logging and slog
	otelService := &observability.Service{Name: serviceName, Version: serviceVersion}
	axiomHandler, shutdownLogger, err := observability.NewOtelHandler(ctx, otelService, observability.WithEndpointURL(endpoint), observability.WithHeaders(headers))
	if err != nil {
		slog.Error("failed to set up axiom log handler", slog.Any("error", err))
		return
	}
	observability.InitializeLogging(axiomHandler)
	defer shutdownLogger(ctx)

	// Setup OpenTelemetry tracing
	axiomProcessor, err := observability.NewOtelSpanProcessor(ctx, observability.WithEndpointURL(endpoint), observability.WithHeaders(headers))
	if err != nil {
		slog.Error("failed to set up axiom span processor", slog.Any("error", err))
		return
	}
	shutdownTracer := observability.InitializeTracing(otelService, axiomProcessor)
	defer shutdownTracer(ctx)

	///////////
	/*axiomDataset := "my-dataset"
	apiKey := "key"
	endpoint := "https://api.axiom.co/v1/traces"
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", apiKey),
	}

	otelService := &observability.Service{Name: serviceName, Version: serviceVersion}

	processor, err := observability.NewOtelSpanProcessor(ctx, endpoint, observability.WithHeaders(headers))
	axiomProcessor, err := observability.NewAxiomSpanProcessor(ctx, axiomDataset, apiKey)

	shutdown := observability.InitializeTracing(otelService, processor, axiomProcessor)
	defer shutdown(ctx)

	tracerProvider := observability.NewTracerProvider(otelService, processor, axiomProcessor)
	otel.SetTracerProvider(tracerProvider)
	defer tracerProvider.Shutdown(ctx)

	////////

	otelHandler, shutdownLogger, err := observability.NewOtelHandler(ctx, otelService, endpoint, observability.WithHeaders(headers), observability.WithLevel(slog.LevelWarn))
	axiomHandler, shutdownAxiomLogger, err := observability.NewAxiomHandler(ctx, otelService, axiomDataset, apiKey, observability.WithLevel(slog.LevelWarn))
	defer shutdownAxiomLogger(ctx)
	consoleHandler := observability.NewConsoleHandler(observability.WithLevel(slog.LevelInfo))
	observability.InitializeLogging(otelHandler, axiomHandler, consoleHandler)

	otelLogProcessor, err := observability.NewOtelLogProcessor(ctx, endpoint)
	loggerProvider := observability.NewLoggingProvider(otelService, otelLogProcessor)
	global.SetLoggerProvider(loggerProvider)
	defer loggerProvider.Shutdown(ctx)

	observability.InitializeLogging(observability.NewConsoleHandler())*/

	///////////

	client := workflows.NewClient(workflows.WithAPIKey(tileboxAPIKey))

	cluster, err := client.Clusters.Get(ctx, "testing-4qgCk4qHH85qR7")
	if err != nil {
		slog.Error("failed to get cluster", slog.Any("error", err))
		return
	}

	job, err := client.Jobs.Submit(ctx, "Axiom Observability", cluster,
		[]workflows.Task{&axiom.MyAxiomTask{}},
	)
	if err != nil {
		slog.Error("Failed to submit job", "error", err)
		return
	}

	slog.Info("Job submitted", slog.String("job_id", job.ID.String()))
}
