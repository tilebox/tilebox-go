// Package datasets provides a client for interacting with Tilebox Datasets.
//
// Documentation: https://docs.tilebox.com/datasets
package datasets // import "github.com/tilebox/tilebox-go/datasets/v1"

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	datasetsv1 "github.com/tilebox/tilebox-go/protogen/datasets/v1"
)

// Dataset represents a Tilebox Dataset.
//
// Documentation: https://docs.tilebox.com/datasets/concepts/datasets
type Dataset struct {
	// ID is the unique identifier of the dataset.
	ID uuid.UUID
	// Type is the type of the dataset.
	Type *datasetsv1.AnnotatedType
	// Name is the name of the dataset.
	Name string
	// Description is a short description of the dataset.
	Description string
	// Slug is the unique slug of the dataset.
	Slug string
}

func (d Dataset) String() string {
	kind := ""
	switch d.Type.GetKind() {
	case datasetsv1.DatasetKind_DATASET_KIND_TEMPORAL:
		kind = "Temporal"
	case datasetsv1.DatasetKind_DATASET_KIND_SPATIOTEMPORAL:
		kind = "SpatioTemporal"
	case datasetsv1.DatasetKind_DATASET_KIND_UNSPECIFIED:
	}

	return fmt.Sprintf("%s [%s Dataset]: %s", d.Name, kind, d.Description)
}

type DatasetClient interface {
	// Get returns a dataset by its slug, e.g. "open_data.copernicus.sentinel1_sar".
	Get(ctx context.Context, slug string) (*Dataset, error)

	// List returns a list of all available datasets.
	List(ctx context.Context) ([]*Dataset, error)
}

var _ DatasetClient = &datasetClient{}

type datasetClient struct {
	service DatasetService
}

func (d datasetClient) Get(ctx context.Context, slug string) (*Dataset, error) {
	response, err := d.service.GetDataset(ctx, slug)
	if err != nil {
		return nil, err
	}

	return protoToDataset(response), nil
}

func (d datasetClient) List(ctx context.Context) ([]*Dataset, error) {
	response, err := d.service.ListDatasets(ctx)
	if err != nil {
		return nil, err
	}

	datasets := make([]*Dataset, len(response.GetDatasets()))
	for i, datasetMessage := range response.GetDatasets() {
		datasets[i] = protoToDataset(datasetMessage)
	}

	return datasets, nil
}

func protoToDataset(d *datasetsv1.Dataset) *Dataset {
	return &Dataset{
		ID:          d.GetId().AsUUID(),
		Type:        d.GetType(),
		Name:        d.GetName(),
		Description: d.GetSummary(),
		Slug:        d.GetSlug(),
	}
}
