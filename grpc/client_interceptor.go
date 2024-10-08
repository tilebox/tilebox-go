package grpc

import (
	"context"

	"connectrpc.com/connect"
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
		request.Header().Set("Authorization", "Bearer "+at.token())
		return next(ctx, request)
	}
}
