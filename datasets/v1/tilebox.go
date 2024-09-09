package datasets

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/tilebox/tilebox-go/observability"
	datasetsv1 "github.com/tilebox/tilebox-go/protogen/go/datasets/v1"
	"github.com/tilebox/tilebox-go/protogen/go/datasets/v1/datasetsv1connect"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type datasetServiceConfig struct {
	tracerProvider trace.TracerProvider
	tracerName     string
}

type DatasetServiceOption func(*datasetServiceConfig)

func WithDatasetServiceTracerProvider(tracerProvider trace.TracerProvider) DatasetServiceOption {
	return func(cfg *datasetServiceConfig) {
		cfg.tracerProvider = tracerProvider
	}
}

func WithDatasetServiceTracerName(tracerName string) DatasetServiceOption {
	return func(cfg *datasetServiceConfig) {
		cfg.tracerName = tracerName
	}
}

type DatasetService struct {
	client datasetsv1connect.TileboxServiceClient
	tracer trace.Tracer
}

func newDatasetServiceConfig(options []DatasetServiceOption) *datasetServiceConfig {
	cfg := &datasetServiceConfig{
		tracerProvider: otel.GetTracerProvider(),    // use the global tracer provider by default
		tracerName:     "tilebox.com/observability", // the default tracer name we use
	}
	for _, option := range options {
		option(cfg)
	}

	return cfg
}

func NewDatasetService(client datasetsv1connect.TileboxServiceClient, options ...DatasetServiceOption) *DatasetService {
	cfg := newDatasetServiceConfig(options)
	return &DatasetService{
		client: client,
		tracer: cfg.tracerProvider.Tracer(cfg.tracerName),
	}
}

func (ts *DatasetService) SaveDatapoints(ctx context.Context, collectionID uuid.UUID, datapoints *datasetsv1.Datapoints, enableUpdate bool) (*datasetsv1.SaveDatapointsResponse, error) {
	return observability.WithSpanResult(ctx, ts.tracer, "dataset/save", func(ctx context.Context) (*datasetsv1.SaveDatapointsResponse, error) {
		res, err := ts.client.SaveDatapoints(ctx, connect.NewRequest(
			&datasetsv1.SaveDatapointsRequest{
				CollectionId: &datasetsv1.ID{
					Uuid: collectionID[:],
				},
				Datapoints:   datapoints,
				EnableUpdate: enableUpdate,
			},
		))

		if err != nil {
			return nil, fmt.Errorf("failed to save datapoints: %w", err)
		}

		return res.Msg, nil
	})
}
