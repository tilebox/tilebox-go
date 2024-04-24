// The internal API a task runner uses to communicate with a workflows-service.

// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: workflows/v1/task.proto

package workflowsv1connect

import (
	connect "connectrpc.com/connect"
	context "context"
	errors "errors"
	v1 "github.com/tilebox/tilebox-go/protogen/go/workflows/v1"
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
	// TaskServiceName is the fully-qualified name of the TaskService service.
	TaskServiceName = "workflows.v1.TaskService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// TaskServiceNextTaskProcedure is the fully-qualified name of the TaskService's NextTask RPC.
	TaskServiceNextTaskProcedure = "/workflows.v1.TaskService/NextTask"
	// TaskServiceTaskFailedProcedure is the fully-qualified name of the TaskService's TaskFailed RPC.
	TaskServiceTaskFailedProcedure = "/workflows.v1.TaskService/TaskFailed"
	// TaskServiceExtendTaskLeaseProcedure is the fully-qualified name of the TaskService's
	// ExtendTaskLease RPC.
	TaskServiceExtendTaskLeaseProcedure = "/workflows.v1.TaskService/ExtendTaskLease"
)

// These variables are the protoreflect.Descriptor objects for the RPCs defined in this package.
var (
	taskServiceServiceDescriptor               = v1.File_workflows_v1_task_proto.Services().ByName("TaskService")
	taskServiceNextTaskMethodDescriptor        = taskServiceServiceDescriptor.Methods().ByName("NextTask")
	taskServiceTaskFailedMethodDescriptor      = taskServiceServiceDescriptor.Methods().ByName("TaskFailed")
	taskServiceExtendTaskLeaseMethodDescriptor = taskServiceServiceDescriptor.Methods().ByName("ExtendTaskLease")
)

// TaskServiceClient is a client for the workflows.v1.TaskService service.
type TaskServiceClient interface {
	// NextTask marks a task as computed and asks for the next task to run.
	// If no task marked as computed is sent, it is assumed that the task runner just started up or was idling so
	// the task server will send a task to run using a work-stealing algorithm.
	// If a task marked as computed is sent, the task server will send a next task to run using a depth first execution
	// algorithm, and only fall back to work-stealing if otherwise no tasks are available.
	// If the next_task_to_run field of the request is not set, a next task will never be returned, but a task
	// can still be marked as computed this way.
	NextTask(context.Context, *connect.Request[v1.NextTaskRequest]) (*connect.Response[v1.NextTaskResponse], error)
	// TaskFailed tells the task server that we have failed to compute a task.
	// The task server will then mark the task as queued or failed, depending on the retry policy,
	// and possibly cancel the job.
	// If a task runner wants to continue executing tasks, it should afterwards fetch a new one using GetTaskToRun.
	TaskFailed(context.Context, *connect.Request[v1.TaskFailedRequest]) (*connect.Response[v1.TaskStateResponse], error)
	// ExtendTaskLease is called by the task runner to extend the lease on a task.
	// On success, the response will contain the new lease expiration time.
	// If the task does not need to be extended, the response will be empty.
	ExtendTaskLease(context.Context, *connect.Request[v1.TaskLeaseRequest]) (*connect.Response[v1.TaskLease], error)
}

