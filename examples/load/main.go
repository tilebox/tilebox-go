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

	// Select a collection
	collection, err := dataset.Collection(ctx, "S1A_EW_GRDM_1S-COG")
	if err != nil {
		log.Fatalf("Failed to get collection: %v", err)
	}

	// Select a time interval
	oneMonthAgo := time.Now().AddDate(0, -1, 0)
	loadInterval := tileboxdatasets.NewTimeInterval(oneMonthAgo, time.Now(), false, false)

	// Load one month of data
	datapoints, err := tileboxdatasets.CollectAs[*datasetsv1.CopernicusDataspaceGranule](collection.Load(ctx, loadInterval, false, false))
	if err != nil {
		log.Fatalf("Failed to load and collect datapoints: %v", err)
	}

	slog.Info("Datapoints loaded", slog.Int("count", len(datapoints)))
	slog.Info("First datapoint",
		slog.String("id", datapoints[0].Meta.ID.String()),
		slog.Time("event time", datapoints[0].Meta.EventTime),
		slog.Time("ingestion time", datapoints[0].Meta.IngestionTime),
		slog.String("granule name", datapoints[0].Data.GetGranuleName()),
		slog.String("processing level", datapoints[0].Data.GetProcessingLevel().String()),
		slog.String("product type", datapoints[0].Data.GetProductType()),
		// and so on...
	)
}
