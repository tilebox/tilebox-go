package field

import (
	"fmt"
	"slices"

	datasetsv1 "github.com/tilebox/tilebox-go/protogen/datasets/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

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

// Int32 returns a new Field with type int32.
func Int32(name string) *Descriptor {
	return &Descriptor{
		name: name,
		info: &typeInfo{
			Type: descriptorpb.FieldDescriptorProto_TYPE_INT32,
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

// Message returns a new Field with the type of the provided protobuf message.
// The server only accepts a specific set of well-known message types.
func Message(name string, message proto.Message) *Descriptor {
	return &Descriptor{
		name: name,
		info: &typeInfo{
			Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE,
			TypeName: protoTypeName(message),
		},
	}
}

// Enum returns a new Field with the type of the provided protobuf enum value.
// The server only accepts a specific set of well-known enum types.
func Enum(name string, enum protoreflect.Enum) *Descriptor {
	typeName := fmt.Sprintf(".%s", enum.Type().Descriptor().FullName())
	return &Descriptor{
		name: name,
		info: &typeInfo{
			Type:     descriptorpb.FieldDescriptorProto_TYPE_ENUM,
			TypeName: &typeName,
		},
	}
}

// protoTypeName returns the fully qualified protobuf type name with a leading dot
func protoTypeName(message proto.Message) *string {
	return new(fmt.Sprintf(".%s", message.ProtoReflect().Descriptor().FullName()))
}

type typeInfo struct {
	Type     descriptorpb.FieldDescriptorProto_Type
	TypeName *string // should be nil for scalar types
}

// FieldRole describes a semantic display role fulfilled by a dataset field.
type FieldRole = datasetsv1.FieldRole

const (
	// RolePrimaryTitle marks the field as the primary human-readable title of a datapoint.
	RolePrimaryTitle FieldRole = datasetsv1.FieldRole_FIELD_ROLE_PRIMARY_TITLE
)

// Descriptor builds a dataset field's protobuf type and annotations.
type Descriptor struct {
	name              string
	info              *typeInfo
	description       string
	exampleValue      string
	sourceJSONPointer *string
	queryable         bool
	jsonSchemaRef     *string
	roles             []FieldRole
	repeated          bool
}

// Description can be used to provide more context and details about the data. Optional.
func (d *Descriptor) Description(description string) *Descriptor {
	d.description = description
	return d
}

// ExampleValue can be used to provide an example value for documentation purposes. Optional.
func (d *Descriptor) ExampleValue(exampleValue string) *Descriptor {
	d.exampleValue = exampleValue
	return d
}

// SourceJSONPointer sets the RFC 6901 JSON Pointer locating this field in the source document. Optional.
func (d *Descriptor) SourceJSONPointer(sourceJSONPointer string) *Descriptor {
	d.sourceJSONPointer = &sourceJSONPointer
	return d
}

// Queryable marks the field for projection into query storage and server-side filtering.
// Only optional fields created with String, Bool, Int32, Int64, Uint64, or Float64 can be queryable.
// Repeated fields and fields of all other types are not supported.
func (d *Descriptor) Queryable() *Descriptor {
	d.queryable = true
	return d
}

// JSONSchemaRef sets the JSON Schema reference advertised for this field. Optional.
func (d *Descriptor) JSONSchemaRef(jsonSchemaRef string) *Descriptor {
	d.jsonSchemaRef = &jsonSchemaRef
	return d
}

// Roles sets the semantic display roles fulfilled by this field. Optional.
func (d *Descriptor) Roles(roles ...FieldRole) *Descriptor {
	d.roles = slices.Clone(roles)
	return d
}

// Repeated indicates that this field is an array. Defaults to false.
func (d *Descriptor) Repeated() *Descriptor {
	d.repeated = true
	return d
}

// Descriptor implements the datasets.Field interface by returning itself.
func (d *Descriptor) Descriptor() *Descriptor {
	return d
}

func (d *Descriptor) ToProto() *datasetsv1.Field {
	label := descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL
	if d.repeated {
		label = descriptorpb.FieldDescriptorProto_LABEL_REPEATED
	}

	return datasetsv1.Field_builder{
		Descriptor: &descriptorpb.FieldDescriptorProto{
			Name:     &d.name,
			Type:     &d.info.Type,
			TypeName: d.info.TypeName,
			Label:    &label,
		},
		Annotation: datasetsv1.FieldAnnotation_builder{
			Description:       d.description,
			ExampleValue:      d.exampleValue,
			SourceJsonPointer: d.sourceJSONPointer,
			Queryable:         d.queryable,
			JsonSchemaRef:     d.jsonSchemaRef,
			Roles:             d.roles,
		}.Build(),
	}.Build()
}
