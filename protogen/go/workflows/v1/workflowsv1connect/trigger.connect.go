// The external API for managing triggers in the Workflows service.

// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: workflows/v1/trigger.proto

package workflowsv1connect

import (
	connect "connectrpc.com/connect"
	context "context"
	errors "errors"
	v1 "github.com/tilebox/tilebox-go/protogen/go/workflows/v1"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
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
	// TriggerServiceName is the fully-qualified name of the TriggerService service.
	TriggerServiceName = "workflows.v1.TriggerService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// TriggerServiceListBucketsProcedure is the fully-qualified name of the TriggerService's
	// ListBuckets RPC.
	TriggerServiceListBucketsProcedure = "/workflows.v1.TriggerService/ListBuckets"
	// TriggerServiceListRecurrentTasksProcedure is the fully-qualified name of the TriggerService's
	// ListRecurrentTasks RPC.
	TriggerServiceListRecurrentTasksProcedure = "/workflows.v1.TriggerService/ListRecurrentTasks"
	// TriggerServiceGetRecurrentTaskProcedure is the fully-qualified name of the TriggerService's
	// GetRecurrentTask RPC.
	TriggerServiceGetRecurrentTaskProcedure = "/workflows.v1.TriggerService/GetRecurrentTask"
	// TriggerServiceCreateRecurrentTaskProcedure is the fully-qualified name of the TriggerService's
	// CreateRecurrentTask RPC.
	TriggerServiceCreateRecurrentTaskProcedure = "/workflows.v1.TriggerService/CreateRecurrentTask"
	// TriggerServiceUpdateRecurrentTaskProcedure is the fully-qualified name of the TriggerService's
	// UpdateRecurrentTask RPC.
	TriggerServiceUpdateRecurrentTaskProcedure = "/workflows.v1.TriggerService/UpdateRecurrentTask"
	// TriggerServiceDeleteRecurrentTaskProcedure is the fully-qualified name of the TriggerService's
	// DeleteRecurrentTask RPC.
	TriggerServiceDeleteRecurrentTaskProcedure = "/workflows.v1.TriggerService/DeleteRecurrentTask"
)

// These variables are the protoreflect.Descriptor objects for the RPCs defined in this package.
var (
	triggerServiceServiceDescriptor                   = v1.File_workflows_v1_trigger_proto.Services().ByName("TriggerService")
	triggerServiceListBucketsMethodDescriptor         = triggerServiceServiceDescriptor.Methods().ByName("ListBuckets")
	triggerServiceListRecurrentTasksMethodDescriptor  = triggerServiceServiceDescriptor.Methods().ByName("ListRecurrentTasks")
	triggerServiceGetRecurrentTaskMethodDescriptor    = triggerServiceServiceDescriptor.Methods().ByName("GetRecurrentTask")
	triggerServiceCreateRecurrentTaskMethodDescriptor = triggerServiceServiceDescriptor.Methods().ByName("CreateRecurrentTask")
	triggerServiceUpdateRecurrentTaskMethodDescriptor = triggerServiceServiceDescriptor.Methods().ByName("UpdateRecurrentTask")
	triggerServiceDeleteRecurrentTaskMethodDescriptor = triggerServiceServiceDescriptor.Methods().ByName("DeleteRecurrentTask")
)

// TriggerServiceClient is a client for the workflows.v1.TriggerService service.
type TriggerServiceClient interface {
	// ListBuckets lists all the storage buckets that are available for use as bucket triggers.
	ListBuckets(context.Context, *connect.Request[emptypb.Empty]) (*connect.Response[v1.Buckets], error)
	// ListRecurrentTasks lists all the recurrent tasks that are currently registered in a namespace.
	ListRecurrentTasks(context.Context, *connect.Request[emptypb.Empty]) (*connect.Response[v1.RecurrentTask], error)
	// GetRecurrentTask gets a recurrent task by its ID.
	GetRecurrentTask(context.Context, *connect.Request[v1.UUID]) (*connect.Response[v1.RecurrentTask], error)
	// CreateRecurrentTask creates a new recurrent task in a namespace.
	CreateRecurrentTask(context.Context, *connect.Request[v1.RecurrentTask]) (*connect.Response[v1.RecurrentTask], error)
	// UpdateRecurrentTask updates a recurrent task in a namespace.
	UpdateRecurrentTask(context.Context, *connect.Request[v1.RecurrentTask]) (*connect.Response[v1.RecurrentTask], error)
	// DeleteRecurrentTask deletes a recurrent task from a namespace.
	DeleteRecurrentTask(context.Context, *connect.Request[v1.UUID]) (*connect.Response[emptypb.Empty], error)
}

