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
	// These generated STAC, datasets/v1, Duration, and Timestamp imports register
	// supported dataset field type descriptors in protoregistry.GlobalFiles.
	_ "github.com/tilebox/tilebox-go/protogen/datasets/stac/v1"
	datasetsv1 "github.com/tilebox/tilebox-go/protogen/datasets/v1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
	_ "google.golang.org/protobuf/types/known/durationpb"
	_ "google.golang.org/protobuf/types/known/timestamppb"
)

// DatapointDescriptor contains the resolved protobuf descriptor needed to decode raw datapoint bytes.
type DatapointDescriptor struct {
	// MessageDescriptor is the protobuf message descriptor for datapoints in a dataset.
	MessageDescriptor protoreflect.MessageDescriptor
	resolver          *dynamicpb.Types
	enumValues        map[protoreflect.FieldNumber]map[protoreflect.EnumNumber]string
}

// NewDatapointDescriptor creates a reusable datapoint descriptor from a loaded dataset.
func NewDatapointDescriptor(dataset *Dataset) (*DatapointDescriptor, error) {
	if dataset == nil || dataset.Type == nil || dataset.Type.GetDescriptorSet() == nil || len(dataset.Type.GetDescriptorSet().GetFile()) == 0 {
		return nil, errors.New("dataset does not include a protobuf descriptor")
	}

	descriptorSet := descriptorSetWithDependencies(dataset.Type.GetDescriptorSet())
	files, err := protodesc.NewFiles(descriptorSet)
	if err != nil {
		return nil, fmt.Errorf("failed to build dataset protobuf descriptors: %w", err)
	}

	messageFullName, err := datasetMessageFullName(dataset.Type)
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

	enumValues, err := customEnumValues(messageDescriptor, dataset.Type.GetFieldAnnotations())
	if err != nil {
		return nil, err
	}

	return &DatapointDescriptor{
		MessageDescriptor: messageDescriptor,
		resolver:          dynamicpb.NewTypes(files),
		enumValues:        enumValues,
	}, nil
}

// EnumNumber resolves a custom enum's public value to its protobuf number for ingestion.
// Field names are protobuf field names. Matching is case-sensitive; generated symbols
// and the reserved zero member are not public values. Repeated enum fields use the same lookup.
func (d *DatapointDescriptor) EnumNumber(fieldName, value string) (protoreflect.EnumNumber, error) {
	if d == nil || d.MessageDescriptor == nil {
		return 0, errors.New("datapoint descriptor is required")
	}
	field := d.MessageDescriptor.Fields().ByName(protoreflect.Name(fieldName))
	if field == nil || d.enumValues[field.Number()] == nil {
		return 0, fmt.Errorf("field %q is not a custom enum", fieldName)
	}
	for number, publicValue := range d.enumValues[field.Number()] {
		if publicValue == value {
			return number, nil
		}
	}
	return 0, fmt.Errorf("field %q: unknown custom enum value %q", fieldName, value)
}

