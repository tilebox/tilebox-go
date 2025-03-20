package datasets // import "github.com/tilebox/tilebox-go/datasets/v1"

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

type DatasetService interface {
	GetDataset(ctx context.Context, slug string) (*datasetsv1.Dataset, error)
	ListDatasets(ctx context.Context) (*datasetsv1.ListDatasetsResponse, error)
}

var _ DatasetService = &datasetService{}

type datasetService struct {
	datasetClient datasetsv1connect.DatasetServiceClient
	tracer        trace.Tracer
}

func newDatasetsService(datasetClient datasetsv1connect.DatasetServiceClient, tracer trace.Tracer) DatasetService {
	return &datasetService{
		datasetClient: datasetClient,
		tracer:        tracer,
	}
}

func (s *datasetService) GetDataset(ctx context.Context, slug string) (*datasetsv1.Dataset, error) {
	return observability.WithSpanResult(ctx, s.tracer, "datasets/get", func(ctx context.Context) (*datasetsv1.Dataset, error) {
		res, err := s.datasetClient.GetDataset(ctx, connect.NewRequest(
			&datasetsv1.GetDatasetRequest{
				Slug: slug,
			},
		))
		if err != nil {
			return nil, fmt.Errorf("failed to get dataset: %w", err)
		}

		return res.Msg, nil
	})
}

