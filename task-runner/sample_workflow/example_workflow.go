// Package sample_workflow files are not located in the task-runner package since we want to test unexported field
// behavior as well
package sample_workflow

import (
	"context"
	workflowv1 "github.com/tilebox/tilebox-go/task-runner/protogen/go/workflow/v1"
	"github.com/tilebox/tilebox-go/task-runner/services"
	"log/slog"
)

type SampleTask struct {
	Message      string
	Depth        int
	BranchFactor int
}

func (t *SampleTask) Identifier() services.TaskIdentifier {
	return services.NewTaskIdentifier("tilebox.com/task-runner/SampleTask", "v1.0")
}

func (t *SampleTask) Execute(ctx context.Context) error {
	slog.InfoContext(ctx, "Spawning a Tree", "message", t.Message, "depth", t.Depth)

	err := services.SubmitSubtasks(ctx, &SpawnWorkflowTreeTask{
		workflowv1.SpawnWorkflowTreeTask{
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
	workflowv1.SpawnWorkflowTreeTask
}

func (n *SpawnWorkflowTreeTask) Identifier() services.TaskIdentifier {
	return services.NewTaskIdentifier("tilebox.com/task-runner/SpawnWorkflowTreeTask", "v1.0")
}

func (n *SpawnWorkflowTreeTask) Execute(ctx context.Context) error {
	if n.CurrentLevel >= (n.Depth - 1) {
		return nil
	}

	subtasks := make([]services.Task, n.BranchFactor)
	for i := 0; i < int(n.BranchFactor); i++ {
		subtasks[i] = &SpawnWorkflowTreeTask{
			workflowv1.SpawnWorkflowTreeTask{
				CurrentLevel: n.CurrentLevel + 1,
				Depth:        n.Depth,
				BranchFactor: n.BranchFactor,
			},
		}
	}

	return services.SubmitSubtasks(ctx, subtasks...)
}
