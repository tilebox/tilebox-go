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

	// Select a dataset
	dataset, err := client.Datasets.Get(ctx, "open_data.copernicus.sentinel1_sar")
	if err != nil {
		log.Fatalf("Failed to get dataset: %v", err)
	}

	// Create a collection
	collection, err := client.Collections.Create(ctx, dataset.ID, "My First Collection")
	if err != nil {
		log.Fatalf("Failed to create collection: %v", err)
	}

	// Create datapoints
	datapoints := []*tileboxdatasets.Datapoint{
		tileboxdatasets.NewDatapoint(time.Now(), &datasetsv1.CopernicusDataspaceGranule{GranuleName: "Granule 1"}),
		tileboxdatasets.NewDatapoint(time.Now().Add(-5*time.Hour), &datasetsv1.CopernicusDataspaceGranule{GranuleName: "Past Granule 2"}),
	}

	// Ingest datapoints
	ingestResponse, err := client.Datapoints.Ingest(ctx, collection.ID, datapoints, false)
	if err != nil {
		log.Fatalf("Failed to ingest datapoints: %v", err)
	}
	slog.Info("Ingested datapoints", slog.Int64("created", ingestResponse.NumCreated))

	readings := make([]*tileboxdatasets.Datapoint, 0, 10)
	for range 10 {
		datapoint := &datasetsv1.CopernicusDataspaceGranule{
			GranuleName: "granule",
			ProductType: "product",
		}

		readings = append(readings, tileboxdatasets.NewDatapoint(time.Now(), datapoint))
	}
	_, err = client.Datapoints.Ingest(ctx, collection.ID, readings, false)
	if err != nil {
		log.Fatalf("Failed to ingest datapoints: %v", err)
	}

	// Delete datapoints again
	deleteResponse, err := client.Datapoints.DeleteIDs(ctx, collection.ID, ingestResponse.DatapointIDs)
	if err != nil {
		log.Fatalf("Failed to delete datapoints: %v", err)
	}
	slog.Info("Deleted datapoints", slog.Int64("deleted", deleteResponse.NumDeleted))
}
