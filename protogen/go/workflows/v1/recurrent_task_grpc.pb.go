// The external API for managing recurrent tasks in the Workflows service.

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: workflows/v1/recurrent_task.proto

package workflowsv1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	RecurrentTaskService_ListStorageLocations_FullMethodName  = "/workflows.v1.RecurrentTaskService/ListStorageLocations"
	RecurrentTaskService_GetStorageLocation_FullMethodName    = "/workflows.v1.RecurrentTaskService/GetStorageLocation"
	RecurrentTaskService_CreateStorageLocation_FullMethodName = "/workflows.v1.RecurrentTaskService/CreateStorageLocation"
	RecurrentTaskService_DeleteStorageLocation_FullMethodName = "/workflows.v1.RecurrentTaskService/DeleteStorageLocation"
	RecurrentTaskService_ListRecurrentTasks_FullMethodName    = "/workflows.v1.RecurrentTaskService/ListRecurrentTasks"
	RecurrentTaskService_GetRecurrentTask_FullMethodName      = "/workflows.v1.RecurrentTaskService/GetRecurrentTask"
	RecurrentTaskService_CreateRecurrentTask_FullMethodName   = "/workflows.v1.RecurrentTaskService/CreateRecurrentTask"
	RecurrentTaskService_UpdateRecurrentTask_FullMethodName   = "/workflows.v1.RecurrentTaskService/UpdateRecurrentTask"
	RecurrentTaskService_DeleteRecurrentTask_FullMethodName   = "/workflows.v1.RecurrentTaskService/DeleteRecurrentTask"
	RecurrentTaskService_DeleteAutomation_FullMethodName      = "/workflows.v1.RecurrentTaskService/DeleteAutomation"
)

