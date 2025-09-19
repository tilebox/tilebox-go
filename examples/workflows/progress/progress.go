package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/tilebox/tilebox-go/workflows/v1"
)

type ProgressRootTask struct {
	N int
}

func (t *ProgressRootTask) Execute(ctx context.Context) error {
	// report that 10 units of work need to be done
	err := workflows.DefaultProgress().Add(ctx, uint64(t.N))
	if err != nil {
		return err
	}

	for range t.N {
		_, err := workflows.SubmitSubtask(ctx, &ProgressSubTask{})
		if err != nil {
			return err
		}
	}

	return nil
}

type ProgressSubTask struct{}

func (t *ProgressSubTask) Execute(ctx context.Context) error {
	// report that one unit of work has been completed
	err := workflows.DefaultProgress().Done(ctx, 1)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})))

	ctx := context.Background()
	client := workflows.NewClient()

	job, err := client.Jobs.Submit(ctx, "progress-example",
		[]workflows.Task{&ProgressRootTask{N: 5}},
	)
	if err != nil {
		slog.Error("Failed to submit job", slog.Any("error", err))
		return
	}

	slog.Info("Job submitted", slog.String("job_id", job.ID.String()))
	if len(job.Progress) == 0 {
		slog.Info("No job progress yet")
	}

	runner, err := client.NewTaskRunner(ctx)
	if err != nil {
		slog.Error("failed to create task runner", slog.Any("error", err))
		return
	}

	err = runner.RegisterTasks(&ProgressRootTask{}, &ProgressSubTask{})
	if err != nil {
		slog.Error("failed to register tasks", slog.Any("error", err))
		return
	}

	runner.RunAll(ctx)

	job, err = client.Jobs.Get(ctx, job.ID) // fetch up-to-date job progress
	if err != nil {
		slog.Error("Failed to get job", slog.Any("error", err))
		return
	}

	slog.Info("Job completed", slog.String("job_id", job.ID.String()))
	for _, progress := range job.Progress {
		slog.Info("Job progress", slog.String("label", progress.Label), slog.Uint64("total", progress.Total), slog.Uint64("done", progress.Done))
	}
}
