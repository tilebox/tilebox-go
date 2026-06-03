package workflows // import "github.com/tilebox/tilebox-go/workflows/v1"

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	workflowsv1 "github.com/tilebox/tilebox-go/protogen/workflows/v1"
)

// Workflow represents a logical grouping for a set of tasks.
type Workflow struct {
	// Slug is the unique identifier of the workflow within the namespace.
	Slug string
	// Name is the human-readable name of the workflow.
	Name string
	// Description is the workflow description.
	Description string
	// Releases are the immutable releases published for this workflow.
	Releases []*WorkflowRelease
}

// WorkflowRelease represents an immutable release of a workflow.
type WorkflowRelease struct {
	// ID is the unique identifier of the workflow release.
	ID uuid.UUID
	// Artifact is the artifact associated with this release.
	Artifact *Artifact
	// Content describes the files and task implementations included in this release.
	Content *ReleaseContent
	// CreatedAt is the time the workflow release was created.
	CreatedAt time.Time
}

// Artifact represents the artifact associated with a workflow release.
type Artifact struct {
	// ID is the unique identifier of the artifact.
	ID uuid.UUID
	// Digest is the SHA-256 digest of the artifact.
	Digest string
}

// ReleaseContent represents the content included in a workflow release artifact.
type ReleaseContent struct {
	// Fingerprint is the SHA-256 fingerprint of the release content.
	Fingerprint string
	// Tasks are the task implementations included in this release content.
	Tasks []TaskIdentifier
	// Files are the files and directories included in the release artifact.
	Files []*Path
	// RunnerObjectPath is the Python module/object path to the runner instance.
	RunnerObjectPath string
	// CommandOverride is an optional custom command for starting a worker runtime.
	CommandOverride []string
}

// Path represents a file or directory path in release content.
type Path struct {
	// Path is the relative path.
	Path string
	// Directory reports whether this path is a directory.
	Directory bool
	// Children are nested paths when this path is a directory.
	Children []*Path
}

// WorkflowReleaseDeployment is the result of deploying or undeploying a workflow release.
type WorkflowReleaseDeployment struct {
	// Release is the workflow release that was deployed or undeployed.
	Release *WorkflowRelease
	// Clusters are the clusters affected by the deployment operation.
	Clusters []*Cluster
}

type workflowOptions struct {
	description string
}

// WorkflowOption configures workflow create requests.
type WorkflowOption func(*workflowOptions)

// WithDescription sets the workflow description.
func WithDescription(description string) WorkflowOption {
	return func(options *workflowOptions) {
		options.description = description
	}
}

func applyWorkflowOptions(options ...WorkflowOption) workflowOptions {
	var applied workflowOptions
	for _, option := range options {
		option(&applied)
	}
	return applied
}

type WorkflowClient interface {
	// Create creates a new workflow with the given name.
	Create(ctx context.Context, name string, options ...WorkflowOption) (*Workflow, error)

	// List returns all workflows.
	List(ctx context.Context) ([]*Workflow, error)

	// Get returns a workflow by its slug.
	Get(ctx context.Context, slug string) (*Workflow, error)

	// PublishRelease publishes a new immutable release for a workflow.
	PublishRelease(ctx context.Context, workflowSlug string, artifactID uuid.UUID, content *ReleaseContent) (*WorkflowRelease, error)

	// DeployRelease deploys a workflow release to clusters.
	DeployRelease(ctx context.Context, workflowSlug string, releaseID uuid.UUID, clusterSlugs []string) (*WorkflowReleaseDeployment, error)

	// UndeployRelease undeploys a workflow release from clusters.
	UndeployRelease(ctx context.Context, workflowSlug string, releaseID uuid.UUID, clusterSlugs []string) (*WorkflowReleaseDeployment, error)
}

var _ WorkflowClient = &workflowClient{}

type workflowClient struct {
	service WorkflowService
}

func (c workflowClient) Create(ctx context.Context, name string, options ...WorkflowOption) (*Workflow, error) {
	appliedOptions := applyWorkflowOptions(options...)
	response, err := c.service.CreateWorkflow(ctx, name, appliedOptions.description)
	if err != nil {
		return nil, err
	}

	return protoToWorkflow(response), nil
}

func (c workflowClient) List(ctx context.Context) ([]*Workflow, error) {
	response, err := c.service.ListWorkflows(ctx)
	if err != nil {
		return nil, err
	}

	workflows := make([]*Workflow, len(response.GetWorkflows()))
	for i, workflow := range response.GetWorkflows() {
		workflows[i] = protoToWorkflow(workflow)
	}

	return workflows, nil
}

func (c workflowClient) Get(ctx context.Context, slug string) (*Workflow, error) {
	response, err := c.service.GetWorkflow(ctx, slug)
	if err != nil {
		return nil, err
	}

	return protoToWorkflow(response), nil
}

func (c workflowClient) PublishRelease(ctx context.Context, workflowSlug string, artifactID uuid.UUID, content *ReleaseContent) (*WorkflowRelease, error) {
	response, err := c.service.PublishWorkflowRelease(ctx, workflowSlug, artifactID, content)
	if err != nil {
		return nil, err
	}

	return protoToWorkflowRelease(response), nil
}

func (c workflowClient) DeployRelease(ctx context.Context, workflowSlug string, releaseID uuid.UUID, clusterSlugs []string) (*WorkflowReleaseDeployment, error) {
	response, err := c.service.DeployWorkflowRelease(ctx, workflowSlug, releaseID, clusterSlugs)
	if err != nil {
		return nil, err
	}

	return protoToWorkflowReleaseDeployment(response.GetRelease(), response.GetClusters()), nil
}

