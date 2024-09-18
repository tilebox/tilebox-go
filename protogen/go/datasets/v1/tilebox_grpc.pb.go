// The Tilebox service provides access to datasets.

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: datasets/v1/tilebox.proto

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
	TileboxService_GetDataset_FullMethodName               = "/datasets.v1.TileboxService/GetDataset"
	TileboxService_UpdateDatasetDescription_FullMethodName = "/datasets.v1.TileboxService/UpdateDatasetDescription"
	TileboxService_ListDatasets_FullMethodName             = "/datasets.v1.TileboxService/ListDatasets"
	TileboxService_GetCollections_FullMethodName           = "/datasets.v1.TileboxService/GetCollections"
	TileboxService_GetCollectionByName_FullMethodName      = "/datasets.v1.TileboxService/GetCollectionByName"
	TileboxService_GetDatasetForInterval_FullMethodName    = "/datasets.v1.TileboxService/GetDatasetForInterval"
	TileboxService_GetDatapointByID_FullMethodName         = "/datasets.v1.TileboxService/GetDatapointByID"
	TileboxService_IngestDatapoints_FullMethodName         = "/datasets.v1.TileboxService/IngestDatapoints"
	TileboxService_DeleteDatapoints_FullMethodName         = "/datasets.v1.TileboxService/DeleteDatapoints"
)

// TileboxServiceClient is the client API for TileboxService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// TileboxService is the service definition for the Tilebox datasets service, which provides access to datasets
type TileboxServiceClient interface {
	GetDataset(ctx context.Context, in *GetDatasetRequest, opts ...grpc.CallOption) (*Dataset, error)
	UpdateDatasetDescription(ctx context.Context, in *UpdateDatasetDescriptionRequest, opts ...grpc.CallOption) (*Dataset, error)
	ListDatasets(ctx context.Context, in *ListDatasetsRequest, opts ...grpc.CallOption) (*ListDatasetsResponse, error)
	GetCollections(ctx context.Context, in *GetCollectionsRequest, opts ...grpc.CallOption) (*Collections, error)
	GetCollectionByName(ctx context.Context, in *GetCollectionByNameRequest, opts ...grpc.CallOption) (*CollectionInfo, error)
	GetDatasetForInterval(ctx context.Context, in *GetDatasetForIntervalRequest, opts ...grpc.CallOption) (*Datapoints, error)
	GetDatapointByID(ctx context.Context, in *GetDatapointByIdRequest, opts ...grpc.CallOption) (*Datapoint, error)
	IngestDatapoints(ctx context.Context, in *IngestDatapointsRequest, opts ...grpc.CallOption) (*IngestDatapointsResponse, error)
	DeleteDatapoints(ctx context.Context, in *DeleteDatapointsRequest, opts ...grpc.CallOption) (*DeleteDatapointsResponse, error)
}

type tileboxServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTileboxServiceClient(cc grpc.ClientConnInterface) TileboxServiceClient {
	return &tileboxServiceClient{cc}
}

