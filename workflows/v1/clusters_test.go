package workflows

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

func Test_clusterClient_Get(t *testing.T) {
	ctx := context.Background()
	client, err := NewReplayClient(t, "cluster")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	cluster, err := client.Clusters.Get(ctx, "dev-cluster-ESugE7S4cwADVK")
	require.NoError(t, err)

	assert.Equal(t, "dev-cluster", cluster.Name)
	assert.Equal(t, "dev-cluster-ESugE7S4cwADVK", cluster.Slug)
}

func Test_clusterClient_List(t *testing.T) {
	ctx := context.Background()
	client, err := NewReplayClient(t, "clusters")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	clusters, err := client.Clusters.List(ctx)
	require.NoError(t, err)

	cluster := clusters[0]
	assert.Equal(t, "dev-cluster", cluster.Name)
	assert.Equal(t, "dev-cluster-ESugE7S4cwADVK", cluster.Slug)
}
