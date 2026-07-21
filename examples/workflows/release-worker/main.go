package main

import (
	"context"
	"log"
	"log/slog"

	workflows "github.com/tilebox/tilebox-go/workflows/v1"
)

type ProcessScene struct {
	SceneID string `json:"scene_id"`
}

func (*ProcessScene) Identifier() workflows.TaskIdentifier {
	return workflows.NewTaskIdentifier("example.com/scenes/process", "v1.0")
}

func (task *ProcessScene) Execute(ctx context.Context) error {
	slog.InfoContext(ctx, "processing scene", slog.String("scene_id", task.SceneID))
	_, err := workflows.SubmitSubtask(ctx, &WriteSummary{SceneID: task.SceneID})
	return err
}

type WriteSummary struct {
	SceneID string `json:"scene_id"`
}

func (*WriteSummary) Identifier() workflows.TaskIdentifier {
	return workflows.NewTaskIdentifier("example.com/scenes/write-summary", "v1.0")
}

func (task *WriteSummary) Execute(ctx context.Context) error {
	slog.InfoContext(ctx, "writing scene summary", slog.String("scene_id", task.SceneID))
	return nil
}

func main() {
	worker := workflows.NewWorker()
	if err := worker.RegisterTasks(
		&ProcessScene{},
		&WriteSummary{},
	); err != nil {
		log.Fatal(err)
	}
	if err := worker.Serve(context.Background()); err != nil {
		log.Fatal(err)
	}
}
