// The external API for managing recurrent tasks in the Workflows service.

// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: workflows/v1/recurrent_task.proto

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
	// RecurrentTaskServiceName is the fully-qualified name of the RecurrentTaskService service.
	RecurrentTaskServiceName = "workflows.v1.RecurrentTaskService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// RecurrentTaskServiceListStorageLocationsProcedure is the fully-qualified name of the
	// RecurrentTaskService's ListStorageLocations RPC.
	RecurrentTaskServiceListStorageLocationsProcedure = "/workflows.v1.RecurrentTaskService/ListStorageLocations"
	// RecurrentTaskServiceGetStorageLocationProcedure is the fully-qualified name of the
	// RecurrentTaskService's GetStorageLocation RPC.
	RecurrentTaskServiceGetStorageLocationProcedure = "/workflows.v1.RecurrentTaskService/GetStorageLocation"
	// RecurrentTaskServiceCreateStorageLocationProcedure is the fully-qualified name of the
	// RecurrentTaskService's CreateStorageLocation RPC.
	RecurrentTaskServiceCreateStorageLocationProcedure = "/workflows.v1.RecurrentTaskService/CreateStorageLocation"
	// RecurrentTaskServiceDeleteStorageLocationProcedure is the fully-qualified name of the
	// RecurrentTaskService's DeleteStorageLocation RPC.
	RecurrentTaskServiceDeleteStorageLocationProcedure = "/workflows.v1.RecurrentTaskService/DeleteStorageLocation"
	// RecurrentTaskServiceListRecurrentTasksProcedure is the fully-qualified name of the
	// RecurrentTaskService's ListRecurrentTasks RPC.
	RecurrentTaskServiceListRecurrentTasksProcedure = "/workflows.v1.RecurrentTaskService/ListRecurrentTasks"
	// RecurrentTaskServiceGetRecurrentTaskProcedure is the fully-qualified name of the
	// RecurrentTaskService's GetRecurrentTask RPC.
	RecurrentTaskServiceGetRecurrentTaskProcedure = "/workflows.v1.RecurrentTaskService/GetRecurrentTask"
	// RecurrentTaskServiceCreateRecurrentTaskProcedure is the fully-qualified name of the
	// RecurrentTaskService's CreateRecurrentTask RPC.
	RecurrentTaskServiceCreateRecurrentTaskProcedure = "/workflows.v1.RecurrentTaskService/CreateRecurrentTask"
	// RecurrentTaskServiceUpdateRecurrentTaskProcedure is the fully-qualified name of the
	// RecurrentTaskService's UpdateRecurrentTask RPC.
	RecurrentTaskServiceUpdateRecurrentTaskProcedure = "/workflows.v1.RecurrentTaskService/UpdateRecurrentTask"
	// RecurrentTaskServiceDeleteRecurrentTaskProcedure is the fully-qualified name of the
	// RecurrentTaskService's DeleteRecurrentTask RPC.
	RecurrentTaskServiceDeleteRecurrentTaskProcedure = "/workflows.v1.RecurrentTaskService/DeleteRecurrentTask"
	// RecurrentTaskServiceDeleteAutomationProcedure is the fully-qualified name of the
	// RecurrentTaskService's DeleteAutomation RPC.
	RecurrentTaskServiceDeleteAutomationProcedure = "/workflows.v1.RecurrentTaskService/DeleteAutomation"
)

