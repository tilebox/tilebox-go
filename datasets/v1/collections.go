package datasets // import "github.com/tilebox/tilebox-go/datasets/v1"

import (
	"context"
	"errors"
	"fmt"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	datasetsv1 "github.com/tilebox/tilebox-go/protogen/datasets/v1"
	"github.com/tilebox/tilebox-go/query"
)

// Collection represents a Tilebox Dataset collection.
//
// Documentation: https://docs.tilebox.com/datasets/concepts/collections
type Collection struct {
	// ID is the unique identifier of the collection.
	ID uuid.UUID
	// Name is the name of the collection.
	Name string
	// Availability is the time interval for which data is available.
	Availability *query.TimeInterval
	// Count is the number of datapoints in the collection.
	Count uint64
}

func (c Collection) String() string {
	availability := "<availability unknown>"
	if c.Availability != nil {
		availability = c.Availability.String()
	}

	count := ""
	if c.Count > 0 {
		count = fmt.Sprintf(" (%d data points)", c.Count)
	}

	return fmt.Sprintf("Collection %s: %s%s", c.Name, availability, count)
}

type CollectionClient interface {
	// Create creates a new collection in the dataset with the given name.
	Create(ctx context.Context, datasetID uuid.UUID, name string) (*Collection, error)

	// Get returns a collection by its name.
	Get(ctx context.Context, datasetID uuid.UUID, name string) (*Collection, error)

	// GetOrCreate returns a collection by its name, creating it if it does not exist.
	GetOrCreate(ctx context.Context, datasetID uuid.UUID, name string) (*Collection, error)

	// Delete deletes a collection by its id.
	Delete(ctx context.Context, datasetID uuid.UUID, collectionID uuid.UUID) error

	// List returns a list of all available collections in the dataset.
	List(ctx context.Context, datasetID uuid.UUID) ([]*Collection, error)
}

var _ CollectionClient = &collectionClient{}

type collectionClient struct {
	service CollectionService
}

func (c collectionClient) Create(ctx context.Context, datasetID uuid.UUID, name string) (*Collection, error) {
	response, err := c.service.CreateCollection(ctx, datasetID, name)
	if err != nil {
		return nil, err
	}

	return protoToCollection(response), nil
}

func (c collectionClient) Get(ctx context.Context, datasetID uuid.UUID, name string) (*Collection, error) {
	response, err := c.service.GetCollectionByName(ctx, datasetID, name)
	if err != nil {
		return nil, err
	}

	return protoToCollection(response), nil
}

func (c collectionClient) GetOrCreate(ctx context.Context, datasetID uuid.UUID, name string) (*Collection, error) {
	collection, err := c.Get(ctx, datasetID, name)
	if err != nil {
		var connectErr *connect.Error
		if errors.As(err, &connectErr) && connectErr.Code() == connect.CodeNotFound {
			return c.Create(ctx, datasetID, name)
		}
	}

	return collection, err
}

func (c collectionClient) Delete(ctx context.Context, datasetID uuid.UUID, collectionID uuid.UUID) error {
	return c.service.DeleteCollection(ctx, datasetID, collectionID)
}

func (c collectionClient) List(ctx context.Context, datasetID uuid.UUID) ([]*Collection, error) {
	response, err := c.service.ListCollections(ctx, datasetID)
	if err != nil {
		return nil, err
	}

	collections := make([]*Collection, len(response.GetData()))
	for i, collectionMessage := range response.GetData() {
		collections[i] = protoToCollection(collectionMessage)
	}

	return collections, nil
}

func protoToCollection(c *datasetsv1.CollectionInfo) *Collection {
	return &Collection{
		ID:           c.GetCollection().GetId().AsUUID(),
		Name:         c.GetCollection().GetName(),
		Availability: query.ProtoToTimeInterval(c.GetAvailability()),
		Count:        c.GetCount(),
	}
}
