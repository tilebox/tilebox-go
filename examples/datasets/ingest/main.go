package main

import (
	"context"
	"log"
	"log/slog"
	"time"

	"github.com/tilebox/tilebox-go/datasets/v1"
	examplesv1 "github.com/tilebox/tilebox-go/protogen/go/examples/v1"
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
		log.Fatalf("Failed to get dataset: %v", err)
	}

	// Create a collection
	collection, err := client.Collections.Create(ctx, dataset.ID, "My First Collection")
	if err != nil {
		log.Fatalf("Failed to create collection: %v", err)
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
		log.Fatalf("Failed to ingest datapoints: %v", err)
	}
	slog.Info("Ingested datapoints", slog.Int64("created", ingestResponse.NumCreated))

	// Delete datapoints again
	numDeleted, err := client.Datapoints.DeleteIDs(ctx, collection.ID, ingestResponse.DatapointIDs)
	if err != nil {
		log.Fatalf("Failed to delete datapoints: %v", err)
	}
	slog.Info("Deleted datapoints", slog.Int64("deleted", numDeleted))
}
