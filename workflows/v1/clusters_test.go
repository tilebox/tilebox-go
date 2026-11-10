package workflows

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"connectrpc.com/connect/v2"
	"connectrpc.com/connect/v2/connecthttp"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tileboxv1 "github.com/tilebox/tilebox-go/protogen/tilebox/v1"
	workflowsv1 "github.com/tilebox/tilebox-go/protogen/workflows/v1"
	"github.com/tilebox/tilebox-go/protogen/workflows/v1/workflowsv1connect"
	"github.com/tilebox/tilebox-go/workflows/v1/cluster"
	"go.opentelemetry.io/otel/trace/noop"
)

type runnerEnvironmentAPI struct {
	workflowsv1connect.UnimplementedWorkflowsServiceHandler

	get func(*workflowsv1.GetClusterRunnerEnvironmentRequest) (*workflowsv1.GetClusterRunnerEnvironmentResponse, error)
}

func (s *runnerEnvironmentAPI) GetClusterRunnerEnvironment(_ context.Context, request *workflowsv1.GetClusterRunnerEnvironmentRequest) (*workflowsv1.GetClusterRunnerEnvironmentResponse, error) {
	return s.get(request)
}

func TestClusterRunnerEnvironment(t *testing.T) {
	id, releaseID := uuid.New(), uuid.New()
	for _, slug := range []string{"", "workers"} {
		t.Run("cluster="+slug, func(t *testing.T) {
			service := &runnerEnvironmentAPI{get: func(request *workflowsv1.GetClusterRunnerEnvironmentRequest) (*workflowsv1.GetClusterRunnerEnvironmentResponse, error) {
				assert.Equal(t, slug, request.GetClusterSlug())
				return workflowsv1.GetClusterRunnerEnvironmentResponse_builder{
					Cluster: workflowsv1.Cluster_builder{Slug: slug, DeployedReleases: []*workflowsv1.Workflow{
						workflowsv1.Workflow_builder{Releases: []*workflowsv1.WorkflowRelease{
							workflowsv1.WorkflowRelease_builder{Id: tileboxv1.NewUUID(releaseID)}.Build(),
						}}.Build(),
					}}.Build(),
					EnvironmentVariables: []*workflowsv1.RunnerEnvironmentVariable{
						workflowsv1.RunnerEnvironmentVariable_builder{Id: tileboxv1.NewUUID(id), Name: "TOKEN", Value: " literal\n$VALUE=a ", IsSecret: true}.Build(),
					},
				}.Build(), nil
			}}
			server := connect.NewServer()
			workflowsv1connect.RegisterWorkflowsServiceHandler(server, service)
			mux := http.NewServeMux()
			connecthttp.Mount(mux, server)
			httpServer := httptest.NewServer(mux)
			t.Cleanup(httpServer.Close)
			client := NewClient(connect.NewClient(connecthttp.NewTransport(httpServer.Client(), httpServer.URL)), WithDisableTracing())
			result, err := client.Clusters.GetRunnerEnvironment(t.Context(), slug)
			require.NoError(t, err)
			assert.Equal(t, slug, result.Cluster.Slug)
			require.Len(t, result.Cluster.DeployedWorkflows, 1)
			require.Len(t, result.Cluster.DeployedWorkflows[0].Releases, 1)
			assert.Equal(t, releaseID, result.Cluster.DeployedWorkflows[0].Releases[0].ID)
			assert.Equal(t, []RunnerEnvironmentVariable{{ID: id, Name: "TOKEN", Value: " literal\n$VALUE=a ", IsSecret: true}}, result.EnvironmentVariables)

			service.get = func(*workflowsv1.GetClusterRunnerEnvironmentRequest) (*workflowsv1.GetClusterRunnerEnvironmentResponse, error) {
				return workflowsv1.GetClusterRunnerEnvironmentResponse_builder{Cluster: workflowsv1.Cluster_builder{Slug: slug}.Build()}.Build(), nil
			}
			result, err = client.Clusters.GetRunnerEnvironment(t.Context(), slug)
			require.NoError(t, err)
			assert.Empty(t, result.EnvironmentVariables)
			for _, code := range []connect.Code{connect.CodeUnimplemented, connect.CodePermissionDenied} {
				service.get = func(*workflowsv1.GetClusterRunnerEnvironmentRequest) (*workflowsv1.GetClusterRunnerEnvironmentResponse, error) {
					return nil, connect.NewError(code, "unavailable")
				}
				_, err = client.Clusters.GetRunnerEnvironment(t.Context(), slug)
				assert.Equal(t, code, connect.CodeOf(err))
			}
		})
	}
}

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
	createClusterRequest           *workflowsv1.CreateClusterRequest
	updateClusterRequest           *workflowsv1.UpdateClusterRequest
	updateWorkflowRequest          *workflowsv1.UpdateWorkflowRequest
	listWorkflowsRequest           *workflowsv1.ListWorkflowsRequest
	listPublicWorkflowsRequest     *workflowsv1.ListPublicWorkflowsRequest
	undeployWorkflowReleaseRequest *workflowsv1.UndeployWorkflowReleaseRequest
}

