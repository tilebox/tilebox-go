package datasets

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/paulmach/orb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	stacv1 "github.com/tilebox/tilebox-go/protogen/datasets/stac/v1"
	datasetsv1 "github.com/tilebox/tilebox-go/protogen/datasets/v1"
	examplesv1 "github.com/tilebox/tilebox-go/protogen/examples/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestNewDatapointDescriptor(t *testing.T) {
	descriptor, err := NewDatapointDescriptor(exampleDataset())

	require.NoError(t, err)
	require.NotNil(t, descriptor)
	assert.Equal(t, "tilebox.v1.Sentinel2Msi", string(descriptor.MessageDescriptor.FullName()))
}

func TestNewDatapointDescriptorResolvesSTACImports(t *testing.T) {
	descriptor, err := NewDatapointDescriptor(stacDataset())

	require.NoError(t, err)
	require.NotNil(t, descriptor)
	assets := descriptor.MessageDescriptor.Fields().ByName("assets")
	require.NotNil(t, assets)
	assert.Equal(t, protoreflect.FullName("datasets.stac.v1.Assets"), assets.Message().FullName())
}

func TestUnmarshalDatapoint(t *testing.T) {
	descriptor, err := NewDatapointDescriptor(exampleDataset())
	require.NoError(t, err)

	datapointID := uuid.MustParse("01941f29-c650-202f-6495-c71dd2118fb1")
	geometry := orb.Point{16, 48}
	timestamp := time.Date(2025, time.January, 1, 0, 0, 19, 24_000_000, time.UTC)
	datapoint := examplesv1.Sentinel2Msi_builder{
		Id:              datasetsv1.NewUUID(datapointID),
		Time:            timestamppb.New(timestamp),
		Geometry:        datasetsv1.NewGeometry(geometry),
		GranuleName:     new("S2B_MSIL1C_20250101T000019_N0511_R073_T57QWV_20250101T010340.SAFE"),
		ProcessingLevel: new(datasetsv1.ProcessingLevel_PROCESSING_LEVEL_L1C),
		FlightDirection: new(datasetsv1.FlightDirection_FLIGHT_DIRECTION_ASCENDING),
		AcquisitionMode: new(datasetsv1.AcquisitionMode_ACQUISITION_MODE_NOBS),
	}.Build()
	rawDatapoint, err := proto.Marshal(datapoint)
	require.NoError(t, err)

	got, err := UnmarshalDatapoint(descriptor, rawDatapoint)

	require.NoError(t, err)
	assert.Equal(t, timestamp, got["time"])
	assert.Equal(t, geometry, got["geometry"])
	assert.Equal(t, "S2B_MSIL1C_20250101T000019_N0511_R073_T57QWV_20250101T010340.SAFE", got["granule_name"])
	assert.Equal(t, "L1C", got["processing_level"])
	assert.Equal(t, "ASCENDING", got["flight_direction"])
	assert.Equal(t, "NOBS", got["acquisition_mode"])
	assert.Equal(t, datapointID, got["id"])
	assert.NotContains(t, got, "cloud_cover")
	assert.NotContains(t, got, "file_size")
	assert.NotContains(t, got, "updated")
}

func TestUnmarshalDatapointConvertsDuration(t *testing.T) {
	descriptor, err := NewDatapointDescriptor(durationDataset())
	require.NoError(t, err)

	duration := 1500 * time.Millisecond
	message := dynamicpb.NewMessage(descriptor.MessageDescriptor)
	field := descriptor.MessageDescriptor.Fields().ByName("elapsed")
	message.Set(field, protoreflect.ValueOfMessage(durationpb.New(duration).ProtoReflect()))
	rawDatapoint, err := proto.Marshal(message)
	require.NoError(t, err)

	got, err := UnmarshalDatapoint(descriptor, rawDatapoint)

	require.NoError(t, err)
	assert.Equal(t, duration, got["elapsed"])
}

func TestDatapointDecoderUnmarshalValidatesDescriptor(t *testing.T) {
	_, err := DatapointDecoder{}.Unmarshal(nil, nil)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "datapoint descriptor is required")
}

