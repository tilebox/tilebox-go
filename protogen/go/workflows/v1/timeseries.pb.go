// Timeseries types for workflows

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.0
// 	protoc        (unknown)
// source: workflows/v1/timeseries.proto

package workflowsv1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// TimeInterval is a message that represents a time interval, with a start and end time.
// copied from datasets/v1/datasets.proto until we figure out how to import across packages
type TimeInterval struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StartTime *timestamppb.Timestamp `protobuf:"bytes,1,opt,name=start_time,json=startTime,proto3" json:"start_time,omitempty"` // Start time of the interval.
	EndTime   *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=end_time,json=endTime,proto3" json:"end_time,omitempty"`       // End time of the interval.
	// Flag indicating whether the start time is exclusive. If true, the start time is not included in the interval.
	StartExclusive bool `protobuf:"varint,3,opt,name=start_exclusive,json=startExclusive,proto3" json:"start_exclusive,omitempty"`
	// Flag indicating whether the end time is inclusive. If true, the end time is included in the interval.
	EndInclusive bool `protobuf:"varint,4,opt,name=end_inclusive,json=endInclusive,proto3" json:"end_inclusive,omitempty"`
}

func (x *TimeInterval) Reset() {
	*x = TimeInterval{}
	if protoimpl.UnsafeEnabled {
		mi := &file_workflows_v1_timeseries_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TimeInterval) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TimeInterval) ProtoMessage() {}

func (x *TimeInterval) ProtoReflect() protoreflect.Message {
	mi := &file_workflows_v1_timeseries_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TimeInterval.ProtoReflect.Descriptor instead.
func (*TimeInterval) Descriptor() ([]byte, []int) {
	return file_workflows_v1_timeseries_proto_rawDescGZIP(), []int{0}
}

func (x *TimeInterval) GetStartTime() *timestamppb.Timestamp {
	if x != nil {
		return x.StartTime
	}
	return nil
}

func (x *TimeInterval) GetEndTime() *timestamppb.Timestamp {
	if x != nil {
		return x.EndTime
	}
	return nil
}

func (x *TimeInterval) GetStartExclusive() bool {
	if x != nil {
		return x.StartExclusive
	}
	return false
}

func (x *TimeInterval) GetEndInclusive() bool {
	if x != nil {
		return x.EndInclusive
	}
	return false
}

// DatapointInterval is a message that represents a datapoint interval, with a start and end datapoint.
type DatapointInterval struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Start *UUID `protobuf:"bytes,1,opt,name=start,proto3" json:"start,omitempty"`
	End   *UUID `protobuf:"bytes,2,opt,name=end,proto3" json:"end,omitempty"`
}

func (x *DatapointInterval) Reset() {
	*x = DatapointInterval{}
	if protoimpl.UnsafeEnabled {
		mi := &file_workflows_v1_timeseries_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DatapointInterval) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DatapointInterval) ProtoMessage() {}

func (x *DatapointInterval) ProtoReflect() protoreflect.Message {
	mi := &file_workflows_v1_timeseries_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DatapointInterval.ProtoReflect.Descriptor instead.
func (*DatapointInterval) Descriptor() ([]byte, []int) {
	return file_workflows_v1_timeseries_proto_rawDescGZIP(), []int{1}
}

func (x *DatapointInterval) GetStart() *UUID {
	if x != nil {
		return x.Start
	}
	return nil
}

func (x *DatapointInterval) GetEnd() *UUID {
	if x != nil {
		return x.End
	}
	return nil
}

// TimeseriesDatasetChunk is a message that represents a chunk of a timeseries dataset.
type TimeseriesDatasetChunk struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DatasetId             *UUID              `protobuf:"bytes,1,opt,name=dataset_id,json=datasetId,proto3" json:"dataset_id,omitempty"`
	CollectionId          *UUID              `protobuf:"bytes,2,opt,name=collection_id,json=collectionId,proto3" json:"collection_id,omitempty"`
	TimeInterval          *TimeInterval      `protobuf:"bytes,3,opt,name=time_interval,json=timeInterval,proto3" json:"time_interval,omitempty"`
	DatapointInterval     *DatapointInterval `protobuf:"bytes,4,opt,name=datapoint_interval,json=datapointInterval,proto3" json:"datapoint_interval,omitempty"`
	BranchFactor          int32              `protobuf:"varint,5,opt,name=branch_factor,json=branchFactor,proto3" json:"branch_factor,omitempty"`
	ChunkSize             int32              `protobuf:"varint,6,opt,name=chunk_size,json=chunkSize,proto3" json:"chunk_size,omitempty"`
	DatapointsPer_365Days int64              `protobuf:"varint,7,opt,name=datapoints_per_365_days,json=datapointsPer365Days,proto3" json:"datapoints_per_365_days,omitempty"`
}

