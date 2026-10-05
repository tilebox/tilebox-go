package workflows

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tileboxv1 "github.com/tilebox/tilebox-go/protogen/tilebox/v1"
	workflowsv1 "github.com/tilebox/tilebox-go/protogen/workflows/v1"
)

func TestAutomationStorageLocationReference(t *testing.T) {
	storageLocationID := uuid.MustParse("019e4f3c-4646-7312-b8fe-2e7fa83c1546")
	automation := protoToAutomation(workflowsv1.AutomationPrototype_builder{
		StorageEventTriggers: []*workflowsv1.StorageEventTrigger{
			workflowsv1.StorageEventTrigger_builder{
				GlobPattern: "data/*.tif",
				StorageLocation: workflowsv1.StorageLocation_builder{
					Id: tileboxv1.NewUUID(storageLocationID), Name: "Imagery",
					Reference: workflowsv1.StorageLocationReference_builder{
						Type:      workflowsv1.StorageType_STORAGE_TYPE_GCS_BUCKET,
						GcsBucket: workflowsv1.GCSBucketReference_builder{Bucket: "bucket", ProjectId: "project", Location: "EU"}.Build(),
					}.Build(),
				}.Build(),
			}.Build(),
		},
	}.Build())
	require.Len(t, automation.StorageEventTriggers, 1)
	assert.Equal(t, "data/*.tif", automation.StorageEventTriggers[0].GlobPattern)
	location := automation.StorageEventTriggers[0].StorageLocation
	assert.Equal(t, storageLocationID, location.ID)
	assert.Equal(t, "Imagery", location.Name)
	assert.Equal(t, &StorageLocationReference{
		Type: StorageTypeGCS, GCSBucket: &GCSBucketReference{Bucket: "bucket", ProjectID: "project", Location: "EU"},
	}, location.Reference)
}
