package query

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	datasetsv1 "github.com/tilebox/tilebox-go/protogen/go/datasets/v1"
	"pgregory.net/rapid"
)

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
	got := timeInterval.ToProtoDatapointInterval()

	assert.Nil(t, got)
}

func TestTimeInterval_ToProtoTimeInterval(t *testing.T) {
	tests := []struct {
		name         string
		timeInterval *TimeInterval
		want         *datasetsv1.TimeInterval
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
	genTimeInterval := rapid.Custom(func(t *rapid.T) *TimeInterval {
		return &TimeInterval{
			Start:          time.Unix(rapid.Int64().Draw(t, "Start"), 0).UTC(),
			End:            time.Unix(rapid.Int64().Draw(t, "End"), 0).UTC(),
			StartExclusive: rapid.Bool().Draw(t, "StartExclusive"),
			EndInclusive:   rapid.Bool().Draw(t, "EndInclusive"),
		}
	})

	rapid.Check(t, func(t *rapid.T) {
		input := genTimeInterval.Draw(t, "time interval")

		got := ProtoToTimeInterval(input.ToProtoTimeInterval())

		assert.Equal(t, input.Start, got.Start)
		assert.Equal(t, input.End, got.End)
		assert.Equal(t, input.StartExclusive, got.StartExclusive)
		assert.Equal(t, input.EndInclusive, got.EndInclusive)
	})
}
