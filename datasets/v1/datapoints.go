package datasets

import (
	"context"
	"errors"
	"fmt"
	"iter"
	"reflect"
	"strings"
	"time"

	"github.com/google/uuid"
	datasetsv1 "github.com/tilebox/tilebox-go/protogen/go/datasets/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type DatapointsClient interface {
	Load(ctx context.Context, collectionID uuid.UUID, interval LoadInterval, options ...LoadOption) iter.Seq2[*RawDatapoint, error]
	LoadInto(ctx context.Context, collectionID uuid.UUID, interval LoadInterval, datapoints any, options ...LoadOption) error
	Ingest(ctx context.Context, collectionID uuid.UUID, data []*Datapoint, allowExisting bool) (*IngestResponse, error)
	Delete(ctx context.Context, collectionID uuid.UUID, data []*Datapoint) (*DeleteResponse, error)
	DeleteIDs(ctx context.Context, collectionID uuid.UUID, datapointIDs []uuid.UUID) (*DeleteResponse, error)
}

var _ DatapointsClient = &datapointClient{}

type datapointClient struct {
	dataIngestionService DataIngestionService
	dataAccessService    DataAccessService
}

// loadConfig contains the configuration for a Load request.
type loadConfig struct {
	skipData bool
	skipMeta bool
}

// LoadOption is an interface for configuring a Load request.
type LoadOption func(*loadConfig)

// WithSkipData skips the data when loading datapoints.
// It is an optional flag for omitting the actual datapoint data from the response.
// If set, no datapoint data will be returned.
//
// Defaults to false.
func WithSkipData() LoadOption {
	return func(cfg *loadConfig) {
		cfg.skipData = true
	}
}

// WithSkipMeta skips the metadata when loading datapoints.
// It is an optional flag for omitting the metadata from the response.
// If set, no metadata will be returned.
//
// Defaults to false.
func WithSkipMeta() LoadOption {
	return func(cfg *loadConfig) {
		cfg.skipMeta = true
	}
}

func newLoadConfig(options []LoadOption) *loadConfig {
	cfg := &loadConfig{
		skipData: false,
		skipMeta: false,
	}
	for _, option := range options {
		option(cfg)
	}

	return cfg
}

// Load loads datapoints from a collection.
//
// interval specifies the time or data point interval for which data should be loaded.
//
// WithSkipData and WithSkipMeta can be used to skip the data or metadata when loading datapoints.
// If both WithSkipData and WithSkipMeta are specified, the response will only consist of a list of datapoint IDs without any
// additional data or metadata.
//
// The datapoints are loaded in a lazy manner, and returned as a sequence of RawDatapoint.
// The output sequence can be transformed into typed Datapoint using CollectAs or As functions.
//
// Documentation: https://docs.tilebox.com/datasets/loading-data
func (d datapointClient) Load(ctx context.Context, collectionID uuid.UUID, interval LoadInterval, options ...LoadOption) iter.Seq2[*RawDatapoint, error] {
	cfg := newLoadConfig(options)

	return func(yield func(*RawDatapoint, error) bool) {
		var page *datasetsv1.Pagination // nil for the first request

		timeInterval := interval.ToProtoTimeInterval()
		datapointInterval := interval.ToProtoDatapointInterval()

		if timeInterval == nil && datapointInterval == nil {
			yield(nil, errors.New("time interval and datapoint interval cannot both be nil"))
			return
		}

		for {
			datapointsMessage, err := d.dataAccessService.GetDatasetForInterval(ctx, collectionID, timeInterval, datapointInterval, page, cfg.skipData, cfg.skipMeta)
			if err != nil {
				yield(nil, err)
				return
			}

			// if skipMeta is true, datapointsMessage.GetMeta() is not nil and contains the datapoint ids
			for i, dp := range datapointsMessage.GetMeta() {
				datapointID, err := uuid.Parse(dp.GetId())
				if err != nil {
					yield(nil, fmt.Errorf("failed to parse datapoint id from response: %w", err))
					return
				}

				meta := &DatapointMetadata{
					ID: datapointID,
				}
				var data []byte

				if !cfg.skipMeta {
					meta.EventTime = dp.GetEventTime().AsTime()
					meta.IngestionTime = dp.GetIngestionTime().AsTime()
				}

				if !cfg.skipData {
					data = datapointsMessage.GetData().GetValue()[i]
				}

				datapoint := &RawDatapoint{
					Meta: meta,
					Data: data,
				}
				if !yield(datapoint, nil) {
					return
				}
			}

			page = datapointsMessage.GetNextPage()
			if page == nil {
				break
			}
		}
	}
}

