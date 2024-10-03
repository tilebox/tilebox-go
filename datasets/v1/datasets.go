package datasets

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

type Client struct {
	service Service
}

func NewClient(options ...ClientOption) *Client {
	cfg := newClientConfig(options)
	connectClient := newConnectClient(datasetsv1connect.NewTileboxServiceClient, cfg)

	return &Client{
		service: newDatasetsService(connectClient, cfg.tracerProvider.Tracer("tilebox.com/observability")),
	}
}

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

func (c *Client) Dataset(ctx context.Context, slug string) (*Dataset, error) {
	// FIXME: use new endpoint with slug
	return nil, errors.New("not implemented")
	/*response, err := c.service.GetDataset(ctx, slug)
	if err != nil {
		return nil, err
	}

	dataset, err := protoToDataset(response, c.service)
	if err != nil {
		return nil, fmt.Errorf("failed to convert dataset from response: %w", err)
	}

	return dataset, nil*/
}

type Dataset struct {
	ID      uuid.UUID
	Slug    string
	Name    string
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
		// Slug:     d.GetSlug(), FIXME
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

type Collection struct {
	ID           uuid.UUID
	Name         string
	Availability TimeInterval
	Count        uint64

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

type IngestResponse struct {
	NumCreated   int64
	NumExisting  int64
	DatapointIDs []uuid.UUID
}

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

	return &IngestResponse{
		NumCreated:  response.GetNumCreated(),
		NumExisting: response.GetNumExisting(),
		// DatapointIDs: response.GetDatapointIds(), FIXME
	}, nil
}

type DeleteResponse struct {
	NumDeleted int64
}

func (c *Collection) Delete(ctx context.Context, data []*RawDatapoint) (*DeleteResponse, error) {
	datapointIDs := make([]uuid.UUID, len(data))
	for i, datapoint := range data {
		datapointIDs[i] = datapoint.Meta.ID
	}

	return c.DeleteIDs(ctx, datapointIDs)
}

func (c *Collection) DeleteIDs(ctx context.Context, datapointIDs []uuid.UUID) (*DeleteResponse, error) {
	response, err := c.service.DeleteDatapoints(ctx, c.ID, datapointIDs)
	if err != nil {
		return nil, err
	}

	return &DeleteResponse{
		NumDeleted: response.GetNumDeleted(),
	}, nil
}

type DatapointMetadata struct {
	ID            uuid.UUID
	EventTime     time.Time
	IngestionTime time.Time
}

type Datapoint[T proto.Message] struct {
	Meta *DatapointMetadata
	Data T
}

type RawDatapoint struct {
	Meta *DatapointMetadata
	Data []byte
}

func ValidateAs[T proto.Message](data []*Datapoint[T]) ([]*RawDatapoint, error) {
	rawDatapoints := make([]*RawDatapoint, len(data))
	for i, datapoint := range data {
		message, err := proto.Marshal(datapoint.Data)
		if err != nil {
			return nil, errors.New("failed to marshal datapoint data")
		}

		rawDatapoints[i] = &RawDatapoint{
			Meta: datapoint.Meta,
			Data: message,
		}
	}
	return rawDatapoints, nil
}

func CollectAs[T proto.Message](seq iter.Seq2[*RawDatapoint, error]) ([]*Datapoint[T], error) {
	return Collect(As[T](seq))
}

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