// RecurrentTaskServiceClient is the client API for RecurrentTaskService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// RecurrentTaskService is a service for managing recurrent tasks. Currently, we support two types of triggers for
// recurrent tasks:
// - Bucket triggers, which triggers tasks when an object is uploaded to a storage bucket that matches a glob pattern
// - Cron triggers, which triggers tasks on a schedule
type RecurrentTaskServiceClient interface {
	// ListStorageLocations lists all the storage buckets that are available for use as bucket triggers.
	ListStorageLocations(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*StorageLocations, error)
	// GetStorageLocation gets a storage location by its ID.
	GetStorageLocation(ctx context.Context, in *UUID, opts ...grpc.CallOption) (*StorageLocation, error)
	// CreateStorageLocation creates a new storage bucket.
	CreateStorageLocation(ctx context.Context, in *StorageLocation, opts ...grpc.CallOption) (*StorageLocation, error)
	// DeleteStorageLocation deletes a storage location.
	DeleteStorageLocation(ctx context.Context, in *UUID, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// ListRecurrentTasks lists all the recurrent tasks that are currently registered in a namespace.
	ListRecurrentTasks(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*RecurrentTasks, error)
	// GetRecurrentTask gets a recurrent task by its ID.
	GetRecurrentTask(ctx context.Context, in *UUID, opts ...grpc.CallOption) (*RecurrentTaskPrototype, error)
	// CreateRecurrentTask creates a new recurrent task in a namespace.
	CreateRecurrentTask(ctx context.Context, in *RecurrentTaskPrototype, opts ...grpc.CallOption) (*RecurrentTaskPrototype, error)
	// UpdateRecurrentTask updates a recurrent task in a namespace.
	UpdateRecurrentTask(ctx context.Context, in *RecurrentTaskPrototype, opts ...grpc.CallOption) (*RecurrentTaskPrototype, error)
	// DeleteRecurrentTask deletes a recurrent task from a namespace.
	DeleteRecurrentTask(ctx context.Context, in *UUID, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// DeleteAutomation deletes an automation from a namespace.
	DeleteAutomation(ctx context.Context, in *DeleteAutomationRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type recurrentTaskServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRecurrentTaskServiceClient(cc grpc.ClientConnInterface) RecurrentTaskServiceClient {
	return &recurrentTaskServiceClient{cc}
}

func (c *recurrentTaskServiceClient) ListStorageLocations(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*StorageLocations, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(StorageLocations)
	err := c.cc.Invoke(ctx, RecurrentTaskService_ListStorageLocations_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *recurrentTaskServiceClient) GetStorageLocation(ctx context.Context, in *UUID, opts ...grpc.CallOption) (*StorageLocation, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(StorageLocation)
	err := c.cc.Invoke(ctx, RecurrentTaskService_GetStorageLocation_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *recurrentTaskServiceClient) CreateStorageLocation(ctx context.Context, in *StorageLocation, opts ...grpc.CallOption) (*StorageLocation, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(StorageLocation)
	err := c.cc.Invoke(ctx, RecurrentTaskService_CreateStorageLocation_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *recurrentTaskServiceClient) DeleteStorageLocation(ctx context.Context, in *UUID, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, RecurrentTaskService_DeleteStorageLocation_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *recurrentTaskServiceClient) ListRecurrentTasks(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*RecurrentTasks, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RecurrentTasks)
	err := c.cc.Invoke(ctx, RecurrentTaskService_ListRecurrentTasks_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *recurrentTaskServiceClient) GetRecurrentTask(ctx context.Context, in *UUID, opts ...grpc.CallOption) (*RecurrentTaskPrototype, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RecurrentTaskPrototype)
	err := c.cc.Invoke(ctx, RecurrentTaskService_GetRecurrentTask_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *recurrentTaskServiceClient) CreateRecurrentTask(ctx context.Context, in *RecurrentTaskPrototype, opts ...grpc.CallOption) (*RecurrentTaskPrototype, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RecurrentTaskPrototype)
	err := c.cc.Invoke(ctx, RecurrentTaskService_CreateRecurrentTask_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *recurrentTaskServiceClient) UpdateRecurrentTask(ctx context.Context, in *RecurrentTaskPrototype, opts ...grpc.CallOption) (*RecurrentTaskPrototype, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RecurrentTaskPrototype)
	err := c.cc.Invoke(ctx, RecurrentTaskService_UpdateRecurrentTask_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *recurrentTaskServiceClient) DeleteRecurrentTask(ctx context.Context, in *UUID, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, RecurrentTaskService_DeleteRecurrentTask_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *recurrentTaskServiceClient) DeleteAutomation(ctx context.Context, in *DeleteAutomationRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, RecurrentTaskService_DeleteAutomation_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RecurrentTaskServiceServer is the server API for RecurrentTaskService service.
// All implementations must embed UnimplementedRecurrentTaskServiceServer
// for forward compatibility.
//
// RecurrentTaskService is a service for managing recurrent tasks. Currently, we support two types of triggers for
// recurrent tasks:
// - Bucket triggers, which triggers tasks when an object is uploaded to a storage bucket that matches a glob pattern
// - Cron triggers, which triggers tasks on a schedule
type RecurrentTaskServiceServer interface {
	// ListStorageLocations lists all the storage buckets that are available for use as bucket triggers.
	ListStorageLocations(context.Context, *emptypb.Empty) (*StorageLocations, error)
	// GetStorageLocation gets a storage location by its ID.
	GetStorageLocation(context.Context, *UUID) (*StorageLocation, error)
	// CreateStorageLocation creates a new storage bucket.
	CreateStorageLocation(context.Context, *StorageLocation) (*StorageLocation, error)
	// DeleteStorageLocation deletes a storage location.
	DeleteStorageLocation(context.Context, *UUID) (*emptypb.Empty, error)
	// ListRecurrentTasks lists all the recurrent tasks that are currently registered in a namespace.
	ListRecurrentTasks(context.Context, *emptypb.Empty) (*RecurrentTasks, error)
	// GetRecurrentTask gets a recurrent task by its ID.
	GetRecurrentTask(context.Context, *UUID) (*RecurrentTaskPrototype, error)
	// CreateRecurrentTask creates a new recurrent task in a namespace.
	CreateRecurrentTask(context.Context, *RecurrentTaskPrototype) (*RecurrentTaskPrototype, error)
	// UpdateRecurrentTask updates a recurrent task in a namespace.
	UpdateRecurrentTask(context.Context, *RecurrentTaskPrototype) (*RecurrentTaskPrototype, error)
	// DeleteRecurrentTask deletes a recurrent task from a namespace.
	DeleteRecurrentTask(context.Context, *UUID) (*emptypb.Empty, error)
	// DeleteAutomation deletes an automation from a namespace.
	DeleteAutomation(context.Context, *DeleteAutomationRequest) (*emptypb.Empty, error)
	mustEmbedUnimplementedRecurrentTaskServiceServer()
}

// UnimplementedRecurrentTaskServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedRecurrentTaskServiceServer struct{}

func (UnimplementedRecurrentTaskServiceServer) ListStorageLocations(context.Context, *emptypb.Empty) (*StorageLocations, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListStorageLocations not implemented")
}
func (UnimplementedRecurrentTaskServiceServer) GetStorageLocation(context.Context, *UUID) (*StorageLocation, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStorageLocation not implemented")
}
func (UnimplementedRecurrentTaskServiceServer) CreateStorageLocation(context.Context, *StorageLocation) (*StorageLocation, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateStorageLocation not implemented")
}
func (UnimplementedRecurrentTaskServiceServer) DeleteStorageLocation(context.Context, *UUID) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteStorageLocation not implemented")
}
func (UnimplementedRecurrentTaskServiceServer) ListRecurrentTasks(context.Context, *emptypb.Empty) (*RecurrentTasks, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListRecurrentTasks not implemented")
}
func (UnimplementedRecurrentTaskServiceServer) GetRecurrentTask(context.Context, *UUID) (*RecurrentTaskPrototype, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRecurrentTask not implemented")
}
func (UnimplementedRecurrentTaskServiceServer) CreateRecurrentTask(context.Context, *RecurrentTaskPrototype) (*RecurrentTaskPrototype, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateRecurrentTask not implemented")
}
func (UnimplementedRecurrentTaskServiceServer) UpdateRecurrentTask(context.Context, *RecurrentTaskPrototype) (*RecurrentTaskPrototype, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateRecurrentTask not implemented")
}
func (UnimplementedRecurrentTaskServiceServer) DeleteRecurrentTask(context.Context, *UUID) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteRecurrentTask not implemented")
}
func (UnimplementedRecurrentTaskServiceServer) DeleteAutomation(context.Context, *DeleteAutomationRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteAutomation not implemented")
}
func (UnimplementedRecurrentTaskServiceServer) mustEmbedUnimplementedRecurrentTaskServiceServer() {}
func (UnimplementedRecurrentTaskServiceServer) testEmbeddedByValue()                              {}

// UnsafeRecurrentTaskServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RecurrentTaskServiceServer will
// result in compilation errors.
type UnsafeRecurrentTaskServiceServer interface {
	mustEmbedUnimplementedRecurrentTaskServiceServer()
}

func RegisterRecurrentTaskServiceServer(s grpc.ServiceRegistrar, srv RecurrentTaskServiceServer) {
	// If the following call pancis, it indicates UnimplementedRecurrentTaskServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&RecurrentTaskService_ServiceDesc, srv)
}

func _RecurrentTaskService_ListStorageLocations_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RecurrentTaskServiceServer).ListStorageLocations(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RecurrentTaskService_ListStorageLocations_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RecurrentTaskServiceServer).ListStorageLocations(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _RecurrentTaskService_GetStorageLocation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UUID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RecurrentTaskServiceServer).GetStorageLocation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RecurrentTaskService_GetStorageLocation_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RecurrentTaskServiceServer).GetStorageLocation(ctx, req.(*UUID))
	}
	return interceptor(ctx, in, info, handler)
}

func _RecurrentTaskService_CreateStorageLocation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StorageLocation)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RecurrentTaskServiceServer).CreateStorageLocation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RecurrentTaskService_CreateStorageLocation_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RecurrentTaskServiceServer).CreateStorageLocation(ctx, req.(*StorageLocation))
	}
	return interceptor(ctx, in, info, handler)
}

