package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/wkt"
	tileboxdatasets "github.com/tilebox/tilebox-go/datasets/v1"
	testv1 "github.com/tilebox/tilebox-go/protogen-test/tilebox/v1"
	"github.com/tilebox/tilebox-go/query"
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
	sinceJan2001 := query.NewTimeInterval(jan2001, time.Now())

	// Query one month of data
	var datapoints []*testv1.Modis
	err = client.Datapoints.QueryInto(ctx, []uuid.UUID{collection.ID}, &datapoints, tileboxdatasets.WithTemporalExtent(sinceJan2001))
	if err != nil {
		log.Fatalf("Failed to query datapoints: %v", err)
	}

	slog.Info("Found datapoints", slog.Int("count", len(datapoints)))
	if len(datapoints) > 0 {
		slog.Info("First datapoint",
			slog.String("id", datapoints[0].GetId().AsUUID().String()),
			slog.Time("event time", datapoints[0].GetTime().AsTime()),
			slog.Time("ingestion time", datapoints[0].GetIngestionTime().AsTime()),
			slog.String("geometry", wkt.MarshalString(datapoints[0].GetGeometry().AsGeometry())),
			slog.String("granule name", datapoints[0].GetGranuleName()),
			// slog.String("processing level", datapoints[0].GetProcessingLevel().String()),
			// slog.String("product type", datapoints[0].GetProductType()),
			// and so on...
		)
	}

	// Perform a spatial query
	colorado := orb.Polygon{
		{{-109.05, 37.09}, {-102.06, 37.09}, {-102.06, 41.59}, {-109.05, 41.59}, {-109.05, 37.09}},
	}
	var datapointsOverColorado []*testv1.Modis
	err = client.Datapoints.QueryInto(ctx,
		[]uuid.UUID{collection.ID}, &datapointsOverColorado,
		tileboxdatasets.WithTemporalExtent(sinceJan2001),
		tileboxdatasets.WithSpatialExtent(colorado),
	)
	if err != nil {
		log.Fatalf("Failed to query datapoints: %v", err)
	}

	slog.Info("Found datapoints over Colorado", slog.Int("count", len(datapoints)))
	if len(datapointsOverColorado) > 0 {
		slog.Info("First datapoint over Colorado",
			slog.String("id", datapoints[0].GetId().AsUUID().String()),
			slog.Time("event time", datapoints[0].GetTime().AsTime()),
			slog.Time("ingestion time", datapoints[0].GetIngestionTime().AsTime()),
			slog.String("geometry", wkt.MarshalString(datapoints[0].GetGeometry().AsGeometry())),
			slog.String("granule name", datapoints[0].GetGranuleName()),
			// slog.String("processing level", datapoints[0].GetProcessingLevel().String()),
			// slog.String("product type", datapoints[0].GetProductType()),
			// and so on...
		)
	}
}
