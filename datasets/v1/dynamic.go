package datasets

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
	datasetsv1 "github.com/tilebox/tilebox-go/protogen/datasets/v1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// DatapointDescriptor contains the resolved protobuf descriptor needed to decode raw datapoint bytes.
type DatapointDescriptor struct {
	// MessageDescriptor is the protobuf message descriptor for datapoints in a dataset.
	MessageDescriptor protoreflect.MessageDescriptor
	resolver          *dynamicpb.Types
}

// NewDatapointDescriptor creates a reusable datapoint descriptor from a loaded dataset.
func NewDatapointDescriptor(dataset *Dataset) (*DatapointDescriptor, error) {
	if dataset == nil || dataset.Type == nil || dataset.Type.GetDescriptorSet() == nil || len(dataset.Type.GetDescriptorSet().GetFile()) == 0 {
		return nil, errors.New("dataset does not include a protobuf descriptor")
	}

	descriptorSet := descriptorSetWithWellKnownTypes(dataset.Type.GetDescriptorSet())
	files, err := protodesc.NewFiles(descriptorSet)
	if err != nil {
		return nil, fmt.Errorf("failed to build dataset protobuf descriptors: %w", err)
	}

	messageFullName, err := datasetMessageFullName(dataset.Type.GetDescriptorSet())
	if err != nil {
		return nil, err
	}
	descriptor, err := files.FindDescriptorByName(messageFullName)
	if err != nil {
		return nil, fmt.Errorf("failed to find dataset protobuf message descriptor %q: %w", messageFullName, err)
	}
	messageDescriptor, ok := descriptor.(protoreflect.MessageDescriptor)
	if !ok {
		return nil, fmt.Errorf("dataset protobuf descriptor %q is not a message descriptor", messageFullName)
	}

	return &DatapointDescriptor{
		MessageDescriptor: messageDescriptor,
		resolver:          dynamicpb.NewTypes(files),
	}, nil
}

// DatapointDecoder decodes raw protobuf datapoints into JSON-like maps.
//
// The options mirror protojson.UnmarshalOptions.
type DatapointDecoder struct {
	// If AllowPartial is set, input for messages that will result in missing required fields will not return an error.
	AllowPartial bool
	// If DiscardUnknown is set, unknown fields are ignored.
	DiscardUnknown bool
	// Resolver is used for looking up message and extension types. If nil, the resolver from the DatapointDescriptor is used.
	Resolver interface {
		protoregistry.MessageTypeResolver
		protoregistry.ExtensionTypeResolver
	}
	// RecursionLimit limits how deeply messages may be nested. If zero, a default limit is applied.
	RecursionLimit int
}

// UnmarshalDatapoint decodes raw protobuf datapoint bytes into a JSON-like map.
func UnmarshalDatapoint(descriptor *DatapointDescriptor, data []byte) (map[string]any, error) {
	return DatapointDecoder{}.Unmarshal(descriptor, data)
}

// Unmarshal decodes raw protobuf datapoint bytes into a JSON-like map.
func (d DatapointDecoder) Unmarshal(descriptor *DatapointDescriptor, data []byte) (map[string]any, error) {
	if descriptor == nil || descriptor.MessageDescriptor == nil {
		return nil, errors.New("datapoint descriptor is required")
	}

	resolver := d.resolver(descriptor)
	message := dynamicpb.NewMessage(descriptor.MessageDescriptor)
	err := proto.UnmarshalOptions{
		AllowPartial:   d.AllowPartial,
		DiscardUnknown: d.DiscardUnknown,
		Resolver:       resolver,
		RecursionLimit: d.RecursionLimit,
	}.Unmarshal(data, message)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal datapoint: %w", err)
	}

	jsonData, err := protojson.MarshalOptions{
		AllowPartial:  d.AllowPartial,
		Resolver:      resolver,
		UseProtoNames: true,
	}.Marshal(message)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal datapoint as JSON: %w", err)
	}

	var datapoint map[string]any
	if err := json.Unmarshal(jsonData, &datapoint); err != nil {
		return nil, fmt.Errorf("failed to unmarshal datapoint JSON: %w", err)
	}
	if err := convertSpecialTypes(descriptor.MessageDescriptor, datapoint); err != nil {
		return nil, err
	}
	return datapoint, nil
}

func (d DatapointDecoder) resolver(descriptor *DatapointDescriptor) interface {
	protoregistry.MessageTypeResolver
	protoregistry.ExtensionTypeResolver
} {
	if d.Resolver != nil {
		return d.Resolver
	}
	return descriptor.resolver
}

func descriptorSetWithWellKnownTypes(datasetDescriptorSet *descriptorpb.FileDescriptorSet) *descriptorpb.FileDescriptorSet {
	descriptorSet := &descriptorpb.FileDescriptorSet{}
	seen := map[string]bool{}
	addFile := func(file *descriptorpb.FileDescriptorProto) {
		if file == nil || seen[file.GetName()] {
			return
		}
		seen[file.GetName()] = true
		descriptorSet.File = append(descriptorSet.File, file)
	}

	addFile(protodesc.ToFileDescriptorProto(durationpb.File_google_protobuf_duration_proto))
	addFile(protodesc.ToFileDescriptorProto(timestamppb.File_google_protobuf_timestamp_proto))
	addFile(protodesc.ToFileDescriptorProto(datasetsv1.File_datasets_v1_well_known_types_proto))
	for _, file := range datasetDescriptorSet.GetFile() {
		addFile(file)
	}
	return descriptorSet
}

