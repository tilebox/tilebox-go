package main

import (
	"context"
	"log/slog"

	pbtask "github.com/tilebox/tilebox-go/examples/workflows/protobuf-task"
	"github.com/tilebox/tilebox-go/workflows/v1"
)

func main() {
	ctx := context.Background()
	client := workflows.NewClient()

	taskRunner, err := client.NewTaskRunner(ctx)
	if err != nil {
		slog.Error("failed to create task runner", slog.Any("error", err))
		return
	}

	err = taskRunner.RegisterTasks(
		&pbtask.SampleTask{},            // regular struct task
		&pbtask.SpawnWorkflowTreeTask{}, // protobuf task
	)
	if err != nil {
		slog.Error("failed to register tasks", slog.Any("error", err))
		return // exit the program if we can't register one of the tasks
	}

	taskRunner.RunForever(ctx)
}
