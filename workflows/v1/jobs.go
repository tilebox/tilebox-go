package workflows // import "github.com/tilebox/tilebox-go/workflows/v1"

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"iter"
	"time"

	"github.com/google/uuid"
	workflowsv1 "github.com/tilebox/tilebox-go/protogen/go/workflows/v1"
	"google.golang.org/protobuf/proto"
)

// Job represents a Tilebox Workflows job.
//
// Documentation: https://docs.tilebox.com/workflows/concepts/jobs
type Job struct {
	// ID is the unique identifier of the job.
	ID uuid.UUID
	// Name is the name of the job.
	Name string
	// Canceled indicates whether the job has been canceled.
	Canceled bool
	// State is the current state of the job.
	State JobState
	// SubmittedAt is the time the job was submitted.
	SubmittedAt time.Time
	// StartedAt is the time the job started running.
	StartedAt time.Time
	// TaskSummaries is the task summaries of the job.
	TaskSummaries []*TaskSummary
	// AutomationID is the ID of the automation that submitted the job.
	AutomationID uuid.UUID
}

// JobState is the state of a Job.
type JobState int32

// JobState values.
const (
	_            JobState = iota
	JobQueued             // The job is queued and waiting to be run.
	JobStarted            // At least one task of the job has been started.
	JobCompleted          // All tasks of the job have been completed.
)

// TaskSummary is a summary of a task.
type TaskSummary struct {
	// ID is the unique identifier of the task.
	ID uuid.UUID
	// Display is the display message of the task.
	Display string
	// State is the state of the task.
	State TaskState
	// ParentID is the ID of the parent task.
	ParentID uuid.UUID
	// DependsOn is the list of IDs of the tasks that this task depends on.
	DependsOn []uuid.UUID
	// StartedAt is the time the task started.
	StartedAt time.Time
	// StoppedAt is the time the task stopped.
	StoppedAt time.Time
}

// TaskState is the state of a Task.
type TaskState int32

// TaskState values.
const (
	_             TaskState = iota
	TaskQueued              // The task is queued and waiting to be run.
	TaskRunning             // The task is currently running on some task runner.
	TaskComputed            // The task has been computed and the output is available.
	TaskFailed              // The task has failed.
	TaskCancelled           // The task has been cancelled due to user request.
)

type JobClient interface {
	Submit(ctx context.Context, jobName string, clusterSlug string, maxRetries int, tasks ...Task) (*Job, error)
	Get(ctx context.Context, jobID uuid.UUID) (*Job, error)
	Retry(ctx context.Context, jobID uuid.UUID) (int64, error)
	Cancel(ctx context.Context, jobID uuid.UUID) error
	List(ctx context.Context, idInterval *workflowsv1.IDInterval) iter.Seq2[*Job, error]
	ValidateJob(jobName string, clusterSlug string, maxRetries int, tasks ...Task) (*workflowsv1.SubmitJobRequest, error)
}

var _ JobClient = &jobClient{}

type jobClient struct {
	service JobService
}

// Submit submits a job to the specified cluster.
//
// Documentation: https://docs.tilebox.com/workflows/concepts/jobs#submission
func (c jobClient) Submit(ctx context.Context, jobName string, clusterSlug string, maxRetries int, tasks ...Task) (*Job, error) {
	jobRequest, err := c.ValidateJob(jobName, clusterSlug, maxRetries, tasks...)
	if err != nil {
		return nil, err
	}

	response, err := c.service.SubmitJob(ctx, jobRequest)
	if err != nil {
		return nil, err
	}

	job, err := protoToJob(response)
	if err != nil {
		return nil, fmt.Errorf("failed to convert job from response: %w", err)
	}

	return job, nil
}

// Get returns a job by its ID.
func (c jobClient) Get(ctx context.Context, jobID uuid.UUID) (*Job, error) {
	response, err := c.service.GetJob(ctx, jobID)
	if err != nil {
		return nil, err
	}

	job, err := protoToJob(response)
	if err != nil {
		return nil, fmt.Errorf("failed to convert job from response: %w", err)
	}

	return job, nil
}

// Retry retries a job.
//
// Documentation: https://docs.tilebox.com/workflows/concepts/jobs#retry-handling
func (c jobClient) Retry(ctx context.Context, jobID uuid.UUID) (int64, error) {
	response, err := c.service.RetryJob(ctx, jobID)
	if err != nil {
		return 0, err
	}

	return response.GetNumTasksRescheduled(), nil
}

// Cancel cancels a job.
//
// Documentation: https://docs.tilebox.com/workflows/concepts/jobs#cancellation
func (c jobClient) Cancel(ctx context.Context, jobID uuid.UUID) error {
	return c.service.CancelJob(ctx, jobID)
}

