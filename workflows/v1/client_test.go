package workflows

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tilebox/tilebox-go/client"
	"github.com/tilebox/tilebox-go/internal/grpc"
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
		WithDisableTracing(),
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
		WithDisableTracing(),
	)
}

func TestClientMetadataOverride(t *testing.T) {
	metadata := client.Metadata{Name: "cli", Version: "v1.2.3"}
	cfg := newClientConfig([]ClientOption{WithClientMetadata(metadata)})

	require.Equal(t, metadata, cfg.clientMetadata)
}
