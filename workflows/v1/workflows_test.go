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
	"github.com/tilebox/tilebox-go/workflows/v1/workflow"
	"go.opentelemetry.io/otel/trace/noop"
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

	workflows, err := Collect(client.List(ctx))
	require.NoError(t, err)

	require.Len(t, workflows, 2)
	assert.Equal(t, "one", workflows[0].Slug)
	assert.Equal(t, "two", workflows[1].Slug)
}

func TestWorkflowClient_ListAndListPublic_Paginate(t *testing.T) {
	nextID := uuid.New()
	owned := &pagingWorkflowService{owned: []*workflowsv1.ListWorkflowsResponse{
		workflowsv1.ListWorkflowsResponse_builder{Workflows: []*workflowsv1.Workflow{workflowsv1.Workflow_builder{Slug: "one"}.Build()}, NextPage: tileboxv1.Pagination_builder{StartingAfter: tileboxv1.NewUUID(nextID)}.Build()}.Build(),
		workflowsv1.ListWorkflowsResponse_builder{Workflows: []*workflowsv1.Workflow{workflowsv1.Workflow_builder{Slug: "two"}.Build()}}.Build(),
	}}
	client := workflowClient{service: owned}

	workflows, err := Collect(client.List(context.Background()))
	require.NoError(t, err)
	assert.Equal(t, []string{"one", "two"}, []string{workflows[0].Slug, workflows[1].Slug})
	require.Len(t, owned.ownedPages, 2)
	assert.Equal(t, nextID, owned.ownedPages[1].GetStartingAfter().AsUUID())

	public := &pagingWorkflowService{public: []*workflowsv1.ListPublicWorkflowsResponse{
		workflowsv1.ListPublicWorkflowsResponse_builder{Workflows: []*workflowsv1.Workflow{workflowsv1.Workflow_builder{Slug: "tilebox:example"}.Build()}, NextPage: tileboxv1.Pagination_builder{StartingAfter: tileboxv1.NewUUID(nextID)}.Build()}.Build(),
		workflowsv1.ListPublicWorkflowsResponse_builder{Workflows: []*workflowsv1.Workflow{workflowsv1.Workflow_builder{Slug: "other:example"}.Build()}}.Build(),
	}}
	client.service = public
	workflows, err = Collect(client.ListPublic(context.Background()))
	require.NoError(t, err)
	require.Len(t, workflows, 2)
	assert.Equal(t, "tilebox:example", workflows[0].Slug)
	assert.Equal(t, "other:example", workflows[1].Slug)
	require.Len(t, public.publicPages, 2)
	assert.Equal(t, nextID, public.publicPages[1].GetStartingAfter().AsUUID())
	assert.Empty(t, workflows[0].Releases)
}

func TestWorkflowClient_ListIterators(t *testing.T) {
	for _, public := range []bool{false, true} {
		for _, test := range []struct {
			name      string
			limit     int64
			stopEarly bool
			want      []string
			calls     int
		}{
			{name: "unlimited", want: []string{"one", "two", "three", "four"}, calls: 2},
			{name: "negative limit", limit: -1, want: []string{"one", "two", "three", "four"}, calls: 2},
			{name: "within page", limit: 1, want: []string{"one"}, calls: 1},
			{name: "page boundary", limit: 2, want: []string{"one", "two"}, calls: 1},
			{name: "across pages", limit: 3, want: []string{"one", "two", "three"}, calls: 2},
			{name: "early stop", stopEarly: true, want: []string{"one"}, calls: 1},
		} {
			prefix := "owned/"
			if public {
				prefix = "public/"
			}
			t.Run(prefix+test.name, func(t *testing.T) {
				startID, nextID := uuid.New(), uuid.New()
				pages := [][]*workflowsv1.Workflow{
					{workflowsv1.Workflow_builder{Slug: "one"}.Build(), workflowsv1.Workflow_builder{Slug: "two"}.Build()},
					{workflowsv1.Workflow_builder{Slug: "three"}.Build(), workflowsv1.Workflow_builder{Slug: "four"}.Build()},
				}
				service := &pagingWorkflowService{}
				for i, items := range pages {
					var next *tileboxv1.Pagination
					if i == 0 {
						next = tileboxv1.Pagination_builder{StartingAfter: tileboxv1.NewUUID(nextID)}.Build()
					}
					service.owned = append(service.owned, workflowsv1.ListWorkflowsResponse_builder{Workflows: items, NextPage: next}.Build())
					service.public = append(service.public, workflowsv1.ListPublicWorkflowsResponse_builder{Workflows: items, NextPage: next}.Build())
				}
				client := workflowClient{service: service}
				list := client.List
				requests := &service.ownedPages
				if public {
					list = client.ListPublic
					requests = &service.publicPages
				}
				sequence := list(context.Background(), workflow.WithLimit(test.limit), workflow.WithCursor(workflow.NewCursor(startID)))
				require.Empty(t, *requests, "creating an iterator must not fetch pages")
				var slugs []string
				for item, err := range sequence {
					require.NoError(t, err)
					slugs = append(slugs, item.Slug)
					if test.stopEarly {
						break
					}
				}
				assert.Equal(t, test.want, slugs)
				require.Len(t, *requests, test.calls)
				assert.Equal(t, startID, (*requests)[0].GetStartingAfter().AsUUID())
				if test.limit > 0 {
					assert.Equal(t, test.limit, (*requests)[0].GetLimit())
				}
				if test.calls == 2 {
					assert.Equal(t, nextID, (*requests)[1].GetStartingAfter().AsUUID())
					if test.limit > 0 {
						assert.Equal(t, test.limit-2, (*requests)[1].GetLimit())
					}
				}
			})
		}
	}
}

