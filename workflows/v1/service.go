package workflows // import "github.com/tilebox/tilebox-go/workflows/v1"

import (
	"context"
	"fmt"
	"time"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/tilebox/tilebox-go/observability"
	workflowsv1 "github.com/tilebox/tilebox-go/protogen/go/workflows/v1"
	"github.com/tilebox/tilebox-go/protogen/go/workflows/v1/workflowsv1connect"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type _automationService interface {
	CreateStorageLocation(ctx context.Context, location string, storageType workflowsv1.StorageType) (*workflowsv1.StorageLocation, error)
	GetStorageLocation(ctx context.Context, storageLocationID uuid.UUID) (*workflowsv1.StorageLocation, error)
	DeleteStorageLocation(ctx context.Context, storageLocationID uuid.UUID) error
	ListStorageLocations(ctx context.Context) (*workflowsv1.StorageLocations, error)

	CreateAutomation(ctx context.Context, automation *workflowsv1.AutomationPrototype) (*workflowsv1.AutomationPrototype, error)
	GetAutomation(ctx context.Context, automationID uuid.UUID) (*workflowsv1.AutomationPrototype, error)
	UpdateAutomation(ctx context.Context, automation *workflowsv1.AutomationPrototype) (*workflowsv1.AutomationPrototype, error)
	DeleteAutomation(ctx context.Context, automationID uuid.UUID, cancelJobs bool) error
	ListAutomations(ctx context.Context) (*workflowsv1.Automations, error)
}

var _ _automationService = &automationService{}

type automationService struct {
	automationClient workflowsv1connect.AutomationServiceClient
	tracer           trace.Tracer
}

func (s *automationService) CreateStorageLocation(ctx context.Context, location string, storageType workflowsv1.StorageType) (*workflowsv1.StorageLocation, error) {
	return observability.WithSpanResult(ctx, s.tracer, "workflows/storage_locations/create", func(ctx context.Context) (*workflowsv1.StorageLocation, error) {
		res, err := s.automationClient.CreateStorageLocation(ctx, connect.NewRequest(workflowsv1.StorageLocation_builder{
			Location: location,
			Type:     storageType,
		}.Build()))
		if err != nil {
			return nil, fmt.Errorf("failed to create storage location: %w", err)
		}

		return res.Msg, nil
	})
}

func (s *automationService) GetStorageLocation(ctx context.Context, storageLocationID uuid.UUID) (*workflowsv1.StorageLocation, error) {
	return observability.WithSpanResult(ctx, s.tracer, "workflows/storage_locations/get", func(ctx context.Context) (*workflowsv1.StorageLocation, error) {
		res, err := s.automationClient.GetStorageLocation(ctx, connect.NewRequest(
			uuidToProtobuf(storageLocationID),
		))
		if err != nil {
			return nil, fmt.Errorf("failed to get storage location: %w", err)
		}

		return res.Msg, nil
	})
}

func (s *automationService) DeleteStorageLocation(ctx context.Context, storageLocationID uuid.UUID) error {
	return observability.WithSpan(ctx, s.tracer, "workflows/storage_locations/delete", func(ctx context.Context) error {
		_, err := s.automationClient.DeleteStorageLocation(ctx, connect.NewRequest(
			uuidToProtobuf(storageLocationID),
		))
		if err != nil {
			return fmt.Errorf("failed to delete storage location: %w", err)
		}

		return nil
	})
}

func (s *automationService) ListStorageLocations(ctx context.Context) (*workflowsv1.StorageLocations, error) {
	return observability.WithSpanResult(ctx, s.tracer, "workflows/storage_locations/list", func(ctx context.Context) (*workflowsv1.StorageLocations, error) {
		res, err := s.automationClient.ListStorageLocations(ctx, connect.NewRequest(&emptypb.Empty{}))
		if err != nil {
			return nil, fmt.Errorf("failed to list storage locations: %w", err)
		}

		return res.Msg, nil
	})
}

func (s *automationService) CreateAutomation(ctx context.Context, automation *workflowsv1.AutomationPrototype) (*workflowsv1.AutomationPrototype, error) {
	return observability.WithSpanResult(ctx, s.tracer, "workflows/automations/create", func(ctx context.Context) (*workflowsv1.AutomationPrototype, error) {
		res, err := s.automationClient.CreateAutomation(ctx, connect.NewRequest(automation))
		if err != nil {
			return nil, fmt.Errorf("failed to create automation: %w", err)
		}

		return res.Msg, nil
	})
}

func (s *automationService) GetAutomation(ctx context.Context, automationID uuid.UUID) (*workflowsv1.AutomationPrototype, error) {
	return observability.WithSpanResult(ctx, s.tracer, "workflows/automations/get", func(ctx context.Context) (*workflowsv1.AutomationPrototype, error) {
		res, err := s.automationClient.GetAutomation(ctx, connect.NewRequest(
			uuidToProtobuf(automationID),
		))
		if err != nil {
			return nil, fmt.Errorf("failed to get automation: %w", err)
		}

		return res.Msg, nil
	})
}

func (s *automationService) UpdateAutomation(ctx context.Context, automation *workflowsv1.AutomationPrototype) (*workflowsv1.AutomationPrototype, error) {
	return observability.WithSpanResult(ctx, s.tracer, "workflows/automations/update", func(ctx context.Context) (*workflowsv1.AutomationPrototype, error) {
		res, err := s.automationClient.UpdateAutomation(ctx, connect.NewRequest(automation))
		if err != nil {
			return nil, fmt.Errorf("failed to update automation: %w", err)
		}

		return res.Msg, nil
	})
}

func (s *automationService) DeleteAutomation(ctx context.Context, automationID uuid.UUID, cancelJobs bool) error {
	return observability.WithSpan(ctx, s.tracer, "workflows/automations/delete", func(ctx context.Context) error {
		_, err := s.automationClient.DeleteAutomation(ctx, connect.NewRequest(workflowsv1.DeleteAutomationRequest_builder{
			AutomationId: uuidToProtobuf(automationID),
			CancelJobs:   cancelJobs,
		}.Build()))
		if err != nil {
			return fmt.Errorf("failed to delete automation: %w", err)
		}

		return nil
	})
}

func (s *automationService) ListAutomations(ctx context.Context) (*workflowsv1.Automations, error) {
	return observability.WithSpanResult(ctx, s.tracer, "workflows/automations/list", func(ctx context.Context) (*workflowsv1.Automations, error) {
		res, err := s.automationClient.ListAutomations(ctx, connect.NewRequest(&emptypb.Empty{}))
		if err != nil {
			return nil, fmt.Errorf("failed to list automations: %w", err)
		}

		return res.Msg, nil
	})
}

type JobService interface {
	SubmitJob(ctx context.Context, req *workflowsv1.SubmitJobRequest) (*workflowsv1.Job, error)
	GetJob(ctx context.Context, jobID uuid.UUID) (*workflowsv1.Job, error)
	RetryJob(ctx context.Context, jobID uuid.UUID) (*workflowsv1.RetryJobResponse, error)
	CancelJob(ctx context.Context, jobID uuid.UUID) error
	QueryJobs(ctx context.Context, filters *workflowsv1.QueryFilters, page *workflowsv1.Pagination) (*workflowsv1.QueryJobsResponse, error)
}

var _ JobService = &jobService{}

type jobService struct {
	jobClient workflowsv1connect.JobServiceClient
	tracer    trace.Tracer
}

func newJobService(jobClient workflowsv1connect.JobServiceClient, tracer trace.Tracer) JobService {
	return &jobService{
		jobClient: jobClient,
		tracer:    tracer,
	}
}

func (s *jobService) SubmitJob(ctx context.Context, req *workflowsv1.SubmitJobRequest) (*workflowsv1.Job, error) {
	return observability.WithSpanResult(ctx, s.tracer, "workflows/jobs/submit", func(ctx context.Context) (*workflowsv1.Job, error) {
		traceParent := observability.GetTraceParentOfCurrentSpan(ctx)
		req.SetTraceParent(traceParent)

		res, err := s.jobClient.SubmitJob(ctx, connect.NewRequest(req))
		if err != nil {
			return nil, fmt.Errorf("failed to submit job: %w", err)
		}

		return res.Msg, nil
	})
}

func (s *jobService) GetJob(ctx context.Context, jobID uuid.UUID) (*workflowsv1.Job, error) {
	return observability.WithSpanResult(ctx, s.tracer, "workflows/jobs/get", func(ctx context.Context) (*workflowsv1.Job, error) {
		res, err := s.jobClient.GetJob(ctx, connect.NewRequest(workflowsv1.GetJobRequest_builder{
			JobId: uuidToProtobuf(jobID),
		}.Build()))
		if err != nil {
			return nil, fmt.Errorf("failed to get job: %w", err)
		}

		return res.Msg, nil
	})
}

func (s *jobService) RetryJob(ctx context.Context, jobID uuid.UUID) (*workflowsv1.RetryJobResponse, error) {
	return observability.WithSpanResult(ctx, s.tracer, "workflows/jobs/retry", func(ctx context.Context) (*workflowsv1.RetryJobResponse, error) {
		res, err := s.jobClient.RetryJob(ctx, connect.NewRequest(workflowsv1.RetryJobRequest_builder{
			JobId: uuidToProtobuf(jobID),
		}.Build()))
		if err != nil {
			return nil, fmt.Errorf("failed to retry job: %w", err)
		}

		return res.Msg, nil
	})
}

func (s *jobService) CancelJob(ctx context.Context, jobID uuid.UUID) error {
	return observability.WithSpan(ctx, s.tracer, "workflows/jobs/cancel", func(ctx context.Context) error {
		_, err := s.jobClient.CancelJob(ctx, connect.NewRequest(workflowsv1.CancelJobRequest_builder{
			JobId: uuidToProtobuf(jobID),
		}.Build()))
		if err != nil {
			return fmt.Errorf("failed to cancel job: %w", err)
		}

		return nil
	})
}

func (s *jobService) QueryJobs(ctx context.Context, filters *workflowsv1.QueryFilters, page *workflowsv1.Pagination) (*workflowsv1.QueryJobsResponse, error) {
	return observability.WithSpanResult(ctx, s.tracer, "workflows/jobs/query", func(ctx context.Context) (*workflowsv1.QueryJobsResponse, error) {
		res, err := s.jobClient.QueryJobs(ctx, connect.NewRequest(workflowsv1.QueryJobsRequest_builder{
			Filters: filters,
			Page:    page,
		}.Build()))
		if err != nil {
			return nil, fmt.Errorf("failed to query jobs: %w", err)
		}

		return res.Msg, nil
	})
}

type TaskService interface {
	NextTask(ctx context.Context, computedTask *workflowsv1.ComputedTask, nextTaskToRun *workflowsv1.NextTaskToRun) (*workflowsv1.NextTaskResponse, error)
	TaskFailed(ctx context.Context, taskID uuid.UUID, display string, cancelJob bool) (*workflowsv1.TaskStateResponse, error)
	ExtendTaskLease(ctx context.Context, taskID uuid.UUID, requestedLease time.Duration) (*workflowsv1.TaskLease, error)
}

var _ TaskService = &taskService{}

type taskService struct {
	taskClient workflowsv1connect.TaskServiceClient
	tracer     trace.Tracer
}

func newTaskService(taskClient workflowsv1connect.TaskServiceClient, tracer trace.Tracer) TaskService {
	return &taskService{
		taskClient: taskClient,
		tracer:     tracer,
	}
}

func (s *taskService) NextTask(ctx context.Context, computedTask *workflowsv1.ComputedTask, nextTaskToRun *workflowsv1.NextTaskToRun) (*workflowsv1.NextTaskResponse, error) {
	return observability.WithSpanResult(ctx, s.tracer, "workflows/tasks/next", func(ctx context.Context) (*workflowsv1.NextTaskResponse, error) {
		res, err := s.taskClient.NextTask(ctx, connect.NewRequest(workflowsv1.NextTaskRequest_builder{
			ComputedTask:  computedTask,
			NextTaskToRun: nextTaskToRun,
		}.Build()))
		if err != nil {
			return nil, fmt.Errorf("failed to get next task: %w", err)
		}

		return res.Msg, nil
	})
}

func (s *taskService) TaskFailed(ctx context.Context, taskID uuid.UUID, display string, cancelJob bool) (*workflowsv1.TaskStateResponse, error) {
	return observability.WithSpanResult(ctx, s.tracer, "workflows/tasks/failed", func(ctx context.Context) (*workflowsv1.TaskStateResponse, error) {
		res, err := s.taskClient.TaskFailed(ctx, connect.NewRequest(workflowsv1.TaskFailedRequest_builder{
			TaskId:    uuidToProtobuf(taskID),
			Display:   display,
			CancelJob: cancelJob,
		}.Build()))
		if err != nil {
			return nil, fmt.Errorf("failed to mark task as failed: %w", err)
		}

		return res.Msg, nil
	})
}

func (s *taskService) ExtendTaskLease(ctx context.Context, taskID uuid.UUID, requestedLease time.Duration) (*workflowsv1.TaskLease, error) {
	return observability.WithSpanResult(ctx, s.tracer, "workflows/tasks/lease", func(ctx context.Context) (*workflowsv1.TaskLease, error) {
		res, err := s.taskClient.ExtendTaskLease(ctx, connect.NewRequest(workflowsv1.TaskLeaseRequest_builder{
			TaskId:         uuidToProtobuf(taskID),
			RequestedLease: durationpb.New(requestedLease),
		}.Build()))
		if err != nil {
			return nil, fmt.Errorf("failed to extend task lease: %w", err)
		}

		return res.Msg, nil
	})
}

type WorkflowService interface {
	CreateCluster(ctx context.Context, name string) (*workflowsv1.Cluster, error)
	GetCluster(ctx context.Context, slug string) (*workflowsv1.Cluster, error)
	DeleteCluster(ctx context.Context, slug string) error
	ListClusters(ctx context.Context) (*workflowsv1.ListClustersResponse, error)
}

var _ WorkflowService = &workflowService{}

type workflowService struct {
	workflowClient workflowsv1connect.WorkflowsServiceClient
	tracer         trace.Tracer
}

func newWorkflowService(workflowClient workflowsv1connect.WorkflowsServiceClient, tracer trace.Tracer) WorkflowService {
	return &workflowService{
		workflowClient: workflowClient,
		tracer:         tracer,
	}
}

func (s *workflowService) CreateCluster(ctx context.Context, name string) (*workflowsv1.Cluster, error) {
	return observability.WithSpanResult(ctx, s.tracer, "workflows/clusters/create", func(ctx context.Context) (*workflowsv1.Cluster, error) {
		res, err := s.workflowClient.CreateCluster(ctx, connect.NewRequest(workflowsv1.CreateClusterRequest_builder{
			Name: name,
		}.Build()))
		if err != nil {
			return nil, fmt.Errorf("failed to create cluster: %w", err)
		}

		return res.Msg, nil
	})
}

func (s *workflowService) GetCluster(ctx context.Context, slug string) (*workflowsv1.Cluster, error) {
	return observability.WithSpanResult(ctx, s.tracer, "workflows/clusters/get", func(ctx context.Context) (*workflowsv1.Cluster, error) {
		res, err := s.workflowClient.GetCluster(ctx, connect.NewRequest(workflowsv1.GetClusterRequest_builder{
			ClusterSlug: slug,
		}.Build()))
		if err != nil {
			return nil, fmt.Errorf("failed to get cluster: %w", err)
		}

		return res.Msg, nil
	})
}

func (s *workflowService) DeleteCluster(ctx context.Context, slug string) error {
	return observability.WithSpan(ctx, s.tracer, "workflows/clusters/delete", func(ctx context.Context) error {
		_, err := s.workflowClient.DeleteCluster(ctx, connect.NewRequest(workflowsv1.DeleteClusterRequest_builder{
			ClusterSlug: slug,
		}.Build()))
		if err != nil {
			return fmt.Errorf("failed to delete cluster: %w", err)
		}

		return nil
	})
}

func (s *workflowService) ListClusters(ctx context.Context) (*workflowsv1.ListClustersResponse, error) {
	return observability.WithSpanResult(ctx, s.tracer, "workflows/clusters/list", func(ctx context.Context) (*workflowsv1.ListClustersResponse, error) {
		res, err := s.workflowClient.ListClusters(ctx, connect.NewRequest(&workflowsv1.ListClustersRequest{}))
		if err != nil {
			return nil, fmt.Errorf("failed to list clusters: %w", err)
		}

		return res.Msg, nil
	})
}
