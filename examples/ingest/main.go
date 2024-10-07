package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"time"

	tileboxdatasets "github.com/tilebox/tilebox-go/datasets/v1"
	datasetsv1 "github.com/tilebox/tilebox-go/protogen/go/datasets/v1"
)

func main() {
	ctx := context.Background()

	// Create a Tilebox Datasets client
	client := tileboxdatasets.NewClient(
		tileboxdatasets.WithAPIKey(os.Getenv("TILEBOX_API_KEY")),
	)

	// List all available datasets
	datasets, err := client.Datasets(ctx)
	if err != nil {
		log.Fatalf("Failed to list datasets: %v", err)
	}
	dataset := datasets[1]

	// Create a collection
	collection, err := dataset.CreateCollection(ctx, "My First Collection")
	if err != nil {
		log.Fatalf("Failed to create collection: %v", err)
	}

	// Create datapoints
	datapoints, err := tileboxdatasets.Datapoints(
		tileboxdatasets.NewDatapoint(time.Now(), &datasetsv1.CopernicusDataspaceGranule{GranuleName: "Granule 1"}),
		tileboxdatasets.NewDatapoint(time.Now().Add(-5*time.Hour), &datasetsv1.CopernicusDataspaceGranule{GranuleName: "Past Granule 2"}),
	)
	if err != nil {
		log.Fatalf("Failed to create datapoints: %v", err)
	}

	// Ingest datapoints
	ingestResponse, err := collection.Ingest(ctx, datapoints, false)
	if err != nil {
		log.Fatalf("Failed to ingest datapoints: %v", err)
	}
	slog.Info("Ingested datapoints", slog.Int64("created", ingestResponse.NumCreated))

	// Delete datapoints again
	deleteResponse, err := collection.DeleteIDs(ctx, ingestResponse.DatapointIDs)
	if err != nil {
		log.Fatalf("Failed to delete datapoints: %v", err)
	}
	slog.Info("Deleted datapoints", slog.Int64("deleted", deleteResponse.NumDeleted))
}
