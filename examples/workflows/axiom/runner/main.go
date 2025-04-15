package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/tilebox/tilebox-go/examples/workflows/axiom"
	"github.com/tilebox/tilebox-go/observability"
	"github.com/tilebox/tilebox-go/workflows/v1"
	"go.opentelemetry.io/otel"
)

var (
	serviceName = "task-runner"
	version     = "dev"
)

func main() {
	ctx := context.Background()

	tileboxAPIKey := os.Getenv("TILEBOX_API_KEY")
	axiomAPIKey := os.Getenv("AXIOM_API_KEY")
	axiomTracesDataset := os.Getenv("AXIOM_TRACES_DATASET")
	axiomLogsDataset := os.Getenv("AXIOM_LOGS_DATASET")

	// Setup slog
	axiomLogHandler, shutdownLogger, err := observability.NewAxiomLogger(axiomLogsDataset, axiomAPIKey, slog.LevelDebug, true)
	defer shutdownLogger()
	if err != nil {
		slog.Error("failed to set up axiom log handler", slog.Any("error", err))
		return
	}
	slog.SetDefault(slog.New(axiomLogHandler)) // set the global slog logger to our axiom logger

	// Setup OpenTelemetry tracer provider
	axiomTracerProvider, shutdownTracer, err := observability.NewAxiomTracerProvider(ctx, axiomTracesDataset, axiomAPIKey, serviceName, version)
	defer shutdownTracer()
	if err != nil {
		slog.Error("failed to set up axiom tracer provider", slog.Any("error", err))
		return
	}
	otel.SetTracerProvider(axiomTracerProvider) // set the global tracer provider to our axiom tracer provider

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
