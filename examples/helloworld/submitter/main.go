package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/google/uuid"
	"github.com/tilebox/tilebox-go/examples/helloworld"
	"github.com/tilebox/tilebox-go/workflows/v1"
)

func main() {
	ctx := context.Background()

	jobs := workflows.NewJobService(
		workflows.NewJobClient(
			workflows.WithAPIKey(os.Getenv("TILEBOX_API_KEY")),
		),
	)

	job, err := jobs.Submit(ctx, "hello-world", "testing-4qgCk4qHH85qR7", 0, uuid.Nil,
		&helloworld.HelloTask{
			Name: "Tilebox",
		},
	)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to submit job", "error", err)
		return
	}

	slog.InfoContext(ctx, "Job submitted", "job_id", uuid.Must(uuid.FromBytes(job.GetId().GetUuid())))
}
