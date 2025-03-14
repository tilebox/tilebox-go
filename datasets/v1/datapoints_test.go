package datasets

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	datasetsv1 "github.com/tilebox/tilebox-go/protogen/go/datasets/v1"
	"google.golang.org/protobuf/proto"
)

func NewDatapointClient(n int) (DatapointsClient, error) {
	return &datapointClient{
		dataIngestionService: mockDataIngestionService{},
		dataAccessService:    mockDataAccessService{n: n},
	}, nil
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
	client, err := NewDatapointClient(10)
	require.NoError(t, err)

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
	client, err := NewDatapointClient(1000)
	require.NoError(b, err)

	collectionID := uuid.New()
	interval := NewStandardTimeInterval(time.Now(), time.Now())

	var datapoints []*TypedDatapoint[*datasetsv1.CopernicusDataspaceGranule]
	b.Run("CollectAs", func(b *testing.B) {
		for range b.N {
			err = client.LoadInto(ctx, collectionID, interval, &datapoints)
			require.NoError(b, err)
		}
	})
	resultLoadInto = datapoints
}
