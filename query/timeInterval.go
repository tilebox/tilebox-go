package query // import "github.com/tilebox/tilebox-go/query"

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
	datasetsv1 "github.com/tilebox/tilebox-go/protogen/go/datasets/v1"
	workflowsv1 "github.com/tilebox/tilebox-go/protogen/go/workflows/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// smallestPossibleTimeDelta is the smallest possible time delta that can be represented by time.Duration.
const smallestPossibleTimeDelta = time.Nanosecond

// TemporalExtent is an interface for querying resources and filtering them by a time range.
type TemporalExtent interface {
	ToProtoTimeInterval() *datasetsv1.TimeInterval
	ToProtoDatapointInterval() *datasetsv1.DatapointInterval
	ToProtoIDInterval() *workflowsv1.IDInterval
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

	// We use exclusive for Start and inclusive for End, because that way when both are false
	// we have a half-open interval [Start, End) which is the default behaviour we want to achieve.
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

func NewEmptyTimeInterval() *TimeInterval {
	return &TimeInterval{
		Start:          time.Unix(0, 0),
		End:            time.Unix(0, 0),
		StartExclusive: true,
		EndInclusive:   false,
	}
}

func ProtoToTimeInterval(t *datasetsv1.TimeInterval) *TimeInterval {
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

func (t *TimeInterval) ToProtoIDInterval() *workflowsv1.IDInterval {
	startID, err := newUUIDWithTime(t.Start)
	if err != nil {
		startID = uuid.Nil
	}

	endID, err := newUUIDWithTime(t.End)
	if err != nil {
		endID = uuid.Nil
	}

	return &workflowsv1.IDInterval{
		StartId:        &workflowsv1.UUID{Uuid: startID[:]},
		EndId:          &workflowsv1.UUID{Uuid: endID[:]},
		StartExclusive: t.StartExclusive,
		EndInclusive:   t.EndInclusive,
	}
}

// newUUIDWithTime generates a new uuid.UUID for a given time and random entropy
func newUUIDWithTime(t time.Time) (uuid.UUID, error) {
	uu, err := ulid.New(ulid.Timestamp(t), rand.Reader)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create new ulid: %w", err)
	}
	return uuid.UUID(uu), nil
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
