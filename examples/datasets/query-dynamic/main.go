package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"time"

	"github.com/paulmach/orb"
	"github.com/tilebox/tilebox-go/datasets/v1"
	"github.com/tilebox/tilebox-go/query"
)

func main() {
	ctx := context.Background()

	// Create a Tilebox Datasets client.
	client := datasets.NewClient()

	// Select a dataset. The dataset metadata contains the protobuf descriptor needed to decode datapoints dynamically.
	dataset, err := client.Datasets.Get(ctx, "open_data.copernicus.sentinel2_msi")
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get dataset", slog.Any("error", err))
		return
	}

	// Build the descriptor once and reuse it for every datapoint returned by the query.
	descriptor, err := datasets.NewDatapointDescriptor(dataset)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to build datapoint descriptor", slog.Any("error", err))
		return
	}

	// Select a collection.
	collection, err := client.Collections.Get(ctx, dataset.ID, "S2A_S2MSI1C")
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get collection", slog.Any("error", err))
		return
	}

	// Select a temporal extent.
	startDate := time.Date(2025, 4, 2, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 4, 3, 0, 0, 0, 0, time.UTC)

	// Select a spatial extent.
	area := orb.Polygon{ // area roughly covering the state of Colorado
		{{-109.05, 41.00}, {-102.05, 41.00}, {-102.05, 37.0}, {-109.045, 37.0}, {-109.05, 41.00}},
	}

	// Query raw protobuf datapoint bytes and decode them dynamically into maps.
	datapoints := make([]map[string]any, 0)
	for rawDatapoint, err := range client.Datapoints.Query(ctx,
		dataset.ID,
		datasets.WithCollections(collection),
		datasets.WithTemporalExtent(query.NewTimeInterval(startDate, endDate)),
		datasets.WithSpatialExtent(area),
		datasets.WithLimit(5),
	) {
		if err != nil {
			slog.ErrorContext(ctx, "Failed to query datapoints", slog.Any("error", err))
			return
		}

		datapoint, err := datasets.UnmarshalDatapoint(descriptor, rawDatapoint)
		if err != nil {
			slog.ErrorContext(ctx, "Failed to decode datapoint", slog.Any("error", err))
			return
		}
		datapoints = append(datapoints, datapoint)
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(datapoints); err != nil {
		slog.ErrorContext(ctx, "Failed to write datapoints as JSON", slog.Any("error", err))
	}
}
