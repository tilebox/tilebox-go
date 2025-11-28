package mapreduce

import (
	"context"
	"log/slog"
	"time"

	"github.com/tilebox/tilebox-go/workflows/v1"
	"github.com/tilebox/tilebox-go/workflows/v1/subtask"
)

const (
	NumMapTasks    = 180
	NumReduceTasks = 10
)

// RootTask is the entry point that spawns all map and reduce tasks.
// It creates 180 map tasks and 10 reduce tasks, where each reduce task
// depends on all map tasks completing first.
type RootTask struct{}

func (t *RootTask) Execute(ctx context.Context) error {
	slog.Info("Starting MapReduce job")

	mapTasks := make([]workflows.Task, NumMapTasks)
	for i := range NumMapTasks {
		mapTasks[i] = &MapTask{
			PartitionIndex: i,
		}
	}

	mapFutures, err := workflows.SubmitSubtasks(ctx, mapTasks)
	if err != nil {
		return err
	}

	slog.Info("Submitted map tasks", slog.Int("count", len(mapFutures)))

	mapProgress := workflows.Progress("map")
	_ = mapProgress.Add(ctx, uint64(len(mapFutures)))

	for i := range NumReduceTasks {
		_, err := workflows.SubmitSubtask(ctx, &ReduceTask{
			ReducerIndex: i,
		}, subtask.WithDependencies(mapFutures...))
		if err != nil {
			return err
		}
	}

	slog.Info("Submitted reduce tasks", slog.Int("count", NumReduceTasks))

	reduceProgress := workflows.Progress("reduce")
	_ = reduceProgress.Add(ctx, NumReduceTasks)

	return nil
}

// MapTask processes a single partition of the input data.
type MapTask struct {
	PartitionIndex int
}

func (t *MapTask) Execute(ctx context.Context) error {
	slog.Info("Executing map task", slog.Int("partition", t.PartitionIndex))

	time.Sleep(10 * time.Millisecond)

	progress := workflows.Progress("map")
	_ = progress.Done(ctx, 1)
	return nil
}

// ReduceTask aggregates results from all map tasks.
type ReduceTask struct {
	ReducerIndex int
}

func (t *ReduceTask) Execute(ctx context.Context) error {
	slog.Info("Executing reduce task", slog.Int("reducer", t.ReducerIndex))

	time.Sleep(10 * time.Millisecond)

	progress := workflows.Progress("reduce")
	_ = progress.Done(ctx, 1)
	return nil
}
