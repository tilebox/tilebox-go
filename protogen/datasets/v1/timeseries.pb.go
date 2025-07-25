// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: datasets/v1/timeseries.proto

package datasetsv1

import (
	v1 "github.com/tilebox/tilebox-go/protogen/tilebox/v1"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
	reflect "reflect"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// TimeseriesDatasetChunk is a message that represents a chunk of a timeseries dataset.
// used by workflow tasks that are executed once for each datapoint in a timeseries dataset
type TimeseriesDatasetChunk struct {
	state                            protoimpl.MessageState `protogen:"opaque.v1"`
	xxx_hidden_DatasetId             *v1.ID                 `protobuf:"bytes,1,opt,name=dataset_id,json=datasetId"`
	xxx_hidden_CollectionId          *v1.ID                 `protobuf:"bytes,2,opt,name=collection_id,json=collectionId"`
	xxx_hidden_TimeInterval          *v1.TimeInterval       `protobuf:"bytes,3,opt,name=time_interval,json=timeInterval"`
	xxx_hidden_DatapointInterval     *v1.IDInterval         `protobuf:"bytes,4,opt,name=datapoint_interval,json=datapointInterval"`
	xxx_hidden_BranchFactor          int32                  `protobuf:"varint,5,opt,name=branch_factor,json=branchFactor"`
	xxx_hidden_ChunkSize             int32                  `protobuf:"varint,6,opt,name=chunk_size,json=chunkSize"`
	xxx_hidden_DatapointsPer_365Days int64                  `protobuf:"varint,7,opt,name=datapoints_per_365_days,json=datapointsPer365Days"`
	unknownFields                    protoimpl.UnknownFields
	sizeCache                        protoimpl.SizeCache
}

func (x *TimeseriesDatasetChunk) Reset() {
	*x = TimeseriesDatasetChunk{}
	mi := &file_datasets_v1_timeseries_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TimeseriesDatasetChunk) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TimeseriesDatasetChunk) ProtoMessage() {}

func (x *TimeseriesDatasetChunk) ProtoReflect() protoreflect.Message {
	mi := &file_datasets_v1_timeseries_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (x *TimeseriesDatasetChunk) GetDatasetId() *v1.ID {
	if x != nil {
		return x.xxx_hidden_DatasetId
	}
	return nil
}

func (x *TimeseriesDatasetChunk) GetCollectionId() *v1.ID {
	if x != nil {
		return x.xxx_hidden_CollectionId
	}
	return nil
}

func (x *TimeseriesDatasetChunk) GetTimeInterval() *v1.TimeInterval {
	if x != nil {
		return x.xxx_hidden_TimeInterval
	}
	return nil
}

func (x *TimeseriesDatasetChunk) GetDatapointInterval() *v1.IDInterval {
	if x != nil {
		return x.xxx_hidden_DatapointInterval
	}
	return nil
}

func (x *TimeseriesDatasetChunk) GetBranchFactor() int32 {
	if x != nil {
		return x.xxx_hidden_BranchFactor
	}
	return 0
}

func (x *TimeseriesDatasetChunk) GetChunkSize() int32 {
	if x != nil {
		return x.xxx_hidden_ChunkSize
	}
	return 0
}

func (x *TimeseriesDatasetChunk) GetDatapointsPer_365Days() int64 {
	if x != nil {
		return x.xxx_hidden_DatapointsPer_365Days
	}
	return 0
}

func (x *TimeseriesDatasetChunk) SetDatasetId(v *v1.ID) {
	x.xxx_hidden_DatasetId = v
}

func (x *TimeseriesDatasetChunk) SetCollectionId(v *v1.ID) {
	x.xxx_hidden_CollectionId = v
}

func (x *TimeseriesDatasetChunk) SetTimeInterval(v *v1.TimeInterval) {
	x.xxx_hidden_TimeInterval = v
}

func (x *TimeseriesDatasetChunk) SetDatapointInterval(v *v1.IDInterval) {
	x.xxx_hidden_DatapointInterval = v
}

func (x *TimeseriesDatasetChunk) SetBranchFactor(v int32) {
	x.xxx_hidden_BranchFactor = v
}

func (x *TimeseriesDatasetChunk) SetChunkSize(v int32) {
	x.xxx_hidden_ChunkSize = v
}

func (x *TimeseriesDatasetChunk) SetDatapointsPer_365Days(v int64) {
	x.xxx_hidden_DatapointsPer_365Days = v
}

func (x *TimeseriesDatasetChunk) HasDatasetId() bool {
	if x == nil {
		return false
	}
	return x.xxx_hidden_DatasetId != nil
}

func (x *TimeseriesDatasetChunk) HasCollectionId() bool {
	if x == nil {
		return false
	}
	return x.xxx_hidden_CollectionId != nil
}

func (x *TimeseriesDatasetChunk) HasTimeInterval() bool {
	if x == nil {
		return false
	}
	return x.xxx_hidden_TimeInterval != nil
}

func (x *TimeseriesDatasetChunk) HasDatapointInterval() bool {
	if x == nil {
		return false
	}
	return x.xxx_hidden_DatapointInterval != nil
}

func (x *TimeseriesDatasetChunk) ClearDatasetId() {
	x.xxx_hidden_DatasetId = nil
}

func (x *TimeseriesDatasetChunk) ClearCollectionId() {
	x.xxx_hidden_CollectionId = nil
}

func (x *TimeseriesDatasetChunk) ClearTimeInterval() {
	x.xxx_hidden_TimeInterval = nil
}

func (x *TimeseriesDatasetChunk) ClearDatapointInterval() {
	x.xxx_hidden_DatapointInterval = nil
}

type TimeseriesDatasetChunk_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

	DatasetId             *v1.ID
	CollectionId          *v1.ID
	TimeInterval          *v1.TimeInterval
	DatapointInterval     *v1.IDInterval
	BranchFactor          int32
	ChunkSize             int32
	DatapointsPer_365Days int64
}

