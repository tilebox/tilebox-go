// This file contains the proto messages describing a dataset type. A dataset type is akin to a protobuf message type,
// and can be converted to it.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.4
// 	protoc        (unknown)
// source: datasets/v1/dataset_type.proto

package datasetsv1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	descriptorpb "google.golang.org/protobuf/types/descriptorpb"
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

// DatasetKind is an enum describing the kind of dataset. A dataset kind specifies a set of default fields, that
// cannot be modified and will always be present for every DatasetMessageType of that kind.
type DatasetKind int32

const (
	DatasetKind_DATASET_KIND_UNSPECIFIED DatasetKind = 0
	// A temporal dataset is a dataset that contains a "time" field, which is a timestamp.
	DatasetKind_DATASET_KIND_TEMPORAL DatasetKind = 1
	// A spatio temporal dataset is a dataset that contains a "start_time" and "end_time" field which are timestamps
	// as well as a "geometry" field which is Polygon or MultiPolygon.
	DatasetKind_DATASET_KIND_SPATIOTEMPORAL DatasetKind = 2
)

// Enum value maps for DatasetKind.
var (
	DatasetKind_name = map[int32]string{
		0: "DATASET_KIND_UNSPECIFIED",
		1: "DATASET_KIND_TEMPORAL",
		2: "DATASET_KIND_SPATIOTEMPORAL",
	}
	DatasetKind_value = map[string]int32{
		"DATASET_KIND_UNSPECIFIED":    0,
		"DATASET_KIND_TEMPORAL":       1,
		"DATASET_KIND_SPATIOTEMPORAL": 2,
	}
)

func (x DatasetKind) Enum() *DatasetKind {
	p := new(DatasetKind)
	*p = x
	return p
}

func (x DatasetKind) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (DatasetKind) Descriptor() protoreflect.EnumDescriptor {
	return file_datasets_v1_dataset_type_proto_enumTypes[0].Descriptor()
}

func (DatasetKind) Type() protoreflect.EnumType {
	return &file_datasets_v1_dataset_type_proto_enumTypes[0]
}

func (x DatasetKind) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use DatasetKind.Descriptor instead.
func (DatasetKind) EnumDescriptor() ([]byte, []int) {
	return file_datasets_v1_dataset_type_proto_rawDescGZIP(), []int{0}
}

type Field struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The descriptor contains the name of the field, the type, optional labels (e.g. repeated) and other information.
	// If the type is TYPE_MESSAGE, then the type_name must be a fully qualified name to a well known type, e.g.
	// `datasets.v1.Vec3` or `google.protobuf.Timestamp`.
	Descriptor_ *descriptorpb.FieldDescriptorProto `protobuf:"bytes,1,opt,name=descriptor,proto3" json:"descriptor,omitempty"`
	// An optional description and example value for the field.
	Annotation *FieldAnnotation `protobuf:"bytes,2,opt,name=annotation,proto3" json:"annotation,omitempty"`
	// A flag indicating whether the field should be queryable. This means we will build an index for the field, and
	// allow users to query for certain values server-side.
	Queryable     bool `protobuf:"varint,3,opt,name=queryable,proto3" json:"queryable,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Field) Reset() {
	*x = Field{}
	mi := &file_datasets_v1_dataset_type_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Field) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Field) ProtoMessage() {}

func (x *Field) ProtoReflect() protoreflect.Message {
	mi := &file_datasets_v1_dataset_type_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Field.ProtoReflect.Descriptor instead.
func (*Field) Descriptor() ([]byte, []int) {
	return file_datasets_v1_dataset_type_proto_rawDescGZIP(), []int{0}
}

func (x *Field) GetDescriptor_() *descriptorpb.FieldDescriptorProto {
	if x != nil {
		return x.Descriptor_
	}
	return nil
}

func (x *Field) GetAnnotation() *FieldAnnotation {
	if x != nil {
		return x.Annotation
	}
	return nil
}

func (x *Field) GetQueryable() bool {
	if x != nil {
		return x.Queryable
	}
	return false
}

type DatasetType struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// kind denotes the kind of dataset this type describes. We do not rely on the default fields to be set in our
	// array of fields - since that way users could circumvent those by manually sending requests not from our Console.
	Kind DatasetKind `protobuf:"varint,1,opt,name=kind,proto3,enum=datasets.v1.DatasetKind" json:"kind,omitempty"`
	// A list of fields that this dataset consists of.
	Fields        []*Field `protobuf:"bytes,2,rep,name=fields,proto3" json:"fields,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DatasetType) Reset() {
	*x = DatasetType{}
	mi := &file_datasets_v1_dataset_type_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DatasetType) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DatasetType) ProtoMessage() {}

