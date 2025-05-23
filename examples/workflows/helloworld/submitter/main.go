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

	cluster, err := client.Clusters.Get(ctx, "testing-4qgCk4qHH85qR7")
	if err != nil {
		slog.ErrorContext(ctx, "failed to get cluster", slog.Any("error", err))
		return
	}

	job, err := client.Jobs.Submit(ctx, "hello-world", cluster,
		[]workflows.Task{
			&helloworld.HelloTask{
				Name: "Tilebox",
			},
		},
	)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to submit job", slog.Any("error", err))
		return
	}

	slog.InfoContext(ctx, "Job submitted", slog.String("job_id", job.ID.String()))
}
