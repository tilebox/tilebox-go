// The Tilebox service provides access to datasets.

// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: datasets/v1/tilebox.proto

package datasetsv1connect

import (
	connect "connectrpc.com/connect"
	context "context"
	errors "errors"
	v1 "github.com/tilebox/tilebox-go/protogen/go/datasets/v1"
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
	// TileboxServiceName is the fully-qualified name of the TileboxService service.
	TileboxServiceName = "datasets.v1.TileboxService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// TileboxServiceCreateDatasetProcedure is the fully-qualified name of the TileboxService's
	// CreateDataset RPC.
	TileboxServiceCreateDatasetProcedure = "/datasets.v1.TileboxService/CreateDataset"
	// TileboxServiceGetDatasetProcedure is the fully-qualified name of the TileboxService's GetDataset
	// RPC.
	TileboxServiceGetDatasetProcedure = "/datasets.v1.TileboxService/GetDataset"
	// TileboxServiceUpdateDatasetProcedure is the fully-qualified name of the TileboxService's
	// UpdateDataset RPC.
	TileboxServiceUpdateDatasetProcedure = "/datasets.v1.TileboxService/UpdateDataset"
	// TileboxServiceListDatasetsProcedure is the fully-qualified name of the TileboxService's
	// ListDatasets RPC.
	TileboxServiceListDatasetsProcedure = "/datasets.v1.TileboxService/ListDatasets"
	// TileboxServiceCreateCollectionProcedure is the fully-qualified name of the TileboxService's
	// CreateCollection RPC.
	TileboxServiceCreateCollectionProcedure = "/datasets.v1.TileboxService/CreateCollection"
	// TileboxServiceGetCollectionsProcedure is the fully-qualified name of the TileboxService's
	// GetCollections RPC.
	TileboxServiceGetCollectionsProcedure = "/datasets.v1.TileboxService/GetCollections"
	// TileboxServiceGetCollectionByNameProcedure is the fully-qualified name of the TileboxService's
	// GetCollectionByName RPC.
	TileboxServiceGetCollectionByNameProcedure = "/datasets.v1.TileboxService/GetCollectionByName"
	// TileboxServiceGetDatasetForIntervalProcedure is the fully-qualified name of the TileboxService's
	// GetDatasetForInterval RPC.
	TileboxServiceGetDatasetForIntervalProcedure = "/datasets.v1.TileboxService/GetDatasetForInterval"
	// TileboxServiceGetDatapointByIDProcedure is the fully-qualified name of the TileboxService's
	// GetDatapointByID RPC.
	TileboxServiceGetDatapointByIDProcedure = "/datasets.v1.TileboxService/GetDatapointByID"
	// TileboxServiceIngestDatapointsProcedure is the fully-qualified name of the TileboxService's
	// IngestDatapoints RPC.
	TileboxServiceIngestDatapointsProcedure = "/datasets.v1.TileboxService/IngestDatapoints"
	// TileboxServiceDeleteDatapointsProcedure is the fully-qualified name of the TileboxService's
	// DeleteDatapoints RPC.
	TileboxServiceDeleteDatapointsProcedure = "/datasets.v1.TileboxService/DeleteDatapoints"
)

