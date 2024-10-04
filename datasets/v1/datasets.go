// Package datasets provides a client for interacting with Tilebox Datasets.
//
// Documentation: https://docs.tilebox.com/datasets
package datasets // import "github.com/tilebox/tilebox-go/datasets/v1"

import (
	"context"
	"errors"
	"fmt"
	"iter"
	"reflect"
	"time"

	"github.com/google/uuid"
	datasetsv1 "github.com/tilebox/tilebox-go/protogen/go/datasets/v1"
	"github.com/tilebox/tilebox-go/protogen/go/datasets/v1/datasetsv1connect"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const otelTracerName = "tilebox.com/observability"

// Client is a Tilebox Datasets client.
type Client struct {
	service Service
}

// NewClient creates a new Tilebox Datasets client.
//
// By default, the returned Client is configured with:
//   - "https://api.tilebox.com" as the URL
//   - no API key
//   - a grpc.RetryHTTPClient HTTP client
//   - the global tracer provider
//
// The passed options are used to override these default values and configure the returned Client appropriately.
func NewClient(options ...ClientOption) *Client {
	cfg := newClientConfig(options)
	connectClient := newConnectClient(datasetsv1connect.NewTileboxServiceClient, cfg)

	return &Client{
		service: newDatasetsService(connectClient, cfg.tracerProvider.Tracer(otelTracerName)),
	}
}

// Datasets returns a list of all available datasets.
func (c *Client) Datasets(ctx context.Context) ([]*Dataset, error) {
	response, err := c.service.ListDatasets(ctx)
	if err != nil {
		return nil, err
	}

	datasets := make([]*Dataset, len(response.GetDatasets()))
	for i, d := range response.GetDatasets() {
		dataset, err := protoToDataset(d, c.service)
		if err != nil {
			return nil, fmt.Errorf("failed to convert dataset from response: %w", err)
		}

		datasets[i] = dataset
	}

	return datasets, nil
}

// Dataset returns a dataset by its slug.
/*func (c *Client) Dataset(ctx context.Context, slug string) (*Dataset, error) {
	response, err := c.service.GetDataset(ctx, slug)
	if err != nil {
		return nil, err
	}

	dataset, err := protoToDataset(response, c.service)
	if err != nil {
		return nil, fmt.Errorf("failed to convert dataset from response: %w", err)
	}

	return dataset, nil
}*/

// Dataset represents a Tilebox Time Series Dataset.
//
// Documentation: https://docs.tilebox.com/datasets/timeseries
type Dataset struct {
	// ID is the unique identifier of the dataset.
	ID uuid.UUID
	// Slug is the unique slug of the dataset.
	// Slug string

	// Name is the name of the dataset.
	Name string
	// Summary is a summary of the purpose of the dataset.
	Summary string

	service Service
}

func protoToDataset(d *datasetsv1.Dataset, service Service) (*Dataset, error) {
	id, err := protoToUUID(d.GetId())
	if err != nil {
		return nil, fmt.Errorf("failed to convert dataset id to uuid: %w", err)
	}

	return &Dataset{
		ID: id,
		// Slug:     d.GetSlug(),
		Name:    d.GetName(),
		Summary: d.GetSummary(),
		service: service,
	}, nil
}

func protoToUUID(id *datasetsv1.ID) (uuid.UUID, error) {
	if id == nil || len(id.GetUuid()) == 0 {
		return uuid.Nil, nil
	}

	bytes, err := uuid.FromBytes(id.GetUuid())
	if err != nil {
		return uuid.Nil, err
	}

	return bytes, nil
}

// Collections returns a list of all available collections in the dataset.
func (d *Dataset) Collections(ctx context.Context) ([]*Collection, error) {
	response, err := d.service.GetCollections(ctx, d.ID)
	if err != nil {
		return nil, err
	}

	collections := make([]*Collection, len(response.GetData()))
	for i, c := range response.GetData() {
		collection, err := protoToCollection(c, d.service)
		if err != nil {
			return nil, fmt.Errorf("failed to convert collection from response: %w", err)
		}

		collections[i] = collection
	}

	return collections, nil
}

// Collection returns a collection by its name.
func (d *Dataset) Collection(ctx context.Context, name string) (*Collection, error) {
	response, err := d.service.GetCollectionByName(ctx, d.ID, name)
	if err != nil {
		return nil, err
	}

	collection, err := protoToCollection(response, d.service)
	if err != nil {
		return nil, fmt.Errorf("failed to convert collection from response: %w", err)
	}

	return collection, nil
}

// CreateCollection creates a new collection in the dataset with the given name.
func (d *Dataset) CreateCollection(ctx context.Context, collectionName string) (*Collection, error) {
	response, err := d.service.CreateCollection(ctx, d.ID, collectionName)
	if err != nil {
		return nil, err
	}

	collection, err := protoToCollection(response, d.service)
	if err != nil {
		return nil, fmt.Errorf("failed to convert collection from response: %w", err)
	}

	return collection, nil
}

// Collection represents a Tilebox Time Series Dataset collection.
//
// Documentation: https://docs.tilebox.com/datasets/collections
type Collection struct {
	// ID is the unique identifier of the collection.
	ID uuid.UUID
	// Name is the name of the collection.
	Name string
	// Availability is the time interval for which data is available.
	Availability TimeInterval
	// Count is the number of datapoints in the collection.
	Count uint64

	service Service
}

func protoToCollection(c *datasetsv1.CollectionInfo, service Service) (*Collection, error) {
	id, err := uuid.Parse(c.GetCollection().GetId())
	if err != nil {
		return nil, fmt.Errorf("failed to parse collection id: %w", err)
	}

	return &Collection{
		ID:           id,
		Name:         c.GetCollection().GetName(),
		Availability: *protoToTimeInterval(c.GetAvailability()),
		Count:        c.GetCount(),
		service:      service,
	}, nil
}

// Load loads datapoints from the collection.
//
// loadInterval specifies the time interval for which data should be loaded.
//
// skipData and skipMeta specify whether the data and metadata should be skipped.
// If both skipData and skipMeta are true, the response will only consist of a list of datapoint IDs without any
// additional data or metadata.
//
// The datapoints are loaded in a lazy manner, and returned as a sequence of RawDatapoint.
// The output sequence can be transformed into typed Datapoint using CollectAs or As functions.
//
// Documentation: https://docs.tilebox.com/datasets/loading-data
func (c *Collection) Load(ctx context.Context, loadInterval LoadInterval, skipData, skipMeta bool) iter.Seq2[*RawDatapoint, error] {
	return func(yield func(*RawDatapoint, error) bool) {
		var page *datasetsv1.Pagination // nil for the first request

		timeInterval := loadInterval.ToProtoTimeInterval()
		datapointInterval := loadInterval.ToProtoDatapointInterval()

		if timeInterval == nil && datapointInterval == nil {
			yield(nil, errors.New("time interval and datapoint interval cannot both be nil"))
			return
		}

		for {
			datapointsMessage, err := c.service.GetDatasetForInterval(ctx, c.ID, timeInterval, datapointInterval, page, skipData, skipMeta)
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

				if !skipMeta {
					meta.EventTime = dp.GetEventTime().AsTime()
					meta.IngestionTime = dp.GetIngestionTime().AsTime()
				}

				if !skipData {
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

// IngestResponse contains the response from the Ingest method.
type IngestResponse struct {
	// NumCreated is the number of datapoints that were created.
	NumCreated int64
	// NumExisting is the number of datapoints that were ignored because they already existed.
	NumExisting int64
	// DatapointIDs is the list of all the datapoints IDs in the same order as the datapoints in the request.
	DatapointIDs []uuid.UUID
}

// Ingest ingests datapoints into the collection.
//
// data is a list of datapoints to ingest that should be created using Datapoints.
//
// allowExisting specifies whether to allow existing datapoints as part of the request. If true, datapoints that already
// exist will be ignored, and the number of such existing datapoints will be returned in the response. If false, any
// datapoints that already exist will result in an error. Setting this to true is useful for achieving idempotency (e.g.
// allowing re-ingestion of datapoints that have already been ingested in the past).
func (c *Collection) Ingest(ctx context.Context, data []*RawDatapoint, allowExisting bool) (*IngestResponse, error) {
	datapoints := &datasetsv1.Datapoints{
		Meta: make([]*datasetsv1.DatapointMetadata, len(data)),
		Data: &datasetsv1.RepeatedAny{
			Value: make([][]byte, len(data)),
		},
	}

	for i, datapoint := range data {
		datapoints.GetMeta()[i] = &datasetsv1.DatapointMetadata{
			EventTime: timestamppb.New(datapoint.Meta.EventTime),
		}
		datapoints.GetData().GetValue()[i] = datapoint.Data
	}

	response, err := c.service.IngestDatapoints(ctx, c.ID, datapoints, allowExisting)
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

// Delete deletes datapoints from the collection.
//
// The datapoints are identified by their IDs.
func (c *Collection) Delete(ctx context.Context, data []*RawDatapoint) (*DeleteResponse, error) {
	datapointIDs := make([]uuid.UUID, len(data))
	for i, datapoint := range data {
		datapointIDs[i] = datapoint.Meta.ID
	}

	return c.DeleteIDs(ctx, datapointIDs)
}

// DeleteIDs deletes datapoints from the collection by their IDs.
func (c *Collection) DeleteIDs(ctx context.Context, datapointIDs []uuid.UUID) (*DeleteResponse, error) {
	response, err := c.service.DeleteDatapoints(ctx, c.ID, datapointIDs)
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

// Datapoint represents a datapoint in a collection.
// It contains the metadata and the data itself.
type Datapoint[T proto.Message] struct {
	// Meta contains the metadata of the datapoint.
	Meta *DatapointMetadata
	// Data contains the data of the datapoint.
	Data T
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
func NewDatapoint[T proto.Message](time time.Time, message T) *Datapoint[T] {
	return &Datapoint[T]{
		Meta: &DatapointMetadata{
			EventTime: time,
		},
		Data: message,
	}
}

// Datapoints converts a list of Datapoint to RawDatapoint.
//
// It is used to convert the data before ingesting it into a collection.
func Datapoints[T proto.Message](data ...*Datapoint[T]) ([]*RawDatapoint, error) {
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

// CollectAs converts a sequence of RawDatapoint into a slice of Datapoint with the given type.
func CollectAs[T proto.Message](seq iter.Seq2[*RawDatapoint, error]) ([]*Datapoint[T], error) {
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

// As converts a sequence of RawDatapoint into a sequence of Datapoint with the given type.
func As[T proto.Message](seq iter.Seq2[*RawDatapoint, error]) iter.Seq2[*Datapoint[T], error] {
	var t T
	descriptor := reflect.New(reflect.TypeOf(t).Elem()).Interface().(T).ProtoReflect()

	return func(yield func(*Datapoint[T], error) bool) {
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

			datapoint := &Datapoint[T]{
				Meta: rawDatapoint.Meta,
				Data: data,
			}
			if !yield(datapoint, nil) {
				return
			}
		}
	}
}
