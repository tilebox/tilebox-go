package field // import "github.com/tilebox/tilebox-go/datasets/v1/field"

import (
	"fmt"

	datasetsv1 "github.com/tilebox/tilebox-go/protogen/datasets/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// protoTypeName returns the fully qualified protobuf type name with a leading dot
func protoTypeName(message proto.Message) *string {
	return proto.String(fmt.Sprintf(".%s", message.ProtoReflect().Descriptor().FullName()))
}

type typeInfo struct {
	Type     descriptorpb.FieldDescriptorProto_Type
	TypeName *string // should be nil for scalar types
}

// String returns a new Field with type string.
func String(name string) *Descriptor {
	return &Descriptor{
		name: name,
		info: &typeInfo{
			Type: descriptorpb.FieldDescriptorProto_TYPE_STRING,
		},
	}
}

// Bytes returns a new Field with type bytes.
func Bytes(name string) *Descriptor {
	return &Descriptor{
		name: name,
		info: &typeInfo{
			Type: descriptorpb.FieldDescriptorProto_TYPE_BYTES,
		},
	}
}

// Bool returns a new Field with type bool.
func Bool(name string) *Descriptor {
	return &Descriptor{
		name: name,
		info: &typeInfo{
			Type: descriptorpb.FieldDescriptorProto_TYPE_BOOL,
		},
	}
}

// Int64 returns a new Field with type int64.
func Int64(name string) *Descriptor {
	return &Descriptor{
		name: name,
		info: &typeInfo{
			Type: descriptorpb.FieldDescriptorProto_TYPE_INT64,
		},
	}
}

// Uint64 returns a new Field with type uint64.
func Uint64(name string) *Descriptor {
	return &Descriptor{
		name: name,
		info: &typeInfo{
			Type: descriptorpb.FieldDescriptorProto_TYPE_UINT64,
		},
	}
}

// Float64 returns a new Field with type float64.
func Float64(name string) *Descriptor {
	return &Descriptor{
		name: name,
		info: &typeInfo{
			Type: descriptorpb.FieldDescriptorProto_TYPE_DOUBLE,
		},
	}
}

// Duration returns a new Field with type duration.
func Duration(name string) *Descriptor {
	return &Descriptor{
		name: name,
		info: &typeInfo{
			Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE,
			TypeName: protoTypeName(&durationpb.Duration{}),
		},
	}
}

// Timestamp returns a new Field with type timestamp.
func Timestamp(name string) *Descriptor {
	return &Descriptor{
		name: name,
		info: &typeInfo{
			Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE,
			TypeName: protoTypeName(&timestamppb.Timestamp{}),
		},
	}
}

// UUID returns a new Field with type UUID.
func UUID(name string) *Descriptor {
	return &Descriptor{
		name: name,
		info: &typeInfo{
			Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE,
			TypeName: protoTypeName(&datasetsv1.UUID{}),
		},
	}
}

// Geometry returns a new Field with type Geometry.
func Geometry(name string) *Descriptor {
	return &Descriptor{
		name: name,
		info: &typeInfo{
			Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE,
			TypeName: protoTypeName(&datasetsv1.Geometry{}),
		},
	}
}

// A Descriptor for field configuration.
// The usage is as follows:
//
//	[]field.Descriptor{
//	    field.Int64("field_name"),
//	}
type Descriptor struct {
	name         string
	info         *typeInfo
	description  string
	exampleValue string
	repeated     bool
}

// Description can be used to provide more context and details about the data. Optional.
func (b *Descriptor) Description(description string) *Descriptor {
	b.description = description
	return b
}

// ExampleValue can be used to provide an example value for documentation purposes. Optional.
func (b *Descriptor) ExampleValue(exampleValue string) *Descriptor {
	b.exampleValue = exampleValue
	return b
}

// Repeated indicates that this field is an array. Defaults to false.
func (b *Descriptor) Repeated() *Descriptor {
	b.repeated = true
	return b
}

func (b *Descriptor) ToProto() *datasetsv1.Field {
	label := descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL
	if b.repeated {
		label = descriptorpb.FieldDescriptorProto_LABEL_REPEATED
	}

	return datasetsv1.Field_builder{
		Descriptor: &descriptorpb.FieldDescriptorProto{
			Name:     &b.name,
			Type:     &b.info.Type,
			TypeName: b.info.TypeName,
			Label:    &label,
		},
		Annotation: datasetsv1.FieldAnnotation_builder{
			Description:  b.description,
			ExampleValue: b.exampleValue,
		}.Build(),
		Queryable: false,
	}.Build()
}
