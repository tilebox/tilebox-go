package datasets // import "github.com/tilebox/tilebox-go/datasets/v1"

import (
	"context"
	"errors"
	"fmt"
	"iter"
	"reflect"

	"github.com/google/uuid"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/ewkb"
	datasetsv1 "github.com/tilebox/tilebox-go/protogen/go/datasets/v1"
	"github.com/tilebox/tilebox-go/query"
	"google.golang.org/protobuf/proto"
)

type DatapointClient interface {
	Query(ctx context.Context, collectionIDs []uuid.UUID, options ...QueryOption) iter.Seq2[[]byte, error]
	QueryInto(ctx context.Context, collectionIDs []uuid.UUID, datapoints any, options ...QueryOption) error
	Ingest(ctx context.Context, collectionID uuid.UUID, datapoints any, allowExisting bool) (*IngestResponse, error)
	Delete(ctx context.Context, collectionID uuid.UUID, datapoints any) (*DeleteResponse, error)
	DeleteIDs(ctx context.Context, collectionID uuid.UUID, datapointIDs []uuid.UUID) (*DeleteResponse, error)
}

var _ DatapointClient = &datapointClient{}

type datapointClient struct {
	dataIngestionService DataIngestionService
	dataAccessService    DataAccessService
}

// queryOptions contains the configuration for a Query request.
type queryOptions struct {
	temporalExtent query.TemporalExtent
	spatialExtent  orb.Geometry
	skipData       bool
}

// QueryOption is an interface for configuring a Query request.
type QueryOption func(*queryOptions)

// WithTemporalExtent specifies the time interval for which data should be queried.
// Right now, a temporal extent is required for every query.
func WithTemporalExtent(temporalExtent query.TemporalExtent) QueryOption {
	return func(cfg *queryOptions) {
		cfg.temporalExtent = temporalExtent
	}
}

// WithSpatialExtent specifies the geographical extent in which to query data.
// Optional, if not specified the query will return all results found globally.
func WithSpatialExtent(spatialExtent orb.Geometry) QueryOption {
	return func(cfg *queryOptions) {
		cfg.spatialExtent = spatialExtent
	}
}

// WithSkipData skips the data when querying datapoints.
// It is an optional flag for omitting the actual datapoint data from the response.
// If set, only datapoint IDs will be returned.
//
// Defaults to false.
func WithSkipData() QueryOption {
	return func(cfg *queryOptions) {
		cfg.skipData = true
	}
}

