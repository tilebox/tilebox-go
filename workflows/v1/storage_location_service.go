package workflows

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/tilebox/tilebox-go/observability"
	tileboxv1 "github.com/tilebox/tilebox-go/protogen/tilebox/v1"
	workflowsv1 "github.com/tilebox/tilebox-go/protogen/workflows/v1"
	"github.com/tilebox/tilebox-go/protogen/workflows/v1/workflowsv1connect"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/types/known/emptypb"
)

type _storageLocationService interface {
	CreateStorageLocation(context.Context, *workflowsv1.CreateStorageLocationRequest) (*workflowsv1.StorageLocation, error)
	GetStorageLocation(context.Context, uuid.UUID) (*workflowsv1.StorageLocation, error)
	ListStorageLocations(context.Context) (*workflowsv1.StorageLocations, error)
	UpdateStorageLocation(context.Context, *workflowsv1.UpdateStorageLocationRequest) (*workflowsv1.StorageLocation, error)
	DeleteStorageLocation(context.Context, uuid.UUID) error
	CreateStorageSubscription(context.Context, *workflowsv1.CreateStorageSubscriptionRequest) (*workflowsv1.StorageSubscription, error)
	GetStorageSubscription(context.Context, uuid.UUID) (*workflowsv1.StorageSubscription, error)
	ListStorageSubscriptions(context.Context, uuid.UUID) (*workflowsv1.StorageSubscriptions, error)
	DeleteStorageSubscription(context.Context, uuid.UUID) error
	ListStorageSubscriptionEvents(context.Context, uuid.UUID, *tileboxv1.Pagination) (*workflowsv1.StorageSubscriptionEvents, error)
}

var _ _storageLocationService = &storageLocationService{}

type storageLocationService struct {
	client workflowsv1connect.StorageLocationServiceClient
	tracer trace.Tracer
}

func (s *storageLocationService) CreateStorageLocation(ctx context.Context, request *workflowsv1.CreateStorageLocationRequest) (*workflowsv1.StorageLocation, error) {
	return observability.WithSpanResult(ctx, s.tracer, "workflows/storage_locations/create", func(ctx context.Context) (*workflowsv1.StorageLocation, error) {
		response, err := s.client.CreateStorageLocation(ctx, connect.NewRequest(request))
		if err != nil {
			return nil, fmt.Errorf("failed to create storage location: %w", err)
		}
		return response.Msg, nil
	})
}

func (s *storageLocationService) GetStorageLocation(ctx context.Context, storageLocationID uuid.UUID) (*workflowsv1.StorageLocation, error) {
	return observability.WithSpanResult(ctx, s.tracer, "workflows/storage_locations/get", func(ctx context.Context) (*workflowsv1.StorageLocation, error) {
		response, err := s.client.GetStorageLocation(ctx, connect.NewRequest(tileboxv1.NewUUID(storageLocationID)))
		if err != nil {
			return nil, fmt.Errorf("failed to get storage location: %w", err)
		}
		return response.Msg, nil
	})
}

func (s *storageLocationService) ListStorageLocations(ctx context.Context) (*workflowsv1.StorageLocations, error) {
	return observability.WithSpanResult(ctx, s.tracer, "workflows/storage_locations/list", func(ctx context.Context) (*workflowsv1.StorageLocations, error) {
		response, err := s.client.ListStorageLocations(ctx, connect.NewRequest(&emptypb.Empty{}))
		if err != nil {
			return nil, fmt.Errorf("failed to list storage locations: %w", err)
		}
		return response.Msg, nil
	})
}

func (s *storageLocationService) UpdateStorageLocation(ctx context.Context, request *workflowsv1.UpdateStorageLocationRequest) (*workflowsv1.StorageLocation, error) {
	return observability.WithSpanResult(ctx, s.tracer, "workflows/storage_locations/update", func(ctx context.Context) (*workflowsv1.StorageLocation, error) {
		response, err := s.client.UpdateStorageLocation(ctx, connect.NewRequest(request))
		if err != nil {
			return nil, fmt.Errorf("failed to update storage location: %w", err)
		}
		return response.Msg, nil
	})
}

