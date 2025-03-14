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

	// Select a collection
	collection, err := client.Collections.Get(ctx, dataset.ID, "S1A_EW_GRDM_1S-COG")
	if err != nil {
		log.Fatalf("Failed to get collection: %v", err)
	}

	// Select a time interval
	oneMonthAgo := time.Now().AddDate(0, -1, 0)
	loadInterval := tileboxdatasets.NewStandardTimeInterval(oneMonthAgo, time.Now())

	// Load one month of data
	var datapoints []*tileboxdatasets.TypedDatapoint[*datasetsv1.CopernicusDataspaceGranule]
	err = client.Datapoints.LoadInto(ctx, collection.ID, loadInterval, &datapoints)
	if err != nil {
		log.Fatalf("Failed to load datapoints: %v", err)
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
