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
	"github.com/tilebox/tilebox-go/workflows/v1/job"
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
	// DependsOn is the list of task IDs that this task depends on.
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
	// Submit submits a job to a cluster.
	//
	// Options:
	//   - job.WithMaxRetries: sets the maximum number of times a job can be automatically retried. Defaults to 0.
	//   - job.WithClusterSlug: sets the cluster slug of the cluster where the job will be executed. Defaults to the default cluster.
	//
	// Documentation: https://docs.tilebox.com/workflows/concepts/jobs#submission
	Submit(ctx context.Context, jobName string, tasks []Task, options ...job.SubmitOption) (*Job, error)

	// Get returns a job by its ID.
	Get(ctx context.Context, jobID uuid.UUID) (*Job, error)

	// Retry retries a job.
	//
	// Returns the number of rescheduled tasks.
	//
	// Documentation: https://docs.tilebox.com/workflows/concepts/jobs#retry-handling
	Retry(ctx context.Context, jobID uuid.UUID) (int64, error)

	// Cancel cancels a job.
	//
	// Documentation: https://docs.tilebox.com/workflows/concepts/jobs#cancellation
	Cancel(ctx context.Context, jobID uuid.UUID) error

	// Query returns a list of all jobs within the given interval.
	//
	// Options:
	//   - job.WithTemporalExtent: specifies the time or ID interval for which jobs should be queried (Required)
	//   - job.WithAutomationID: specifies the automation ID to filter jobs by. Only jobs submitted by the specified
	//  automation will be returned. (Optional)
	//
	// The jobs are lazily loaded and returned as a sequence.
	// The output sequence can be transformed into a slice using Collect.
	Query(ctx context.Context, options ...job.QueryOption) iter.Seq2[*Job, error]
}

var _ JobClient = &jobClient{}

type jobClient struct {
	service JobService
}

func (c jobClient) Submit(ctx context.Context, jobName string, tasks []Task, options ...job.SubmitOption) (*Job, error) {
	opts := &job.SubmitOptions{}
	for _, option := range options {
		option(opts)
	}

	jobRequest, err := validateJob(jobName, opts.ClusterSlug, opts.MaxRetries, tasks...)
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

func (c jobClient) Retry(ctx context.Context, jobID uuid.UUID) (int64, error) {
	response, err := c.service.RetryJob(ctx, jobID)
	if err != nil {
		return 0, err
	}

	return response.GetNumTasksRescheduled(), nil
}

func (c jobClient) Cancel(ctx context.Context, jobID uuid.UUID) error {
	return c.service.CancelJob(ctx, jobID)
}

func (c jobClient) Query(ctx context.Context, options ...job.QueryOption) iter.Seq2[*Job, error] {
	opts := &job.QueryOptions{}
	for _, option := range options {
		option(opts)
	}

	if opts.TemporalExtent == nil {
		return func(yield func(*Job, error) bool) {
			// right now we return an error, in the future we might want to support queries without a temporal extent
			yield(nil, errors.New("temporal extent is required"))
		}
	}

	return func(yield func(*Job, error) bool) {
		var page *workflowsv1.Pagination // nil for the first request

		// we already validated that TemporalExtent is not nil
		timeInterval := opts.TemporalExtent.ToProtoWorkflowTimeInterval()
		idInterval := opts.TemporalExtent.ToProtoIDInterval()

		if timeInterval == nil && idInterval == nil {
			yield(nil, errors.New("invalid temporal extent"))
			return
		}

		filters := &workflowsv1.QueryFilters{}
		if timeInterval != nil {
			filters.TemporalExtent = &workflowsv1.QueryFilters_TimeInterval{TimeInterval: timeInterval}
		} else {
			filters.TemporalExtent = &workflowsv1.QueryFilters_IdInterval{IdInterval: idInterval}
		}

		for {
			jobsMessage, err := c.service.QueryJobs(ctx, filters, page)
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

func validateJob(jobName string, clusterSlug string, maxRetries int64, tasks ...Task) (*workflowsv1.SubmitJobRequest, error) {
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
			MaxRetries: maxRetries,
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