func (x *TimeseriesDatasetChunk) Reset() {
	*x = TimeseriesDatasetChunk{}
	if protoimpl.UnsafeEnabled {
		mi := &file_workflows_v1_timeseries_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TimeseriesDatasetChunk) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TimeseriesDatasetChunk) ProtoMessage() {}

func (x *TimeseriesDatasetChunk) ProtoReflect() protoreflect.Message {
	mi := &file_workflows_v1_timeseries_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TimeseriesDatasetChunk.ProtoReflect.Descriptor instead.
func (*TimeseriesDatasetChunk) Descriptor() ([]byte, []int) {
	return file_workflows_v1_timeseries_proto_rawDescGZIP(), []int{2}
}

func (x *TimeseriesDatasetChunk) GetDatasetId() *UUID {
	if x != nil {
		return x.DatasetId
	}
	return nil
}

func (x *TimeseriesDatasetChunk) GetCollectionId() *UUID {
	if x != nil {
		return x.CollectionId
	}
	return nil
}

func (x *TimeseriesDatasetChunk) GetTimeInterval() *TimeInterval {
	if x != nil {
		return x.TimeInterval
	}
	return nil
}

func (x *TimeseriesDatasetChunk) GetDatapointInterval() *DatapointInterval {
	if x != nil {
		return x.DatapointInterval
	}
	return nil
}

func (x *TimeseriesDatasetChunk) GetBranchFactor() int32 {
	if x != nil {
		return x.BranchFactor
	}
	return 0
}

func (x *TimeseriesDatasetChunk) GetChunkSize() int32 {
	if x != nil {
		return x.ChunkSize
	}
	return 0
}

func (x *TimeseriesDatasetChunk) GetDatapointsPer_365Days() int64 {
	if x != nil {
		return x.DatapointsPer_365Days
	}
	return 0
}

// TimeChunk is a message that represents a time interval and a chunk size.
type TimeChunk struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TimeInterval *TimeInterval        `protobuf:"bytes,1,opt,name=time_interval,json=timeInterval,proto3" json:"time_interval,omitempty"`
	ChunkSize    *durationpb.Duration `protobuf:"bytes,2,opt,name=chunk_size,json=chunkSize,proto3" json:"chunk_size,omitempty"`
}

func (x *TimeChunk) Reset() {
	*x = TimeChunk{}
	if protoimpl.UnsafeEnabled {
		mi := &file_workflows_v1_timeseries_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TimeChunk) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TimeChunk) ProtoMessage() {}

