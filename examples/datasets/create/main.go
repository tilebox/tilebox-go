package main

import (
	"context"
	"log/slog"

	"github.com/tilebox/tilebox-go/datasets/v1"
	"github.com/tilebox/tilebox-go/datasets/v1/field"
	stacv1 "github.com/tilebox/tilebox-go/protogen/datasets/stac/v1"
)

func main() {
	ctx := context.Background()

	// Create a Tilebox Datasets client
	client := datasets.NewClient()

	fields := make([]datasets.Field, 0, 5)
	fields = append(fields,
		field.String("granule_name").
			Description("The source granule name used as the primary title of the STAC item.").
			ExampleValue("S2A_MSIL2A_20260521T104031_N0511_R008_T32TQM_20260521T132145").
			SourceJSONPointer("/properties/granule_name").
			Queryable().
			Roles(field.RolePrimaryTitle),
		field.Message("assets", &stacv1.Assets{}).
			Description("The STAC assets associated with the item.").
			SourceJSONPointer("/assets"),
		field.Message("providers", &stacv1.Provider{}).
			Description("The organizations that produced, processed, or hosted the item.").
			SourceJSONPointer("/properties/providers").
			Repeated(),
		field.Float64("eo_cloud_cover").
			Description("The percentage of the item covered by clouds.").
			ExampleValue("12.5").
			SourceJSONPointer("/properties/eo:cloud_cover").
			Queryable().
			JSONSchemaRef("https://stac-extensions.github.io/eo/v2.0.0/schema.json#/definitions/eo:cloud_cover"),
	)

	// Create a STAC-compatible catalog dataset.
	dataset, err := client.Datasets.CreateOrUpdate(ctx,
		datasets.KindSpatiotemporal,
		"sentinel2_l2a_catalog",
		"Sentinel-2 L2A STAC catalog",
		fields,
	)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to create dataset", slog.Any("error", err))
		return
	}
	slog.InfoContext(ctx, "Created dataset", slog.String("dataset_id", dataset.ID.String()))

	// Add another STAC property. Schema updates include all existing custom fields.
	fields = append(fields,
		field.String("proj_code").
			Description("The coordinate reference system of the item geometry or assets.").
			ExampleValue("EPSG:32632").
			SourceJSONPointer("/properties/proj:code").
			JSONSchemaRef("https://stac-extensions.github.io/projection/v2.0.0/schema.json#/definitions/fields/properties/proj:code"),
	)
	dataset, err = client.Datasets.CreateOrUpdate(ctx,
		datasets.KindSpatiotemporal,
		"sentinel2_l2a_catalog",
		"Sentinel-2 L2A STAC catalog",
		fields,
	)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to update dataset", slog.Any("error", err))
		return
	}

	slog.InfoContext(ctx, "Updated dataset, added proj_code field", slog.String("dataset_id", dataset.ID.String()))
}
