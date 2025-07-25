// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: workflows/v1/job.proto

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
	// JobServiceName is the fully-qualified name of the JobService service.
	JobServiceName = "workflows.v1.JobService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// JobServiceSubmitJobProcedure is the fully-qualified name of the JobService's SubmitJob RPC.
	JobServiceSubmitJobProcedure = "/workflows.v1.JobService/SubmitJob"
	// JobServiceGetJobProcedure is the fully-qualified name of the JobService's GetJob RPC.
	JobServiceGetJobProcedure = "/workflows.v1.JobService/GetJob"
	// JobServiceRetryJobProcedure is the fully-qualified name of the JobService's RetryJob RPC.
	JobServiceRetryJobProcedure = "/workflows.v1.JobService/RetryJob"
	// JobServiceCancelJobProcedure is the fully-qualified name of the JobService's CancelJob RPC.
	JobServiceCancelJobProcedure = "/workflows.v1.JobService/CancelJob"
	// JobServiceVisualizeJobProcedure is the fully-qualified name of the JobService's VisualizeJob RPC.
	JobServiceVisualizeJobProcedure = "/workflows.v1.JobService/VisualizeJob"
	// JobServiceQueryJobsProcedure is the fully-qualified name of the JobService's QueryJobs RPC.
	JobServiceQueryJobsProcedure = "/workflows.v1.JobService/QueryJobs"
	// JobServiceGetJobPrototypeProcedure is the fully-qualified name of the JobService's
	// GetJobPrototype RPC.
	JobServiceGetJobPrototypeProcedure = "/workflows.v1.JobService/GetJobPrototype"
	// JobServiceCloneJobProcedure is the fully-qualified name of the JobService's CloneJob RPC.
	JobServiceCloneJobProcedure = "/workflows.v1.JobService/CloneJob"
)

// JobServiceClient is a client for the workflows.v1.JobService service.
type JobServiceClient interface {
	SubmitJob(context.Context, *connect.Request[v1.SubmitJobRequest]) (*connect.Response[v1.Job], error)
	GetJob(context.Context, *connect.Request[v1.GetJobRequest]) (*connect.Response[v1.Job], error)
	RetryJob(context.Context, *connect.Request[v1.RetryJobRequest]) (*connect.Response[v1.RetryJobResponse], error)
	CancelJob(context.Context, *connect.Request[v1.CancelJobRequest]) (*connect.Response[v1.CancelJobResponse], error)
	VisualizeJob(context.Context, *connect.Request[v1.VisualizeJobRequest]) (*connect.Response[v1.Diagram], error)
	QueryJobs(context.Context, *connect.Request[v1.QueryJobsRequest]) (*connect.Response[v1.QueryJobsResponse], error)
	GetJobPrototype(context.Context, *connect.Request[v1.GetJobPrototypeRequest]) (*connect.Response[v1.GetJobPrototypeResponse], error)
	CloneJob(context.Context, *connect.Request[v1.CloneJobRequest]) (*connect.Response[v1.Job], error)
}