func (s *storageLocationService) DeleteStorageLocation(ctx context.Context, storageLocationID uuid.UUID) error {
	return observability.WithSpan(ctx, s.tracer, "workflows/storage_locations/delete", func(ctx context.Context) error {
		_, err := s.client.DeleteStorageLocation(ctx, connect.NewRequest(tileboxv1.NewUUID(storageLocationID)))
		if err != nil {
			return fmt.Errorf("failed to delete storage location: %w", err)
		}
		return nil
	})
}

func (s *storageLocationService) CreateStorageSubscription(ctx context.Context, request *workflowsv1.CreateStorageSubscriptionRequest) (*workflowsv1.StorageSubscription, error) {
	return observability.WithSpanResult(ctx, s.tracer, "workflows/storage_subscriptions/create", func(ctx context.Context) (*workflowsv1.StorageSubscription, error) {
		response, err := s.client.CreateStorageSubscription(ctx, connect.NewRequest(request))
		if err != nil {
			return nil, fmt.Errorf("failed to create storage subscription: %w", err)
		}
		return response.Msg, nil
	})
}

func (s *storageLocationService) GetStorageSubscription(ctx context.Context, subscriptionID uuid.UUID) (*workflowsv1.StorageSubscription, error) {
	return observability.WithSpanResult(ctx, s.tracer, "workflows/storage_subscriptions/get", func(ctx context.Context) (*workflowsv1.StorageSubscription, error) {
		response, err := s.client.GetStorageSubscription(ctx, connect.NewRequest(tileboxv1.NewUUID(subscriptionID)))
		if err != nil {
			return nil, fmt.Errorf("failed to get storage subscription: %w", err)
		}
		return response.Msg, nil
	})
}

func (s *storageLocationService) ListStorageSubscriptions(ctx context.Context, storageLocationID uuid.UUID) (*workflowsv1.StorageSubscriptions, error) {
	return observability.WithSpanResult(ctx, s.tracer, "workflows/storage_subscriptions/list", func(ctx context.Context) (*workflowsv1.StorageSubscriptions, error) {
		response, err := s.client.ListStorageSubscriptions(ctx, connect.NewRequest(workflowsv1.ListStorageSubscriptionsRequest_builder{
			StorageLocationId: tileboxv1.NewUUID(storageLocationID),
		}.Build()))
		if err != nil {
			return nil, fmt.Errorf("failed to list storage subscriptions: %w", err)
		}
		return response.Msg, nil
	})
}

func (s *storageLocationService) DeleteStorageSubscription(ctx context.Context, subscriptionID uuid.UUID) error {
	return observability.WithSpan(ctx, s.tracer, "workflows/storage_subscriptions/delete", func(ctx context.Context) error {
		_, err := s.client.DeleteStorageSubscription(ctx, connect.NewRequest(tileboxv1.NewUUID(subscriptionID)))
		if err != nil {
			return fmt.Errorf("failed to delete storage subscription: %w", err)
		}
		return nil
	})
}

func (s *storageLocationService) ListStorageSubscriptionEvents(ctx context.Context, storageLocationID uuid.UUID, page *tileboxv1.Pagination) (*workflowsv1.StorageSubscriptionEvents, error) {
	return observability.WithSpanResult(ctx, s.tracer, "workflows/storage_subscriptions/list_events", func(ctx context.Context) (*workflowsv1.StorageSubscriptionEvents, error) {
		response, err := s.client.ListStorageSubscriptionEvents(ctx, connect.NewRequest(workflowsv1.ListStorageSubscriptionEventsRequest_builder{
			StorageLocationId: tileboxv1.NewUUID(storageLocationID), Page: page,
		}.Build()))
		if err != nil {
			return nil, fmt.Errorf("failed to list storage subscription events: %w", err)
		}
		return response.Msg, nil
	})
}
