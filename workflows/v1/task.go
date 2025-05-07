package workflows // import "github.com/tilebox/tilebox-go/workflows/v1"

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
)

// TaskIdentifier is the struct that defines the unique identifier of a task.
// It is used to uniquely identify a task and specify its version.
type TaskIdentifier interface {
	Name() string
	Version() string
	Display() string
}

type taskIdentifier struct {
	name    string
	version string
}

var _ TaskIdentifier = &taskIdentifier{}

func NewTaskIdentifier(name, version string) TaskIdentifier {
	return taskIdentifier{
		name:    name,
		version: version,
	}
}

// Name returns the name of the task.
func (t taskIdentifier) Name() string {
	return t.name
}

// Version returns the version of the task.
func (t taskIdentifier) Version() string {
	return t.version
}

// Display returns a human-readable string representation of the task identifier, to be used in graph visualizations.
// Can be overridden during task execution to provide a more descriptive name.
func (t taskIdentifier) Display() string {
	return t.name
}

func (t taskIdentifier) String() string {
	return fmt.Sprintf("%s@%s", t.name, t.version)
}

// Task is the interface for a task that can be submitted to the workflow service.
// It doesn't need to be identifiable or executable, but it can be both.
type Task interface{}

// ExplicitlyIdentifiableTask is the interface for a task that provides a user-defined task identifier.
// The identifier is used to uniquely identify the task and specify its version. If a task is not an
// ExplicitlyIdentifiableTask, the task runner will generate an identifier for it using reflection.
type ExplicitlyIdentifiableTask interface {
	Identifier() TaskIdentifier
}

// ExecutableTask is the interface for a task that can be executed, and therefore be registered with a task runner.
type ExecutableTask interface {
	Execute(ctx context.Context) error
}

func identifierFromTask(task Task) TaskIdentifier {
	if identifiableTask, ok := task.(ExplicitlyIdentifiableTask); ok {
		return identifiableTask.Identifier()
	}
	return &taskIdentifier{
		name:    getStructName(task),
		version: "v0.0", // default version
	}
}

// ValidateIdentifier performs client-side validation on a task identifier.
func ValidateIdentifier(identifier TaskIdentifier) error {
	if identifier.Name() == "" {
		return errors.New("task name is empty")
	}
	if len(identifier.Name()) > 256 {
		return errors.New("task name is too long")
	}
	_, _, err := parseVersion(identifier.Version())
	if err != nil {
		return err
	}
	return nil
}

var versionPattern = regexp.MustCompile(`^v(\d+)\.(\d+)$`)

// parseVersion parses the major and minor version from a string in the format "vMajor.Minor" and returns them as int64.
// If the version string is not in the correct format, an error is returned. It uses a regular expression to parse the version.
func parseVersion(version string) (int64, int64, error) {
	match := versionPattern.FindStringSubmatch(version)
	if len(match) != 3 { // first the whole match, then the two submatches (major and minor version)
		return 0, 0, fmt.Errorf("invalid task version: %s. expected format v<major>.<minor>, e.g. v3.2", version)
	}
	major, err := strconv.Atoi(match[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid major version: %s", match[1])
	}
	minor, err := strconv.Atoi(match[2])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid minor version: %s", match[1])
	}
	return int64(major), int64(minor), nil
}

// getStructName returns the struct name of a task. If the task is a pointer, the name of the pointed-to type is returned.
// This function is used to generate a default identifier name for a task if it doesn't provide an explicit identifier.
func getStructName(task interface{}) string {
	t := reflect.TypeOf(task)
	if t.Kind() == reflect.Pointer {
		return t.Elem().Name()
	}
	return t.Name()
}
