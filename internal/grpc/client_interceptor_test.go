package grpc

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestAddClientMetadataInterceptor(t *testing.T) {
	interceptor := NewAddClientMetadataInterceptor("cli", "v1.2.3")
	request := connect.NewRequest(&emptypb.Empty{})
	request.Header().Set("Authorization", "Bearer token")

	wrapped := interceptor.WrapUnary(func(_ context.Context, got connect.AnyRequest) (connect.AnyResponse, error) {
		require.Equal(t, "cli", got.Header().Get(ClientSourceHeader))
		require.Equal(t, "v1.2.3", got.Header().Get(ClientVersionHeader))
		require.Equal(t, "Bearer token", got.Header().Get("Authorization"))
		return connect.NewResponse(&emptypb.Empty{}), nil
	})

	_, err := wrapped(t.Context(), request)
	require.NoError(t, err)
}

func TestAddClientMetadataInterceptorPreservesWrapperOverride(t *testing.T) {
	interceptor := NewAddClientMetadataInterceptor("go_sdk", "v0.11.1")
	request := connect.NewRequest(&emptypb.Empty{})
	request.Header().Set(ClientSourceHeader, "cli")
	request.Header().Set(ClientVersionHeader, "v1.2.3")

	wrapped := interceptor.WrapUnary(func(_ context.Context, got connect.AnyRequest) (connect.AnyResponse, error) {
		require.Equal(t, "cli", got.Header().Get(ClientSourceHeader))
		require.Equal(t, "v1.2.3", got.Header().Get(ClientVersionHeader))
		return connect.NewResponse(&emptypb.Empty{}), nil
	})

	_, err := wrapped(t.Context(), request)
	require.NoError(t, err)
}

func TestClientVersion(t *testing.T) {
	require.NotEmpty(t, ClientVersion())
}
