package field // import "github.com/tilebox/tilebox-go/datasets/v1/field"

import (
	"fmt"

	datasetsv1 "github.com/tilebox/tilebox-go/protogen/datasets/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// String returns a new Field with type string.
func String(name string) *stringBuilder {
	return &stringBuilder{&Descriptor{
		name: name,
		info: &typeInfo{
			Type: descriptorpb.FieldDescriptorProto_TYPE_STRING,
		},
	}}
}

// Bytes returns a new Field with type bytes.
func Bytes(name string) *bytesBuilder {
	return &bytesBuilder{&Descriptor{
		name: name,
		info: &typeInfo{
			Type: descriptorpb.FieldDescriptorProto_TYPE_BYTES,
		},
	}}
}

// Bool returns a new Field with type bool.
func Bool(name string) *boolBuilder {
	return &boolBuilder{&Descriptor{
		name: name,
		info: &typeInfo{
			Type: descriptorpb.FieldDescriptorProto_TYPE_BOOL,
		},
	}}
}

// Int64 returns a new Field with type int64.
func Int64(name string) *int64Builder {
	return &int64Builder{&Descriptor{
		name: name,
		info: &typeInfo{
			Type: descriptorpb.FieldDescriptorProto_TYPE_INT64,
		},
	}}
}

// Uint64 returns a new Field with type uint64.
func Uint64(name string) *uint64Builder {
	return &uint64Builder{&Descriptor{
		name: name,
		info: &typeInfo{
			Type: descriptorpb.FieldDescriptorProto_TYPE_UINT64,
		},
	}}
}

// Float64 returns a new Field with type float64.
func Float64(name string) *float64Builder {
	return &float64Builder{&Descriptor{
		name: name,
		info: &typeInfo{
			Type: descriptorpb.FieldDescriptorProto_TYPE_DOUBLE,
		},
	}}
}

// Duration returns a new Field with type duration.
func Duration(name string) *durationBuilder {
	return &durationBuilder{&Descriptor{
		name: name,
		info: &typeInfo{
			Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE,
			TypeName: protoTypeName(&durationpb.Duration{}),
		},
	}}
}

// Timestamp returns a new Field with type timestamp.
func Timestamp(name string) *timestampBuilder {
	return &timestampBuilder{&Descriptor{
		name: name,
		info: &typeInfo{
			Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE,
			TypeName: protoTypeName(&timestamppb.Timestamp{}),
		},
	}}
}

// UUID returns a new Field with type UUID.
func UUID(name string) *uuidBuilder {
	return &uuidBuilder{&Descriptor{
		name: name,
		info: &typeInfo{
			Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE,
			TypeName: protoTypeName(&datasetsv1.UUID{}),
		},
	}}
}

// Geometry returns a new Field with type Geometry.
func Geometry(name string) *geometryBuilder {
	return &geometryBuilder{&Descriptor{
		name: name,
		info: &typeInfo{
			Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE,
			TypeName: protoTypeName(&datasetsv1.Geometry{}),
		},
	}}
}

// protoTypeName returns the fully qualified protobuf type name with a leading dot
func protoTypeName(message proto.Message) *string {
	return proto.String(fmt.Sprintf(".%s", message.ProtoReflect().Descriptor().FullName()))
}

type typeInfo struct {
	Type     descriptorpb.FieldDescriptorProto_Type
	TypeName *string // should be nil for scalar types
}

// stringBuilder is the builder for string fields.
type stringBuilder struct {
	desc *Descriptor
}

// Description can be used to provide more context and details about the data. Optional.
func (b *stringBuilder) Description(description string) *stringBuilder {
	b.desc.description = description
	return b
}

// ExampleValue can be used to provide an example value for documentation purposes. Optional.
func (b *stringBuilder) ExampleValue(exampleValue string) *stringBuilder {
	b.desc.exampleValue = exampleValue
	return b
}

// Repeated indicates that this field is an array. Defaults to false.
func (b *stringBuilder) Repeated() *stringBuilder {
	b.desc.repeated = true
	return b
}

// Descriptor implements the datasets.Field interface by returning its descriptor.
func (b *stringBuilder) Descriptor() *Descriptor {
	return b.desc
}

// bytesBuilder is the builder for bytes fields.
type bytesBuilder struct {
	desc *Descriptor
}

// Description can be used to provide more context and details about the data. Optional.
func (b *bytesBuilder) Description(description string) *bytesBuilder {
	b.desc.description = description
	return b
}

// ExampleValue can be used to provide an example value for documentation purposes. Optional.
func (b *bytesBuilder) ExampleValue(exampleValue string) *bytesBuilder {
	b.desc.exampleValue = exampleValue
	return b
}

// Repeated indicates that this field is an array. Defaults to false.
func (b *bytesBuilder) Repeated() *bytesBuilder {
	b.desc.repeated = true
	return b
}

// Descriptor implements the datasets.Field interface by returning its descriptor.
func (b *bytesBuilder) Descriptor() *Descriptor {
	return b.desc
}

// boolBuilder is the builder for bool fields.
type boolBuilder struct {
	desc *Descriptor
}

// Description can be used to provide more context and details about the data. Optional.
func (b *boolBuilder) Description(description string) *boolBuilder {
	b.desc.description = description
	return b
}

// ExampleValue can be used to provide an example value for documentation purposes. Optional.
func (b *boolBuilder) ExampleValue(exampleValue string) *boolBuilder {
	b.desc.exampleValue = exampleValue
	return b
}

// Repeated indicates that this field is an array. Defaults to false.
func (b *boolBuilder) Repeated() *boolBuilder {
	b.desc.repeated = true
	return b
}

// Descriptor implements the datasets.Field interface by returning its descriptor.
func (b *boolBuilder) Descriptor() *Descriptor {
	return b.desc
}

// int64Builder is the builder for int64 fields.
type int64Builder struct {
	desc *Descriptor
}

// Description can be used to provide more context and details about the data. Optional.
func (b *int64Builder) Description(description string) *int64Builder {
	b.desc.description = description
	return b
}

// ExampleValue can be used to provide an example value for documentation purposes. Optional.
func (b *int64Builder) ExampleValue(exampleValue string) *int64Builder {
	b.desc.exampleValue = exampleValue
	return b
}

// Repeated indicates that this field is an array. Defaults to false.
func (b *int64Builder) Repeated() *int64Builder {
	b.desc.repeated = true
	return b
}

// Descriptor implements the datasets.Field interface by returning its descriptor.
func (b *int64Builder) Descriptor() *Descriptor {
	return b.desc
}

// uint64Builder is the builder for uint64 fields.
type uint64Builder struct {
	desc *Descriptor
}

// Description can be used to provide more context and details about the data. Optional.
func (b *uint64Builder) Description(description string) *uint64Builder {
	b.desc.description = description
	return b
}

// ExampleValue can be used to provide an example value for documentation purposes. Optional.
func (b *uint64Builder) ExampleValue(exampleValue string) *uint64Builder {
	b.desc.exampleValue = exampleValue
	return b
}

// Repeated indicates that this field is an array. Defaults to false.
func (b *uint64Builder) Repeated() *uint64Builder {
	b.desc.repeated = true
	return b
}

// Descriptor implements the datasets.Field interface by returning its descriptor.
func (b *uint64Builder) Descriptor() *Descriptor {
	return b.desc
}

// float64Builder is the builder for float64 fields.
type float64Builder struct {
	desc *Descriptor
}

// Description can be used to provide more context and details about the data. Optional.
func (b *float64Builder) Description(description string) *float64Builder {
	b.desc.description = description
	return b
}

// ExampleValue can be used to provide an example value for documentation purposes. Optional.
func (b *float64Builder) ExampleValue(exampleValue string) *float64Builder {
	b.desc.exampleValue = exampleValue
	return b
}

// Repeated indicates that this field is an array. Defaults to false.
func (b *float64Builder) Repeated() *float64Builder {
	b.desc.repeated = true
	return b
}

// Descriptor implements the datasets.Field interface by returning its descriptor.
func (b *float64Builder) Descriptor() *Descriptor {
	return b.desc
}

// durationBuilder is the builder for duration fields.
type durationBuilder struct {
	desc *Descriptor
}

// Description can be used to provide more context and details about the data. Optional.
func (b *durationBuilder) Description(description string) *durationBuilder {
	b.desc.description = description
	return b
}

// ExampleValue can be used to provide an example value for documentation purposes. Optional.
func (b *durationBuilder) ExampleValue(exampleValue string) *durationBuilder {
	b.desc.exampleValue = exampleValue
	return b
}

// Repeated indicates that this field is an array. Defaults to false.
func (b *durationBuilder) Repeated() *durationBuilder {
	b.desc.repeated = true
	return b
}

// Descriptor implements the datasets.Field interface by returning its descriptor.
func (b *durationBuilder) Descriptor() *Descriptor {
	return b.desc
}

// timestampBuilder is the builder for timestamp fields.
type timestampBuilder struct {
	desc *Descriptor
}

// Description can be used to provide more context and details about the data. Optional.
func (b *timestampBuilder) Description(description string) *timestampBuilder {
	b.desc.description = description
	return b
}

// ExampleValue can be used to provide an example value for documentation purposes. Optional.
func (b *timestampBuilder) ExampleValue(exampleValue string) *timestampBuilder {
	b.desc.exampleValue = exampleValue
	return b
}

// Repeated indicates that this field is an array. Defaults to false.
func (b *timestampBuilder) Repeated() *timestampBuilder {
	b.desc.repeated = true
	return b
}

// Descriptor implements the datasets.Field interface by returning its descriptor.
func (b *timestampBuilder) Descriptor() *Descriptor {
	return b.desc
}

// uuidBuilder is the builder for uuid fields.
type uuidBuilder struct {
	desc *Descriptor
}

// Description can be used to provide more context and details about the data. Optional.
func (b *uuidBuilder) Description(description string) *uuidBuilder {
	b.desc.description = description
	return b
}

// ExampleValue can be used to provide an example value for documentation purposes. Optional.
func (b *uuidBuilder) ExampleValue(exampleValue string) *uuidBuilder {
	b.desc.exampleValue = exampleValue
	return b
}

// Repeated indicates that this field is an array. Defaults to false.
func (b *uuidBuilder) Repeated() *uuidBuilder {
	b.desc.repeated = true
	return b
}

// Descriptor implements the datasets.Field interface by returning its descriptor.
func (b *uuidBuilder) Descriptor() *Descriptor {
	return b.desc
}

// geometryBuilder is the builder for geometry fields.
type geometryBuilder struct {
	desc *Descriptor
}

// Description can be used to provide more context and details about the data. Optional.
func (b *geometryBuilder) Description(description string) *geometryBuilder {
	b.desc.description = description
	return b
}

// ExampleValue can be used to provide an example value for documentation purposes. Optional.
func (b *geometryBuilder) ExampleValue(exampleValue string) *geometryBuilder {
	b.desc.exampleValue = exampleValue
	return b
}

// Repeated indicates that this field is an array. Defaults to false.
func (b *geometryBuilder) Repeated() *geometryBuilder {
	b.desc.repeated = true
	return b
}

// Descriptor implements the datasets.Field interface by returning its descriptor.
func (b *geometryBuilder) Descriptor() *Descriptor {
	return b.desc
}

type Descriptor struct {
	name         string
	info         *typeInfo
	description  string
	exampleValue string
	repeated     bool
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
			Description:  d.description,
			ExampleValue: d.exampleValue,
		}.Build(),
		Queryable: false,
	}.Build()
}
