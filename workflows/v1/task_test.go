package workflows

import (
	"reflect"
	"strings"
	"testing"
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
			want: TaskIdentifier{
				Name:    "test",
				Version: "v0.0",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTaskIdentifier(tt.args.name, tt.args.version); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTaskIdentifier() = %v, wantMajor %v", got, tt.want)
			}
		})
	}
}

type emptyTask struct{}

type identifiableTask struct{}

func (t *identifiableTask) Identifier() TaskIdentifier {
	return TaskIdentifier{
		Name:    "myName",
		Version: "v1.2",
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
			want: TaskIdentifier{
				Name:    "emptyTask",
				Version: "v0.0",
			},
		},
		{
			name: "identifier identifiable task",
			args: args{
				task: &identifiableTask{},
			},
			want: TaskIdentifier{
				Name:    "myName",
				Version: "v1.2",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := identifierFromTask(tt.args.task); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("identifierFromTask() = %v, wantMajor %v", got, tt.want)
			}
		})
	}
}

func TestValidateIdentifier(t *testing.T) {
	type args struct {
		identifier TaskIdentifier
	}
	tests := []struct {
		name           string
		args           args
		wantErr        bool
		wantErrMessage string
	}{
		{
			name: "ValidateIdentifier",
			args: args{
				identifier: TaskIdentifier{
					Name:    "test",
					Version: "v0.0",
				},
			},
			wantErr: false,
		},
		{
			name: "ValidateIdentifier name empty",
			args: args{
				identifier: TaskIdentifier{
					Name:    "",
					Version: "v0.0",
				},
			},
			wantErr:        true,
			wantErrMessage: "task name is empty",
		},
		{
			name: "ValidateIdentifier name too long",
			args: args{
				identifier: TaskIdentifier{
					Name:    strings.Repeat("a", 257),
					Version: "v0.0",
				},
			},
			wantErr:        true,
			wantErrMessage: "task name is too long",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateIdentifier(tt.args.identifier)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateIdentifier() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				if !strings.Contains(err.Error(), tt.wantErrMessage) {
					t.Errorf("CreateCluster() error = %v, wantErrMessage %v", err, tt.wantErrMessage)
				}
				return
			}
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
			if gotMajor != tt.wantMajor {
				t.Errorf("parseVersion() gotMajor = %v, wantMajor %v", gotMajor, tt.wantMajor)
			}
			if gotMinor != tt.wantMinor {
				t.Errorf("parseVersion() gotMinor = %v, wantMajor %v", gotMinor, tt.wantMinor)
			}
		})
	}
}
