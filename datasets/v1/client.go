package datasets // import "github.com/tilebox/tilebox-go/datasets/v1"

import (
	"context"
	"net"
	"net/http"
	"strings"

	"connectrpc.com/connect"
	"github.com/tilebox/tilebox-go/grpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

// clientConfig contains the configuration for Tilebox Datasets client.
type clientConfig struct {
	httpClient     connect.HTTPClient
	url            string
	apiKey         string
	connectOptions []connect.ClientOption

	tracerProvider trace.TracerProvider
}

// ClientOption is an interface for configuring a client. Using such options helpers is a
// quite common pattern in Go, as it allows for optional parameters in constructors.
// This concrete implementation here is inspired by how libraries such as axiom-go and connect do their
// configuration.
type ClientOption func(*clientConfig)

// WithHTTPClient sets the connect.HTTPClient to use for the client.
//
// Defaults to grpc.RetryHTTPClient.
func WithHTTPClient(httpClient connect.HTTPClient) ClientOption {
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

// WithConnectClientOptions sets additional options for the connect.HTTPClient.
func WithConnectClientOptions(options ...connect.ClientOption) ClientOption {
	return func(cfg *clientConfig) {
		cfg.connectOptions = append(cfg.connectOptions, options...)
	}
}

// WithDisableTracing disables OpenTelemetry tracing for the client.
func WithDisableTracing() ClientOption {
	return func(cfg *clientConfig) {
		cfg.tracerProvider = noop.NewTracerProvider()
	}
}

func newClientConfig(options []ClientOption) *clientConfig {
	cfg := &clientConfig{
		url:            "https://api.tilebox.com",
		tracerProvider: otel.GetTracerProvider(), // use the global tracer provider by default
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

func newConnectClient[T any](newClientFunc func(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) T, cfg *clientConfig) T {
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
