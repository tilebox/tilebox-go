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

type jobServiceConfig struct {
	tracerProvider trace.TracerProvider
	tracerName     string
}

type JobServiceOption func(*jobServiceConfig)

func WithJobServiceTracerProvider(tracerProvider trace.TracerProvider) JobServiceOption {
	return func(cfg *jobServiceConfig) {
		cfg.tracerProvider = tracerProvider
	}
}

func WithJobServiceTracerName(tracerName string) JobServiceOption {
	return func(cfg *jobServiceConfig) {
		cfg.tracerName = tracerName
	}
}

type JobService struct {
	client workflowsv1connect.JobServiceClient
	tracer trace.Tracer
}

func newJobServiceConfig(options []JobServiceOption) *jobServiceConfig {
	cfg := &jobServiceConfig{
		tracerProvider: otel.GetTracerProvider(),    // use the global tracer provider by default
		tracerName:     "tilebox.com/observability", // the default tracer name we use
	}
	for _, option := range options {
		option(cfg)
	}

	return cfg
}

func NewJobService(client workflowsv1connect.JobServiceClient, options ...JobServiceOption) *JobService {
	cfg := newJobServiceConfig(options)
	return &JobService{
		client: client,
		tracer: cfg.tracerProvider.Tracer(cfg.tracerName),
	}
}

func (js *JobService) Submit(ctx context.Context, jobName, clusterSlug string, maxRetries int, tasks ...Task) (*workflowsv1.Job, error) {
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
