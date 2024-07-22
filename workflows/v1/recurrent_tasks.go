package workflows

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	workflowsv1 "github.com/tilebox/tilebox-go/protogen/go/workflows/v1"
	"github.com/tilebox/tilebox-go/protogen/go/workflows/v1/workflowsv1connect"
	"google.golang.org/protobuf/types/known/emptypb"
)

type RecurrentTaskService struct {
	client workflowsv1connect.RecurrentTaskServiceClient
}

func NewRecurrentTaskService(client workflowsv1connect.RecurrentTaskServiceClient) *RecurrentTaskService {
	return &RecurrentTaskService{
		client: client,
	}
}

func (rs *RecurrentTaskService) ListStorageLocations(ctx context.Context) ([]*workflowsv1.StorageLocation, error) {
	response, err := rs.client.ListStorageLocations(ctx, connect.NewRequest(&emptypb.Empty{}))

	if err != nil {
		return nil, fmt.Errorf("failed to list storage locations: %w", err)
	}

	return response.Msg.GetLocations(), nil
}

func (rs *RecurrentTaskService) GetStorageLocation(ctx context.Context, storageLocationID *workflowsv1.UUID) (*workflowsv1.StorageLocation, error) {
	response, err := rs.client.GetStorageLocation(ctx, connect.NewRequest(storageLocationID))

	if err != nil {
		return nil, fmt.Errorf("failed to get storage location: %w", err)
	}

	return response.Msg, nil
}

func (rs *RecurrentTaskService) CreateStorageLocation(ctx context.Context, location string, storageType workflowsv1.StorageType) (*workflowsv1.StorageLocation, error) {
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

func (rs *RecurrentTaskService) DeleteStorageLocation(ctx context.Context, storageLocationID *workflowsv1.UUID) error {
	_, err := rs.client.DeleteStorageLocation(ctx, connect.NewRequest(storageLocationID))

	if err != nil {
		return fmt.Errorf("failed to delete storage location: %w", err)
	}

	return nil
}
