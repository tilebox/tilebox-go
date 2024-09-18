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

type serviceConfig struct {
	tracerProvider trace.TracerProvider
	tracerName     string
}

type ServiceOption func(*serviceConfig)

func WithServiceTracerProvider(tracerProvider trace.TracerProvider) ServiceOption {
	return func(cfg *serviceConfig) {
		cfg.tracerProvider = tracerProvider
	}
}

func WithServiceTracerName(tracerName string) ServiceOption {
	return func(cfg *serviceConfig) {
		cfg.tracerName = tracerName
	}
}

type Service struct {
	client datasetsv1connect.TileboxServiceClient
	tracer trace.Tracer
}

func newServiceConfig(options []ServiceOption) *serviceConfig {
	cfg := &serviceConfig{
		tracerProvider: otel.GetTracerProvider(),    // use the global tracer provider by default
		tracerName:     "tilebox.com/observability", // the default tracer name we use
	}
	for _, option := range options {
		option(cfg)
	}

	return cfg
}

func NewDatasetsService(client datasetsv1connect.TileboxServiceClient, options ...ServiceOption) *Service {
	cfg := newServiceConfig(options)
	return &Service{
		client: client,
		tracer: cfg.tracerProvider.Tracer(cfg.tracerName),
	}
}

func (s *Service) Ingest(ctx context.Context, collectionID uuid.UUID, datapoints *datasetsv1.Datapoints, allowExisting bool) (*datasetsv1.IngestDatapointsResponse, error) {
	return observability.WithSpanResult(ctx, s.tracer, "datasets/save", func(ctx context.Context) (*datasetsv1.IngestDatapointsResponse, error) {
		res, err := s.client.IngestDatapoints(ctx, connect.NewRequest(
			&datasetsv1.IngestDatapointsRequest{
				CollectionId: &datasetsv1.ID{
					Uuid: collectionID[:],
				},
				Datapoints:    datapoints,
				AllowExisting: allowExisting,
			},
		))

		if err != nil {
			return nil, fmt.Errorf("failed to ingest datapoints: %w", err)
		}

		return res.Msg, nil
	})
}
