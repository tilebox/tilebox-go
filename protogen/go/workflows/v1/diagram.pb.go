// Diagram service for rendering diagrams

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.4
// 	protoc        (unknown)
// source: workflows/v1/diagram.proto

package workflowsv1

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

// Request to render a diagram
type RenderDiagramRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Diagram       string                 `protobuf:"bytes,1,opt,name=diagram,proto3" json:"diagram,omitempty"`                                  // The diagram graph in the D2 syntax
	RenderOptions *RenderOptions         `protobuf:"bytes,2,opt,name=render_options,json=renderOptions,proto3" json:"render_options,omitempty"` // The options for rendering the diagram
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RenderDiagramRequest) Reset() {
	*x = RenderDiagramRequest{}
	mi := &file_workflows_v1_diagram_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RenderDiagramRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RenderDiagramRequest) ProtoMessage() {}

func (x *RenderDiagramRequest) ProtoReflect() protoreflect.Message {
	mi := &file_workflows_v1_diagram_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RenderDiagramRequest.ProtoReflect.Descriptor instead.
func (*RenderDiagramRequest) Descriptor() ([]byte, []int) {
	return file_workflows_v1_diagram_proto_rawDescGZIP(), []int{0}
}

func (x *RenderDiagramRequest) GetDiagram() string {
	if x != nil {
		return x.Diagram
	}
	return ""
}

func (x *RenderDiagramRequest) GetRenderOptions() *RenderOptions {
	if x != nil {
		return x.RenderOptions
	}
	return nil
}

// Options for rendering the diagram
type RenderOptions struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The layout to use for rendering the diagram: https://d2lang.com/tour/layouts/.
	// "dagre" or "elk". Defaults to "dagre"
	Layout string `protobuf:"bytes,1,opt,name=layout,proto3" json:"layout,omitempty"`
	// The D2 theme to use when rendering the diagram: https://d2lang.com/tour/themes/
	ThemeId *int64 `protobuf:"varint,2,opt,name=theme_id,json=themeId,proto3,oneof" json:"theme_id,omitempty"`
	// Whether to render the diagram in a sketchy (hand-drawn) style
	Sketchy bool `protobuf:"varint,3,opt,name=sketchy,proto3" json:"sketchy,omitempty"`
	// The padding around the diagram
	Padding int64 `protobuf:"varint,4,opt,name=padding,proto3" json:"padding,omitempty"`
	// Set explicitly the direction of the diagram: https://d2lang.com/tour/layouts/#direction.
	// "up", "down", "right", "left".
	Direction     string `protobuf:"bytes,5,opt,name=direction,proto3" json:"direction,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RenderOptions) Reset() {
	*x = RenderOptions{}
	mi := &file_workflows_v1_diagram_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RenderOptions) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RenderOptions) ProtoMessage() {}

func (x *RenderOptions) ProtoReflect() protoreflect.Message {
	mi := &file_workflows_v1_diagram_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RenderOptions.ProtoReflect.Descriptor instead.
func (*RenderOptions) Descriptor() ([]byte, []int) {
	return file_workflows_v1_diagram_proto_rawDescGZIP(), []int{1}
}

func (x *RenderOptions) GetLayout() string {
	if x != nil {
		return x.Layout
	}
	return ""
}

func (x *RenderOptions) GetThemeId() int64 {
	if x != nil && x.ThemeId != nil {
		return *x.ThemeId
	}
	return 0
}

func (x *RenderOptions) GetSketchy() bool {
	if x != nil {
		return x.Sketchy
	}
	return false
}

func (x *RenderOptions) GetPadding() int64 {
	if x != nil {
		return x.Padding
	}
	return 0
}

func (x *RenderOptions) GetDirection() string {
	if x != nil {
		return x.Direction
	}
	return ""
}

// A rendered diagram
type Diagram struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Svg           []byte                 `protobuf:"bytes,1,opt,name=svg,proto3" json:"svg,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Diagram) Reset() {
	*x = Diagram{}
	mi := &file_workflows_v1_diagram_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Diagram) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Diagram) ProtoMessage() {}

