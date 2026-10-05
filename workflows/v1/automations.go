package workflows

import (
	"context"

	"github.com/google/uuid"
	tileboxv1 "github.com/tilebox/tilebox-go/protogen/tilebox/v1"
	workflowsv1 "github.com/tilebox/tilebox-go/protogen/workflows/v1"
)

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

func protoIDToUUID(id *tileboxv1.ID) uuid.UUID {
	if id == nil {
		return uuid.Nil
	}
	return id.AsUUID()
}
