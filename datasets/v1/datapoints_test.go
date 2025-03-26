package datasets

import (
	"context"
	"iter"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "github.com/tilebox/tilebox-go/protogen-test/tilebox/v1"
	datasetsv1 "github.com/tilebox/tilebox-go/protogen/go/datasets/v1"
	"google.golang.org/protobuf/proto"
)

func NewDatapointClient(n int) DatapointClient {
	return &datapointClient{
		dataIngestionService: mockDataIngestionService{},
		dataAccessService:    mockDataAccessService{n: n},
	}
}

type mockDataIngestionService struct {
	DataIngestionService
}

func (m mockDataIngestionService) Ingest(_ context.Context, _ uuid.UUID, datapoints [][]byte, _ bool) (*datasetsv1.IngestResponse, error) {
	return &datasetsv1.IngestResponse{
		NumCreated: int64(len(datapoints)),
	}, nil
}

type mockDataAccessService struct {
	n int
	DataAccessService
}

func (m mockDataAccessService) Query(_ context.Context, _ []uuid.UUID, _ *datasetsv1.QueryFilters, _ *datasetsv1.Pagination, _ bool) (*datasetsv1.QueryResultPage, error) {
	data := make([][]byte, m.n)
	for i := range m.n {
		datapoint := &datasetsv1.CopernicusDataspaceGranule{
			GranuleName:     uuid.New().String(),
			ProcessingLevel: datasetsv1.ProcessingLevel_PROCESSING_LEVEL_L1,
			Satellite:       "Sentinel-1",
			FlightDirection: datasetsv1.FlightDirection_FLIGHT_DIRECTION_ASCENDING,
		}

		message, err := proto.Marshal(datapoint)
		if err != nil {
			return nil, err
		}

		data[i] = message
	}

	return &datasetsv1.QueryResultPage{
		Data: &datasetsv1.RepeatedAny{
			Value: data,
		},
		NextPage: nil,
	}, nil
}

func Test_datapointClient_LoadInto(t *testing.T) {
	ctx := context.Background()
	client := NewDatapointClient(10)

	collectionID := uuid.New()
	interval := NewStandardTimeInterval(time.Now(), time.Now())

	type args struct {
		collectionID uuid.UUID
		interval     LoadInterval
		datapoints   any
		options      []LoadOption
	}
	tests := []struct {
		name    string
		args    args
		wantErr string
	}{
		{
			name: "LoadInto",
			args: args{
				collectionID: collectionID,
				interval:     interval,
				datapoints:   &[]*datasetsv1.CopernicusDataspaceGranule{},
				options:      nil,
			},
		},
		{
			name: "LoadInto nil",
			args: args{
				datapoints: nil,
			},
			wantErr: "datapoints must be a pointer, got <nil>",
		},
		{
			name: "LoadInto not a pointer",
			args: args{
				datapoints: collectionID,
			},
			wantErr: "datapoints must be a pointer, got uuid.UUID",
		},
		{
			name: "LoadInto not a slice",
			args: args{
				datapoints: &collectionID,
			},
			wantErr: "datapoints must be a pointer to a slice, got *uuid.UUID",
		},
		{
			name: "LoadInto slice wrong interface",
			args: args{
				datapoints: &[]context.Context{},
			},
			wantErr: "datapoints must be a pointer to a slice of proto.Message, got *[]context.Context",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.LoadInto(ctx, tt.args.collectionID, tt.args.interval, tt.args.datapoints, tt.args.options...)
			if tt.wantErr != "" {
				// we wanted an error, let's check if we got one
				require.Error(t, err, "expected an error, got none")
				assert.Contains(t, err.Error(), tt.wantErr, "error didn't contain expected message: '%s', got error '%s' instead.", tt.wantErr, err.Error())
				return
			}
			// we didn't want an error:
			require.NoError(t, err, "got an unexpected error")

			datapoints := *tt.args.datapoints.(*[]*datasetsv1.CopernicusDataspaceGranule)

			assert.Len(t, datapoints, 10)
			assert.NotNil(t, datapoints[0])
		})
	}
}

// resultLoadInto is used to avoid the compiler optimizing away the benchmark output
var resultLoadInto []*datasetsv1.CopernicusDataspaceGranule

// BenchmarkCollectAs benchmarks the LoadInto method
func Benchmark_LoadInto(b *testing.B) {
	ctx := context.Background()
	client := NewDatapointClient(1000)

	collectionID := uuid.New()
	interval := NewStandardTimeInterval(time.Now(), time.Now())

	var datapoints []*datasetsv1.CopernicusDataspaceGranule
	b.Run("CollectAs", func(b *testing.B) {
		for range b.N {
			err := client.LoadInto(ctx, collectionID, interval, &datapoints)
			require.NoError(b, err)
		}
	})
	resultLoadInto = datapoints
}

