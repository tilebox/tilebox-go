package datasets

import (
	"context"
	"iter"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/paulmach/orb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	datasetsv1 "github.com/tilebox/tilebox-go/protogen/go/datasets/v1"
	examplesv1 "github.com/tilebox/tilebox-go/protogen/go/examples/v1"
	tileboxv1 "github.com/tilebox/tilebox-go/protogen/go/tilebox/v1"
	"github.com/tilebox/tilebox-go/query"
	"google.golang.org/protobuf/proto"
)

func pointer[T any](x T) *T {
	return &x
}

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
	return datasetsv1.IngestResponse_builder{
		NumCreated: int64(len(datapoints)),
	}.Build(), nil
}

type mockDataAccessService struct {
	n int
	DataAccessService
}

func (m mockDataAccessService) Query(_ context.Context, _ []uuid.UUID, _ *datasetsv1.QueryFilters, _ *tileboxv1.Pagination, _ bool) (*datasetsv1.QueryResultPage, error) {
	data := make([][]byte, m.n)
	for i := range m.n {
		datapoint := examplesv1.Sentinel2Msi_builder{
			GranuleName:     pointer(uuid.New().String()),
			ProcessingLevel: pointer(datasetsv1.ProcessingLevel_PROCESSING_LEVEL_L1),
			Platform:        pointer("S2B"),
			FlightDirection: pointer(datasetsv1.FlightDirection_FLIGHT_DIRECTION_ASCENDING),
		}.Build()

		message, err := proto.Marshal(datapoint)
		if err != nil {
			return nil, err
		}

		data[i] = message
	}

	return datasetsv1.QueryResultPage_builder{
		Data: datasetsv1.RepeatedAny_builder{
			Value: data,
		}.Build(),
		NextPage: nil,
	}.Build(), nil
}

func Test_QueryOptions(t *testing.T) {
	now := time.Now()
	colorado := orb.Polygon{
		{{-109.05, 37.09}, {-102.06, 37.09}, {-102.06, 41.59}, {-109.05, 41.59}, {-109.05, 37.09}},
	}

	tests := []struct {
		name    string
		options []QueryOption
		want    queryOptions
	}{
		{
			name: "with temporal extent",
			options: []QueryOption{
				WithTemporalExtent(query.NewTimeInterval(
					now,
					now.Add(time.Hour),
				)),
			},
			want: queryOptions{
				temporalExtent: query.NewTimeInterval(
					now,
					now.Add(time.Hour),
				),
			},
		},
		{
			name: "with spatial extent",
			options: []QueryOption{
				WithSpatialExtent(colorado),
			},
			want: queryOptions{
				spatialExtent: &query.SpatialFilter{Geometry: colorado},
			},
		},
		{
			name: "with spatial extent filter",
			options: []QueryOption{
				WithSpatialExtentFilter(&query.SpatialFilter{Geometry: colorado, Mode: datasetsv1.SpatialFilterMode_SPATIAL_FILTER_MODE_INTERSECTS, CoordinateSystem: datasetsv1.SpatialCoordinateSystem_SPATIAL_COORDINATE_SYSTEM_CARTESIAN}),
			},
			want: queryOptions{
				spatialExtent: &query.SpatialFilter{Geometry: colorado, Mode: datasetsv1.SpatialFilterMode_SPATIAL_FILTER_MODE_INTERSECTS, CoordinateSystem: datasetsv1.SpatialCoordinateSystem_SPATIAL_COORDINATE_SYSTEM_CARTESIAN},
			},
		},
		{
			name: "with skip data",
			options: []QueryOption{
				WithSkipData(),
			},
			want: queryOptions{
				skipData: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var options queryOptions
			for _, option := range tt.options {
				option(&options)
			}
			assert.Equal(t, tt.want, options)
		})
	}
}

