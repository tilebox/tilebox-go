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