func (x *TimeChunk) ProtoReflect() protoreflect.Message {
	mi := &file_workflows_v1_timeseries_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TimeChunk.ProtoReflect.Descriptor instead.
func (*TimeChunk) Descriptor() ([]byte, []int) {
	return file_workflows_v1_timeseries_proto_rawDescGZIP(), []int{3}
}

func (x *TimeChunk) GetTimeInterval() *TimeInterval {
	if x != nil {
		return x.TimeInterval
	}
	return nil
}

func (x *TimeChunk) GetChunkSize() *durationpb.Duration {
	if x != nil {
		return x.ChunkSize
	}
	return nil
}

var File_workflows_v1_timeseries_proto protoreflect.FileDescriptor

var file_workflows_v1_timeseries_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x2f, 0x76, 0x31, 0x2f, 0x74,
	0x69, 0x6d, 0x65, 0x73, 0x65, 0x72, 0x69, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x0c, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x2e, 0x76, 0x31, 0x1a, 0x17, 0x77,
	0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x6f, 0x72, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xce, 0x01, 0x0a, 0x0c, 0x54, 0x69, 0x6d, 0x65,
	0x49, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x12, 0x39, 0x0a, 0x0a, 0x73, 0x74, 0x61, 0x72,
	0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x73, 0x74, 0x61, 0x72, 0x74, 0x54,
	0x69, 0x6d, 0x65, 0x12, 0x35, 0x0a, 0x08, 0x65, 0x6e, 0x64, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x52, 0x07, 0x65, 0x6e, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x27, 0x0a, 0x0f, 0x73, 0x74,
	0x61, 0x72, 0x74, 0x5f, 0x65, 0x78, 0x63, 0x6c, 0x75, 0x73, 0x69, 0x76, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x0e, 0x73, 0x74, 0x61, 0x72, 0x74, 0x45, 0x78, 0x63, 0x6c, 0x75, 0x73,
	0x69, 0x76, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x65, 0x6e, 0x64, 0x5f, 0x69, 0x6e, 0x63, 0x6c, 0x75,
	0x73, 0x69, 0x76, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0c, 0x65, 0x6e, 0x64, 0x49,
	0x6e, 0x63, 0x6c, 0x75, 0x73, 0x69, 0x76, 0x65, 0x22, 0x63, 0x0a, 0x11, 0x44, 0x61, 0x74, 0x61,
	0x70, 0x6f, 0x69, 0x6e, 0x74, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x12, 0x28, 0x0a,
	0x05, 0x73, 0x74, 0x61, 0x72, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x77,
	0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x55, 0x49, 0x44,
	0x52, 0x05, 0x73, 0x74, 0x61, 0x72, 0x74, 0x12, 0x24, 0x0a, 0x03, 0x65, 0x6e, 0x64, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73,
	0x2e, 0x76, 0x31, 0x2e, 0x55, 0x55, 0x49, 0x44, 0x52, 0x03, 0x65, 0x6e, 0x64, 0x22, 0x90, 0x03,
	0x0a, 0x16, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x65, 0x72, 0x69, 0x65, 0x73, 0x44, 0x61, 0x74, 0x61,
	0x73, 0x65, 0x74, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x12, 0x31, 0x0a, 0x0a, 0x64, 0x61, 0x74, 0x61,
	0x73, 0x65, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x77,
	0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x55, 0x49, 0x44,
	0x52, 0x09, 0x64, 0x61, 0x74, 0x61, 0x73, 0x65, 0x74, 0x49, 0x64, 0x12, 0x37, 0x0a, 0x0d, 0x63,
	0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x12, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x2e, 0x76,
	0x31, 0x2e, 0x55, 0x55, 0x49, 0x44, 0x52, 0x0c, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x49, 0x64, 0x12, 0x3f, 0x0a, 0x0d, 0x74, 0x69, 0x6d, 0x65, 0x5f, 0x69, 0x6e, 0x74,
	0x65, 0x72, 0x76, 0x61, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x77, 0x6f,
	0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x49,
	0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x52, 0x0c, 0x74, 0x69, 0x6d, 0x65, 0x49, 0x6e, 0x74,
	0x65, 0x72, 0x76, 0x61, 0x6c, 0x12, 0x4e, 0x0a, 0x12, 0x64, 0x61, 0x74, 0x61, 0x70, 0x6f, 0x69,
	0x6e, 0x74, 0x5f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1f, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x2e, 0x76, 0x31,
	0x2e, 0x44, 0x61, 0x74, 0x61, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x76,
	0x61, 0x6c, 0x52, 0x11, 0x64, 0x61, 0x74, 0x61, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x49, 0x6e, 0x74,
	0x65, 0x72, 0x76, 0x61, 0x6c, 0x12, 0x23, 0x0a, 0x0d, 0x62, 0x72, 0x61, 0x6e, 0x63, 0x68, 0x5f,
	0x66, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0c, 0x62, 0x72,
	0x61, 0x6e, 0x63, 0x68, 0x46, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x68,
	0x75, 0x6e, 0x6b, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09,
	0x63, 0x68, 0x75, 0x6e, 0x6b, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x35, 0x0a, 0x17, 0x64, 0x61, 0x74,
	0x61, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x5f, 0x70, 0x65, 0x72, 0x5f, 0x33, 0x36, 0x35, 0x5f,
	0x64, 0x61, 0x79, 0x73, 0x18, 0x07, 0x20, 0x01, 0x28, 0x03, 0x52, 0x14, 0x64, 0x61, 0x74, 0x61,
	0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x50, 0x65, 0x72, 0x33, 0x36, 0x35, 0x44, 0x61, 0x79, 0x73,
	0x22, 0x86, 0x01, 0x0a, 0x09, 0x54, 0x69, 0x6d, 0x65, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x12, 0x3f,
	0x0a, 0x0d, 0x74, 0x69, 0x6d, 0x65, 0x5f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77,
	0x73, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61,
	0x6c, 0x52, 0x0c, 0x74, 0x69, 0x6d, 0x65, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x12,
	0x38, 0x0a, 0x0a, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x09,
	0x63, 0x68, 0x75, 0x6e, 0x6b, 0x53, 0x69, 0x7a, 0x65, 0x42, 0xb8, 0x01, 0x0a, 0x10, 0x63, 0x6f,
	0x6d, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x2e, 0x76, 0x31, 0x42, 0x0f,
	0x54, 0x69, 0x6d, 0x65, 0x73, 0x65, 0x72, 0x69, 0x65, 0x73, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50,
	0x01, 0x5a, 0x42, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x74, 0x69,
	0x6c, 0x65, 0x62, 0x6f, 0x78, 0x2f, 0x74, 0x69, 0x6c, 0x65, 0x62, 0x6f, 0x78, 0x2d, 0x67, 0x6f,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x67, 0x65, 0x6e, 0x2f, 0x67, 0x6f, 0x2f, 0x77, 0x6f, 0x72,
	0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x2f, 0x76, 0x31, 0x3b, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c,
	0x6f, 0x77, 0x73, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x57, 0x58, 0x58, 0xaa, 0x02, 0x0c, 0x57, 0x6f,
	0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x0c, 0x57, 0x6f, 0x72,
	0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x18, 0x57, 0x6f, 0x72, 0x6b,
	0x66, 0x6c, 0x6f, 0x77, 0x73, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61,
	0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x0d, 0x57, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73,
	0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_workflows_v1_timeseries_proto_rawDescOnce sync.Once
	file_workflows_v1_timeseries_proto_rawDescData = file_workflows_v1_timeseries_proto_rawDesc
)

