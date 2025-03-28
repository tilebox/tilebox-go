package workflows // import "github.com/tilebox/tilebox-go/workflows/v1"

import (
	"context"

	"github.com/google/uuid"
	workflowsv1 "github.com/tilebox/tilebox-go/protogen/go/workflows/v1"
)

type _automationClient interface {
	Create(ctx context.Context, automation *workflowsv1.AutomationPrototype) (*workflowsv1.AutomationPrototype, error)
	Get(ctx context.Context, automationID uuid.UUID) (*workflowsv1.AutomationPrototype, error)
	Update(ctx context.Context, automation *workflowsv1.AutomationPrototype) (*workflowsv1.AutomationPrototype, error)
	Delete(ctx context.Context, automationID uuid.UUID, cancelJobs bool) error
	List(ctx context.Context) ([]*workflowsv1.AutomationPrototype, error)
}

var _ _automationClient = &automationClient{}

type automationClient struct {
	service AutomationService
}

func (d automationClient) Create(ctx context.Context, automation *workflowsv1.AutomationPrototype) (*workflowsv1.AutomationPrototype, error) {
	response, err := d.service.CreateAutomation(ctx, automation)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (d automationClient) Get(ctx context.Context, automationID uuid.UUID) (*workflowsv1.AutomationPrototype, error) {
	response, err := d.service.GetAutomation(ctx, automationID)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (d automationClient) Update(ctx context.Context, automation *workflowsv1.AutomationPrototype) (*workflowsv1.AutomationPrototype, error) {
	response, err := d.service.UpdateAutomation(ctx, automation)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (d automationClient) Delete(ctx context.Context, automationID uuid.UUID, cancelJobs bool) error {
	return d.service.DeleteAutomation(ctx, automationID, cancelJobs)
}

func (d automationClient) List(ctx context.Context) ([]*workflowsv1.AutomationPrototype, error) {
	response, err := d.service.ListAutomations(ctx)
	if err != nil {
		return nil, err
	}

	return response.GetAutomations(), nil
}
