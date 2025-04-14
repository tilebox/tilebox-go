package runner

import (
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/metric/noop"
)

func Test_Options(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	tests := []struct {
		name    string
		options []Option
		want    Options
	}{
		{
			name: "with logger",
			options: []Option{
				WithRunnerLogger(logger),
			},
			want: Options{
				Logger: logger,
			},
		},
		{
			name: "with disable metrics",
			options: []Option{
				WithDisableMetrics(),
			},
			want: Options{
				MeterProvider: noop.NewMeterProvider(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var options Options
			for _, option := range tt.options {
				option(&options)
			}
			assert.Equal(t, tt.want, options)
		})
	}
}
