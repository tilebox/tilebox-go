package workflows

import (
	"context"
	"errors"
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	workflowsv1 "github.com/tilebox/tilebox-go/protogen/workflows/v1"
	"github.com/tilebox/tilebox-go/workflows/v1/cluster"
	"go.opentelemetry.io/otel/trace/noop"
)

func Test_clusterClient_Get(t *testing.T) {
	ctx := context.Background()
	client := NewReplayClient(t, "cluster")

	cluster, err := client.Clusters.Get(ctx, "dev-cluster-ESugE7S4cwADVK")
	require.NoError(t, err)

	assert.Equal(t, "dev-cluster", cluster.Name)
	assert.Equal(t, "dev-cluster-ESugE7S4cwADVK", cluster.Slug)
}

func Test_clusterClient_Create_WithOptions(t *testing.T) {
	ctx := context.Background()
	service := &fakeWorkflowService{
		cluster: workflowsv1.Cluster_builder{
			Slug:        "dev",
			DisplayName: "Dev",
			Description: "Development cluster",
		}.Build(),
	}
	client := clusterClient{service: service}

	cluster, err := client.Create(ctx, "Dev", cluster.WithDescription("Development cluster"), cluster.WithSlug("dev"))
	require.NoError(t, err)

	assert.Equal(t, "Dev", service.createClusterName)
	assert.Equal(t, "Development cluster", service.createClusterDescription)
	assert.Equal(t, "dev", service.createClusterSlug)
	assert.Equal(t, "dev", cluster.Slug)
	assert.Equal(t, "Development cluster", cluster.Description)
}

func Test_clusterClient_Update(t *testing.T) {
	ctx := context.Background()
	service := &fakeWorkflowService{
		cluster: workflowsv1.Cluster_builder{
			Slug:        "dev",
			DisplayName: "Dev",
			Description: "Development cluster",
		}.Build(),
	}
	client := clusterClient{service: service}

	cluster, err := client.Update(ctx, "dev", cluster.WithName("Dev"), cluster.WithDescription("Development cluster"))
	require.NoError(t, err)

	assert.Equal(t, "dev", service.updateClusterSlug)
	require.NotNil(t, service.updateClusterName)
	require.NotNil(t, service.updateClusterDescription)
	assert.Equal(t, "Dev", *service.updateClusterName)
	assert.Equal(t, "Development cluster", *service.updateClusterDescription)
	assert.Equal(t, "Dev", cluster.Name)
	assert.Equal(t, "Development cluster", cluster.Description)
}

func Test_workflowService_CreateCluster_SendsDescriptionAndSlug(t *testing.T) {
	ctx := context.Background()
	connectClient := &fakeWorkflowsConnectClient{}
	service := newWorkflowService(connectClient, noop.NewTracerProvider().Tracer("test"))

	_, err := service.CreateCluster(ctx, "Dev", "Development cluster", "dev")
	require.NoError(t, err)

	require.NotNil(t, connectClient.createClusterRequest)
	assert.Equal(t, "Dev", connectClient.createClusterRequest.GetName())
	assert.Equal(t, "Development cluster", connectClient.createClusterRequest.GetDescription())
	assert.Equal(t, "dev", connectClient.createClusterRequest.GetSlug())
}

func Test_workflowService_UpdateCluster_PreservesOptionalPresence(t *testing.T) {
	ctx := context.Background()
	connectClient := &fakeWorkflowsConnectClient{}
	service := newWorkflowService(connectClient, noop.NewTracerProvider().Tracer("test"))

	_, err := service.UpdateCluster(ctx, "dev", nil, nil)
	require.NoError(t, err)
	require.NotNil(t, connectClient.updateClusterRequest)
	assert.Equal(t, "dev", connectClient.updateClusterRequest.GetClusterSlug())
	assert.False(t, connectClient.updateClusterRequest.HasName())
	assert.False(t, connectClient.updateClusterRequest.HasDescription())

	connectClient.updateClusterRequest = nil
	emptyDescription := ""
	_, err = service.UpdateCluster(ctx, "dev", nil, &emptyDescription)
	require.NoError(t, err)
	require.NotNil(t, connectClient.updateClusterRequest)
	assert.False(t, connectClient.updateClusterRequest.HasName())
	assert.True(t, connectClient.updateClusterRequest.HasDescription())
	assert.Empty(t, connectClient.updateClusterRequest.GetDescription())
}

func Test_clusterClient_List(t *testing.T) {
	ctx := context.Background()
	client := NewReplayClient(t, "clusters")

	clusters, err := client.Clusters.List(ctx)
	require.NoError(t, err)

	require.Len(t, clusters, 2)
	cluster := clusters[1]
	assert.Equal(t, "dev-cluster", cluster.Name)
	assert.Equal(t, "dev-cluster-ESugE7S4cwADVK", cluster.Slug)
}

