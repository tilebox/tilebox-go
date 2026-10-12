package datasets

import (
	"connectrpc.com/connect/v2"
	"github.com/tilebox/tilebox-go/protogen/datasets/v1/datasetsv1connect"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

const otelTracerName = "tilebox.com/observability"

// Client is a Tilebox Datasets client.
type Client struct {
	Datasets    DatasetClient
	Collections CollectionClient
	Datapoints  DatapointClient
}

// NewClient builds Datasets services over a shared Connect client.
// Use client.New for Tilebox defaults or supply your own native Connect client.
// Tracing uses the global tracer provider unless disabled.
func NewClient(rpcClient *connect.Client, options ...ClientOption) *Client {
	cfg := newClientConfig(options)
	datasetConnectClient := datasetsv1connect.NewDatasetServiceClient(rpcClient)
	collectionConnectClient := datasetsv1connect.NewCollectionServiceClient(rpcClient)
	dataAccessConnectClient := datasetsv1connect.NewDataAccessServiceClient(rpcClient)
	dataIngestionConnectClient := datasetsv1connect.NewDataIngestionServiceClient(rpcClient)

	tracer := cfg.tracerProvider.Tracer(otelTracerName)

	return &Client{
		Datasets: &datasetClient{
			service: newDatasetsService(datasetConnectClient, tracer),
		},
		Collections: &collectionClient{
			service: newCollectionService(collectionConnectClient, tracer),
		},
		Datapoints: &datapointClient{
			dataAccessService:    newDataAccessService(dataAccessConnectClient, tracer),
			dataIngestionService: newDataIngestionService(dataIngestionConnectClient, tracer),
		},
	}
}

// clientConfig contains the configuration for Tilebox Datasets client.
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