// NewTriggerServiceClient constructs a client for the workflows.v1.TriggerService service. By
// default, it uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses,
// and sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the
// connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewTriggerServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) TriggerServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &triggerServiceClient{
		listBuckets: connect.NewClient[emptypb.Empty, v1.Buckets](
			httpClient,
			baseURL+TriggerServiceListBucketsProcedure,
			connect.WithSchema(triggerServiceListBucketsMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		listRecurrentTasks: connect.NewClient[emptypb.Empty, v1.RecurrentTask](
			httpClient,
			baseURL+TriggerServiceListRecurrentTasksProcedure,
			connect.WithSchema(triggerServiceListRecurrentTasksMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		getRecurrentTask: connect.NewClient[v1.UUID, v1.RecurrentTask](
			httpClient,
			baseURL+TriggerServiceGetRecurrentTaskProcedure,
			connect.WithSchema(triggerServiceGetRecurrentTaskMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		createRecurrentTask: connect.NewClient[v1.RecurrentTask, v1.RecurrentTask](
			httpClient,
			baseURL+TriggerServiceCreateRecurrentTaskProcedure,
			connect.WithSchema(triggerServiceCreateRecurrentTaskMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		updateRecurrentTask: connect.NewClient[v1.RecurrentTask, v1.RecurrentTask](
			httpClient,
			baseURL+TriggerServiceUpdateRecurrentTaskProcedure,
			connect.WithSchema(triggerServiceUpdateRecurrentTaskMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		deleteRecurrentTask: connect.NewClient[v1.UUID, emptypb.Empty](
			httpClient,
			baseURL+TriggerServiceDeleteRecurrentTaskProcedure,
			connect.WithSchema(triggerServiceDeleteRecurrentTaskMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
	}
}

// triggerServiceClient implements TriggerServiceClient.
type triggerServiceClient struct {
	listBuckets         *connect.Client[emptypb.Empty, v1.Buckets]
	listRecurrentTasks  *connect.Client[emptypb.Empty, v1.RecurrentTask]
	getRecurrentTask    *connect.Client[v1.UUID, v1.RecurrentTask]
	createRecurrentTask *connect.Client[v1.RecurrentTask, v1.RecurrentTask]
	updateRecurrentTask *connect.Client[v1.RecurrentTask, v1.RecurrentTask]
	deleteRecurrentTask *connect.Client[v1.UUID, emptypb.Empty]
}

// ListBuckets calls workflows.v1.TriggerService.ListBuckets.
func (c *triggerServiceClient) ListBuckets(ctx context.Context, req *connect.Request[emptypb.Empty]) (*connect.Response[v1.Buckets], error) {
	return c.listBuckets.CallUnary(ctx, req)
}

// ListRecurrentTasks calls workflows.v1.TriggerService.ListRecurrentTasks.
func (c *triggerServiceClient) ListRecurrentTasks(ctx context.Context, req *connect.Request[emptypb.Empty]) (*connect.Response[v1.RecurrentTask], error) {
	return c.listRecurrentTasks.CallUnary(ctx, req)
}

// GetRecurrentTask calls workflows.v1.TriggerService.GetRecurrentTask.
func (c *triggerServiceClient) GetRecurrentTask(ctx context.Context, req *connect.Request[v1.UUID]) (*connect.Response[v1.RecurrentTask], error) {
	return c.getRecurrentTask.CallUnary(ctx, req)
}

// CreateRecurrentTask calls workflows.v1.TriggerService.CreateRecurrentTask.
func (c *triggerServiceClient) CreateRecurrentTask(ctx context.Context, req *connect.Request[v1.RecurrentTask]) (*connect.Response[v1.RecurrentTask], error) {
	return c.createRecurrentTask.CallUnary(ctx, req)
}

// UpdateRecurrentTask calls workflows.v1.TriggerService.UpdateRecurrentTask.
func (c *triggerServiceClient) UpdateRecurrentTask(ctx context.Context, req *connect.Request[v1.RecurrentTask]) (*connect.Response[v1.RecurrentTask], error) {
	return c.updateRecurrentTask.CallUnary(ctx, req)
}

// DeleteRecurrentTask calls workflows.v1.TriggerService.DeleteRecurrentTask.
func (c *triggerServiceClient) DeleteRecurrentTask(ctx context.Context, req *connect.Request[v1.UUID]) (*connect.Response[emptypb.Empty], error) {
	return c.deleteRecurrentTask.CallUnary(ctx, req)
}

// TriggerServiceHandler is an implementation of the workflows.v1.TriggerService service.
type TriggerServiceHandler interface {
	// ListBuckets lists all the storage buckets that are available for use as bucket triggers.
	ListBuckets(context.Context, *connect.Request[emptypb.Empty]) (*connect.Response[v1.Buckets], error)
	// ListRecurrentTasks lists all the recurrent tasks that are currently registered in a namespace.
	ListRecurrentTasks(context.Context, *connect.Request[emptypb.Empty]) (*connect.Response[v1.RecurrentTask], error)
	// GetRecurrentTask gets a recurrent task by its ID.
	GetRecurrentTask(context.Context, *connect.Request[v1.UUID]) (*connect.Response[v1.RecurrentTask], error)
	// CreateRecurrentTask creates a new recurrent task in a namespace.
	CreateRecurrentTask(context.Context, *connect.Request[v1.RecurrentTask]) (*connect.Response[v1.RecurrentTask], error)
	// UpdateRecurrentTask updates a recurrent task in a namespace.
	UpdateRecurrentTask(context.Context, *connect.Request[v1.RecurrentTask]) (*connect.Response[v1.RecurrentTask], error)
	// DeleteRecurrentTask deletes a recurrent task from a namespace.
	DeleteRecurrentTask(context.Context, *connect.Request[v1.UUID]) (*connect.Response[emptypb.Empty], error)
}

// NewTriggerServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewTriggerServiceHandler(svc TriggerServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	triggerServiceListBucketsHandler := connect.NewUnaryHandler(
		TriggerServiceListBucketsProcedure,
		svc.ListBuckets,
		connect.WithSchema(triggerServiceListBucketsMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	triggerServiceListRecurrentTasksHandler := connect.NewUnaryHandler(
		TriggerServiceListRecurrentTasksProcedure,
		svc.ListRecurrentTasks,
		connect.WithSchema(triggerServiceListRecurrentTasksMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	triggerServiceGetRecurrentTaskHandler := connect.NewUnaryHandler(
		TriggerServiceGetRecurrentTaskProcedure,
		svc.GetRecurrentTask,
		connect.WithSchema(triggerServiceGetRecurrentTaskMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	triggerServiceCreateRecurrentTaskHandler := connect.NewUnaryHandler(
		TriggerServiceCreateRecurrentTaskProcedure,
		svc.CreateRecurrentTask,
		connect.WithSchema(triggerServiceCreateRecurrentTaskMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	triggerServiceUpdateRecurrentTaskHandler := connect.NewUnaryHandler(
		TriggerServiceUpdateRecurrentTaskProcedure,
		svc.UpdateRecurrentTask,
		connect.WithSchema(triggerServiceUpdateRecurrentTaskMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	triggerServiceDeleteRecurrentTaskHandler := connect.NewUnaryHandler(
		TriggerServiceDeleteRecurrentTaskProcedure,
		svc.DeleteRecurrentTask,
		connect.WithSchema(triggerServiceDeleteRecurrentTaskMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	return "/workflows.v1.TriggerService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case TriggerServiceListBucketsProcedure:
			triggerServiceListBucketsHandler.ServeHTTP(w, r)
		case TriggerServiceListRecurrentTasksProcedure:
			triggerServiceListRecurrentTasksHandler.ServeHTTP(w, r)
		case TriggerServiceGetRecurrentTaskProcedure:
			triggerServiceGetRecurrentTaskHandler.ServeHTTP(w, r)
		case TriggerServiceCreateRecurrentTaskProcedure:
			triggerServiceCreateRecurrentTaskHandler.ServeHTTP(w, r)
		case TriggerServiceUpdateRecurrentTaskProcedure:
			triggerServiceUpdateRecurrentTaskHandler.ServeHTTP(w, r)
		case TriggerServiceDeleteRecurrentTaskProcedure:
			triggerServiceDeleteRecurrentTaskHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedTriggerServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedTriggerServiceHandler struct{}

func (UnimplementedTriggerServiceHandler) ListBuckets(context.Context, *connect.Request[emptypb.Empty]) (*connect.Response[v1.Buckets], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("workflows.v1.TriggerService.ListBuckets is not implemented"))
}

func (UnimplementedTriggerServiceHandler) ListRecurrentTasks(context.Context, *connect.Request[emptypb.Empty]) (*connect.Response[v1.RecurrentTask], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("workflows.v1.TriggerService.ListRecurrentTasks is not implemented"))
}

func (UnimplementedTriggerServiceHandler) GetRecurrentTask(context.Context, *connect.Request[v1.UUID]) (*connect.Response[v1.RecurrentTask], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("workflows.v1.TriggerService.GetRecurrentTask is not implemented"))
}

func (UnimplementedTriggerServiceHandler) CreateRecurrentTask(context.Context, *connect.Request[v1.RecurrentTask]) (*connect.Response[v1.RecurrentTask], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("workflows.v1.TriggerService.CreateRecurrentTask is not implemented"))
}

func (UnimplementedTriggerServiceHandler) UpdateRecurrentTask(context.Context, *connect.Request[v1.RecurrentTask]) (*connect.Response[v1.RecurrentTask], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("workflows.v1.TriggerService.UpdateRecurrentTask is not implemented"))
}

func (UnimplementedTriggerServiceHandler) DeleteRecurrentTask(context.Context, *connect.Request[v1.UUID]) (*connect.Response[emptypb.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("workflows.v1.TriggerService.DeleteRecurrentTask is not implemented"))
}