// NewTaskServiceClient constructs a client for the workflows.v1.TaskService service. By default, it
// uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses, and sends
// uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC() or
// connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewTaskServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) TaskServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &taskServiceClient{
		nextTask: connect.NewClient[v1.NextTaskRequest, v1.NextTaskResponse](
			httpClient,
			baseURL+TaskServiceNextTaskProcedure,
			connect.WithSchema(taskServiceNextTaskMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		taskFailed: connect.NewClient[v1.TaskFailedRequest, v1.TaskStateResponse](
			httpClient,
			baseURL+TaskServiceTaskFailedProcedure,
			connect.WithSchema(taskServiceTaskFailedMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		extendTaskLease: connect.NewClient[v1.TaskLeaseRequest, v1.TaskLease](
			httpClient,
			baseURL+TaskServiceExtendTaskLeaseProcedure,
			connect.WithSchema(taskServiceExtendTaskLeaseMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
	}
}

// taskServiceClient implements TaskServiceClient.
type taskServiceClient struct {
	nextTask        *connect.Client[v1.NextTaskRequest, v1.NextTaskResponse]
	taskFailed      *connect.Client[v1.TaskFailedRequest, v1.TaskStateResponse]
	extendTaskLease *connect.Client[v1.TaskLeaseRequest, v1.TaskLease]
}

// NextTask calls workflows.v1.TaskService.NextTask.
func (c *taskServiceClient) NextTask(ctx context.Context, req *connect.Request[v1.NextTaskRequest]) (*connect.Response[v1.NextTaskResponse], error) {
	return c.nextTask.CallUnary(ctx, req)
}

// TaskFailed calls workflows.v1.TaskService.TaskFailed.
func (c *taskServiceClient) TaskFailed(ctx context.Context, req *connect.Request[v1.TaskFailedRequest]) (*connect.Response[v1.TaskStateResponse], error) {
	return c.taskFailed.CallUnary(ctx, req)
}

// ExtendTaskLease calls workflows.v1.TaskService.ExtendTaskLease.
func (c *taskServiceClient) ExtendTaskLease(ctx context.Context, req *connect.Request[v1.TaskLeaseRequest]) (*connect.Response[v1.TaskLease], error) {
	return c.extendTaskLease.CallUnary(ctx, req)
}

// TaskServiceHandler is an implementation of the workflows.v1.TaskService service.
type TaskServiceHandler interface {
	// NextTask marks a task as computed and asks for the next task to run.
	// If no task marked as computed is sent, it is assumed that the task runner just started up or was idling so
	// the task server will send a task to run using a work-stealing algorithm.
	// If a task marked as computed is sent, the task server will send a next task to run using a depth first execution
	// algorithm, and only fall back to work-stealing if otherwise no tasks are available.
	// If the next_task_to_run field of the request is not set, a next task will never be returned, but a task
	// can still be marked as computed this way.
	NextTask(context.Context, *connect.Request[v1.NextTaskRequest]) (*connect.Response[v1.NextTaskResponse], error)
	// TaskFailed tells the task server that we have failed to compute a task.
	// The task server will then mark the task as queued or failed, depending on the retry policy,
	// and possibly cancel the job.
	// If a task runner wants to continue executing tasks, it should afterwards fetch a new one using GetTaskToRun.
	TaskFailed(context.Context, *connect.Request[v1.TaskFailedRequest]) (*connect.Response[v1.TaskStateResponse], error)
	// ExtendTaskLease is called by the task runner to extend the lease on a task.
	// On success, the response will contain the new lease expiration time.
	// If the task does not need to be extended, the response will be empty.
	ExtendTaskLease(context.Context, *connect.Request[v1.TaskLeaseRequest]) (*connect.Response[v1.TaskLease], error)
}

// NewTaskServiceHandler builds an HTTP handler from the service implementation. It returns the path
// on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewTaskServiceHandler(svc TaskServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	taskServiceNextTaskHandler := connect.NewUnaryHandler(
		TaskServiceNextTaskProcedure,
		svc.NextTask,
		connect.WithSchema(taskServiceNextTaskMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	taskServiceTaskFailedHandler := connect.NewUnaryHandler(
		TaskServiceTaskFailedProcedure,
		svc.TaskFailed,
		connect.WithSchema(taskServiceTaskFailedMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	taskServiceExtendTaskLeaseHandler := connect.NewUnaryHandler(
		TaskServiceExtendTaskLeaseProcedure,
		svc.ExtendTaskLease,
		connect.WithSchema(taskServiceExtendTaskLeaseMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	return "/workflows.v1.TaskService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case TaskServiceNextTaskProcedure:
			taskServiceNextTaskHandler.ServeHTTP(w, r)
		case TaskServiceTaskFailedProcedure:
			taskServiceTaskFailedHandler.ServeHTTP(w, r)
		case TaskServiceExtendTaskLeaseProcedure:
			taskServiceExtendTaskLeaseHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedTaskServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedTaskServiceHandler struct{}

func (UnimplementedTaskServiceHandler) NextTask(context.Context, *connect.Request[v1.NextTaskRequest]) (*connect.Response[v1.NextTaskResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("workflows.v1.TaskService.NextTask is not implemented"))
}

func (UnimplementedTaskServiceHandler) TaskFailed(context.Context, *connect.Request[v1.TaskFailedRequest]) (*connect.Response[v1.TaskStateResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("workflows.v1.TaskService.TaskFailed is not implemented"))
}

func (UnimplementedTaskServiceHandler) ExtendTaskLease(context.Context, *connect.Request[v1.TaskLeaseRequest]) (*connect.Response[v1.TaskLease], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("workflows.v1.TaskService.ExtendTaskLease is not implemented"))
}
