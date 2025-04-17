package logger

import (
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Options(t *testing.T) {
	tests := []struct {
		name    string
		options []Option
		want    Options
	}{
		{
			name: "with level",
			options: []Option{
				WithLevel(slog.LevelDebug),
			},
			want: Options{
				Level: slog.LevelDebug,
			},
		},
		{
			name: "with endpoint URL",
			options: []Option{
				WithEndpointURL("https://api.axiom.co/v1/logs"),
			},
			want: Options{
				Endpoint: "https://api.axiom.co/v1/logs",
			},
		},
		{
			name: "with headers",
			options: []Option{
				WithHeaders(map[string]string{
					"Authorization": "Bearer ABC",
					"Something":     "here",
				}),
			},
			want: Options{
				Headers: map[string]string{
					"Authorization": "Bearer ABC",
					"Something":     "here",
				},
			},
		},
		{
			name: "with export interval",
			options: []Option{
				WithExportInterval(20 * time.Second),
			},
			want: Options{
				ExportInterval: 20 * time.Second,
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