func (x *Diagram) ProtoReflect() protoreflect.Message {
	mi := &file_workflows_v1_diagram_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Diagram.ProtoReflect.Descriptor instead.
func (*Diagram) Descriptor() ([]byte, []int) {
	return file_workflows_v1_diagram_proto_rawDescGZIP(), []int{2}
}

func (x *Diagram) GetSvg() []byte {
	if x != nil {
		return x.Svg
	}
	return nil
}

var File_workflows_v1_diagram_proto protoreflect.FileDescriptor

var file_workflows_v1_diagram_proto_rawDesc = string([]byte{
	0x0a, 0x1a, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x2f, 0x76, 0x31, 0x2f, 0x64,
	0x69, 0x61, 0x67, 0x72, 0x61, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x77, 0x6f,
	0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x2e, 0x76, 0x31, 0x1a, 0x17, 0x77, 0x6f, 0x72, 0x6b,
	0x66, 0x6c, 0x6f, 0x77, 0x73, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x74, 0x0a, 0x14, 0x52, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x44, 0x69, 0x61,
	0x67, 0x72, 0x61, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x64,
	0x69, 0x61, 0x67, 0x72, 0x61, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x64, 0x69,
	0x61, 0x67, 0x72, 0x61, 0x6d, 0x12, 0x42, 0x0a, 0x0e, 0x72, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x5f,
	0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e,
	0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x6e,
	0x64, 0x65, 0x72, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x0d, 0x72, 0x65, 0x6e, 0x64,
	0x65, 0x72, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x22, 0xa6, 0x01, 0x0a, 0x0d, 0x52, 0x65,
	0x6e, 0x64, 0x65, 0x72, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x6c,
	0x61, 0x79, 0x6f, 0x75, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6c, 0x61, 0x79,
	0x6f, 0x75, 0x74, 0x12, 0x1e, 0x0a, 0x08, 0x74, 0x68, 0x65, 0x6d, 0x65, 0x5f, 0x69, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x03, 0x48, 0x00, 0x52, 0x07, 0x74, 0x68, 0x65, 0x6d, 0x65, 0x49, 0x64,
	0x88, 0x01, 0x01, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x6b, 0x65, 0x74, 0x63, 0x68, 0x79, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x6b, 0x65, 0x74, 0x63, 0x68, 0x79, 0x12, 0x18, 0x0a,
	0x07, 0x70, 0x61, 0x64, 0x64, 0x69, 0x6e, 0x67, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07,
	0x70, 0x61, 0x64, 0x64, 0x69, 0x6e, 0x67, 0x12, 0x1c, 0x0a, 0x09, 0x64, 0x69, 0x72, 0x65, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x64, 0x69, 0x72, 0x65,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x0b, 0x0a, 0x09, 0x5f, 0x74, 0x68, 0x65, 0x6d, 0x65, 0x5f,
	0x69, 0x64, 0x22, 0x1b, 0x0a, 0x07, 0x44, 0x69, 0x61, 0x67, 0x72, 0x61, 0x6d, 0x12, 0x10, 0x0a,
	0x03, 0x73, 0x76, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x73, 0x76, 0x67, 0x32,
	0x55, 0x0a, 0x0e, 0x44, 0x69, 0x61, 0x67, 0x72, 0x61, 0x6d, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x12, 0x43, 0x0a, 0x06, 0x52, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x12, 0x22, 0x2e, 0x77, 0x6f,
	0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x6e, 0x64, 0x65,
	0x72, 0x44, 0x69, 0x61, 0x67, 0x72, 0x61, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x15, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x44,
	0x69, 0x61, 0x67, 0x72, 0x61, 0x6d, 0x42, 0xb5, 0x01, 0x0a, 0x10, 0x63, 0x6f, 0x6d, 0x2e, 0x77,
	0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x2e, 0x76, 0x31, 0x42, 0x0c, 0x44, 0x69, 0x61,
	0x67, 0x72, 0x61, 0x6d, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x42, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x74, 0x69, 0x6c, 0x65, 0x62, 0x6f, 0x78, 0x2f,
	0x74, 0x69, 0x6c, 0x65, 0x62, 0x6f, 0x78, 0x2d, 0x67, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x67, 0x65, 0x6e, 0x2f, 0x67, 0x6f, 0x2f, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73,
	0x2f, 0x76, 0x31, 0x3b, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x76, 0x31, 0xa2,
	0x02, 0x03, 0x57, 0x58, 0x58, 0xaa, 0x02, 0x0c, 0x57, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77,
	0x73, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x0c, 0x57, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73,
	0x5c, 0x56, 0x31, 0xe2, 0x02, 0x18, 0x57, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x5c,
	0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02,
	0x0d, 0x57, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_workflows_v1_diagram_proto_rawDescOnce sync.Once
	file_workflows_v1_diagram_proto_rawDescData []byte
)

func file_workflows_v1_diagram_proto_rawDescGZIP() []byte {
	file_workflows_v1_diagram_proto_rawDescOnce.Do(func() {
		file_workflows_v1_diagram_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_workflows_v1_diagram_proto_rawDesc), len(file_workflows_v1_diagram_proto_rawDesc)))
	})
	return file_workflows_v1_diagram_proto_rawDescData
}

var file_workflows_v1_diagram_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_workflows_v1_diagram_proto_goTypes = []any{
	(*RenderDiagramRequest)(nil), // 0: workflows.v1.RenderDiagramRequest
	(*RenderOptions)(nil),        // 1: workflows.v1.RenderOptions
	(*Diagram)(nil),              // 2: workflows.v1.Diagram
}
var file_workflows_v1_diagram_proto_depIdxs = []int32{
	1, // 0: workflows.v1.RenderDiagramRequest.render_options:type_name -> workflows.v1.RenderOptions
	0, // 1: workflows.v1.DiagramService.Render:input_type -> workflows.v1.RenderDiagramRequest
	2, // 2: workflows.v1.DiagramService.Render:output_type -> workflows.v1.Diagram
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_workflows_v1_diagram_proto_init() }
func file_workflows_v1_diagram_proto_init() {
	if File_workflows_v1_diagram_proto != nil {
		return
	}
	file_workflows_v1_core_proto_init()
	file_workflows_v1_diagram_proto_msgTypes[1].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_workflows_v1_diagram_proto_rawDesc), len(file_workflows_v1_diagram_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_workflows_v1_diagram_proto_goTypes,
		DependencyIndexes: file_workflows_v1_diagram_proto_depIdxs,
		MessageInfos:      file_workflows_v1_diagram_proto_msgTypes,
	}.Build()
	File_workflows_v1_diagram_proto = out.File
	file_workflows_v1_diagram_proto_goTypes = nil
	file_workflows_v1_diagram_proto_depIdxs = nil
}