// RecurrentTaskServiceClient is a client for the workflows.v1.RecurrentTaskService service.
type RecurrentTaskServiceClient interface {
	// ListStorageLocations lists all the storage buckets that are available for use as bucket triggers.
	ListStorageLocations(context.Context, *connect.Request[emptypb.Empty]) (*connect.Response[v1.StorageLocations], error)
	// GetStorageLocation gets a storage location by its ID.
	GetStorageLocation(context.Context, *connect.Request[v1.UUID]) (*connect.Response[v1.StorageLocation], error)
	// CreateStorageLocation creates a new storage bucket.
	CreateStorageLocation(context.Context, *connect.Request[v1.StorageLocation]) (*connect.Response[v1.StorageLocation], error)
	// DeleteStorageLocation deletes a storage location.
	DeleteStorageLocation(context.Context, *connect.Request[v1.UUID]) (*connect.Response[emptypb.Empty], error)
	// ListRecurrentTasks lists all the recurrent tasks that are currently registered in a namespace.
	ListRecurrentTasks(context.Context, *connect.Request[emptypb.Empty]) (*connect.Response[v1.RecurrentTasks], error)
	// GetRecurrentTask gets a recurrent task by its ID.
	GetRecurrentTask(context.Context, *connect.Request[v1.UUID]) (*connect.Response[v1.RecurrentTaskPrototype], error)
	// CreateRecurrentTask creates a new recurrent task in a namespace.
	CreateRecurrentTask(context.Context, *connect.Request[v1.RecurrentTaskPrototype]) (*connect.Response[v1.RecurrentTaskPrototype], error)
	// UpdateRecurrentTask updates a recurrent task in a namespace.
	UpdateRecurrentTask(context.Context, *connect.Request[v1.RecurrentTaskPrototype]) (*connect.Response[v1.RecurrentTaskPrototype], error)
	// DeleteRecurrentTask deletes a recurrent task from a namespace.
	DeleteRecurrentTask(context.Context, *connect.Request[v1.UUID]) (*connect.Response[emptypb.Empty], error)
	// DeleteAutomation deletes an automation from a namespace.
	DeleteAutomation(context.Context, *connect.Request[v1.DeleteAutomationRequest]) (*connect.Response[emptypb.Empty], error)
}

// NewRecurrentTaskServiceClient constructs a client for the workflows.v1.RecurrentTaskService
// service. By default, it uses the Connect protocol with the binary Protobuf Codec, asks for
// gzipped responses, and sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply
// the connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewRecurrentTaskServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) RecurrentTaskServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	recurrentTaskServiceMethods := v1.File_workflows_v1_recurrent_task_proto.Services().ByName("RecurrentTaskService").Methods()
	return &recurrentTaskServiceClient{
		listStorageLocations: connect.NewClient[emptypb.Empty, v1.StorageLocations](
			httpClient,
			baseURL+RecurrentTaskServiceListStorageLocationsProcedure,
			connect.WithSchema(recurrentTaskServiceMethods.ByName("ListStorageLocations")),
			connect.WithClientOptions(opts...),
		),
		getStorageLocation: connect.NewClient[v1.UUID, v1.StorageLocation](
			httpClient,
			baseURL+RecurrentTaskServiceGetStorageLocationProcedure,
			connect.WithSchema(recurrentTaskServiceMethods.ByName("GetStorageLocation")),
			connect.WithClientOptions(opts...),
		),
		createStorageLocation: connect.NewClient[v1.StorageLocation, v1.StorageLocation](
			httpClient,
			baseURL+RecurrentTaskServiceCreateStorageLocationProcedure,
			connect.WithSchema(recurrentTaskServiceMethods.ByName("CreateStorageLocation")),
			connect.WithClientOptions(opts...),
		),
		deleteStorageLocation: connect.NewClient[v1.UUID, emptypb.Empty](
			httpClient,
			baseURL+RecurrentTaskServiceDeleteStorageLocationProcedure,
			connect.WithSchema(recurrentTaskServiceMethods.ByName("DeleteStorageLocation")),
			connect.WithClientOptions(opts...),
		),
		listRecurrentTasks: connect.NewClient[emptypb.Empty, v1.RecurrentTasks](
			httpClient,
			baseURL+RecurrentTaskServiceListRecurrentTasksProcedure,
			connect.WithSchema(recurrentTaskServiceMethods.ByName("ListRecurrentTasks")),
			connect.WithClientOptions(opts...),
		),
		getRecurrentTask: connect.NewClient[v1.UUID, v1.RecurrentTaskPrototype](
			httpClient,
			baseURL+RecurrentTaskServiceGetRecurrentTaskProcedure,
			connect.WithSchema(recurrentTaskServiceMethods.ByName("GetRecurrentTask")),
			connect.WithClientOptions(opts...),
		),
		createRecurrentTask: connect.NewClient[v1.RecurrentTaskPrototype, v1.RecurrentTaskPrototype](
			httpClient,
			baseURL+RecurrentTaskServiceCreateRecurrentTaskProcedure,
			connect.WithSchema(recurrentTaskServiceMethods.ByName("CreateRecurrentTask")),
			connect.WithClientOptions(opts...),
		),
		updateRecurrentTask: connect.NewClient[v1.RecurrentTaskPrototype, v1.RecurrentTaskPrototype](
			httpClient,
			baseURL+RecurrentTaskServiceUpdateRecurrentTaskProcedure,
			connect.WithSchema(recurrentTaskServiceMethods.ByName("UpdateRecurrentTask")),
			connect.WithClientOptions(opts...),
		),
		deleteRecurrentTask: connect.NewClient[v1.UUID, emptypb.Empty](
			httpClient,
			baseURL+RecurrentTaskServiceDeleteRecurrentTaskProcedure,
			connect.WithSchema(recurrentTaskServiceMethods.ByName("DeleteRecurrentTask")),
			connect.WithClientOptions(opts...),
		),
		deleteAutomation: connect.NewClient[v1.DeleteAutomationRequest, emptypb.Empty](
			httpClient,
			baseURL+RecurrentTaskServiceDeleteAutomationProcedure,
			connect.WithSchema(recurrentTaskServiceMethods.ByName("DeleteAutomation")),
			connect.WithClientOptions(opts...),
		),
	}
}

