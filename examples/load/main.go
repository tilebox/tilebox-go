package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"time"

	tileboxdatasets "github.com/tilebox/tilebox-go/datasets/v1"
	testv1 "github.com/tilebox/tilebox-go/protogen-test/tilebox/v1"
)

func main() {
	ctx := context.Background()

	// Create a Tilebox Datasets client
	client := tileboxdatasets.NewClient(
		tileboxdatasets.WithAPIKey(os.Getenv("TILEBOX_API_KEY")),
	)

	// Select a dataset
	dataset, err := client.Datasets.Get(ctx, "tilebox.modis")
	if err != nil {
		log.Fatalf("Failed to get dataset: %v", err)
	}

	// Select a collection
	collection, err := client.Collections.Get(ctx, dataset.ID, "MCD12Q1")
	if err != nil {
		log.Fatalf("Failed to get collection: %v", err)
	}

	// Select a time interval
	jan2001 := time.Date(2001, time.January, 1, 0, 0, 0, 0, time.UTC)
	loadInterval := tileboxdatasets.NewStandardTimeInterval(jan2001, time.Now())

	// Load one month of data
	var datapoints []*testv1.Modis
	err = client.Datapoints.LoadInto(ctx, collection.ID, loadInterval, &datapoints)
	if err != nil {
		log.Fatalf("Failed to load datapoints: %v", err)
	}

	slog.Info("Datapoints loaded", slog.Int("count", len(datapoints)))
	if len(datapoints) == 0 {
		return
	}

	slog.Info("First datapoint",
		slog.String("id", datapoints[0].GetId().AsUUID().String()),
		slog.Time("event time", datapoints[0].GetTime().AsTime()),
		slog.Time("ingestion time", datapoints[0].GetIngestionTime().AsTime()),
		slog.String("granule name", datapoints[0].GetGranuleName()),
		// slog.String("processing level", datapoints[0].GetProcessingLevel().String()),
		// slog.String("product type", datapoints[0].GetProductType()),
		// and so on...
	)
}
