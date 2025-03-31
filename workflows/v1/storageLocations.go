package workflows // import "github.com/tilebox/tilebox-go/workflows/v1"

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	workflowsv1 "github.com/tilebox/tilebox-go/protogen/go/workflows/v1"
)

// StorageLocation represents a kind of storage that can contain data files or objects and be used as a trigger source.
//
// Documentation: https://docs.tilebox.com/workflows/near-real-time/storage-events#storage-locations
type StorageLocation struct {
	// ID is the unique identifier of the storage location.
	ID uuid.UUID
	// Location is the unique identifier of the storage location in the storage system.
	Location string
	// Type is the type of the storage location, e.g. GCS, S3, FS.
	Type StorageType
}

// StorageType is a type of storage location.
type StorageType int32

// StorageType values.
const (
	_              StorageType = iota
	StorageTypeGCS             // Google Cloud Storage
	StorageTypeS3              // Amazon Web Services S3
	StorageTypeFS              // Local filesystem
)

type StorageLocationClient interface {
	Create(ctx context.Context, location string, storageType StorageType) (*StorageLocation, error)
	Get(ctx context.Context, storageLocationID uuid.UUID) (*StorageLocation, error)
	Delete(ctx context.Context, storageLocationID uuid.UUID) error
	List(ctx context.Context) ([]*StorageLocation, error)
}

var _ StorageLocationClient = &storageLocationClient{}

type storageLocationClient struct {
	service AutomationService
}

func (d storageLocationClient) Create(ctx context.Context, location string, storageType StorageType) (*StorageLocation, error) {
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

func (d storageLocationClient) Get(ctx context.Context, storageLocationID uuid.UUID) (*StorageLocation, error) {
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

func (d storageLocationClient) List(ctx context.Context) ([]*StorageLocation, error) {
	response, err := d.service.ListStorageLocations(ctx)
	if err != nil {
		return nil, err
	}

	storageLocations := make([]*StorageLocation, len(response.GetLocations()))
	for i, s := range response.GetLocations() {
		storageLocation, err := protoToStorageLocation(s)
		if err != nil {
			return nil, fmt.Errorf("failed to convert storageLocation from response: %w", err)
		}
		storageLocations[i] = storageLocation
	}

	return storageLocations, nil
}

func protoToStorageLocation(s *workflowsv1.StorageLocation) (*StorageLocation, error) {
	id, err := protoToUUID(s.GetId())
	if err != nil {
		return nil, fmt.Errorf("failed to parse storage location id: %w", err)
	}

	return &StorageLocation{
		ID:       id,
		Location: s.GetLocation(),
		Type:     StorageType(s.GetType()),
	}, nil
}
