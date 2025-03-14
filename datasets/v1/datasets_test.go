package datasets

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tilebox/tilebox-go/grpc"
)

const recordingDirectory = "testdata/recordings"

func NewRecordClient(tb testing.TB, filename string) (*Client, error) {
	err := os.MkdirAll(recordingDirectory, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("failed to create recording directory: %w", err)
	}

	file, err := os.Create(fmt.Sprintf("%s/%s.rpcs.bin", recordingDirectory, filename))
	if err != nil {
		return nil, fmt.Errorf("failed to create replay file: %w", err)
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
	), nil
}

func NewReplayClient(tb testing.TB, filename string) (*Client, error) {
	file, err := os.Open(fmt.Sprintf("%s/%s.rpcs.bin", recordingDirectory, filename))
	if err != nil {
		return nil, fmt.Errorf("failed to open replay file: %w", err)
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
	), nil
}

func TestClient_Datasets_List(t *testing.T) {
	ctx := context.Background()
	client, err := NewReplayClient(t, "datasets")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	datasets, err := client.Datasets.List(ctx)
	require.NoError(t, err)

	dataset := datasets[0]
	assert.Equal(t, "ERS SAR Granules", dataset.Name)
	assert.Equal(t, "49f17988-9f1c-446e-be2a-f949875b8274", dataset.ID.String())
}

func TestClient_Datasets_Get(t *testing.T) {
	ctx := context.Background()
	client, err := NewReplayClient(t, "dataset")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	dataset, err := client.Datasets.Get(ctx, "open_data.asf.ers_sar")
	require.NoError(t, err)

	assert.Equal(t, "ERS SAR Granules", dataset.Name)
	assert.Equal(t, "49f17988-9f1c-446e-be2a-f949875b8274", dataset.ID.String())
}