func _RecurrentTaskService_DeleteStorageLocation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UUID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RecurrentTaskServiceServer).DeleteStorageLocation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RecurrentTaskService_DeleteStorageLocation_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RecurrentTaskServiceServer).DeleteStorageLocation(ctx, req.(*UUID))
	}
	return interceptor(ctx, in, info, handler)
}

func _RecurrentTaskService_ListRecurrentTasks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RecurrentTaskServiceServer).ListRecurrentTasks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RecurrentTaskService_ListRecurrentTasks_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RecurrentTaskServiceServer).ListRecurrentTasks(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _RecurrentTaskService_GetRecurrentTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UUID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RecurrentTaskServiceServer).GetRecurrentTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RecurrentTaskService_GetRecurrentTask_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RecurrentTaskServiceServer).GetRecurrentTask(ctx, req.(*UUID))
	}
	return interceptor(ctx, in, info, handler)
}

func _RecurrentTaskService_CreateRecurrentTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RecurrentTaskPrototype)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RecurrentTaskServiceServer).CreateRecurrentTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RecurrentTaskService_CreateRecurrentTask_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RecurrentTaskServiceServer).CreateRecurrentTask(ctx, req.(*RecurrentTaskPrototype))
	}
	return interceptor(ctx, in, info, handler)
}

