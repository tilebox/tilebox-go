package main

import (
	"context"
	"log"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/wkt"
	"github.com/tilebox/tilebox-go/datasets/v1"
	examplesv1 "github.com/tilebox/tilebox-go/protogen/examples/v1"
	"github.com/tilebox/tilebox-go/query"
)

func main() {
	ctx := context.Background()

	// Create a Tilebox Datasets client
	client := datasets.NewClient()

	// Select a dataset
	dataset, err := client.Datasets.Get(ctx, "open_data.copernicus.sentinel2_msi")
	if err != nil {
		log.Fatalf("Failed to get dataset: %v", err)
	}

	// Select a collection
	collection, err := client.Collections.Get(ctx, dataset.ID, "S2A_S2MSI1C")
	if err != nil {
		log.Fatalf("Failed to get collection: %v", err)
	}

	// Select a temporal extent
	startDate := time.Date(2025, 4, 2, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 4, 3, 0, 0, 0, 0, time.UTC)

	// Select a spatial extent
	area := orb.Polygon{ // area roughly covering the state of Colorado
		{{-109.05, 41.00}, {-102.05, 41.00}, {-102.05, 37.0}, {-109.045, 37.0}, {-109.05, 41.00}},
	}

	// Perform a spatial-temporal query
	// Sentinel2Msi type is generated using tilebox-generate
	var foundDatapoints []*examplesv1.Sentinel2Msi
	err = client.Datapoints.QueryInto(ctx,
		[]uuid.UUID{collection.ID},
		&foundDatapoints,
		datasets.WithTemporalExtent(query.NewTimeInterval(startDate, endDate)),
		datasets.WithSpatialExtent(area),
	)
	if err != nil {
		log.Fatalf("Failed to query datapoints: %v", err)
	}

	slog.Info("Found datapoints over Colorado in March 2025", slog.Int("count", len(foundDatapoints)))
	if len(foundDatapoints) > 0 {
		slog.Info("First datapoint over Colorado",
			slog.String("id", foundDatapoints[0].GetId().AsUUID().String()),
			slog.Time("event time", foundDatapoints[0].GetTime().AsTime()),
			slog.Time("ingestion time", foundDatapoints[0].GetIngestionTime().AsTime()),
			slog.String("geometry", wkt.MarshalString(foundDatapoints[0].GetGeometry().AsGeometry())),
			slog.String("granule name", foundDatapoints[0].GetGranuleName()),
			slog.String("processing level", foundDatapoints[0].GetProcessingLevel().String()),
			slog.String("product type", foundDatapoints[0].GetProductType()),
			// and so on...
		)
	}
}
