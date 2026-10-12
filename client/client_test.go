package client_test

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"connectrpc.com/connect/v2"
	"connectrpc.com/connect/v2/connecthttp"
	"connectrpc.com/connect/v2/connectinprocess"
	"github.com/stretchr/testify/require"
	accounts "github.com/tilebox/tilebox-go/accounts/v1"
	"github.com/tilebox/tilebox-go/client"
	datasets "github.com/tilebox/tilebox-go/datasets/v1"
	accountspb "github.com/tilebox/tilebox-go/protogen/accounts/v1alpha1"
	"github.com/tilebox/tilebox-go/protogen/accounts/v1alpha1/accountsv1alpha1connect"
	datasetspb "github.com/tilebox/tilebox-go/protogen/datasets/v1"
	"github.com/tilebox/tilebox-go/protogen/datasets/v1/datasetsv1connect"
	workflowspb "github.com/tilebox/tilebox-go/protogen/workflows/v1"
	"github.com/tilebox/tilebox-go/protogen/workflows/v1/workflowsv1connect"
	workflows "github.com/tilebox/tilebox-go/workflows/v1"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type contractServices struct {
	datasetsv1connect.UnimplementedDatasetServiceHandler
	workflowsv1connect.UnimplementedWorkflowsServiceHandler

	headers []http.Header
}

func (s *contractServices) GetAccountDetails(ctx context.Context, _ *accountspb.GetAccountDetailsRequest) (*accountspb.AccountDetails, error) {
	s.record(ctx)
	return accountspb.AccountDetails_builder{UserName: "shared-client"}.Build(), nil
}

func (s *contractServices) ListDatasets(ctx context.Context, _ *datasetspb.ListDatasetsRequest) (*datasetspb.ListDatasetsResponse, error) {
	s.record(ctx)
	return datasetspb.ListDatasetsResponse_builder{}.Build(), nil
}

func (s *contractServices) ListClusters(ctx context.Context, _ *workflowspb.ListClustersRequest) (*workflowspb.ListClustersResponse, error) {
	s.record(ctx)
	return workflowspb.ListClustersResponse_builder{
		Clusters: []*workflowspb.Cluster{workflowspb.Cluster_builder{Slug: "shared-cluster"}.Build()},
	}.Build(), nil
}

func (s *contractServices) record(ctx context.Context) {
	info, _ := connect.CallInfoForServerContext(ctx)
	headers := make(http.Header)
	for key, values := range info.RequestHeader().All() {
		headers[key] = append([]string(nil), values...)
	}
	s.headers = append(s.headers, headers)
}

type contractRoundTripper func(*http.Request) (*http.Response, error)

func (f contractRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}

func TestSharedNativeClientAcrossDomains(t *testing.T) {
	services := &contractServices{}
	server := connect.NewServer()
	accountsv1alpha1connect.RegisterAccountServiceHandler(server, services)
	datasetsv1connect.RegisterDatasetServiceHandler(server, services)
	workflowsv1connect.RegisterWorkflowsServiceHandler(server, services)
	var httpCalls atomic.Int32
	forbiddenHTTP := contractRoundTripper(func(*http.Request) (*http.Response, error) {
		httpCalls.Add(1)
		return nil, errors.New("native transport must not use HTTP")
	})
	interceptor := func(name string) connect.ClientInterceptor {
		return func(next connect.ClientFunc) connect.ClientFunc {
			return func(ctx context.Context, spec connect.Spec) (connect.ClientStream, error) {
				info, _ := connect.CallInfoForClientContext(ctx)
				headers := info.RequestHeader()
				headers.Set("X-Order", headers.Get("X-Order")+name)
				// These must be replaced by the SDK interceptor, which runs last.
				headers.Set("Authorization", "user-"+name)
				headers.Set(client.ClientHeader, "user-"+name)
				return next(ctx, spec)
			}
		}
	}
	rpc := client.New(
		client.WithTransport(connectinprocess.New(server)),
		client.WithURL(":invalid HTTP URL"),
		client.WithHTTPClient(&http.Client{Transport: forbiddenHTTP}),
		client.WithHTTPTransport(forbiddenHTTP),
		client.WithRetryPolicy(client.RetryPolicy{MaxRetries: -1}),
		client.WithTransportOptions(connecthttp.WithReadMaxBytes(1)),
		client.WithAPIKey("shared-key"),
		client.WithClientMetadata(client.Metadata{Name: "contract"}),
		client.WithInterceptors(interceptor("first")),
		client.WithInterceptors(interceptor("second")),
	)
	account, err := accounts.NewClient(rpc).Account.GetAccountDetails(t.Context())
	require.NoError(t, err)
	require.Equal(t, "shared-client", account.GetUserName())
	listedDatasets, err := datasets.NewClient(rpc).Datasets.List(t.Context())
	require.NoError(t, err)
	require.Empty(t, listedDatasets)
	clusters, err := workflows.NewClient(rpc).Clusters.List(t.Context())
	require.NoError(t, err)
	require.Len(t, clusters, 1)
	require.Equal(t, "shared-cluster", clusters[0].Slug)
	require.Len(t, services.headers, 3)
	for _, headers := range services.headers {
		require.Equal(t, "firstsecond", headers.Get("X-Order"))
		require.Equal(t, "Bearer shared-key", headers.Get("Authorization"))
		require.Equal(t, `name="contract"`, headers.Get(client.ClientHeader))
	}
	require.Zero(t, httpCalls.Load())
}

