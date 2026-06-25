package workflows

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tileboxv1 "github.com/tilebox/tilebox-go/protogen/tilebox/v1"
	workflowsv1 "github.com/tilebox/tilebox-go/protogen/workflows/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestWorkflowClient_Create(t *testing.T) {
	ctx := context.Background()
	service := &fakeWorkflowService{
		workflow: workflowsv1.Workflow_builder{
			Slug:        "agentic-workflow",
			Name:        "Agentic Workflow",
			Description: "Description",
		}.Build(),
	}
	client := workflowClient{service: service}

	workflow, err := client.Create(ctx, "Agentic Workflow", WithDescription("Description"))
	require.NoError(t, err)

	assert.Equal(t, "Agentic Workflow", service.createName)
	assert.Equal(t, "Description", service.createDescription)
	assert.Equal(t, "agentic-workflow", workflow.Slug)
	assert.Equal(t, "Agentic Workflow", workflow.Name)
	assert.Equal(t, "Description", workflow.Description)
}

func TestWorkflowClient_List(t *testing.T) {
	ctx := context.Background()
	service := &fakeWorkflowService{
		listWorkflowsResponse: workflowsv1.ListWorkflowsResponse_builder{
			Workflows: []*workflowsv1.Workflow{
				workflowsv1.Workflow_builder{Slug: "one", Name: "One"}.Build(),
				workflowsv1.Workflow_builder{Slug: "two", Name: "Two"}.Build(),
			},
		}.Build(),
	}
	client := workflowClient{service: service}

	workflows, err := client.List(ctx)
	require.NoError(t, err)

	require.Len(t, workflows, 2)
	assert.Equal(t, "one", workflows[0].Slug)
	assert.Equal(t, "two", workflows[1].Slug)
}

func TestWorkflowClient_Get(t *testing.T) {
	ctx := context.Background()
	service := &fakeWorkflowService{
		workflow: workflowsv1.Workflow_builder{Slug: "agentic-workflow", Name: "Agentic Workflow"}.Build(),
	}
	client := workflowClient{service: service}

	workflow, err := client.Get(ctx, "agentic-workflow")
	require.NoError(t, err)

	assert.Equal(t, "agentic-workflow", service.getSlug)
	assert.Equal(t, "Agentic Workflow", workflow.Name)
}

func TestWorkflowClient_Delete(t *testing.T) {
	ctx := context.Background()
	service := &fakeWorkflowService{}
	client := workflowClient{service: service}

	err := client.Delete(ctx, "agentic-workflow")
	require.NoError(t, err)

	assert.Equal(t, "agentic-workflow", service.deleteWorkflowSlug)
}

func TestWorkflowClient_PublishRelease(t *testing.T) {
	ctx := context.Background()
	releaseID := uuid.New()
	artifactID := uuid.New()
	digest := strings.Repeat("a", 64)
	fingerprint := strings.Repeat("b", 64)
	createdAt := time.Date(2026, time.May, 29, 12, 0, 0, 0, time.UTC)
	service := &fakeWorkflowService{
		workflowRelease: workflowsv1.WorkflowRelease_builder{
			Id: tileboxv1.NewUUID(releaseID),
			Artifact: workflowsv1.Artifact_builder{
				Id:     tileboxv1.NewUUID(artifactID),
				Digest: digest,
			}.Build(),
			Content: releaseContentToProtoMust(t, &ReleaseContent{
				Fingerprint: fingerprint,
				Tasks:       []TaskIdentifier{NewTaskIdentifier("tilebox.com/task/Review", "v1.0")},
				Files: []*Path{{
					Path:      ".",
					Directory: true,
					Children:  []*Path{{Path: "main.py"}},
				}},
				RunnerObjectPath: "my_module.my_runner:runner",
				CommandOverride:  []string{"python", "main.py"},
			}),
			CreatedAt: timestamppb.New(createdAt),
		}.Build(),
	}
	client := workflowClient{service: service}
	content := &ReleaseContent{
		Fingerprint:      fingerprint,
		Tasks:            []TaskIdentifier{NewTaskIdentifier("tilebox.com/task/Review", "v1.0")},
		Files:            []*Path{{Path: ".", Directory: true, Children: []*Path{{Path: "main.py"}}}},
		RunnerObjectPath: "my_module.my_runner:runner",
		CommandOverride:  []string{"python", "main.py"},
	}

	release, err := client.PublishRelease(ctx, "agentic-workflow", artifactID, content)
	require.NoError(t, err)

	assert.Equal(t, "agentic-workflow", service.publishWorkflowSlug)
	assert.Equal(t, artifactID, service.publishArtifactID)
	assert.Equal(t, content, service.publishContent)
	assert.Equal(t, releaseID, release.ID)
	require.NotNil(t, release.Artifact)
	assert.Equal(t, artifactID, release.Artifact.ID)
	assert.Equal(t, digest, release.Artifact.Digest)
	require.NotNil(t, release.Content)
	assert.Equal(t, fingerprint, release.Content.Fingerprint)
	require.Len(t, release.Content.Tasks, 1)
	assert.Equal(t, "tilebox.com/task/Review", release.Content.Tasks[0].Name())
	assert.Equal(t, "v1.0", release.Content.Tasks[0].Version())
	require.Len(t, release.Content.Files, 1)
	assert.Equal(t, ".", release.Content.Files[0].Path)
	assert.True(t, release.Content.Files[0].Directory)
	require.Len(t, release.Content.Files[0].Children, 1)
	assert.Equal(t, "main.py", release.Content.Files[0].Children[0].Path)
	assert.Equal(t, "my_module.my_runner:runner", release.Content.RunnerObjectPath)
	assert.Equal(t, []string{"python", "main.py"}, release.Content.CommandOverride)
	assert.Equal(t, createdAt, release.CreatedAt)
}

