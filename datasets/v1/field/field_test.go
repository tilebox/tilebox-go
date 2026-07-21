package field

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	stacv1 "github.com/tilebox/tilebox-go/protogen/datasets/stac/v1"
	datasetsv1 "github.com/tilebox/tilebox-go/protogen/datasets/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

func Test_Descriptor(t *testing.T) {
	tests := []struct {
		name       string
		descriptor *Descriptor
		want       *Descriptor
	}{
		{
			name:       "uuid",
			descriptor: UUID("id").Descriptor(),
			want: &Descriptor{
				name: "id",
				info: &typeInfo{
					Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE,
					TypeName: proto.String(".datasets.v1.UUID"),
				},
				description:  "",
				exampleValue: "",
				repeated:     false,
			},
		},
		{
			name:       "string with description",
			descriptor: String("test").Description("my description").Descriptor(),
			want: &Descriptor{
				name: "test",
				info: &typeInfo{
					Type:     descriptorpb.FieldDescriptorProto_TYPE_STRING,
					TypeName: nil,
				},
				description:  "my description",
				exampleValue: "",
				repeated:     false,
			},
		},
		{
			name:       "string with example value",
			descriptor: String("test").ExampleValue("my example value").Descriptor(),
			want: &Descriptor{
				name: "test",
				info: &typeInfo{
					Type:     descriptorpb.FieldDescriptorProto_TYPE_STRING,
					TypeName: nil,
				},
				description:  "",
				exampleValue: "my example value",
				repeated:     false,
			},
		},
		{
			name:       "string repeated",
			descriptor: String("test").Repeated().Descriptor(),
			want: &Descriptor{
				name: "test",
				info: &typeInfo{
					Type:     descriptorpb.FieldDescriptorProto_TYPE_STRING,
					TypeName: nil,
				},
				description:  "",
				exampleValue: "",
				repeated:     true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.descriptor)
		})
	}
}

func TestDescriptor_ToProtoWithAnnotations(t *testing.T) {
	field := String("stac_id").
		Description("Source STAC item ID").
		ExampleValue("S2A_001").
		SourceJSONPointer("/id").
		Queryable().
		JSONSchemaRef("https://schemas.stacspec.org/v1.1.0/item-spec/json-schema/item.json#/id").
		Roles(RolePrimaryTitle).
		ToProto()

	annotation := field.GetAnnotation()
	require.NotNil(t, annotation)
	assert.Equal(t, "Source STAC item ID", annotation.GetDescription())
	assert.Equal(t, "S2A_001", annotation.GetExampleValue())
	assert.True(t, annotation.HasSourceJsonPointer())
	assert.Equal(t, "/id", annotation.GetSourceJsonPointer())
	assert.True(t, annotation.GetQueryable())
	assert.True(t, annotation.HasJsonSchemaRef())
	assert.Equal(t, "https://schemas.stacspec.org/v1.1.0/item-spec/json-schema/item.json#/id", annotation.GetJsonSchemaRef())
	assert.Equal(t, []datasetsv1.FieldRole{datasetsv1.FieldRole_FIELD_ROLE_PRIMARY_TITLE}, annotation.GetRoles())
}

func TestDescriptor_ToProtoWithSTACTypes(t *testing.T) {
	tests := []struct {
		name         string
		descriptor   *Descriptor
		wantType     descriptorpb.FieldDescriptorProto_Type
		wantTypeName string
	}{
		{
			name:         "message",
			descriptor:   Message("assets", &stacv1.Assets{}),
			wantType:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE,
			wantTypeName: ".datasets.stac.v1.Assets",
		},
		{
			name:         "enum",
			descriptor:   Enum("orbit_state", stacv1.SatelliteOrbitState_SATELLITE_ORBIT_STATE_ASCENDING),
			wantType:     descriptorpb.FieldDescriptorProto_TYPE_ENUM,
			wantTypeName: ".datasets.stac.v1.SatelliteOrbitState",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := tt.descriptor.ToProto().GetDescriptor()
			require.NotNil(t, field)
			assert.Equal(t, tt.wantType, field.GetType())
			assert.Equal(t, tt.wantTypeName, field.GetTypeName())
		})
	}
}
