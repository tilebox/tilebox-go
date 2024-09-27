package datasets

import (
	"context"
	"errors"
	"fmt"
	"iter"
	"runtime/debug"
	"strings"
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

type LoadInterval interface {
	ToProtoTimeInterval() *datasetsv1.TimeInterval
	ToProtoDatapointInterval() *datasetsv1.DatapointInterval
}

var _ LoadInterval = &TimeInterval{}

type TimeInterval struct {
	start time.Time
	end   time.Time
}

func NewTimeInterval(start, end time.Time) *TimeInterval {
	return &TimeInterval{
		start: start,
		end:   end,
	}
}

func (t *TimeInterval) ToProtoTimeInterval() *datasetsv1.TimeInterval {
	return &datasetsv1.TimeInterval{
		StartTime:      timestamppb.New(t.start),
		EndTime:        timestamppb.New(t.end),
		StartExclusive: false,
		EndInclusive:   true,
	}
}

func (t *TimeInterval) ToProtoDatapointInterval() *datasetsv1.DatapointInterval {
	return nil
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

type Service interface {
	GetDataset(ctx context.Context, datasetID uuid.UUID) (*datasetsv1.Dataset, error)
	ListDatasets(ctx context.Context) (*datasetsv1.ListDatasetsResponse, error)
	CreateCollection(ctx context.Context, datasetID uuid.UUID, collectionName string) (*datasetsv1.CollectionInfo, error)
	GetCollections(ctx context.Context, datasetID uuid.UUID) (*datasetsv1.Collections, error)
	GetCollectionByName(ctx context.Context, datasetID uuid.UUID, collectionName string) (*datasetsv1.CollectionInfo, error)
	Load(ctx context.Context, collectionID uuid.UUID, loadInterval LoadInterval, skipData bool, skipMeta bool) iter.Seq2[*Datapoint, error]
	GetDatasetForInterval(ctx context.Context, collectionID uuid.UUID, timeInterval *datasetsv1.TimeInterval, datapointInterval *datasetsv1.DatapointInterval, page *datasetsv1.Pagination, skipData bool, skipMeta bool) (*datasetsv1.Datapoints, error)
	GetDatapointByID(ctx context.Context, collectionID uuid.UUID, datapointID uuid.UUID, skipData bool) (*datasetsv1.Datapoint, error)
	Ingest(ctx context.Context, collectionID uuid.UUID, datapoints *datasetsv1.Datapoints, allowExisting bool) (*datasetsv1.IngestDatapointsResponse, error)
	Delete(ctx context.Context, collectionID uuid.UUID, datapointIDs []uuid.UUID) (*datasetsv1.DeleteDatapointsResponse, error)
}

var _ Service = &service{}

type service struct {
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

func NewDatasetsService(client datasetsv1connect.TileboxServiceClient, options ...ServiceOption) Service {
	cfg := newServiceConfig(options)
	return &service{
		client: client,
		tracer: cfg.tracerProvider.Tracer(cfg.tracerName),
	}
}

func (s *service) GetDataset(ctx context.Context, datasetID uuid.UUID) (*datasetsv1.Dataset, error) {
	return observability.WithSpanResult(ctx, s.tracer, "datasets/get_dataset", func(ctx context.Context) (*datasetsv1.Dataset, error) {
		res, err := s.client.GetDataset(ctx, connect.NewRequest(
			&datasetsv1.GetDatasetRequest{
				DatasetId: datasetID.String(),
			},
		))

		if err != nil {
			return nil, fmt.Errorf("failed to get dataset: %w", err)
		}

		return res.Msg, nil
	})
}

func (s *service) ListDatasets(ctx context.Context) (*datasetsv1.ListDatasetsResponse, error) {
	return observability.WithSpanResult(ctx, s.tracer, "datasets/list_datasets", func(ctx context.Context) (*datasetsv1.ListDatasetsResponse, error) {
		res, err := s.client.ListDatasets(ctx, connect.NewRequest(
			&datasetsv1.ListDatasetsRequest{
				ClientInfo: clientInfo(),
			},
		))

		if err != nil {
			return nil, fmt.Errorf("failed to list datasets: %w", err)
		}

		return res.Msg, nil
	})
}

func clientInfo() *datasetsv1.ClientInfo {
	var packages []*datasetsv1.Package
	buildInfo, ok := debug.ReadBuildInfo()
	if ok {
		for _, dep := range buildInfo.Deps {
			if strings.HasPrefix(dep.Path, "github.com/tilebox/") {
				packages = append(packages, &datasetsv1.Package{
					Name:    dep.Path,
					Version: dep.Version,
				})
			}
		}
	}

	return &datasetsv1.ClientInfo{
		Name:        "Go",
		Environment: "Tilebox Go Client",
		Packages:    packages,
	}
}

func (s *service) CreateCollection(ctx context.Context, datasetID uuid.UUID, collectionName string) (*datasetsv1.CollectionInfo, error) {
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

func (s *service) GetCollections(ctx context.Context, datasetID uuid.UUID) (*datasetsv1.Collections, error) {
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

func (s *service) GetCollectionByName(ctx context.Context, datasetID uuid.UUID, collectionName string) (*datasetsv1.CollectionInfo, error) {
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

func (s *service) Load(ctx context.Context, collectionID uuid.UUID, loadInterval LoadInterval, skipData, skipMeta bool) iter.Seq2[*Datapoint, error] {
	return func(yield func(*Datapoint, error) bool) {
		var page *datasetsv1.Pagination // nil for the first request

		timeInterval := loadInterval.ToProtoTimeInterval()
		datapointInterval := loadInterval.ToProtoDatapointInterval()

		if timeInterval == nil && datapointInterval == nil {
			yield(nil, errors.New("time interval and datapoint interval cannot both be nil"))
			return
		}

		for {
			datapointsMessage, err := s.GetDatasetForInterval(ctx, collectionID, timeInterval, datapointInterval, page, skipData, skipMeta)
			if err != nil {
				yield(nil, err)
				return
			}

			// if skipMeta is true, datapointsMessage.GetMeta() is not nil and contains the datapoint ids
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

func (s *service) GetDatasetForInterval(ctx context.Context, collectionID uuid.UUID, timeInterval *datasetsv1.TimeInterval, datapointInterval *datasetsv1.DatapointInterval, page *datasetsv1.Pagination, skipData, skipMeta bool) (*datasetsv1.Datapoints, error) {
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

func (s *service) GetDatapointByID(ctx context.Context, collectionID uuid.UUID, datapointID uuid.UUID, skipData bool) (*datasetsv1.Datapoint, error) {
	return observability.WithSpanResult(ctx, s.tracer, "datasets/get_datapoint_by_id", func(ctx context.Context) (*datasetsv1.Datapoint, error) {
		res, err := s.client.GetDatapointByID(ctx, connect.NewRequest(
			&datasetsv1.GetDatapointByIdRequest{
				CollectionId: collectionID.String(),
				Id:           datapointID.String(),
				SkipData:     skipData,
			},
		))

		if err != nil {
			return nil, fmt.Errorf("failed to get datapoint by id: %w", err)
		}

		return res.Msg, nil
	})
}

func (s *service) Ingest(ctx context.Context, collectionID uuid.UUID, datapoints *datasetsv1.Datapoints, allowExisting bool) (*datasetsv1.IngestDatapointsResponse, error) {
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

func (s *service) Delete(ctx context.Context, collectionID uuid.UUID, datapointIDs []uuid.UUID) (*datasetsv1.DeleteDatapointsResponse, error) {
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