// TileboxServiceClient is a client for the datasets.v1.TileboxService service.
type TileboxServiceClient interface {
	CreateDataset(context.Context, *connect.Request[v1.CreateDatasetRequest]) (*connect.Response[v1.Dataset], error)
	GetDataset(context.Context, *connect.Request[v1.GetDatasetRequest]) (*connect.Response[v1.Dataset], error)
	UpdateDataset(context.Context, *connect.Request[v1.UpdateDatasetRequest]) (*connect.Response[v1.Dataset], error)
	ListDatasets(context.Context, *connect.Request[v1.ListDatasetsRequest]) (*connect.Response[v1.ListDatasetsResponse], error)
	CreateCollection(context.Context, *connect.Request[v1.CreateCollectionRequest]) (*connect.Response[v1.CollectionInfo], error)
	GetCollections(context.Context, *connect.Request[v1.GetCollectionsRequest]) (*connect.Response[v1.Collections], error)
	GetCollectionByName(context.Context, *connect.Request[v1.GetCollectionByNameRequest]) (*connect.Response[v1.CollectionInfo], error)
	GetDatasetForInterval(context.Context, *connect.Request[v1.GetDatasetForIntervalRequest]) (*connect.Response[v1.Datapoints], error)
	GetDatapointByID(context.Context, *connect.Request[v1.GetDatapointByIdRequest]) (*connect.Response[v1.Datapoint], error)
	IngestDatapoints(context.Context, *connect.Request[v1.IngestDatapointsRequest]) (*connect.Response[v1.IngestDatapointsResponse], error)
	DeleteDatapoints(context.Context, *connect.Request[v1.DeleteDatapointsRequest]) (*connect.Response[v1.DeleteDatapointsResponse], error)
}

// NewTileboxServiceClient constructs a client for the datasets.v1.TileboxService service. By
// default, it uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses,
// and sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the
// connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewTileboxServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) TileboxServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	tileboxServiceMethods := v1.File_datasets_v1_tilebox_proto.Services().ByName("TileboxService").Methods()
	return &tileboxServiceClient{
		createDataset: connect.NewClient[v1.CreateDatasetRequest, v1.Dataset](
			httpClient,
			baseURL+TileboxServiceCreateDatasetProcedure,
			connect.WithSchema(tileboxServiceMethods.ByName("CreateDataset")),
			connect.WithClientOptions(opts...),
		),
		getDataset: connect.NewClient[v1.GetDatasetRequest, v1.Dataset](
			httpClient,
			baseURL+TileboxServiceGetDatasetProcedure,
			connect.WithSchema(tileboxServiceMethods.ByName("GetDataset")),
			connect.WithClientOptions(opts...),
		),
		updateDataset: connect.NewClient[v1.UpdateDatasetRequest, v1.Dataset](
			httpClient,
			baseURL+TileboxServiceUpdateDatasetProcedure,
			connect.WithSchema(tileboxServiceMethods.ByName("UpdateDataset")),
			connect.WithClientOptions(opts...),
		),
		listDatasets: connect.NewClient[v1.ListDatasetsRequest, v1.ListDatasetsResponse](
			httpClient,
			baseURL+TileboxServiceListDatasetsProcedure,
			connect.WithSchema(tileboxServiceMethods.ByName("ListDatasets")),
			connect.WithClientOptions(opts...),
		),
		createCollection: connect.NewClient[v1.CreateCollectionRequest, v1.CollectionInfo](
			httpClient,
			baseURL+TileboxServiceCreateCollectionProcedure,
			connect.WithSchema(tileboxServiceMethods.ByName("CreateCollection")),
			connect.WithClientOptions(opts...),
		),
		getCollections: connect.NewClient[v1.GetCollectionsRequest, v1.Collections](
			httpClient,
			baseURL+TileboxServiceGetCollectionsProcedure,
			connect.WithSchema(tileboxServiceMethods.ByName("GetCollections")),
			connect.WithClientOptions(opts...),
		),
		getCollectionByName: connect.NewClient[v1.GetCollectionByNameRequest, v1.CollectionInfo](
			httpClient,
			baseURL+TileboxServiceGetCollectionByNameProcedure,
			connect.WithSchema(tileboxServiceMethods.ByName("GetCollectionByName")),
			connect.WithClientOptions(opts...),
		),
		getDatasetForInterval: connect.NewClient[v1.GetDatasetForIntervalRequest, v1.Datapoints](
			httpClient,
			baseURL+TileboxServiceGetDatasetForIntervalProcedure,
			connect.WithSchema(tileboxServiceMethods.ByName("GetDatasetForInterval")),
			connect.WithClientOptions(opts...),
		),
		getDatapointByID: connect.NewClient[v1.GetDatapointByIdRequest, v1.Datapoint](
			httpClient,
			baseURL+TileboxServiceGetDatapointByIDProcedure,
			connect.WithSchema(tileboxServiceMethods.ByName("GetDatapointByID")),
			connect.WithClientOptions(opts...),
		),
		ingestDatapoints: connect.NewClient[v1.IngestDatapointsRequest, v1.IngestDatapointsResponse](
			httpClient,
			baseURL+TileboxServiceIngestDatapointsProcedure,
			connect.WithSchema(tileboxServiceMethods.ByName("IngestDatapoints")),
			connect.WithClientOptions(opts...),
		),
		deleteDatapoints: connect.NewClient[v1.DeleteDatapointsRequest, v1.DeleteDatapointsResponse](
			httpClient,
			baseURL+TileboxServiceDeleteDatapointsProcedure,
			connect.WithSchema(tileboxServiceMethods.ByName("DeleteDatapoints")),
			connect.WithClientOptions(opts...),
		),
	}
}

