package workflows

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"connectrpc.com/connect"
	"github.com/tilebox/tilebox-go/observability"
	workflowsv1 "github.com/tilebox/tilebox-go/protogen/go/workflows/v1"
	"github.com/tilebox/tilebox-go/protogen/go/workflows/v1/workflowsv1connect"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/proto"
)

type JobService struct {
	client workflowsv1connect.JobServiceClient
	tracer trace.Tracer
}

func NewJobService(client workflowsv1connect.JobServiceClient) *JobService {
	return &JobService{
		client: client,
		tracer: otel.Tracer("tilebox.com/observability"),
	}
}

func (js *JobService) Submit(ctx context.Context, jobName, clusterSlug string, tasks ...Task) (*workflowsv1.Job, error) {
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
				Name:    identifier.Name,
				Version: identifier.Version,
			},
			Input:   subtaskInput,
			Display: identifier.Name,
		})
	}

	return observability.WithSpanResult(ctx, js.tracer, fmt.Sprintf("job/%s", jobName), func(ctx context.Context) (*workflowsv1.Job, error) {
		traceParent := observability.GetTraceParentOfCurrentSpan(ctx)

		job, err := js.client.SubmitJob(ctx, connect.NewRequest(
			&workflowsv1.SubmitJobRequest{
				Tasks:       rootTasks,
				JobName:     jobName,
				TraceParent: traceParent,
			}))

		if err != nil {
			return nil, fmt.Errorf("failed to submit job: %w", err)
		}

		return job.Msg, nil
	})
}
