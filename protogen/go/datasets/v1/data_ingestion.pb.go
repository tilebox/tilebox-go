// Proto messages and service definition for the APIs related to loading and quering data.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        (unknown)
// source: datasets/v1/data_ingestion.proto

package datasetsv1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Legacy ingest request
// IngestDatapointsRequest is used to ingest one or multiple datapoints into a collection.
type IngestDatapointsRequest struct {
	state        protoimpl.MessageState `protogen:"open.v1"`
	CollectionId *ID                    `protobuf:"bytes,1,opt,name=collection_id,json=collectionId,proto3" json:"collection_id,omitempty"` // The collection to insert the datapoints into.
	Datapoints   *Datapoints            `protobuf:"bytes,2,opt,name=datapoints,proto3" json:"datapoints,omitempty"`
	// Whether to allow existing datapoints as part of the request. If true, datapoints that already exist will be
	// ignored, and the number of such existing datapoints will be returned in the response. If false, any datapoints
	// that already exist will result in an error. Setting this to true is useful for achieving idempotency (e.g.
	// allowing re-ingestion of datapoints that have already been ingested in the past).
	AllowExisting bool `protobuf:"varint,3,opt,name=allow_existing,json=allowExisting,proto3" json:"allow_existing,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *IngestDatapointsRequest) Reset() {
	*x = IngestDatapointsRequest{}
	mi := &file_datasets_v1_data_ingestion_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *IngestDatapointsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IngestDatapointsRequest) ProtoMessage() {}

func (x *IngestDatapointsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_datasets_v1_data_ingestion_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IngestDatapointsRequest.ProtoReflect.Descriptor instead.
func (*IngestDatapointsRequest) Descriptor() ([]byte, []int) {
	return file_datasets_v1_data_ingestion_proto_rawDescGZIP(), []int{0}
}

func (x *IngestDatapointsRequest) GetCollectionId() *ID {
	if x != nil {
		return x.CollectionId
	}
	return nil
}

func (x *IngestDatapointsRequest) GetDatapoints() *Datapoints {
	if x != nil {
		return x.Datapoints
	}
	return nil
}

func (x *IngestDatapointsRequest) GetAllowExisting() bool {
	if x != nil {
		return x.AllowExisting
	}
	return false
}

// IngestRequest is used to ingest one or multiple datapoints into a collection.
type IngestRequest struct {
	state        protoimpl.MessageState `protogen:"open.v1"`
	CollectionId *ID                    `protobuf:"bytes,1,opt,name=collection_id,json=collectionId,proto3" json:"collection_id,omitempty"` // The collection to insert the datapoints into.
	// The datapoints to insert. The values here are encoded protobuf messages. The type of the message is determined
	// by the type of the dataset that the specified collection belongs to.
	Values [][]byte `protobuf:"bytes,2,rep,name=values,proto3" json:"values,omitempty"`
	// Whether to allow existing datapoints as part of the request. If true, datapoints that already exist will be
	// ignored, and the number of such existing datapoints will be returned in the response. If false, any datapoints
	// that already exist will result in an error. Setting this to true is useful for achieving idempotency (e.g.
	// allowing re-ingestion of datapoints that have already been ingested in the past).
	AllowExisting bool `protobuf:"varint,3,opt,name=allow_existing,json=allowExisting,proto3" json:"allow_existing,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *IngestRequest) Reset() {
	*x = IngestRequest{}
	mi := &file_datasets_v1_data_ingestion_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *IngestRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IngestRequest) ProtoMessage() {}

func (x *IngestRequest) ProtoReflect() protoreflect.Message {
	mi := &file_datasets_v1_data_ingestion_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IngestRequest.ProtoReflect.Descriptor instead.
func (*IngestRequest) Descriptor() ([]byte, []int) {
	return file_datasets_v1_data_ingestion_proto_rawDescGZIP(), []int{1}
}

func (x *IngestRequest) GetCollectionId() *ID {
	if x != nil {
		return x.CollectionId
	}
	return nil
}

func (x *IngestRequest) GetValues() [][]byte {
	if x != nil {
		return x.Values
	}
	return nil
}

func (x *IngestRequest) GetAllowExisting() bool {
	if x != nil {
		return x.AllowExisting
	}
	return false
}

// IngestResponse is the response to a IngestRequest, indicating the number of datapoints that were
// ingested as well as the generated ids for those datapoints.
type IngestResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	NumCreated    int64                  `protobuf:"varint,1,opt,name=num_created,json=numCreated,proto3" json:"num_created,omitempty"`      // The number of datapoints that were created.
	NumExisting   int64                  `protobuf:"varint,2,opt,name=num_existing,json=numExisting,proto3" json:"num_existing,omitempty"`   // The number of datapoints that were ignored because they already existed.
	DatapointIds  []*ID                  `protobuf:"bytes,3,rep,name=datapoint_ids,json=datapointIds,proto3" json:"datapoint_ids,omitempty"` // The ids of the datapoints in the same order as the datapoints in the request.
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *IngestResponse) Reset() {
	*x = IngestResponse{}
	mi := &file_datasets_v1_data_ingestion_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *IngestResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IngestResponse) ProtoMessage() {}