func file_workflows_v1_timeseries_proto_rawDescGZIP() []byte {
	file_workflows_v1_timeseries_proto_rawDescOnce.Do(func() {
		file_workflows_v1_timeseries_proto_rawDescData = protoimpl.X.CompressGZIP(file_workflows_v1_timeseries_proto_rawDescData)
	})
	return file_workflows_v1_timeseries_proto_rawDescData
}

var file_workflows_v1_timeseries_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_workflows_v1_timeseries_proto_goTypes = []interface{}{
	(*TimeInterval)(nil),           // 0: workflows.v1.TimeInterval
	(*DatapointInterval)(nil),      // 1: workflows.v1.DatapointInterval
	(*TimeseriesDatasetChunk)(nil), // 2: workflows.v1.TimeseriesDatasetChunk
	(*TimeChunk)(nil),              // 3: workflows.v1.TimeChunk
	(*timestamppb.Timestamp)(nil),  // 4: google.protobuf.Timestamp
	(*UUID)(nil),                   // 5: workflows.v1.UUID
	(*durationpb.Duration)(nil),    // 6: google.protobuf.Duration
}
var file_workflows_v1_timeseries_proto_depIdxs = []int32{
	4,  // 0: workflows.v1.TimeInterval.start_time:type_name -> google.protobuf.Timestamp
	4,  // 1: workflows.v1.TimeInterval.end_time:type_name -> google.protobuf.Timestamp
	5,  // 2: workflows.v1.DatapointInterval.start:type_name -> workflows.v1.UUID
	5,  // 3: workflows.v1.DatapointInterval.end:type_name -> workflows.v1.UUID
	5,  // 4: workflows.v1.TimeseriesDatasetChunk.dataset_id:type_name -> workflows.v1.UUID
	5,  // 5: workflows.v1.TimeseriesDatasetChunk.collection_id:type_name -> workflows.v1.UUID
	0,  // 6: workflows.v1.TimeseriesDatasetChunk.time_interval:type_name -> workflows.v1.TimeInterval
	1,  // 7: workflows.v1.TimeseriesDatasetChunk.datapoint_interval:type_name -> workflows.v1.DatapointInterval
	0,  // 8: workflows.v1.TimeChunk.time_interval:type_name -> workflows.v1.TimeInterval
	6,  // 9: workflows.v1.TimeChunk.chunk_size:type_name -> google.protobuf.Duration
	10, // [10:10] is the sub-list for method output_type
	10, // [10:10] is the sub-list for method input_type
	10, // [10:10] is the sub-list for extension type_name
	10, // [10:10] is the sub-list for extension extendee
	0,  // [0:10] is the sub-list for field type_name
}

func init() { file_workflows_v1_timeseries_proto_init() }
func file_workflows_v1_timeseries_proto_init() {
	if File_workflows_v1_timeseries_proto != nil {
		return
	}
	file_workflows_v1_core_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_workflows_v1_timeseries_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TimeInterval); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_workflows_v1_timeseries_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DatapointInterval); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_workflows_v1_timeseries_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TimeseriesDatasetChunk); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_workflows_v1_timeseries_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TimeChunk); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_workflows_v1_timeseries_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_workflows_v1_timeseries_proto_goTypes,
		DependencyIndexes: file_workflows_v1_timeseries_proto_depIdxs,
		MessageInfos:      file_workflows_v1_timeseries_proto_msgTypes,
	}.Build()
	File_workflows_v1_timeseries_proto = out.File
	file_workflows_v1_timeseries_proto_rawDesc = nil
	file_workflows_v1_timeseries_proto_goTypes = nil
	file_workflows_v1_timeseries_proto_depIdxs = nil
}
