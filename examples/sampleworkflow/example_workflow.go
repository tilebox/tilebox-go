package sampleworkflow

import (
	"context"
	"log/slog"

	examplesv1 "github.com/tilebox/tilebox-go/protogen/go/examples/v1"
	"github.com/tilebox/tilebox-go/workflows/v1"
)

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

	err := workflows.SubmitSubtasks(ctx, &SpawnWorkflowTreeTask{
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

type SpawnWorkflowTreeTask struct {
	examplesv1.SpawnWorkflowTreeTask
}

func (n *SpawnWorkflowTreeTask) Identifier() workflows.TaskIdentifier {
	return workflows.NewTaskIdentifier("tilebox.com/task-runner/SpawnWorkflowTreeTask", "v1.0")
}

func (n *SpawnWorkflowTreeTask) Execute(ctx context.Context) error {
	if n.CurrentLevel >= (n.Depth - 1) {
		return nil
	}

	subtasks := make([]workflows.Task, n.BranchFactor)
	for i := range n.BranchFactor {
		subtasks[i] = &SpawnWorkflowTreeTask{
			examplesv1.SpawnWorkflowTreeTask{
				CurrentLevel: n.CurrentLevel + 1,
				Depth:        n.Depth,
				BranchFactor: n.BranchFactor,
			},
		}
	}

	return workflows.SubmitSubtasks(ctx, subtasks...)
}
