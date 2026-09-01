package accounts

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/tilebox/tilebox-go/observability"
	accountsv1alpha1 "github.com/tilebox/tilebox-go/protogen/accounts/v1alpha1"
	"github.com/tilebox/tilebox-go/protogen/accounts/v1alpha1/accountsv1alpha1connect"
	"go.opentelemetry.io/otel/trace"
)

// BillingClient provides access to account billing information.
type BillingClient interface {
	// GetActivePlan returns the active subscription plan for the authenticated account.
	GetActivePlan(ctx context.Context) (*accountsv1alpha1.Plan, error)

	// GetUsageReport returns the current usage report for the authenticated account.
	//
	// Options:
	//   - WithHistoryDays: includes historical values for the requested number of days.
	GetUsageReport(ctx context.Context, options ...UsageReportOption) (*accountsv1alpha1.UsageReport, error)
}

var _ BillingClient = &billingClient{}

type billingClient struct {
	connectClient accountsv1alpha1connect.BillingServiceClient
	tracer        trace.Tracer
}

func (c *billingClient) GetActivePlan(ctx context.Context) (*accountsv1alpha1.Plan, error) {
	return observability.WithSpanResult(ctx, c.tracer, "accounts/billing/active_plan/get", func(ctx context.Context) (*accountsv1alpha1.Plan, error) {
		response, err := c.connectClient.GetActivePlan(ctx, connect.NewRequest(
			accountsv1alpha1.GetActivePlanRequest_builder{}.Build(),
		))
		if err != nil {
			return nil, fmt.Errorf("failed to get active plan: %w", err)
		}

		return response.Msg, nil
	})
}

func (c *billingClient) GetUsageReport(ctx context.Context, options ...UsageReportOption) (*accountsv1alpha1.UsageReport, error) {
	usageReportOptions := newUsageReportOptions(options)
	return observability.WithSpanResult(ctx, c.tracer, "accounts/billing/usage_report/get", func(ctx context.Context) (*accountsv1alpha1.UsageReport, error) {
		response, err := c.connectClient.GetUsageReport(ctx, connect.NewRequest(
			accountsv1alpha1.GetUsageReportRequest_builder{
				HistoryDays: usageReportOptions.historyDays,
			}.Build(),
		))
		if err != nil {
			return nil, fmt.Errorf("failed to get usage report: %w", err)
		}

		return response.Msg, nil
	})
}

type usageReportOptions struct {
	historyDays uint64
}

// UsageReportOption configures a usage report request.
type UsageReportOption func(*usageReportOptions)

// WithHistoryDays includes historical usage values for the requested number of days.
// The API supports up to 365 days.
func WithHistoryDays(historyDays uint64) UsageReportOption {
	return func(options *usageReportOptions) {
		options.historyDays = historyDays
	}
}

func newUsageReportOptions(options []UsageReportOption) usageReportOptions {
	var result usageReportOptions
	for _, option := range options {
		option(&result)
	}
	return result
}