func TestWorkflowClient_List_StopsOnLaterPageError(t *testing.T) {
	for _, public := range []bool{false, true} {
		service := &pagingWorkflowService{err: errors.New("second page failed")}
		next := tileboxv1.Pagination_builder{StartingAfter: tileboxv1.NewUUID(uuid.New())}.Build()
		service.owned = []*workflowsv1.ListWorkflowsResponse{workflowsv1.ListWorkflowsResponse_builder{NextPage: next}.Build()}
		service.public = []*workflowsv1.ListPublicWorkflowsResponse{workflowsv1.ListPublicWorkflowsResponse_builder{NextPage: next}.Build()}
		client := workflowClient{service: service}
		list := client.List
		if public {
			list = client.ListPublic
		}
		count := 0
		for item, err := range list(context.Background()) {
			assert.Nil(t, item)
			require.ErrorIs(t, err, service.err)
			count++
		}
		assert.Equal(t, 1, count, "yield an error once, even when the consumer continues")
	}
}

func TestWorkflowService_ListTransport(t *testing.T) {
	connectClient := &fakeWorkflowsConnectClient{}
	service := newWorkflowService(connectClient, noop.NewTracerProvider().Tracer("test"))
	client := workflowClient{service: service}
	cursorID := uuid.New()
	_, err := client.ListPage(context.Background(), workflow.WithLimit(37), workflow.WithCursor(workflow.NewCursor(cursorID)))
	require.NoError(t, err)
	assert.Equal(t, int64(37), connectClient.listWorkflowsRequest.GetPage().GetLimit())
	assert.Equal(t, cursorID, connectClient.listWorkflowsRequest.GetPage().GetStartingAfter().AsUUID())
	_, err = client.ListPublicPage(context.Background(), workflow.WithLimit(13), workflow.WithCursor(workflow.NewCursor(cursorID)))
	require.NoError(t, err)
	assert.Equal(t, int64(13), connectClient.listPublicWorkflowsRequest.GetPage().GetLimit())
	assert.Equal(t, cursorID, connectClient.listPublicWorkflowsRequest.GetPage().GetStartingAfter().AsUUID())
}

