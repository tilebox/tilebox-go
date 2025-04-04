package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/tilebox/tilebox-go/examples/helloworld"
	"github.com/tilebox/tilebox-go/workflows/v1"
)

func main() {
	ctx := context.Background()

	client := workflows.NewClient(
		workflows.WithAPIKey(os.Getenv("TILEBOX_API_KEY")),
	)

	job, err := client.Jobs.Submit(ctx, "hello-world", "testing-4qgCk4qHH85qR7", 0,
		&helloworld.HelloTask{
			Name: "Tilebox",
		},
	)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to submit job", slog.Any("error", err))
		return
	}

	slog.InfoContext(ctx, "Job submitted", slog.String("job_id", job.ID.String()))
}
