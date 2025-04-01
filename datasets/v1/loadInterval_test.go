package datasets

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	datasetsv1 "github.com/tilebox/tilebox-go/protogen/go/datasets/v1"
	"pgregory.net/rapid"
)

func TestNewEmptyTimeInterval(t *testing.T) {
	got := newEmptyTimeInterval()

	assert.Equal(t, time.Duration(0), got.Duration())
}

func TestTimeInterval_Duration(t *testing.T) {
	now := time.Now()
	fiveMinutesEarlier := now.Add(-5 * time.Minute)
	oneHourLater := now.Add(time.Hour)

	tests := []struct {
		name         string
		timeInterval *TimeInterval
		want         time.Duration
	}{
		{
			name: "NewTimeInterval",
			timeInterval: NewTimeInterval(
				now,
				oneHourLater,
				false,
				false,
			),
			want: time.Hour,
		},
		{
			name: "NewTimeInterval with start exclusive",
			timeInterval: NewTimeInterval(
				fiveMinutesEarlier,
				now,
				true,
				false,
			),
			want: 5 * time.Minute,
		},
		{
			name: "NewTimeInterval with end inclusive",
			timeInterval: NewTimeInterval(
				now,
				oneHourLater,
				false,
				true,
			),
			want: time.Hour,
		},
		{
			name: "NewTimeInterval with negative duration",
			timeInterval: NewTimeInterval(
				oneHourLater,
				now,
				false,
				false,
			),
			want: -time.Hour,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.timeInterval.Duration()
			assert.Equal(t, tt.want, got)
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
			name: "ToHalfOpen",
			timeInterval: NewTimeInterval(
				now,
				oneHourLater,
				false,
				false,
			),
		},
		{
			name: "ToHalfOpen with start exclusive",
			timeInterval: NewTimeInterval(
				now,
				oneHourLater,
				true,
				false,
			),
		},
		{
			name: "ToHalfOpen with end inclusive",
			timeInterval: NewTimeInterval(
				now,
				oneHourLater,
				false,
				true,
			),
		},
		{
			name: "ToHalfOpen with both true",
			timeInterval: NewTimeInterval(
				oneHourLater,
				now,
				true,
				true,
			),
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
	timeInterval := NewTimeInterval(
		time.Now(),
		time.Now().Add(time.Hour),
		false,
		false,
	)
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
			name: "ToProtoTimeInterval",
			timeInterval: NewTimeInterval(
				time.Now(),
				time.Now().Add(time.Hour),
				false,
				false,
			),
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

		got := protoToTimeInterval(input.ToProtoTimeInterval())

		assert.Equal(t, input.Start, got.Start)
		assert.Equal(t, input.End, got.End)
		assert.Equal(t, input.StartExclusive, got.StartExclusive)
		assert.Equal(t, input.EndInclusive, got.EndInclusive)
	})
}