// List returns a list of all jobs within the given interval.
//
// The jobs are loaded in a lazy manner, and returned as a sequence.
// The output sequence can be transformed into a slice using Collect.
func (c jobClient) List(ctx context.Context, idInterval *workflowsv1.IDInterval) iter.Seq2[*Job, error] {
	return func(yield func(*Job, error) bool) {
		var page *workflowsv1.Pagination // nil for the first request

		for {
			jobsMessage, err := c.service.ListJobs(ctx, idInterval, page)
			if err != nil {
				yield(nil, err)
				return
			}

			for _, jobMessage := range jobsMessage.GetJobs() {
				job, err := protoToJob(jobMessage)
				if err != nil {
					yield(nil, fmt.Errorf("failed to convert job from response: %w", err))
					return
				}

				if !yield(job, nil) {
					return
				}
			}

			page = jobsMessage.GetNextPage()
			if page == nil {
				break
			}
		}
	}
}

func (c jobClient) ValidateJob(jobName string, clusterSlug string, maxRetries int, tasks ...Task) (*workflowsv1.SubmitJobRequest, error) {
	if len(tasks) == 0 {
		return nil, errors.New("no tasks to submit")
	}

	rootTasks := make([]*workflowsv1.TaskSubmission, 0)

	for _, task := range tasks {
		var subtaskInput []byte
		var err error

		identifier := identifierFromTask(task)
		err = ValidateIdentifier(identifier)
		if err != nil {
			return nil, fmt.Errorf("task has invalid task identifier: %w", err)
		}

		taskProto, isProtobuf := task.(proto.Message)
		if isProtobuf {
			subtaskInput, err = proto.Marshal(taskProto)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal protobuf task: %w", err)
			}
		} else {
			subtaskInput, err = json.Marshal(task)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal task: %w", err)
			}
		}

		rootTasks = append(rootTasks, &workflowsv1.TaskSubmission{
			ClusterSlug: clusterSlug,
			Identifier: &workflowsv1.TaskIdentifier{
				Name:    identifier.Name(),
				Version: identifier.Version(),
			},
			Input:      subtaskInput,
			Display:    identifier.Display(),
			MaxRetries: int64(maxRetries),
		})
	}

	return &workflowsv1.SubmitJobRequest{
		Tasks:   rootTasks,
		JobName: jobName,
	}, nil
}

func protoToJob(job *workflowsv1.Job) (*Job, error) {
	id, err := protoToUUID(job.GetId())
	if err != nil {
		return nil, fmt.Errorf("failed to parse job id: %w", err)
	}

	automationID, err := protoToUUID(job.GetAutomationId())
	if err != nil {
		return nil, fmt.Errorf("failed to parse automation id: %w", err)
	}

	taskSummaries := make([]*TaskSummary, len(job.GetTaskSummaries()))
	for i, taskSummary := range job.GetTaskSummaries() {
		taskSummaries[i], err = protoToTaskSummary(taskSummary)
		if err != nil {
			return nil, fmt.Errorf("failed to convert task summary from response: %w", err)
		}
	}

	return &Job{
		ID:            id,
		Name:          job.GetName(),
		Canceled:      job.GetCanceled(),
		State:         JobState(job.GetState()),
		SubmittedAt:   job.GetSubmittedAt().AsTime(),
		StartedAt:     job.GetStartedAt().AsTime(),
		TaskSummaries: taskSummaries,
		AutomationID:  automationID,
	}, nil
}

func protoToTaskSummary(t *workflowsv1.TaskSummary) (*TaskSummary, error) {
	taskSummaryID, err := protoToUUID(t.GetId())
	if err != nil {
		return nil, fmt.Errorf("failed to parse task summary id: %w", err)
	}

	parentID, err := protoToUUID(t.GetParentId())
	if err != nil {
		return nil, fmt.Errorf("failed to parse task summary parent id: %w", err)
	}

	dependsOn := make([]uuid.UUID, len(t.GetDependsOn()))
	for j, dependsOnID := range t.GetDependsOn() {
		dependsOn[j], err = protoToUUID(dependsOnID)
		if err != nil {
			return nil, fmt.Errorf("failed to parse task summary depends on id: %w", err)
		}
	}

	return &TaskSummary{
		ID:        taskSummaryID,
		Display:   t.GetDisplay(),
		State:     TaskState(t.GetState()),
		ParentID:  parentID,
		DependsOn: dependsOn,
		StartedAt: t.GetStartedAt().AsTime(),
		StoppedAt: t.GetStoppedAt().AsTime(),
	}, nil
}

// Collect converts any sequence into a slice.
//
// It returns an error if any of the elements in the sequence has a non-nil error.
func Collect[K any](seq iter.Seq2[K, error]) ([]K, error) {
	s := make([]K, 0)

	for k, err := range seq {
		if err != nil {
			return nil, err
		}
		s = append(s, k)
	}
	return s, nil
}
