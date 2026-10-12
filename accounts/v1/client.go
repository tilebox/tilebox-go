// Package accounts provides a client for interacting with Tilebox Accounts.
package accounts

import (
	"os"

	"connectrpc.com/connect/v2"
	"connectrpc.com/connect/v2/connecthttp"
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
//   - a pooled HTTP/2-capable client with transient-error retries
//   - the global tracer provider
//
// The passed options are used to override these default values and configure the returned Client appropriately.
func NewClient(options ...ClientOption) *Client {
	cfg := newClientConfig(options)
	rpcClient := connect.NewClient(
		connecthttp.NewTransport(cfg.httpClient, cfg.url, cfg.transportOptions...),
		cfg.interceptors...,
	)
	accountConnectClient := accountsv1alpha1connect.NewAccountServiceClient(rpcClient)
	billingConnectClient := accountsv1alpha1connect.NewBillingServiceClient(rpcClient)
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
	httpClient       connecthttp.HTTPClient
	url              string
	apiKey           string
	clientMetadata   client.Metadata
	transportOptions []connecthttp.Option
	interceptors     []connect.ClientInterceptor

	tracerProvider trace.TracerProvider
}

// ClientOption configures a client.
type ClientOption func(*clientConfig)

// WithHTTPClient sets the HTTP client, including its connection pool and retry policy.
//
// By default, HTTPS uses HTTP/2 with HTTP/1.1 fallback and transient-error retries.
func WithHTTPClient(httpClient connecthttp.HTTPClient) ClientOption {
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
		cfg.httpClient, cfg.url = grpc.NewHTTPClient(cfg.url)
	}
	cfg.interceptors = append(cfg.interceptors, grpc.NewClientHeadersInterceptor(cfg.apiKey, cfg.clientMetadata.HeaderValue()))

	return cfg
}
