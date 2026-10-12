package accounts

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"connectrpc.com/connect/v2"
	"connectrpc.com/connect/v2/connecthttp"
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

func TestClientInterceptorsAndCallMetadata(t *testing.T) {
	addHeader := func(value string) connect.ClientInterceptor {
		return func(next connect.ClientFunc) connect.ClientFunc {
			return func(ctx context.Context, spec connect.Spec) (connect.ClientStream, error) {
				info, _ := connect.CallInfoForClientContext(ctx)
				info.RequestHeader().Set("X-Caller", info.RequestHeader().Get("X-Caller")+value)
				info.RequestHeader().Set("Authorization", "Bearer interceptor-key")
				return next(ctx, spec)
			}
		}
	}
	service := &fakeAccountsService{}
	client := newTestClient(t, service, WithHTTPClient(nil),
		WithInterceptors(addHeader("A")), WithInterceptors(addHeader("B")),
		WithTransportOptions(connecthttp.WithSendGzip()),
	)
	ctx, info := connect.NewClientContext(t.Context())
	info.RequestHeader().Set("X-Caller", "caller-")
	details, err := client.Account.GetAccountDetails(ctx)
	require.NoError(t, err)
	assert.Equal(t, "Test User", details.GetUserName())
	assert.Equal(t, "Bearer test-api-key", service.authorization)
	assert.Equal(t, "caller-AB", info.ResponseHeader().Get("X-Caller"))
	assert.Equal(t, "complete", info.ResponseTrailer().Get("X-Result"))
}

func TestClientResponseLimit(t *testing.T) {
	for _, test := range []struct {
		name    string
		options []ClientOption
	}{
		{name: "default"},
		{name: "larger limit", options: []ClientOption{WithTransportOptions(connecthttp.WithReadMaxBytes(8 << 20))}},
	} {
		t.Run(test.name, func(t *testing.T) {
			service := &fakeAccountsService{userName: strings.Repeat("x", (4<<20)+1)}
			client := newTestClient(t, service, test.options...)
			details, err := client.Account.GetAccountDetails(t.Context())
			if test.name == "default" {
				require.Error(t, err)
				assert.Equal(t, connect.CodeResourceExhausted, connect.CodeOf(err))
			} else {
				require.NoError(t, err)
				assert.Equal(t, service.userName, details.GetUserName())
			}
		})
	}
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
	rpcServer := connect.NewServer()
	accountsv1alpha1connect.RegisterAccountServiceHandler(rpcServer, service)
	accountsv1alpha1connect.RegisterBillingServiceHandler(rpcServer, service)
	connecthttp.Mount(mux, rpcServer)

	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	options = append([]ClientOption{
		WithURL(server.URL),
		WithHTTPClient(server.Client()),
		WithAPIKey("test-api-key"),
		WithDisableTracing(),
	}, options...)
	return NewClient(options...)
}

type fakeAccountsService struct {
	authorization  string
	clientMetadata string
	historyDays    uint64
	userName       string
}

func (s *fakeAccountsService) GetAccountDetails(ctx context.Context, _ *accountsv1alpha1.GetAccountDetailsRequest) (*accountsv1alpha1.AccountDetails, error) {
	info, _ := connect.CallInfoForServerContext(ctx)
	s.authorization = info.RequestHeader().Get("Authorization")
	s.clientMetadata = info.RequestHeader().Get("Tilebox-Client")
	info.ResponseHeader().Set("X-Caller", info.RequestHeader().Get("X-Caller"))
	info.ResponseTrailer().Set("X-Result", "complete")
	userName := s.userName
	if userName == "" {
		userName = "Test User"
	}
	return accountsv1alpha1.AccountDetails_builder{
		UserName:         userName,
		OrganizationSlug: "test-organization",
	}.Build(), nil
}

func (s *fakeAccountsService) GetActivePlan(ctx context.Context, _ *accountsv1alpha1.GetActivePlanRequest) (*accountsv1alpha1.Plan, error) {
	info, _ := connect.CallInfoForServerContext(ctx)
	s.authorization = info.RequestHeader().Get("Authorization")
	return accountsv1alpha1.Plan_builder{
		Tier: accountsv1alpha1.SubscriptionTier_SUBSCRIPTION_TIER_PAID,
	}.Build(), nil
}

func (s *fakeAccountsService) GetUsageReport(ctx context.Context, request *accountsv1alpha1.GetUsageReportRequest) (*accountsv1alpha1.UsageReport, error) {
	info, _ := connect.CallInfoForServerContext(ctx)
	s.authorization = info.RequestHeader().Get("Authorization")
	s.historyDays = request.GetHistoryDays()
	return accountsv1alpha1.UsageReport_builder{
		Metrics: []*accountsv1alpha1.UsageMetric{
			accountsv1alpha1.UsageMetric_builder{Key: "storage_bytes"}.Build(),
		},
	}.Build(), nil
}
