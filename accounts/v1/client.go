// Package accounts provides a client for interacting with Tilebox Accounts.
package accounts

import (
	"connectrpc.com/connect/v2"
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

// NewClient builds Accounts services over a shared Connect client.
// Use client.New for Tilebox defaults or supply your own native Connect client.
// Tracing uses the global tracer provider unless disabled.
func NewClient(rpcClient *connect.Client, options ...ClientOption) *Client {
	cfg := newClientConfig(options)
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
		tracerProvider: otel.GetTracerProvider(),
	}
	for _, option := range options {
		option(cfg)
	}

	return cfg
}
