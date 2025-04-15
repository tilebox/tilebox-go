package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/tilebox/tilebox-go/examples/workflows/helloworld"
	"github.com/tilebox/tilebox-go/workflows/v1"
)

func main() {
	ctx := context.Background()

	client := workflows.NewClient(
		workflows.WithAPIKey(os.Getenv("TILEBOX_API_KEY")),
	)

	cluster, err := client.Clusters.Get(ctx, "testing-4qgCk4qHH85qR7")
	if err != nil {
		slog.ErrorContext(ctx, "failed to get cluster", slog.Any("error", err))
		return
	}

	runner, err := client.NewTaskRunner(cluster)
	if err != nil {
		slog.ErrorContext(ctx, "failed to create task runner", slog.Any("error", err))
		return
	}

	err = runner.RegisterTasks(&helloworld.HelloTask{})
	if err != nil {
		slog.ErrorContext(ctx, "failed to register tasks", slog.Any("error", err))
		return
	}

	runner.RunForever(context.Background())
}
