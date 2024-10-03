package datasets

import (
	"time"

	datasetsv1 "github.com/tilebox/tilebox-go/protogen/go/datasets/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const smallestPossibleTimeDelta = time.Microsecond

type LoadInterval interface {
	ToProtoTimeInterval() *datasetsv1.TimeInterval
	ToProtoDatapointInterval() *datasetsv1.DatapointInterval
}

var _ LoadInterval = &TimeInterval{}

type TimeInterval struct {
	start          time.Time
	end            time.Time
	startExclusive bool
	endInclusive   bool
}

func NewTimeInterval(start, end time.Time, startExclusive, endInclusive bool) *TimeInterval {
	return &TimeInterval{
		start:          start,
		end:            end,
		startExclusive: startExclusive,
		endInclusive:   endInclusive,
	}
}

func NewEmptyTimeInterval() *TimeInterval {
	return &TimeInterval{
		start:          time.Unix(0, 0),
		end:            time.Unix(0, 0),
		startExclusive: true,
		endInclusive:   false,
	}
}

func protoToTimeInterval(t *datasetsv1.TimeInterval) *TimeInterval {
	if t == nil {
		return NewEmptyTimeInterval()
	}

	return &TimeInterval{
		start:          t.GetStartTime().AsTime(),
		end:            t.GetEndTime().AsTime(),
		startExclusive: t.GetStartExclusive(),
		endInclusive:   t.GetEndInclusive(),
	}
}

func (t *TimeInterval) ToProtoTimeInterval() *datasetsv1.TimeInterval {
	return &datasetsv1.TimeInterval{
		StartTime:      timestamppb.New(t.start),
		EndTime:        timestamppb.New(t.end),
		StartExclusive: t.startExclusive,
		EndInclusive:   t.endInclusive,
	}
}

func (t *TimeInterval) ToProtoDatapointInterval() *datasetsv1.DatapointInterval {
	return nil
}

func (t *TimeInterval) ToHalfOpen() *TimeInterval {
	start := t.start
	if t.startExclusive {
		start = start.Add(smallestPossibleTimeDelta)
	}
	end := t.end
	if t.endInclusive {
		end = end.Add(smallestPossibleTimeDelta)
	}

	return &TimeInterval{
		start:          start,
		end:            end,
		startExclusive: false,
		endInclusive:   false,
	}
}

func (t *TimeInterval) Duration() time.Duration {
	return t.end.Sub(t.start)
}
