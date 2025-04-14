package subtask

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Options(t *testing.T) {
	tests := []struct {
		name    string
		options []SubmitOption
		want    SubmitOptions
	}{
		{
			name: "with no dependencies",
			options: []SubmitOption{
				WithDependencies(),
			},
			want: SubmitOptions{},
		},
		{
			name: "with dependencies",
			options: []SubmitOption{
				WithDependencies(1, 3, 2),
			},
			want: SubmitOptions{
				Dependencies: []FutureTask{1, 3, 2},
			},
		},
		{
			name: "with cluster slug",
			options: []SubmitOption{
				WithClusterSlug("slug"),
			},
			want: SubmitOptions{
				ClusterSlug: "slug",
			},
		},
		{
			name: "with max retries",
			options: []SubmitOption{
				WithMaxRetries(7),
			},
			want: SubmitOptions{
				MaxRetries: 7,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var options SubmitOptions
			for _, option := range tt.options {
				option(&options)
			}
			assert.Equal(t, tt.want, options)
		})
	}
}
