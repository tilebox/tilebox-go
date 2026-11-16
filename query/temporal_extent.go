package query

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
// A zero Start or End leaves that bound unspecified. Dataset queries resolve omitted
// bounds from the selected collections' availability. The zero-value interval leaves
// both bounds unspecified; use NewEmptyTimeInterval for an empty interval.
// Explicit bounds can be exclusive or inclusive. Flags on omitted bounds are ignored.
type TimeInterval struct {
	// Start is the start time of the interval, or zero to omit the lower bound.
	Start time.Time
	// End is the end time of the interval, or zero to omit the upper bound.
	End time.Time

	// We use exclusive for Start and inclusive for End, because that way when both are false,
	// we have a half-open interval [Start, End) which is the default behavior we want to achieve.
	StartExclusive bool
	EndInclusive   bool
}

// NewTimeInterval creates a new half-open TimeInterval which includes the start time but excludes the end time.
// A zero start or end leaves that bound unspecified.
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

// NewEmptyTimeInterval creates an interval containing no times, unlike the zero-value
// TimeInterval, which leaves both bounds unspecified.
func NewEmptyTimeInterval() *TimeInterval {
	return &TimeInterval{
		Start:          time.Unix(0, 0).UTC(),
		End:            time.Unix(0, 0).UTC(),
		StartExclusive: true,
		EndInclusive:   false,
	}
}

// ProtoToTimeInterval preserves omitted bounds in a non-nil message.
// A nil message represents an empty interval, as used for absent collection availability.
func ProtoToTimeInterval(t *tileboxv1.TimeInterval) *TimeInterval {
	if t == nil {
		return NewEmptyTimeInterval()
	}

	interval := &TimeInterval{
		StartExclusive: t.GetStartExclusive(),
		EndInclusive:   t.GetEndInclusive(),
	}
	if t.HasStartTime() {
		interval.Start = t.GetStartTime().AsTime()
	}
	if t.HasEndTime() {
		interval.End = t.GetEndTime().AsTime()
	}
	return interval
}

func (t *TimeInterval) ToProtoTimeInterval() *tileboxv1.TimeInterval {
	interval := tileboxv1.TimeInterval_builder{
		StartExclusive: t.StartExclusive,
		EndInclusive:   t.EndInclusive,
	}.Build()
	if !t.Start.IsZero() {
		interval.SetStartTime(timestamppb.New(t.Start))
	}
	if !t.End.IsZero() {
		interval.SetEndTime(timestamppb.New(t.End))
	}
	return interval
}

func (t *TimeInterval) ToProtoIDInterval() *tileboxv1.IDInterval {
	return nil
}

// ToHalfOpen converts the TimeInterval to a half-open interval [Start, End).
// Omitted bounds remain unspecified.
func (t *TimeInterval) ToHalfOpen() *TimeInterval {
	start := t.Start
	if !start.IsZero() && t.StartExclusive {
		start = start.Add(smallestPossibleTimeDelta)
	}
	end := t.End
	if !end.IsZero() && t.EndInclusive {
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

	start := t.Start.String()
	startCh := "["
	if t.Start.IsZero() {
		start = ".."
		startCh = "("
	} else if t.StartExclusive {
		startCh = "("
	}

	end := t.End.String()
	endCh := ")"
	if t.End.IsZero() {
		end = ".."
	} else if t.EndInclusive {
		endCh = "]"
	}

	return fmt.Sprintf("%s%s, %s%s", startCh, start, end, endCh)
}
