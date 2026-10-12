package workflows

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"strings"
	"time"

	"connectrpc.com/connect/v2"
	"connectrpc.com/connect/v2/connecthttp"
	"github.com/tilebox/tilebox-go/client"
	"github.com/tilebox/tilebox-go/internal/grpc"
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

// NewClient creates a new Tilebox Workflows client.
//
// By default, the returned Client is configured with:
//   - "https://api.tilebox.com" as the URL
//   - environment variable TILEBOX_API_KEY as the API key
//   - a pooled HTTP/2-capable client with transient-error retries
//   - the global tracer provider
//   - Tilebox OpenTelemetry exporters when an API key and HTTP(S) API URL are configured
//
// The passed options are used to override these default values and configure the returned Client appropriately.
func NewClient(options ...ClientOption) *Client {
	cfg := newClientConfig(options)
	configureTileboxTelemetry(context.Background(), cfg)

	rpcClient := connect.NewClient(
		connecthttp.NewTransport(cfg.httpClient, cfg.url, cfg.transportOptions...),
		cfg.interceptors...,
	)
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
			service:    storageLocationService,
			httpClient: cfg.notificationHTTPClient, baseURL: cfg.url, apiKey: cfg.apiKey,
			clientMetadata: cfg.clientMetadata.HeaderValue(),
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
	httpClient             connecthttp.HTTPClient
	notificationHTTPClient connecthttp.HTTPClient
	url                    string
	apiKey                 string
	clientMetadata         client.Metadata
	transportOptions       []connecthttp.Option
	interceptors           []connect.ClientInterceptor

	tracerProvider trace.TracerProvider
	disableTracing bool
	disableLogging bool
}

// ClientOption is an interface for configuring a client. Using such options helpers is a
// quite common pattern in Go, as it allows for optional parameters in constructors.
// This concrete implementation here is inspired by how libraries such as axiom-go and connect do their
// configuration.
type ClientOption func(*clientConfig)

// WithHTTPClient sets the HTTP client, including its connection pool and retry policy.
//
// By default, HTTPS uses HTTP/2 with HTTP/1.1 fallback and transient-error retries.
func WithHTTPClient(httpClient connecthttp.HTTPClient) ClientOption {
	return func(cfg *clientConfig) {
		cfg.httpClient = httpClient
	}
}

// WithURL sets the URL of the Tilebox Datasets service.
//
// Defaults to "https://api.tilebox.com".
func WithURL(url string) ClientOption {
	return func(cfg *clientConfig) {
		cfg.url = url
	}
}

// WithAPIKey sets the API key to use for the client.
//
// Defaults to no API key.
func WithAPIKey(apiKey string) ClientOption {
	return func(cfg *clientConfig) {
		cfg.apiKey = apiKey
	}
}

// WithClientMetadata replaces the automatically detected metadata sent with each request.
// Wrappers such as the Tilebox CLI can use this to identify themselves as the client.
func WithClientMetadata(metadata client.Metadata) ClientOption {
	return func(cfg *clientConfig) {
		cfg.clientMetadata = metadata
	}
}

// WithTransportOptions configures the RPC wire protocol, compression, and message limits.
// Defaults to binary Connect with a 4 MiB response message limit.
func WithTransportOptions(options ...connecthttp.Option) ClientOption {
	return func(cfg *clientConfig) {
		cfg.transportOptions = append(cfg.transportOptions, options...)
	}
}

// WithInterceptors adds client interceptors in outermost-first order, before the
// SDK sets authentication and client metadata headers.
func WithInterceptors(interceptors ...connect.ClientInterceptor) ClientOption {
	return func(cfg *clientConfig) {
		cfg.interceptors = append(cfg.interceptors, interceptors...)
	}
}

// WithDisableTracing disables OpenTelemetry tracing for the client.
func WithDisableTracing() ClientOption {
	return func(cfg *clientConfig) {
		cfg.tracerProvider = noop.NewTracerProvider()
		cfg.disableTracing = true
	}
}

// WithDisableLogging disables Tilebox OpenTelemetry log export setup for the client.
func WithDisableLogging() ClientOption {
	return func(cfg *clientConfig) {
		cfg.disableLogging = true
	}
}

func newClientConfig(options []ClientOption) *clientConfig {
	cfg := &clientConfig{
		url:            "https://api.tilebox.com",
		apiKey:         os.Getenv("TILEBOX_API_KEY"),
		clientMetadata: client.DefaultMetadata(),
		tracerProvider: otel.GetTracerProvider(), // use the global tracer provider by default
	}
	for _, option := range options {
		option(cfg)
	}

	cfg.notificationHTTPClient = cfg.httpClient
	if cfg.httpClient == nil {
		cfg.httpClient, cfg.url = grpc.NewHTTPClient(cfg.url)
		// Notifications are not idempotent. Share the pool, but never retry them.
		cfg.notificationHTTPClient = grpc.NonRetryHTTPClient(cfg.httpClient, 30*time.Second)
	}
	cfg.interceptors = append(cfg.interceptors, grpc.NewClientHeadersInterceptor(cfg.apiKey, cfg.clientMetadata.HeaderValue()))

	return cfg
}