func datasetMessageFullName(descriptorSet *descriptorpb.FileDescriptorSet) (protoreflect.FullName, error) {
	for _, file := range descriptorSet.GetFile() {
		if len(file.GetMessageType()) == 0 {
			continue
		}

		messageName := file.GetMessageType()[0].GetName()
		if messageName == "" {
			return "", errors.New("dataset protobuf descriptor includes a message type without a name")
		}
		if file.GetPackage() == "" {
			return protoreflect.FullName(messageName), nil
		}
		return protoreflect.FullName(strings.Join([]string{file.GetPackage(), messageName}, ".")), nil
	}
	return "", errors.New("dataset protobuf descriptor does not include a message type")
}

func convertSpecialTypes(descriptor protoreflect.MessageDescriptor, datapoint map[string]any) error {
	fields := descriptor.Fields()
	for i := range fields.Len() {
		field := fields.Get(i)
		if field.IsMap() || (field.Kind() != protoreflect.MessageKind && field.Kind() != protoreflect.EnumKind) {
			continue
		}

		name := string(field.Name())
		value, ok := datapoint[name]
		if !ok {
			continue
		}

		converted, err := convertSpecialType(field, value)
		if err != nil {
			return fmt.Errorf("failed to convert field %q: %w", name, err)
		}
		datapoint[name] = converted
	}
	return nil
}

func convertSpecialType(field protoreflect.FieldDescriptor, value any) (any, error) {
	if field.IsList() {
		values, ok := value.([]any)
		if !ok {
			return nil, fmt.Errorf("expected repeated field to be an array, got %T", value)
		}

		converted := make([]any, len(values))
		for i, item := range values {
			convertedItem, err := convertSpecialScalar(field, item)
			if err != nil {
				return nil, fmt.Errorf("failed to convert item %d: %w", i, err)
			}
			converted[i] = convertedItem
		}
		return converted, nil
	}

	return convertSpecialScalar(field, value)
}

func convertSpecialScalar(field protoreflect.FieldDescriptor, value any) (any, error) {
	if field.Kind() == protoreflect.EnumKind {
		return convertEnum(field.Enum(), value)
	}

	switch field.Message().FullName() {
	case "datasets.v1.Geometry":
		return convertGeometry(value)
	case "datasets.v1.UUID":
		return convertUUID(value)
	case "google.protobuf.Duration":
		return convertDuration(value)
	case "google.protobuf.Timestamp":
		return convertTimestamp(value)
	default:
		return value, nil
	}
}

func convertEnum(descriptor protoreflect.EnumDescriptor, value any) (any, error) {
	name, ok := value.(string)
	if !ok {
		return nil, fmt.Errorf("expected enum to be a string, got %T", value)
	}
	prefix := screamingSnakeCase(string(descriptor.Name())) + "_"
	return strings.TrimPrefix(name, prefix), nil
}

func screamingSnakeCase(value string) string {
	var builder strings.Builder
	for i, r := range value {
		if i > 0 && unicode.IsUpper(r) {
			builder.WriteByte('_')
		}
		builder.WriteRune(unicode.ToUpper(r))
	}
	return builder.String()
}

func convertGeometry(value any) (any, error) {
	message, ok := value.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("expected Geometry to be an object, got %T", value)
	}
	wkb, err := bytesField(message, "wkb")
	if err != nil {
		return nil, err
	}
	return datasetsv1.Geometry_builder{Wkb: wkb}.Build().AsGeometry(), nil
}

func convertUUID(value any) (any, error) {
	message, ok := value.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("expected UUID to be an object, got %T", value)
	}
	data, err := bytesField(message, "uuid")
	if err != nil {
		return nil, err
	}
	id, err := uuid.FromBytes(data)
	if err != nil {
		return nil, fmt.Errorf("invalid UUID: %w", err)
	}
	return id, nil
}

func convertDuration(value any) (any, error) {
	duration, ok := value.(string)
	if !ok {
		return nil, fmt.Errorf("expected Duration to be a string, got %T", value)
	}
	parsed, err := time.ParseDuration(duration)
	if err != nil {
		return nil, fmt.Errorf("invalid Duration: %w", err)
	}
	return parsed, nil
}

func convertTimestamp(value any) (any, error) {
	timestamp, ok := value.(string)
	if !ok {
		return nil, fmt.Errorf("expected Timestamp to be a string, got %T", value)
	}
	parsed, err := time.Parse(time.RFC3339Nano, timestamp)
	if err != nil {
		return nil, fmt.Errorf("invalid Timestamp: %w", err)
	}
	return parsed, nil
}

func bytesField(message map[string]any, name string) ([]byte, error) {
	value, ok := message[name]
	if !ok {
		return nil, fmt.Errorf("missing %q", name)
	}
	encoded, ok := value.(string)
	if !ok {
		return nil, fmt.Errorf("expected %q to be a base64 string, got %T", name, value)
	}
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("invalid %q: %w", name, err)
	}
	return data, nil
}
