package main

import (
	"context"
	"log/slog"

	"github.com/tilebox/tilebox-go/datasets/v1"
	"github.com/tilebox/tilebox-go/datasets/v1/field"
)

func main() {
	ctx := context.Background()

	// Create a Tilebox Datasets client
	client := datasets.NewClient()

	// Create a dataset
	dataset, err := client.Datasets.CreateOrUpdate(ctx,
		datasets.KindSpatiotemporal,
		"Personal Landsat-8 catalog",
		"my_landsat8_oli_tirs",
		[]datasets.Field{
			field.String("granule_name").Description("The name of the Earth View Granule").ExampleValue("LC08_L1GT_037077_20220830_20220910_02_T2"),
			field.String("platform").Description("Landsat satellite number.").ExampleValue("LANDSAT_8"),
			field.Int64("proj_shape").Repeated().Description("Raster shape").ExampleValue("[7681, 7591]"),
		},
	)
	if err != nil {
		slog.Error("Failed to create dataset", slog.Any("error", err))
		return
	}
	slog.InfoContext(ctx, "Created dataset", slog.String("dataset_id", dataset.ID.String()))

	// Add a field to the dataset
	dataset, err = client.Datasets.CreateOrUpdate(ctx,
		datasets.KindSpatiotemporal,
		"Personal Landsat-8 catalog",
		"my_landsat8_oli_tirs",
		[]datasets.Field{
			field.String("granule_name").Description("The name of the Earth View Granule").ExampleValue("LC08_L1GT_037077_20220830_20220910_02_T2"),
			field.String("platform").Description("Landsat satellite number.").ExampleValue("LANDSAT_8"),
			field.Int64("proj_shape").Repeated().Description("Raster shape").ExampleValue("[7681, 7591]"),
			field.Uint64("proj_epsg").Description("Projection EPSG code").ExampleValue("32610"),
		},
	)
	if err != nil {
		slog.Error("Failed to update dataset", slog.Any("error", err))
		return
	}

	slog.InfoContext(ctx, "Updated dataset, added proj_epsg field", slog.String("dataset_id", dataset.ID.String()))
}
