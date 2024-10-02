package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	tileboxdatasets "github.com/tilebox/tilebox-go/datasets/v1"
	datasetsv1 "github.com/tilebox/tilebox-go/protogen/go/datasets/v1"
)

func main() {
	ctx := context.Background()

	start := time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)
	loadInterval := tileboxdatasets.NewTimeInterval(start, end)

	client := tileboxdatasets.NewClient(
		tileboxdatasets.WithURL("https://api.tilebox.dev"),
		tileboxdatasets.WithAPIKey(os.Getenv("TILEBOX_STAGING_API_KEY")),
	)

	datasets, err := client.Datasets(ctx)
	if err != nil {
		panic(err)
	}
	/*dataset, err := client.Dataset(ctx, "open_data.copernicus.sentinel1_sar")
	if err != nil {
		panic(err)
	}*/
	dataset := datasets[1]

	_, err = dataset.Collections(ctx)
	if err != nil {
		panic(err)
	}
	/*collection2, err := dataset.CreateCollection(ctx, "corentin")
	if err != nil {
		panic(err)
	}*/
	collection, err := dataset.Collection(ctx, "S1A_EW_GRDM_1S")
	if err != nil {
		panic(err)
	}

	data := collection.Load(ctx, loadInterval, false, false) // raw bytes

	// then depending on what the user want to do
	dataRaw, err := tileboxdatasets.Collect(data)
	if err != nil {
		panic(err)
	}

	dataCopernicus, err := tileboxdatasets.CollectAs[*datasetsv1.CopernicusDataspaceGranule](data)
	if err != nil {
		panic(err)
	}

	/*_, err = collection2.Ingest(ctx, dataRaw, false)
	if err != nil {
		panic(err)
	}*/

	/*validated, err := tileboxdatasets.ValidateAs[*datasetsv1.CopernicusDataspaceGranule](dataCopernicus)
	if err != nil {
		panic(err)
	}*/

	_, err = collection.Delete(ctx, dataRaw)
	if err != nil {
		panic(err)
	}

	slog.InfoContext(ctx, "Datapoints loaded", slog.Int("count", len(dataCopernicus)))
	slog.InfoContext(ctx, "First datapoint",
		slog.String("id", dataCopernicus[0].Meta.ID.String()),
		slog.Time("eventTime", dataCopernicus[0].Meta.EventTime),
		slog.Time("ingestionTime", dataCopernicus[0].Meta.IngestionTime),
	)
}
