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
// This concrete implementation here is inspired by how connect does its client configuration.
type ClientOption interface {
	applyToClient(config *clientConfig)
}

type httpClientOption struct {
	httpClient connect.HTTPClient
}

func WithHTTPClient(httpClient connect.HTTPClient) ClientOption {
	return &httpClientOption{httpClient}
}

func (o *httpClientOption) applyToClient(config *clientConfig) {
	config.httpClient = o.httpClient
}

type apiKeyOption struct {
	apiKey string
}

type urlOption struct {
	url string
}

func WithURL(url string) ClientOption {
	return &urlOption{url}
}

func (o *urlOption) applyToClient(config *clientConfig) {
	config.url = o.url
}

func WithAPIKey(apiKey string) ClientOption {
	return &apiKeyOption{apiKey}
}

func (o *apiKeyOption) applyToClient(config *clientConfig) {
	config.apiKey = o.apiKey
}

type connectOptions struct {
	options []connect.ClientOption
}

func WithConnectClientOptions(options ...connect.ClientOption) ClientOption {
	return &connectOptions{options}
}

func (o *connectOptions) applyToClient(config *clientConfig) {
	config.connectOptions = append(config.connectOptions, o.options...)
}

func newClientConfig(options []ClientOption) clientConfig {
	cfg := clientConfig{
		httpClient: grpc.RetryHTTPClient(),
		url:        "https://api.tilebox.com",
	}
	for _, opt := range options {
		opt.applyToClient(&cfg)
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
