package accounts // import "github.com/tilebox/tilebox-go/accounts/v1"

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/tilebox/tilebox-go/observability"
	accountsv1alpha1 "github.com/tilebox/tilebox-go/protogen/accounts/v1alpha1"
	"github.com/tilebox/tilebox-go/protogen/accounts/v1alpha1/accountsv1alpha1connect"
	"go.opentelemetry.io/otel/trace"
)

// AccountClient provides access to account details for the authenticated credential.
type AccountClient interface {
	// GetAccountDetails returns details about the account associated with the authenticated credential.
	GetAccountDetails(ctx context.Context) (*accountsv1alpha1.AccountDetails, error)
}

var _ AccountClient = &accountClient{}

type accountClient struct {
	connectClient accountsv1alpha1connect.AccountServiceClient
	tracer        trace.Tracer
}

func (c *accountClient) GetAccountDetails(ctx context.Context) (*accountsv1alpha1.AccountDetails, error) {
	return observability.WithSpanResult(ctx, c.tracer, "accounts/details/get", func(ctx context.Context) (*accountsv1alpha1.AccountDetails, error) {
		response, err := c.connectClient.GetAccountDetails(ctx, connect.NewRequest(
			accountsv1alpha1.GetAccountDetailsRequest_builder{}.Build(),
		))
		if err != nil {
			return nil, fmt.Errorf("failed to get account details: %w", err)
		}

		return response.Msg, nil
	})
}
