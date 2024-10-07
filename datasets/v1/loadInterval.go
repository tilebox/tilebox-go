package datasets // import "github.com/tilebox/tilebox-go/datasets/v1"

import (
	"time"

	datasetsv1 "github.com/tilebox/tilebox-go/protogen/go/datasets/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// smallestPossibleTimeDelta is the smallest possible time delta that can be represented by time.Duration.
const smallestPossibleTimeDelta = time.Nanosecond

// LoadInterval is an interface for loading data from a collection.
type LoadInterval interface {
	ToProtoTimeInterval() *datasetsv1.TimeInterval
	ToProtoDatapointInterval() *datasetsv1.DatapointInterval
}

var _ LoadInterval = &TimeInterval{}

// TimeInterval represents a time interval.
//
// Both the Start and End time can be exclusive or inclusive.
type TimeInterval struct {
	// Start is the start time of the interval.
	Start time.Time
	// End is the end time of the interval.
	End time.Time

	// We use exclusive for Start and inclusive for End, because that way when both are false
	// we have a half-open interval [Start, End) which is the default behaviour we want to achieve.
	StartExclusive bool
	EndInclusive   bool
}

// NewTimeInterval creates a new TimeInterval.
func NewTimeInterval(start, end time.Time, startExclusive, endInclusive bool) *TimeInterval {
	return &TimeInterval{
		Start:          start,
		End:            end,
		StartExclusive: startExclusive,
		EndInclusive:   endInclusive,
	}
}

// NewStandardTimeInterval creates a new TimeInterval that represents a standard interval with inclusive start and exclusive end.
func NewStandardTimeInterval(start, end time.Time) *TimeInterval {
	return NewTimeInterval(start, end, false, false)
}

func newEmptyTimeInterval() *TimeInterval {
	return &TimeInterval{
		Start:          time.Unix(0, 0),
		End:            time.Unix(0, 0),
		StartExclusive: true,
		EndInclusive:   false,
	}
}

func protoToTimeInterval(t *datasetsv1.TimeInterval) *TimeInterval {
	if t == nil {
		return newEmptyTimeInterval()
	}

	return &TimeInterval{
		Start:          t.GetStartTime().AsTime(),
		End:            t.GetEndTime().AsTime(),
		StartExclusive: t.GetStartExclusive(),
		EndInclusive:   t.GetEndInclusive(),
	}
}

func (t *TimeInterval) ToProtoTimeInterval() *datasetsv1.TimeInterval {
	return &datasetsv1.TimeInterval{
		StartTime:      timestamppb.New(t.Start),
		EndTime:        timestamppb.New(t.End),
		StartExclusive: t.StartExclusive,
		EndInclusive:   t.EndInclusive,
	}
}

func (t *TimeInterval) ToProtoDatapointInterval() *datasetsv1.DatapointInterval {
	return nil
}

// ToHalfOpen converts the TimeInterval to a half-open interval [Start, End).
func (t *TimeInterval) ToHalfOpen() *TimeInterval {
	start := t.Start
	if t.StartExclusive {
		start = start.Add(smallestPossibleTimeDelta)
	}
	end := t.End
	if t.EndInclusive {
		end = end.Add(smallestPossibleTimeDelta)
	}

	return &TimeInterval{
		Start:          start,
		End:            end,
		StartExclusive: false,
		EndInclusive:   false,
	}
}

// Duration returns the duration between the start and end times of the TimeInterval.
func (t *TimeInterval) Duration() time.Duration {
	return t.End.Sub(t.Start)
}
