// Proto messages and service definition for the APIs related to querying data.

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: datasets/v1/data_access.proto

package datasetsv1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	DataAccessService_GetDatasetForInterval_FullMethodName = "/datasets.v1.DataAccessService/GetDatasetForInterval"
	DataAccessService_GetDatapointByID_FullMethodName      = "/datasets.v1.DataAccessService/GetDatapointByID"
	DataAccessService_QueryByID_FullMethodName             = "/datasets.v1.DataAccessService/QueryByID"
	DataAccessService_Query_FullMethodName                 = "/datasets.v1.DataAccessService/Query"
)

// DataAccessServiceClient is the client API for DataAccessService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// DataAccessService provides data access and querying capabilities for Tilebox datasets.
type DataAccessServiceClient interface {
	// GetDatasetForInterval returns a list of data points for a given time interval and collection.
	GetDatasetForInterval(ctx context.Context, in *GetDatasetForIntervalRequest, opts ...grpc.CallOption) (*DatapointPage, error)
	// GetDatapointByID returns a single datapoint by its ID.
	GetDatapointByID(ctx context.Context, in *GetDatapointByIdRequest, opts ...grpc.CallOption) (*Datapoint, error)
	// QueryByID returns a single data point by its ID.
	QueryByID(ctx context.Context, in *QueryByIDRequest, opts ...grpc.CallOption) (*Any, error)
	// Query returns a list of data points matching the given query filters.
	Query(ctx context.Context, in *QueryRequest, opts ...grpc.CallOption) (*QueryResultPage, error)
}

type dataAccessServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDataAccessServiceClient(cc grpc.ClientConnInterface) DataAccessServiceClient {
	return &dataAccessServiceClient{cc}
}

func (c *dataAccessServiceClient) GetDatasetForInterval(ctx context.Context, in *GetDatasetForIntervalRequest, opts ...grpc.CallOption) (*DatapointPage, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DatapointPage)
	err := c.cc.Invoke(ctx, DataAccessService_GetDatasetForInterval_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dataAccessServiceClient) GetDatapointByID(ctx context.Context, in *GetDatapointByIdRequest, opts ...grpc.CallOption) (*Datapoint, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Datapoint)
	err := c.cc.Invoke(ctx, DataAccessService_GetDatapointByID_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dataAccessServiceClient) QueryByID(ctx context.Context, in *QueryByIDRequest, opts ...grpc.CallOption) (*Any, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Any)
	err := c.cc.Invoke(ctx, DataAccessService_QueryByID_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dataAccessServiceClient) Query(ctx context.Context, in *QueryRequest, opts ...grpc.CallOption) (*QueryResultPage, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(QueryResultPage)
	err := c.cc.Invoke(ctx, DataAccessService_Query_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DataAccessServiceServer is the server API for DataAccessService service.
// All implementations must embed UnimplementedDataAccessServiceServer
// for forward compatibility.
//
// DataAccessService provides data access and querying capabilities for Tilebox datasets.
type DataAccessServiceServer interface {
	// GetDatasetForInterval returns a list of data points for a given time interval and collection.
	GetDatasetForInterval(context.Context, *GetDatasetForIntervalRequest) (*DatapointPage, error)
	// GetDatapointByID returns a single datapoint by its ID.
	GetDatapointByID(context.Context, *GetDatapointByIdRequest) (*Datapoint, error)
	// QueryByID returns a single data point by its ID.
	QueryByID(context.Context, *QueryByIDRequest) (*Any, error)
	// Query returns a list of data points matching the given query filters.
	Query(context.Context, *QueryRequest) (*QueryResultPage, error)
	mustEmbedUnimplementedDataAccessServiceServer()
}

// UnimplementedDataAccessServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedDataAccessServiceServer struct{}

func (UnimplementedDataAccessServiceServer) GetDatasetForInterval(context.Context, *GetDatasetForIntervalRequest) (*DatapointPage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDatasetForInterval not implemented")
}
func (UnimplementedDataAccessServiceServer) GetDatapointByID(context.Context, *GetDatapointByIdRequest) (*Datapoint, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDatapointByID not implemented")
}
func (UnimplementedDataAccessServiceServer) QueryByID(context.Context, *QueryByIDRequest) (*Any, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryByID not implemented")
}
func (UnimplementedDataAccessServiceServer) Query(context.Context, *QueryRequest) (*QueryResultPage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Query not implemented")
}
func (UnimplementedDataAccessServiceServer) mustEmbedUnimplementedDataAccessServiceServer() {}
func (UnimplementedDataAccessServiceServer) testEmbeddedByValue()                           {}

// UnsafeDataAccessServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DataAccessServiceServer will
// result in compilation errors.
type UnsafeDataAccessServiceServer interface {
	mustEmbedUnimplementedDataAccessServiceServer()
}

func RegisterDataAccessServiceServer(s grpc.ServiceRegistrar, srv DataAccessServiceServer) {
	// If the following call pancis, it indicates UnimplementedDataAccessServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&DataAccessService_ServiceDesc, srv)
}

func _DataAccessService_GetDatasetForInterval_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDatasetForIntervalRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataAccessServiceServer).GetDatasetForInterval(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DataAccessService_GetDatasetForInterval_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataAccessServiceServer).GetDatasetForInterval(ctx, req.(*GetDatasetForIntervalRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DataAccessService_GetDatapointByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDatapointByIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataAccessServiceServer).GetDatapointByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DataAccessService_GetDatapointByID_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataAccessServiceServer).GetDatapointByID(ctx, req.(*GetDatapointByIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DataAccessService_QueryByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryByIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataAccessServiceServer).QueryByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DataAccessService_QueryByID_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataAccessServiceServer).QueryByID(ctx, req.(*QueryByIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DataAccessService_Query_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataAccessServiceServer).Query(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DataAccessService_Query_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataAccessServiceServer).Query(ctx, req.(*QueryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// DataAccessService_ServiceDesc is the grpc.ServiceDesc for DataAccessService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DataAccessService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "datasets.v1.DataAccessService",
	HandlerType: (*DataAccessServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetDatasetForInterval",
			Handler:    _DataAccessService_GetDatasetForInterval_Handler,
		},
		{
			MethodName: "GetDatapointByID",
			Handler:    _DataAccessService_GetDatapointByID_Handler,
		},
		{
			MethodName: "QueryByID",
			Handler:    _DataAccessService_QueryByID_Handler,
		},
		{
			MethodName: "Query",
			Handler:    _DataAccessService_Query_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "datasets/v1/data_access.proto",
}
