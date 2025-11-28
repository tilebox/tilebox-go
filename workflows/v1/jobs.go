package workflows // import "github.com/tilebox/tilebox-go/workflows/v1"

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"iter"
	"time"

	"github.com/google/uuid"
	tileboxv1 "github.com/tilebox/tilebox-go/protogen/tilebox/v1"
	workflowsv1 "github.com/tilebox/tilebox-go/protogen/workflows/v1"
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
	// State is the current state of the job.
	State job.State
	// SubmittedAt is the time the job was submitted.
	SubmittedAt time.Time
	// TaskSummaries is the task summaries of the job.
	TaskSummaries []*TaskSummary
	// AutomationID is the ID of the automation that submitted the job.
	AutomationID uuid.UUID
	// Progress is a list of progress indicators for the job.
	Progress []*ProgressIndicator
	// ExecutionStats contains execution statistics for the job.
	ExecutionStats *ExecutionStats
}

// TaskSummary is a summary of a task.
type TaskSummary struct {
	// ID is the unique identifier of the task.
	ID uuid.UUID
	// Display is the label message of the task.
	Display string
	// State is the state of the task.
	State TaskState
	// ParentID is the ID of the parent task.
	ParentID uuid.UUID
	// StartedAt is the time the task started.
	StartedAt time.Time
	// StoppedAt is the time the task stopped.
	StoppedAt time.Time
}

// TaskState is the state of a Task.
type TaskState int32

type ProgressIndicator struct {
	Label string
	Total uint64
	Done  uint64
}

// ExecutionStats contains execution statistics for a job.
type ExecutionStats struct {
	// FirstTaskStartedAt is the time the first task of the job was started.
	FirstTaskStartedAt time.Time
	// LastTaskStoppedAt is the time the last task of the job was stopped.
	LastTaskStoppedAt time.Time
	// ComputeTime is the total compute time of the job, as sum of all task compute times.
	ComputeTime time.Duration
	// ElapsedTime is the elapsed time of the job (wall time), which is the time from first task started to last task stopped.
	ElapsedTime time.Duration
	// Parallelism is the parallelism factor of the job, which is the average number of tasks running at any given time.
	Parallelism float64
	// TotalTasks is the total number of tasks of the job.
	TotalTasks uint64
	// TasksByState is the number of tasks by their state.
	TasksByState []*TaskStateCount
}

// TaskStateCount represents the count of tasks in a particular state.
type TaskStateCount struct {
	// State is the task state.
	State TaskState
	// Count is the number of tasks in this state.
	Count uint64
}

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
	//   - job.WithAutomationIDs: specifies the automation IDs to filter jobs by. Only jobs submitted by the specified
	//  automation will be returned. (Optional)
	//   - job.WithJobState: specifies the job state to filter jobs by. Only jobs with the specified state will be returned. (Optional)
	//   - job.WithJobName: specifies the job name to filter jobs by. Only jobs with the specified name will be returned. (Optional)
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

	return protoToJob(response), nil
}

