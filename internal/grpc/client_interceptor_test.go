package grpc

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestAddClientMetadataInterceptor(t *testing.T) {
	headerValue := `name="cli", version="v1.2.3"`
	interceptor := NewAddClientMetadataInterceptor(headerValue)
	request := connect.NewRequest(&emptypb.Empty{})
	request.Header().Set("Authorization", "Bearer token")

	wrapped := interceptor.WrapUnary(func(_ context.Context, got connect.AnyRequest) (connect.AnyResponse, error) {
		require.Equal(t, headerValue, got.Header().Get(ClientHeader))
		require.Equal(t, "Bearer token", got.Header().Get("Authorization"))
		return connect.NewResponse(&emptypb.Empty{}), nil
	})

	_, err := wrapped(t.Context(), request)
	require.NoError(t, err)
}
