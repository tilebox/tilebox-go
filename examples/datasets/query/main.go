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
	examplesv1 "github.com/tilebox/tilebox-go/protogen/go/examples/v1"
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
	startDate := time.Date(2025, time.March, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, time.April, 1, 0, 0, 0, 0, time.UTC)
	march2025 := query.NewTimeInterval(startDate, endDate)

	// Select a spatial extent
	colorado := orb.Polygon{
		{{-109.05, 37.09}, {-102.06, 37.09}, {-102.06, 41.59}, {-109.05, 41.59}, {-109.05, 37.09}},
	}

	// Perform a spatial-temporal query
	// Sentinel2Msi type is generated using tilebox-generate
	var datapointsOverColorado []*examplesv1.Sentinel2Msi
	err = client.Datapoints.QueryInto(ctx,
		[]uuid.UUID{collection.ID}, &datapointsOverColorado,
		datasets.WithTemporalExtent(march2025),
		datasets.WithSpatialExtent(colorado),
	)
	if err != nil {
		log.Fatalf("Failed to query datapoints: %v", err)
	}

	slog.Info("Found datapoints over Colorado in March 2025", slog.Int("count", len(datapointsOverColorado)))
	if len(datapointsOverColorado) > 0 {
		slog.Info("First datapoint over Colorado",
			slog.String("id", datapointsOverColorado[0].GetId().AsUUID().String()),
			slog.Time("event time", datapointsOverColorado[0].GetTime().AsTime()),
			slog.Time("ingestion time", datapointsOverColorado[0].GetIngestionTime().AsTime()),
			slog.String("geometry", wkt.MarshalString(datapointsOverColorado[0].GetGeometry().AsGeometry())),
			slog.String("granule name", datapointsOverColorado[0].GetGranuleName()),
			slog.String("processing level", datapointsOverColorado[0].GetProcessingLevel().String()),
			slog.String("product type", datapointsOverColorado[0].GetProductType()),
			// and so on...
		)
	}
}
