package datasets

import (
	"context"
	"fmt"
	"runtime/debug"
	"strings"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/tilebox/tilebox-go/observability"
	datasetsv1 "github.com/tilebox/tilebox-go/protogen/go/datasets/v1"
	"github.com/tilebox/tilebox-go/protogen/go/datasets/v1/datasetsv1connect"
	"go.opentelemetry.io/otel/trace"
)

type Service interface {
	GetDataset(ctx context.Context, datasetID uuid.UUID) (*datasetsv1.Dataset, error)
	ListDatasets(ctx context.Context) (*datasetsv1.ListDatasetsResponse, error)
	CreateCollection(ctx context.Context, datasetID uuid.UUID, collectionName string) (*datasetsv1.CollectionInfo, error)
	GetCollections(ctx context.Context, datasetID uuid.UUID) (*datasetsv1.Collections, error)
	GetCollectionByName(ctx context.Context, datasetID uuid.UUID, collectionName string) (*datasetsv1.CollectionInfo, error)
	GetDatasetForInterval(ctx context.Context, collectionID uuid.UUID, timeInterval *datasetsv1.TimeInterval, datapointInterval *datasetsv1.DatapointInterval, page *datasetsv1.Pagination, skipData bool, skipMeta bool) (*datasetsv1.Datapoints, error)
	GetDatapointByID(ctx context.Context, collectionID uuid.UUID, datapointID uuid.UUID, skipData bool) (*datasetsv1.Datapoint, error)
	IngestDatapoints(ctx context.Context, collectionID uuid.UUID, datapoints *datasetsv1.Datapoints, allowExisting bool) (*datasetsv1.IngestDatapointsResponse, error)
	DeleteDatapoints(ctx context.Context, collectionID uuid.UUID, datapointIDs []uuid.UUID) (*datasetsv1.DeleteDatapointsResponse, error)
}

var _ Service = &service{}

type service struct {
	client datasetsv1connect.TileboxServiceClient
	tracer trace.Tracer
}

func newDatasetsService(client datasetsv1connect.TileboxServiceClient, tracer trace.Tracer) Service {
	return &service{
		client: client,
		tracer: tracer,
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

func (s *service) IngestDatapoints(ctx context.Context, collectionID uuid.UUID, datapoints *datasetsv1.Datapoints, allowExisting bool) (*datasetsv1.IngestDatapointsResponse, error) {
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

func (s *service) DeleteDatapoints(ctx context.Context, collectionID uuid.UUID, datapointIDs []uuid.UUID) (*datasetsv1.DeleteDatapointsResponse, error) {
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