func (b0 TimeseriesDatasetChunk_builder) Build() *TimeseriesDatasetChunk {
	m0 := &TimeseriesDatasetChunk{}
	b, x := &b0, m0
	_, _ = b, x
	x.xxx_hidden_DatasetId = b.DatasetId
	x.xxx_hidden_CollectionId = b.CollectionId
	x.xxx_hidden_TimeInterval = b.TimeInterval
	x.xxx_hidden_DatapointInterval = b.DatapointInterval
	x.xxx_hidden_BranchFactor = b.BranchFactor
	x.xxx_hidden_ChunkSize = b.ChunkSize
	x.xxx_hidden_DatapointsPer_365Days = b.DatapointsPer_365Days
	return m0
}

// TimeChunk is a message that represents a time interval and a chunk size.
// used by workflow tasks that are executed once for each chunk of size chunk_size in a time interval
// e.g. for a time interval of 100 days, and a chunk size of 1 day, such a workflow will spawn a tree of
// eventually 100 leaf tasks
type TimeChunk struct {
	state                   protoimpl.MessageState `protogen:"opaque.v1"`
	xxx_hidden_TimeInterval *v1.TimeInterval       `protobuf:"bytes,1,opt,name=time_interval,json=timeInterval"`
	xxx_hidden_ChunkSize    *durationpb.Duration   `protobuf:"bytes,2,opt,name=chunk_size,json=chunkSize"`
	unknownFields           protoimpl.UnknownFields
	sizeCache               protoimpl.SizeCache
}

func (x *TimeChunk) Reset() {
	*x = TimeChunk{}
	mi := &file_datasets_v1_timeseries_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TimeChunk) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TimeChunk) ProtoMessage() {}

func (x *TimeChunk) ProtoReflect() protoreflect.Message {
	mi := &file_datasets_v1_timeseries_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (x *TimeChunk) GetTimeInterval() *v1.TimeInterval {
	if x != nil {
		return x.xxx_hidden_TimeInterval
	}
	return nil
}

func (x *TimeChunk) GetChunkSize() *durationpb.Duration {
	if x != nil {
		return x.xxx_hidden_ChunkSize
	}
	return nil
}

func (x *TimeChunk) SetTimeInterval(v *v1.TimeInterval) {
	x.xxx_hidden_TimeInterval = v
}

func (x *TimeChunk) SetChunkSize(v *durationpb.Duration) {
	x.xxx_hidden_ChunkSize = v
}

func (x *TimeChunk) HasTimeInterval() bool {
	if x == nil {
		return false
	}
	return x.xxx_hidden_TimeInterval != nil
}

func (x *TimeChunk) HasChunkSize() bool {
	if x == nil {
		return false
	}
	return x.xxx_hidden_ChunkSize != nil
}

func (x *TimeChunk) ClearTimeInterval() {
	x.xxx_hidden_TimeInterval = nil
}

func (x *TimeChunk) ClearChunkSize() {
	x.xxx_hidden_ChunkSize = nil
}

type TimeChunk_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

	TimeInterval *v1.TimeInterval
	ChunkSize    *durationpb.Duration
}

func (b0 TimeChunk_builder) Build() *TimeChunk {
	m0 := &TimeChunk{}
	b, x := &b0, m0
	_, _ = b, x
	x.xxx_hidden_TimeInterval = b.TimeInterval
	x.xxx_hidden_ChunkSize = b.ChunkSize
	return m0
}

