package workflows

import (
	"context"
	"encoding/json"
	"iter"

	"github.com/google/uuid"
	tileboxv1 "github.com/tilebox/tilebox-go/protogen/tilebox/v1"
	workflowsv1 "github.com/tilebox/tilebox-go/protogen/workflows/v1"
	"github.com/tilebox/tilebox-go/workflows/v1/storagelocation"
)

// StorageType is the kind of storage location used by automation triggers.
type StorageType string

const (
	StorageTypeUnspecified StorageType = "unspecified"
	StorageTypeGCS         StorageType = "gcs"
	StorageTypeS3          StorageType = "s3"
	StorageTypeFS          StorageType = "fs"
	StorageTypeAzureBlob   StorageType = "azure_blob"
)

func (t StorageType) String() string {
	return string(t)
}

func (t StorageType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

// StorageLocation identifies a bucket, container, or directory whose objects can trigger automations.
type StorageLocation struct {
	ID        uuid.UUID
	Name      string
	Reference *StorageLocationReference
}

// StorageLocationReference contains immutable provider coordinates.
// Set Type and the corresponding provider field when creating a location.
type StorageLocationReference struct {
	Type        StorageType
	AWSS3Bucket *AWSS3BucketReference
	GCSBucket   *GCSBucketReference
	AzureBlob   *AzureBlobReference
	Filesystem  *FilesystemReference
}

// AWSS3BucketReference identifies an S3 bucket and its region.
type AWSS3BucketReference struct {
	Bucket string
	Region string
}

// GCSBucketReference identifies a GCS bucket and its owning project.
type GCSBucketReference struct {
	Bucket    string
	ProjectID string
	// Location is an optional region, dual-region, or multi-region.
	Location string
}

// AzureBlobReference identifies a container in an Azure storage account.
type AzureBlobReference struct {
	// StorageAccountResourceID is the full ARM resource ID, not a URL or connection string.
	StorageAccountResourceID string
	Container                string
	// Region is the optional primary region of the storage account.
	Region string
}

// FilesystemReference identifies a directory monitored by a filesystem notifier.
type FilesystemReference struct {
	Path string
}

// StorageLocationClient manages storage locations and their notification subscriptions.
type StorageLocationClient interface {
	// Create registers an existing bucket, container, or directory; it does not create cloud resources.
	Create(ctx context.Context, name string, reference *StorageLocationReference) (*StorageLocation, error)
	// Get returns a storage location by ID.
	Get(ctx context.Context, storageLocationID uuid.UUID) (*StorageLocation, error)
	// List returns all storage locations.
	List(ctx context.Context) ([]*StorageLocation, error)
	// Update changes the display name. Provider coordinates are immutable.
	Update(ctx context.Context, storageLocationID uuid.UUID, name string) (*StorageLocation, error)
	// Delete removes the location and its subscriptions from Tilebox without deleting cloud resources.
	Delete(ctx context.Context, storageLocationID uuid.UUID) error
	// CreateSubscription registers notification delivery and returns provider setup values.
	// Configure the provider to send notifications to the returned endpoint afterwards.
	CreateSubscription(ctx context.Context, storageLocationID uuid.UUID, config StorageSubscriptionConfig) (*StorageSubscription, error)
	// GetSubscription returns a subscription by its ID.
	GetSubscription(ctx context.Context, subscriptionID uuid.UUID) (*StorageSubscription, error)
	// ListSubscriptions returns subscriptions for a storage location.
	ListSubscriptions(ctx context.Context, storageLocationID uuid.UUID) ([]*StorageSubscription, error)
	// DeleteSubscription stops accepting notifications without deleting provider resources.
	DeleteSubscription(ctx context.Context, subscriptionID uuid.UUID) error
	// ListSubscriptionEvents lazily iterates over a location's events, newest receipt first.
	ListSubscriptionEvents(ctx context.Context, storageLocationID uuid.UUID, options ...storagelocation.ListOption) iter.Seq2[*StorageSubscriptionEvent, error]
	// ListSubscriptionEventsPage returns one page of a location's events, newest receipt first.
	ListSubscriptionEventsPage(ctx context.Context, storageLocationID uuid.UUID, options ...storagelocation.ListOption) (*StorageSubscriptionEventPage, error)
}

var _ StorageLocationClient = &storageLocationClient{}

type storageLocationClient struct {
	service _storageLocationService
}

func (c storageLocationClient) Create(ctx context.Context, name string, reference *StorageLocationReference) (*StorageLocation, error) {
	response, err := c.service.CreateStorageLocation(ctx, workflowsv1.CreateStorageLocationRequest_builder{
		Name: name, Reference: reference.toProto(),
	}.Build())
	if err != nil {
		return nil, err
	}
	return protoToStorageLocation(response), nil
}

func (c storageLocationClient) Get(ctx context.Context, storageLocationID uuid.UUID) (*StorageLocation, error) {
	response, err := c.service.GetStorageLocation(ctx, storageLocationID)
	if err != nil {
		return nil, err
	}
	return protoToStorageLocation(response), nil
}

func (c storageLocationClient) List(ctx context.Context) ([]*StorageLocation, error) {
	response, err := c.service.ListStorageLocations(ctx)
	if err != nil {
		return nil, err
	}
	locations := make([]*StorageLocation, len(response.GetLocations()))
	for i, location := range response.GetLocations() {
		locations[i] = protoToStorageLocation(location)
	}
	return locations, nil
}

func (c storageLocationClient) Update(ctx context.Context, storageLocationID uuid.UUID, name string) (*StorageLocation, error) {
	response, err := c.service.UpdateStorageLocation(ctx, workflowsv1.UpdateStorageLocationRequest_builder{
		StorageLocationId: tileboxv1.NewUUID(storageLocationID), Name: name,
	}.Build())
	if err != nil {
		return nil, err
	}
	return protoToStorageLocation(response), nil
}

func (c storageLocationClient) Delete(ctx context.Context, storageLocationID uuid.UUID) error {
	return c.service.DeleteStorageLocation(ctx, storageLocationID)
}

func (c storageLocationClient) CreateSubscription(ctx context.Context, storageLocationID uuid.UUID, config StorageSubscriptionConfig) (*StorageSubscription, error) {
	response, err := c.service.CreateStorageSubscription(ctx, config.toProto(storageLocationID))
	if err != nil {
		return nil, err
	}
	return protoToStorageSubscription(response), nil
}

func (c storageLocationClient) GetSubscription(ctx context.Context, subscriptionID uuid.UUID) (*StorageSubscription, error) {
	response, err := c.service.GetStorageSubscription(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}
	return protoToStorageSubscription(response), nil
}

func (c storageLocationClient) ListSubscriptions(ctx context.Context, storageLocationID uuid.UUID) ([]*StorageSubscription, error) {
	response, err := c.service.ListStorageSubscriptions(ctx, storageLocationID)
	if err != nil {
		return nil, err
	}
	subscriptions := make([]*StorageSubscription, len(response.GetSubscriptions()))
	for i, subscription := range response.GetSubscriptions() {
		subscriptions[i] = protoToStorageSubscription(subscription)
	}
	return subscriptions, nil
}

func (c storageLocationClient) DeleteSubscription(ctx context.Context, subscriptionID uuid.UUID) error {
	return c.service.DeleteStorageSubscription(ctx, subscriptionID)
}

func (c storageLocationClient) ListSubscriptionEvents(ctx context.Context, storageLocationID uuid.UUID, options ...storagelocation.ListOption) iter.Seq2[*StorageSubscriptionEvent, error] {
	applied := storagelocation.NewListOptions(options...)
	return func(yield func(*StorageSubscriptionEvent, error) bool) {
		cursor := applied.Cursor
		remaining := applied.Limit
		for {
			page, err := c.ListSubscriptionEventsPage(ctx, storageLocationID, storagelocation.WithCursor(cursor), storagelocation.WithLimit(remaining))
			if err != nil {
				yield(nil, err)
				return
			}
			for _, event := range page.Events {
				if !yield(event, nil) {
					return
				}
				if applied.Limit > 0 {
					remaining--
					if remaining == 0 {
						return
					}
				}
			}
			cursor = page.NextCursor
			if cursor == nil {
				return
			}
		}
	}
}

func (c storageLocationClient) ListSubscriptionEventsPage(ctx context.Context, storageLocationID uuid.UUID, options ...storagelocation.ListOption) (*StorageSubscriptionEventPage, error) {
	applied := storagelocation.NewListOptions(options...)
	// The API accepts at most 100 events per page, even when the iterator's total limit is larger.
	limit := min(applied.Limit, 100)
	response, err := c.service.ListStorageSubscriptionEvents(ctx, storageLocationID, paginationFromOptions(limit, applied.Cursor))
	if err != nil {
		return nil, err
	}
	events := make([]*StorageSubscriptionEvent, len(response.GetEvents()))
	for i, event := range response.GetEvents() {
		events[i] = protoToStorageSubscriptionEvent(event)
	}
	return &StorageSubscriptionEventPage{Events: events, NextCursor: cursorFromPagination(response.GetNextPage())}, nil
}

func protoToStorageLocation(location *workflowsv1.StorageLocation) *StorageLocation {
	if location == nil {
		return nil
	}
	return &StorageLocation{
		ID: protoIDToUUID(location.GetId()), Name: location.GetName(),
		Reference: protoToStorageLocationReference(location.GetReference()),
	}
}

func protoToStorageLocationReference(reference *workflowsv1.StorageLocationReference) *StorageLocationReference {
	if reference == nil {
		return nil
	}
	result := &StorageLocationReference{Type: protoToStorageType(reference.GetType())}
	if bucket := reference.GetAwsS3Bucket(); bucket != nil {
		result.AWSS3Bucket = &AWSS3BucketReference{Bucket: bucket.GetBucket(), Region: bucket.GetRegion()}
	}
	if bucket := reference.GetGcsBucket(); bucket != nil {
		result.GCSBucket = &GCSBucketReference{Bucket: bucket.GetBucket(), ProjectID: bucket.GetProjectId(), Location: bucket.GetLocation()}
	}
	if blob := reference.GetAzureBlob(); blob != nil {
		result.AzureBlob = &AzureBlobReference{StorageAccountResourceID: blob.GetStorageAccountResourceId(), Container: blob.GetContainer(), Region: blob.GetRegion()}
	}
	if filesystem := reference.GetFilesystem(); filesystem != nil {
		result.Filesystem = &FilesystemReference{Path: filesystem.GetPath()}
	}
	return result
}

func (r *StorageLocationReference) toProto() *workflowsv1.StorageLocationReference {
	if r == nil {
		return nil
	}
	result := workflowsv1.StorageLocationReference_builder{Type: r.Type.toProto()}
	if r.AWSS3Bucket != nil {
		result.AwsS3Bucket = workflowsv1.AWSS3BucketReference_builder{Bucket: r.AWSS3Bucket.Bucket, Region: r.AWSS3Bucket.Region}.Build()
	}
	if r.GCSBucket != nil {
		result.GcsBucket = workflowsv1.GCSBucketReference_builder{Bucket: r.GCSBucket.Bucket, ProjectId: r.GCSBucket.ProjectID, Location: r.GCSBucket.Location}.Build()
	}
	if r.AzureBlob != nil {
		result.AzureBlob = workflowsv1.AzureBlobReference_builder{StorageAccountResourceId: r.AzureBlob.StorageAccountResourceID, Container: r.AzureBlob.Container, Region: r.AzureBlob.Region}.Build()
	}
	if r.Filesystem != nil {
		result.Filesystem = workflowsv1.FilesystemReference_builder{Path: r.Filesystem.Path}.Build()
	}
	return result.Build()
}

func protoToStorageType(storageType workflowsv1.StorageType) StorageType {
	switch storageType {
	case workflowsv1.StorageType_STORAGE_TYPE_UNSPECIFIED:
		return StorageTypeUnspecified
	case workflowsv1.StorageType_STORAGE_TYPE_GCS_BUCKET:
		return StorageTypeGCS
	case workflowsv1.StorageType_STORAGE_TYPE_AWS_S3_BUCKET:
		return StorageTypeS3
	case workflowsv1.StorageType_STORAGE_TYPE_FILESYSTEM:
		return StorageTypeFS
	case workflowsv1.StorageType_STORAGE_TYPE_AZURE_BLOB:
		return StorageTypeAzureBlob
	default:
		return StorageTypeUnspecified
	}
}

func (t StorageType) toProto() workflowsv1.StorageType {
	switch t {
	case StorageTypeGCS:
		return workflowsv1.StorageType_STORAGE_TYPE_GCS_BUCKET
	case StorageTypeS3:
		return workflowsv1.StorageType_STORAGE_TYPE_AWS_S3_BUCKET
	case StorageTypeFS:
		return workflowsv1.StorageType_STORAGE_TYPE_FILESYSTEM
	case StorageTypeAzureBlob:
		return workflowsv1.StorageType_STORAGE_TYPE_AZURE_BLOB
	case StorageTypeUnspecified:
		return workflowsv1.StorageType_STORAGE_TYPE_UNSPECIFIED
	default:
		return workflowsv1.StorageType_STORAGE_TYPE_UNSPECIFIED
	}
}