var contractEchoSpec = connect.Spec{Procedure: "/test.Client/Echo", StreamType: connect.StreamTypeUnary}

func contractEchoHandler() http.Handler {
	server := connect.NewServer()
	server.Register(connect.Method{Spec: contractEchoSpec, Handler: func(_ context.Context, _ connect.Spec, stream connect.ServerStream) error {
		var request wrapperspb.StringValue
		if err := stream.Receive(&request); err != nil {
			return err
		}
		return stream.Send(&request)
	}})
	mux := http.NewServeMux()
	connecthttp.Mount(mux, server)
	return mux
}

func TestDefaultClientIgnoresWrappedDefaultTransport(t *testing.T) {
	server := httptest.NewServer(contractEchoHandler())
	t.Cleanup(server.Close)
	var wrapperCalls atomic.Int32
	original := http.DefaultTransport
	http.DefaultTransport = contractRoundTripper(func(request *http.Request) (*http.Response, error) {
		wrapperCalls.Add(1)
		return original.RoundTrip(request)
	})
	t.Cleanup(func() { http.DefaultTransport = original })
	var rpc *connect.Client
	require.NotPanics(t, func() { rpc = client.New(client.WithURL(server.URL)) })
	var response wrapperspb.StringValue
	require.NoError(t, rpc.CallUnary(t.Context(), contractEchoSpec, wrapperspb.String("default"), &response))
	require.Equal(t, "default", response.GetValue())
	require.Zero(t, wrapperCalls.Load())
}

func TestClientHTTPRetryConfiguration(t *testing.T) {
	for _, test := range []struct {
		name      string
		retries   int
		custom    string
		failures  int32
		attempts  int32
		wantError bool
	}{
		{name: "zero disables retries", failures: 1, attempts: 1, wantError: true},
		{name: "one retry succeeds", retries: 1, failures: 1, attempts: 2},
		{name: "three retries exhaust", retries: 3, failures: 4, attempts: 4, wantError: true},
		{name: "custom transport retains retries", retries: 2, custom: "transport", failures: 2, attempts: 3},
		{name: "full HTTP client bypasses retries", retries: 3, custom: "client", failures: 1, attempts: 1, wantError: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			var attempts, transportCalls atomic.Int32
			echo := contractEchoHandler()
			server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				if attempts.Add(1) <= test.failures {
					writer.WriteHeader(http.StatusServiceUnavailable)
					return
				}
				echo.ServeHTTP(writer, request)
			}))
			t.Cleanup(server.Close)
			options := []client.Option{
				client.WithURL(server.URL),
				client.WithRetryPolicy(client.RetryPolicy{MaxRetries: test.retries, MinWait: time.Millisecond, MaxWait: time.Millisecond}),
			}
			transport := client.NewHTTPTransport()
			t.Cleanup(transport.CloseIdleConnections)
			countingTransport := contractRoundTripper(func(request *http.Request) (*http.Response, error) {
				transportCalls.Add(1)
				return transport.RoundTrip(request)
			})
			if test.custom != "" {
				options = append(options, client.WithHTTPTransport(countingTransport))
			}
			if test.custom == "client" {
				options = append(options, client.WithHTTPClient(&http.Client{Transport: countingTransport}))
			}
			var response wrapperspb.StringValue
			err := client.New(options...).CallUnary(t.Context(), contractEchoSpec, wrapperspb.String("retry payload"), &response)
			if test.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, "retry payload", response.GetValue())
			}
			require.Equal(t, test.attempts, attempts.Load())
			if test.custom != "" {
				require.Equal(t, test.attempts, transportCalls.Load())
			}
		})
	}
}

func TestClientUnixSocketRetries(t *testing.T) {
	directory, err := os.MkdirTemp("/tmp", "rpc-") //nolint:usetesting // t.TempDir paths can exceed the Unix socket path limit on macOS.
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, os.RemoveAll(directory)) })
	socketPath := filepath.Join(directory, "rpc.sock")
	listener, err := (&net.ListenConfig{}).Listen(t.Context(), "unix", socketPath)
	require.NoError(t, err)
	var attempts atomic.Int32
	echo := contractEchoHandler()
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if attempts.Add(1) == 1 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		echo.ServeHTTP(w, r)
	}))
	require.NoError(t, server.Listener.Close())
	server.Listener = listener
	server.Start()
	t.Cleanup(server.Close)
	rpc := client.New(client.WithURL(socketPath))
	var response wrapperspb.StringValue
	require.NoError(t, rpc.CallUnary(t.Context(), contractEchoSpec, wrapperspb.String("socket payload"), &response))
	require.Equal(t, "socket payload", response.GetValue())
	require.EqualValues(t, 2, attempts.Load())

	t.Run("custom HTTP client", func(t *testing.T) {
		transport := client.NewHTTPTransport()
		t.Cleanup(transport.CloseIdleConnections)
		transport.Proxy = nil
		transport.DialContext = func(ctx context.Context, _, _ string) (net.Conn, error) {
			return (&net.Dialer{}).DialContext(ctx, "unix", socketPath)
		}
		httpClient := &http.Client{Transport: transport, Timeout: time.Second}
		rpc := client.New(client.WithURL(socketPath), client.WithHTTPClient(httpClient))
		var response wrapperspb.StringValue
		require.NoError(t, rpc.CallUnary(t.Context(), contractEchoSpec, wrapperspb.String("custom socket payload"), &response))
		require.Equal(t, "custom socket payload", response.GetValue())
		require.EqualValues(t, 3, attempts.Load())
	})
}