func Test_datapointClient_GetInto(t *testing.T) {
	ctx := context.Background()
	client := NewReplayClient(t, "datapoint_getinto")

	dataset, err := client.Datasets.Get(ctx, "open_data.copernicus.sentinel2_msi")
	require.NoError(t, err)

	collection, err := client.Collections.Get(ctx, dataset.ID, "S2B_S2MSI1C")
	require.NoError(t, err)

	datapointID := uuid.MustParse("01941f29-c650-202f-6495-c71dd2118fb1")

	t.Run("GetInto", func(t *testing.T) {
		var datapoint examplesv1.Sentinel2Msi
		err := client.Datapoints.GetInto(ctx, []uuid.UUID{collection.ID}, datapointID, &datapoint)
		require.NoError(t, err)

		assert.Equal(t, "01941f29-c650-202f-6495-c71dd2118fb1", uuid.Must(uuid.FromBytes(datapoint.GetId().GetUuid())).String())
		assert.Equal(t, "2025-01-01 00:00:19.024 +0000 UTC", datapoint.GetTime().AsTime().String())
		assert.Equal(t, "S2B_MSIL1C_20250101T000019_N0511_R073_T57QWV_20250101T010340.SAFE", datapoint.GetGranuleName())
	})

	t.Run("GetInto WithSkipData", func(t *testing.T) {
		var datapoint examplesv1.Sentinel2Msi
		err := client.Datapoints.GetInto(ctx, []uuid.UUID{collection.ID}, datapointID, &datapoint, WithSkipData())
		require.NoError(t, err)

		assert.Equal(t, "01941f29-c650-202f-6495-c71dd2118fb1", uuid.Must(uuid.FromBytes(datapoint.GetId().GetUuid())).String())
		assert.Empty(t, datapoint.GetGranuleName())
	})
}

func Test_datapointClient_QueryInto(t *testing.T) {
	ctx := context.Background()
	client := NewDatapointClient(10)

	collectionID := uuid.New()
	timeInterval := query.NewTimeInterval(time.Now(), time.Now())

	type args struct {
		collectionID uuid.UUID
		interval     query.TemporalExtent
		datapoints   any
	}
	tests := []struct {
		name    string
		args    args
		wantErr string
	}{
		{
			name: "QueryInto",
			args: args{
				collectionID: collectionID,
				interval:     timeInterval,
				datapoints:   &[]*examplesv1.Sentinel2Msi{},
			},
		},
		{
			name: "QueryInto nil",
			args: args{
				datapoints: nil,
			},
			wantErr: "datapoints must be a pointer, got <nil>",
		},
		{
			name: "QueryInto not a pointer",
			args: args{
				datapoints: collectionID,
			},
			wantErr: "datapoints must be a pointer, got uuid.UUID",
		},
		{
			name: "QueryInto not a slice",
			args: args{
				datapoints: &collectionID,
			},
			wantErr: "datapoints must be a pointer to a slice, got *uuid.UUID",
		},
		{
			name: "QueryInto slice wrong interface",
			args: args{
				datapoints: &[]context.Context{},
			},
			wantErr: "datapoints must be a pointer to a slice of proto.Message, got *[]context.Context",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.QueryInto(ctx, []uuid.UUID{tt.args.collectionID}, tt.args.datapoints, WithTemporalExtent(tt.args.interval))
			if tt.wantErr != "" {
				// we wanted an error, let's check if we got one
				require.Error(t, err, "expected an error, got none")
				assert.Contains(t, err.Error(), tt.wantErr, "error didn't contain expected message: '%s', got error '%s' instead.", tt.wantErr, err.Error())
				return
			}
			// we didn't want an error:
			require.NoError(t, err, "got an unexpected error")

			datapoints := *tt.args.datapoints.(*[]*examplesv1.Sentinel2Msi)

			assert.Len(t, datapoints, 10)
			assert.NotNil(t, datapoints[0])
		})
	}
}

// resultQueryInto is used to avoid the compiler optimizing away the benchmark output
var resultQueryInto []*examplesv1.Sentinel2Msi

// BenchmarkCollectAs benchmarks the QueryInto method
func Benchmark_QueryInto(b *testing.B) {
	ctx := context.Background()
	client := NewDatapointClient(1000)

	collectionID := uuid.New()
	timeInterval := query.NewTimeInterval(time.Now(), time.Now())

	var datapoints []*examplesv1.Sentinel2Msi
	b.Run("CollectAs", func(b *testing.B) {
		for range b.N {
			err := client.QueryInto(ctx, []uuid.UUID{collectionID}, &datapoints, WithTemporalExtent(timeInterval))
			require.NoError(b, err)
		}
	})
	resultQueryInto = datapoints
}

