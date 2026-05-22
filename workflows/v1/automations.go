package workflows // import "github.com/tilebox/tilebox-go/workflows/v1"

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	tileboxv1 "github.com/tilebox/tilebox-go/protogen/tilebox/v1"
	workflowsv1 "github.com/tilebox/tilebox-go/protogen/workflows/v1"
)

// StorageType is the kind of storage location used by automation triggers.
type StorageType string

const (
	StorageTypeUnspecified StorageType = "unspecified"
	StorageTypeGCS         StorageType = "gcs"
	StorageTypeS3          StorageType = "s3"
	StorageTypeFS          StorageType = "fs"
)

func (t StorageType) String() string {
	return string(t)
}

func (t StorageType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

// StorageLocation is a storage location available for automation storage event triggers.
type StorageLocation struct {
	// ID is the unique identifier of the storage location.
	ID uuid.UUID
	// Location is the storage-system-specific location identifier.
	Location string
	// Type is the kind of storage location.
	Type StorageType
}

// Automation represents an automation prototype that can submit tasks from storage or cron triggers.
type Automation struct {
	// ID is the unique identifier of the automation.
	ID uuid.UUID
	// Name is the human-readable name of the automation.
	Name string
	// Prototype is the task submission prototype that the automation submits.
	Prototype *AutomationTaskPrototype
	// StorageEventTriggers are triggers that submit the task for matching storage events.
	StorageEventTriggers []*StorageEventTrigger
	// CronTriggers are triggers that submit the task on a schedule.
	CronTriggers []*CronTrigger
	// Disabled reports whether the automation is paused.
	Disabled bool
}

// AutomationTaskPrototype is the single task submitted by an automation.
type AutomationTaskPrototype struct {
	// ClusterSlug is the cluster where the task should run.
	ClusterSlug string
	// Identifier identifies the task implementation.
	Identifier TaskIdentifier
	// Display is a human-readable task label.
	Display string
	// Dependencies are task dependency indexes.
	Dependencies []int64
	// MaxRetries is the maximum number of automatic retries for the task.
	MaxRetries int64
	// Input is the serialized task input.
	Input []byte
}

// StorageEventTrigger submits an automation task when a matching object is created in a storage location.
type StorageEventTrigger struct {
	// ID is the unique identifier of the trigger.
	ID uuid.UUID
	// StorageLocation is the storage location watched by this trigger.
	StorageLocation *StorageLocation
	// GlobPattern matches objects/files in the storage location.
	GlobPattern string
}

// CronTrigger submits an automation task on a cron schedule.
type CronTrigger struct {
	// ID is the unique identifier of the trigger.
	ID uuid.UUID
	// Schedule is the cron schedule for the trigger.
	Schedule string
}

type AutomationClient interface {
	// List returns all automations.
	List(ctx context.Context) ([]*Automation, error)

	// Get returns an automation by ID.
	Get(ctx context.Context, automationID uuid.UUID) (*Automation, error)

	// GetStorageLocation returns a storage location by ID.
	GetStorageLocation(ctx context.Context, storageLocationID uuid.UUID) (*StorageLocation, error)

	// ListStorageLocations returns all storage locations available for automation triggers.
	ListStorageLocations(ctx context.Context) ([]*StorageLocation, error)
}

var _ AutomationClient = &automationClient{}

type automationClient struct {
	service _automationService
}

func (c automationClient) List(ctx context.Context) ([]*Automation, error) {
	response, err := c.service.ListAutomations(ctx)
	if err != nil {
		return nil, err
	}

	automations := make([]*Automation, len(response.GetAutomations()))
	for i, automation := range response.GetAutomations() {
		automations[i] = protoToAutomation(automation)
	}

	return automations, nil
}

func (c automationClient) Get(ctx context.Context, automationID uuid.UUID) (*Automation, error) {
	response, err := c.service.GetAutomation(ctx, automationID)
	if err != nil {
		return nil, err
	}

	return protoToAutomation(response), nil
}

func (c automationClient) GetStorageLocation(ctx context.Context, storageLocationID uuid.UUID) (*StorageLocation, error) {
	response, err := c.service.GetStorageLocation(ctx, storageLocationID)
	if err != nil {
		return nil, err
	}

	return protoToStorageLocation(response), nil
}

func (c automationClient) ListStorageLocations(ctx context.Context) ([]*StorageLocation, error) {
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

func protoToAutomation(automation *workflowsv1.AutomationPrototype) *Automation {
	if automation == nil {
		return nil
	}

	storageEventTriggers := make([]*StorageEventTrigger, len(automation.GetStorageEventTriggers()))
	for i, trigger := range automation.GetStorageEventTriggers() {
		storageEventTriggers[i] = protoToStorageEventTrigger(trigger)
	}

	cronTriggers := make([]*CronTrigger, len(automation.GetCronTriggers()))
	for i, trigger := range automation.GetCronTriggers() {
		cronTriggers[i] = protoToCronTrigger(trigger)
	}

	return &Automation{
		ID:                   protoIDToUUID(automation.GetId()),
		Name:                 automation.GetName(),
		Prototype:            protoToAutomationTaskPrototype(automation.GetPrototype()),
		StorageEventTriggers: storageEventTriggers,
		CronTriggers:         cronTriggers,
		Disabled:             automation.GetDisabled(),
	}
}

func protoToAutomationTaskPrototype(prototype *workflowsv1.SingleTaskSubmission) *AutomationTaskPrototype {
	if prototype == nil {
		return nil
	}

	var identifier TaskIdentifier
	if protoIdentifier := prototype.GetIdentifier(); protoIdentifier != nil {
		identifier = NewTaskIdentifier(protoIdentifier.GetName(), protoIdentifier.GetVersion())
	}

	return &AutomationTaskPrototype{
		ClusterSlug:  prototype.GetClusterSlug(),
		Identifier:   identifier,
		Display:      prototype.GetDisplay(),
		Dependencies: prototype.GetDependencies(),
		MaxRetries:   prototype.GetMaxRetries(),
		Input:        prototype.GetInput(),
	}
}

func protoToStorageEventTrigger(trigger *workflowsv1.StorageEventTrigger) *StorageEventTrigger {
	if trigger == nil {
		return nil
	}
	return &StorageEventTrigger{
		ID:              protoIDToUUID(trigger.GetId()),
		StorageLocation: protoToStorageLocation(trigger.GetStorageLocation()),
		GlobPattern:     trigger.GetGlobPattern(),
	}
}

func protoToCronTrigger(trigger *workflowsv1.CronTrigger) *CronTrigger {
	if trigger == nil {
		return nil
	}
	return &CronTrigger{
		ID:       protoIDToUUID(trigger.GetId()),
		Schedule: trigger.GetSchedule(),
	}
}

func protoToStorageLocation(location *workflowsv1.StorageLocation) *StorageLocation {
	if location == nil {
		return nil
	}
	return &StorageLocation{
		ID:       protoIDToUUID(location.GetId()),
		Location: location.GetLocation(),
		Type:     protoToStorageType(location.GetType()),
	}
}

func protoToStorageType(storageType workflowsv1.StorageType) StorageType {
	switch storageType {
	case workflowsv1.StorageType_STORAGE_TYPE_UNSPECIFIED:
		return StorageTypeUnspecified
	case workflowsv1.StorageType_STORAGE_TYPE_GCS:
		return StorageTypeGCS
	case workflowsv1.StorageType_STORAGE_TYPE_S3:
		return StorageTypeS3
	case workflowsv1.StorageType_STORAGE_TYPE_FS:
		return StorageTypeFS
	default:
		return StorageTypeUnspecified
	}
}

func protoIDToUUID(id *tileboxv1.ID) uuid.UUID {
	if id == nil {
		return uuid.Nil
	}
	return id.AsUUID()
}
