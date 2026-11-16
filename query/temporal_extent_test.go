package query

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tileboxv1 "github.com/tilebox/tilebox-go/protogen/tilebox/v1"
	"google.golang.org/protobuf/proto"
	"pgregory.net/rapid"
)

func TestTimeInterval_String(t *testing.T) {
	genTimeInterval := rapid.OneOf(
		rapid.Just(NewEmptyTimeInterval()),
		rapid.Custom(func(t *rapid.T) *TimeInterval {
			return &TimeInterval{
				Start:          time.Unix(rapid.Int64().Draw(t, "Start"), 0).UTC(),
				End:            time.Unix(rapid.Int64().Draw(t, "End"), 0).UTC(),
				StartExclusive: rapid.Bool().Draw(t, "StartExclusive"),
				EndInclusive:   rapid.Bool().Draw(t, "EndInclusive"),
			}
		}),
	)

	rapid.Check(t, func(t *rapid.T) {
		input := genTimeInterval.Draw(t, "time interval")
		got := input.String()

		if input.Equal(NewEmptyTimeInterval()) {
			assert.Contains(t, got, "empty")
			return
		}

		if input.StartExclusive {
			assert.True(t, strings.HasPrefix(got, "("))
		} else {
			assert.True(t, strings.HasPrefix(got, "["))
		}

		if input.EndInclusive {
			assert.True(t, strings.HasSuffix(got, "]"))
		} else {
			assert.True(t, strings.HasSuffix(got, ")"))
		}

		assert.Contains(t, got, input.Start.String())
		assert.Contains(t, got, input.End.String())
	})
}

func TestTimeInterval_Equals(t *testing.T) {
	now := time.Now()
	oneHourLater := now.Add(time.Hour)

	tests := []struct {
		name          string
		timeIntervalA *TimeInterval
		timeIntervalB *TimeInterval
		wantEqual     bool
	}{
		{
			name:          "interval should equal itself",
			timeIntervalA: NewTimeInterval(now, oneHourLater),
			timeIntervalB: NewTimeInterval(now, oneHourLater),
			wantEqual:     true,
		},
		{
			name:          "closed interval should equal its half open counterpart",
			timeIntervalA: &TimeInterval{Start: now, End: oneHourLater.Add(-smallestPossibleTimeDelta), StartExclusive: false, EndInclusive: true},
			timeIntervalB: NewTimeInterval(now, oneHourLater),
			wantEqual:     true,
		},
		{
			name:          "open interval should equal its half open counterpart",
			timeIntervalA: &TimeInterval{Start: now.Add(-smallestPossibleTimeDelta), End: oneHourLater, StartExclusive: true, EndInclusive: false},
			timeIntervalB: NewTimeInterval(now, oneHourLater),
			wantEqual:     true,
		},
		{
			name:          "open closed interval should equal its half open counterpart",
			timeIntervalA: &TimeInterval{Start: now.Add(-smallestPossibleTimeDelta), End: oneHourLater.Add(-smallestPossibleTimeDelta), StartExclusive: true, EndInclusive: true},
			timeIntervalB: NewTimeInterval(now, oneHourLater),
			wantEqual:     true,
		},
		{
			name:          "NewPointInTime",
			timeIntervalA: &TimeInterval{Start: now, End: now, StartExclusive: false, EndInclusive: true},
			timeIntervalB: NewPointInTime(now),
			wantEqual:     true,
		},
		{
			name:          "missing bounds ignore inclusion flags",
			timeIntervalA: &TimeInterval{StartExclusive: true, EndInclusive: true},
			timeIntervalB: &TimeInterval{},
			wantEqual:     true,
		},
		{
			name:          "unbounded is not empty",
			timeIntervalA: &TimeInterval{},
			timeIntervalB: NewEmptyTimeInterval(),
		},
		{
			name:          "missing start is not the epoch",
			timeIntervalA: &TimeInterval{End: now},
			timeIntervalB: NewTimeInterval(time.Unix(0, 0).UTC(), now),
		},
		{
			name:          "missing end is not the epoch",
			timeIntervalA: &TimeInterval{Start: now},
			timeIntervalB: NewTimeInterval(now, time.Unix(0, 0).UTC()),
		},
		{
			name:          "start-only exclusive equals shifted inclusive start",
			timeIntervalA: &TimeInterval{Start: now, StartExclusive: true, EndInclusive: true},
			timeIntervalB: &TimeInterval{Start: now.Add(time.Nanosecond)},
			wantEqual:     true,
		},
		{
			name:          "end-only inclusive equals shifted exclusive end",
			timeIntervalA: &TimeInterval{End: now, StartExclusive: true, EndInclusive: true},
			timeIntervalB: &TimeInterval{End: now.Add(time.Nanosecond)},
			wantEqual:     true,
		},
		{
			name:          "explicit start inclusion still matters",
			timeIntervalA: &TimeInterval{Start: now, StartExclusive: true},
			timeIntervalB: &TimeInterval{Start: now},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantEqual, tt.timeIntervalA.Equal(tt.timeIntervalB))
		})
	}
}

