package datasets

import (
	"time"

	datasetsv1 "github.com/tilebox/tilebox-go/protogen/go/datasets/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type LoadInterval interface {
	ToProtoTimeInterval() *datasetsv1.TimeInterval
	ToProtoDatapointInterval() *datasetsv1.DatapointInterval
}

var _ LoadInterval = &TimeInterval{}

type TimeInterval struct {
	start time.Time
	end   time.Time
}

func NewTimeInterval(start, end time.Time) *TimeInterval {
	return &TimeInterval{
		start: start,
		end:   end,
	}
}

func protoToTimeInterval(t *datasetsv1.TimeInterval) *TimeInterval {
	return &TimeInterval{
		start: t.GetStartTime().AsTime(),
		end:   t.GetEndTime().AsTime(),
	}
}

func (t *TimeInterval) ToProtoTimeInterval() *datasetsv1.TimeInterval {
	return &datasetsv1.TimeInterval{
		StartTime:      timestamppb.New(t.start),
		EndTime:        timestamppb.New(t.end),
		StartExclusive: false,
		EndInclusive:   true,
	}
}

func (t *TimeInterval) ToProtoDatapointInterval() *datasetsv1.DatapointInterval {
	return nil
}
