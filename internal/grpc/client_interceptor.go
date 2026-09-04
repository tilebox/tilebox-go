package grpc

import (
	"context"

	"connectrpc.com/connect"
)

const ClientHeader = "Tilebox-Client"

type addAuthTokenInterceptor struct {
	connect.Interceptor

	token func() string
}

func NewAddAuthTokenInterceptor(token func() string) connect.Interceptor {
	return &addAuthTokenInterceptor{token: token}
}

func (at *addAuthTokenInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, request connect.AnyRequest) (connect.AnyResponse, error) {
		request.Header().Set("Authorization", "Bearer "+at.token())
		return next(ctx, request)
	}
}

type addClientMetadataInterceptor struct {
	connect.Interceptor

	headerValue string
}

func NewAddClientMetadataInterceptor(headerValue string) connect.Interceptor {
	return &addClientMetadataInterceptor{headerValue: headerValue}
}

func (i *addClientMetadataInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, request connect.AnyRequest) (connect.AnyResponse, error) {
		if i.headerValue != "" {
			request.Header().Set(ClientHeader, i.headerValue)
		}
		return next(ctx, request)
	}
}