func (x *IngestResponse) ProtoReflect() protoreflect.Message {
	mi := &file_datasets_v1_data_ingestion_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IngestResponse.ProtoReflect.Descriptor instead.
func (*IngestResponse) Descriptor() ([]byte, []int) {
	return file_datasets_v1_data_ingestion_proto_rawDescGZIP(), []int{2}
}

func (x *IngestResponse) GetNumCreated() int64 {
	if x != nil {
		return x.NumCreated
	}
	return 0
}

func (x *IngestResponse) GetNumExisting() int64 {
	if x != nil {
		return x.NumExisting
	}
	return 0
}

func (x *IngestResponse) GetDatapointIds() []*ID {
	if x != nil {
		return x.DatapointIds
	}
	return nil
}

// DeleteRequest is used to delete multiple datapoints from a collection.
type DeleteRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	CollectionId  *ID                    `protobuf:"bytes,1,opt,name=collection_id,json=collectionId,proto3" json:"collection_id,omitempty"`
	DatapointIds  []*ID                  `protobuf:"bytes,2,rep,name=datapoint_ids,json=datapointIds,proto3" json:"datapoint_ids,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteRequest) Reset() {
	*x = DeleteRequest{}
	mi := &file_datasets_v1_data_ingestion_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteRequest) ProtoMessage() {}

func (x *DeleteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_datasets_v1_data_ingestion_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteRequest.ProtoReflect.Descriptor instead.
func (*DeleteRequest) Descriptor() ([]byte, []int) {
	return file_datasets_v1_data_ingestion_proto_rawDescGZIP(), []int{3}
}

func (x *DeleteRequest) GetCollectionId() *ID {
	if x != nil {
		return x.CollectionId
	}
	return nil
}

func (x *DeleteRequest) GetDatapointIds() []*ID {
	if x != nil {
		return x.DatapointIds
	}
	return nil
}

// DeleteResponse is used to indicate that multiple datapoints were deleted.
type DeleteResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	NumDeleted    int64                  `protobuf:"varint,1,opt,name=num_deleted,json=numDeleted,proto3" json:"num_deleted,omitempty"` // The number of datapoints that were deleted.
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteResponse) Reset() {
	*x = DeleteResponse{}
	mi := &file_datasets_v1_data_ingestion_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteResponse) ProtoMessage() {}

func (x *DeleteResponse) ProtoReflect() protoreflect.Message {
	mi := &file_datasets_v1_data_ingestion_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteResponse.ProtoReflect.Descriptor instead.
func (*DeleteResponse) Descriptor() ([]byte, []int) {
	return file_datasets_v1_data_ingestion_proto_rawDescGZIP(), []int{4}
}

func (x *DeleteResponse) GetNumDeleted() int64 {
	if x != nil {
		return x.NumDeleted
	}
	return 0
}

var File_datasets_v1_data_ingestion_proto protoreflect.FileDescriptor

var file_datasets_v1_data_ingestion_proto_rawDesc = string([]byte{
	0x0a, 0x20, 0x64, 0x61, 0x74, 0x61, 0x73, 0x65, 0x74, 0x73, 0x2f, 0x76, 0x31, 0x2f, 0x64, 0x61,
	0x74, 0x61, 0x5f, 0x69, 0x6e, 0x67, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x0b, 0x64, 0x61, 0x74, 0x61, 0x73, 0x65, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x1a,
	0x16, 0x64, 0x61, 0x74, 0x61, 0x73, 0x65, 0x74, 0x73, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x6f, 0x72,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xaf, 0x01, 0x0a, 0x17, 0x49, 0x6e, 0x67, 0x65,
	0x73, 0x74, 0x44, 0x61, 0x74, 0x61, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x34, 0x0a, 0x0d, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x64, 0x61, 0x74,
	0x61, 0x73, 0x65, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x49, 0x44, 0x52, 0x0c, 0x63, 0x6f, 0x6c,
	0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x37, 0x0a, 0x0a, 0x64, 0x61, 0x74,
	0x61, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e,
	0x64, 0x61, 0x74, 0x61, 0x73, 0x65, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x61, 0x74, 0x61,
	0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x52, 0x0a, 0x64, 0x61, 0x74, 0x61, 0x70, 0x6f, 0x69, 0x6e,
	0x74, 0x73, 0x12, 0x25, 0x0a, 0x0e, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x5f, 0x65, 0x78, 0x69, 0x73,
	0x74, 0x69, 0x6e, 0x67, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0d, 0x61, 0x6c, 0x6c, 0x6f,
	0x77, 0x45, 0x78, 0x69, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x22, 0x84, 0x01, 0x0a, 0x0d, 0x49, 0x6e,
	0x67, 0x65, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x34, 0x0a, 0x0d, 0x63,
	0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x73, 0x65, 0x74, 0x73, 0x2e, 0x76, 0x31,
	0x2e, 0x49, 0x44, 0x52, 0x0c, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x49,
	0x64, 0x12, 0x16, 0x0a, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28,
	0x0c, 0x52, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x12, 0x25, 0x0a, 0x0e, 0x61, 0x6c, 0x6c,
	0x6f, 0x77, 0x5f, 0x65, 0x78, 0x69, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x0d, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x45, 0x78, 0x69, 0x73, 0x74, 0x69, 0x6e, 0x67,
	0x22, 0x8a, 0x01, 0x0a, 0x0e, 0x49, 0x6e, 0x67, 0x65, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x6e, 0x75, 0x6d, 0x5f, 0x63, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x6e, 0x75, 0x6d, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x64, 0x12, 0x21, 0x0a, 0x0c, 0x6e, 0x75, 0x6d, 0x5f, 0x65, 0x78, 0x69, 0x73,
	0x74, 0x69, 0x6e, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x6e, 0x75, 0x6d, 0x45,
	0x78, 0x69, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x12, 0x34, 0x0a, 0x0d, 0x64, 0x61, 0x74, 0x61, 0x70,
	0x6f, 0x69, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0f,
	0x2e, 0x64, 0x61, 0x74, 0x61, 0x73, 0x65, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x49, 0x44, 0x52,
	0x0c, 0x64, 0x61, 0x74, 0x61, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x49, 0x64, 0x73, 0x22, 0x7b, 0x0a,
	0x0d, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x34,
	0x0a, 0x0d, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x73, 0x65, 0x74, 0x73,
	0x2e, 0x76, 0x31, 0x2e, 0x49, 0x44, 0x52, 0x0c, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x49, 0x64, 0x12, 0x34, 0x0a, 0x0d, 0x64, 0x61, 0x74, 0x61, 0x70, 0x6f, 0x69, 0x6e,
	0x74, 0x5f, 0x69, 0x64, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x64, 0x61,
	0x74, 0x61, 0x73, 0x65, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x49, 0x44, 0x52, 0x0c, 0x64, 0x61,
	0x74, 0x61, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x49, 0x64, 0x73, 0x22, 0x31, 0x0a, 0x0e, 0x44, 0x65,
	0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1f, 0x0a, 0x0b,
	0x6e, 0x75, 0x6d, 0x5f, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x0a, 0x6e, 0x75, 0x6d, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x32, 0xf9, 0x01,
	0x0a, 0x14, 0x44, 0x61, 0x74, 0x61, 0x49, 0x6e, 0x67, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x57, 0x0a, 0x10, 0x49, 0x6e, 0x67, 0x65, 0x73, 0x74,
	0x44, 0x61, 0x74, 0x61, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x12, 0x24, 0x2e, 0x64, 0x61, 0x74,
	0x61, 0x73, 0x65, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x49, 0x6e, 0x67, 0x65, 0x73, 0x74, 0x44,
	0x61, 0x74, 0x61, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x1b, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x73, 0x65, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x49,
	0x6e, 0x67, 0x65, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12,
	0x43, 0x0a, 0x06, 0x49, 0x6e, 0x67, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x2e, 0x64, 0x61, 0x74, 0x61,
	0x73, 0x65, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x49, 0x6e, 0x67, 0x65, 0x73, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x73, 0x65, 0x74, 0x73,
	0x2e, 0x76, 0x31, 0x2e, 0x49, 0x6e, 0x67, 0x65, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x00, 0x12, 0x43, 0x0a, 0x06, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x12, 0x1a,
	0x2e, 0x64, 0x61, 0x74, 0x61, 0x73, 0x65, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x65, 0x6c,
	0x65, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x64, 0x61, 0x74,
	0x61, 0x73, 0x65, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0xb4, 0x01, 0x0a, 0x0f, 0x63, 0x6f,
	0x6d, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x73, 0x65, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x42, 0x12, 0x44,
	0x61, 0x74, 0x61, 0x49, 0x6e, 0x67, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x50, 0x72, 0x6f, 0x74,
	0x6f, 0x50, 0x01, 0x5a, 0x40, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x74, 0x69, 0x6c, 0x65, 0x62, 0x6f, 0x78, 0x2f, 0x74, 0x69, 0x6c, 0x65, 0x62, 0x6f, 0x78, 0x2d,
	0x67, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x67, 0x65, 0x6e, 0x2f, 0x67, 0x6f, 0x2f, 0x64,
	0x61, 0x74, 0x61, 0x73, 0x65, 0x74, 0x73, 0x2f, 0x76, 0x31, 0x3b, 0x64, 0x61, 0x74, 0x61, 0x73,
	0x65, 0x74, 0x73, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x44, 0x58, 0x58, 0xaa, 0x02, 0x0b, 0x44, 0x61,
	0x74, 0x61, 0x73, 0x65, 0x74, 0x73, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x0b, 0x44, 0x61, 0x74, 0x61,
	0x73, 0x65, 0x74, 0x73, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x17, 0x44, 0x61, 0x74, 0x61, 0x73, 0x65,
	0x74, 0x73, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74,
	0x61, 0xea, 0x02, 0x0c, 0x44, 0x61, 0x74, 0x61, 0x73, 0x65, 0x74, 0x73, 0x3a, 0x3a, 0x56, 0x31,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_datasets_v1_data_ingestion_proto_rawDescOnce sync.Once
	file_datasets_v1_data_ingestion_proto_rawDescData []byte
)

func file_datasets_v1_data_ingestion_proto_rawDescGZIP() []byte {
	file_datasets_v1_data_ingestion_proto_rawDescOnce.Do(func() {
		file_datasets_v1_data_ingestion_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_datasets_v1_data_ingestion_proto_rawDesc), len(file_datasets_v1_data_ingestion_proto_rawDesc)))
	})
	return file_datasets_v1_data_ingestion_proto_rawDescData
}

var file_datasets_v1_data_ingestion_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_datasets_v1_data_ingestion_proto_goTypes = []any{
	(*IngestDatapointsRequest)(nil), // 0: datasets.v1.IngestDatapointsRequest
	(*IngestRequest)(nil),           // 1: datasets.v1.IngestRequest
	(*IngestResponse)(nil),          // 2: datasets.v1.IngestResponse
	(*DeleteRequest)(nil),           // 3: datasets.v1.DeleteRequest
	(*DeleteResponse)(nil),          // 4: datasets.v1.DeleteResponse
	(*ID)(nil),                      // 5: datasets.v1.ID
	(*Datapoints)(nil),              // 6: datasets.v1.Datapoints
}
var file_datasets_v1_data_ingestion_proto_depIdxs = []int32{
	5, // 0: datasets.v1.IngestDatapointsRequest.collection_id:type_name -> datasets.v1.ID
	6, // 1: datasets.v1.IngestDatapointsRequest.datapoints:type_name -> datasets.v1.Datapoints
	5, // 2: datasets.v1.IngestRequest.collection_id:type_name -> datasets.v1.ID
	5, // 3: datasets.v1.IngestResponse.datapoint_ids:type_name -> datasets.v1.ID
	5, // 4: datasets.v1.DeleteRequest.collection_id:type_name -> datasets.v1.ID
	5, // 5: datasets.v1.DeleteRequest.datapoint_ids:type_name -> datasets.v1.ID
	0, // 6: datasets.v1.DataIngestionService.IngestDatapoints:input_type -> datasets.v1.IngestDatapointsRequest
	1, // 7: datasets.v1.DataIngestionService.Ingest:input_type -> datasets.v1.IngestRequest
	3, // 8: datasets.v1.DataIngestionService.Delete:input_type -> datasets.v1.DeleteRequest
	2, // 9: datasets.v1.DataIngestionService.IngestDatapoints:output_type -> datasets.v1.IngestResponse
	2, // 10: datasets.v1.DataIngestionService.Ingest:output_type -> datasets.v1.IngestResponse
	4, // 11: datasets.v1.DataIngestionService.Delete:output_type -> datasets.v1.DeleteResponse
	9, // [9:12] is the sub-list for method output_type
	6, // [6:9] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_datasets_v1_data_ingestion_proto_init() }
func file_datasets_v1_data_ingestion_proto_init() {
	if File_datasets_v1_data_ingestion_proto != nil {
		return
	}
	file_datasets_v1_core_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_datasets_v1_data_ingestion_proto_rawDesc), len(file_datasets_v1_data_ingestion_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_datasets_v1_data_ingestion_proto_goTypes,
		DependencyIndexes: file_datasets_v1_data_ingestion_proto_depIdxs,
		MessageInfos:      file_datasets_v1_data_ingestion_proto_msgTypes,
	}.Build()
	File_datasets_v1_data_ingestion_proto = out.File
	file_datasets_v1_data_ingestion_proto_goTypes = nil
	file_datasets_v1_data_ingestion_proto_depIdxs = nil
}
