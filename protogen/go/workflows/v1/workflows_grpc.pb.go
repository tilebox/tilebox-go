// The externally facing API allowing users to interact with workflows.

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             (unknown)
// source: workflows/v1/workflows.proto

package workflowsv1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.62.0 or later.
const _ = grpc.SupportPackageIsVersion8

const (
	WorkflowsService_CreateCluster_FullMethodName = "/workflows.v1.WorkflowsService/CreateCluster"
	WorkflowsService_GetCluster_FullMethodName    = "/workflows.v1.WorkflowsService/GetCluster"
	WorkflowsService_DeleteCluster_FullMethodName = "/workflows.v1.WorkflowsService/DeleteCluster"
	WorkflowsService_ListClusters_FullMethodName  = "/workflows.v1.WorkflowsService/ListClusters"
)

// WorkflowsServiceClient is the client API for WorkflowsService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// A service for managing workflows.
type WorkflowsServiceClient interface {
	CreateCluster(ctx context.Context, in *CreateClusterRequest, opts ...grpc.CallOption) (*Cluster, error)
	GetCluster(ctx context.Context, in *GetClusterRequest, opts ...grpc.CallOption) (*Cluster, error)
	DeleteCluster(ctx context.Context, in *DeleteClusterRequest, opts ...grpc.CallOption) (*DeleteClusterResponse, error)
	ListClusters(ctx context.Context, in *ListClustersRequest, opts ...grpc.CallOption) (*ListClustersResponse, error)
}

type workflowsServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewWorkflowsServiceClient(cc grpc.ClientConnInterface) WorkflowsServiceClient {
	return &workflowsServiceClient{cc}
}

func (c *workflowsServiceClient) CreateCluster(ctx context.Context, in *CreateClusterRequest, opts ...grpc.CallOption) (*Cluster, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Cluster)
	err := c.cc.Invoke(ctx, WorkflowsService_CreateCluster_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *workflowsServiceClient) GetCluster(ctx context.Context, in *GetClusterRequest, opts ...grpc.CallOption) (*Cluster, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Cluster)
	err := c.cc.Invoke(ctx, WorkflowsService_GetCluster_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *workflowsServiceClient) DeleteCluster(ctx context.Context, in *DeleteClusterRequest, opts ...grpc.CallOption) (*DeleteClusterResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteClusterResponse)
	err := c.cc.Invoke(ctx, WorkflowsService_DeleteCluster_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *workflowsServiceClient) ListClusters(ctx context.Context, in *ListClustersRequest, opts ...grpc.CallOption) (*ListClustersResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListClustersResponse)
	err := c.cc.Invoke(ctx, WorkflowsService_ListClusters_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// WorkflowsServiceServer is the server API for WorkflowsService service.
// All implementations must embed UnimplementedWorkflowsServiceServer
// for forward compatibility
//
// A service for managing workflows.
type WorkflowsServiceServer interface {
	CreateCluster(context.Context, *CreateClusterRequest) (*Cluster, error)
	GetCluster(context.Context, *GetClusterRequest) (*Cluster, error)
	DeleteCluster(context.Context, *DeleteClusterRequest) (*DeleteClusterResponse, error)
	ListClusters(context.Context, *ListClustersRequest) (*ListClustersResponse, error)
	mustEmbedUnimplementedWorkflowsServiceServer()
}

// UnimplementedWorkflowsServiceServer must be embedded to have forward compatible implementations.
type UnimplementedWorkflowsServiceServer struct {
}

func (UnimplementedWorkflowsServiceServer) CreateCluster(context.Context, *CreateClusterRequest) (*Cluster, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateCluster not implemented")
}
func (UnimplementedWorkflowsServiceServer) GetCluster(context.Context, *GetClusterRequest) (*Cluster, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCluster not implemented")
}
func (UnimplementedWorkflowsServiceServer) DeleteCluster(context.Context, *DeleteClusterRequest) (*DeleteClusterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteCluster not implemented")
}
func (UnimplementedWorkflowsServiceServer) ListClusters(context.Context, *ListClustersRequest) (*ListClustersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListClusters not implemented")
}
func (UnimplementedWorkflowsServiceServer) mustEmbedUnimplementedWorkflowsServiceServer() {}

// UnsafeWorkflowsServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to WorkflowsServiceServer will
// result in compilation errors.
type UnsafeWorkflowsServiceServer interface {
	mustEmbedUnimplementedWorkflowsServiceServer()
}

func RegisterWorkflowsServiceServer(s grpc.ServiceRegistrar, srv WorkflowsServiceServer) {
	s.RegisterService(&WorkflowsService_ServiceDesc, srv)
}

func _WorkflowsService_CreateCluster_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateClusterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WorkflowsServiceServer).CreateCluster(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WorkflowsService_CreateCluster_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WorkflowsServiceServer).CreateCluster(ctx, req.(*CreateClusterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WorkflowsService_GetCluster_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetClusterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WorkflowsServiceServer).GetCluster(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WorkflowsService_GetCluster_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WorkflowsServiceServer).GetCluster(ctx, req.(*GetClusterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WorkflowsService_DeleteCluster_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteClusterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WorkflowsServiceServer).DeleteCluster(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WorkflowsService_DeleteCluster_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WorkflowsServiceServer).DeleteCluster(ctx, req.(*DeleteClusterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WorkflowsService_ListClusters_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListClustersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WorkflowsServiceServer).ListClusters(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WorkflowsService_ListClusters_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WorkflowsServiceServer).ListClusters(ctx, req.(*ListClustersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// WorkflowsService_ServiceDesc is the grpc.ServiceDesc for WorkflowsService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var WorkflowsService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "workflows.v1.WorkflowsService",
	HandlerType: (*WorkflowsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateCluster",
			Handler:    _WorkflowsService_CreateCluster_Handler,
		},
		{
			MethodName: "GetCluster",
			Handler:    _WorkflowsService_GetCluster_Handler,
		},
		{
			MethodName: "DeleteCluster",
			Handler:    _WorkflowsService_DeleteCluster_Handler,
		},
		{
			MethodName: "ListClusters",
			Handler:    _WorkflowsService_ListClusters_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "workflows/v1/workflows.proto",
}
