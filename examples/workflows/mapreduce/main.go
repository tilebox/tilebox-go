package main

import (
	"context"
	"log/slog"

	"github.com/tilebox/tilebox-go/workflows/v1"
	"github.com/tilebox/tilebox-go/workflows/v1/subtask"
)

// RootTask is the entry point that spawns all map and reduce tasks.
type RootTask struct {
	NumMapTasks    uint64
	NumReduceTasks uint64
}

func (t *RootTask) Execute(ctx context.Context) error {
	slog.InfoContext(ctx, "Starting MapReduce job")

	mapTasks := make([]workflows.Task, t.NumMapTasks)
	for i := range t.NumMapTasks {
		mapTasks[i] = &MapTask{Index: i}
	}

	mapFutures, err := workflows.SubmitSubtasks(ctx, mapTasks)
	if err != nil {
		return err
	}

	slog.InfoContext(ctx, "Submitted map tasks", slog.Int("count", len(mapFutures)))

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

	slog.InfoContext(ctx, "Submitted reduce tasks", slog.Uint64("count", t.NumReduceTasks))

	reduceProgress := workflows.Progress("reduce")
	_ = reduceProgress.Add(ctx, t.NumReduceTasks)

	return nil
}

// MapTask processes a single item of input data.
type MapTask struct {
	Index uint64
}

func (t *MapTask) Execute(ctx context.Context) error {
	slog.InfoContext(ctx, "Executing map task", slog.Uint64("index", t.Index))
	_ = workflows.Progress("map").Done(ctx, 1)
	return nil
}

// ReduceTask aggregates results from map tasks.
type ReduceTask struct {
	ReducerIndex uint64
}

func (t *ReduceTask) Execute(ctx context.Context) error {
	slog.InfoContext(ctx, "Executing reduce task", slog.Uint64("reducer_index", t.ReducerIndex))
	_ = workflows.Progress("reduce").Done(ctx, 1)
	return nil
}

func main() {
	ctx := context.Background()
	workflows.ConfigureConsoleLogging(slog.LevelInfo)
	client := workflows.NewClient()

	job, err := client.Jobs.Submit(ctx, "mapreduce-example",
		[]workflows.Task{
			&RootTask{
				NumMapTasks:    100,
				NumReduceTasks: 8,
			},
		},
	)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to submit job", slog.Any("error", err))
		return
	}

	slog.InfoContext(ctx, "Job submitted", slog.String("job_id", job.ID.String()))

	runner, err := client.NewTaskRunner(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "failed to create task runner", slog.Any("error", err))
		return
	}

	if err := runner.RegisterTasks(
		&RootTask{},
		&MapTask{},
		&ReduceTask{},
	); err != nil {
		slog.ErrorContext(ctx, "failed to register tasks", slog.Any("error", err))
		return
	}

	runner.RunForever(ctx)
}
