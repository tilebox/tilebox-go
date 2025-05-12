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

	cluster, err := client.Clusters.Get(ctx, "testing-4qgCk4qHH85qR7")
	if err != nil {
		slog.ErrorContext(ctx, "failed to get cluster", slog.Any("error", err))
		return
	}

	job, err := client.Jobs.Submit(ctx, "spawn-workflow-tree", cluster,
		[]workflows.Task{
			&pbtask.SampleTask{ // protobuf task
				Message:      "hello go runner!",
				Depth:        8,
				BranchFactor: 4,
			},
		},
	)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to submit job", "error", err)
		return
	}

	slog.InfoContext(ctx, "Job submitted", slog.String("job_id", job.ID.String()))
}