// Query datapoints from one or more collections of the same dataset.
//
// Options:
//   - WithTemporalExtent: specifies the time or data point interval for which data should be loaded. (Required)
//   - WithSpatialExtent: specifies the spatial extent for which data should be loaded. (Optional)
//   - WithSkipData: can be used to skip the actual data when loading datapoints, only returning matching datapoint IDs. (Optional)
//
// The datapoints are loaded in a lazy manner, and returned as a sequence of bytes.
// The output sequence can be transformed into a proto.Message using `CollectAs`/`As`.
//
// Example usage:
//
//	for datapointBytes, err := range client.Datapoints.Query(ctx, collectionIDs, WithTemporalExtent(timeInterval), WithSpatialExtent(geometry)) {
//	  if err != nil {
//	    // handle error
//	  }
//	  datapoint := &datasetsv1.CopernicusDataspaceGranule{}
//	  err = proto.Unmarshal(datapointBytes, datapoint)
//	  if err != nil {
//	    // handle unmarshal error
//	  }
//	  // do something with the datapoint
//	}
//
// Documentation: https://docs.tilebox.com/datasets/query
func (d datapointClient) Query(ctx context.Context, collectionIDs []uuid.UUID, options ...QueryOption) iter.Seq2[[]byte, error] {
	cfg := &queryOptions{
		skipData: false,
	}
	for _, option := range options {
		option(cfg)
	}

	if cfg.temporalExtent == nil {
		return func(yield func([]byte, error) bool) {
			// right now we return an error, in the future we might want to support queries without a temporal extent
			yield(nil, errors.New("temporal extent is required"))
		}
	}

	return func(yield func([]byte, error) bool) {
		var page *datasetsv1.Pagination // nil for the first request

		// we already validated that temporalExtent is not nil
		timeInterval := cfg.temporalExtent.ToProtoTimeInterval()
		datapointInterval := cfg.temporalExtent.ToProtoDatapointInterval()

		if timeInterval == nil && datapointInterval == nil {
			yield(nil, errors.New("invalid temporal extent"))
			return
		}

		filters := &datasetsv1.QueryFilters{}
		if timeInterval != nil {
			filters.TemporalExtent = &datasetsv1.QueryFilters_TimeInterval{TimeInterval: timeInterval}
		} else {
			filters.TemporalExtent = &datasetsv1.QueryFilters_DatapointInterval{DatapointInterval: datapointInterval}
		}

		geometry := cfg.spatialExtent
		if geometry != nil {
			wkb, err := ewkb.Marshal(geometry, ewkb.DefaultSRID)
			if err != nil {
				yield(nil, fmt.Errorf("invalid geometry: failed to marshal geometry as wkb: %w", err))
				return
			}
			filters.SpatialExtent = &datasetsv1.SpatialFilter{
				Geometry: &datasetsv1.Geometry{Wkb: wkb},
			}
		}

		for {
			datapointsMessage, err := d.dataAccessService.Query(ctx, collectionIDs, filters, page, cfg.skipData)
			if err != nil {
				yield(nil, err)
				return
			}

			for _, data := range datapointsMessage.GetData().GetValue() {
				if !yield(data, nil) {
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

// QueryInto queries datapoints from one or more collections of the same dataset into a slice of datapoints of a
// compatible proto.Message type.
//
// QueryInto is a convenience function for Query, when no pagination or manual iteration is required.
//
// Example usage:
// var datapoints []*datasetsv1.CopernicusDataspaceGranule
// err := client.Datapoints.QueryInto(ctx, collectionIDs, &datapoints, WithTemporalExtent(timeInterval))
func (d datapointClient) QueryInto(ctx context.Context, collectionIDs []uuid.UUID, datapoints any, options ...QueryOption) error {
	err := validateDatapoints(datapoints)
	if err != nil {
		return err // already a nice validation error
	}

	slice := reflect.Indirect(reflect.ValueOf(datapoints))
	datapointType := slice.Type().Elem().Elem()

	rawDatapoints, err := Collect(d.Query(ctx, collectionIDs, options...))
	if err != nil {
		return err
	}

	slice.Set(reflect.MakeSlice(slice.Type(), len(rawDatapoints), len(rawDatapoints)))

	for i, rawDatapoint := range rawDatapoints {
		datapoint := reflect.New(datapointType)

		err = proto.Unmarshal(rawDatapoint, datapoint.Interface().(proto.Message))
		if err != nil {
			return fmt.Errorf("failed to unmarshal datapoint: %w", err)
		}

		slice.Index(i).Set(datapoint)
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
func (d datapointClient) Ingest(ctx context.Context, collectionID uuid.UUID, datapoints any, allowExisting bool) (*IngestResponse, error) {
	err := validateDatapoints(datapoints)
	if err != nil {
		return nil, fmt.Errorf("failed to validate datapoints: %w", err)
	}

	slice := reflect.Indirect(reflect.ValueOf(datapoints))

	marshaledDatapoints := make([][]byte, slice.Len())
	for i := range slice.Len() {
		datapoint := slice.Index(i).Interface().(proto.Message)
		marshal, err := proto.Marshal(datapoint)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal datapoint: %w", err)
		}

		marshaledDatapoints[i] = marshal
	}

	response, err := d.dataIngestionService.Ingest(ctx, collectionID, marshaledDatapoints, allowExisting)
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

func validateDatapoints(datapoints any) error {
	rv := reflect.ValueOf(datapoints)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return fmt.Errorf("datapoints must be a pointer, got %v", reflect.TypeOf(datapoints))
	}
	slice := reflect.Indirect(rv)
	if slice.Kind() != reflect.Slice {
		return fmt.Errorf("datapoints must be a pointer to a slice, got %v", reflect.TypeOf(datapoints))
	}
	if slice.Type().Elem().Kind() != reflect.Ptr {
		return fmt.Errorf("datapoints must be a pointer to a slice of proto.Message, got %v", reflect.TypeOf(datapoints))
	}
	messageType := reflect.TypeOf((*proto.Message)(nil)).Elem()
	if !slice.Type().Elem().Implements(messageType) {
		return fmt.Errorf("datapoints must be a pointer to a slice of proto.Message, got %v", reflect.TypeOf(datapoints))
	}
	return nil
}

// DeleteResponse contains the response from the Delete method.
type DeleteResponse struct {
	// NumDeleted is the number of datapoints that were deleted.
	NumDeleted int64
}

// Delete datapoints from a collection.
//
// The datapoints are identified by their IDs.
func (d datapointClient) Delete(ctx context.Context, collectionID uuid.UUID, datapoints any) (*DeleteResponse, error) {
	err := validateDatapoints(datapoints)
	if err != nil {
		return nil, fmt.Errorf("failed to validate datapoints: %w", err)
	}

	slice := reflect.Indirect(reflect.ValueOf(datapoints))

	datapointIDs := make([]uuid.UUID, slice.Len())
	for i := range slice.Len() {
		datapoint := slice.Index(i).Interface().(proto.Message)
		idFieldDescriptor := datapoint.ProtoReflect().Descriptor().Fields().ByName("id")
		if idFieldDescriptor == nil {
			return nil, errors.New("failed to find id field in datapoint")
		}

		idField := datapoint.ProtoReflect().Get(idFieldDescriptor).Message().Interface().(*datasetsv1.UUID)
		id, err := uuid.FromBytes(idField.GetUuid())
		if err != nil {
			return nil, fmt.Errorf("failed to parse datapoint id: %w", err)
		}
		datapointIDs[i] = id
	}

	return d.DeleteIDs(ctx, collectionID, datapointIDs)
}

// DeleteIDs deletes datapoints from a collection by their IDs.
func (d datapointClient) DeleteIDs(ctx context.Context, collectionID uuid.UUID, datapointIDs []uuid.UUID) (*DeleteResponse, error) {
	response, err := d.dataIngestionService.Delete(ctx, collectionID, datapointIDs)
	if err != nil {
		return nil, err
	}

	return &DeleteResponse{
		NumDeleted: response.GetNumDeleted(),
	}, nil
}

// CollectAs converts a sequence of bytes into a slice of proto.Message.
func CollectAs[T proto.Message](seq iter.Seq2[[]byte, error]) ([]T, error) {
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

// As converts a sequence of bytes into a sequence of proto.Message.
func As[T proto.Message](seq iter.Seq2[[]byte, error]) iter.Seq2[T, error] {
	var t T
	descriptor := reflect.New(reflect.TypeOf(t).Elem()).Interface().(T).ProtoReflect()

	return func(yield func(T, error) bool) {
		for rawDatapoint, err := range seq {
			if err != nil {
				yield(t, err)
				return
			}

			datapoint := descriptor.New().Interface().(T)
			err = proto.Unmarshal(rawDatapoint, datapoint)
			if err != nil {
				yield(t, fmt.Errorf("failed to unmarshal datapoint: %w", err))
				return
			}

			if !yield(datapoint, nil) {
				return
			}
		}
	}
}