// recurrentTaskServiceClient implements RecurrentTaskServiceClient.
type recurrentTaskServiceClient struct {
	listStorageLocations  *connect.Client[emptypb.Empty, v1.StorageLocations]
	getStorageLocation    *connect.Client[v1.UUID, v1.StorageLocation]
	createStorageLocation *connect.Client[v1.StorageLocation, v1.StorageLocation]
	deleteStorageLocation *connect.Client[v1.UUID, emptypb.Empty]
	listRecurrentTasks    *connect.Client[emptypb.Empty, v1.RecurrentTasks]
	getRecurrentTask      *connect.Client[v1.UUID, v1.RecurrentTaskPrototype]
	createRecurrentTask   *connect.Client[v1.RecurrentTaskPrototype, v1.RecurrentTaskPrototype]
	updateRecurrentTask   *connect.Client[v1.RecurrentTaskPrototype, v1.RecurrentTaskPrototype]
	deleteRecurrentTask   *connect.Client[v1.UUID, emptypb.Empty]
	deleteAutomation      *connect.Client[v1.DeleteAutomationRequest, emptypb.Empty]
}

// ListStorageLocations calls workflows.v1.RecurrentTaskService.ListStorageLocations.
func (c *recurrentTaskServiceClient) ListStorageLocations(ctx context.Context, req *connect.Request[emptypb.Empty]) (*connect.Response[v1.StorageLocations], error) {
	return c.listStorageLocations.CallUnary(ctx, req)
}

// GetStorageLocation calls workflows.v1.RecurrentTaskService.GetStorageLocation.
func (c *recurrentTaskServiceClient) GetStorageLocation(ctx context.Context, req *connect.Request[v1.UUID]) (*connect.Response[v1.StorageLocation], error) {
	return c.getStorageLocation.CallUnary(ctx, req)
}

// CreateStorageLocation calls workflows.v1.RecurrentTaskService.CreateStorageLocation.
func (c *recurrentTaskServiceClient) CreateStorageLocation(ctx context.Context, req *connect.Request[v1.StorageLocation]) (*connect.Response[v1.StorageLocation], error) {
	return c.createStorageLocation.CallUnary(ctx, req)
}

// DeleteStorageLocation calls workflows.v1.RecurrentTaskService.DeleteStorageLocation.
func (c *recurrentTaskServiceClient) DeleteStorageLocation(ctx context.Context, req *connect.Request[v1.UUID]) (*connect.Response[emptypb.Empty], error) {
	return c.deleteStorageLocation.CallUnary(ctx, req)
}

// ListRecurrentTasks calls workflows.v1.RecurrentTaskService.ListRecurrentTasks.
func (c *recurrentTaskServiceClient) ListRecurrentTasks(ctx context.Context, req *connect.Request[emptypb.Empty]) (*connect.Response[v1.RecurrentTasks], error) {
	return c.listRecurrentTasks.CallUnary(ctx, req)
}

