package datasets

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	datasetsv1 "github.com/tilebox/tilebox-go/protogen/go/datasets/v1"
	"google.golang.org/protobuf/proto"
)

type mockService struct {
	meta []*datasetsv1.DatapointMetadata
	data [][]byte

	Service
}

func NewMockService(tb testing.TB, n int) Service {
	// generate some mock data
	meta := make([]*datasetsv1.DatapointMetadata, n)
	data := make([][]byte, n)
	for i := range n {
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
			tb.Fatalf("failed to marshal datapoint: %v", err)
		}

		data[i] = message
	}

	return &mockService{
		meta: meta,
		data: data,
	}
}

func (s *mockService) GetDatasetForInterval(_ context.Context, _ uuid.UUID, _ *datasetsv1.TimeInterval, _ *datasetsv1.DatapointInterval, _ *datasetsv1.Pagination, _ bool, _ bool) (*datasetsv1.Datapoints, error) {
	return &datasetsv1.Datapoints{
		Meta: s.meta,
		Data: &datasetsv1.RepeatedAny{
			Value: s.data,
		},
	}, nil
}

// result is used to avoid the compiler optimizing away the benchmark output
var result []*Datapoint[*datasetsv1.CopernicusDataspaceGranule]

// BenchmarkCollectAsLoad benchmarks the CollectAs + Load functions
// It is used to benchmark the cost of reflection and proto.Marshal inside CollectAs
func BenchmarkCollectAsLoad(b *testing.B) {
	ctx := context.Background()
	loadInterval := NewTimeInterval(time.Now(), time.Now(), false, false) // dummy load interval

	collection := &Collection{
		service: NewMockService(b, 1000),
	}

	var r []*Datapoint[*datasetsv1.CopernicusDataspaceGranule] // used to avoid the compiler optimizing the output
	b.Run("CollectAs", func(b *testing.B) {
		for range b.N {
			data := collection.Load(ctx, loadInterval, false, false)
			r, _ = CollectAs[*datasetsv1.CopernicusDataspaceGranule](data)
		}
	})
	result = r

	b.Run("Marshal and no reflection", func(b *testing.B) {
		for range b.N {
			data := collection.Load(ctx, loadInterval, false, false)
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
			data := collection.Load(ctx, loadInterval, false, false)
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
