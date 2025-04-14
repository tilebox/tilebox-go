package job

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Options(t *testing.T) {
	tests := []struct {
		name    string
		options []SubmitOption
		want    SubmitJobOptions
	}{
		{
			name: "with max retries",
			options: []SubmitOption{
				WithMaxRetries(7),
			},
			want: SubmitJobOptions{
				MaxRetries: 7,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var options SubmitJobOptions
			for _, option := range tt.options {
				option(&options)
			}
			assert.Equal(t, tt.want, options)
		})
	}
}