// GetRecurrentTask calls workflows.v1.RecurrentTaskService.GetRecurrentTask.
func (c *recurrentTaskServiceClient) GetRecurrentTask(ctx context.Context, req *connect.Request[v1.UUID]) (*connect.Response[v1.RecurrentTaskPrototype], error) {
	return c.getRecurrentTask.CallUnary(ctx, req)
}

// CreateRecurrentTask calls workflows.v1.RecurrentTaskService.CreateRecurrentTask.
func (c *recurrentTaskServiceClient) CreateRecurrentTask(ctx context.Context, req *connect.Request[v1.RecurrentTaskPrototype]) (*connect.Response[v1.RecurrentTaskPrototype], error) {
	return c.createRecurrentTask.CallUnary(ctx, req)
}

// UpdateRecurrentTask calls workflows.v1.RecurrentTaskService.UpdateRecurrentTask.
func (c *recurrentTaskServiceClient) UpdateRecurrentTask(ctx context.Context, req *connect.Request[v1.RecurrentTaskPrototype]) (*connect.Response[v1.RecurrentTaskPrototype], error) {
	return c.updateRecurrentTask.CallUnary(ctx, req)
}

// DeleteRecurrentTask calls workflows.v1.RecurrentTaskService.DeleteRecurrentTask.
func (c *recurrentTaskServiceClient) DeleteRecurrentTask(ctx context.Context, req *connect.Request[v1.UUID]) (*connect.Response[emptypb.Empty], error) {
	return c.deleteRecurrentTask.CallUnary(ctx, req)
}

// DeleteAutomation calls workflows.v1.RecurrentTaskService.DeleteAutomation.
func (c *recurrentTaskServiceClient) DeleteAutomation(ctx context.Context, req *connect.Request[v1.DeleteAutomationRequest]) (*connect.Response[emptypb.Empty], error) {
	return c.deleteAutomation.CallUnary(ctx, req)
}

// RecurrentTaskServiceHandler is an implementation of the workflows.v1.RecurrentTaskService
// service.
type RecurrentTaskServiceHandler interface {
	// ListStorageLocations lists all the storage buckets that are available for use as bucket triggers.
	ListStorageLocations(context.Context, *connect.Request[emptypb.Empty]) (*connect.Response[v1.StorageLocations], error)
	// GetStorageLocation gets a storage location by its ID.
	GetStorageLocation(context.Context, *connect.Request[v1.UUID]) (*connect.Response[v1.StorageLocation], error)
	// CreateStorageLocation creates a new storage bucket.
	CreateStorageLocation(context.Context, *connect.Request[v1.StorageLocation]) (*connect.Response[v1.StorageLocation], error)
	// DeleteStorageLocation deletes a storage location.
	DeleteStorageLocation(context.Context, *connect.Request[v1.UUID]) (*connect.Response[emptypb.Empty], error)
	// ListRecurrentTasks lists all the recurrent tasks that are currently registered in a namespace.
	ListRecurrentTasks(context.Context, *connect.Request[emptypb.Empty]) (*connect.Response[v1.RecurrentTasks], error)
	// GetRecurrentTask gets a recurrent task by its ID.
	GetRecurrentTask(context.Context, *connect.Request[v1.UUID]) (*connect.Response[v1.RecurrentTaskPrototype], error)
	// CreateRecurrentTask creates a new recurrent task in a namespace.
	CreateRecurrentTask(context.Context, *connect.Request[v1.RecurrentTaskPrototype]) (*connect.Response[v1.RecurrentTaskPrototype], error)
	// UpdateRecurrentTask updates a recurrent task in a namespace.
	UpdateRecurrentTask(context.Context, *connect.Request[v1.RecurrentTaskPrototype]) (*connect.Response[v1.RecurrentTaskPrototype], error)
	// DeleteRecurrentTask deletes a recurrent task from a namespace.
	DeleteRecurrentTask(context.Context, *connect.Request[v1.UUID]) (*connect.Response[emptypb.Empty], error)
	// DeleteAutomation deletes an automation from a namespace.
	DeleteAutomation(context.Context, *connect.Request[v1.DeleteAutomationRequest]) (*connect.Response[emptypb.Empty], error)
}

