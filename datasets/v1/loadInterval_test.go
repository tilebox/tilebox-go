package datasets

import (
	"testing"
	"time"

	datasetsv1 "github.com/tilebox/tilebox-go/protogen/go/datasets/v1"
	"pgregory.net/rapid"
)

func TestNewEmptyTimeInterval(t *testing.T) {
	got := newEmptyTimeInterval()

	if got.Duration() != 0 {
		t.Errorf("newEmptyTimeInterval() Duration = %v, want 0", got.Duration())
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

			if got.StartExclusive {
				t.Errorf("ToHalfOpen() StartExclusive = %v, want false", got.StartExclusive)
			}
			if got.EndInclusive {
				t.Errorf("ToHalfOpen() EndInclusive = %v, want false", got.EndInclusive)
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

			if !got.GetStartTime().AsTime().Equal(tt.timeInterval.Start) {
				t.Errorf("ToProtoTimeInterval() Start = %v, want %v", got.GetStartTime().AsTime(), tt.timeInterval.Start)
			}
			if !got.GetEndTime().AsTime().Equal(tt.timeInterval.End) {
				t.Errorf("ToProtoTimeInterval() End = %v, want %v", got.GetEndTime().AsTime(), tt.timeInterval.End)
			}
			if got.GetStartExclusive() != tt.timeInterval.StartExclusive {
				t.Errorf("ToProtoTimeInterval() StartExclusive = %v, want %v", got.GetStartExclusive(), tt.timeInterval.StartExclusive)
			}
			if got.GetEndInclusive() != tt.timeInterval.EndInclusive {
				t.Errorf("ToProtoTimeInterval() EndInclusive = %v, want %v", got.GetEndInclusive(), tt.timeInterval.EndInclusive)
			}
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

		if got.Start != input.Start {
			t.Errorf("protoToTimeInterval() Start = %v, want %v", got.Start, input.Start)
		}
		if got.End != input.End {
			t.Errorf("protoToTimeInterval() End = %v, want %v", got.End, input.End)
		}
		if got.StartExclusive != input.StartExclusive {
			t.Errorf("protoToTimeInterval() StartExclusive = %v, want %v", got.StartExclusive, input.StartExclusive)
		}
		if got.EndInclusive != input.EndInclusive {
			t.Errorf("protoToTimeInterval() EndInclusive = %v, want %v", got.EndInclusive, input.EndInclusive)
		}
	})
}