func (c workflowClient) UndeployRelease(ctx context.Context, workflowSlug string, releaseID uuid.UUID, clusterSlugs []string) (*WorkflowReleaseDeployment, error) {
	response, err := c.service.UndeployWorkflowRelease(ctx, workflowSlug, releaseID, clusterSlugs)
	if err != nil {
		return nil, err
	}

	return protoToWorkflowReleaseDeployment(response.GetRelease(), response.GetClusters()), nil
}

func protoToWorkflow(workflow *workflowsv1.Workflow) *Workflow {
	if workflow == nil {
		return nil
	}

	releases := make([]*WorkflowRelease, len(workflow.GetReleases()))
	for i, release := range workflow.GetReleases() {
		releases[i] = protoToWorkflowRelease(release)
	}

	return &Workflow{
		Slug:        workflow.GetSlug(),
		Name:        workflow.GetName(),
		Description: workflow.GetDescription(),
		Releases:    releases,
	}
}

func protoToWorkflowRelease(release *workflowsv1.WorkflowRelease) *WorkflowRelease {
	if release == nil {
		return nil
	}

	var createdAt time.Time
	if release.GetCreatedAt() != nil {
		createdAt = release.GetCreatedAt().AsTime()
	}

	return &WorkflowRelease{
		ID:        protoIDToUUID(release.GetId()),
		Artifact:  protoToArtifact(release.GetArtifact()),
		Content:   protoToReleaseContent(release.GetContent()),
		CreatedAt: createdAt,
	}
}

func protoToArtifact(artifact *workflowsv1.Artifact) *Artifact {
	if artifact == nil {
		return nil
	}
	return &Artifact{
		ID:     protoIDToUUID(artifact.GetId()),
		Digest: artifact.GetDigest(),
	}
}

func protoToReleaseContent(content *workflowsv1.ReleaseContent) *ReleaseContent {
	if content == nil {
		return nil
	}

	tasks := make([]TaskIdentifier, len(content.GetTasks()))
	for i, task := range content.GetTasks() {
		tasks[i] = protoToTaskIdentifier(task)
	}

	return &ReleaseContent{
		Fingerprint:      content.GetFingerprint(),
		Tasks:            tasks,
		Files:            protoToPaths(content.GetFiles()),
		RunnerObjectPath: content.GetRunnerObjectPath(),
		CommandOverride:  content.GetCommandOverride(),
	}
}

func protoToPaths(paths []*workflowsv1.Path) []*Path {
	converted := make([]*Path, len(paths))
	for i, path := range paths {
		converted[i] = protoToPath(path)
	}
	return converted
}

func protoToPath(path *workflowsv1.Path) *Path {
	if path == nil {
		return nil
	}

	children := make([]*Path, len(path.GetChildren()))
	for i, child := range path.GetChildren() {
		children[i] = protoToPath(child)
	}

	return &Path{
		Path:      path.GetPath(),
		Directory: path.GetDirectory(),
		Children:  children,
	}
}

func protoToWorkflowReleaseDeployment(release *workflowsv1.WorkflowRelease, clusters []*workflowsv1.Cluster) *WorkflowReleaseDeployment {
	deployedClusters := make([]*Cluster, len(clusters))
	for i, cluster := range clusters {
		deployedClusters[i] = protoToCluster(cluster)
	}

	return &WorkflowReleaseDeployment{
		Release:  protoToWorkflowRelease(release),
		Clusters: deployedClusters,
	}
}

func protoToTaskIdentifier(identifier *workflowsv1.TaskIdentifier) TaskIdentifier {
	if identifier == nil {
		return nil
	}
	return NewTaskIdentifier(identifier.GetName(), identifier.GetVersion())
}

func taskIdentifiersToProto(tasks []TaskIdentifier) ([]*workflowsv1.TaskIdentifier, error) {
	protoTasks := make([]*workflowsv1.TaskIdentifier, len(tasks))
	for i, task := range tasks {
		if task == nil {
			return nil, fmt.Errorf("task identifier at index %d is nil", i)
		}
		if err := ValidateIdentifier(task); err != nil {
			return nil, err
		}
		protoTasks[i] = workflowsv1.TaskIdentifier_builder{
			Name:    task.Name(),
			Version: task.Version(),
		}.Build()
	}
	return protoTasks, nil
}

func releaseContentToProto(content *ReleaseContent) (*workflowsv1.ReleaseContent, error) {
	if content == nil {
		return nil, fmt.Errorf("release content is nil")
	}

	tasks, err := taskIdentifiersToProto(content.Tasks)
	if err != nil {
		return nil, err
	}

	return workflowsv1.ReleaseContent_builder{
		Fingerprint:      content.Fingerprint,
		Tasks:            tasks,
		Files:            pathsToProto(content.Files),
		RunnerObjectPath: content.RunnerObjectPath,
		CommandOverride:  content.CommandOverride,
	}.Build(), nil
}

func pathsToProto(paths []*Path) []*workflowsv1.Path {
	converted := make([]*workflowsv1.Path, len(paths))
	for i, path := range paths {
		converted[i] = pathToProto(path)
	}
	return converted
}

func pathToProto(path *Path) *workflowsv1.Path {
	if path == nil {
		return nil
	}

	children := make([]*workflowsv1.Path, len(path.Children))
	for i, child := range path.Children {
		children[i] = pathToProto(child)
	}

	return workflowsv1.Path_builder{
		Path:      path.Path,
		Directory: path.Directory,
		Children:  children,
	}.Build()
}