func (c *tileboxServiceClient) GetDataset(ctx context.Context, in *GetDatasetRequest, opts ...grpc.CallOption) (*Dataset, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Dataset)
	err := c.cc.Invoke(ctx, TileboxService_GetDataset_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tileboxServiceClient) UpdateDatasetDescription(ctx context.Context, in *UpdateDatasetDescriptionRequest, opts ...grpc.CallOption) (*Dataset, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Dataset)
	err := c.cc.Invoke(ctx, TileboxService_UpdateDatasetDescription_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tileboxServiceClient) ListDatasets(ctx context.Context, in *ListDatasetsRequest, opts ...grpc.CallOption) (*ListDatasetsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListDatasetsResponse)
	err := c.cc.Invoke(ctx, TileboxService_ListDatasets_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tileboxServiceClient) GetCollections(ctx context.Context, in *GetCollectionsRequest, opts ...grpc.CallOption) (*Collections, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Collections)
	err := c.cc.Invoke(ctx, TileboxService_GetCollections_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tileboxServiceClient) GetCollectionByName(ctx context.Context, in *GetCollectionByNameRequest, opts ...grpc.CallOption) (*CollectionInfo, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CollectionInfo)
	err := c.cc.Invoke(ctx, TileboxService_GetCollectionByName_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tileboxServiceClient) GetDatasetForInterval(ctx context.Context, in *GetDatasetForIntervalRequest, opts ...grpc.CallOption) (*Datapoints, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Datapoints)
	err := c.cc.Invoke(ctx, TileboxService_GetDatasetForInterval_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tileboxServiceClient) GetDatapointByID(ctx context.Context, in *GetDatapointByIdRequest, opts ...grpc.CallOption) (*Datapoint, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Datapoint)
	err := c.cc.Invoke(ctx, TileboxService_GetDatapointByID_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tileboxServiceClient) IngestDatapoints(ctx context.Context, in *IngestDatapointsRequest, opts ...grpc.CallOption) (*IngestDatapointsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(IngestDatapointsResponse)
	err := c.cc.Invoke(ctx, TileboxService_IngestDatapoints_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tileboxServiceClient) DeleteDatapoints(ctx context.Context, in *DeleteDatapointsRequest, opts ...grpc.CallOption) (*DeleteDatapointsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteDatapointsResponse)
	err := c.cc.Invoke(ctx, TileboxService_DeleteDatapoints_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TileboxServiceServer is the server API for TileboxService service.
// All implementations must embed UnimplementedTileboxServiceServer
// for forward compatibility.
//
// TileboxService is the service definition for the Tilebox datasets service, which provides access to datasets
type TileboxServiceServer interface {
	GetDataset(context.Context, *GetDatasetRequest) (*Dataset, error)
	UpdateDatasetDescription(context.Context, *UpdateDatasetDescriptionRequest) (*Dataset, error)
	ListDatasets(context.Context, *ListDatasetsRequest) (*ListDatasetsResponse, error)
	GetCollections(context.Context, *GetCollectionsRequest) (*Collections, error)
	GetCollectionByName(context.Context, *GetCollectionByNameRequest) (*CollectionInfo, error)
	GetDatasetForInterval(context.Context, *GetDatasetForIntervalRequest) (*Datapoints, error)
	GetDatapointByID(context.Context, *GetDatapointByIdRequest) (*Datapoint, error)
	IngestDatapoints(context.Context, *IngestDatapointsRequest) (*IngestDatapointsResponse, error)
	DeleteDatapoints(context.Context, *DeleteDatapointsRequest) (*DeleteDatapointsResponse, error)
	mustEmbedUnimplementedTileboxServiceServer()
}

// UnimplementedTileboxServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedTileboxServiceServer struct{}

func (UnimplementedTileboxServiceServer) GetDataset(context.Context, *GetDatasetRequest) (*Dataset, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDataset not implemented")
}
func (UnimplementedTileboxServiceServer) UpdateDatasetDescription(context.Context, *UpdateDatasetDescriptionRequest) (*Dataset, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateDatasetDescription not implemented")
}
func (UnimplementedTileboxServiceServer) ListDatasets(context.Context, *ListDatasetsRequest) (*ListDatasetsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListDatasets not implemented")
}
func (UnimplementedTileboxServiceServer) GetCollections(context.Context, *GetCollectionsRequest) (*Collections, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCollections not implemented")
}
func (UnimplementedTileboxServiceServer) GetCollectionByName(context.Context, *GetCollectionByNameRequest) (*CollectionInfo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCollectionByName not implemented")
}
func (UnimplementedTileboxServiceServer) GetDatasetForInterval(context.Context, *GetDatasetForIntervalRequest) (*Datapoints, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDatasetForInterval not implemented")
}
func (UnimplementedTileboxServiceServer) GetDatapointByID(context.Context, *GetDatapointByIdRequest) (*Datapoint, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDatapointByID not implemented")
}
func (UnimplementedTileboxServiceServer) IngestDatapoints(context.Context, *IngestDatapointsRequest) (*IngestDatapointsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IngestDatapoints not implemented")
}
func (UnimplementedTileboxServiceServer) DeleteDatapoints(context.Context, *DeleteDatapointsRequest) (*DeleteDatapointsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteDatapoints not implemented")
}
func (UnimplementedTileboxServiceServer) mustEmbedUnimplementedTileboxServiceServer() {}
func (UnimplementedTileboxServiceServer) testEmbeddedByValue()                        {}

// UnsafeTileboxServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TileboxServiceServer will
// result in compilation errors.
type UnsafeTileboxServiceServer interface {
	mustEmbedUnimplementedTileboxServiceServer()
}

func RegisterTileboxServiceServer(s grpc.ServiceRegistrar, srv TileboxServiceServer) {
	// If the following call pancis, it indicates UnimplementedTileboxServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&TileboxService_ServiceDesc, srv)
}

func _TileboxService_GetDataset_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDatasetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TileboxServiceServer).GetDataset(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TileboxService_GetDataset_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TileboxServiceServer).GetDataset(ctx, req.(*GetDatasetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TileboxService_UpdateDatasetDescription_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateDatasetDescriptionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TileboxServiceServer).UpdateDatasetDescription(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TileboxService_UpdateDatasetDescription_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TileboxServiceServer).UpdateDatasetDescription(ctx, req.(*UpdateDatasetDescriptionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TileboxService_ListDatasets_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListDatasetsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TileboxServiceServer).ListDatasets(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TileboxService_ListDatasets_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TileboxServiceServer).ListDatasets(ctx, req.(*ListDatasetsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TileboxService_GetCollections_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCollectionsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TileboxServiceServer).GetCollections(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TileboxService_GetCollections_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TileboxServiceServer).GetCollections(ctx, req.(*GetCollectionsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TileboxService_GetCollectionByName_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCollectionByNameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TileboxServiceServer).GetCollectionByName(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TileboxService_GetCollectionByName_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TileboxServiceServer).GetCollectionByName(ctx, req.(*GetCollectionByNameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TileboxService_GetDatasetForInterval_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDatasetForIntervalRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TileboxServiceServer).GetDatasetForInterval(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TileboxService_GetDatasetForInterval_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TileboxServiceServer).GetDatasetForInterval(ctx, req.(*GetDatasetForIntervalRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TileboxService_GetDatapointByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDatapointByIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TileboxServiceServer).GetDatapointByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TileboxService_GetDatapointByID_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TileboxServiceServer).GetDatapointByID(ctx, req.(*GetDatapointByIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TileboxService_IngestDatapoints_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IngestDatapointsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TileboxServiceServer).IngestDatapoints(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TileboxService_IngestDatapoints_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TileboxServiceServer).IngestDatapoints(ctx, req.(*IngestDatapointsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TileboxService_DeleteDatapoints_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteDatapointsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TileboxServiceServer).DeleteDatapoints(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TileboxService_DeleteDatapoints_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TileboxServiceServer).DeleteDatapoints(ctx, req.(*DeleteDatapointsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// TileboxService_ServiceDesc is the grpc.ServiceDesc for TileboxService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TileboxService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "datasets.v1.TileboxService",
	HandlerType: (*TileboxServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetDataset",
			Handler:    _TileboxService_GetDataset_Handler,
		},
		{
			MethodName: "UpdateDatasetDescription",
			Handler:    _TileboxService_UpdateDatasetDescription_Handler,
		},
		{
			MethodName: "ListDatasets",
			Handler:    _TileboxService_ListDatasets_Handler,
		},
		{
			MethodName: "GetCollections",
			Handler:    _TileboxService_GetCollections_Handler,
		},
		{
			MethodName: "GetCollectionByName",
			Handler:    _TileboxService_GetCollectionByName_Handler,
		},
		{
			MethodName: "GetDatasetForInterval",
			Handler:    _TileboxService_GetDatasetForInterval_Handler,
		},
		{
			MethodName: "GetDatapointByID",
			Handler:    _TileboxService_GetDatapointByID_Handler,
		},
		{
			MethodName: "IngestDatapoints",
			Handler:    _TileboxService_IngestDatapoints_Handler,
		},
		{
			MethodName: "DeleteDatapoints",
			Handler:    _TileboxService_DeleteDatapoints_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "datasets/v1/tilebox.proto",
}
