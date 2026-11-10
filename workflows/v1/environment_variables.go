package workflows

import (
	"context"
	"fmt"
	"iter"

	"github.com/google/uuid"
	"github.com/tilebox/tilebox-go/observability"
	tileboxv1 "github.com/tilebox/tilebox-go/protogen/tilebox/v1"
	workflowsv1 "github.com/tilebox/tilebox-go/protogen/workflows/v1"
	"github.com/tilebox/tilebox-go/protogen/workflows/v1/workflowsv1connect"
	"github.com/tilebox/tilebox-go/query"
	"go.opentelemetry.io/otel/trace"
)

// ClusterScope selects all current and future clusters or a set of cluster slugs.
// AllClusters and a nonempty ClusterSlugs selection are mutually exclusive.
// Create and update require one of these selections. Reads may return an empty
// selection after its clusters are deleted; it applies nowhere, never everywhere.
// A selection containing the empty slug selects the default cluster.
type ClusterScope struct {
	AllClusters  bool
	ClusterSlugs []string
}

// EnvironmentVariable is one shared value with a stable ID and cluster scope.
// Never log Value or SecretPreview. A preview is not a usable value.
type EnvironmentVariable struct {
	ID            uuid.UUID
	Name          string
	IsSecret      bool
	Scope         ClusterScope
	Value         *string
	SecretPreview *string
	Description   string
}

// CreateEnvironmentVariableParams creates an entry without replacing existing entries.
// Name must be uppercase, Value must be nonempty, and Scope must explicitly select clusters.
type CreateEnvironmentVariableParams struct {
	Name        string
	Value       string
	IsSecret    bool
	Scope       ClusterScope
	Description string
}

// UpdateEnvironmentVariableParams changes only supplied fields. Scope replaces the
// entire selection; false declassifies a secret, and an empty Description clears it.
// A supplied Scope must select all clusters or at least one cluster slug.
// A supplied Name must be uppercase and Value must be nonempty. At least one field must be supplied.
type UpdateEnvironmentVariableParams struct {
	Name        *string
	Value       *string
	IsSecret    *bool
	Scope       *ClusterScope
	Description *string
}

// ListEnvironmentVariablesParams configures inventory or effective-cluster reads.
type ListEnvironmentVariablesParams struct {
	// Nil lists inventory; a pointer to "" selects the default cluster.
	ClusterSlug    *string
	IncludeSecrets bool
	// Limit is a page size for ListPage and a total limit for List. Zero uses defaults.
	Limit  int64
	Cursor *query.Cursor
}

// EnvironmentVariablePage is one page of variables, newest first.
type EnvironmentVariablePage struct {
	EnvironmentVariables []*EnvironmentVariable
	NextCursor           *query.Cursor
}

// EnvironmentVariableClient manages variables in the authenticated namespace.
type EnvironmentVariableClient interface {
	// Create and Update return metadata only, without values or previews.
	Create(context.Context, CreateEnvironmentVariableParams) (*EnvironmentVariable, error)
	Get(ctx context.Context, id uuid.UUID, includeSecret bool) (*EnvironmentVariable, error)
	List(context.Context, ListEnvironmentVariablesParams) iter.Seq2[*EnvironmentVariable, error]
	ListPage(context.Context, ListEnvironmentVariablesParams) (*EnvironmentVariablePage, error)
	Update(context.Context, uuid.UUID, UpdateEnvironmentVariableParams) (*EnvironmentVariable, error)
	Delete(context.Context, uuid.UUID) error
}

type environmentVariableClient struct {
	service workflowsv1connect.ResourceServiceClient
	tracer  trace.Tracer
}

var _ EnvironmentVariableClient = (*environmentVariableClient)(nil)

func (c *environmentVariableClient) Create(ctx context.Context, params CreateEnvironmentVariableParams) (*EnvironmentVariable, error) {
	return observability.WithSpanResult(ctx, c.tracer, "workflows/variables/create", func(ctx context.Context) (*EnvironmentVariable, error) {
		response, err := c.service.CreateEnvironmentVariable(ctx, workflowsv1.CreateEnvironmentVariableRequest_builder{
			Name: params.Name, Value: params.Value, IsSecret: params.IsSecret, Scope: params.Scope.toProto(), Description: params.Description,
		}.Build())
		if err != nil {
			return nil, fmt.Errorf("failed to create environment variable: %w", err)
		}
		return protoToEnvironmentVariable(response), nil
	})
}

