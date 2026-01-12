package query // import "github.com/tilebox/tilebox-go/query"

import (
	"fmt"
	"time"

	tileboxv1 "github.com/tilebox/tilebox-go/protogen/tilebox/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// smallestPossibleTimeDelta is the smallest possible time delta that can be represented by time.Duration.
const smallestPossibleTimeDelta = time.Nanosecond

// TemporalExtent is an interface for types that can be converted to a temporal extent, to be used in queries.
type TemporalExtent interface {
	ToProtoTimeInterval() *tileboxv1.TimeInterval
	ToProtoIDInterval() *tileboxv1.IDInterval
}

var _ TemporalExtent = &TimeInterval{}

// TimeInterval represents a time interval.
//
// Both the Start and End time can be exclusive or inclusive.
type TimeInterval struct {
	// Start is the start time of the interval.
	Start time.Time
	// End is the end time of the interval.
	End time.Time

	// We use exclusive for Start and inclusive for End, because that way when both are false,
	// we have a half-open interval [Start, End) which is the default behavior we want to achieve.
	StartExclusive bool
	EndInclusive   bool
}

// NewTimeInterval creates a new half-open TimeInterval which includes the start time but excludes the end time.
// In case you want to control the inclusivity of the start and end time, consider instantiating a TimeInterval struct
// directly.
func NewTimeInterval(start, end time.Time) *TimeInterval {
	return &TimeInterval{
		Start:          start,
		End:            end,
		StartExclusive: false,
		EndInclusive:   false,
	}
}

// NewPointInTime creates a new TimeInterval that represents a single point in time.
func NewPointInTime(t time.Time) *TimeInterval {
	return &TimeInterval{
		Start:          t,
		End:            t,
		StartExclusive: false,
		EndInclusive:   true,
	}
}

func NewEmptyTimeInterval() *TimeInterval {
	return &TimeInterval{
		Start:          time.Unix(0, 0).UTC(),
		End:            time.Unix(0, 0).UTC(),
		StartExclusive: true,
		EndInclusive:   false,
	}
}

func ProtoToTimeInterval(t *tileboxv1.TimeInterval) *TimeInterval {
	if t == nil {
		return NewEmptyTimeInterval()
	}

	return &TimeInterval{
		Start:          t.GetStartTime().AsTime(),
		End:            t.GetEndTime().AsTime(),
		StartExclusive: t.GetStartExclusive(),
		EndInclusive:   t.GetEndInclusive(),
	}
}

func (t *TimeInterval) ToProtoTimeInterval() *tileboxv1.TimeInterval {
	return tileboxv1.TimeInterval_builder{
		StartTime:      timestamppb.New(t.Start),
		EndTime:        timestamppb.New(t.End),
		StartExclusive: t.StartExclusive,
		EndInclusive:   t.EndInclusive,
	}.Build()
}

func (t *TimeInterval) ToProtoIDInterval() *tileboxv1.IDInterval {
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

func (t *TimeInterval) Equal(other *TimeInterval) bool {
	this := t.ToHalfOpen()
	other = other.ToHalfOpen()
	return this.Start.Equal(other.Start) && this.End.Equal(other.End)
}

func (t *TimeInterval) String() string {
	if t.Equal(NewEmptyTimeInterval()) {
		return "<empty>"
	}

	startCh := "["
	if t.StartExclusive {
		startCh = "("
	}

	endCh := ")"
	if t.EndInclusive {
		endCh = "]"
	}

	return fmt.Sprintf("%s%s, %s%s", startCh, t.Start, t.End, endCh)
}
