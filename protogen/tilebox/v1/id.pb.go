// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: tilebox/v1/id.proto

package tileboxv1

import (
	_ "buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Bytes field (in message)
type ID struct {
	state           protoimpl.MessageState `protogen:"opaque.v1"`
	xxx_hidden_Uuid []byte                 `protobuf:"bytes,1,opt,name=uuid"`
	unknownFields   protoimpl.UnknownFields
	sizeCache       protoimpl.SizeCache
}

func (x *ID) Reset() {
	*x = ID{}
	mi := &file_tilebox_v1_id_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ID) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ID) ProtoMessage() {}

func (x *ID) ProtoReflect() protoreflect.Message {
	mi := &file_tilebox_v1_id_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (x *ID) GetUuid() []byte {
	if x != nil {
		return x.xxx_hidden_Uuid
	}
	return nil
}

func (x *ID) SetUuid(v []byte) {
	if v == nil {
		v = []byte{}
	}
	x.xxx_hidden_Uuid = v
}

type ID_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

	Uuid []byte
}

func (b0 ID_builder) Build() *ID {
	m0 := &ID{}
	b, x := &b0, m0
	_, _ = b, x
	x.xxx_hidden_Uuid = b.Uuid
	return m0
}

var File_tilebox_v1_id_proto protoreflect.FileDescriptor

const file_tilebox_v1_id_proto_rawDesc = "" +
	"\n" +
	"\x13tilebox/v1/id.proto\x12\n" +
	"tilebox.v1\x1a\x1bbuf/validate/validate.proto\"3\n" +
	"\x02ID\x12-\n" +
	"\x04uuid\x18\x01 \x01(\fB\x19\xbaH\x16z\x14J\x10\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00h\x10R\x04uuidB\xa4\x01\n" +
	"\x0ecom.tilebox.v1B\aIdProtoP\x01Z;github.com/tilebox/tilebox-go/protogen/tilebox/v1;tileboxv1\xa2\x02\x03TXX\xaa\x02\n" +
	"Tilebox.V1\xca\x02\n" +
	"Tilebox\\V1\xe2\x02\x16Tilebox\\V1\\GPBMetadata\xea\x02\vTilebox::V1\x92\x03\x02\b\x02b\beditionsp\xe8\a"

var file_tilebox_v1_id_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_tilebox_v1_id_proto_goTypes = []any{
	(*ID)(nil), // 0: tilebox.v1.ID
}
var file_tilebox_v1_id_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_tilebox_v1_id_proto_init() }
func file_tilebox_v1_id_proto_init() {
	if File_tilebox_v1_id_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_tilebox_v1_id_proto_rawDesc), len(file_tilebox_v1_id_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_tilebox_v1_id_proto_goTypes,
		DependencyIndexes: file_tilebox_v1_id_proto_depIdxs,
		MessageInfos:      file_tilebox_v1_id_proto_msgTypes,
	}.Build()
	File_tilebox_v1_id_proto = out.File
	file_tilebox_v1_id_proto_goTypes = nil
	file_tilebox_v1_id_proto_depIdxs = nil
}
