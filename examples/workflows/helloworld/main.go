package main

import (
	"context"
	"log/slog"

	"github.com/tilebox/tilebox-go/workflows/v1"
)

type HelloTask struct {
	Name string
}

// The Execute method isn't needed to submit a task but is required to run a task on a task runner.
func (t *HelloTask) Execute(ctx context.Context) error {
	slog.InfoContext(ctx, "Hello World!", slog.String("Name", t.Name))
	return nil
}

func main() {
	ctx := context.Background()
	workflows.ConfigureConsoleLogging(slog.LevelInfo)
	client := workflows.NewClient()

	job, err := client.Jobs.Submit(ctx, "hello-world",
		[]workflows.Task{
			&HelloTask{
				Name: "Tilebox",
			},
		},
	)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to submit job", slog.Any("error", err))
		return
	}

	slog.InfoContext(ctx, "Job submitted", slog.String("job_id", job.ID.String()))

	runner, err := client.NewTaskRunner(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "failed to create task runner", slog.Any("error", err))
		return
	}

	if err := runner.RegisterTasks(&HelloTask{}); err != nil {
		slog.ErrorContext(ctx, "failed to register tasks", slog.Any("error", err))
		return
	}

	runner.RunForever(ctx)
}
