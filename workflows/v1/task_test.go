package workflows

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTaskIdentifier(t *testing.T) {
	type args struct {
		name    string
		version string
	}
	tests := []struct {
		name string
		args args
		want TaskIdentifier
	}{
		{
			name: "NewTaskIdentifier",
			args: args{
				name:    "test",
				version: "v0.0",
			},
			want: taskIdentifier{
				name:    "test",
				version: "v0.0",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewTaskIdentifier(tt.args.name, tt.args.version)

			assert.Equal(t, tt.want.Name(), got.Name())
			assert.Equal(t, tt.want.Version(), got.Version())
		})
	}
}

type emptyTask struct{}

type identifiableTask struct{}

func (t *identifiableTask) Identifier() TaskIdentifier {
	return taskIdentifier{
		name:    "myName",
		version: "v1.2",
	}
}

func Test_identifierFromTask(t *testing.T) {
	type args struct {
		task Task
	}
	tests := []struct {
		name string
		args args
		want TaskIdentifier
	}{
		{
			name: "identifier empty task",
			args: args{
				task: &emptyTask{},
			},
			want: taskIdentifier{
				name:    "emptyTask",
				version: "v0.0",
			},
		},
		{
			name: "identifier identifiable task",
			args: args{
				task: &identifiableTask{},
			},
			want: taskIdentifier{
				name:    "myName",
				version: "v1.2",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := identifierFromTask(tt.args.task)

			assert.Equal(t, tt.want.Name(), got.Name())
			assert.Equal(t, tt.want.Version(), got.Version())
		})
	}
}

func TestValidateIdentifier(t *testing.T) {
	type args struct {
		identifier TaskIdentifier
	}
	tests := []struct {
		name    string
		args    args
		wantErr string
	}{
		{
			name: "ValidateIdentifier",
			args: args{
				identifier: taskIdentifier{
					name:    "test",
					version: "v0.0",
				},
			},
		},
		{
			name: "ValidateIdentifier name empty",
			args: args{
				identifier: taskIdentifier{
					name:    "",
					version: "v0.0",
				},
			},
			wantErr: "task name is empty",
		},
		{
			name: "ValidateIdentifier name too long",
			args: args{
				identifier: taskIdentifier{
					name:    strings.Repeat("a", 257),
					version: "v0.0",
				},
			},
			wantErr: "task name is too long",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateIdentifier(tt.args.identifier)
			if tt.wantErr != "" {
				// we wanted an error, let's check if we got one
				require.Error(t, err, "expected an error, got none")
				assert.Contains(t, err.Error(), tt.wantErr, "error didn't contain expected message: '%s', got error '%s' instead.", tt.wantErr, err.Error())
				return
			}
			// we didn't want an error:
			require.NoError(t, err, "got an unexpected error")
		})
	}
}

func Test_parseVersion(t *testing.T) {
	type args struct {
		version string
	}
	tests := []struct {
		name      string
		args      args
		wantMajor int64
		wantMinor int64
		wantErr   bool
	}{
		{
			name: "parseVersion v0.0",
			args: args{
				version: "v0.0",
			},
			wantMajor: 0,
			wantMinor: 0,
			wantErr:   false,
		},
		{
			name: "parseVersion v2.3",
			args: args{
				version: "v2.3",
			},
			wantMajor: 2,
			wantMinor: 3,
			wantErr:   false,
		},
		{
			name: "parseVersion wrong format",
			args: args{
				version: "2.3",
			},
			wantErr: true,
		},
		{
			name: "parseVersion wrong major",
			args: args{
				version: "vA.3",
			},
			wantErr: true,
		},
		{
			name: "parseVersion wrong minor",
			args: args{
				version: "v2.A",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMajor, gotMinor, err := parseVersion(tt.args.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantMajor, gotMajor)
			assert.Equal(t, tt.wantMinor, gotMinor)
		})
	}
}
