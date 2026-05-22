package workflows

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tileboxv1 "github.com/tilebox/tilebox-go/protogen/tilebox/v1"
	workflowsv1 "github.com/tilebox/tilebox-go/protogen/workflows/v1"
)

type mockAutomationService struct {
	_automationService

	getStorageLocationID uuid.UUID
	storageLocation      *workflowsv1.StorageLocation
}

func (m *mockAutomationService) GetStorageLocation(_ context.Context, storageLocationID uuid.UUID) (*workflowsv1.StorageLocation, error) {
	m.getStorageLocationID = storageLocationID
	return m.storageLocation, nil
}

func TestAutomationClientGetStorageLocation(t *testing.T) {
	storageLocationID := uuid.MustParse("019e4f3c-4646-7312-b8fe-2e7fa83c1546")
	service := &mockAutomationService{
		storageLocation: workflowsv1.StorageLocation_builder{
			Id:       tileboxv1.NewUUID(storageLocationID),
			Location: "gs://bucket",
			Type:     workflowsv1.StorageType_STORAGE_TYPE_GCS,
		}.Build(),
	}
	client := &automationClient{service: service}

	location, err := client.GetStorageLocation(context.Background(), storageLocationID)

	require.NoError(t, err)
	require.NotNil(t, location)
	assert.Equal(t, storageLocationID, service.getStorageLocationID)
	assert.Equal(t, storageLocationID, location.ID)
	assert.Equal(t, "gs://bucket", location.Location)
	assert.Equal(t, StorageTypeGCS, location.Type)
}
