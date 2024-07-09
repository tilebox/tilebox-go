package workflows

import (
	"context"
	"net"
	"net/http"
	"strings"

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
		url: "https://api.tilebox.com",
	}
	for _, option := range options {
		option(cfg)
	}

	// if no http client is set by the user, we use a default one
	if cfg.httpClient == nil {
		// if the URL looks like an HTTP URL, we use a retrying HTTP client
		if strings.HasPrefix(cfg.url, "https://") || strings.HasPrefix(cfg.url, "http://") {
			cfg.httpClient = grpc.RetryHTTPClient()
		} else { // we connect to a unix socket
			address := cfg.url // we copy the url to a temporary variable so that we can modify cfg.url after
			dial := func(context.Context, string, string) (net.Conn, error) {
				return net.Dial("unix", address)
			}
			transport := &http.Transport{DialContext: dial}
			cfg.httpClient = &http.Client{Transport: transport}
			cfg.url = "http://localhost" // connect requires a dummy url starting with http://
		}
	}

	return cfg
}

func NewTaskClient(options ...ClientOption) workflowsv1connect.TaskServiceClient {
	return newConnectClient(workflowsv1connect.NewTaskServiceClient, options...)
}

func NewJobClient(options ...ClientOption) workflowsv1connect.JobServiceClient {
	return newConnectClient(workflowsv1connect.NewJobServiceClient, options...)
}

func NewRecurrentTaskClient(options ...ClientOption) workflowsv1connect.RecurrentTaskServiceClient {
	return newConnectClient(workflowsv1connect.NewRecurrentTaskServiceClient, options...)
}

func newConnectClient[T any](newClientFunc func(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) T, options ...ClientOption) T {
	cfg := newClientConfig(options)

	// for on-orbit deployments, we also support requests without an API key, so we don't need to add an interceptor
	// for those cases
	interceptors := make([]connect.Interceptor, 0)
	if cfg.apiKey != "" {
		interceptors = append(interceptors, grpc.NewAddAuthTokenInterceptor(func() string {
			return cfg.apiKey
		}))
	}

	return newClientFunc(
		cfg.httpClient,
		cfg.url,
		connect.WithClientOptions(cfg.connectOptions...),
		connect.WithInterceptors(interceptors...),
	)
}
