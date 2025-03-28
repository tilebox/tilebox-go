package workflows // import "github.com/tilebox/tilebox-go/workflows/v1"

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	workflowsv1 "github.com/tilebox/tilebox-go/protogen/go/workflows/v1"
)

// _storageLocation represents a kind of storage that can contain data files or objects and be used as a trigger source.
//
// Documentation: https://docs.tilebox.com/workflows/near-real-time/storage-events#storage-locations
type _storageLocation struct {
	// ID is the unique identifier of the storage location.
	ID uuid.UUID
	// Location is the unique identifier of the storage location in the storage system.
	Location string
	// Type is the type of the storage location, e.g. GCS, S3, FS.
	Type _storageType
}

// _storageType is a type of storage location.
type _storageType int32

type _storageLocationClient interface {
	Create(ctx context.Context, location string, storageType _storageType) (*_storageLocation, error)
	Get(ctx context.Context, storageLocationID uuid.UUID) (*_storageLocation, error)
	Delete(ctx context.Context, storageLocationID uuid.UUID) error
	List(ctx context.Context) ([]*_storageLocation, error)
}

var _ _storageLocationClient = &storageLocationClient{}

type storageLocationClient struct {
	service AutomationService
}

func (d storageLocationClient) Create(ctx context.Context, location string, storageType _storageType) (*_storageLocation, error) {
	response, err := d.service.CreateStorageLocation(ctx, location, workflowsv1.StorageType(storageType))
	if err != nil {
		return nil, err
	}

	storageLocation, err := protoToStorageLocation(response)
	if err != nil {
		return nil, fmt.Errorf("failed to convert storageLocation from response: %w", err)
	}

	return storageLocation, nil
}

func (d storageLocationClient) Get(ctx context.Context, storageLocationID uuid.UUID) (*_storageLocation, error) {
	response, err := d.service.GetStorageLocation(ctx, storageLocationID)
	if err != nil {
		return nil, err
	}

	storageLocation, err := protoToStorageLocation(response)
	if err != nil {
		return nil, fmt.Errorf("failed to convert storageLocation from response: %w", err)
	}

	return storageLocation, nil
}

func (d storageLocationClient) Delete(ctx context.Context, storageLocationID uuid.UUID) error {
	return d.service.DeleteStorageLocation(ctx, storageLocationID)
}

func (d storageLocationClient) List(ctx context.Context) ([]*_storageLocation, error) {
	response, err := d.service.ListStorageLocations(ctx)
	if err != nil {
		return nil, err
	}

	storageLocations := make([]*_storageLocation, len(response.GetLocations()))
	for i, s := range response.GetLocations() {
		storageLocation, err := protoToStorageLocation(s)
		if err != nil {
			return nil, fmt.Errorf("failed to convert storageLocation from response: %w", err)
		}
		storageLocations[i] = storageLocation
	}

	return storageLocations, nil
}

func protoToStorageLocation(s *workflowsv1.StorageLocation) (*_storageLocation, error) {
	id, err := protoToUUID(s.GetId())
	if err != nil {
		return nil, fmt.Errorf("failed to parse storage location id: %w", err)
	}

	return &_storageLocation{
		ID:       id,
		Location: s.GetLocation(),
		Type:     _storageType(s.GetType()),
	}, nil
}
