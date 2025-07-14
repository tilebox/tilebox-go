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
			datasetsv1.GetDatasetRequest_builder{
				Slug: slug,
			}.Build(),
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
			datasetsv1.ListDatasetsRequest_builder{
				ClientInfo: clientInfo(),
			}.Build(),
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
				packages = append(packages, datasetsv1.Package_builder{
					Name:    dep.Path,
					Version: dep.Version,
				}.Build())
			}
		}
	}

	return datasetsv1.ClientInfo_builder{
		Name:        "Go",
		Environment: "Tilebox Go Client",
		Packages:    packages,
	}.Build()
}

type CollectionService interface {
	CreateCollection(ctx context.Context, datasetID uuid.UUID, collectionName string) (*datasetsv1.CollectionInfo, error)
	GetCollectionByName(ctx context.Context, datasetID uuid.UUID, collectionName string) (*datasetsv1.CollectionInfo, error)
	DeleteCollection(ctx context.Context, datasetID uuid.UUID, collectionID uuid.UUID) error
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

func (s *collectionService) CreateCollection(ctx context.Context, datasetID uuid.UUID, name string) (*datasetsv1.CollectionInfo, error) {
	return observability.WithSpanResult(ctx, s.tracer, "datasets/collections/create", func(ctx context.Context) (*datasetsv1.CollectionInfo, error) {
		res, err := s.collectionClient.CreateCollection(ctx, connect.NewRequest(
			datasetsv1.CreateCollectionRequest_builder{
				DatasetId: uuidToProtobuf(datasetID),
				Name:      name,
			}.Build(),
		))
		if err != nil {
			return nil, fmt.Errorf("failed to create collection: %w", err)
		}

		return res.Msg, nil
	})
}

func (s *collectionService) GetCollectionByName(ctx context.Context, datasetID uuid.UUID, name string) (*datasetsv1.CollectionInfo, error) {
	return observability.WithSpanResult(ctx, s.tracer, "datasets/collections/get", func(ctx context.Context) (*datasetsv1.CollectionInfo, error) {
		res, err := s.collectionClient.GetCollectionByName(ctx, connect.NewRequest(
			datasetsv1.GetCollectionByNameRequest_builder{
				CollectionName:   name,
				WithAvailability: true,
				WithCount:        true,
				DatasetId:        uuidToProtobuf(datasetID),
			}.Build(),
		))
		if err != nil {
			return nil, fmt.Errorf("failed to get collection: %w", err)
		}

		return res.Msg, nil
	})
}

func (s *collectionService) DeleteCollection(ctx context.Context, datasetID uuid.UUID, collectionID uuid.UUID) error {
	return observability.WithSpan(ctx, s.tracer, "datasets/collections/delete", func(ctx context.Context) error {
		_, err := s.collectionClient.DeleteCollection(ctx, connect.NewRequest(
			datasetsv1.DeleteCollectionRequest_builder{
				CollectionId: uuidToProtobuf(collectionID),
				DatasetId:    uuidToProtobuf(datasetID),
			}.Build(),
		))
		if err != nil {
			return fmt.Errorf("failed to delete collection: %w", err)
		}

		return nil
	})
}

func (s *collectionService) ListCollections(ctx context.Context, datasetID uuid.UUID) (*datasetsv1.CollectionInfos, error) {
	return observability.WithSpanResult(ctx, s.tracer, "datasets/collections/list", func(ctx context.Context) (*datasetsv1.CollectionInfos, error) {
		res, err := s.collectionClient.ListCollections(ctx, connect.NewRequest(
			datasetsv1.ListCollectionsRequest_builder{
				DatasetId:        uuidToProtobuf(datasetID),
				WithAvailability: true,
				WithCount:        true,
			}.Build(),
		))
		if err != nil {
			return nil, fmt.Errorf("failed to list collections: %w", err)
		}

		return res.Msg, nil
	})
}

type DataAccessService interface {
	Query(ctx context.Context, collectionIDs []uuid.UUID, filters *datasetsv1.QueryFilters, page *datasetsv1.Pagination, skipData bool) (*datasetsv1.QueryResultPage, error)
	QueryByID(ctx context.Context, collectionIDs []uuid.UUID, datapointID uuid.UUID, skipData bool) (*datasetsv1.Any, error)
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

func (s *dataAccessService) Query(ctx context.Context, collectionIDs []uuid.UUID, filters *datasetsv1.QueryFilters, page *datasetsv1.Pagination, skipData bool) (*datasetsv1.QueryResultPage, error) {
	return observability.WithSpanResult(ctx, s.tracer, "datasets/datapoints/query", func(ctx context.Context) (*datasetsv1.QueryResultPage, error) {
		res, err := s.dataAccessClient.Query(ctx, connect.NewRequest(
			datasetsv1.QueryRequest_builder{
				CollectionIds: uuidsToProtobuf(collectionIDs),
				Filters:       filters,
				Page:          page,
				SkipData:      skipData,
			}.Build(),
		))
		if err != nil {
			return nil, fmt.Errorf("failed to query datpoints: %w", err)
		}

		return res.Msg, nil
	})
}

func (s *dataAccessService) QueryByID(ctx context.Context, collectionIDs []uuid.UUID, datapointID uuid.UUID, skipData bool) (*datasetsv1.Any, error) {
	return observability.WithSpanResult(ctx, s.tracer, "datasets/datapoints/get", func(ctx context.Context) (*datasetsv1.Any, error) {
		res, err := s.dataAccessClient.QueryByID(ctx, connect.NewRequest(
			datasetsv1.QueryByIDRequest_builder{
				CollectionIds: uuidsToProtobuf(collectionIDs),
				Id:            uuidToProtobuf(datapointID),
				SkipData:      skipData,
			}.Build(),
		))
		if err != nil {
			return nil, fmt.Errorf("failed to query datapoint by id: %w", err)
		}

		return res.Msg, nil
	})
}

type DataIngestionService interface {
	Ingest(ctx context.Context, collectionID uuid.UUID, datapoints [][]byte, allowExisting bool) (*datasetsv1.IngestResponse, error)
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

func (s *dataIngestionService) Ingest(ctx context.Context, collectionID uuid.UUID, datapoints [][]byte, allowExisting bool) (*datasetsv1.IngestResponse, error) {
	return observability.WithSpanResult(ctx, s.tracer, "datasets/datapoints/ingest", func(ctx context.Context) (*datasetsv1.IngestResponse, error) {
		res, err := s.dataIngestionClient.Ingest(ctx, connect.NewRequest(
			datasetsv1.IngestRequest_builder{
				CollectionId:  uuidToProtobuf(collectionID),
				Values:        datapoints,
				AllowExisting: allowExisting,
			}.Build(),
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
			datasetsv1.DeleteRequest_builder{
				CollectionId: uuidToProtobuf(collectionID),
				DatapointIds: lo.Map(datapointIDs, func(datapointID uuid.UUID, _ int) *datasetsv1.ID {
					return uuidToProtobuf(datapointID)
				}),
			}.Build(),
		))
		if err != nil {
			return nil, fmt.Errorf("failed to delete datapoints: %w", err)
		}

		return res.Msg, nil
	})
}