func Test_datapointClient_Load(t *testing.T) {
	ctx := context.Background()
	client, err := NewReplayClient(t, "load")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	dataset, err := client.Datasets.Get(ctx, "tilebox.modis")
	require.NoError(t, err)

	collection, err := client.Collections.Get(ctx, dataset.ID, "MCD12Q1")
	require.NoError(t, err)
	assert.Equal(t, "MCD12Q1", collection.Name)

	jan2001 := time.Date(2001, time.January, 1, 0, 0, 0, 0, time.UTC)
	interval := NewStandardTimeInterval(jan2001, jan2001.AddDate(0, 0, 7))

	t.Run("CollectAs", func(t *testing.T) {
		datapoints, err := CollectAs[*v1.Modis](client.Datapoints.Load(ctx, collection.ID, interval))
		require.NoError(t, err)

		assert.Len(t, datapoints, 315)
		assert.Equal(t, "00e3c7a7-3400-00ad-770d-e7789458d06d", uuid.Must(uuid.FromBytes(datapoints[0].GetId().GetUuid())).String())
		assert.Equal(t, "2001-01-01 00:00:00 +0000 UTC", datapoints[0].GetTime().AsTime().String())
		assert.Equal(t, "MCD12Q1.A2001001.h13v12.061.2022146061358.hdf", datapoints[0].GetGranuleName())
	})

	t.Run("CollectAs WithSkipData", func(t *testing.T) {
		datapoints, err := CollectAs[*v1.Modis](client.Datapoints.Load(ctx, collection.ID, interval, WithSkipData()))
		require.NoError(t, err)

		assert.Len(t, datapoints, 315)
		assert.Equal(t, "00e3c7a7-3400-00ad-770d-e7789458d06d", uuid.Must(uuid.FromBytes(datapoints[0].GetId().GetUuid())).String())
		assert.Empty(t, datapoints[0].GetGranuleName())
	})
}

type mockService struct {
	data [][]byte

	DatapointClient
}

func NewMockDatapointClient(tb testing.TB, n int) DatapointClient {
	// generate some mock data
	data := make([][]byte, n)
	for i := range n {
		datapoint := &datasetsv1.CopernicusDataspaceGranule{
			GranuleName:     uuid.New().String(),
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
		data: data,
	}
}

func (s *mockService) Load(_ context.Context, _ uuid.UUID, _ LoadInterval, _ ...LoadOption) iter.Seq2[[]byte, error] {
	return func(yield func([]byte, error) bool) {
		for _, data := range s.data {
			if !yield(data, nil) {
				return
			}
		}
	}
}

// result is used to avoid the compiler optimizing away the benchmark output
var result []*datasetsv1.CopernicusDataspaceGranule

// BenchmarkCollectAs benchmarks the CollectAs function
// It is used to benchmark the cost of reflection and proto.Marshal inside CollectAs
func BenchmarkCollectAs(b *testing.B) {
	ctx := context.Background()
	collectionID := uuid.New()             // dummy collection ID
	loadInterval := newEmptyTimeInterval() // dummy load interval

	client := NewClient()
	client.Datapoints = NewMockDatapointClient(b, 1000)

	var r []*datasetsv1.CopernicusDataspaceGranule // used to avoid the compiler optimizing the output
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

				err = proto.Unmarshal(datapoint, r)
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
			datapoints := make([][]byte, 0)
			for datapoint, err := range data {
				if err != nil {
					b.Fatalf("failed to load datapoint: %v", err)
				}
				datapoints = append(datapoints, datapoint)
			}
		}
	})
}

func Test_datapointClient_Ingest(t *testing.T) {
	ctx := context.Background()
	client := NewDatapointClient(10)

	collectionID := uuid.New()

	type args struct {
		collectionID  uuid.UUID
		datapoints    any
		allowExisting bool
	}
	tests := []struct {
		name    string
		args    args
		want    *IngestResponse
		wantErr string
	}{
		{
			name: "Ingest",
			args: args{
				collectionID:  collectionID,
				datapoints:    &[]*datasetsv1.CopernicusDataspaceGranule{},
				allowExisting: false,
			},
		},
		{
			name: "Ingest nil",
			args: args{
				datapoints: nil,
			},
			wantErr: "datapoints must be a pointer, got <nil>",
		},
		{
			name: "Ingest not a pointer",
			args: args{
				datapoints: collectionID,
			},
			wantErr: "datapoints must be a pointer, got uuid.UUID",
		},
		{
			name: "Ingest not a slice",
			args: args{
				datapoints: &collectionID,
			},
			wantErr: "datapoints must be a pointer to a slice, got *uuid.UUID",
		},
		{
			name: "Ingest slice wrong interface",
			args: args{
				datapoints: &[]context.Context{},
			},
			wantErr: "datapoints must be a pointer to a slice of proto.Message, got *[]context.Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.Ingest(ctx, tt.args.collectionID, tt.args.datapoints, tt.args.allowExisting)
			if tt.wantErr != "" {
				// we wanted an error, let's check if we got one
				require.Error(t, err, "expected an error, got none")
				assert.Contains(t, err.Error(), tt.wantErr, "error didn't contain expected message: '%s', got error '%s' instead.", tt.wantErr, err.Error())
				return
			}
			// we didn't want an error:
			require.NoError(t, err, "got an unexpected error")

			assert.NotNil(t, got)
		})
	}
}