type fakeWorkflowsConnectClient struct {
	createClusterRequest  *workflowsv1.CreateClusterRequest
	updateClusterRequest  *workflowsv1.UpdateClusterRequest
	updateWorkflowRequest *workflowsv1.UpdateWorkflowRequest
}

func (c *fakeWorkflowsConnectClient) CreateCluster(_ context.Context, req *connect.Request[workflowsv1.CreateClusterRequest]) (*connect.Response[workflowsv1.Cluster], error) {
	c.createClusterRequest = req.Msg
	return connect.NewResponse(workflowsv1.Cluster_builder{Slug: "dev"}.Build()), nil
}

func (c *fakeWorkflowsConnectClient) UpdateCluster(_ context.Context, req *connect.Request[workflowsv1.UpdateClusterRequest]) (*connect.Response[workflowsv1.Cluster], error) {
	c.updateClusterRequest = req.Msg
	return connect.NewResponse(workflowsv1.Cluster_builder{Slug: "dev"}.Build()), nil
}

func (c *fakeWorkflowsConnectClient) GetCluster(context.Context, *connect.Request[workflowsv1.GetClusterRequest]) (*connect.Response[workflowsv1.Cluster], error) {
	return nil, errors.New("not implemented")
}

func (c *fakeWorkflowsConnectClient) DeleteCluster(context.Context, *connect.Request[workflowsv1.DeleteClusterRequest]) (*connect.Response[workflowsv1.DeleteClusterResponse], error) {
	return nil, errors.New("not implemented")
}

func (c *fakeWorkflowsConnectClient) ListClusters(context.Context, *connect.Request[workflowsv1.ListClustersRequest]) (*connect.Response[workflowsv1.ListClustersResponse], error) {
	return nil, errors.New("not implemented")
}

func (c *fakeWorkflowsConnectClient) CreateWorkflow(context.Context, *connect.Request[workflowsv1.CreateWorkflowRequest]) (*connect.Response[workflowsv1.Workflow], error) {
	return nil, errors.New("not implemented")
}

func (c *fakeWorkflowsConnectClient) ListWorkflows(context.Context, *connect.Request[workflowsv1.ListWorkflowsRequest]) (*connect.Response[workflowsv1.ListWorkflowsResponse], error) {
	return nil, errors.New("not implemented")
}

func (c *fakeWorkflowsConnectClient) GetWorkflow(context.Context, *connect.Request[workflowsv1.GetWorkflowRequest]) (*connect.Response[workflowsv1.Workflow], error) {
	return nil, errors.New("not implemented")
}

func (c *fakeWorkflowsConnectClient) UpdateWorkflow(_ context.Context, req *connect.Request[workflowsv1.UpdateWorkflowRequest]) (*connect.Response[workflowsv1.Workflow], error) {
	c.updateWorkflowRequest = req.Msg
	return connect.NewResponse(workflowsv1.Workflow_builder{Slug: "agentic-workflow"}.Build()), nil
}

func (c *fakeWorkflowsConnectClient) DeleteWorkflow(context.Context, *connect.Request[workflowsv1.DeleteWorkflowRequest]) (*connect.Response[workflowsv1.DeleteWorkflowResponse], error) {
	return nil, errors.New("not implemented")
}

func (c *fakeWorkflowsConnectClient) PublishWorkflowRelease(context.Context, *connect.Request[workflowsv1.PublishWorkflowReleaseRequest]) (*connect.Response[workflowsv1.WorkflowRelease], error) {
	return nil, errors.New("not implemented")
}

func (c *fakeWorkflowsConnectClient) UnpublishWorkflowRelease(context.Context, *connect.Request[workflowsv1.UnpublishWorkflowReleaseRequest]) (*connect.Response[workflowsv1.UnpublishWorkflowReleaseResponse], error) {
	return nil, errors.New("not implemented")
}

func (c *fakeWorkflowsConnectClient) DeployWorkflowRelease(context.Context, *connect.Request[workflowsv1.DeployWorkflowReleaseRequest]) (*connect.Response[workflowsv1.DeployWorkflowReleaseResponse], error) {
	return nil, errors.New("not implemented")
}

func (c *fakeWorkflowsConnectClient) UndeployWorkflowRelease(context.Context, *connect.Request[workflowsv1.UndeployWorkflowReleaseRequest]) (*connect.Response[workflowsv1.UndeployWorkflowReleaseResponse], error) {
	return nil, errors.New("not implemented")
}