// NewJobServiceClient constructs a client for the workflows.v1.JobService service. By default, it
// uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses, and sends
// uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC() or
// connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewJobServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) JobServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	jobServiceMethods := v1.File_workflows_v1_job_proto.Services().ByName("JobService").Methods()
	return &jobServiceClient{
		submitJob: connect.NewClient[v1.SubmitJobRequest, v1.Job](
			httpClient,
			baseURL+JobServiceSubmitJobProcedure,
			connect.WithSchema(jobServiceMethods.ByName("SubmitJob")),
			connect.WithClientOptions(opts...),
		),
		getJob: connect.NewClient[v1.GetJobRequest, v1.Job](
			httpClient,
			baseURL+JobServiceGetJobProcedure,
			connect.WithSchema(jobServiceMethods.ByName("GetJob")),
			connect.WithClientOptions(opts...),
		),
		retryJob: connect.NewClient[v1.RetryJobRequest, v1.RetryJobResponse](
			httpClient,
			baseURL+JobServiceRetryJobProcedure,
			connect.WithSchema(jobServiceMethods.ByName("RetryJob")),
			connect.WithClientOptions(opts...),
		),
		cancelJob: connect.NewClient[v1.CancelJobRequest, v1.CancelJobResponse](
			httpClient,
			baseURL+JobServiceCancelJobProcedure,
			connect.WithSchema(jobServiceMethods.ByName("CancelJob")),
			connect.WithClientOptions(opts...),
		),
		visualizeJob: connect.NewClient[v1.VisualizeJobRequest, v1.Diagram](
			httpClient,
			baseURL+JobServiceVisualizeJobProcedure,
			connect.WithSchema(jobServiceMethods.ByName("VisualizeJob")),
			connect.WithClientOptions(opts...),
		),
		queryJobs: connect.NewClient[v1.QueryJobsRequest, v1.QueryJobsResponse](
			httpClient,
			baseURL+JobServiceQueryJobsProcedure,
			connect.WithSchema(jobServiceMethods.ByName("QueryJobs")),
			connect.WithClientOptions(opts...),
		),
		getJobPrototype: connect.NewClient[v1.GetJobPrototypeRequest, v1.GetJobPrototypeResponse](
			httpClient,
			baseURL+JobServiceGetJobPrototypeProcedure,
			connect.WithSchema(jobServiceMethods.ByName("GetJobPrototype")),
			connect.WithClientOptions(opts...),
		),
		cloneJob: connect.NewClient[v1.CloneJobRequest, v1.Job](
			httpClient,
			baseURL+JobServiceCloneJobProcedure,
			connect.WithSchema(jobServiceMethods.ByName("CloneJob")),
			connect.WithClientOptions(opts...),
		),
	}
}

// jobServiceClient implements JobServiceClient.
type jobServiceClient struct {
	submitJob       *connect.Client[v1.SubmitJobRequest, v1.Job]
	getJob          *connect.Client[v1.GetJobRequest, v1.Job]
	retryJob        *connect.Client[v1.RetryJobRequest, v1.RetryJobResponse]
	cancelJob       *connect.Client[v1.CancelJobRequest, v1.CancelJobResponse]
	visualizeJob    *connect.Client[v1.VisualizeJobRequest, v1.Diagram]
	queryJobs       *connect.Client[v1.QueryJobsRequest, v1.QueryJobsResponse]
	getJobPrototype *connect.Client[v1.GetJobPrototypeRequest, v1.GetJobPrototypeResponse]
	cloneJob        *connect.Client[v1.CloneJobRequest, v1.Job]
}

// SubmitJob calls workflows.v1.JobService.SubmitJob.
func (c *jobServiceClient) SubmitJob(ctx context.Context, req *connect.Request[v1.SubmitJobRequest]) (*connect.Response[v1.Job], error) {
	return c.submitJob.CallUnary(ctx, req)
}

// GetJob calls workflows.v1.JobService.GetJob.
func (c *jobServiceClient) GetJob(ctx context.Context, req *connect.Request[v1.GetJobRequest]) (*connect.Response[v1.Job], error) {
	return c.getJob.CallUnary(ctx, req)
}

// RetryJob calls workflows.v1.JobService.RetryJob.
func (c *jobServiceClient) RetryJob(ctx context.Context, req *connect.Request[v1.RetryJobRequest]) (*connect.Response[v1.RetryJobResponse], error) {
	return c.retryJob.CallUnary(ctx, req)
}

// CancelJob calls workflows.v1.JobService.CancelJob.
func (c *jobServiceClient) CancelJob(ctx context.Context, req *connect.Request[v1.CancelJobRequest]) (*connect.Response[v1.CancelJobResponse], error) {
	return c.cancelJob.CallUnary(ctx, req)
}