// NewRecurrentTaskServiceHandler builds an HTTP handler from the service implementation. It returns
// the path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewRecurrentTaskServiceHandler(svc RecurrentTaskServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	recurrentTaskServiceMethods := v1.File_workflows_v1_recurrent_task_proto.Services().ByName("RecurrentTaskService").Methods()
	recurrentTaskServiceListStorageLocationsHandler := connect.NewUnaryHandler(
		RecurrentTaskServiceListStorageLocationsProcedure,
		svc.ListStorageLocations,
		connect.WithSchema(recurrentTaskServiceMethods.ByName("ListStorageLocations")),
		connect.WithHandlerOptions(opts...),
	)
	recurrentTaskServiceGetStorageLocationHandler := connect.NewUnaryHandler(
		RecurrentTaskServiceGetStorageLocationProcedure,
		svc.GetStorageLocation,
		connect.WithSchema(recurrentTaskServiceMethods.ByName("GetStorageLocation")),
		connect.WithHandlerOptions(opts...),
	)
	recurrentTaskServiceCreateStorageLocationHandler := connect.NewUnaryHandler(
		RecurrentTaskServiceCreateStorageLocationProcedure,
		svc.CreateStorageLocation,
		connect.WithSchema(recurrentTaskServiceMethods.ByName("CreateStorageLocation")),
		connect.WithHandlerOptions(opts...),
	)
	recurrentTaskServiceDeleteStorageLocationHandler := connect.NewUnaryHandler(
		RecurrentTaskServiceDeleteStorageLocationProcedure,
		svc.DeleteStorageLocation,
		connect.WithSchema(recurrentTaskServiceMethods.ByName("DeleteStorageLocation")),
		connect.WithHandlerOptions(opts...),
	)
	recurrentTaskServiceListRecurrentTasksHandler := connect.NewUnaryHandler(
		RecurrentTaskServiceListRecurrentTasksProcedure,
		svc.ListRecurrentTasks,
		connect.WithSchema(recurrentTaskServiceMethods.ByName("ListRecurrentTasks")),
		connect.WithHandlerOptions(opts...),
	)
	recurrentTaskServiceGetRecurrentTaskHandler := connect.NewUnaryHandler(
		RecurrentTaskServiceGetRecurrentTaskProcedure,
		svc.GetRecurrentTask,
		connect.WithSchema(recurrentTaskServiceMethods.ByName("GetRecurrentTask")),
		connect.WithHandlerOptions(opts...),
	)
	recurrentTaskServiceCreateRecurrentTaskHandler := connect.NewUnaryHandler(
		RecurrentTaskServiceCreateRecurrentTaskProcedure,
		svc.CreateRecurrentTask,
		connect.WithSchema(recurrentTaskServiceMethods.ByName("CreateRecurrentTask")),
		connect.WithHandlerOptions(opts...),
	)
	recurrentTaskServiceUpdateRecurrentTaskHandler := connect.NewUnaryHandler(
		RecurrentTaskServiceUpdateRecurrentTaskProcedure,
		svc.UpdateRecurrentTask,
		connect.WithSchema(recurrentTaskServiceMethods.ByName("UpdateRecurrentTask")),
		connect.WithHandlerOptions(opts...),
	)
	recurrentTaskServiceDeleteRecurrentTaskHandler := connect.NewUnaryHandler(
		RecurrentTaskServiceDeleteRecurrentTaskProcedure,
		svc.DeleteRecurrentTask,
		connect.WithSchema(recurrentTaskServiceMethods.ByName("DeleteRecurrentTask")),
		connect.WithHandlerOptions(opts...),
	)
	recurrentTaskServiceDeleteAutomationHandler := connect.NewUnaryHandler(
		RecurrentTaskServiceDeleteAutomationProcedure,
		svc.DeleteAutomation,
		connect.WithSchema(recurrentTaskServiceMethods.ByName("DeleteAutomation")),
		connect.WithHandlerOptions(opts...),
	)
	return "/workflows.v1.RecurrentTaskService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case RecurrentTaskServiceListStorageLocationsProcedure:
			recurrentTaskServiceListStorageLocationsHandler.ServeHTTP(w, r)
		case RecurrentTaskServiceGetStorageLocationProcedure:
			recurrentTaskServiceGetStorageLocationHandler.ServeHTTP(w, r)
		case RecurrentTaskServiceCreateStorageLocationProcedure:
			recurrentTaskServiceCreateStorageLocationHandler.ServeHTTP(w, r)
		case RecurrentTaskServiceDeleteStorageLocationProcedure:
			recurrentTaskServiceDeleteStorageLocationHandler.ServeHTTP(w, r)
		case RecurrentTaskServiceListRecurrentTasksProcedure:
			recurrentTaskServiceListRecurrentTasksHandler.ServeHTTP(w, r)
		case RecurrentTaskServiceGetRecurrentTaskProcedure:
			recurrentTaskServiceGetRecurrentTaskHandler.ServeHTTP(w, r)
		case RecurrentTaskServiceCreateRecurrentTaskProcedure:
			recurrentTaskServiceCreateRecurrentTaskHandler.ServeHTTP(w, r)
		case RecurrentTaskServiceUpdateRecurrentTaskProcedure:
			recurrentTaskServiceUpdateRecurrentTaskHandler.ServeHTTP(w, r)
		case RecurrentTaskServiceDeleteRecurrentTaskProcedure:
			recurrentTaskServiceDeleteRecurrentTaskHandler.ServeHTTP(w, r)
		case RecurrentTaskServiceDeleteAutomationProcedure:
			recurrentTaskServiceDeleteAutomationHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedRecurrentTaskServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedRecurrentTaskServiceHandler struct{}

