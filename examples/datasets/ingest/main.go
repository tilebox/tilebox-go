package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/tilebox/tilebox-go/datasets/v1"
	examplesv1 "github.com/tilebox/tilebox-go/protogen/examples/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func main() {
	ctx := context.Background()

	// Create a Tilebox Datasets client
	client := datasets.NewClient()

	// Select a dataset
	dataset, err := client.Datasets.Get(ctx, "open_data.copernicus.sentinel2_msi")
	if err != nil {
		slog.Error("Failed to get dataset", slog.Any("error", err))
		return
	}

	// Create a collection
	collection, err := client.Collections.Create(ctx, dataset.ID, "My First Collection")
	if err != nil {
		slog.Error("Failed to create collection", slog.Any("error", err))
		return
	}

	// Create datapoints
	// Sentinel2Msi type is generated using tilebox-generate
	datapoints := []*examplesv1.Sentinel2Msi{
		examplesv1.Sentinel2Msi_builder{
			Time:        timestamppb.New(time.Now()),
			GranuleName: proto.String("Granule 1"),
		}.Build(),
		examplesv1.Sentinel2Msi_builder{
			Time:        timestamppb.New(time.Now().Add(-5 * time.Hour)),
			GranuleName: proto.String("Past Granule 2"),
		}.Build(),
	}

	// Ingest datapoints
	ingestResponse, err := client.Datapoints.Ingest(ctx, collection.ID, &datapoints, false)
	if err != nil {
		slog.Error("Failed to ingest datapoints", slog.Any("error", err))
		return
	}
	slog.Info("Ingested datapoints", slog.Int64("created", ingestResponse.NumCreated))

	// Delete datapoints again
	numDeleted, err := client.Datapoints.DeleteIDs(ctx, collection.ID, ingestResponse.DatapointIDs)
	if err != nil {
		slog.Error("Failed to delete datapoints", slog.Any("error", err))
		return
	}
	slog.Info("Deleted datapoints", slog.Int64("deleted", numDeleted))
}