func Test_datapointClient_Query(t *testing.T) {
	ctx := context.Background()
	client := NewReplayClient(t, "load")

	dataset, err := client.Datasets.Get(ctx, "open_data.copernicus.sentinel2_msi")
	require.NoError(t, err)

	collection, err := client.Collections.Get(ctx, dataset.ID, "S2B_S2MSI1C")
	require.NoError(t, err)
	assert.Equal(t, "S2B_S2MSI1C", collection.Name)

	jan2025 := time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)
	timeInterval := query.NewTimeInterval(jan2025, jan2025.Add(1*time.Hour))

	t.Run("CollectAs", func(t *testing.T) {
		datapoints, err := CollectAs[*examplesv1.Sentinel2Msi](client.Datapoints.Query(ctx, []uuid.UUID{collection.ID}, WithTemporalExtent(timeInterval)))
		require.NoError(t, err)

		assert.Len(t, datapoints, 437)
		assert.Equal(t, "01941f29-c650-202f-6495-c71dd2118fb1", uuid.Must(uuid.FromBytes(datapoints[0].GetId().GetUuid())).String())
		assert.Equal(t, "2025-01-01 00:00:19.024 +0000 UTC", datapoints[0].GetTime().AsTime().String())
		assert.Equal(t, "S2B_MSIL1C_20250101T000019_N0511_R073_T57QWV_20250101T010340.SAFE", datapoints[0].GetGranuleName())
	})

	t.Run("CollectAs WithSkipData", func(t *testing.T) {
		datapoints, err := CollectAs[*examplesv1.Sentinel2Msi](client.Datapoints.Query(ctx, []uuid.UUID{collection.ID}, WithTemporalExtent(timeInterval), WithSkipData()))
		require.NoError(t, err)

		assert.Len(t, datapoints, 437)
		assert.Equal(t, "01941f29-c650-202f-6495-c71dd2118fb1", uuid.Must(uuid.FromBytes(datapoints[0].GetId().GetUuid())).String())
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
		datapoint := examplesv1.Sentinel2Msi_builder{
			GranuleName:     pointer(uuid.New().String()),
			ProcessingLevel: pointer(datasetsv1.ProcessingLevel_PROCESSING_LEVEL_L1),
			Platform:        pointer("S2B"),
			FlightDirection: pointer(datasetsv1.FlightDirection_FLIGHT_DIRECTION_ASCENDING),
		}.Build()

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

func (s *mockService) Query(_ context.Context, _ []uuid.UUID, _ ...QueryOption) iter.Seq2[[]byte, error] {
	return func(yield func([]byte, error) bool) {
		for _, data := range s.data {
			if !yield(data, nil) {
				return
			}
		}
	}
}

// result is used to avoid the compiler optimizing away the benchmark output
var result []*examplesv1.Sentinel2Msi

// BenchmarkCollectAs benchmarks the CollectAs function
// It is used to benchmark the cost of reflection and proto.Marshal inside CollectAs
func BenchmarkCollectAs(b *testing.B) {
	ctx := context.Background()
	collectionID := uuid.New()                    // dummy collection ID
	queryInterval := query.NewEmptyTimeInterval() // dummy time interval

	client := NewClient()
	client.Datapoints = NewMockDatapointClient(b, 1000)

	var r []*examplesv1.Sentinel2Msi // used to avoid the compiler optimizing the output
	b.Run("CollectAs", func(b *testing.B) {
		for range b.N {
			data := client.Datapoints.Query(ctx, []uuid.UUID{collectionID}, WithTemporalExtent(queryInterval))
			r, _ = CollectAs[*examplesv1.Sentinel2Msi](data)
		}
	})
	result = r

	b.Run("Marshal and no reflection", func(b *testing.B) {
		for range b.N {
			data := client.Datapoints.Query(ctx, []uuid.UUID{collectionID}, WithTemporalExtent(queryInterval))
			datapoints := make([]*examplesv1.Sentinel2Msi, 0)
			for datapoint, err := range data {
				if err != nil {
					b.Fatalf("failed to load datapoint: %v", err)
				}
				r := &examplesv1.Sentinel2Msi{}

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
			data := client.Datapoints.Query(ctx, []uuid.UUID{collectionID}, WithTemporalExtent(queryInterval))
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
				datapoints:    &[]*examplesv1.Sentinel2Msi{},
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