func TestTimeInterval_ToHalfOpen(t *testing.T) {
	now := time.Now()
	oneHourLater := now.Add(time.Hour)

	tests := []struct {
		name         string
		timeInterval *TimeInterval
	}{
		{
			name: "ToHalfOpen no-op",
			timeInterval: &TimeInterval{
				Start:          now,
				End:            oneHourLater,
				StartExclusive: false,
				EndInclusive:   false,
			},
		},
		{
			name: "Left half open interval ToHalfOpen",
			timeInterval: &TimeInterval{
				Start:          now,
				End:            oneHourLater,
				StartExclusive: true,
				EndInclusive:   false,
			},
		},
		{
			name: "Open interval to ToHalfOpen",
			timeInterval: &TimeInterval{
				Start:          now,
				End:            oneHourLater,
				StartExclusive: false,
				EndInclusive:   true,
			},
		},
		{
			name: "Closed interval ToHalfOpen",
			timeInterval: &TimeInterval{
				Start:          now,
				End:            oneHourLater,
				StartExclusive: true,
				EndInclusive:   true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.timeInterval.ToHalfOpen()

			assert.False(t, got.StartExclusive)
			assert.False(t, got.EndInclusive)
		})
	}
}

func TestTimeInterval_ToProtoDatapointInterval(t *testing.T) {
	timeInterval := NewTimeInterval(time.Now(), time.Now().Add(time.Hour))
	got := timeInterval.ToProtoIDInterval()

	assert.Nil(t, got)
}

func TestTimeInterval_ToProtoTimeInterval(t *testing.T) {
	tests := []struct {
		name         string
		timeInterval *TimeInterval
		want         *tileboxv1.TimeInterval
	}{
		{
			name:         "ToProtoTimeInterval",
			timeInterval: NewTimeInterval(time.Now(), time.Now().Add(time.Hour)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.timeInterval.ToProtoTimeInterval()

			assert.True(t, got.GetStartTime().AsTime().Equal(tt.timeInterval.Start))
			assert.True(t, got.GetEndTime().AsTime().Equal(tt.timeInterval.End))
			assert.Equal(t, tt.timeInterval.StartExclusive, got.GetStartExclusive())
			assert.Equal(t, tt.timeInterval.EndInclusive, got.GetEndInclusive())
		})
	}
}

func Test_protoToTimeIntervalRoundtrip(t *testing.T) {
	genTime := rapid.OneOf(
		rapid.Just(time.Time{}),
		rapid.Just(time.Unix(0, 0).UTC()),
		rapid.Custom(func(t *rapid.T) time.Time {
			return time.Unix(rapid.Int64Range(-2208988800, 4102444800).Draw(t, "seconds"), rapid.Int64Range(0, 999999999).Draw(t, "nanoseconds")).UTC()
		}),
	)
	genTimeInterval := rapid.Custom(func(t *rapid.T) *TimeInterval {
		return &TimeInterval{
			Start:          genTime.Draw(t, "Start"),
			End:            genTime.Draw(t, "End"),
			StartExclusive: rapid.Bool().Draw(t, "StartExclusive"),
			EndInclusive:   rapid.Bool().Draw(t, "EndInclusive"),
		}
	})

	rapid.Check(t, func(t *rapid.T) {
		input := genTimeInterval.Draw(t, "time interval")

		encoded, err := proto.Marshal(input.ToProtoTimeInterval())
		require.NoError(t, err)
		wire := &tileboxv1.TimeInterval{}
		require.NoError(t, proto.Unmarshal(encoded, wire))
		assert.Equal(t, !input.Start.IsZero(), wire.HasStartTime())
		assert.Equal(t, !input.End.IsZero(), wire.HasEndTime())
		got := ProtoToTimeInterval(wire)

		assert.Equal(t, input.Start, got.Start)
		assert.Equal(t, input.End, got.End)
		assert.Equal(t, input.StartExclusive, got.StartExclusive)
		assert.Equal(t, input.EndInclusive, got.EndInclusive)
	})
}

func TestTimeInterval_OpenBounds(t *testing.T) {
	bound := time.Date(2026, time.May, 3, 12, 30, 0, 0, time.UTC)
	for _, tc := range []struct {
		name     string
		interval *TimeInterval
		halfOpen *TimeInterval
		text     string
	}{
		{"both omitted", &TimeInterval{}, &TimeInterval{}, "(.., ..)"},
		{"omitted flags ignored", &TimeInterval{StartExclusive: true, EndInclusive: true}, &TimeInterval{}, "(.., ..)"},
		{"start only", &TimeInterval{Start: bound}, &TimeInterval{Start: bound}, "[2026-05-03 12:30:00 +0000 UTC, ..)"},
		{"exclusive start", &TimeInterval{Start: bound, StartExclusive: true, EndInclusive: true}, &TimeInterval{Start: bound.Add(time.Nanosecond)}, "(2026-05-03 12:30:00 +0000 UTC, ..)"},
		{"end only", &TimeInterval{End: bound}, &TimeInterval{End: bound}, "(.., 2026-05-03 12:30:00 +0000 UTC)"},
		{"inclusive end", &TimeInterval{End: bound, StartExclusive: true, EndInclusive: true}, &TimeInterval{End: bound.Add(time.Nanosecond)}, "(.., 2026-05-03 12:30:00 +0000 UTC]"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.halfOpen, tc.interval.ToHalfOpen())
			assert.Equal(t, tc.text, tc.interval.String())
		})
	}

	assert.Equal(t, NewEmptyTimeInterval(), ProtoToTimeInterval(nil))
	assert.Equal(t, &TimeInterval{}, ProtoToTimeInterval(&tileboxv1.TimeInterval{}))
	assert.Equal(t, "<empty>", NewEmptyTimeInterval().String())
}