// LoadInto loads datapoints from a collection into a slice of datapoints.
// LoadInto is a convenience function for Load.
//
// Example usage:
// var datapoints []*tileboxdatasets.TypedDatapoint[*datasetsv1.CopernicusDataspaceGranule]
// err := client.Datapoints.LoadInto(ctx, collection.ID, loadInterval, &datapoints)
func (d datapointClient) LoadInto(ctx context.Context, collectionID uuid.UUID, interval LoadInterval, datapoints any, options ...LoadOption) error {
	rv := reflect.ValueOf(datapoints)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return fmt.Errorf("datapoints must be a pointer, got %v", reflect.TypeOf(datapoints))
	}
	slice := reflect.Indirect(rv)
	if slice.Kind() != reflect.Slice {
		return fmt.Errorf("datapoints must be a pointer to a slice, got %v", reflect.TypeOf(datapoints))
	}
	if slice.Type().Elem().Kind() != reflect.Pointer {
		return fmt.Errorf("datapoints must be a pointer to a slice of *TypedDatapoint, got %v", reflect.TypeOf(datapoints))
	}
	if !strings.HasPrefix(slice.Type().Elem().Elem().Name(), "TypedDatapoint[") {
		return fmt.Errorf("datapoints must be a pointer to a slice of *TypedDatapoint, got %v", reflect.TypeOf(datapoints))
	}

	typedDatapointType := slice.Type().Elem().Elem()
	dataField, ok := typedDatapointType.FieldByName("Data")
	if !ok {
		return errors.New("datapoints does not contain a field named Data")
	}
	datapointType := dataField.Type.Elem()

	rawDatapoints, err := Collect(d.Load(ctx, collectionID, interval, options...))
	if err != nil {
		return err
	}

	slice.Set(reflect.MakeSlice(slice.Type(), len(rawDatapoints), len(rawDatapoints)))

	for i, rawDatapoint := range rawDatapoints {
		typedDatapoint := reflect.New(typedDatapointType)
		reflect.Indirect(typedDatapoint).FieldByName("Meta").Set(reflect.ValueOf(rawDatapoint.Meta))
		reflect.Indirect(typedDatapoint).FieldByName("Data").Set(reflect.New(datapointType))

		err := proto.Unmarshal(rawDatapoint.Data, reflect.Indirect(typedDatapoint).FieldByName("Data").Interface().(proto.Message))
		if err != nil {
			return fmt.Errorf("failed to unmarshal datapoint data: %w", err)
		}

		slice.Index(i).Set(typedDatapoint)
	}
	return nil
}

// IngestResponse contains the response from the Ingest method.
type IngestResponse struct {
	// NumCreated is the number of datapoints that were created.
	NumCreated int64
	// NumExisting is the number of datapoints that were ignored because they already existed.
	NumExisting int64
	// DatapointIDs is the list of all the datapoints IDs in the same order as the datapoints in the request.
	DatapointIDs []uuid.UUID
}

// Ingest datapoints into a collection.
//
// data is a list of datapoints to ingest that should be created using Datapoints.
//
// allowExisting specifies whether to allow existing datapoints as part of the request. If true, datapoints that already
// exist will be ignored, and the number of such existing datapoints will be returned in the response. If false, any
// datapoints that already exist will result in an error. Setting this to true is useful for achieving idempotency (e.g.
// allowing re-ingestion of datapoints that have already been ingested in the past).
func (d datapointClient) Ingest(ctx context.Context, collectionID uuid.UUID, data []*Datapoint, allowExisting bool) (*IngestResponse, error) {
	rawDatapoints, err := toRawDatapoints(data)
	if err != nil {
		return nil, fmt.Errorf("failed to convert datapoints to raw datapoints: %w", err)
	}

	datapoints := &datasetsv1.Datapoints{
		Meta: make([]*datasetsv1.DatapointMetadata, len(data)),
		Data: &datasetsv1.RepeatedAny{
			Value: make([][]byte, len(data)),
		},
	}

	for i, datapoint := range rawDatapoints {
		datapoints.GetMeta()[i] = &datasetsv1.DatapointMetadata{
			EventTime: timestamppb.New(datapoint.Meta.EventTime),
		}
		datapoints.GetData().GetValue()[i] = datapoint.Data
	}

	response, err := d.dataIngestionService.IngestDatapoints(ctx, collectionID, datapoints, allowExisting)
	if err != nil {
		return nil, err
	}

	datapointIDs := make([]uuid.UUID, len(response.GetDatapointIds()))
	for i, id := range response.GetDatapointIds() {
		datapointID, err := protoToUUID(id)
		if err != nil {
			return nil, fmt.Errorf("failed to convert datapoint id from response: %w", err)
		}
		datapointIDs[i] = datapointID
	}

	return &IngestResponse{
		NumCreated:   response.GetNumCreated(),
		NumExisting:  response.GetNumExisting(),
		DatapointIDs: datapointIDs,
	}, nil
}

