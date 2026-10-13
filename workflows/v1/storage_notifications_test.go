package workflows

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"connectrpc.com/connect/v2"
	"connectrpc.com/connect/v2/connecthttp"
	"connectrpc.com/connect/v2/connectinprocess"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tilebox/tilebox-go/client"
	tileboxv1 "github.com/tilebox/tilebox-go/protogen/tilebox/v1"
	workflowsv1 "github.com/tilebox/tilebox-go/protogen/workflows/v1"
	"github.com/tilebox/tilebox-go/protogen/workflows/v1/workflowsv1connect"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestNotifyObjectCreatedSharedTransport(t *testing.T) {
	for _, test := range []struct {
		name      string
		inProcess bool
		policy    *client.RetryPolicy
		failures  int32
		attempts  int32
		wantError bool
	}{
		{name: "default HTTP retries", failures: 2, attempts: 3},
		{name: "configured HTTP retries", policy: &client.RetryPolicy{MaxRetries: 1, MinWait: time.Millisecond, MaxWait: time.Millisecond}, failures: 1, attempts: 2},
		{name: "disabled HTTP retries", policy: &client.RetryPolicy{}, failures: 1, attempts: 1, wantError: true},
		{name: "in-process", inProcess: true, attempts: 1},
	} {
		t.Run(test.name, func(t *testing.T) {
			id := uuid.New()
			expected := workflowsv1.NotifyObjectCreatedRequest_builder{
				StorageSubscriptionId: tileboxv1.NewUUID(id), ObjectPath: "incoming/scene.tif",
			}.Build()
			var attempts atomic.Int32
			service := storageLocationRPCService{call: func(ctx context.Context, request proto.Message) (proto.Message, error) {
				info, ok := connect.CallInfoForServerContext(ctx)
				assert.True(t, ok)
				assert.Equal(t, "Bearer shared-key", info.RequestHeader().Get("Authorization"))
				assert.Equal(t, `name="test-cli"`, info.RequestHeader().Get(client.ClientHeader))
				assert.Equal(t, "shared-interceptor", info.RequestHeader().Get("X-Test"))
				assert.True(t, proto.Equal(expected, request), "request: %v", request)
				if attempts.Add(1) <= test.failures {
					return nil, connect.NewError(connect.CodeUnavailable, "try again")
				}
				return &emptypb.Empty{}, nil
			}}
			rpcServer := connect.NewServer()
			workflowsv1connect.RegisterStorageLocationServiceHandler(rpcServer, service)
			options := []client.Option{
				client.WithAPIKey("shared-key"),
				client.WithClientMetadata(client.Metadata{Name: "test-cli"}),
				client.WithInterceptors(func(next connect.ClientFunc) connect.ClientFunc {
					return func(ctx context.Context, spec connect.Spec) (connect.ClientStream, error) {
						info, _ := connect.CallInfoForClientContext(ctx)
						info.RequestHeader().Set("X-Test", "shared-interceptor")
						return next(ctx, spec)
					}
				}),
			}
			if test.inProcess {
				options = append(options, client.WithTransport(connectinprocess.New(rpcServer)))
			} else {
				mux := http.NewServeMux()
				connecthttp.Mount(mux, rpcServer)
				server := httptest.NewServer(mux)
				t.Cleanup(server.Close)
				options = append(options, client.WithURL(server.URL))
			}
			if test.policy != nil {
				options = append(options, client.WithRetryPolicy(*test.policy))
			}
			c := NewClient(client.New(options...))
			err := c.StorageLocations.NotifyObjectCreated(t.Context(), id, "incoming/scene.tif", time.Time{})
			if test.wantError {
				require.Error(t, err)
				assert.Equal(t, connect.CodeUnavailable, connect.CodeOf(err))
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, test.attempts, attempts.Load())
		})
	}
}

func TestNotifyObjectCreatedCancellation(t *testing.T) {
	started := make(chan struct{})
	service := storageLocationRPCService{call: func(ctx context.Context, _ proto.Message) (proto.Message, error) {
		close(started)
		<-ctx.Done()
		return nil, ctx.Err()
	}}
	rpcServer := connect.NewServer()
	workflowsv1connect.RegisterStorageLocationServiceHandler(rpcServer, service)
	mux := http.NewServeMux()
	connecthttp.Mount(mux, rpcServer)
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)
	c := NewClient(client.New(client.WithURL(server.URL)))
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	done := make(chan error, 1)
	go func() { done <- c.StorageLocations.NotifyObjectCreated(ctx, uuid.New(), "file.tif", time.Time{}) }()
	select {
	case <-started:
	case <-time.After(5 * time.Second):
		t.Fatal("request did not start")
	}
	cancel()
	select {
	case err := <-done:
		require.Error(t, err)
		assert.Equal(t, connect.CodeCanceled, connect.CodeOf(err))
		require.ErrorContains(t, err, "delivery may have succeeded")
	case <-time.After(5 * time.Second):
		t.Fatal("request did not cancel")
	}
}
