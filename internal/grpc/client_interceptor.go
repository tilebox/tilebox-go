package grpc

import (
	"context"
	"runtime/debug"

	"connectrpc.com/connect"
)

const (
	ClientSourceHeader  = "Tilebox-Client-Source"
	ClientVersionHeader = "Tilebox-Client-Version"
	modulePath          = "github.com/tilebox/tilebox-go"
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

type addClientMetadataInterceptor struct {
	connect.Interceptor

	source  string
	version string
}

func NewAddClientMetadataInterceptor(source, version string) connect.Interceptor {
	return &addClientMetadataInterceptor{source: source, version: version}
}

func (i *addClientMetadataInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, request connect.AnyRequest) (connect.AnyResponse, error) {
		if request.Header().Get(ClientSourceHeader) == "" {
			request.Header().Set(ClientSourceHeader, i.source)
		}
		if request.Header().Get(ClientVersionHeader) == "" {
			request.Header().Set(ClientVersionHeader, i.version)
		}
		return next(ctx, request)
	}
}

// ClientVersion returns the tilebox-go module version embedded in the running binary.
func ClientVersion() string {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return "dev"
	}
	if buildInfo.Main.Path == modulePath && buildInfo.Main.Version != "(devel)" {
		return buildInfo.Main.Version
	}
	for _, dependency := range buildInfo.Deps {
		if dependency.Path == modulePath {
			return dependency.Version
		}
	}
	return "dev"
}
