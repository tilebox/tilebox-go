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
	serviceName = "task-runner"
	version     = "dev"
)

func main() {
	ctx := context.Background()

	tileboxAPIKey := os.Getenv("TILEBOX_API_KEY")
	endpoint := "" // FIXME: defaults to localhost:4318
	headers := map[string]string{
		"Authorization": "Bearer <apikey>", // FIXME: if required
	}

	// Setup OpenTelemetry logging and slog
	otelService := &observability.Service{Name: serviceName, Version: version}
	otelHandler, shutdownLogger, err := observability.NewOtelHandler(ctx, otelService, observability.WithEndpointURL(endpoint), observability.WithHeaders(headers))
	if err != nil {
		slog.Error("failed to set up otel log handler", slog.Any("error", err))
		return
	}
	observability.InitializeLogging(otelHandler, observability.NewConsoleHandler())
	defer shutdownLogger(ctx)

	// Setup OpenTelemetry tracing
	otelProcessor, err := observability.NewOtelSpanProcessor(ctx, observability.WithEndpointURL(endpoint), observability.WithHeaders(headers))
	if err != nil {
		slog.Error("failed to set up otel span processor", slog.Any("error", err))
		return
	}
	shutdownTracer := observability.InitializeTracing(otelService, otelProcessor)
	defer shutdownTracer(ctx)

	client := workflows.NewClient(workflows.WithAPIKey(tileboxAPIKey))

	cluster, err := client.Clusters.Get(ctx, "testing-4qgCk4qHH85qR7")
	if err != nil {
		slog.Error("failed to get cluster", slog.Any("error", err))
		return
	}

	taskRunner, err := client.NewTaskRunner(cluster)
	if err != nil {
		slog.Error("failed to create task runner", slog.Any("error", err))
		return
	}

	err = taskRunner.RegisterTasks(&axiom.MyAxiomTask{})
	if err != nil {
		slog.Error("failed to register tasks", slog.Any("error", err))
		return
	}

	taskRunner.RunForever(ctx)
}
