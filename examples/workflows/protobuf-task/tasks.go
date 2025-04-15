package protobuftask

import (
	"context"
	"log/slog"

	examplesv1 "github.com/tilebox/tilebox-go/protogen/go/examples/v1"
	"github.com/tilebox/tilebox-go/workflows/v1"
)

// SampleTask is a regular struct task that submits a protobuf task.
type SampleTask struct {
	Message      string
	Depth        int
	BranchFactor int
}

func (t *SampleTask) Identifier() workflows.TaskIdentifier {
	return workflows.NewTaskIdentifier("tilebox.com/task-runner/SampleTask", "v1.0")
}

func (t *SampleTask) Execute(ctx context.Context) error {
	slog.InfoContext(ctx, "Spawning a Tree", slog.String("message", t.Message), slog.Int("depth", t.Depth))

	_, err := workflows.SubmitSubtask(ctx, &SpawnWorkflowTreeTask{
		examplesv1.SpawnWorkflowTreeTask{
			CurrentLevel: 0,
			Depth:        int64(t.Depth),
			BranchFactor: int64(t.BranchFactor),
		},
	})
	if err != nil {
		return err
	}

	return nil
}

// SpawnWorkflowTreeTask embeds a protobuf message
type SpawnWorkflowTreeTask struct {
	examplesv1.SpawnWorkflowTreeTask
}

func (t *SpawnWorkflowTreeTask) Identifier() workflows.TaskIdentifier {
	return workflows.NewTaskIdentifier("tilebox.com/task-runner/SpawnWorkflowTreeTask", "v1.0")
}

func (t *SpawnWorkflowTreeTask) Execute(ctx context.Context) error {
	if t.GetCurrentLevel() >= (t.GetDepth() - 1) {
		return nil
	}

	subtasks := make([]workflows.Task, t.GetBranchFactor())
	for i := range t.GetBranchFactor() {
		subtasks[i] = &SpawnWorkflowTreeTask{
			examplesv1.SpawnWorkflowTreeTask{
				CurrentLevel: t.GetCurrentLevel() + 1,
				Depth:        t.GetDepth(),
				BranchFactor: t.GetBranchFactor(),
			},
		}
	}

	_, err := workflows.SubmitSubtasks(ctx, subtasks)
	return err
}