func (x *DatasetType) ProtoReflect() protoreflect.Message {
	mi := &file_datasets_v1_dataset_type_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DatasetType.ProtoReflect.Descriptor instead.
func (*DatasetType) Descriptor() ([]byte, []int) {
	return file_datasets_v1_dataset_type_proto_rawDescGZIP(), []int{1}
}

func (x *DatasetType) GetKind() DatasetKind {
	if x != nil {
		return x.Kind
	}
	return DatasetKind_DATASET_KIND_UNSPECIFIED
}

func (x *DatasetType) GetFields() []*Field {
	if x != nil {
		return x.Fields
	}
	return nil
}

var File_datasets_v1_dataset_type_proto protoreflect.FileDescriptor

var file_datasets_v1_dataset_type_proto_rawDesc = string([]byte{
	0x0a, 0x1e, 0x64, 0x61, 0x74, 0x61, 0x73, 0x65, 0x74, 0x73, 0x2f, 0x76, 0x31, 0x2f, 0x64, 0x61,
	0x74, 0x61, 0x73, 0x65, 0x74, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x0b, 0x64, 0x61, 0x74, 0x61, 0x73, 0x65, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x1a, 0x20, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64,
	0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x16, 0x64, 0x61, 0x74, 0x61, 0x73, 0x65, 0x74, 0x73, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x6f, 0x72,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xaa, 0x01, 0x0a, 0x05, 0x46, 0x69, 0x65, 0x6c,
	0x64, 0x12, 0x45, 0x0a, 0x0a, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x44, 0x65, 0x73,
	0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x52, 0x0a, 0x64, 0x65,
	0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x12, 0x3c, 0x0a, 0x0a, 0x61, 0x6e, 0x6e, 0x6f,
	0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x64,
	0x61, 0x74, 0x61, 0x73, 0x65, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64,
	0x41, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0a, 0x61, 0x6e, 0x6e, 0x6f,
	0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1c, 0x0a, 0x09, 0x71, 0x75, 0x65, 0x72, 0x79, 0x61,
	0x62, 0x6c, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x09, 0x71, 0x75, 0x65, 0x72, 0x79,
	0x61, 0x62, 0x6c, 0x65, 0x22, 0x67, 0x0a, 0x0b, 0x44, 0x61, 0x74, 0x61, 0x73, 0x65, 0x74, 0x54,
	0x79, 0x70, 0x65, 0x12, 0x2c, 0x0a, 0x04, 0x6b, 0x69, 0x6e, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x18, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x73, 0x65, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x2e,
	0x44, 0x61, 0x74, 0x61, 0x73, 0x65, 0x74, 0x4b, 0x69, 0x6e, 0x64, 0x52, 0x04, 0x6b, 0x69, 0x6e,
	0x64, 0x12, 0x2a, 0x0a, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x12, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x73, 0x65, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x2e,
	0x46, 0x69, 0x65, 0x6c, 0x64, 0x52, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x2a, 0x67, 0x0a,
	0x0b, 0x44, 0x61, 0x74, 0x61, 0x73, 0x65, 0x74, 0x4b, 0x69, 0x6e, 0x64, 0x12, 0x1c, 0x0a, 0x18,
	0x44, 0x41, 0x54, 0x41, 0x53, 0x45, 0x54, 0x5f, 0x4b, 0x49, 0x4e, 0x44, 0x5f, 0x55, 0x4e, 0x53,
	0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x19, 0x0a, 0x15, 0x44, 0x41,
	0x54, 0x41, 0x53, 0x45, 0x54, 0x5f, 0x4b, 0x49, 0x4e, 0x44, 0x5f, 0x54, 0x45, 0x4d, 0x50, 0x4f,
	0x52, 0x41, 0x4c, 0x10, 0x01, 0x12, 0x1f, 0x0a, 0x1b, 0x44, 0x41, 0x54, 0x41, 0x53, 0x45, 0x54,
	0x5f, 0x4b, 0x49, 0x4e, 0x44, 0x5f, 0x53, 0x50, 0x41, 0x54, 0x49, 0x4f, 0x54, 0x45, 0x4d, 0x50,
	0x4f, 0x52, 0x41, 0x4c, 0x10, 0x02, 0x42, 0xb2, 0x01, 0x0a, 0x0f, 0x63, 0x6f, 0x6d, 0x2e, 0x64,
	0x61, 0x74, 0x61, 0x73, 0x65, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x42, 0x10, 0x44, 0x61, 0x74, 0x61,
	0x73, 0x65, 0x74, 0x54, 0x79, 0x70, 0x65, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x40,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x74, 0x69, 0x6c, 0x65, 0x62,
	0x6f, 0x78, 0x2f, 0x74, 0x69, 0x6c, 0x65, 0x62, 0x6f, 0x78, 0x2d, 0x67, 0x6f, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x67, 0x65, 0x6e, 0x2f, 0x67, 0x6f, 0x2f, 0x64, 0x61, 0x74, 0x61, 0x73, 0x65,
	0x74, 0x73, 0x2f, 0x76, 0x31, 0x3b, 0x64, 0x61, 0x74, 0x61, 0x73, 0x65, 0x74, 0x73, 0x76, 0x31,
	0xa2, 0x02, 0x03, 0x44, 0x58, 0x58, 0xaa, 0x02, 0x0b, 0x44, 0x61, 0x74, 0x61, 0x73, 0x65, 0x74,
	0x73, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x0b, 0x44, 0x61, 0x74, 0x61, 0x73, 0x65, 0x74, 0x73, 0x5c,
	0x56, 0x31, 0xe2, 0x02, 0x17, 0x44, 0x61, 0x74, 0x61, 0x73, 0x65, 0x74, 0x73, 0x5c, 0x56, 0x31,
	0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x0c, 0x44,
	0x61, 0x74, 0x61, 0x73, 0x65, 0x74, 0x73, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
})

var (
	file_datasets_v1_dataset_type_proto_rawDescOnce sync.Once
	file_datasets_v1_dataset_type_proto_rawDescData []byte
)

func file_datasets_v1_dataset_type_proto_rawDescGZIP() []byte {
	file_datasets_v1_dataset_type_proto_rawDescOnce.Do(func() {
		file_datasets_v1_dataset_type_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_datasets_v1_dataset_type_proto_rawDesc), len(file_datasets_v1_dataset_type_proto_rawDesc)))
	})
	return file_datasets_v1_dataset_type_proto_rawDescData
}

