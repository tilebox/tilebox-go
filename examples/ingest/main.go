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

	dataCopernicus := []*tileboxdatasets.Datapoint[*datasetsv1.CopernicusDataspaceGranule]{
		{
			Meta: &tileboxdatasets.DatapointMetadata{
				EventTime: time.Now(),
			},
			Data: &datasetsv1.CopernicusDataspaceGranule{
				GranuleName: "Granule 1",
			},
		},
	}

	// Validate the datapoints and convert it to the right type
	validatedData, err := tileboxdatasets.ValidateAs(dataCopernicus)
	if err != nil {
		log.Fatalf("Failed to validate datapoints as copernicus granules: %v", err)
	}

	// Ingest datapoints
	ingestResponse, err := collection.Ingest(ctx, validatedData, false)
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
