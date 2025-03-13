package datasets

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	datasetsv1 "github.com/tilebox/tilebox-go/protogen/go/datasets/v1"
)

// Collection represents a Tilebox Dataset collection.
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
}

type CollectionClient interface {
	Create(ctx context.Context, datasetID uuid.UUID, collectionName string) (*Collection, error)
	Get(ctx context.Context, datasetID uuid.UUID, name string) (*Collection, error)
	List(ctx context.Context, datasetID uuid.UUID) ([]*Collection, error)
}

var _ CollectionClient = &collectionClient{}

type collectionClient struct {
	service CollectionService
}

// Create creates a new collection in the dataset with the given name.
func (c collectionClient) Create(ctx context.Context, datasetID uuid.UUID, collectionName string) (*Collection, error) {
	response, err := c.service.CreateCollection(ctx, datasetID, collectionName)
	if err != nil {
		return nil, err
	}

	collection, err := protoToCollection(response)
	if err != nil {
		return nil, fmt.Errorf("failed to convert collection from response: %w", err)
	}

	return collection, nil
}

// Get returns a collection by its name.
func (c collectionClient) Get(ctx context.Context, datasetID uuid.UUID, name string) (*Collection, error) {
	response, err := c.service.GetCollectionByName(ctx, datasetID, name)
	if err != nil {
		return nil, err
	}

	collection, err := protoToCollection(response)
	if err != nil {
		return nil, fmt.Errorf("failed to convert collection from response: %w", err)
	}

	return collection, nil
}

// List returns a list of all available collections in the dataset.
func (c collectionClient) List(ctx context.Context, datasetID uuid.UUID) ([]*Collection, error) {
	response, err := c.service.ListCollections(ctx, datasetID)
	if err != nil {
		return nil, err
	}

	collections := make([]*Collection, len(response.GetData()))
	for i, collectionMessage := range response.GetData() {
		collection, err := protoToCollection(collectionMessage)
		if err != nil {
			return nil, fmt.Errorf("failed to convert collection from response: %w", err)
		}

		collections[i] = collection
	}

	return collections, nil
}

func protoToCollection(c *datasetsv1.CollectionInfo) (*Collection, error) {
	id, err := uuid.Parse(c.GetCollection().GetId())
	if err != nil {
		return nil, fmt.Errorf("failed to parse collection id: %w", err)
	}

	return &Collection{
		ID:           id,
		Name:         c.GetCollection().GetName(),
		Availability: *protoToTimeInterval(c.GetAvailability()),
		Count:        c.GetCount(),
	}, nil
}
