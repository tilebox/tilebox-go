package main

import (
	"context"
	"log/slog"

	"github.com/tilebox/tilebox-go/examples/workflows/helloworld"
	"github.com/tilebox/tilebox-go/workflows/v1"
)

func main() {
	ctx := context.Background()
	client := workflows.NewClient()

	runner, err := client.NewTaskRunner(ctx)
	if err != nil {
		slog.Error("failed to create task runner", slog.Any("error", err))
		return
	}

	err = runner.RegisterTasks(&helloworld.HelloTask{})
	if err != nil {
		slog.Error("failed to register tasks", slog.Any("error", err))
		return
	}

	runner.RunForever(ctx)
}
