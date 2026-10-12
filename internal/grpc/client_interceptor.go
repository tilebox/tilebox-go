package grpc

import (
	"context"

	"connectrpc.com/connect/v2"
)

const ClientHeader = "Tilebox-Client"

// NewClientHeadersInterceptor sets SDK-owned headers for every RPC type.
// An empty API key leaves caller-provided authentication intact for deployments
// with custom authentication or no authentication.
func NewClientHeadersInterceptor(apiKey, metadata string) connect.ClientInterceptor {
	return func(next connect.ClientFunc) connect.ClientFunc {
		return func(ctx context.Context, spec connect.Spec) (connect.ClientStream, error) {
			info, _ := connect.CallInfoForClientContext(ctx)
			if apiKey != "" {
				info.RequestHeader().Set("Authorization", "Bearer "+apiKey)
			}
			if metadata != "" {
				info.RequestHeader().Set(ClientHeader, metadata)
			}
			return next(ctx, spec)
		}
	}
}