// VisualizeJob calls workflows.v1.JobService.VisualizeJob.
func (c *jobServiceClient) VisualizeJob(ctx context.Context, req *connect.Request[v1.VisualizeJobRequest]) (*connect.Response[v1.Diagram], error) {
	return c.visualizeJob.CallUnary(ctx, req)
}

// QueryJobs calls workflows.v1.JobService.QueryJobs.
func (c *jobServiceClient) QueryJobs(ctx context.Context, req *connect.Request[v1.QueryJobsRequest]) (*connect.Response[v1.QueryJobsResponse], error) {
	return c.queryJobs.CallUnary(ctx, req)
}

// GetJobPrototype calls workflows.v1.JobService.GetJobPrototype.
func (c *jobServiceClient) GetJobPrototype(ctx context.Context, req *connect.Request[v1.GetJobPrototypeRequest]) (*connect.Response[v1.GetJobPrototypeResponse], error) {
	return c.getJobPrototype.CallUnary(ctx, req)
}

// CloneJob calls workflows.v1.JobService.CloneJob.
func (c *jobServiceClient) CloneJob(ctx context.Context, req *connect.Request[v1.CloneJobRequest]) (*connect.Response[v1.Job], error) {
	return c.cloneJob.CallUnary(ctx, req)
}

// JobServiceHandler is an implementation of the workflows.v1.JobService service.
type JobServiceHandler interface {
	SubmitJob(context.Context, *connect.Request[v1.SubmitJobRequest]) (*connect.Response[v1.Job], error)
	GetJob(context.Context, *connect.Request[v1.GetJobRequest]) (*connect.Response[v1.Job], error)
	RetryJob(context.Context, *connect.Request[v1.RetryJobRequest]) (*connect.Response[v1.RetryJobResponse], error)
	CancelJob(context.Context, *connect.Request[v1.CancelJobRequest]) (*connect.Response[v1.CancelJobResponse], error)
	VisualizeJob(context.Context, *connect.Request[v1.VisualizeJobRequest]) (*connect.Response[v1.Diagram], error)
	QueryJobs(context.Context, *connect.Request[v1.QueryJobsRequest]) (*connect.Response[v1.QueryJobsResponse], error)
	GetJobPrototype(context.Context, *connect.Request[v1.GetJobPrototypeRequest]) (*connect.Response[v1.GetJobPrototypeResponse], error)
	CloneJob(context.Context, *connect.Request[v1.CloneJobRequest]) (*connect.Response[v1.Job], error)
}

