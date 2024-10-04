package datasets // import "github.com/tilebox/tilebox-go/datasets/v1"

import (
	"time"

	datasetsv1 "github.com/tilebox/tilebox-go/protogen/go/datasets/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const smallestPossibleTimeDelta = time.Microsecond

// LoadInterval is an interface for loading data from a collection.
type LoadInterval interface {
	ToProtoTimeInterval() *datasetsv1.TimeInterval
	ToProtoDatapointInterval() *datasetsv1.DatapointInterval
}

var _ LoadInterval = &TimeInterval{}

// TimeInterval represents a time interval.
//
// Both the start and end time can be exclusive or inclusive.
type TimeInterval struct {
	// start is the start time of the interval.
	start time.Time
	// end is the end time of the interval.
	end time.Time

	// We use exclusive for start and inclusive for end, because that way when both are false
	// we have a half-open interval [start, end) which is the default behaviour we want to achieve.
	startExclusive bool
	endInclusive   bool
}

// NewTimeInterval creates a new TimeInterval.
func NewTimeInterval(start, end time.Time, startExclusive, endInclusive bool) *TimeInterval {
	return &TimeInterval{
		start:          start,
		end:            end,
		startExclusive: startExclusive,
		endInclusive:   endInclusive,
	}
}

// NewEmptyTimeInterval creates a new TimeInterval that represents an empty interval.
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

// ToHalfOpen converts the TimeInterval to a half-open interval [start, end).
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

// Duration returns the duration between the start and end times of the TimeInterval.
func (t *TimeInterval) Duration() time.Duration {
	return t.end.Sub(t.start)
}
