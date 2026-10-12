package grpc

import (
	"context"
	"io"
	"testing"

	"connectrpc.com/connect/v2"
	"connectrpc.com/connect/v2/connectinprocess"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestClientHeadersInterceptor(t *testing.T) {
	for _, test := range []struct {
		name, apiKey, metadata, wantAuth, wantMetadata string
	}{
		{"SDK headers", "sdk-key", `name="cli"`, "Bearer sdk-key", `name="cli"`},
		{"caller headers", "", "", "Bearer caller-key", `name="caller"`},
	} {
		t.Run(test.name, func(t *testing.T) {
			for _, streamType := range []connect.StreamType{
				connect.StreamTypeUnary, connect.StreamTypeServer, connect.StreamTypeClient, connect.StreamTypeBidi,
			} {
				t.Run(streamType.String(), func(t *testing.T) {
					t.Parallel()
					spec := connect.Spec{Procedure: "/test.Headers/Echo", StreamType: streamType}
					server := connect.NewServer()
					server.Register(connect.Method{Spec: spec, Handler: func(ctx context.Context, _ connect.Spec, stream connect.ServerStream) error {
						info, _ := connect.CallInfoForServerContext(ctx)
						var request wrapperspb.StringValue
						if err := stream.Receive(&request); err != nil {
							return err
						}
						info.ResponseHeader().Set("X-Response", "header")
						info.ResponseTrailer().Set("X-Response", "trailer")
						return stream.Send(wrapperspb.String(info.RequestHeader().Get("Authorization") + "|" +
							info.RequestHeader().Get(ClientHeader) + "|" + info.RequestHeader().Get("X-Custom") + "|" + request.GetValue()))
					}})
					client := connect.NewClient(connectinprocess.New(server), NewClientHeadersInterceptor(test.apiKey, test.metadata))
					ctx, info := connect.NewClientContext(t.Context())
					info.RequestHeader().Set("Authorization", "Bearer caller-key")
					info.RequestHeader().Set(ClientHeader, `name="caller"`)
					info.RequestHeader().Set("X-Custom", "preserved")
					request := wrapperspb.String("payload")
					var response wrapperspb.StringValue
					if streamType == connect.StreamTypeUnary {
						require.NoError(t, client.CallUnary(ctx, spec, request, &response))
					} else {
						var stream connect.ClientStream
						var err error
						if streamType == connect.StreamTypeServer {
							stream, err = client.CallServerStream(ctx, spec, request)
						} else {
							stream, err = client.CallClientStream(ctx, spec)
						}
						require.NoError(t, err)
						defer stream.Close()
						if streamType != connect.StreamTypeServer {
							require.NoError(t, stream.Send(request))
							require.NoError(t, stream.CloseSend())
						}
						require.NoError(t, stream.Receive(&response))
						require.ErrorIs(t, stream.Receive(&wrapperspb.StringValue{}), io.EOF)
					}
					require.Equal(t, test.wantAuth+"|"+test.wantMetadata+"|preserved|payload", response.GetValue())
					require.Equal(t, "header", info.ResponseHeader().Get("X-Response"))
					require.Equal(t, "trailer", info.ResponseTrailer().Get("X-Response"))
				})
			}
		})
	}
}