var file_datasets_v1_dataset_type_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_datasets_v1_dataset_type_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_datasets_v1_dataset_type_proto_goTypes = []any{
	(DatasetKind)(0),    // 0: datasets.v1.DatasetKind
	(*Field)(nil),       // 1: datasets.v1.Field
	(*DatasetType)(nil), // 2: datasets.v1.DatasetType
	(*descriptorpb.FieldDescriptorProto)(nil), // 3: google.protobuf.FieldDescriptorProto
	(*FieldAnnotation)(nil),                   // 4: datasets.v1.FieldAnnotation
}
var file_datasets_v1_dataset_type_proto_depIdxs = []int32{
	3, // 0: datasets.v1.Field.descriptor:type_name -> google.protobuf.FieldDescriptorProto
	4, // 1: datasets.v1.Field.annotation:type_name -> datasets.v1.FieldAnnotation
	0, // 2: datasets.v1.DatasetType.kind:type_name -> datasets.v1.DatasetKind
	1, // 3: datasets.v1.DatasetType.fields:type_name -> datasets.v1.Field
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_datasets_v1_dataset_type_proto_init() }
func file_datasets_v1_dataset_type_proto_init() {
	if File_datasets_v1_dataset_type_proto != nil {
		return
	}
	file_datasets_v1_core_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_datasets_v1_dataset_type_proto_rawDesc), len(file_datasets_v1_dataset_type_proto_rawDesc)),
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_datasets_v1_dataset_type_proto_goTypes,
		DependencyIndexes: file_datasets_v1_dataset_type_proto_depIdxs,
		EnumInfos:         file_datasets_v1_dataset_type_proto_enumTypes,
		MessageInfos:      file_datasets_v1_dataset_type_proto_msgTypes,
	}.Build()
	File_datasets_v1_dataset_type_proto = out.File
	file_datasets_v1_dataset_type_proto_goTypes = nil
	file_datasets_v1_dataset_type_proto_depIdxs = nil
}
