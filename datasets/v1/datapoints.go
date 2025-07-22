package datasets // import "github.com/tilebox/tilebox-go/datasets/v1"

import (
	"context"
	"errors"
	"fmt"
	"iter"
	"reflect"

	"github.com/google/uuid"
	"github.com/paulmach/orb"
	datasetsv1 "github.com/tilebox/tilebox-go/protogen/datasets/v1"
	tileboxv1 "github.com/tilebox/tilebox-go/protogen/tilebox/v1"
	"github.com/tilebox/tilebox-go/query"
	"google.golang.org/protobuf/proto"
)

type DatapointClient interface {
	// GetInto get a single datapoint by its ID from one or more collections of the same dataset into a proto.Message.
	//
	// Options:
	//   - WithSkipData: can be used to skip the actual data, only returning required datapoint fields. (Optional)
	//
	// Example usage:
	//
	//	var datapoint v1.Sentinel1Sar
	//	err = client.Datapoints.GetInto(ctx, collectionIDs, datapointID, &datapoint)
	GetInto(ctx context.Context, collectionIDs []uuid.UUID, datapointID uuid.UUID, datapoint proto.Message, options ...QueryOption) error

	// Query datapoints from one or more collections of the same dataset.
	//
	// Options:
	//   - WithTemporalExtent: specifies the time or data point interval for which data should be loaded. (Required)
	//   - WithSpatialExtent: specifies the spatial extent for which data should be loaded. (Optional)
	//   - WithSkipData: can be used to skip the actual data when loading datapoints, only returning required datapoint fields. (Optional)
	//
	// The datapoints are lazily loaded and returned as a sequence of bytes.
	// The output sequence can be transformed into a proto.Message using CollectAs / As.
	//
	// Example usage:
	//
	//	for datapointBytes, err := range client.Datapoints.Query(ctx, collectionIDs, WithTemporalExtent(timeInterval), WithSpatialExtent(geometry)) {
	//	  if err != nil {
	//	    // handle error
	//	  }
	//	  datapoint := &v1.Sentinel1Sar{}
	//	  err = proto.Unmarshal(datapointBytes, datapoint)
	//	  if err != nil {
	//	    // handle unmarshal error
	//	  }
	//	  // do something with the datapoint
	//	}
	//
	// Documentation: https://docs.tilebox.com/datasets/query
	Query(ctx context.Context, collectionIDs []uuid.UUID, options ...QueryOption) iter.Seq2[[]byte, error]

	// QueryInto queries datapoints from one or more collections of the same dataset into a slice of datapoints of a
	// compatible proto.Message type.
	//
	// QueryInto is a convenience function for Query, when no manual pagination or custom iteration is required.
	//
	// Example usage:
	//
	//	var datapoints []*v1.Sentinel1Sar
	//	err := client.Datapoints.QueryInto(ctx, collectionIDs, &datapoints, WithTemporalExtent(timeInterval))
	QueryInto(ctx context.Context, collectionIDs []uuid.UUID, datapoints any, options ...QueryOption) error

	// Ingest datapoints into a collection.
	//
	// data is a list of datapoints to ingest that should be created using Datapoints.
	//
	// allowExisting specifies whether to allow existing datapoints as part of the request. If true, datapoints that already
	// exist will be ignored, and the number of such existing datapoints will be returned in the response. If false, any
	// datapoints that already exist will result in an error. Setting this to true is useful for achieving idempotency (e.g.
	// allowing re-ingestion of datapoints that have already been ingested in the past).
	Ingest(ctx context.Context, collectionID uuid.UUID, datapoints any, allowExisting bool) (*IngestResponse, error)

	// Delete datapoints from a collection.
	//
	// The datapoints are identified by their IDs.
	//
	// Returns the number of deleted datapoints.
	Delete(ctx context.Context, collectionID uuid.UUID, datapoints any) (int64, error)

	// DeleteIDs deletes datapoints from a collection by their IDs.
	//
	// Returns the number of deleted datapoints.
	DeleteIDs(ctx context.Context, collectionID uuid.UUID, datapointIDs []uuid.UUID) (int64, error)
}

var _ DatapointClient = &datapointClient{}

type datapointClient struct {
	dataIngestionService DataIngestionService
	dataAccessService    DataAccessService
}

