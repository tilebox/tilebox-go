package accounts

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tilebox/tilebox-go/client"
	accountsv1alpha1 "github.com/tilebox/tilebox-go/protogen/accounts/v1alpha1"
	"github.com/tilebox/tilebox-go/protogen/accounts/v1alpha1/accountsv1alpha1connect"
)

func TestClient_GetAccountDetails(t *testing.T) {
	service := &fakeAccountsService{}
	client := newTestClient(t, service)

	details, err := client.Account.GetAccountDetails(context.Background())
	require.NoError(t, err)

	assert.Equal(t, "Test User", details.GetUserName())
	assert.Equal(t, "test-organization", details.GetOrganizationSlug())
	assert.Equal(t, "Bearer test-api-key", service.authorization)
	assert.Contains(t, service.clientMetadata, `name="go"`)
	assert.Contains(t, service.clientMetadata, `runtime="go"`)
}

func TestClientMetadataCanBeOverridden(t *testing.T) {
	service := &fakeAccountsService{}
	client := newTestClient(t, service, WithClientMetadata(client.Metadata{Name: "cli", Version: "v1.2.3"}))

	_, err := client.Account.GetAccountDetails(t.Context())
	require.NoError(t, err)

	assert.Equal(t, `name="cli", version="v1.2.3"`, service.clientMetadata)
}

func TestClient_GetActivePlan(t *testing.T) {
	service := &fakeAccountsService{}
	client := newTestClient(t, service)

	plan, err := client.Billing.GetActivePlan(context.Background())
	require.NoError(t, err)

	assert.Equal(t, accountsv1alpha1.SubscriptionTier_SUBSCRIPTION_TIER_PAID, plan.GetTier())
	assert.Equal(t, "Bearer test-api-key", service.authorization)
}

func TestClient_GetUsageReport(t *testing.T) {
	tests := []struct {
		name                string
		options             []UsageReportOption
		expectedHistoryDays uint64
	}{
		{
			name: "current usage",
		},
		{
			name:                "usage with history",
			options:             []UsageReportOption{WithHistoryDays(30)},
			expectedHistoryDays: 30,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &fakeAccountsService{}
			client := newTestClient(t, service)

			report, err := client.Billing.GetUsageReport(context.Background(), tt.options...)
			require.NoError(t, err)

			require.Len(t, report.GetMetrics(), 1)
			assert.Equal(t, "storage_bytes", report.GetMetrics()[0].GetKey())
			assert.Equal(t, tt.expectedHistoryDays, service.historyDays)
			assert.Equal(t, "Bearer test-api-key", service.authorization)
		})
	}
}

func newTestClient(t *testing.T, service *fakeAccountsService, options ...ClientOption) *Client {
	t.Helper()

	mux := http.NewServeMux()
	accountPath, accountHandler := accountsv1alpha1connect.NewAccountServiceHandler(service)
	billingPath, billingHandler := accountsv1alpha1connect.NewBillingServiceHandler(service)
	mux.Handle(accountPath, accountHandler)
	mux.Handle(billingPath, billingHandler)

	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	options = append(options,
		WithURL(server.URL),
		WithHTTPClient(server.Client()),
		WithAPIKey("test-api-key"),
		WithDisableTracing(),
	)
	return NewClient(options...)
}

type fakeAccountsService struct {
	authorization  string
	clientMetadata string
	historyDays    uint64
}

func (s *fakeAccountsService) GetAccountDetails(_ context.Context, request *connect.Request[accountsv1alpha1.GetAccountDetailsRequest]) (*connect.Response[accountsv1alpha1.AccountDetails], error) {
	s.authorization = request.Header().Get("Authorization")
	s.clientMetadata = request.Header().Get("Tilebox-Client")
	return connect.NewResponse(accountsv1alpha1.AccountDetails_builder{
		UserName:         "Test User",
		OrganizationSlug: "test-organization",
	}.Build()), nil
}

func (s *fakeAccountsService) GetActivePlan(_ context.Context, request *connect.Request[accountsv1alpha1.GetActivePlanRequest]) (*connect.Response[accountsv1alpha1.Plan], error) {
	s.authorization = request.Header().Get("Authorization")
	return connect.NewResponse(accountsv1alpha1.Plan_builder{
		Tier: accountsv1alpha1.SubscriptionTier_SUBSCRIPTION_TIER_PAID,
	}.Build()), nil
}

func (s *fakeAccountsService) GetUsageReport(_ context.Context, request *connect.Request[accountsv1alpha1.GetUsageReportRequest]) (*connect.Response[accountsv1alpha1.UsageReport], error) {
	s.authorization = request.Header().Get("Authorization")
	s.historyDays = request.Msg.GetHistoryDays()
	return connect.NewResponse(accountsv1alpha1.UsageReport_builder{
		Metrics: []*accountsv1alpha1.UsageMetric{
			accountsv1alpha1.UsageMetric_builder{Key: "storage_bytes"}.Build(),
		},
	}.Build()), nil
}
