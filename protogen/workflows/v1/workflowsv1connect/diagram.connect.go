// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: workflows/v1/diagram.proto

package workflowsv1connect

import (
	connect "connectrpc.com/connect"
	context "context"
	errors "errors"
	v1 "github.com/tilebox/tilebox-go/protogen/workflows/v1"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect.IsAtLeastVersion1_13_0

const (
	// DiagramServiceName is the fully-qualified name of the DiagramService service.
	DiagramServiceName = "workflows.v1.DiagramService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// DiagramServiceRenderProcedure is the fully-qualified name of the DiagramService's Render RPC.
	DiagramServiceRenderProcedure = "/workflows.v1.DiagramService/Render"
)

// DiagramServiceClient is a client for the workflows.v1.DiagramService service.
type DiagramServiceClient interface {
	Render(context.Context, *connect.Request[v1.RenderDiagramRequest]) (*connect.Response[v1.Diagram], error)
}

// NewDiagramServiceClient constructs a client for the workflows.v1.DiagramService service. By
// default, it uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses,
// and sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the
// connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewDiagramServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) DiagramServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	diagramServiceMethods := v1.File_workflows_v1_diagram_proto.Services().ByName("DiagramService").Methods()
	return &diagramServiceClient{
		render: connect.NewClient[v1.RenderDiagramRequest, v1.Diagram](
			httpClient,
			baseURL+DiagramServiceRenderProcedure,
			connect.WithSchema(diagramServiceMethods.ByName("Render")),
			connect.WithClientOptions(opts...),
		),
	}
}

// diagramServiceClient implements DiagramServiceClient.
type diagramServiceClient struct {
	render *connect.Client[v1.RenderDiagramRequest, v1.Diagram]
}

// Render calls workflows.v1.DiagramService.Render.
func (c *diagramServiceClient) Render(ctx context.Context, req *connect.Request[v1.RenderDiagramRequest]) (*connect.Response[v1.Diagram], error) {
	return c.render.CallUnary(ctx, req)
}

// DiagramServiceHandler is an implementation of the workflows.v1.DiagramService service.
type DiagramServiceHandler interface {
	Render(context.Context, *connect.Request[v1.RenderDiagramRequest]) (*connect.Response[v1.Diagram], error)
}

// NewDiagramServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewDiagramServiceHandler(svc DiagramServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	diagramServiceMethods := v1.File_workflows_v1_diagram_proto.Services().ByName("DiagramService").Methods()
	diagramServiceRenderHandler := connect.NewUnaryHandler(
		DiagramServiceRenderProcedure,
		svc.Render,
		connect.WithSchema(diagramServiceMethods.ByName("Render")),
		connect.WithHandlerOptions(opts...),
	)
	return "/workflows.v1.DiagramService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case DiagramServiceRenderProcedure:
			diagramServiceRenderHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedDiagramServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedDiagramServiceHandler struct{}

func (UnimplementedDiagramServiceHandler) Render(context.Context, *connect.Request[v1.RenderDiagramRequest]) (*connect.Response[v1.Diagram], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("workflows.v1.DiagramService.Render is not implemented"))
}
