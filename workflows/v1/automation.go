package workflows

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	workflowsv1 "github.com/tilebox/tilebox-go/protogen/go/workflows/v1"
	"github.com/tilebox/tilebox-go/protogen/go/workflows/v1/workflowsv1connect"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AutomationService struct {
	client workflowsv1connect.AutomationServiceClient
}

func NewAutomationService(client workflowsv1connect.AutomationServiceClient) *AutomationService {
	return &AutomationService{
		client: client,
	}
}

func (rs *AutomationService) ListStorageLocations(ctx context.Context) ([]*workflowsv1.StorageLocation, error) {
	response, err := rs.client.ListStorageLocations(ctx, connect.NewRequest(&emptypb.Empty{}))

	if err != nil {
		return nil, fmt.Errorf("failed to list storage locations: %w", err)
	}

	return response.Msg.GetLocations(), nil
}

func (rs *AutomationService) GetStorageLocation(ctx context.Context, storageLocationID *workflowsv1.UUID) (*workflowsv1.StorageLocation, error) {
	response, err := rs.client.GetStorageLocation(ctx, connect.NewRequest(storageLocationID))

	if err != nil {
		return nil, fmt.Errorf("failed to get storage location: %w", err)
	}

	return response.Msg, nil
}

func (rs *AutomationService) CreateStorageLocation(ctx context.Context, location string, storageType workflowsv1.StorageType) (*workflowsv1.StorageLocation, error) {
	req := connect.NewRequest(&workflowsv1.StorageLocation{
		Location: location,
		Type:     storageType,
	})
	response, err := rs.client.CreateStorageLocation(ctx, req)

	if err != nil {
		return nil, fmt.Errorf("failed to create storage location: %w", err)
	}

	return response.Msg, nil
}

func (rs *AutomationService) DeleteStorageLocation(ctx context.Context, storageLocationID *workflowsv1.UUID) error {
	_, err := rs.client.DeleteStorageLocation(ctx, connect.NewRequest(storageLocationID))

	if err != nil {
		return fmt.Errorf("failed to delete storage location: %w", err)
	}

	return nil
}