func (c *fakeWorkflowsConnectClient) CreateCluster(_ context.Context, req *workflowsv1.CreateClusterRequest) (*workflowsv1.Cluster, error) {
	c.createClusterRequest = req
	return workflowsv1.Cluster_builder{Slug: "dev"}.Build(), nil
}

func (c *fakeWorkflowsConnectClient) UpdateCluster(_ context.Context, req *workflowsv1.UpdateClusterRequest) (*workflowsv1.Cluster, error) {
	c.updateClusterRequest = req
	return workflowsv1.Cluster_builder{Slug: "dev"}.Build(), nil
}

func (c *fakeWorkflowsConnectClient) GetCluster(context.Context, *workflowsv1.GetClusterRequest) (*workflowsv1.Cluster, error) {
	return nil, errors.New("not implemented")
}

func (c *fakeWorkflowsConnectClient) GetClusterRunnerEnvironment(context.Context, *workflowsv1.GetClusterRunnerEnvironmentRequest) (*workflowsv1.GetClusterRunnerEnvironmentResponse, error) {
	return nil, errors.New("not implemented")
}

func (c *fakeWorkflowsConnectClient) DeleteCluster(context.Context, *workflowsv1.DeleteClusterRequest) (*workflowsv1.DeleteClusterResponse, error) {
	return nil, errors.New("not implemented")
}

func (c *fakeWorkflowsConnectClient) ListClusters(context.Context, *workflowsv1.ListClustersRequest) (*workflowsv1.ListClustersResponse, error) {
	return nil, errors.New("not implemented")
}

func (c *fakeWorkflowsConnectClient) CreateWorkflow(context.Context, *workflowsv1.CreateWorkflowRequest) (*workflowsv1.Workflow, error) {
	return nil, errors.New("not implemented")
}

func (c *fakeWorkflowsConnectClient) ListWorkflows(_ context.Context, req *workflowsv1.ListWorkflowsRequest) (*workflowsv1.ListWorkflowsResponse, error) {
	c.listWorkflowsRequest = req
	return &workflowsv1.ListWorkflowsResponse{}, nil
}

func (c *fakeWorkflowsConnectClient) ListPublicWorkflows(_ context.Context, req *workflowsv1.ListPublicWorkflowsRequest) (*workflowsv1.ListPublicWorkflowsResponse, error) {
	c.listPublicWorkflowsRequest = req
	return &workflowsv1.ListPublicWorkflowsResponse{}, nil
}

func (c *fakeWorkflowsConnectClient) GetWorkflow(context.Context, *workflowsv1.GetWorkflowRequest) (*workflowsv1.Workflow, error) {
	return nil, errors.New("not implemented")
}

func (c *fakeWorkflowsConnectClient) UpdateWorkflow(_ context.Context, req *workflowsv1.UpdateWorkflowRequest) (*workflowsv1.Workflow, error) {
	c.updateWorkflowRequest = req
	return workflowsv1.Workflow_builder{Reference: "agentic-workflow"}.Build(), nil
}

func (c *fakeWorkflowsConnectClient) DeleteWorkflow(context.Context, *workflowsv1.DeleteWorkflowRequest) (*workflowsv1.DeleteWorkflowResponse, error) {
	return nil, errors.New("not implemented")
}

func (c *fakeWorkflowsConnectClient) PublishWorkflowRelease(context.Context, *workflowsv1.PublishWorkflowReleaseRequest) (*workflowsv1.WorkflowRelease, error) {
	return nil, errors.New("not implemented")
}

func (c *fakeWorkflowsConnectClient) UnpublishWorkflowRelease(context.Context, *workflowsv1.UnpublishWorkflowReleaseRequest) (*workflowsv1.UnpublishWorkflowReleaseResponse, error) {
	return nil, errors.New("not implemented")
}

func (c *fakeWorkflowsConnectClient) DeployWorkflowRelease(context.Context, *workflowsv1.DeployWorkflowReleaseRequest) (*workflowsv1.DeployWorkflowReleaseResponse, error) {
	return nil, errors.New("not implemented")
}

func (c *fakeWorkflowsConnectClient) UndeployWorkflowRelease(_ context.Context, req *workflowsv1.UndeployWorkflowReleaseRequest) (*workflowsv1.UndeployWorkflowReleaseResponse, error) {
	c.undeployWorkflowReleaseRequest = req
	return &workflowsv1.UndeployWorkflowReleaseResponse{}, nil
}
