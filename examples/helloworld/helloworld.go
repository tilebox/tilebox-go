package helloworld

import (
	"context"
	"log/slog"
)

type HelloTask struct {
	Name string
}

// The Execute method isn't needed to submit a task.
func (t *HelloTask) Execute(context.Context) error {
	slog.Info("Hello World!", "Name", t.Name)
	return nil
}
