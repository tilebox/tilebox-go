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

	job, err := client.Jobs.Submit(ctx, "spawn-workflow-tree",
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