func (s *datasetService) ListDatasets(ctx context.Context) (*datasetsv1.ListDatasetsResponse, error) {
	return observability.WithSpanResult(ctx, s.tracer, "datasets/list", func(ctx context.Context) (*datasetsv1.ListDatasetsResponse, error) {
		res, err := s.datasetClient.ListDatasets(ctx, connect.NewRequest(
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

type CollectionService interface {
	CreateCollection(ctx context.Context, datasetID uuid.UUID, collectionName string) (*datasetsv1.CollectionInfo, error)
	GetCollectionByName(ctx context.Context, datasetID uuid.UUID, collectionName string) (*datasetsv1.CollectionInfo, error)
	ListCollections(ctx context.Context, datasetID uuid.UUID) (*datasetsv1.CollectionInfos, error)
}

var _ CollectionService = &collectionService{}

type collectionService struct {
	collectionClient datasetsv1connect.CollectionServiceClient
	tracer           trace.Tracer
}

func newCollectionService(collectionClient datasetsv1connect.CollectionServiceClient, tracer trace.Tracer) CollectionService {
	return &collectionService{
		collectionClient: collectionClient,
		tracer:           tracer,
	}
}

func (s *collectionService) CreateCollection(ctx context.Context, datasetID uuid.UUID, collectionName string) (*datasetsv1.CollectionInfo, error) {
	return observability.WithSpanResult(ctx, s.tracer, "datasets/collections/create", func(ctx context.Context) (*datasetsv1.CollectionInfo, error) {
		res, err := s.collectionClient.CreateCollection(ctx, connect.NewRequest(
			&datasetsv1.CreateCollectionRequest{
				DatasetId: uuidToProtobuf(datasetID),
				Name:      collectionName,
			},
		))
		if err != nil {
			return nil, fmt.Errorf("failed to create collection: %w", err)
		}

		return res.Msg, nil
	})
}

func (s *collectionService) GetCollectionByName(ctx context.Context, datasetID uuid.UUID, collectionName string) (*datasetsv1.CollectionInfo, error) {
	return observability.WithSpanResult(ctx, s.tracer, "datasets/collections/get", func(ctx context.Context) (*datasetsv1.CollectionInfo, error) {
		res, err := s.collectionClient.GetCollectionByName(ctx, connect.NewRequest(
			&datasetsv1.GetCollectionByNameRequest{
				CollectionName:   collectionName,
				WithAvailability: true,
				WithCount:        true,
				DatasetId:        uuidToProtobuf(datasetID),
			},
		))
		if err != nil {
			return nil, fmt.Errorf("failed to get collections: %w", err)
		}

		return res.Msg, nil
	})
}

func (s *collectionService) ListCollections(ctx context.Context, datasetID uuid.UUID) (*datasetsv1.CollectionInfos, error) {
	return observability.WithSpanResult(ctx, s.tracer, "datasets/collections/list", func(ctx context.Context) (*datasetsv1.CollectionInfos, error) {
		res, err := s.collectionClient.ListCollections(ctx, connect.NewRequest(
			&datasetsv1.ListCollectionsRequest{
				DatasetId:        uuidToProtobuf(datasetID),
				WithAvailability: true,
				WithCount:        true,
			},
		))
		if err != nil {
			return nil, fmt.Errorf("failed to list collections: %w", err)
		}

		return res.Msg, nil
	})
}

type DataAccessService interface {
	GetDatasetForInterval(ctx context.Context, collectionID uuid.UUID, timeInterval *datasetsv1.TimeInterval, datapointInterval *datasetsv1.DatapointInterval, page *datasetsv1.Pagination, skipData bool, skipMeta bool) (*datasetsv1.DatapointPage, error)
	GetDatapointByID(ctx context.Context, collectionID uuid.UUID, datapointID uuid.UUID, skipData bool) (*datasetsv1.Datapoint, error)
}

var _ DataAccessService = &dataAccessService{}

type dataAccessService struct {
	dataAccessClient datasetsv1connect.DataAccessServiceClient
	tracer           trace.Tracer
}

func newDataAccessService(dataAccessClient datasetsv1connect.DataAccessServiceClient, tracer trace.Tracer) DataAccessService {
	return &dataAccessService{
		dataAccessClient: dataAccessClient,
		tracer:           tracer,
	}
}

func (s *dataAccessService) GetDatasetForInterval(ctx context.Context, collectionID uuid.UUID, timeInterval *datasetsv1.TimeInterval, datapointInterval *datasetsv1.DatapointInterval, page *datasetsv1.Pagination, skipData, skipMeta bool) (*datasetsv1.DatapointPage, error) {
	return observability.WithSpanResult(ctx, s.tracer, "datasets/datapoints/load", func(ctx context.Context) (*datasetsv1.DatapointPage, error) {
		res, err := s.dataAccessClient.GetDatasetForInterval(ctx, connect.NewRequest(
			&datasetsv1.GetDatasetForIntervalRequest{
				CollectionId:      collectionID.String(),
				TimeInterval:      timeInterval,
				DatapointInterval: datapointInterval,
				Page:              paginationToLegacyPagination(page),
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

func (s *dataAccessService) GetDatapointByID(ctx context.Context, collectionID uuid.UUID, datapointID uuid.UUID, skipData bool) (*datasetsv1.Datapoint, error) {
	return observability.WithSpanResult(ctx, s.tracer, "datasets/datapoints/get", func(ctx context.Context) (*datasetsv1.Datapoint, error) {
		res, err := s.dataAccessClient.GetDatapointByID(ctx, connect.NewRequest(
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

type DataIngestionService interface {
	// IngestDatapoints is the legacy name for where datapoints are separated into Meta/Data. change this to Ingest() in the future
	IngestDatapoints(ctx context.Context, collectionID uuid.UUID, datapoints *datasetsv1.Datapoints, allowExisting bool) (*datasetsv1.IngestResponse, error)
	Delete(ctx context.Context, collectionID uuid.UUID, datapointIDs []uuid.UUID) (*datasetsv1.DeleteResponse, error)
}

var _ DataIngestionService = &dataIngestionService{}

type dataIngestionService struct {
	dataIngestionClient datasetsv1connect.DataIngestionServiceClient
	tracer              trace.Tracer
}

func newDataIngestionService(dataIngestionClient datasetsv1connect.DataIngestionServiceClient, tracer trace.Tracer) DataIngestionService {
	return &dataIngestionService{
		dataIngestionClient: dataIngestionClient,
		tracer:              tracer,
	}
}

func (s *dataIngestionService) IngestDatapoints(ctx context.Context, collectionID uuid.UUID, datapoints *datasetsv1.Datapoints, allowExisting bool) (*datasetsv1.IngestResponse, error) {
	return observability.WithSpanResult(ctx, s.tracer, "datasets/datapoints/ingest", func(ctx context.Context) (*datasetsv1.IngestResponse, error) {
		res, err := s.dataIngestionClient.IngestDatapoints(ctx, connect.NewRequest(
			&datasetsv1.IngestDatapointsRequest{
				CollectionId:  uuidToProtobuf(collectionID),
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

func (s *dataIngestionService) Delete(ctx context.Context, collectionID uuid.UUID, datapointIDs []uuid.UUID) (*datasetsv1.DeleteResponse, error) {
	return observability.WithSpanResult(ctx, s.tracer, "datasets/datapoints/delete", func(ctx context.Context) (*datasetsv1.DeleteResponse, error) {
		res, err := s.dataIngestionClient.Delete(ctx, connect.NewRequest(
			&datasetsv1.DeleteRequest{
				CollectionId: uuidToProtobuf(collectionID),
				DatapointIds: lo.Map(datapointIDs, func(datapointID uuid.UUID, _ int) *datasetsv1.ID {
					return uuidToProtobuf(datapointID)
				}),
			},
		))
		if err != nil {
			return nil, fmt.Errorf("failed to delete datapoints: %w", err)
		}

		return res.Msg, nil
	})
}

func paginationToLegacyPagination(pagination *datasetsv1.Pagination) *datasetsv1.LegacyPagination {
	if pagination == nil {
		return nil
	}
	p := &datasetsv1.LegacyPagination{
		Limit: pagination.Limit, //nolint:protogetter // do not use GetLimit() since that converts nil to 0
	}
	id, err := protoToUUID(pagination.GetStartingAfter())
	if err != nil {
		id = uuid.Nil
	}
	if id != uuid.Nil {
		idAsString := id.String()
		p.StartingAfter = &idAsString
	}
	return p
}

func paginationFromLegacyPagination(pagination *datasetsv1.LegacyPagination) *datasetsv1.Pagination {
	if pagination == nil {
		return nil
	}
	p := &datasetsv1.Pagination{
		Limit: pagination.Limit, //nolint:protogetter // do not use GetLimit() since that converts nil to 0
	}
	if pagination.StartingAfter != nil {
		id, err := uuid.Parse(pagination.GetStartingAfter())
		if err == nil && id != uuid.Nil {
			p.StartingAfter = uuidToProtobuf(id)
		}
	}
	return p
}
