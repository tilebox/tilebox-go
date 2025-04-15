package main

import (
	"context"
	"log/slog"
	"os"

	protobuftask "github.com/tilebox/tilebox-go/examples/workflows/protobuf-task"
	"github.com/tilebox/tilebox-go/workflows/v1"
)

func main() {
	ctx := context.Background()

	client := workflows.NewClient(workflows.WithAPIKey(os.Getenv("TILEBOX_API_KEY")))

	cluster, err := client.Clusters.Get(ctx, "testing-4qgCk4qHH85qR7")
	if err != nil {
		slog.ErrorContext(ctx, "failed to get cluster", slog.Any("error", err))
		return
	}

	taskRunner, err := client.NewTaskRunner(cluster)
	if err != nil {
		slog.Error("failed to create task runner", slog.Any("error", err))
		return
	}

	err = taskRunner.RegisterTasks(
		&protobuftask.SampleTask{},            // regular struct task
		&protobuftask.SpawnWorkflowTreeTask{}, // protobuf task
	)
	if err != nil {
		slog.Error("failed to register tasks", slog.Any("error", err))
		return // exit the program if we can't register one of the tasks
	}

	taskRunner.RunForever(ctx)
}
