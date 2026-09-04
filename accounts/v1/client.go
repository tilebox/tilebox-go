// Package accounts provides a client for interacting with Tilebox Accounts.
package accounts

import (
	"context"
	"net"
	"net/http"
	"os"
	"strings"

	"connectrpc.com/connect"
	"github.com/tilebox/tilebox-go/client"
	"github.com/tilebox/tilebox-go/internal/grpc"
	"github.com/tilebox/tilebox-go/protogen/accounts/v1alpha1/accountsv1alpha1connect"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

const otelTracerName = "tilebox.com/observability"

// Client is a Tilebox Accounts client.
type Client struct {
	Account AccountClient
	Billing BillingClient
}

// NewClient creates a new Tilebox Accounts client.
//
// By default, the returned Client is configured with:
//   - "https://api.tilebox.com" as the URL
//   - environment variable TILEBOX_API_KEY as the API key
//   - a grpc.RetryHTTPClient HTTP client
//   - the global tracer provider
//
// The passed options are used to override these default values and configure the returned Client appropriately.
func NewClient(options ...ClientOption) *Client {
	cfg := newClientConfig(options)
	accountConnectClient := newConnectClient(accountsv1alpha1connect.NewAccountServiceClient, cfg)
	billingConnectClient := newConnectClient(accountsv1alpha1connect.NewBillingServiceClient, cfg)
	tracer := cfg.tracerProvider.Tracer(otelTracerName)

	return &Client{
		Account: &accountClient{
			connectClient: accountConnectClient,
			tracer:        tracer,
		},
		Billing: &billingClient{
			connectClient: billingConnectClient,
			tracer:        tracer,
		},
	}
}

// clientConfig contains the configuration for a Tilebox Accounts client.
type clientConfig struct {
	httpClient     connect.HTTPClient
	url            string
	apiKey         string
	clientMetadata client.Metadata
	connectOptions []connect.ClientOption

	tracerProvider trace.TracerProvider
}

// ClientOption configures a client.
type ClientOption func(*clientConfig)

// WithHTTPClient sets the connect.HTTPClient to use for the client.
//
// Defaults to grpc.RetryHTTPClient.
func WithHTTPClient(httpClient connect.HTTPClient) ClientOption {
	return func(cfg *clientConfig) {
		cfg.httpClient = httpClient
	}
}

// WithURL sets the URL of the Tilebox Accounts service.
//
// Defaults to "https://api.tilebox.com".
func WithURL(url string) ClientOption {
	return func(cfg *clientConfig) {
		cfg.url = url
	}
}

// WithAPIKey sets the API key to use for the client.
//
// Defaults to the TILEBOX_API_KEY environment variable.
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
		apiKey:         os.Getenv("TILEBOX_API_KEY"),
		clientMetadata: client.DefaultMetadata(),
		tracerProvider: otel.GetTracerProvider(),
	}
	for _, option := range options {
		option(cfg)
	}

	if cfg.httpClient == nil {
		if strings.HasPrefix(cfg.url, "https://") || strings.HasPrefix(cfg.url, "http://") {
			cfg.httpClient = grpc.RetryHTTPClient()
		} else {
			address := cfg.url
			dial := func(ctx context.Context, _ string, _ string) (net.Conn, error) {
				var dialer net.Dialer
				return dialer.DialContext(ctx, "unix", address)
			}
			transport := &http.Transport{DialContext: dial}
			cfg.httpClient = &http.Client{Transport: transport}
			cfg.url = "http://localhost"
		}
	}

	return cfg
}

func newConnectClient[T any](newClientFunc func(httpClient connect.HTTPClient, baseURL string, options ...connect.ClientOption) T, cfg *clientConfig) T {
	interceptors := []connect.Interceptor{grpc.NewAddClientMetadataInterceptor(cfg.clientMetadata.HeaderValue())}
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
