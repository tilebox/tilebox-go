package workflows

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"connectrpc.com/connect/v2"
	"github.com/tilebox/tilebox-go/protogen/workflows/v1/workflowsv1connect"
	"github.com/tilebox/tilebox-go/workflows/v1/runner"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

const (
	otelTracerName = "tilebox.com/observability"
	otelMeterName  = otelTracerName
)

// Client is a Tilebox Workflows client.
type Client struct {
	Jobs             JobClient
	Clusters         ClusterClient
	Workflows        WorkflowClient
	Automations      AutomationClient
	StorageLocations StorageLocationClient

	// Used by NewTaskRunner
	taskService TaskService
	tracer      trace.Tracer
}

// NewClient builds Workflows services over a shared Connect client.
// Use client.New for Tilebox defaults or supply your own native Connect client.
// Tracing uses the global tracer provider unless disabled. To export telemetry
// to Tilebox, call ConfigureTileboxTelemetry before constructing this client.
func NewClient(rpcClient *connect.Client, options ...ClientOption) *Client {
	cfg := newClientConfig(options)
	jobConnectClient := workflowsv1connect.NewJobServiceClient(rpcClient)
	telemetryConnectClient := workflowsv1connect.NewTelemetryQueryServiceClient(rpcClient)
	taskConnectClient := workflowsv1connect.NewTaskServiceClient(rpcClient)
	workflowConnectClient := workflowsv1connect.NewWorkflowsServiceClient(rpcClient)
	automationConnectClient := workflowsv1connect.NewAutomationServiceClient(rpcClient)
	storageLocationConnectClient := workflowsv1connect.NewStorageLocationServiceClient(rpcClient)

	tracer := cfg.tracerProvider.Tracer(otelTracerName)
	workflowService := newWorkflowService(workflowConnectClient, tracer)
	automationService := &automationService{automationClient: automationConnectClient, tracer: tracer}
	storageLocationService := &storageLocationService{client: storageLocationConnectClient, tracer: tracer}

	return &Client{
		Jobs:        &jobClient{service: newJobService(jobConnectClient, tracer), telemetryService: newTelemetryService(telemetryConnectClient, tracer)},
		Clusters:    &clusterClient{service: workflowService},
		Workflows:   &workflowClient{service: workflowService},
		Automations: &automationClient{service: automationService},
		StorageLocations: &storageLocationClient{
			service: storageLocationService,
		},

		taskService: newTaskService(taskConnectClient, tracer),
		tracer:      tracer,
	}
}

// NewTaskRunner creates a new TaskRunner.
//
// Options:
//   - runner.WithClusterSlug: sets the cluster for which to listen tasks to. Default to the default cluster.
//   - runner.WithRunnerLogger: sets the logger to use for the task runner. Defaults to slog.Default().
//   - runner.WithDisableMetrics: disables OpenTelemetry metrics for the task runner.
func (c *Client) NewTaskRunner(ctx context.Context, options ...runner.Option) (*TaskRunner, error) {
	return newTaskRunner(ctx, c.taskService, c.Clusters, c.tracer, options...)
}

// NewPollingTaskRunner creates a polling task runner for a resolved cluster and custom task executor.
func (c *Client) NewPollingTaskRunner(cluster *Cluster, executor TaskExecutor, logger *slog.Logger) (*PollingTaskRunner, error) {
	if cluster == nil {
		return nil, errors.New("cluster is required")
	}
	if strings.TrimSpace(cluster.Slug) == "" {
		return nil, errors.New("cluster slug is required")
	}
	return NewPollingTaskRunner(c.taskService, cluster.Slug, executor, logger), nil
}

// clientConfig contains the configuration for Tilebox Workflows client.
type clientConfig struct {
	tracerProvider trace.TracerProvider
}

// ClientOption configures a client.
type ClientOption func(*clientConfig)

// WithDisableTracing disables OpenTelemetry tracing for the client.
func WithDisableTracing() ClientOption {
	return func(cfg *clientConfig) {
		cfg.tracerProvider = noop.NewTracerProvider()
	}
}

func newClientConfig(options []ClientOption) *clientConfig {
	cfg := &clientConfig{
		tracerProvider: otel.GetTracerProvider(), // use the global tracer provider by default
	}
	for _, option := range options {
		option(cfg)
	}

	return cfg
}
