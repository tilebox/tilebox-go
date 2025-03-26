// Package datasets provides a client for interacting with Tilebox Datasets.
//
// Documentation: https://docs.tilebox.com/datasets
package datasets // import "github.com/tilebox/tilebox-go/datasets/v1"

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	datasetsv1 "github.com/tilebox/tilebox-go/protogen/go/datasets/v1"
)

// Dataset represents a Tilebox Time Series Dataset.
//
// Documentation: https://docs.tilebox.com/datasets/timeseries
type Dataset struct {
	// ID is the unique identifier of the dataset.
	ID uuid.UUID
	// Type is the type of the dataset.
	Type *datasetsv1.AnnotatedType
	// Name is the name of the dataset.
	Name string
	// Summary is a summary of the purpose of the dataset.
	Summary string
}

type DatasetClient interface {
	Get(ctx context.Context, slug string) (*Dataset, error)
	List(ctx context.Context) ([]*Dataset, error)
}

var _ DatasetClient = &datasetClient{}

type datasetClient struct {
	service DatasetService
}

// Get returns a dataset by its slug, e.g. "open_data.copernicus.sentinel1_sar".
func (d datasetClient) Get(ctx context.Context, slug string) (*Dataset, error) {
	response, err := d.service.GetDataset(ctx, slug)
	if err != nil {
		return nil, err
	}

	dataset, err := protoToDataset(response)
	if err != nil {
		return nil, fmt.Errorf("failed to convert dataset from response: %w", err)
	}

	return dataset, nil
}

// List returns a list of all available datasets.
func (d datasetClient) List(ctx context.Context) ([]*Dataset, error) {
	response, err := d.service.ListDatasets(ctx)
	if err != nil {
		return nil, err
	}

	datasets := make([]*Dataset, len(response.GetDatasets()))
	for i, datasetMessage := range response.GetDatasets() {
		dataset, err := protoToDataset(datasetMessage)
		if err != nil {
			return nil, fmt.Errorf("failed to convert dataset from response: %w", err)
		}

		datasets[i] = dataset
	}

	return datasets, nil
}

func protoToDataset(d *datasetsv1.Dataset) (*Dataset, error) {
	id, err := protoToUUID(d.GetId())
	if err != nil {
		return nil, fmt.Errorf("failed to convert dataset id to uuid: %w", err)
	}

	return &Dataset{
		ID:      id,
		Type:    d.GetType(),
		Name:    d.GetName(),
		Summary: d.GetSummary(),
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

func uuidToProtobuf(id uuid.UUID) *datasetsv1.ID {
	if id == uuid.Nil {
		return nil
	}

	return &datasetsv1.ID{Uuid: id[:]}
}

func uuidsToProtobuf(ids []uuid.UUID) []*datasetsv1.ID {
	pbIDs := make([]*datasetsv1.ID, 0, len(ids))
	for _, id := range ids {
		pbIDs = append(pbIDs, uuidToProtobuf(id))
	}
	return pbIDs
}