func TestWorkflowClient_UndeployRelease_OptionalIDTransport(t *testing.T) {
	releaseID := uuid.New()
	for _, test := range []struct {
		name    string
		options []workflow.UndeployOption
		wantID  *uuid.UUID
	}{
		{name: "all releases"},
		{name: "specific release", options: []workflow.UndeployOption{workflow.WithReleaseID(releaseID)}, wantID: &releaseID},
		{name: "explicit zero is not omission", options: []workflow.UndeployOption{workflow.WithReleaseID(uuid.Nil)}, wantID: &uuid.Nil},
	} {
		t.Run(test.name, func(t *testing.T) {
			transport := &fakeWorkflowsConnectClient{}
			client := workflowClient{service: newWorkflowService(transport, noop.NewTracerProvider().Tracer("test"))}
			_, err := client.UndeployRelease(context.Background(), "tilebox:example", []string{"dev", "prod"}, test.options...)
			require.NoError(t, err)
			request := transport.undeployWorkflowReleaseRequest
			require.NotNil(t, request)
			assert.Equal(t, "tilebox:example", request.GetWorkflowSlug())
			assert.Equal(t, []string{"dev", "prod"}, request.GetClusterSlugs())
			assert.Equal(t, test.wantID != nil, request.HasReleaseId())
			if test.wantID != nil {
				assert.Equal(t, *test.wantID, request.GetReleaseId().AsUUID())
			}
		})
	}
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

func TestWorkflowClient_Update(t *testing.T) {
	ctx := context.Background()
	service := &fakeWorkflowService{
		workflow: workflowsv1.Workflow_builder{
			Slug:        "agentic-workflow",
			Name:        "Agentic Workflow",
			Description: "Updated description",
		}.Build(),
	}
	client := workflowClient{service: service}

	updatedWorkflow, err := client.Update(ctx, "agentic-workflow", workflow.WithName("Agentic Workflow"), workflow.WithDescription("Updated description"))
	require.NoError(t, err)

	assert.Equal(t, "agentic-workflow", service.updateWorkflowSlug)
	require.NotNil(t, service.updateWorkflowName)
	require.NotNil(t, service.updateWorkflowDescription)
	assert.Equal(t, "Agentic Workflow", *service.updateWorkflowName)
	assert.Equal(t, "Updated description", *service.updateWorkflowDescription)
	assert.Equal(t, "agentic-workflow", updatedWorkflow.Slug)
	assert.Equal(t, "Agentic Workflow", updatedWorkflow.Name)
	assert.Equal(t, "Updated description", updatedWorkflow.Description)
}

func Test_workflowService_UpdateWorkflow_PreservesOptionalPresence(t *testing.T) {
	ctx := context.Background()
	connectClient := &fakeWorkflowsConnectClient{}
	service := newWorkflowService(connectClient, noop.NewTracerProvider().Tracer("test"))

	_, err := service.UpdateWorkflow(ctx, "agentic-workflow", nil, nil)
	require.NoError(t, err)
	require.NotNil(t, connectClient.updateWorkflowRequest)
	assert.Equal(t, "agentic-workflow", connectClient.updateWorkflowRequest.GetWorkflowSlug())
	assert.False(t, connectClient.updateWorkflowRequest.HasName())
	assert.False(t, connectClient.updateWorkflowRequest.HasDescription())

	connectClient.updateWorkflowRequest = nil
	emptyDescription := ""
	_, err = service.UpdateWorkflow(ctx, "agentic-workflow", nil, &emptyDescription)
	require.NoError(t, err)
	require.NotNil(t, connectClient.updateWorkflowRequest)
	assert.False(t, connectClient.updateWorkflowRequest.HasName())
	assert.True(t, connectClient.updateWorkflowRequest.HasDescription())
	assert.Empty(t, connectClient.updateWorkflowRequest.GetDescription())
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
			Clusters: []*workflowsv1.Cluster{
				workflowsv1.Cluster_builder{Slug: "dev", DisplayName: "Dev", Description: "Development cluster"}.Build(),
			},
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
	require.Len(t, release.Clusters, 1)
	assert.Equal(t, "dev", release.Clusters[0].Slug)
	assert.Equal(t, "Dev", release.Clusters[0].Name)
	assert.Equal(t, "Development cluster", release.Clusters[0].Description)
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

	deployment, err := client.UndeployRelease(ctx, "agentic-workflow", []string{"dev"}, workflow.WithReleaseID(releaseID))
	require.NoError(t, err)

	assert.Equal(t, "agentic-workflow", service.undeployWorkflowSlug)
	require.NotNil(t, service.undeployReleaseID)
	assert.Equal(t, releaseID, *service.undeployReleaseID)
	assert.Equal(t, []string{"dev"}, service.undeployClusterSlugs)
	require.NotNil(t, deployment.Release)
	assert.Equal(t, releaseID, deployment.Release.ID)
	require.Len(t, deployment.Clusters, 1)
	assert.Equal(t, "dev", deployment.Clusters[0].Slug)
	assert.Equal(t, "Dev", deployment.Clusters[0].Name)
}

func TestWorkflowClient_UndeployRelease_PreservesNilRelease(t *testing.T) {
	service := &fakeWorkflowService{
		undeployWorkflowReleaseResponse: workflowsv1.UndeployWorkflowReleaseResponse_builder{
			Clusters: []*workflowsv1.Cluster{workflowsv1.Cluster_builder{Slug: "dev"}.Build()},
		}.Build(),
	}

	deployment, err := (workflowClient{service: service}).UndeployRelease(context.Background(), "agentic-workflow", []string{"dev"})
	require.NoError(t, err)
	assert.Nil(t, service.undeployReleaseID)
	assert.Nil(t, deployment.Release)
	require.Len(t, deployment.Clusters, 1)
}

func TestProtoToCluster_MapsDeployedWorkflows(t *testing.T) {
	releaseID := uuid.New()
	cluster := protoToCluster(workflowsv1.Cluster_builder{
		Slug:        "dev",
		DisplayName: "Dev",
		Description: "Development cluster",
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
	assert.Equal(t, "Development cluster", cluster.Description)
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
	cluster                         *workflowsv1.Cluster
	workflow                        *workflowsv1.Workflow
	listWorkflowsResponse           *workflowsv1.ListWorkflowsResponse
	listPublicWorkflowsResponse     *workflowsv1.ListPublicWorkflowsResponse
	workflowRelease                 *workflowsv1.WorkflowRelease
	deployWorkflowReleaseResponse   *workflowsv1.DeployWorkflowReleaseResponse
	undeployWorkflowReleaseResponse *workflowsv1.UndeployWorkflowReleaseResponse
	err                             error

	createClusterName        string
	createClusterDescription string
	createClusterSlug        string
	updateClusterSlug        string
	updateClusterName        *string
	updateClusterDescription *string

	createName                string
	createDescription         string
	getSlug                   string
	updateWorkflowSlug        string
	updateWorkflowName        *string
	updateWorkflowDescription *string
	deleteWorkflowSlug        string

	publishWorkflowSlug string
	publishArtifactID   uuid.UUID
	publishContent      *ReleaseContent

	unpublishWorkflowSlug string
	unpublishReleaseID    uuid.UUID

	deployWorkflowSlug string
	deployReleaseID    uuid.UUID
	deployClusterSlugs []string

	undeployWorkflowSlug string
	undeployReleaseID    *uuid.UUID
	undeployClusterSlugs []string
}

type pagingWorkflowService struct {
	WorkflowService

	owned       []*workflowsv1.ListWorkflowsResponse
	public      []*workflowsv1.ListPublicWorkflowsResponse
	ownedPages  []*tileboxv1.Pagination
	publicPages []*tileboxv1.Pagination
	err         error
}

func (s *pagingWorkflowService) ListWorkflows(_ context.Context, page *tileboxv1.Pagination) (*workflowsv1.ListWorkflowsResponse, error) {
	s.ownedPages = append(s.ownedPages, page)
	if len(s.owned) == 0 {
		return nil, s.err
	}
	response := s.owned[0]
	s.owned = s.owned[1:]
	return response, nil
}

func (s *pagingWorkflowService) ListPublicWorkflows(_ context.Context, page *tileboxv1.Pagination) (*workflowsv1.ListPublicWorkflowsResponse, error) {
	s.publicPages = append(s.publicPages, page)
	if len(s.public) == 0 {
		return nil, s.err
	}
	response := s.public[0]
	s.public = s.public[1:]
	return response, nil
}

func (s *fakeWorkflowService) CreateCluster(_ context.Context, name, description, slug string) (*workflowsv1.Cluster, error) {
	s.createClusterName = name
	s.createClusterDescription = description
	s.createClusterSlug = slug
	return s.cluster, s.err
}

func (s *fakeWorkflowService) UpdateCluster(_ context.Context, slug string, name, description *string) (*workflowsv1.Cluster, error) {
	s.updateClusterSlug = slug
	s.updateClusterName = name
	s.updateClusterDescription = description
	return s.cluster, s.err
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

func (s *fakeWorkflowService) ListWorkflows(context.Context, *tileboxv1.Pagination) (*workflowsv1.ListWorkflowsResponse, error) {
	return s.listWorkflowsResponse, s.err
}

func (s *fakeWorkflowService) ListPublicWorkflows(context.Context, *tileboxv1.Pagination) (*workflowsv1.ListPublicWorkflowsResponse, error) {
	return s.listPublicWorkflowsResponse, s.err
}

func (s *fakeWorkflowService) GetWorkflow(_ context.Context, slug string) (*workflowsv1.Workflow, error) {
	s.getSlug = slug
	return s.workflow, s.err
}

func (s *fakeWorkflowService) UpdateWorkflow(_ context.Context, slug string, name, description *string) (*workflowsv1.Workflow, error) {
	s.updateWorkflowSlug = slug
	s.updateWorkflowName = name
	s.updateWorkflowDescription = description
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

func (s *fakeWorkflowService) UndeployWorkflowRelease(_ context.Context, workflowSlug string, releaseID *uuid.UUID, clusterSlugs []string) (*workflowsv1.UndeployWorkflowReleaseResponse, error) {
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
