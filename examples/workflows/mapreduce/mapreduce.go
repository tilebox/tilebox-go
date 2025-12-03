package mapreduce

import (
	"context"
	"log/slog"

	"github.com/tilebox/tilebox-go/workflows/v1"
	"github.com/tilebox/tilebox-go/workflows/v1/subtask"
)

// RootTask is the entry point that spawns all map and reduce tasks.
type RootTask struct {
	NumMapTasks    int
	NumReduceTasks int
}

func (t *RootTask) Execute(ctx context.Context) error {
	slog.Info("Starting MapReduce job")

	mapTasks := make([]workflows.Task, t.NumMapTasks)
	for i := range t.NumMapTasks {
		mapTasks[i] = &MapTask{Index: i}
	}

	mapFutures, err := workflows.SubmitSubtasks(ctx, mapTasks)
	if err != nil {
		return err
	}

	slog.Info("Submitted map tasks", slog.Int("count", len(mapFutures)))

	mapProgress := workflows.Progress("map")
	_ = mapProgress.Add(ctx, uint64(len(mapFutures)))

	for i := range t.NumReduceTasks {
		_, err := workflows.SubmitSubtask(ctx, &ReduceTask{
			ReducerIndex: i,
		}, subtask.WithDependencies(mapFutures...))
		if err != nil {
			return err
		}
	}

	slog.Info("Submitted reduce tasks", slog.Int("count", t.NumReduceTasks))

	reduceProgress := workflows.Progress("reduce")
	_ = reduceProgress.Add(ctx, uint64(t.NumReduceTasks))

	return nil
}

// MapTask processes a single item of input data.
type MapTask struct {
	Index int
}

func (t *MapTask) Execute(ctx context.Context) error {
	slog.Info("Executing map task", slog.Int("index", t.Index))
	_ = workflows.Progress("map").Done(ctx, 1)
	return nil
}

// ReduceTask aggregates results from map tasks.
type ReduceTask struct {
	ReducerIndex int
}

func (t *ReduceTask) Execute(ctx context.Context) error {
	slog.Info("Executing reduce task", slog.Int("reducer_index", t.ReducerIndex))
	_ = workflows.Progress("reduce").Done(ctx, 1)
	return nil
}