func _RecurrentTaskService_UpdateRecurrentTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RecurrentTaskPrototype)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RecurrentTaskServiceServer).UpdateRecurrentTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RecurrentTaskService_UpdateRecurrentTask_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RecurrentTaskServiceServer).UpdateRecurrentTask(ctx, req.(*RecurrentTaskPrototype))
	}
	return interceptor(ctx, in, info, handler)
}

func _RecurrentTaskService_DeleteRecurrentTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UUID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RecurrentTaskServiceServer).DeleteRecurrentTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RecurrentTaskService_DeleteRecurrentTask_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RecurrentTaskServiceServer).DeleteRecurrentTask(ctx, req.(*UUID))
	}
	return interceptor(ctx, in, info, handler)
}

func _RecurrentTaskService_DeleteAutomation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteAutomationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RecurrentTaskServiceServer).DeleteAutomation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RecurrentTaskService_DeleteAutomation_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RecurrentTaskServiceServer).DeleteAutomation(ctx, req.(*DeleteAutomationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RecurrentTaskService_ServiceDesc is the grpc.ServiceDesc for RecurrentTaskService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RecurrentTaskService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "workflows.v1.RecurrentTaskService",
	HandlerType: (*RecurrentTaskServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListStorageLocations",
			Handler:    _RecurrentTaskService_ListStorageLocations_Handler,
		},
		{
			MethodName: "GetStorageLocation",
			Handler:    _RecurrentTaskService_GetStorageLocation_Handler,
		},
		{
			MethodName: "CreateStorageLocation",
			Handler:    _RecurrentTaskService_CreateStorageLocation_Handler,
		},
		{
			MethodName: "DeleteStorageLocation",
			Handler:    _RecurrentTaskService_DeleteStorageLocation_Handler,
		},
		{
			MethodName: "ListRecurrentTasks",
			Handler:    _RecurrentTaskService_ListRecurrentTasks_Handler,
		},
		{
			MethodName: "GetRecurrentTask",
			Handler:    _RecurrentTaskService_GetRecurrentTask_Handler,
		},
		{
			MethodName: "CreateRecurrentTask",
			Handler:    _RecurrentTaskService_CreateRecurrentTask_Handler,
		},
		{
			MethodName: "UpdateRecurrentTask",
			Handler:    _RecurrentTaskService_UpdateRecurrentTask_Handler,
		},
		{
			MethodName: "DeleteRecurrentTask",
			Handler:    _RecurrentTaskService_DeleteRecurrentTask_Handler,
		},
		{
			MethodName: "DeleteAutomation",
			Handler:    _RecurrentTaskService_DeleteAutomation_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "workflows/v1/recurrent_task.proto",
}
