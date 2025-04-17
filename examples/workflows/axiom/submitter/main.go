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
	serviceName = "submitter"
	version     = "dev"
)

func main() {
	ctx := context.Background()

	tileboxAPIKey := os.Getenv("TILEBOX_API_KEY")
	axiomAPIKey := os.Getenv("AXIOM_API_KEY")
	axiomTracesDataset := os.Getenv("AXIOM_TRACES_DATASET")
	axiomLogsDataset := os.Getenv("AXIOM_LOGS_DATASET")

	// Setup OpenTelemetry logging and slog
	otelService := &observability.Service{Name: serviceName, Version: version}
	axiomHandler, shutdownLogger, err := observability.NewAxiomHandler(ctx, otelService, axiomLogsDataset, axiomAPIKey)
	if err != nil {
		slog.Error("failed to set up axiom log handler", slog.Any("error", err))
		return
	}
	observability.InitializeLogging(axiomHandler, observability.NewConsoleHandler())
	defer shutdownLogger(ctx)

	// Setup OpenTelemetry tracing
	axiomProcessor, err := observability.NewAxiomSpanProcessor(ctx, axiomTracesDataset, axiomAPIKey)
	if err != nil {
		slog.Error("failed to set up axiom span processor", slog.Any("error", err))
		return
	}
	shutdownTracer := observability.InitializeTracing(otelService, axiomProcessor)
	defer shutdownTracer(ctx)

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