func (c *environmentVariableClient) Get(ctx context.Context, id uuid.UUID, includeSecret bool) (*EnvironmentVariable, error) {
	return observability.WithSpanResult(ctx, c.tracer, "workflows/variables/get", func(ctx context.Context) (*EnvironmentVariable, error) {
		response, err := c.service.GetEnvironmentVariable(ctx, workflowsv1.GetEnvironmentVariableRequest_builder{
			Id: tileboxv1.NewUUID(id), IncludeSecret: includeSecret,
		}.Build())
		if err != nil {
			return nil, fmt.Errorf("failed to get environment variable: %w", err)
		}
		return protoToEnvironmentVariable(response), nil
	})
}

func (c *environmentVariableClient) List(ctx context.Context, options ListEnvironmentVariablesParams) iter.Seq2[*EnvironmentVariable, error] {
	return func(yield func(*EnvironmentVariable, error) bool) {
		params := options
		remaining := params.Limit
		for {
			pageParams := params
			pageParams.Limit = min(remaining, 100)
			page, err := c.ListPage(ctx, pageParams)
			if err != nil {
				yield(nil, err)
				return
			}
			for _, variable := range page.EnvironmentVariables {
				if !yield(variable, nil) {
					return
				}
				if params.Limit > 0 {
					remaining--
					if remaining == 0 {
						return
					}
				}
			}
			params.Cursor = page.NextCursor
			if params.Cursor == nil {
				return
			}
		}
	}
}

func (c *environmentVariableClient) ListPage(ctx context.Context, params ListEnvironmentVariablesParams) (*EnvironmentVariablePage, error) {
	return observability.WithSpanResult(ctx, c.tracer, "workflows/variables/list", func(ctx context.Context) (*EnvironmentVariablePage, error) {
		response, err := c.service.ListEnvironmentVariables(ctx, workflowsv1.ListEnvironmentVariablesRequest_builder{
			ClusterSlug: params.ClusterSlug, IncludeSecrets: params.IncludeSecrets, Page: paginationFromOptions(params.Limit, params.Cursor),
		}.Build())
		if err != nil {
			return nil, fmt.Errorf("failed to list environment variables: %w", err)
		}
		variables := make([]*EnvironmentVariable, 0, len(response.GetEnvironmentVariables()))
		for _, variable := range response.GetEnvironmentVariables() {
			variables = append(variables, protoToEnvironmentVariable(variable))
		}
		return &EnvironmentVariablePage{EnvironmentVariables: variables, NextCursor: cursorFromPagination(response.GetNextPage())}, nil
	})
}

func (c *environmentVariableClient) Update(ctx context.Context, id uuid.UUID, params UpdateEnvironmentVariableParams) (*EnvironmentVariable, error) {
	return observability.WithSpanResult(ctx, c.tracer, "workflows/variables/update", func(ctx context.Context) (*EnvironmentVariable, error) {
		response, err := c.service.UpdateEnvironmentVariable(ctx, workflowsv1.UpdateEnvironmentVariableRequest_builder{
			Id: tileboxv1.NewUUID(id), Name: params.Name, Value: params.Value, IsSecret: params.IsSecret, Scope: params.Scope.toProto(), Description: params.Description,
		}.Build())
		if err != nil {
			return nil, fmt.Errorf("failed to update environment variable: %w", err)
		}
		return protoToEnvironmentVariable(response), nil
	})
}

func (c *environmentVariableClient) Delete(ctx context.Context, id uuid.UUID) error {
	return observability.WithSpan(ctx, c.tracer, "workflows/variables/delete", func(ctx context.Context) error {
		_, err := c.service.DeleteEnvironmentVariable(ctx, workflowsv1.DeleteEnvironmentVariableRequest_builder{Id: tileboxv1.NewUUID(id)}.Build())
		if err != nil {
			return fmt.Errorf("failed to delete environment variable: %w", err)
		}
		return nil
	})
}

func (s *ClusterScope) toProto() *workflowsv1.ClusterScope {
	if s == nil {
		return nil
	}
	return workflowsv1.ClusterScope_builder{AllClusters: s.AllClusters, ClusterSlugs: s.ClusterSlugs}.Build()
}

func protoToEnvironmentVariable(value *workflowsv1.EnvironmentVariable) *EnvironmentVariable {
	result := &EnvironmentVariable{
		ID: value.GetId().AsUUID(), Name: value.GetName(), IsSecret: value.GetIsSecret(), Description: value.GetDescription(),
		Scope: ClusterScope{AllClusters: value.GetScope().GetAllClusters(), ClusterSlugs: value.GetScope().GetClusterSlugs()},
	}
	if value.HasValue() {
		result.Value = new(value.GetValue())
	}
	if value.HasSecretPreview() {
		result.SecretPreview = new(value.GetSecretPreview())
	}
	return result
}