func customEnumValues(message protoreflect.MessageDescriptor, annotations []*datasetsv1.FieldAnnotation) (map[protoreflect.FieldNumber]map[protoreflect.EnumNumber]string, error) {
	result := make(map[protoreflect.FieldNumber]map[protoreflect.EnumNumber]string)
	for i, annotation := range annotations {
		values := annotation.GetEnumValues()
		if len(values) == 0 {
			continue
		}
		if len(annotations) != message.Fields().Len() {
			return nil, errors.New("custom enum annotations must align with datapoint fields")
		}
		field := message.Fields().Get(i)
		if field.Kind() != protoreflect.EnumKind {
			return nil, fmt.Errorf("field %q has custom enum values but is not an enum", field.Name())
		}
		members := field.Enum().Values()
		mapping := make(map[protoreflect.EnumNumber]string, len(values))
		seen := make(map[string]bool, len(values))
		index := 0
		for j := range members.Len() {
			number := members.Get(j).Number()
			if number == 0 {
				continue
			}
			if index >= len(values) {
				return nil, fmt.Errorf("field %q: enum annotations do not match descriptor members", field.Name())
			}
			value := values[index].GetValue()
			if value == "" || seen[value] {
				return nil, fmt.Errorf("field %q: custom enum values must be nonempty and unique", field.Name())
			}
			mapping[number] = value
			seen[value] = true
			index++
		}
		if index != len(values) {
			return nil, fmt.Errorf("field %q: enum annotations do not match descriptor members", field.Name())
		}
		result[field.Number()] = mapping
	}
	return result, nil
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
// Custom enums become their public string values. Missing fields are omitted;
// custom enum numbers without a public value, including zero, return an error.
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
	if err := convertSpecialTypes(descriptor, datapoint); err != nil {
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

// descriptorSetWithDependencies adds imports needed by the dataset descriptor.
// SDK-supported dependencies are registered in protoregistry.GlobalFiles by the
// generated protobuf imports above; application-registered dependencies can be resolved too.
func descriptorSetWithDependencies(datasetDescriptorSet *descriptorpb.FileDescriptorSet) *descriptorpb.FileDescriptorSet {
	descriptorSet := &descriptorpb.FileDescriptorSet{}
	seen := map[string]bool{}
	datasetFiles := make(map[string]*descriptorpb.FileDescriptorProto, len(datasetDescriptorSet.GetFile()))
	for _, file := range datasetDescriptorSet.GetFile() {
		datasetFiles[file.GetName()] = file
	}

	var addFile func(file *descriptorpb.FileDescriptorProto)
	addFile = func(file *descriptorpb.FileDescriptorProto) {
		if file == nil || seen[file.GetName()] {
			return
		}
		seen[file.GetName()] = true
		for _, dependency := range file.GetDependency() {
			if datasetFile := datasetFiles[dependency]; datasetFile != nil {
				addFile(datasetFile)
				continue
			}
			registeredFile, err := protoregistry.GlobalFiles.FindFileByPath(dependency)
			if err == nil {
				addFile(protodesc.ToFileDescriptorProto(registeredFile))
			}
		}
		descriptorSet.File = append(descriptorSet.File, file)
	}

	for _, file := range datasetDescriptorSet.GetFile() {
		addFile(file)
	}
	return descriptorSet
}

func datasetMessageFullName(annotatedType *datasetsv1.AnnotatedType) (protoreflect.FullName, error) {
	if typeURL := annotatedType.GetTypeUrl(); typeURL != "" {
		name := protoreflect.FullName(typeURL[strings.LastIndexByte(typeURL, '/')+1:])
		if !name.IsValid() {
			return "", fmt.Errorf("invalid dataset type URL %q", typeURL)
		}
		return name, nil
	}
	// Older dataset descriptors may omit type_url.
	for _, file := range annotatedType.GetDescriptorSet().GetFile() {
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

func convertSpecialTypes(descriptor *DatapointDescriptor, datapoint map[string]any) error {
	fields := descriptor.MessageDescriptor.Fields()
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

		converted, err := convertSpecialType(field, value, descriptor.enumValues[field.Number()])
		if err != nil {
			return fmt.Errorf("failed to convert field %q: %w", name, err)
		}
		datapoint[name] = converted
	}
	return nil
}

func convertSpecialType(field protoreflect.FieldDescriptor, value any, enumValues map[protoreflect.EnumNumber]string) (any, error) {
	if field.IsList() {
		values, ok := value.([]any)
		if !ok {
			return nil, fmt.Errorf("expected repeated field to be an array, got %T", value)
		}

		converted := make([]any, len(values))
		for i, item := range values {
			convertedItem, err := convertSpecialScalar(field, item, enumValues)
			if err != nil {
				return nil, fmt.Errorf("failed to convert item %d: %w", i, err)
			}
			converted[i] = convertedItem
		}
		return converted, nil
	}

	return convertSpecialScalar(field, value, enumValues)
}

func convertSpecialScalar(field protoreflect.FieldDescriptor, value any, enumValues map[protoreflect.EnumNumber]string) (any, error) {
	if field.Kind() == protoreflect.EnumKind {
		return convertEnum(field.Enum(), value, enumValues)
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

func convertEnum(descriptor protoreflect.EnumDescriptor, value any, enumValues map[protoreflect.EnumNumber]string) (any, error) {
	name, ok := value.(string)
	if !ok {
		return nil, fmt.Errorf("expected enum to be a string, got %T", value)
	}
	if enumValues != nil {
		member := descriptor.Values().ByName(protoreflect.Name(name))
		if member != nil {
			if publicValue, found := enumValues[member.Number()]; found {
				return publicValue, nil
			}
		}
		return nil, fmt.Errorf("enum member %q has no public value", name)
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