// tileboxServiceClient implements TileboxServiceClient.
type tileboxServiceClient struct {
	createDataset         *connect.Client[v1.CreateDatasetRequest, v1.Dataset]
	getDataset            *connect.Client[v1.GetDatasetRequest, v1.Dataset]
	updateDataset         *connect.Client[v1.UpdateDatasetRequest, v1.Dataset]
	listDatasets          *connect.Client[v1.ListDatasetsRequest, v1.ListDatasetsResponse]
	createCollection      *connect.Client[v1.CreateCollectionRequest, v1.CollectionInfo]
	getCollections        *connect.Client[v1.GetCollectionsRequest, v1.Collections]
	getCollectionByName   *connect.Client[v1.GetCollectionByNameRequest, v1.CollectionInfo]
	getDatasetForInterval *connect.Client[v1.GetDatasetForIntervalRequest, v1.Datapoints]
	getDatapointByID      *connect.Client[v1.GetDatapointByIdRequest, v1.Datapoint]
	ingestDatapoints      *connect.Client[v1.IngestDatapointsRequest, v1.IngestDatapointsResponse]
	deleteDatapoints      *connect.Client[v1.DeleteDatapointsRequest, v1.DeleteDatapointsResponse]
}

// CreateDataset calls datasets.v1.TileboxService.CreateDataset.
func (c *tileboxServiceClient) CreateDataset(ctx context.Context, req *connect.Request[v1.CreateDatasetRequest]) (*connect.Response[v1.Dataset], error) {
	return c.createDataset.CallUnary(ctx, req)
}

// GetDataset calls datasets.v1.TileboxService.GetDataset.
func (c *tileboxServiceClient) GetDataset(ctx context.Context, req *connect.Request[v1.GetDatasetRequest]) (*connect.Response[v1.Dataset], error) {
	return c.getDataset.CallUnary(ctx, req)
}

// UpdateDataset calls datasets.v1.TileboxService.UpdateDataset.
func (c *tileboxServiceClient) UpdateDataset(ctx context.Context, req *connect.Request[v1.UpdateDatasetRequest]) (*connect.Response[v1.Dataset], error) {
	return c.updateDataset.CallUnary(ctx, req)
}

// ListDatasets calls datasets.v1.TileboxService.ListDatasets.
func (c *tileboxServiceClient) ListDatasets(ctx context.Context, req *connect.Request[v1.ListDatasetsRequest]) (*connect.Response[v1.ListDatasetsResponse], error) {
	return c.listDatasets.CallUnary(ctx, req)
}

// CreateCollection calls datasets.v1.TileboxService.CreateCollection.
func (c *tileboxServiceClient) CreateCollection(ctx context.Context, req *connect.Request[v1.CreateCollectionRequest]) (*connect.Response[v1.CollectionInfo], error) {
	return c.createCollection.CallUnary(ctx, req)
}

