// Package workflows provides a client for interacting with Tilebox Workflows.
//
// Documentation: https://docs.tilebox.com/workflows
package workflows // import "github.com/tilebox/tilebox-go/workflows/v1"

import (
	"context"

	workflowsv1 "github.com/tilebox/tilebox-go/protogen/go/workflows/v1"
)

// Cluster represents a Tilebox Workflows cluster.
//
// Documentation: https://docs.tilebox.com/workflows/concepts/clusters
type Cluster struct {
	// Slug is the unique identifier of the cluster within the namespace.
	Slug string
	// Name is the display name of the cluster.
	Name string
	// Deletable is true when the cluster can be deleted.
	Deletable bool
}

type ClusterClient interface {
	// Create creates a new cluster with the given name.
	Create(ctx context.Context, name string) (*Cluster, error)

	// Get returns a cluster by its slug.
	Get(ctx context.Context, slug string) (*Cluster, error)

	// Delete deletes a cluster by its slug.
	Delete(ctx context.Context, slug string) error

	// List returns a list of all available clusters.
	List(ctx context.Context) ([]*Cluster, error)
}

var _ ClusterClient = &clusterClient{}

type clusterClient struct {
	service WorkflowService
}

func (c clusterClient) Create(ctx context.Context, name string) (*Cluster, error) {
	response, err := c.service.CreateCluster(ctx, name)
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
	return &Cluster{
		Slug:      c.GetSlug(),
		Name:      c.GetDisplayName(),
		Deletable: c.GetDeletable(),
	}
}
