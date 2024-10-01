package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/tilebox/tilebox-go/datasets/v1"
)

func main() {
	ctx := context.Background()

	datasetClient := datasets.NewDatasetsService(
		datasets.NewDatasetClient(
			datasets.WithAPIKey(os.Getenv("TILEBOX_API_KEY")),
		),
	)

	collectionID := uuid.MustParse("90bfcc54-1a39-4b00-9692-45f302c53ec5") // Sentinel-2 S2A_S2MSI1C

	firstAugust2024 := time.Date(2024, time.August, 1, 0, 0, 0, 0, time.UTC)
	firstSept2024 := time.Date(2024, time.September, 1, 0, 0, 0, 0, time.UTC)
	loadInterval := datasets.NewTimeInterval(firstAugust2024, firstSept2024)

	datapoints, err := datasets.Collect(datasetClient.Load(ctx, collectionID, loadInterval, true, false))
	if err != nil {
		slog.ErrorContext(ctx, "Failed to load data", "error", err)
		return
	}

	slog.InfoContext(ctx, "Datapoints loaded", slog.Int("count", len(datapoints)))
	slog.InfoContext(ctx, "First datapoint",
		slog.String("id", datapoints[0].Meta.ID.String()),
		slog.Time("eventTime", datapoints[0].Meta.EventTime),
		slog.Time("ingestionTime", datapoints[0].Meta.IngestionTime),
	)
}
