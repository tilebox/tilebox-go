// Package workflows provides a client for interacting with Tilebox Workflows.
//
// Documentation: https://docs.tilebox.com/workflows
package workflows

import (
	"context"

	"github.com/google/uuid"
	workflowsv1 "github.com/tilebox/tilebox-go/protogen/workflows/v1"
	"github.com/tilebox/tilebox-go/workflows/v1/cluster"
)

// Cluster represents a Tilebox Workflows cluster.
//
// Documentation: https://docs.tilebox.com/workflows/concepts/clusters
type Cluster struct {
	// Slug is the unique identifier of the cluster within the namespace.
	Slug string
	// Name is the label name of the cluster.
	Name string
	// Description is the cluster description.
	Description string
	// Deletable is true when the cluster can be deleted.
	Deletable bool
	// DeployedWorkflows are the workflows currently deployed to this cluster. Each workflow contains the currently deployed release for that workflow.
	DeployedWorkflows []*Workflow
}

// ClusterRunnerEnvironment contains deployments and resolved runtime variables.
// It includes full secret values: never log or persist it. Deployments and values
// may reflect concurrent changes on the server. Each response replaces the entire
// managed environment, including when EnvironmentVariables is empty.
type ClusterRunnerEnvironment struct {
	Cluster              *Cluster
	EnvironmentVariables []RunnerEnvironmentVariable
}

// RunnerEnvironmentVariable is a literal runtime value, including revealed secrets.
type RunnerEnvironmentVariable struct {
	ID       uuid.UUID
	Name     string
	Value    string
	IsSecret bool
}

type ClusterClient interface {
	// Create creates a new cluster with the given name.
	Create(ctx context.Context, name string, options ...cluster.CreateOption) (*Cluster, error)

	// Update updates an existing cluster by its slug.
	Update(ctx context.Context, slug string, options ...cluster.UpdateOption) (*Cluster, error)

	// Get returns a cluster by its slug.
	Get(ctx context.Context, slug string) (*Cluster, error)

	// GetRunnerEnvironment returns deployments and resolved variables, including secrets.
	// An empty slug selects the default cluster.
	GetRunnerEnvironment(ctx context.Context, slug string) (*ClusterRunnerEnvironment, error)

	// Delete deletes a cluster by its slug.
	Delete(ctx context.Context, slug string) error

	// List returns a list of all available clusters.
	List(ctx context.Context) ([]*Cluster, error)
}

var _ ClusterClient = &clusterClient{}

type clusterClient struct {
	service WorkflowService
}

func (c clusterClient) Create(ctx context.Context, name string, options ...cluster.CreateOption) (*Cluster, error) {
	appliedOptions := cluster.NewCreateOptions(options...)
	response, err := c.service.CreateCluster(ctx, name, appliedOptions.Description, appliedOptions.Slug)
	if err != nil {
		return nil, err
	}

	return protoToCluster(response), nil
}

func (c clusterClient) Update(ctx context.Context, slug string, options ...cluster.UpdateOption) (*Cluster, error) {
	appliedOptions := cluster.NewUpdateOptions(options...)
	response, err := c.service.UpdateCluster(ctx, slug, appliedOptions.Name, appliedOptions.Description)
	if err != nil {
		return nil, err
	}

	return protoToCluster(response), nil
}

func (c clusterClient) Get(ctx context.Context, slug string) (*Cluster, error) {
	response, err := c.service.GetCluster(ctx, slug)
	if err != nil {
		return nil, err
	}

	return protoToCluster(response), nil
}

func (c clusterClient) Delete(ctx context.Context, slug string) error {
	return c.service.DeleteCluster(ctx, slug)
}

func (c clusterClient) GetRunnerEnvironment(ctx context.Context, slug string) (*ClusterRunnerEnvironment, error) {
	response, err := c.service.GetClusterRunnerEnvironment(ctx, slug)
	if err != nil {
		return nil, err
	}
	environment := &ClusterRunnerEnvironment{
		Cluster:              protoToCluster(response.GetCluster()),
		EnvironmentVariables: make([]RunnerEnvironmentVariable, len(response.GetEnvironmentVariables())),
	}
	for i, variable := range response.GetEnvironmentVariables() {
		environment.EnvironmentVariables[i] = RunnerEnvironmentVariable{
			ID: variable.GetId().AsUUID(), Name: variable.GetName(),
			Value: variable.GetValue(), IsSecret: variable.GetIsSecret(),
		}
	}
	return environment, nil
}

func (c clusterClient) List(ctx context.Context) ([]*Cluster, error) {
	response, err := c.service.ListClusters(ctx)
	if err != nil {
		return nil, err
	}

	clusters := make([]*Cluster, len(response.GetClusters()))
	for i, c := range response.GetClusters() {
		clusters[i] = protoToCluster(c)
	}

	return clusters, nil
}

func protoToCluster(c *workflowsv1.Cluster) *Cluster {
	deployedWorkflows := make([]*Workflow, len(c.GetDeployedReleases()))
	for i, workflow := range c.GetDeployedReleases() {
		deployedWorkflows[i] = protoToWorkflow(workflow)
	}

	return &Cluster{
		Slug:              c.GetSlug(),
		Name:              c.GetDisplayName(),
		Description:       c.GetDescription(),
		Deletable:         c.GetDeletable(),
		DeployedWorkflows: deployedWorkflows,
	}
}