// GetCollections calls datasets.v1.TileboxService.GetCollections.
func (c *tileboxServiceClient) GetCollections(ctx context.Context, req *connect.Request[v1.GetCollectionsRequest]) (*connect.Response[v1.Collections], error) {
	return c.getCollections.CallUnary(ctx, req)
}

// GetCollectionByName calls datasets.v1.TileboxService.GetCollectionByName.
func (c *tileboxServiceClient) GetCollectionByName(ctx context.Context, req *connect.Request[v1.GetCollectionByNameRequest]) (*connect.Response[v1.CollectionInfo], error) {
	return c.getCollectionByName.CallUnary(ctx, req)
}

// GetDatasetForInterval calls datasets.v1.TileboxService.GetDatasetForInterval.
func (c *tileboxServiceClient) GetDatasetForInterval(ctx context.Context, req *connect.Request[v1.GetDatasetForIntervalRequest]) (*connect.Response[v1.Datapoints], error) {
	return c.getDatasetForInterval.CallUnary(ctx, req)
}

// GetDatapointByID calls datasets.v1.TileboxService.GetDatapointByID.
func (c *tileboxServiceClient) GetDatapointByID(ctx context.Context, req *connect.Request[v1.GetDatapointByIdRequest]) (*connect.Response[v1.Datapoint], error) {
	return c.getDatapointByID.CallUnary(ctx, req)
}

// IngestDatapoints calls datasets.v1.TileboxService.IngestDatapoints.
func (c *tileboxServiceClient) IngestDatapoints(ctx context.Context, req *connect.Request[v1.IngestDatapointsRequest]) (*connect.Response[v1.IngestDatapointsResponse], error) {
	return c.ingestDatapoints.CallUnary(ctx, req)
}

// DeleteDatapoints calls datasets.v1.TileboxService.DeleteDatapoints.
func (c *tileboxServiceClient) DeleteDatapoints(ctx context.Context, req *connect.Request[v1.DeleteDatapointsRequest]) (*connect.Response[v1.DeleteDatapointsResponse], error) {
	return c.deleteDatapoints.CallUnary(ctx, req)
}

// TileboxServiceHandler is an implementation of the datasets.v1.TileboxService service.
type TileboxServiceHandler interface {
	CreateDataset(context.Context, *connect.Request[v1.CreateDatasetRequest]) (*connect.Response[v1.Dataset], error)
	GetDataset(context.Context, *connect.Request[v1.GetDatasetRequest]) (*connect.Response[v1.Dataset], error)
	UpdateDataset(context.Context, *connect.Request[v1.UpdateDatasetRequest]) (*connect.Response[v1.Dataset], error)
	ListDatasets(context.Context, *connect.Request[v1.ListDatasetsRequest]) (*connect.Response[v1.ListDatasetsResponse], error)
	CreateCollection(context.Context, *connect.Request[v1.CreateCollectionRequest]) (*connect.Response[v1.CollectionInfo], error)
	GetCollections(context.Context, *connect.Request[v1.GetCollectionsRequest]) (*connect.Response[v1.Collections], error)
	GetCollectionByName(context.Context, *connect.Request[v1.GetCollectionByNameRequest]) (*connect.Response[v1.CollectionInfo], error)
	GetDatasetForInterval(context.Context, *connect.Request[v1.GetDatasetForIntervalRequest]) (*connect.Response[v1.Datapoints], error)
	GetDatapointByID(context.Context, *connect.Request[v1.GetDatapointByIdRequest]) (*connect.Response[v1.Datapoint], error)
	IngestDatapoints(context.Context, *connect.Request[v1.IngestDatapointsRequest]) (*connect.Response[v1.IngestDatapointsResponse], error)
	DeleteDatapoints(context.Context, *connect.Request[v1.DeleteDatapointsRequest]) (*connect.Response[v1.DeleteDatapointsResponse], error)
}

// NewTileboxServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewTileboxServiceHandler(svc TileboxServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	tileboxServiceMethods := v1.File_datasets_v1_tilebox_proto.Services().ByName("TileboxService").Methods()
	tileboxServiceCreateDatasetHandler := connect.NewUnaryHandler(
		TileboxServiceCreateDatasetProcedure,
		svc.CreateDataset,
		connect.WithSchema(tileboxServiceMethods.ByName("CreateDataset")),
		connect.WithHandlerOptions(opts...),
	)
	tileboxServiceGetDatasetHandler := connect.NewUnaryHandler(
		TileboxServiceGetDatasetProcedure,
		svc.GetDataset,
		connect.WithSchema(tileboxServiceMethods.ByName("GetDataset")),
		connect.WithHandlerOptions(opts...),
	)
	tileboxServiceUpdateDatasetHandler := connect.NewUnaryHandler(
		TileboxServiceUpdateDatasetProcedure,
		svc.UpdateDataset,
		connect.WithSchema(tileboxServiceMethods.ByName("UpdateDataset")),
		connect.WithHandlerOptions(opts...),
	)
	tileboxServiceListDatasetsHandler := connect.NewUnaryHandler(
		TileboxServiceListDatasetsProcedure,
		svc.ListDatasets,
		connect.WithSchema(tileboxServiceMethods.ByName("ListDatasets")),
		connect.WithHandlerOptions(opts...),
	)
	tileboxServiceCreateCollectionHandler := connect.NewUnaryHandler(
		TileboxServiceCreateCollectionProcedure,
		svc.CreateCollection,
		connect.WithSchema(tileboxServiceMethods.ByName("CreateCollection")),
		connect.WithHandlerOptions(opts...),
	)
	tileboxServiceGetCollectionsHandler := connect.NewUnaryHandler(
		TileboxServiceGetCollectionsProcedure,
		svc.GetCollections,
		connect.WithSchema(tileboxServiceMethods.ByName("GetCollections")),
		connect.WithHandlerOptions(opts...),
	)
	tileboxServiceGetCollectionByNameHandler := connect.NewUnaryHandler(
		TileboxServiceGetCollectionByNameProcedure,
		svc.GetCollectionByName,
		connect.WithSchema(tileboxServiceMethods.ByName("GetCollectionByName")),
		connect.WithHandlerOptions(opts...),
	)
	tileboxServiceGetDatasetForIntervalHandler := connect.NewUnaryHandler(
		TileboxServiceGetDatasetForIntervalProcedure,
		svc.GetDatasetForInterval,
		connect.WithSchema(tileboxServiceMethods.ByName("GetDatasetForInterval")),
		connect.WithHandlerOptions(opts...),
	)
	tileboxServiceGetDatapointByIDHandler := connect.NewUnaryHandler(
		TileboxServiceGetDatapointByIDProcedure,
		svc.GetDatapointByID,
		connect.WithSchema(tileboxServiceMethods.ByName("GetDatapointByID")),
		connect.WithHandlerOptions(opts...),
	)
	tileboxServiceIngestDatapointsHandler := connect.NewUnaryHandler(
		TileboxServiceIngestDatapointsProcedure,
		svc.IngestDatapoints,
		connect.WithSchema(tileboxServiceMethods.ByName("IngestDatapoints")),
		connect.WithHandlerOptions(opts...),
	)
	tileboxServiceDeleteDatapointsHandler := connect.NewUnaryHandler(
		TileboxServiceDeleteDatapointsProcedure,
		svc.DeleteDatapoints,
		connect.WithSchema(tileboxServiceMethods.ByName("DeleteDatapoints")),
		connect.WithHandlerOptions(opts...),
	)
	return "/datasets.v1.TileboxService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case TileboxServiceCreateDatasetProcedure:
			tileboxServiceCreateDatasetHandler.ServeHTTP(w, r)
		case TileboxServiceGetDatasetProcedure:
			tileboxServiceGetDatasetHandler.ServeHTTP(w, r)
		case TileboxServiceUpdateDatasetProcedure:
			tileboxServiceUpdateDatasetHandler.ServeHTTP(w, r)
		case TileboxServiceListDatasetsProcedure:
			tileboxServiceListDatasetsHandler.ServeHTTP(w, r)
		case TileboxServiceCreateCollectionProcedure:
			tileboxServiceCreateCollectionHandler.ServeHTTP(w, r)
		case TileboxServiceGetCollectionsProcedure:
			tileboxServiceGetCollectionsHandler.ServeHTTP(w, r)
		case TileboxServiceGetCollectionByNameProcedure:
			tileboxServiceGetCollectionByNameHandler.ServeHTTP(w, r)
		case TileboxServiceGetDatasetForIntervalProcedure:
			tileboxServiceGetDatasetForIntervalHandler.ServeHTTP(w, r)
		case TileboxServiceGetDatapointByIDProcedure:
			tileboxServiceGetDatapointByIDHandler.ServeHTTP(w, r)
		case TileboxServiceIngestDatapointsProcedure:
			tileboxServiceIngestDatapointsHandler.ServeHTTP(w, r)
		case TileboxServiceDeleteDatapointsProcedure:
			tileboxServiceDeleteDatapointsHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedTileboxServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedTileboxServiceHandler struct{}

