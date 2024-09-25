package datasets

import (
	"context"
	"fmt"
	"iter"
	"time"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/tilebox/tilebox-go/observability"
	datasetsv1 "github.com/tilebox/tilebox-go/protogen/go/datasets/v1"
	"github.com/tilebox/tilebox-go/protogen/go/datasets/v1/datasetsv1connect"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/types/known/timestamppb"
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

func (s *Service) CreateCollection(ctx context.Context, datasetID uuid.UUID, collectionName string) (*datasetsv1.CollectionInfo, error) {
	return observability.WithSpanResult(ctx, s.tracer, "datasets/create_collection", func(ctx context.Context) (*datasetsv1.CollectionInfo, error) {
		res, err := s.client.CreateCollection(ctx, connect.NewRequest(
			&datasetsv1.CreateCollectionRequest{
				DatasetId: &datasetsv1.ID{
					Uuid: datasetID[:],
				},
				Name: collectionName,
			},
		))

		if err != nil {
			return nil, fmt.Errorf("failed to create collection: %w", err)
		}

		return res.Msg, nil
	})
}

func (s *Service) GetCollections(ctx context.Context, datasetID uuid.UUID) (*datasetsv1.Collections, error) {
	return observability.WithSpanResult(ctx, s.tracer, "datasets/get_collections", func(ctx context.Context) (*datasetsv1.Collections, error) {
		res, err := s.client.GetCollections(ctx, connect.NewRequest(
			&datasetsv1.GetCollectionsRequest{
				DatasetId: &datasetsv1.ID{
					Uuid: datasetID[:],
				},
				WithAvailability: true,
				WithCount:        true,
			},
		))

		if err != nil {
			return nil, fmt.Errorf("failed to get collections: %w", err)
		}

		return res.Msg, nil
	})
}

func (s *Service) GetCollectionByName(ctx context.Context, datasetID uuid.UUID, collectionName string) (*datasetsv1.CollectionInfo, error) {
	return observability.WithSpanResult(ctx, s.tracer, "datasets/get_collection_by_name", func(ctx context.Context) (*datasetsv1.CollectionInfo, error) {
		res, err := s.client.GetCollectionByName(ctx, connect.NewRequest(
			&datasetsv1.GetCollectionByNameRequest{
				CollectionName:   collectionName,
				WithAvailability: true,
				WithCount:        true,
				DatasetId: &datasetsv1.ID{
					Uuid: datasetID[:],
				},
			},
		))

		if err != nil {
			return nil, fmt.Errorf("failed to get collections: %w", err)
		}

		return res.Msg, nil
	})
}

type DatapointMetadata struct {
	ID            uuid.UUID
	EventTime     time.Time
	IngestionTime time.Time
}

type Datapoint struct {
	Meta *DatapointMetadata
	// data proto.Message
}

func (s *Service) Load(ctx context.Context, collectionID uuid.UUID, start time.Time, end time.Time, skipData, skipMeta bool) iter.Seq2[*Datapoint, error] {
	return func(yield func(*Datapoint, error) bool) {
		var page *datasetsv1.Pagination // nil for the first request

		timeInterval := &datasetsv1.TimeInterval{
			StartTime:      timestamppb.New(start),
			EndTime:        timestamppb.New(end),
			StartExclusive: false,
			EndInclusive:   true,
		}

		for {
			datapointsMessage, err := s.GetDatasetForInterval(ctx, collectionID, timeInterval, nil, page, skipData, skipMeta)
			if err != nil {
				yield(nil, err)
				return
			}

			for _, dp := range datapointsMessage.GetMeta() {
				datapointID, err := uuid.Parse(dp.GetId())
				if err != nil {
					yield(nil, err)
					return
				}

				meta := &DatapointMetadata{
					ID: datapointID,
				}

				if !skipMeta {
					meta.EventTime = dp.GetEventTime().AsTime()
					meta.IngestionTime = dp.GetIngestionTime().AsTime()
				}

				// if !skipData {
				// FIXME
				// }

				datapoint := &Datapoint{
					Meta: meta,
					// data: data,
				}
				if !yield(datapoint, nil) {
					return
				}
			}

			page = datapointsMessage.GetNextPage()
			if page == nil {
				break
			}
		}
	}
}

func Collect[K any](seq iter.Seq2[K, error]) ([]K, error) {
	s := make([]K, 0)

	for k, err := range seq {
		if err != nil {
			return nil, err
		}
		s = append(s, k)
	}
	return s, nil
}

func (s *Service) GetDatasetForInterval(ctx context.Context, collectionID uuid.UUID, timeInterval *datasetsv1.TimeInterval, datapointInterval *datasetsv1.DatapointInterval, page *datasetsv1.Pagination, skipData, skipMeta bool) (*datasetsv1.Datapoints, error) {
	return observability.WithSpanResult(ctx, s.tracer, "datasets/get_dataset_for_interval", func(ctx context.Context) (*datasetsv1.Datapoints, error) {
		res, err := s.client.GetDatasetForInterval(ctx, connect.NewRequest(
			&datasetsv1.GetDatasetForIntervalRequest{
				CollectionId:      collectionID.String(),
				TimeInterval:      timeInterval,
				DatapointInterval: datapointInterval,
				Page:              page,
				SkipData:          skipData,
				SkipMeta:          skipMeta,
			},
		))

		if err != nil {
			return nil, fmt.Errorf("failed to get dataset for interval: %w", err)
		}

		return res.Msg, nil
	})
}

func (s *Service) Ingest(ctx context.Context, collectionID uuid.UUID, datapoints *datasetsv1.Datapoints, allowExisting bool) (*datasetsv1.IngestDatapointsResponse, error) {
	return observability.WithSpanResult(ctx, s.tracer, "datasets/ingest_datapoints", func(ctx context.Context) (*datasetsv1.IngestDatapointsResponse, error) {
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

func (s *Service) Delete(ctx context.Context, collectionID uuid.UUID, datapointIDs []uuid.UUID) (*datasetsv1.DeleteDatapointsResponse, error) {
	return observability.WithSpanResult(ctx, s.tracer, "datasets/delete_datapoints", func(ctx context.Context) (*datasetsv1.DeleteDatapointsResponse, error) {
		res, err := s.client.DeleteDatapoints(ctx, connect.NewRequest(
			&datasetsv1.DeleteDatapointsRequest{
				CollectionId: &datasetsv1.ID{
					Uuid: collectionID[:],
				},
				DatapointIds: lo.Map(datapointIDs, func(datapointID uuid.UUID, _ int) *datasetsv1.ID {
					return &datasetsv1.ID{
						Uuid: datapointID[:],
					}
				}),
			},
		))

		if err != nil {
			return nil, fmt.Errorf("failed to delete datapoints: %w", err)
		}

		return res.Msg, nil
	})
}