var File_datasets_v1_timeseries_proto protoreflect.FileDescriptor

const file_datasets_v1_timeseries_proto_rawDesc = "" +
	"\n" +
	"\x1cdatasets/v1/timeseries.proto\x12\vdatasets.v1\x1a\x1egoogle/protobuf/duration.proto\x1a\x13tilebox/v1/id.proto\x1a\x16tilebox/v1/query.proto\"\xfd\x02\n" +
	"\x16TimeseriesDatasetChunk\x12-\n" +
	"\n" +
	"dataset_id\x18\x01 \x01(\v2\x0e.tilebox.v1.IDR\tdatasetId\x123\n" +
	"\rcollection_id\x18\x02 \x01(\v2\x0e.tilebox.v1.IDR\fcollectionId\x12=\n" +
	"\rtime_interval\x18\x03 \x01(\v2\x18.tilebox.v1.TimeIntervalR\ftimeInterval\x12E\n" +
	"\x12datapoint_interval\x18\x04 \x01(\v2\x16.tilebox.v1.IDIntervalR\x11datapointInterval\x12#\n" +
	"\rbranch_factor\x18\x05 \x01(\x05R\fbranchFactor\x12\x1d\n" +
	"\n" +
	"chunk_size\x18\x06 \x01(\x05R\tchunkSize\x125\n" +
	"\x17datapoints_per_365_days\x18\a \x01(\x03R\x14datapointsPer365Days\"\x84\x01\n" +
	"\tTimeChunk\x12=\n" +
	"\rtime_interval\x18\x01 \x01(\v2\x18.tilebox.v1.TimeIntervalR\ftimeInterval\x128\n" +
	"\n" +
	"chunk_size\x18\x02 \x01(\v2\x19.google.protobuf.DurationR\tchunkSizeB\xb3\x01\n" +
	"\x0fcom.datasets.v1B\x0fTimeseriesProtoP\x01Z=github.com/tilebox/tilebox-go/protogen/datasets/v1;datasetsv1\xa2\x02\x03DXX\xaa\x02\vDatasets.V1\xca\x02\vDatasets\\V1\xe2\x02\x17Datasets\\V1\\GPBMetadata\xea\x02\fDatasets::V1\x92\x03\x02\b\x02b\beditionsp\xe8\a"

var file_datasets_v1_timeseries_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_datasets_v1_timeseries_proto_goTypes = []any{
	(*TimeseriesDatasetChunk)(nil), // 0: datasets.v1.TimeseriesDatasetChunk
	(*TimeChunk)(nil),              // 1: datasets.v1.TimeChunk
	(*v1.ID)(nil),                  // 2: tilebox.v1.ID
	(*v1.TimeInterval)(nil),        // 3: tilebox.v1.TimeInterval
	(*v1.IDInterval)(nil),          // 4: tilebox.v1.IDInterval
	(*durationpb.Duration)(nil),    // 5: google.protobuf.Duration
}
var file_datasets_v1_timeseries_proto_depIdxs = []int32{
	2, // 0: datasets.v1.TimeseriesDatasetChunk.dataset_id:type_name -> tilebox.v1.ID
	2, // 1: datasets.v1.TimeseriesDatasetChunk.collection_id:type_name -> tilebox.v1.ID
	3, // 2: datasets.v1.TimeseriesDatasetChunk.time_interval:type_name -> tilebox.v1.TimeInterval
	4, // 3: datasets.v1.TimeseriesDatasetChunk.datapoint_interval:type_name -> tilebox.v1.IDInterval
	3, // 4: datasets.v1.TimeChunk.time_interval:type_name -> tilebox.v1.TimeInterval
	5, // 5: datasets.v1.TimeChunk.chunk_size:type_name -> google.protobuf.Duration
	6, // [6:6] is the sub-list for method output_type
	6, // [6:6] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_datasets_v1_timeseries_proto_init() }
func file_datasets_v1_timeseries_proto_init() {
	if File_datasets_v1_timeseries_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_datasets_v1_timeseries_proto_rawDesc), len(file_datasets_v1_timeseries_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_datasets_v1_timeseries_proto_goTypes,
		DependencyIndexes: file_datasets_v1_timeseries_proto_depIdxs,
		MessageInfos:      file_datasets_v1_timeseries_proto_msgTypes,
	}.Build()
	File_datasets_v1_timeseries_proto = out.File
	file_datasets_v1_timeseries_proto_goTypes = nil
	file_datasets_v1_timeseries_proto_depIdxs = nil
}
