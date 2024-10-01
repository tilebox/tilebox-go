package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/tilebox/tilebox-go/examples/helloworld"
	"github.com/tilebox/tilebox-go/workflows/v1"
)

func main() {
	runner, err := workflows.NewTaskRunner(
		workflows.NewTaskClient(
			workflows.WithAPIKey(os.Getenv("TILEBOX_API_KEY")),
		),
		workflows.WithCluster("testing-4qgCk4qHH85qR7"),
	)
	if err != nil {
		slog.Error("failed to create task runner", "error", err)
		return
	}

	err = runner.RegisterTasks(
		&helloworld.HelloTask{},
	)
	if err != nil {
		slog.Error("failed to register task", "error", err)
		return
	}

	runner.Run(context.Background())
}
