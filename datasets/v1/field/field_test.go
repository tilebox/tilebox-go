package field

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
			descriptor: UUID("id"),
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
			descriptor: String("test").Description("my description"),
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
			descriptor: String("test").ExampleValue("my example value"),
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
			descriptor: String("test").Repeated(),
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