func TestWorkflowClient_UnpublishRelease(t *testing.T) {
	ctx := context.Background()
	releaseID := uuid.New()
	service := &fakeWorkflowService{}
	client := workflowClient{service: service}

	err := client.UnpublishRelease(ctx, "agentic-workflow", releaseID)
	require.NoError(t, err)

	assert.Equal(t, "agentic-workflow", service.unpublishWorkflowSlug)
	assert.Equal(t, releaseID, service.unpublishReleaseID)
}

func TestWorkflowClient_DeployRelease(t *testing.T) {
	ctx := context.Background()
	releaseID := uuid.New()
	service := &fakeWorkflowService{
		deployWorkflowReleaseResponse: workflowsv1.DeployWorkflowReleaseResponse_builder{
			Release: workflowsv1.WorkflowRelease_builder{Id: tileboxv1.NewUUID(releaseID)}.Build(),
			Clusters: []*workflowsv1.Cluster{
				workflowsv1.Cluster_builder{Slug: "dev", DisplayName: "Dev"}.Build(),
			},
		}.Build(),
	}
	client := workflowClient{service: service}

	deployment, err := client.DeployRelease(ctx, "agentic-workflow", releaseID, []string{"dev"})
	require.NoError(t, err)

	assert.Equal(t, "agentic-workflow", service.deployWorkflowSlug)
	assert.Equal(t, releaseID, service.deployReleaseID)
	assert.Equal(t, []string{"dev"}, service.deployClusterSlugs)
	require.NotNil(t, deployment.Release)
	assert.Equal(t, releaseID, deployment.Release.ID)
	require.Len(t, deployment.Clusters, 1)
	assert.Equal(t, "dev", deployment.Clusters[0].Slug)
	assert.Equal(t, "Dev", deployment.Clusters[0].Name)
}

func TestWorkflowClient_UndeployRelease(t *testing.T) {
	ctx := context.Background()
	releaseID := uuid.New()
	service := &fakeWorkflowService{
		undeployWorkflowReleaseResponse: workflowsv1.UndeployWorkflowReleaseResponse_builder{
			Release: workflowsv1.WorkflowRelease_builder{Id: tileboxv1.NewUUID(releaseID)}.Build(),
			Clusters: []*workflowsv1.Cluster{
				workflowsv1.Cluster_builder{Slug: "dev", DisplayName: "Dev"}.Build(),
			},
		}.Build(),
	}
	client := workflowClient{service: service}

	deployment, err := client.UndeployRelease(ctx, "agentic-workflow", releaseID, []string{"dev"})
	require.NoError(t, err)

	assert.Equal(t, "agentic-workflow", service.undeployWorkflowSlug)
	assert.Equal(t, releaseID, service.undeployReleaseID)
	assert.Equal(t, []string{"dev"}, service.undeployClusterSlugs)
	require.NotNil(t, deployment.Release)
	assert.Equal(t, releaseID, deployment.Release.ID)
	require.Len(t, deployment.Clusters, 1)
	assert.Equal(t, "dev", deployment.Clusters[0].Slug)
	assert.Equal(t, "Dev", deployment.Clusters[0].Name)
}

func TestProtoToCluster_MapsDeployedWorkflows(t *testing.T) {
	releaseID := uuid.New()
	cluster := protoToCluster(workflowsv1.Cluster_builder{
		Slug:        "dev",
		DisplayName: "Dev",
		Deletable:   true,
		DeployedReleases: []*workflowsv1.Workflow{
			workflowsv1.Workflow_builder{
				Slug: "agentic-workflow",
				Name: "Agentic Workflow",
				Releases: []*workflowsv1.WorkflowRelease{
					workflowsv1.WorkflowRelease_builder{Id: tileboxv1.NewUUID(releaseID)}.Build(),
				},
			}.Build(),
		},
	}.Build())

	assert.Equal(t, "dev", cluster.Slug)
	assert.Equal(t, "Dev", cluster.Name)
	assert.True(t, cluster.Deletable)
	require.Len(t, cluster.DeployedWorkflows, 1)
	assert.Equal(t, "agentic-workflow", cluster.DeployedWorkflows[0].Slug)
	require.Len(t, cluster.DeployedWorkflows[0].Releases, 1)
	assert.Equal(t, releaseID, cluster.DeployedWorkflows[0].Releases[0].ID)
}