// DeleteResponse contains the response from the Delete method.
type DeleteResponse struct {
	// NumDeleted is the number of datapoints that were deleted.
	NumDeleted int64
}

// Delete datapoints from a collection.
//
// The datapoints are identified by their IDs.
func (d datapointClient) Delete(ctx context.Context, collectionID uuid.UUID, data []*Datapoint) (*DeleteResponse, error) {
	datapointIDs := make([]uuid.UUID, len(data))
	for i, datapoint := range data {
		datapointIDs[i] = datapoint.Meta.ID
	}

	return d.DeleteIDs(ctx, collectionID, datapointIDs)
}

// DeleteIDs deletes datapoints from a collection by their IDs.
func (d datapointClient) DeleteIDs(ctx context.Context, collectionID uuid.UUID, datapointIDs []uuid.UUID) (*DeleteResponse, error) {
	response, err := d.dataIngestionService.DeleteDatapoints(ctx, collectionID, datapointIDs)
	if err != nil {
		return nil, err
	}

	return &DeleteResponse{
		NumDeleted: response.GetNumDeleted(),
	}, nil
}

// DatapointMetadata contains metadata for a datapoint.
type DatapointMetadata struct {
	// ID is the unique identifier of the datapoint.
	ID uuid.UUID
	// EventTime is the time when the datapoint was created.
	EventTime time.Time
	// IngestionTime is the time when the datapoint was ingested into Tilebox.
	IngestionTime time.Time
}

type TypedDatapoint[T proto.Message] struct {
	// Meta contains the metadata of the datapoint.
	Meta *DatapointMetadata
	// Data contains the data of the datapoint.
	Data T
}

// Datapoint represents a datapoint in a collection.
// It contains the metadata and the data itself.
type Datapoint struct {
	// Meta contains the metadata of the datapoint.
	Meta *DatapointMetadata
	// Data contains the data of the datapoint.
	Data proto.Message
}

// RawDatapoint is an internal representation of a datapoint.
//
// It can be transformed into a Datapoint using CollectAs or As functions.
type RawDatapoint struct {
	// Meta contains the metadata of the datapoint.
	Meta *DatapointMetadata
	// Data contains the data of the datapoint in an internal raw format.
	Data []byte
}

// NewDatapoint creates a new datapoint with the given time and message.
func NewDatapoint(time time.Time, message proto.Message) *Datapoint {
	return &Datapoint{
		Meta: &DatapointMetadata{
			EventTime: time,
		},
		Data: message,
	}
}

// Datapoints converts a list of Datapoint to RawDatapoint.
//
// It is used to convert the data before ingesting it into a collection.
func toRawDatapoints(data []*Datapoint) ([]*RawDatapoint, error) {
	rawDatapoints := make([]*RawDatapoint, len(data))
	for i, datapoint := range data {
		message, err := proto.Marshal(datapoint.Data)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal datapoint data: %w", err)
		}

		rawDatapoints[i] = &RawDatapoint{
			Meta: datapoint.Meta,
			Data: message,
		}
	}
	return rawDatapoints, nil
}

// CollectAs converts a sequence of RawDatapoint into a slice of TypedDatapoint with the given type.
func CollectAs[T proto.Message](seq iter.Seq2[*RawDatapoint, error]) ([]*TypedDatapoint[T], error) {
	return Collect(As[T](seq))
}

// Collect converts any sequence into a slice.
//
// It returns an error if any of the elements in the sequence has a non-nil error.
func Collect[K any](seq iter.Seq2[K, error]) ([]K, error) {
	s := make([]K, 0)

	for k, err := range seq {
		if err != nil {
			return nil, err
		}
		s = append(s, k)
	}
	return s, nil
}

// As converts a sequence of RawDatapoint into a sequence of TypedDatapoint with the given type.
func As[T proto.Message](seq iter.Seq2[*RawDatapoint, error]) iter.Seq2[*TypedDatapoint[T], error] {
	var t T
	descriptor := reflect.New(reflect.TypeOf(t).Elem()).Interface().(T).ProtoReflect()

	return func(yield func(*TypedDatapoint[T], error) bool) {
		for rawDatapoint, err := range seq {
			if err != nil {
				yield(nil, err)
				return
			}

			data := descriptor.New().Interface().(T)
			err = proto.Unmarshal(rawDatapoint.Data, data)
			if err != nil {
				yield(nil, fmt.Errorf("failed to unmarshal datapoint data: %w", err))
				return
			}

			datapoint := &TypedDatapoint[T]{
				Meta: rawDatapoint.Meta,
				Data: data,
			}
			if !yield(datapoint, nil) {
				return
			}
		}
	}
}