func (UnimplementedTileboxServiceHandler) CreateDataset(context.Context, *connect.Request[v1.CreateDatasetRequest]) (*connect.Response[v1.Dataset], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("datasets.v1.TileboxService.CreateDataset is not implemented"))
}

func (UnimplementedTileboxServiceHandler) GetDataset(context.Context, *connect.Request[v1.GetDatasetRequest]) (*connect.Response[v1.Dataset], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("datasets.v1.TileboxService.GetDataset is not implemented"))
}

func (UnimplementedTileboxServiceHandler) UpdateDataset(context.Context, *connect.Request[v1.UpdateDatasetRequest]) (*connect.Response[v1.Dataset], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("datasets.v1.TileboxService.UpdateDataset is not implemented"))
}

func (UnimplementedTileboxServiceHandler) ListDatasets(context.Context, *connect.Request[v1.ListDatasetsRequest]) (*connect.Response[v1.ListDatasetsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("datasets.v1.TileboxService.ListDatasets is not implemented"))
}

func (UnimplementedTileboxServiceHandler) CreateCollection(context.Context, *connect.Request[v1.CreateCollectionRequest]) (*connect.Response[v1.CollectionInfo], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("datasets.v1.TileboxService.CreateCollection is not implemented"))
}

func (UnimplementedTileboxServiceHandler) GetCollections(context.Context, *connect.Request[v1.GetCollectionsRequest]) (*connect.Response[v1.Collections], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("datasets.v1.TileboxService.GetCollections is not implemented"))
}

func (UnimplementedTileboxServiceHandler) GetCollectionByName(context.Context, *connect.Request[v1.GetCollectionByNameRequest]) (*connect.Response[v1.CollectionInfo], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("datasets.v1.TileboxService.GetCollectionByName is not implemented"))
}

func (UnimplementedTileboxServiceHandler) GetDatasetForInterval(context.Context, *connect.Request[v1.GetDatasetForIntervalRequest]) (*connect.Response[v1.Datapoints], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("datasets.v1.TileboxService.GetDatasetForInterval is not implemented"))
}

func (UnimplementedTileboxServiceHandler) GetDatapointByID(context.Context, *connect.Request[v1.GetDatapointByIdRequest]) (*connect.Response[v1.Datapoint], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("datasets.v1.TileboxService.GetDatapointByID is not implemented"))
}

func (UnimplementedTileboxServiceHandler) IngestDatapoints(context.Context, *connect.Request[v1.IngestDatapointsRequest]) (*connect.Response[v1.IngestDatapointsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("datasets.v1.TileboxService.IngestDatapoints is not implemented"))
}

func (UnimplementedTileboxServiceHandler) DeleteDatapoints(context.Context, *connect.Request[v1.DeleteDatapointsRequest]) (*connect.Response[v1.DeleteDatapointsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("datasets.v1.TileboxService.DeleteDatapoints is not implemented"))
}
