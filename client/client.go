// Package client constructs shared Connect clients for Tilebox APIs.
package client

import (
	"net/http"
	"os"
	"strings"

	"connectrpc.com/connect/v2"
	"connectrpc.com/connect/v2/connecthttp"
)

// New creates a Connect client that can be shared by Accounts, Datasets, and
// Workflows. Defaults are https://api.tilebox.com, TILEBOX_API_KEY, SDK metadata,
// and binary Connect over pooled HTTP/2 (HTTP/1.1 fallback) with transient retries.
// Construction does not configure global telemetry providers or loggers.
func New(options ...Option) *connect.Client {
	cfg := config{
		url: "https://api.tilebox.com", apiKey: os.Getenv("TILEBOX_API_KEY"),
		metadata: DefaultMetadata(), retryPolicy: DefaultRetryPolicy(),
		transportOptions: []connecthttp.Option{connecthttp.WithReadMaxBytes(0)},
	}
	for _, option := range options {
		option(&cfg)
	}
	if cfg.transport == nil {
		var socketPath string
		if !strings.HasPrefix(cfg.url, "https://") && !strings.HasPrefix(cfg.url, "http://") {
			socketPath, cfg.url = cfg.url, "http://localhost"
		}
		if cfg.httpClient == nil {
			cfg.httpClient = newHTTPClient(socketPath, cfg.httpTransport, cfg.retryPolicy)
		}
		cfg.transport = connecthttp.NewTransport(cfg.httpClient, cfg.url, cfg.transportOptions...)
	}
	return connect.NewClient(cfg.transport,
		append(cfg.interceptors, newClientHeadersInterceptor(cfg.apiKey, cfg.metadata.HeaderValue()))...)
}

type config struct {
	url              string
	apiKey           string
	metadata         Metadata
	transport        connect.Transport
	httpClient       connecthttp.HTTPClient
	httpTransport    http.RoundTripper
	retryPolicy      RetryPolicy
	transportOptions []connecthttp.Option
	interceptors     []connect.ClientInterceptor
}

// Option configures a shared Connect client.
type Option func(*config)

// WithURL sets the HTTP(S) API URL or Unix socket path.
func WithURL(url string) Option {
	return func(cfg *config) { cfg.url = url }
}

// WithAPIKey overrides TILEBOX_API_KEY. An empty key preserves caller-supplied authentication.
func WithAPIKey(apiKey string) Option {
	return func(cfg *config) { cfg.apiKey = apiKey }
}

// WithClientMetadata replaces the automatically detected SDK metadata.
func WithClientMetadata(metadata Metadata) Option {
	return func(cfg *config) { cfg.metadata = metadata }
}

// WithTransport supplies a native Connect transport, for example connectinprocess.New.
// It bypasses all HTTP configuration, including URL, retries, and transport options.
// Authentication, metadata, and client interceptors still apply.
func WithTransport(transport connect.Transport) Option {
	return func(cfg *config) { cfg.transport = transport }
}

// WithHTTPTransport replaces the connection transport while retaining SDK retries.
// The transport is used unchanged; callers own its TLS, proxy, and protocol settings.
// Start with NewHTTPTransport to retain the SDK's connection pool defaults.
// For Unix sockets, a custom transport must supply its own socket dialer.
func WithHTTPTransport(transport http.RoundTripper) Option {
	return func(cfg *config) { cfg.httpTransport = transport }
}

// WithHTTPClient replaces the complete HTTP client, including retries, redirects,
// and timeouts. WithHTTPTransport and WithRetryPolicy do not apply to this client.
// For Unix sockets, the client must supply its own socket dialer; the RPC base URL
// is still normalized to http://localhost.
func WithHTTPClient(httpClient connecthttp.HTTPClient) Option {
	return func(cfg *config) { cfg.httpClient = httpClient }
}

// WithRetryPolicy configures HTTP retries for the default client. A zero policy
// disables retries. Invalid limits panic when the default HTTP client is constructed.
// Retries require a replayable request body; streaming requests are never buffered.
func WithRetryPolicy(policy RetryPolicy) Option {
	return func(cfg *config) { cfg.retryPolicy = policy }
}

// WithTransportOptions configures the HTTP wire protocol, compression, and message
// limits. Defaults are binary Connect and no response message size limit.
func WithTransportOptions(options ...connecthttp.Option) Option {
	return func(cfg *config) { cfg.transportOptions = append(cfg.transportOptions, options...) }
}

// WithInterceptors adds native interceptors in outermost-first order, before the
// SDK sets authentication and client metadata headers.
func WithInterceptors(interceptors ...connect.ClientInterceptor) Option {
	return func(cfg *config) { cfg.interceptors = append(cfg.interceptors, interceptors...) }
}