func (UnimplementedRecurrentTaskServiceHandler) ListStorageLocations(context.Context, *connect.Request[emptypb.Empty]) (*connect.Response[v1.StorageLocations], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("workflows.v1.RecurrentTaskService.ListStorageLocations is not implemented"))
}

func (UnimplementedRecurrentTaskServiceHandler) GetStorageLocation(context.Context, *connect.Request[v1.UUID]) (*connect.Response[v1.StorageLocation], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("workflows.v1.RecurrentTaskService.GetStorageLocation is not implemented"))
}

func (UnimplementedRecurrentTaskServiceHandler) CreateStorageLocation(context.Context, *connect.Request[v1.StorageLocation]) (*connect.Response[v1.StorageLocation], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("workflows.v1.RecurrentTaskService.CreateStorageLocation is not implemented"))
}

func (UnimplementedRecurrentTaskServiceHandler) DeleteStorageLocation(context.Context, *connect.Request[v1.UUID]) (*connect.Response[emptypb.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("workflows.v1.RecurrentTaskService.DeleteStorageLocation is not implemented"))
}

func (UnimplementedRecurrentTaskServiceHandler) ListRecurrentTasks(context.Context, *connect.Request[emptypb.Empty]) (*connect.Response[v1.RecurrentTasks], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("workflows.v1.RecurrentTaskService.ListRecurrentTasks is not implemented"))
}

func (UnimplementedRecurrentTaskServiceHandler) GetRecurrentTask(context.Context, *connect.Request[v1.UUID]) (*connect.Response[v1.RecurrentTaskPrototype], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("workflows.v1.RecurrentTaskService.GetRecurrentTask is not implemented"))
}

func (UnimplementedRecurrentTaskServiceHandler) CreateRecurrentTask(context.Context, *connect.Request[v1.RecurrentTaskPrototype]) (*connect.Response[v1.RecurrentTaskPrototype], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("workflows.v1.RecurrentTaskService.CreateRecurrentTask is not implemented"))
}

func (UnimplementedRecurrentTaskServiceHandler) UpdateRecurrentTask(context.Context, *connect.Request[v1.RecurrentTaskPrototype]) (*connect.Response[v1.RecurrentTaskPrototype], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("workflows.v1.RecurrentTaskService.UpdateRecurrentTask is not implemented"))
}

func (UnimplementedRecurrentTaskServiceHandler) DeleteRecurrentTask(context.Context, *connect.Request[v1.UUID]) (*connect.Response[emptypb.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("workflows.v1.RecurrentTaskService.DeleteRecurrentTask is not implemented"))
}

func (UnimplementedRecurrentTaskServiceHandler) DeleteAutomation(context.Context, *connect.Request[v1.DeleteAutomationRequest]) (*connect.Response[emptypb.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("workflows.v1.RecurrentTaskService.DeleteAutomation is not implemented"))
}