func (c jobClient) Get(ctx context.Context, jobID uuid.UUID) (*Job, error) {
	response, err := c.service.GetJob(ctx, jobID)
	if err != nil {
		return nil, err
	}

	return protoToJob(response), nil
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
		var page *tileboxv1.Pagination // nil for the first request

		// we already validated that TemporalExtent is not nil
		timeInterval := opts.TemporalExtent.ToProtoTimeInterval()
		idInterval := opts.TemporalExtent.ToProtoIDInterval()

		if timeInterval == nil && idInterval == nil {
			yield(nil, errors.New("invalid temporal extent"))
			return
		}

		var automationIDs []*tileboxv1.ID
		if opts.AutomationID != uuid.Nil {
			automationIDs = append(automationIDs, tileboxv1.NewUUID(opts.AutomationID))
		}
		for _, id := range opts.AutomationIDs {
			automationIDs = append(automationIDs, tileboxv1.NewUUID(id))
		}

		var states []workflowsv1.JobState
		for _, state := range opts.States {
			states = append(states, workflowsv1.JobState(int32(state)))
		}

		filters := workflowsv1.QueryFilters_builder{
			TimeInterval:  timeInterval,
			IdInterval:    idInterval,
			AutomationIds: automationIDs,
			States:        states,
			Name:          opts.Name,
		}.Build()

		for {
			jobsMessage, err := c.service.QueryJobs(ctx, filters, page)
			if err != nil {
				yield(nil, err)
				return
			}

			for _, jobMessage := range jobsMessage.GetJobs() {
				if !yield(protoToJob(jobMessage), nil) {
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

	submissions := make([]*taskSubmission, 0, len(tasks))

	for _, task := range tasks {
		var input []byte
		var err error

		identifier := identifierFromTask(task)
		err = ValidateIdentifier(identifier)
		if err != nil {
			return nil, fmt.Errorf("task has invalid task identifier: %w", err)
		}

		taskProto, isProtobuf := task.(proto.Message)
		if isProtobuf {
			input, err = proto.Marshal(taskProto)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal protobuf task: %w", err)
			}
		} else {
			input, err = json.Marshal(task)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal task: %w", err)
			}
		}

		submissions = append(submissions, &taskSubmission{
			clusterSlug: clusterSlug,
			identifier:  identifier,
			input:       input,
			maxRetries:  maxRetries,
		})
	}

	return workflowsv1.SubmitJobRequest_builder{
		Tasks:   mergeTasksToSubmissions(submissions),
		JobName: jobName,
	}.Build(), nil
}

func protoToJob(jobMessage *workflowsv1.Job) *Job {
	taskSummaries := make([]*TaskSummary, 0, len(jobMessage.GetTaskSummaries()))
	for _, taskSummary := range jobMessage.GetTaskSummaries() {
		taskSummaries = append(taskSummaries, protoToTaskSummary(taskSummary))
	}

	progressIndicators := make([]*ProgressIndicator, 0, len(jobMessage.GetProgress()))
	for _, progress := range jobMessage.GetProgress() {
		progressIndicators = append(progressIndicators, protoToProgressIndicator(progress))
	}

	return &Job{
		ID:             jobMessage.GetId().AsUUID(),
		Name:           jobMessage.GetName(),
		State:          job.State(jobMessage.GetState()),
		SubmittedAt:    jobMessage.GetSubmittedAt().AsTime(),
		TaskSummaries:  taskSummaries,
		AutomationID:   jobMessage.GetAutomationId().AsUUID(),
		Progress:       progressIndicators,
		ExecutionStats: protoToExecutionStats(jobMessage.GetExecutionStats()),
	}
}

func protoToExecutionStats(es *workflowsv1.ExecutionStats) *ExecutionStats {
	if es == nil {
		return nil
	}
	tasksByState := make([]*TaskStateCount, 0, len(es.GetTasksByState()))
	for _, tsc := range es.GetTasksByState() {
		tasksByState = append(tasksByState, &TaskStateCount{
			State: TaskState(tsc.GetState()),
			Count: tsc.GetCount(),
		})
	}
	return &ExecutionStats{
		FirstTaskStartedAt: es.GetFirstTaskStartedAt().AsTime(),
		LastTaskStoppedAt:  es.GetLastTaskStoppedAt().AsTime(),
		ComputeTime:        es.GetComputeTime().AsDuration(),
		ElapsedTime:        es.GetElapsedTime().AsDuration(),
		Parallelism:        es.GetParallelism(),
		TotalTasks:         es.GetTotalTasks(),
		TasksByState:       tasksByState,
	}
}

func protoToTaskSummary(t *workflowsv1.TaskSummary) *TaskSummary {
	return &TaskSummary{
		ID:        t.GetId().AsUUID(),
		Display:   t.GetDisplay(),
		State:     TaskState(t.GetState()),
		ParentID:  t.GetParentId().AsUUID(),
		StartedAt: t.GetStartedAt().AsTime(),
		StoppedAt: t.GetStoppedAt().AsTime(),
	}
}

func protoToProgressIndicator(p *workflowsv1.Progress) *ProgressIndicator {
	return &ProgressIndicator{
		Label: p.GetLabel(),
		Total: p.GetTotal(),
		Done:  p.GetDone(),
	}
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
