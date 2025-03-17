package datasets

import (
	"context"
	"iter"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

type mockDataAccessService struct {
	n int
	DataAccessService
}

func (m mockDataAccessService) GetDatasetForInterval(_ context.Context, _ uuid.UUID, _ *datasetsv1.TimeInterval, _ *datasetsv1.DatapointInterval, _ *datasetsv1.Pagination, _ bool, _ bool) (*datasetsv1.DatapointPage, error) {
	meta := make([]*datasetsv1.DatapointMetadata, m.n)
	data := make([][]byte, m.n)
	for i := range m.n {
		id := uuid.New().String()
		meta[i] = &datasetsv1.DatapointMetadata{
			Id: &id,
		}

		datapoint := &datasetsv1.CopernicusDataspaceGranule{
			GranuleName:     id,
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

	return &datasetsv1.DatapointPage{
		Meta: meta,
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
				datapoints:   &[]*TypedDatapoint[*datasetsv1.CopernicusDataspaceGranule]{},
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
			name: "LoadInto slice wrong type and no pointer",
			args: args{
				datapoints: &[]uuid.UUID{},
			},
			wantErr: "datapoints must be a pointer to a slice of *TypedDatapoint, got *[]uuid.UUID",
		},
		{
			name: "LoadInto slice pointer but wrong type",
			args: args{
				datapoints: &[]*uuid.UUID{},
			},
			wantErr: "datapoints must be a pointer to a slice of *TypedDatapoint, got *[]*uuid.UUID",
		},
		{
			name: "LoadInto slice good type but not pointer",
			args: args{
				datapoints: &[]TypedDatapoint[*datasetsv1.CopernicusDataspaceGranule]{},
			},
			wantErr: "datapoints must be a pointer to a slice of *TypedDatapoint, got *[]datasets.TypedDatapoint[*github.com/tilebox/tilebox-go/protogen/go/datasets/v1.CopernicusDataspaceGranule]",
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

			datapoints := *tt.args.datapoints.(*[]*TypedDatapoint[*datasetsv1.CopernicusDataspaceGranule])

			assert.Len(t, datapoints, 10)
			assert.NotNil(t, datapoints[0].Meta)
			assert.NotNil(t, datapoints[0].Data)
		})
	}
}

// resultLoadInto is used to avoid the compiler optimizing away the benchmark output
var resultLoadInto []*TypedDatapoint[*datasetsv1.CopernicusDataspaceGranule]

// BenchmarkCollectAs benchmarks the LoadInto method
func Benchmark_LoadInto(b *testing.B) {
	ctx := context.Background()
	client := NewDatapointClient(1000)

	collectionID := uuid.New()
	interval := NewStandardTimeInterval(time.Now(), time.Now())

	var datapoints []*TypedDatapoint[*datasetsv1.CopernicusDataspaceGranule]
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

	DatapointClient
}

func NewMockDatapointClient(tb testing.TB, n int) DatapointClient {
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
var result []*TypedDatapoint[*datasetsv1.CopernicusDataspaceGranule]

// BenchmarkCollectAs benchmarks the CollectAs function
// It is used to benchmark the cost of reflection and proto.Marshal inside CollectAs
func BenchmarkCollectAs(b *testing.B) {
	ctx := context.Background()
	collectionID := uuid.New()             // dummy collection ID
	loadInterval := newEmptyTimeInterval() // dummy load interval

	client := NewClient()
	client.Datapoints = NewMockDatapointClient(b, 1000)

	var r []*TypedDatapoint[*datasetsv1.CopernicusDataspaceGranule] // used to avoid the compiler optimizing the output
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
