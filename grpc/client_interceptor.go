package grpc

import (
	"connectrpc.com/connect"
	"context"
)

type addAuthTokenInterceptor struct {
	connect.Interceptor
	token func() string
}

func NewAddAuthTokenInterceptor(token func() string) connect.Interceptor {
	return &addAuthTokenInterceptor{token: token}
}

func (at *addAuthTokenInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, request connect.AnyRequest) (connect.AnyResponse, error) {
		request.Header().Set("authorization", "Bearer "+at.token())
		return next(ctx, request)
	}
}