func TestTaskIdentifiersToProto_ValidatesIdentifiers(t *testing.T) {
	_, err := taskIdentifiersToProto([]TaskIdentifier{NewTaskIdentifier("task", "invalid")})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid task version")
}

func TestTaskIdentifiersToProto_RejectsNilIdentifiers(t *testing.T) {
	_, err := taskIdentifiersToProto([]TaskIdentifier{nil})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "task identifier at index 0 is nil")
}

func TestReleaseContentToProto_RejectsNilContent(t *testing.T) {
	_, err := releaseContentToProto(nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "release content is nil")
}

func TestProtoToWorkflow_HandlesNil(t *testing.T) {
	assert.Nil(t, protoToWorkflow(nil))
	assert.Nil(t, protoToWorkflowRelease(nil))
	assert.Nil(t, protoToArtifact(nil))
	assert.Nil(t, protoToReleaseContent(nil))
	assert.Nil(t, protoToPath(nil))
	assert.Nil(t, protoToTaskIdentifier(nil))
}

type fakeWorkflowService struct {
	workflow                        *workflowsv1.Workflow
	listWorkflowsResponse           *workflowsv1.ListWorkflowsResponse
	workflowRelease                 *workflowsv1.WorkflowRelease
	deployWorkflowReleaseResponse   *workflowsv1.DeployWorkflowReleaseResponse
	undeployWorkflowReleaseResponse *workflowsv1.UndeployWorkflowReleaseResponse
	err                             error

	createName         string
	createDescription  string
	getSlug            string
	deleteWorkflowSlug string

	publishWorkflowSlug string
	publishArtifactID   uuid.UUID
	publishContent      *ReleaseContent

	unpublishWorkflowSlug string
	unpublishReleaseID    uuid.UUID

	deployWorkflowSlug string
	deployReleaseID    uuid.UUID
	deployClusterSlugs []string

	undeployWorkflowSlug string
	undeployReleaseID    uuid.UUID
	undeployClusterSlugs []string
}

func (s *fakeWorkflowService) CreateCluster(context.Context, string) (*workflowsv1.Cluster, error) {
	return nil, errors.New("not implemented")
}

func (s *fakeWorkflowService) GetCluster(context.Context, string) (*workflowsv1.Cluster, error) {
	return nil, errors.New("not implemented")
}

func (s *fakeWorkflowService) DeleteCluster(context.Context, string) error {
	return errors.New("not implemented")
}

func (s *fakeWorkflowService) ListClusters(context.Context) (*workflowsv1.ListClustersResponse, error) {
	return nil, errors.New("not implemented")
}

func (s *fakeWorkflowService) CreateWorkflow(_ context.Context, name, description string) (*workflowsv1.Workflow, error) {
	s.createName = name
	s.createDescription = description
	return s.workflow, s.err
}

func (s *fakeWorkflowService) ListWorkflows(context.Context) (*workflowsv1.ListWorkflowsResponse, error) {
	return s.listWorkflowsResponse, s.err
}

func (s *fakeWorkflowService) GetWorkflow(_ context.Context, slug string) (*workflowsv1.Workflow, error) {
	s.getSlug = slug
	return s.workflow, s.err
}

func (s *fakeWorkflowService) DeleteWorkflow(_ context.Context, slug string) error {
	s.deleteWorkflowSlug = slug
	return s.err
}

func (s *fakeWorkflowService) PublishWorkflowRelease(_ context.Context, workflowSlug string, artifactID uuid.UUID, content *ReleaseContent) (*workflowsv1.WorkflowRelease, error) {
	s.publishWorkflowSlug = workflowSlug
	s.publishArtifactID = artifactID
	s.publishContent = content
	return s.workflowRelease, s.err
}

func (s *fakeWorkflowService) UnpublishWorkflowRelease(_ context.Context, workflowSlug string, releaseID uuid.UUID) error {
	s.unpublishWorkflowSlug = workflowSlug
	s.unpublishReleaseID = releaseID
	return s.err
}

func (s *fakeWorkflowService) DeployWorkflowRelease(_ context.Context, workflowSlug string, releaseID uuid.UUID, clusterSlugs []string) (*workflowsv1.DeployWorkflowReleaseResponse, error) {
	s.deployWorkflowSlug = workflowSlug
	s.deployReleaseID = releaseID
	s.deployClusterSlugs = clusterSlugs
	return s.deployWorkflowReleaseResponse, s.err
}

func (s *fakeWorkflowService) UndeployWorkflowRelease(_ context.Context, workflowSlug string, releaseID uuid.UUID, clusterSlugs []string) (*workflowsv1.UndeployWorkflowReleaseResponse, error) {
	s.undeployWorkflowSlug = workflowSlug
	s.undeployReleaseID = releaseID
	s.undeployClusterSlugs = clusterSlugs
	return s.undeployWorkflowReleaseResponse, s.err
}

func releaseContentToProtoMust(tb testing.TB, content *ReleaseContent) *workflowsv1.ReleaseContent {
	tb.Helper()
	protoContent, err := releaseContentToProto(content)
	require.NoError(tb, err)
	return protoContent
}
