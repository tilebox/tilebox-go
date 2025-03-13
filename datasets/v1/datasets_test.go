package datasets

import (
	"context"
	"fmt"
	"iter"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tilebox/tilebox-go/grpc"
	datasetsv1 "github.com/tilebox/tilebox-go/protogen/go/datasets/v1"
	"google.golang.org/protobuf/proto"
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

func TestClient_Collections_List(t *testing.T) {
	ctx := context.Background()
	client, err := NewReplayClient(t, "collections")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	dataset, err := client.Datasets.Get(ctx, "open_data.asf.ers_sar")
	require.NoError(t, err)

	collections, err := client.Collections.List(ctx, dataset.ID)
	require.NoError(t, err)

	names := lo.Map(collections, func(c *Collection, _ int) string {
		return c.Name
	})
	assert.Equal(t, []string{"ERS-1", "ERS-2"}, names)
}

func TestClient_Collections_Get(t *testing.T) {
	ctx := context.Background()
	client, err := NewReplayClient(t, "collection")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	dataset, err := client.Datasets.Get(ctx, "open_data.asf.ers_sar")
	require.NoError(t, err)

	collection, err := client.Collections.Get(ctx, dataset.ID, "ERS-2")
	require.NoError(t, err)

	assert.Equal(t, "ERS-2", collection.Name)
	assert.Equal(t, "c408f2b8-0488-4528-9fb7-a18361df639f", collection.ID.String())
	assert.Equal(t, "1995-10-01 03:13:03 +0000 UTC", collection.Availability.Start.String())
	assert.NotZero(t, collection.Count)
}

func TestClient_Data_Load(t *testing.T) {
	ctx := context.Background()
	client, err := NewReplayClient(t, "load")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	dataset, err := client.Datasets.Get(ctx, "open_data.asf.ers_sar")
	require.NoError(t, err)

	collection, err := client.Collections.Get(ctx, dataset.ID, "ERS-2")
	require.NoError(t, err)
	assert.Equal(t, "ERS-2", collection.Name)

	jan2000 := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
	interval := NewStandardTimeInterval(jan2000, jan2000.AddDate(0, 0, 7))

	t.Run("CollectAs", func(t *testing.T) {
		datapoints, err := CollectAs[*datasetsv1.ASFSarGranule](client.Datapoints.Load(ctx, collection.ID, interval))
		require.NoError(t, err)

		assert.Len(t, datapoints, 298)
		assert.Equal(t, "00dc7952-6c90-f40d-9b5e-3b72d4e790e0", datapoints[0].Meta.ID.String())
		assert.Equal(t, "2000-01-03 19:37:30 +0000 UTC", datapoints[0].Meta.EventTime.String())
		assert.Equal(t, "E2_24602_STD_F619", datapoints[0].Data.GetGranuleName())
	})

	t.Run("CollectAs WithSkipData", func(t *testing.T) {
		datapoints, err := CollectAs[*datasetsv1.ASFSarGranule](client.Datapoints.Load(ctx, collection.ID, interval, WithSkipData()))
		require.NoError(t, err)

		assert.Len(t, datapoints, 298)
		assert.Equal(t, "00dc7952-6c90-f40d-9b5e-3b72d4e790e0", datapoints[0].Meta.ID.String())
		assert.Empty(t, datapoints[0].Data.GetGranuleName())
	})

	t.Run("CollectAs WithSkipMeta", func(t *testing.T) {
		datapoints, err := CollectAs[*datasetsv1.ASFSarGranule](client.Datapoints.Load(ctx, collection.ID, interval, WithSkipMeta()))
		require.NoError(t, err)

		assert.Len(t, datapoints, 298)
		assert.Equal(t, "00dc7952-6c90-f40d-9b5e-3b72d4e790e0", datapoints[0].Meta.ID.String())
		assert.Empty(t, datapoints[0].Meta.EventTime)
		assert.Equal(t, "E2_24602_STD_F619", datapoints[0].Data.GetGranuleName())
	})

	t.Run("CollectAs WithSkipData WithSkipMeta", func(t *testing.T) {
		datapoints, err := CollectAs[*datasetsv1.ASFSarGranule](client.Datapoints.Load(ctx, collection.ID, interval, WithSkipData(), WithSkipMeta()))
		require.NoError(t, err)

		assert.Len(t, datapoints, 298)
		assert.Equal(t, "00dc7952-6c90-f40d-9b5e-3b72d4e790e0", datapoints[0].Meta.ID.String())
		assert.Empty(t, datapoints[0].Meta.EventTime)
		assert.Empty(t, datapoints[0].Data.GetGranuleName())
	})
}

type mockService struct {
	meta []*DatapointMetadata
	data [][]byte

	DatapointsClient
}

func NewMockDatapointsClient(tb testing.TB, n int) DatapointsClient {
	// generate some mock data
	meta := make([]*DatapointMetadata, n)
	data := make([][]byte, n)
	for i := range n {
		id := uuid.New()
		meta[i] = &DatapointMetadata{
			ID: id,
		}

		datapoint := &datasetsv1.CopernicusDataspaceGranule{
			GranuleName:     id.String(),
			ProcessingLevel: datasetsv1.ProcessingLevel_PROCESSING_LEVEL_L1,
			Satellite:       "Sentinel-1",
			FlightDirection: datasetsv1.FlightDirection_FLIGHT_DIRECTION_ASCENDING,
		}

		message, err := proto.Marshal(datapoint)
		if err != nil {
			tb.Fatalf("failed to marshal datapoint: %v", err)
		}

		data[i] = message
	}

	return &mockService{
		meta: meta,
		data: data,
	}
}

func (s *mockService) Load(_ context.Context, _ uuid.UUID, _ LoadInterval, _ ...LoadOption) iter.Seq2[*RawDatapoint, error] {
	return func(yield func(*RawDatapoint, error) bool) {
		for i := range s.meta {
			datapoint := &RawDatapoint{
				Meta: s.meta[i],
				Data: s.data[i],
			}
			if !yield(datapoint, nil) {
				return
			}
		}
	}
}

// result is used to avoid the compiler optimizing away the benchmark output
var result []*Datapoint[*datasetsv1.CopernicusDataspaceGranule]

// BenchmarkCollectAs benchmarks the CollectAs function
// It is used to benchmark the cost of reflection and proto.Marshal inside CollectAs
func BenchmarkCollectAs(b *testing.B) {
	ctx := context.Background()
	collectionID := uuid.New()             // dummy collection ID
	loadInterval := newEmptyTimeInterval() // dummy load interval

	client := NewClient()
	client.Datapoints = NewMockDatapointsClient(b, 1000)

	var r []*Datapoint[*datasetsv1.CopernicusDataspaceGranule] // used to avoid the compiler optimizing the output
	b.Run("CollectAs", func(b *testing.B) {
		for range b.N {
			data := client.Datapoints.Load(ctx, collectionID, loadInterval)
			r, _ = CollectAs[*datasetsv1.CopernicusDataspaceGranule](data)
		}
	})
	result = r

	b.Run("Marshal and no reflection", func(b *testing.B) {
		for range b.N {
			data := client.Datapoints.Load(ctx, collectionID, loadInterval)
			datapoints := make([]*datasetsv1.CopernicusDataspaceGranule, 0)
			for datapoint, err := range data {
				if err != nil {
					b.Fatalf("failed to load datapoint: %v", err)
				}
				r := &datasetsv1.CopernicusDataspaceGranule{}

				err = proto.Unmarshal(datapoint.Data, r)
				if err != nil {
					b.Fatalf("failed to unmarshal datapoint: %v", err)
				}
				datapoints = append(datapoints, r)
			}
		}
	})

	b.Run("No marshal and no reflection", func(b *testing.B) {
		for range b.N {
			data := client.Datapoints.Load(ctx, collectionID, loadInterval)
			datapoints := make([]*RawDatapoint, 0)
			for datapoint, err := range data {
				if err != nil {
					b.Fatalf("failed to load datapoint: %v", err)
				}
				datapoints = append(datapoints, datapoint)
			}
		}
	})
}
