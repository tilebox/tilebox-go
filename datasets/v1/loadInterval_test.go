package datasets

import (
	"testing"
	"time"

	datasetsv1 "github.com/tilebox/tilebox-go/protogen/go/datasets/v1"
	"pgregory.net/rapid"
)

func TestNewEmptyTimeInterval(t *testing.T) {
	got := NewEmptyTimeInterval()

	if got.Duration() != 0 {
		t.Errorf("NewEmptyTimeInterval() Duration = %v, want 0", got.Duration())
	}
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
			if got := tt.timeInterval.Duration(); got != tt.want {
				t.Errorf("Duration() = %v, want %v", got, tt.want)
			}
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

			if got.startExclusive {
				t.Errorf("ToHalfOpen() startExclusive = %v, want false", got.startExclusive)
			}
			if got.endInclusive {
				t.Errorf("ToHalfOpen() endInclusive = %v, want false", got.endInclusive)
			}
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

	if got != nil {
		t.Errorf("ToProtoDatapointInterval() = %v, want nil", got)
	}
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

			if !got.GetStartTime().AsTime().Equal(tt.timeInterval.start) {
				t.Errorf("ToProtoTimeInterval() start = %v, want %v", got.GetStartTime().AsTime(), tt.timeInterval.start)
			}
			if !got.GetEndTime().AsTime().Equal(tt.timeInterval.end) {
				t.Errorf("ToProtoTimeInterval() end = %v, want %v", got.GetEndTime().AsTime(), tt.timeInterval.end)
			}
			if got.GetStartExclusive() != tt.timeInterval.startExclusive {
				t.Errorf("ToProtoTimeInterval() startExclusive = %v, want %v", got.GetStartExclusive(), tt.timeInterval.startExclusive)
			}
			if got.GetEndInclusive() != tt.timeInterval.endInclusive {
				t.Errorf("ToProtoTimeInterval() endInclusive = %v, want %v", got.GetEndInclusive(), tt.timeInterval.endInclusive)
			}
		})
	}
}

func Test_protoToTimeIntervalRoundtrip(t *testing.T) {
	genTimeInterval := rapid.Custom(func(t *rapid.T) *TimeInterval {
		return &TimeInterval{
			start:          time.Unix(rapid.Int64().Draw(t, "start"), 0).UTC(),
			end:            time.Unix(rapid.Int64().Draw(t, "end"), 0).UTC(),
			startExclusive: rapid.Bool().Draw(t, "startExclusive"),
			endInclusive:   rapid.Bool().Draw(t, "endInclusive"),
		}
	})

	rapid.Check(t, func(t *rapid.T) {
		input := genTimeInterval.Draw(t, "time interval")

		got := protoToTimeInterval(input.ToProtoTimeInterval())

		if got.start != input.start {
			t.Errorf("protoToTimeInterval() start = %v, want %v", got.start, input.start)
		}
		if got.end != input.end {
			t.Errorf("protoToTimeInterval() end = %v, want %v", got.end, input.end)
		}
		if got.startExclusive != input.startExclusive {
			t.Errorf("protoToTimeInterval() startExclusive = %v, want %v", got.startExclusive, input.startExclusive)
		}
		if got.endInclusive != input.endInclusive {
			t.Errorf("protoToTimeInterval() endInclusive = %v, want %v", got.endInclusive, input.endInclusive)
		}
	})
}