// NewJobServiceHandler builds an HTTP handler from the service implementation. It returns the path
// on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewJobServiceHandler(svc JobServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	jobServiceMethods := v1.File_workflows_v1_job_proto.Services().ByName("JobService").Methods()
	jobServiceSubmitJobHandler := connect.NewUnaryHandler(
		JobServiceSubmitJobProcedure,
		svc.SubmitJob,
		connect.WithSchema(jobServiceMethods.ByName("SubmitJob")),
		connect.WithHandlerOptions(opts...),
	)
	jobServiceGetJobHandler := connect.NewUnaryHandler(
		JobServiceGetJobProcedure,
		svc.GetJob,
		connect.WithSchema(jobServiceMethods.ByName("GetJob")),
		connect.WithHandlerOptions(opts...),
	)
	jobServiceRetryJobHandler := connect.NewUnaryHandler(
		JobServiceRetryJobProcedure,
		svc.RetryJob,
		connect.WithSchema(jobServiceMethods.ByName("RetryJob")),
		connect.WithHandlerOptions(opts...),
	)
	jobServiceCancelJobHandler := connect.NewUnaryHandler(
		JobServiceCancelJobProcedure,
		svc.CancelJob,
		connect.WithSchema(jobServiceMethods.ByName("CancelJob")),
		connect.WithHandlerOptions(opts...),
	)
	jobServiceVisualizeJobHandler := connect.NewUnaryHandler(
		JobServiceVisualizeJobProcedure,
		svc.VisualizeJob,
		connect.WithSchema(jobServiceMethods.ByName("VisualizeJob")),
		connect.WithHandlerOptions(opts...),
	)
	jobServiceQueryJobsHandler := connect.NewUnaryHandler(
		JobServiceQueryJobsProcedure,
		svc.QueryJobs,
		connect.WithSchema(jobServiceMethods.ByName("QueryJobs")),
		connect.WithHandlerOptions(opts...),
	)
	jobServiceGetJobPrototypeHandler := connect.NewUnaryHandler(
		JobServiceGetJobPrototypeProcedure,
		svc.GetJobPrototype,
		connect.WithSchema(jobServiceMethods.ByName("GetJobPrototype")),
		connect.WithHandlerOptions(opts...),
	)
	jobServiceCloneJobHandler := connect.NewUnaryHandler(
		JobServiceCloneJobProcedure,
		svc.CloneJob,
		connect.WithSchema(jobServiceMethods.ByName("CloneJob")),
		connect.WithHandlerOptions(opts...),
	)
	return "/workflows.v1.JobService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case JobServiceSubmitJobProcedure:
			jobServiceSubmitJobHandler.ServeHTTP(w, r)
		case JobServiceGetJobProcedure:
			jobServiceGetJobHandler.ServeHTTP(w, r)
		case JobServiceRetryJobProcedure:
			jobServiceRetryJobHandler.ServeHTTP(w, r)
		case JobServiceCancelJobProcedure:
			jobServiceCancelJobHandler.ServeHTTP(w, r)
		case JobServiceVisualizeJobProcedure:
			jobServiceVisualizeJobHandler.ServeHTTP(w, r)
		case JobServiceQueryJobsProcedure:
			jobServiceQueryJobsHandler.ServeHTTP(w, r)
		case JobServiceGetJobPrototypeProcedure:
			jobServiceGetJobPrototypeHandler.ServeHTTP(w, r)
		case JobServiceCloneJobProcedure:
			jobServiceCloneJobHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedJobServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedJobServiceHandler struct{}

func (UnimplementedJobServiceHandler) SubmitJob(context.Context, *connect.Request[v1.SubmitJobRequest]) (*connect.Response[v1.Job], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("workflows.v1.JobService.SubmitJob is not implemented"))
}

func (UnimplementedJobServiceHandler) GetJob(context.Context, *connect.Request[v1.GetJobRequest]) (*connect.Response[v1.Job], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("workflows.v1.JobService.GetJob is not implemented"))
}

func (UnimplementedJobServiceHandler) RetryJob(context.Context, *connect.Request[v1.RetryJobRequest]) (*connect.Response[v1.RetryJobResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("workflows.v1.JobService.RetryJob is not implemented"))
}

func (UnimplementedJobServiceHandler) CancelJob(context.Context, *connect.Request[v1.CancelJobRequest]) (*connect.Response[v1.CancelJobResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("workflows.v1.JobService.CancelJob is not implemented"))
}

func (UnimplementedJobServiceHandler) VisualizeJob(context.Context, *connect.Request[v1.VisualizeJobRequest]) (*connect.Response[v1.Diagram], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("workflows.v1.JobService.VisualizeJob is not implemented"))
}

func (UnimplementedJobServiceHandler) QueryJobs(context.Context, *connect.Request[v1.QueryJobsRequest]) (*connect.Response[v1.QueryJobsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("workflows.v1.JobService.QueryJobs is not implemented"))
}

func (UnimplementedJobServiceHandler) GetJobPrototype(context.Context, *connect.Request[v1.GetJobPrototypeRequest]) (*connect.Response[v1.GetJobPrototypeResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("workflows.v1.JobService.GetJobPrototype is not implemented"))
}

func (UnimplementedJobServiceHandler) CloneJob(context.Context, *connect.Request[v1.CloneJobRequest]) (*connect.Response[v1.Job], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("workflows.v1.JobService.CloneJob is not implemented"))
}