func TestCustomEnumDatapoints(t *testing.T) {
	descriptor, err := NewDatapointDescriptor(customEnumDataset())
	require.NoError(t, err)
	assert.Equal(t, protoreflect.FullName("test.Datapoint"), descriptor.MessageDescriptor.FullName())

	for _, tt := range []struct {
		value  string
		number protoreflect.EnumNumber
		raw    []byte
	}{
		{value: "L1", number: 17, raw: []byte{0x38, 17}},
		{value: "L2-SP", number: 4, raw: []byte{0x38, 4}},
	} {
		t.Run(tt.value, func(t *testing.T) {
			number, err := descriptor.EnumNumber("processing_level", tt.value)
			require.NoError(t, err)
			assert.Equal(t, tt.number, number)
			message := dynamicpb.NewMessage(descriptor.MessageDescriptor)
			message.Set(descriptor.MessageDescriptor.Fields().ByName("processing_level"), protoreflect.ValueOfEnum(number))
			raw, err := proto.Marshal(message)
			require.NoError(t, err)
			// Field 7 is a varint. Verify encoding independently of the decoder.
			assert.Equal(t, tt.raw, raw)
			got, err := UnmarshalDatapoint(descriptor, raw)
			require.NoError(t, err)
			assert.Equal(t, map[string]any{"processing_level": tt.value}, got)
		})
	}

	for _, tt := range []struct {
		name    string
		raw     []byte
		want    map[string]any
		wantErr string
	}{
		{name: "absent", want: map[string]any{}},
		{name: "repeated", raw: []byte{0x5a, 0x03, 4, 17, 4}, want: map[string]any{"levels": []any{"L2-SP", "L1", "L2-SP"}}},
		{name: "zero", raw: []byte{0x38, 0}, wantErr: "has no public value"},
		{name: "unknown number", raw: []byte{0x38, 99}, wantErr: "failed to convert field"},
		{name: "repeated zero", raw: []byte{0x5a, 0x02, 17, 0}, wantErr: "failed to convert item 1"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnmarshalDatapoint(descriptor, tt.raw)
			if tt.wantErr != "" {
				require.ErrorContains(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}

	for _, value := range []string{"l1", "L2", "GENERATED_FIRST", "GENERATED_ZERO", ""} {
		_, err := descriptor.EnumNumber("processing_level", value)
		require.ErrorContains(t, err, "unknown custom enum value")
	}
	for _, fieldName := range []string{"unknown", "title"} {
		_, err := descriptor.EnumNumber(fieldName, "L1")
		require.ErrorContains(t, err, "is not a custom enum")
	}
	number, err := descriptor.EnumNumber("levels", "L2-SP")
	require.NoError(t, err)
	assert.Equal(t, protoreflect.EnumNumber(4), number)
}

func TestNewDatapointDescriptorTypeURL(t *testing.T) {
	for _, tt := range []struct {
		typeURL string
		wantErr string
	}{
		{typeURL: "test.Datapoint"},
		{typeURL: "type.googleapis.com/test.Datapoint"},
		{typeURL: "example.com/test.Missing", wantErr: "failed to find dataset protobuf message descriptor"},
		{typeURL: "example.com/test.Level", wantErr: "is not a message descriptor"},
		{typeURL: "example.com/", wantErr: "invalid dataset type URL"},
	} {
		t.Run(tt.typeURL, func(t *testing.T) {
			dataset := customEnumDataset()
			dataset.Type.SetTypeUrl(tt.typeURL)
			descriptor, err := NewDatapointDescriptor(dataset)
			if tt.wantErr != "" {
				require.ErrorContains(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, protoreflect.FullName("test.Datapoint"), descriptor.MessageDescriptor.FullName())
		})
	}
}

func TestNewDatapointDescriptorInvalidEnumAnnotations(t *testing.T) {
	for _, tt := range []struct {
		name   string
		modify func(*datasetsv1.AnnotatedType)
	}{
		{name: "missing placeholder", modify: func(schema *datasetsv1.AnnotatedType) {
			schema.SetFieldAnnotations(schema.GetFieldAnnotations()[1:])
		}},
		{name: "wrong field type", modify: func(schema *datasetsv1.AnnotatedType) {
			schema.GetFieldAnnotations()[0].SetEnumValues(schema.GetFieldAnnotations()[1].GetEnumValues())
		}},
		{name: "too few values", modify: func(schema *datasetsv1.AnnotatedType) {
			annotation := schema.GetFieldAnnotations()[1]
			annotation.SetEnumValues(annotation.GetEnumValues()[:1])
		}},
		{name: "too many values", modify: func(schema *datasetsv1.AnnotatedType) {
			annotation := schema.GetFieldAnnotations()[1]
			annotation.SetEnumValues(append(annotation.GetEnumValues(), datasetsv1.EnumValue_builder{Value: "L3"}.Build()))
		}},
		{name: "duplicate value", modify: func(schema *datasetsv1.AnnotatedType) {
			schema.GetFieldAnnotations()[1].GetEnumValues()[1].SetValue("L1")
		}},
		{name: "empty value", modify: func(schema *datasetsv1.AnnotatedType) {
			schema.GetFieldAnnotations()[1].GetEnumValues()[0].SetValue("")
		}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			dataset := customEnumDataset()
			tt.modify(dataset.Type)
			_, err := NewDatapointDescriptor(dataset)
			require.Error(t, err)
		})
	}
}

func customEnumDataset() *Dataset {
	annotation := datasetsv1.FieldAnnotation_builder{EnumValues: []*datasetsv1.EnumValue{
		datasetsv1.EnumValue_builder{Value: "L1"}.Build(),
		datasetsv1.EnumValue_builder{Value: "L2-SP"}.Build(),
	}}.Build()
	return &Dataset{Type: datasetsv1.AnnotatedType_builder{
		TypeUrl: "type.googleapis.com/test.Datapoint",
		FieldAnnotations: []*datasetsv1.FieldAnnotation{
			{}, annotation, proto.CloneOf(annotation),
		},
		DescriptorSet: &descriptorpb.FileDescriptorSet{File: []*descriptorpb.FileDescriptorProto{
			{
				Name: new("enums.proto"), Package: new("test"), Syntax: new("editions"), Edition: descriptorpb.Edition_EDITION_2023.Enum(),
				MessageType: []*descriptorpb.DescriptorProto{{Name: new("Unrelated")}},
				EnumType: []*descriptorpb.EnumDescriptorProto{{
					Name: new("Level"),
					Value: []*descriptorpb.EnumValueDescriptorProto{
						{Name: new("GENERATED_ZERO"), Number: new(int32(0))},
						{Name: new("GENERATED_FIRST"), Number: new(int32(17))},
						{Name: new("GENERATED_SECOND"), Number: new(int32(4))},
					},
				}},
			},
			{
				Name: new("datapoint.proto"), Package: new("test"), Syntax: new("editions"), Edition: descriptorpb.Edition_EDITION_2023.Enum(),
				Dependency: []string{"enums.proto"},
				MessageType: []*descriptorpb.DescriptorProto{
					{Name: new("AlsoUnrelated")},
					{Name: new("Datapoint"), Field: []*descriptorpb.FieldDescriptorProto{
						// Annotations follow declaration order, not field-number order.
						{Name: new("title"), Number: new(int32(20)), Type: descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum()},
						{Name: new("processing_level"), Number: new(int32(7)), Type: descriptorpb.FieldDescriptorProto_TYPE_ENUM.Enum(), TypeName: new(".test.Level")},
						{Name: new("levels"), Number: new(int32(11)), Type: descriptorpb.FieldDescriptorProto_TYPE_ENUM.Enum(), TypeName: new(".test.Level"), Label: descriptorpb.FieldDescriptorProto_LABEL_REPEATED.Enum()},
					}},
				},
			},
		}},
	}.Build()}
}

func exampleDataset() *Dataset {
	return &Dataset{
		Type: datasetsv1.AnnotatedType_builder{
			DescriptorSet: &descriptorpb.FileDescriptorSet{File: []*descriptorpb.FileDescriptorProto{
				protodesc.ToFileDescriptorProto(examplesv1.File_tilebox_v1_Sentinel2Msi_proto),
			}},
		}.Build(),
	}
}

func durationDataset() *Dataset {
	return &Dataset{
		Type: datasetsv1.AnnotatedType_builder{
			DescriptorSet: &descriptorpb.FileDescriptorSet{File: []*descriptorpb.FileDescriptorProto{
				{
					Name:       new("tilebox/v1/DurationDatapoint.proto"),
					Package:    new("tilebox.v1"),
					Dependency: []string{"google/protobuf/duration.proto"},
					MessageType: []*descriptorpb.DescriptorProto{
						{
							Name: new("DurationDatapoint"),
							Field: []*descriptorpb.FieldDescriptorProto{
								{
									Name:     new("elapsed"),
									Number:   proto.Int32(1),
									Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
									Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
									TypeName: new(".google.protobuf.Duration"),
								},
							},
						},
					},
				},
			}},
		}.Build(),
	}
}

func stacDataset() *Dataset {
	return &Dataset{
		Type: datasetsv1.AnnotatedType_builder{
			DescriptorSet: &descriptorpb.FileDescriptorSet{File: []*descriptorpb.FileDescriptorProto{
				{
					Name:       new("tilebox/v1/STACDatapoint.proto"),
					Package:    new("tilebox.v1"),
					Dependency: []string{stacv1.File_datasets_stac_v1_asset_proto.Path()},
					MessageType: []*descriptorpb.DescriptorProto{
						{
							Name: new("STACDatapoint"),
							Field: []*descriptorpb.FieldDescriptorProto{
								{
									Name:     new("assets"),
									Number:   proto.Int32(1),
									Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
									Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
									TypeName: new(".datasets.stac.v1.Assets"),
								},
							},
						},
					},
				},
			}},
		}.Build(),
	}
}
