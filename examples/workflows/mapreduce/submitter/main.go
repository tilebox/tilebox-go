package main

import (
	"context"
	"log/slog"

	"github.com/tilebox/tilebox-go/examples/workflows/mapreduce"
	"github.com/tilebox/tilebox-go/workflows/v1"
)

func main() {
	ctx := context.Background()
	client := workflows.NewClient()

	job, err := client.Jobs.Submit(ctx, "mapreduce-example",
		[]workflows.Task{
			&mapreduce.RootTask{
				NumMapTasks:    100,
				NumReduceTasks: 8,
			},
		},
	)
	if err != nil {
		slog.Error("Failed to submit job", slog.Any("error", err))
		return
	}

	slog.Info("Job submitted", slog.String("job_id", job.ID.String()))
}
