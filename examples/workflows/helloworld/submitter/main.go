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

	job, err := client.Jobs.Submit(ctx, "hello-world",
		[]workflows.Task{
			&helloworld.HelloTask{
				Name: "Tilebox",
			},
		},
	)
	if err != nil {
		slog.Error("Failed to submit job", slog.Any("error", err))
		return
	}

	slog.Info("Job submitted", slog.String("job_id", job.ID.String()))
}
