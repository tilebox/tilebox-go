package datasets

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tilebox/tilebox-go/datasets/v1/field"
	"github.com/tilebox/tilebox-go/internal/grpc"
	datasetsv1 "github.com/tilebox/tilebox-go/protogen/datasets/v1"
	"pgregory.net/rapid"
)

const recordingDirectory = "testdata/recordings"

func NewRecordClient(tb testing.TB, filename string) *Client {
	err := os.MkdirAll(recordingDirectory, os.ModePerm)
	if err != nil {
		tb.Fatalf("failed to create recording directory: %v", err)
	}

	file, err := os.Create(fmt.Sprintf("%s/%s.rpcs.bin", recordingDirectory, filename))
	if err != nil {
		tb.Fatalf("failed to create replay file: %v", err)
	}
	tb.Cleanup(func() {
		_ = file.Close()
	})

	httpClient := &http.Client{
		Transport: grpc.NewRecordRoundTripper(file),
	}

	apiKey := os.Getenv("TILEBOX_OPENDATA_ONLY_API_KEY")
	if apiKey == "" {
		tb.Fatalf("TILEBOX_OPENDATA_ONLY_API_KEY is not set")
	}

	return NewClient(
		WithURL("https://api.tilebox.com"),
		WithAPIKey(apiKey),
		WithHTTPClient(httpClient),
	)
}

func NewReplayClient(tb testing.TB, filename string) *Client {
	file, err := os.Open(fmt.Sprintf("%s/%s.rpcs.bin", recordingDirectory, filename))
	if err != nil {
		tb.Fatalf("failed to open replay file: %v", err)
	}
	tb.Cleanup(func() {
		_ = file.Close()
	})

	httpClient := &http.Client{
		Transport: grpc.NewReplayRoundTripper(file),
	}

	return NewClient(
		WithURL("https://api.tilebox.com"), // url/key doesn't matter
		WithAPIKey("key"),
		WithHTTPClient(httpClient),
	)
}

func TestDataset_String(t *testing.T) {
	genDataset := rapid.Custom(func(t *rapid.T) *Dataset {
		return &Dataset{
			ID: uuid.New(),
			Type: datasetsv1.AnnotatedType_builder{
				Kind: rapid.OneOf(
					rapid.Just(datasetsv1.DatasetKind_DATASET_KIND_TEMPORAL),
					rapid.Just(datasetsv1.DatasetKind_DATASET_KIND_SPATIOTEMPORAL),
				).Draw(t, "Kind"),
			}.Build(),
			Name:        rapid.String().Draw(t, "Name"),
			Description: rapid.String().Draw(t, "Description"),
		}
	})

	rapid.Check(t, func(t *rapid.T) {
		input := genDataset.Draw(t, "dataset")
		got := input.String()

		assert.Contains(t, got, input.Name)
		assert.NotContains(t, got, input.ID.String())
		assert.Contains(t, got, input.Description)

		if input.Type.GetKind() == datasetsv1.DatasetKind_DATASET_KIND_TEMPORAL {
			assert.Contains(t, got, "Temporal")
		} else {
			assert.Contains(t, got, "SpatioTemporal")
		}
	})
}

func TestClient_Datasets_List(t *testing.T) {
	ctx := context.Background()
	client := NewReplayClient(t, "datasets")

	datasets, err := client.Datasets.List(ctx)
	require.NoError(t, err)

	dataset := datasets[0]
	assert.Equal(t, "ERS SAR Granules", dataset.Name)
	assert.Equal(t, "49f17988-9f1c-446e-be2a-f949875b8274", dataset.ID.String())
}

func TestClient_Datasets_Get(t *testing.T) {
	ctx := context.Background()
	client := NewReplayClient(t, "dataset")

	dataset, err := client.Datasets.Get(ctx, "open_data.asf.ers_sar")
	require.NoError(t, err)

	assert.Equal(t, "ERS SAR Granules", dataset.Name)
	assert.Equal(t, "49f17988-9f1c-446e-be2a-f949875b8274", dataset.ID.String())
}

func TestClient_Datasets_CreateOrUpdate(t *testing.T) {
	ctx := context.Background()
	client := NewReplayClient(t, "dataset_create_or_update")

	// create
	dataset, err := client.Datasets.CreateOrUpdate(ctx,
		KindSpatiotemporal,
		"Test Dataset",
		"test_dataset",
		[]Field{
			field.String("test_field").Description("Test field").ExampleValue("test value"),
		},
	)
	require.NoError(t, err)

	assert.Equal(t, "Test Dataset", dataset.Name)
	require.Len(t, dataset.Type.GetDescriptorSet().GetFile(), 1)
	require.Len(t, dataset.Type.GetDescriptorSet().GetFile()[0].GetMessageType(), 1)
	messageType := dataset.Type.GetDescriptorSet().GetFile()[0].GetMessageType()[0]
	assert.Equal(t, "TestDataset", messageType.GetName())
	require.Len(t, messageType.GetField(), 5)
	assert.Equal(t, "time", messageType.GetField()[0].GetName())
	assert.Equal(t, "id", messageType.GetField()[1].GetName())
	assert.Equal(t, "ingestion_time", messageType.GetField()[2].GetName())
	assert.Equal(t, "geometry", messageType.GetField()[3].GetName())
	assert.Equal(t, "test_field", messageType.GetField()[4].GetName())

	// update
	dataset, err = client.Datasets.CreateOrUpdate(ctx,
		KindSpatiotemporal,
		"Updated Test Dataset",
		"test_dataset",
		[]Field{
			field.String("test_field").Description("Test field").ExampleValue("test value"),
			field.Int64("updated_field").Description("Updated field").ExampleValue("12"),
		},
	)
	require.NoError(t, err)

	assert.Equal(t, "Updated Test Dataset", dataset.Name)
	require.Len(t, dataset.Type.GetDescriptorSet().GetFile(), 1)
	require.Len(t, dataset.Type.GetDescriptorSet().GetFile()[0].GetMessageType(), 1)
	messageType = dataset.Type.GetDescriptorSet().GetFile()[0].GetMessageType()[0]
	require.Len(t, messageType.GetField(), 6)
	assert.Equal(t, "time", messageType.GetField()[0].GetName())
	assert.Equal(t, "id", messageType.GetField()[1].GetName())
	assert.Equal(t, "ingestion_time", messageType.GetField()[2].GetName())
	assert.Equal(t, "geometry", messageType.GetField()[3].GetName())
	assert.Equal(t, "test_field", messageType.GetField()[4].GetName())
	assert.Equal(t, "updated_field", messageType.GetField()[5].GetName())
}