// queryOptions contains the configuration for a Query request.
type queryOptions struct {
	temporalExtent query.TemporalExtent
	spatialExtent  query.SpatialExtent
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
// Convenience wrapper for WithSpatialExtentFilter applying a default spatial filter mode and coordinate system.
func WithSpatialExtent(geometry orb.Geometry) QueryOption {
	return WithSpatialExtentFilter(&query.SpatialFilter{Geometry: geometry})
}

// WithSpatialExtentFilter specifies a geographical extent as well as an intersection mode and coordinate system to use
// for spatial filtering when querying data.
func WithSpatialExtentFilter(spatialExtent query.SpatialExtent) QueryOption {
	return func(cfg *queryOptions) {
		cfg.spatialExtent = spatialExtent
	}
}

// WithSkipData skips the data when querying datapoints.
// It is an optional flag for omitting the actual datapoint data from the response.
// If set, only the required datapoint fields will be returned.
//
// Defaults to false.
func WithSkipData() QueryOption {
	return func(cfg *queryOptions) {
		cfg.skipData = true
	}
}

func (d datapointClient) GetInto(ctx context.Context, collectionIDs []uuid.UUID, datapointID uuid.UUID, datapoint proto.Message, options ...QueryOption) error {
	cfg := &queryOptions{
		skipData: false,
	}
	for _, option := range options {
		option(cfg)
	}

	rawDatapoint, err := d.dataAccessService.QueryByID(ctx, collectionIDs, datapointID, cfg.skipData)
	if err != nil {
		return err
	}

	err = proto.Unmarshal(rawDatapoint.GetValue(), datapoint)
	if err != nil {
		return fmt.Errorf("failed to unmarshal datapoint: %w", err)
	}

	return nil
}

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
		var page *tileboxv1.Pagination // nil for the first request

		// we already validated that temporalExtent is not nil
		timeInterval := cfg.temporalExtent.ToProtoTimeInterval()
		datapointInterval := cfg.temporalExtent.ToProtoIDInterval()

		if timeInterval == nil && datapointInterval == nil {
			yield(nil, errors.New("invalid temporal extent"))
			return
		}

		filters := datasetsv1.QueryFilters_builder{
			TimeInterval:      timeInterval,
			DatapointInterval: datapointInterval,
		}.Build()

		if cfg.spatialExtent != nil {
			spatialExtent, err := cfg.spatialExtent.ToProtoSpatialFilter()
			if err != nil {
				yield(nil, err)
				return
			}
			filters.SetSpatialExtent(spatialExtent)
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

const (
	ingestChunkSize = 8192
	deleteChunkSize = 8192
)

// IngestResponse contains the response from the Ingest method.
type IngestResponse struct {
	// NumCreated is the number of datapoints that were created.
	NumCreated int64
	// NumExisting is the number of datapoints that were ignored because they already existed.
	NumExisting int64
	// DatapointIDs is the list of all the datapoints IDs in the same order as the datapoints in the request.
	DatapointIDs []uuid.UUID
}

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

	numCreated := int64(0)
	numExisting := int64(0)
	datapointIDs := make([]uuid.UUID, 0, len(marshaledDatapoints))

	for i := 0; i < len(marshaledDatapoints); i += ingestChunkSize {
		chunk := marshaledDatapoints[i:min(i+ingestChunkSize, len(marshaledDatapoints))]
		response, err := d.dataIngestionService.Ingest(ctx, collectionID, chunk, allowExisting)
		if err != nil {
			return nil, err
		}
		numCreated += response.GetNumCreated()
		numExisting += response.GetNumExisting()
		for _, id := range response.GetDatapointIds() {
			datapointIDs = append(datapointIDs, id.AsUUID())
		}
	}

	return &IngestResponse{
		NumCreated:   numCreated,
		NumExisting:  numExisting,
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
	messageType := reflect.TypeFor[proto.Message]()
	if !slice.Type().Elem().Implements(messageType) {
		return fmt.Errorf("datapoints must be a pointer to a slice of proto.Message, got %v", reflect.TypeOf(datapoints))
	}
	return nil
}

func (d datapointClient) Delete(ctx context.Context, collectionID uuid.UUID, datapoints any) (int64, error) {
	err := validateDatapoints(datapoints)
	if err != nil {
		return 0, fmt.Errorf("failed to validate datapoints: %w", err)
	}

	slice := reflect.Indirect(reflect.ValueOf(datapoints))

	datapointIDs := make([]uuid.UUID, slice.Len())
	for i := range slice.Len() {
		datapoint := slice.Index(i).Interface().(proto.Message)
		idFieldDescriptor := datapoint.ProtoReflect().Descriptor().Fields().ByName("id")
		if idFieldDescriptor == nil {
			return 0, errors.New("failed to find id field in datapoint")
		}

		idField := datapoint.ProtoReflect().Get(idFieldDescriptor).Message().Interface().(*datasetsv1.UUID)
		id, err := uuid.FromBytes(idField.GetUuid())
		if err != nil {
			return 0, fmt.Errorf("failed to parse datapoint id: %w", err)
		}
		datapointIDs[i] = id
	}

	return d.DeleteIDs(ctx, collectionID, datapointIDs)
}

func (d datapointClient) DeleteIDs(ctx context.Context, collectionID uuid.UUID, datapointIDs []uuid.UUID) (int64, error) {
	numDeleted := int64(0)

	for i := 0; i < len(datapointIDs); i += deleteChunkSize {
		chunk := datapointIDs[i:min(i+deleteChunkSize, len(datapointIDs))]
		chunkNumDeleted, err := d.DeleteIDs(ctx, collectionID, chunk)
		if err != nil {
			return 0, err
		}
		numDeleted += chunkNumDeleted
	}

	response, err := d.dataIngestionService.Delete(ctx, collectionID, datapointIDs)
	if err != nil {
		return 0, err
	}

	return response.GetNumDeleted(), nil
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
