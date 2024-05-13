package workflows

import (
	"connectrpc.com/connect"
	"github.com/tilebox/tilebox-go/grpc"
	"github.com/tilebox/tilebox-go/protogen/go/workflows/v1/workflowsv1connect"
)

// clientConfig contains the configuration for a gRPC client to a workflows service.
type clientConfig struct {
	httpClient     connect.HTTPClient
	url            string
	apiKey         string
	connectOptions []connect.ClientOption
}

// ClientOption is an interface for configuring a client. Using such options helpers is a
// quite common pattern in Go, as it allows for optional parameters in constructors.
// This concrete implementation here is inspired by how libraries such as axiom-go and connect do their
// configuration.
type ClientOption func(*clientConfig)

func WithHTTPClient(httpClient connect.HTTPClient) ClientOption {
	return func(cfg *clientConfig) {
		cfg.httpClient = httpClient
	}
}

func WithURL(url string) ClientOption {
	return func(cfg *clientConfig) {
		cfg.url = url
	}
}

func WithAPIKey(apiKey string) ClientOption {
	return func(cfg *clientConfig) {
		cfg.apiKey = apiKey
	}
}

func WithConnectClientOptions(options ...connect.ClientOption) ClientOption {
	return func(cfg *clientConfig) {
		cfg.connectOptions = append(cfg.connectOptions, options...)
	}
}

func newClientConfig(options []ClientOption) *clientConfig {
	cfg := &clientConfig{
		httpClient: grpc.RetryHTTPClient(),
		url:        "https://api.tilebox.com",
	}
	for _, option := range options {
		option(cfg)
	}
	return cfg
}

func NewTaskClient(options ...ClientOption) workflowsv1connect.TaskServiceClient {
	cfg := newClientConfig(options)

	return workflowsv1connect.NewTaskServiceClient(
		cfg.httpClient,
		cfg.url,
		connect.WithClientOptions(cfg.connectOptions...),
		connect.WithInterceptors(
			grpc.NewAddAuthTokenInterceptor(func() string {
				return cfg.apiKey
			})),
	)
}

func NewJobClient(options ...ClientOption) workflowsv1connect.JobServiceClient {
	cfg := newClientConfig(options)

	return workflowsv1connect.NewJobServiceClient(
		cfg.httpClient,
		cfg.url,
		connect.WithClientOptions(cfg.connectOptions...),
		connect.WithInterceptors(
			grpc.NewAddAuthTokenInterceptor(func() string {
				return cfg.apiKey
			})),
	)
}
