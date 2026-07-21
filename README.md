<h1 align="center">
  <img src="https://storage.googleapis.com/tbx-web-assets-2bad228/banners/tilebox-banner.svg" alt="Tilebox Logo">
  <br>
</h1>

<div align="center">
  <a href="https://pkg.go.dev/github.com/tilebox/tilebox-go">
    <img src="https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square&color=f43f5e" alt="Go.dev reference badge"/>
  </a>
  <a href="https://github.com/tilebox/tilebox-go/blob/main/LICENSE">
    <img src="https://img.shields.io/github/license/tilebox/tilebox-go.svg?style=flat-square&color=f43f5e" alt="MIT License"/>
  </a>
  <a href="https://github.com/tilebox/tilebox-go/actions">
    <img src="https://img.shields.io/github/actions/workflow/status/tilebox/tilebox-go/main.yml?style=flat-square&color=f43f5e" alt="Build Status"/>
  </a>
  <a href="https://tilebox.com/discord">
    <img src="https://img.shields.io/badge/Discord-%235865F2.svg?style=flat-square&logo=discord&logoColor=white" alt="Join us on Discord"/>
  </a>
</div>

<p align="center">
  <a href="https://docs.tilebox.com/introduction"><b>Documentation</b></a>
  |
  <a href="https://console.tilebox.com/"><b>Console</b></a>
  |
  <a href="https://examples.tilebox.com/"><b>Example Gallery</b></a>
</p>

# Tilebox Go

Go library for [Tilebox](https://tilebox.com), a lightweight space data management and orchestration software - on ground and in orbit.

## Installation

Run the following command to add the library to your project:

```bash
go get github.com/tilebox/tilebox-go
```

For Tilebox datasets type generation, you will need to install [tilebox-generate](https://github.com/tilebox/tilebox-generate) command-line tool.

## Examples

For examples on how to use the library, see the [examples](examples) directory.

## Usage

### Filtering Dataset Queries

Fields marked queryable in a dataset schema can be filtered with fluent Boolean and numeric expressions. Multiple
expressions passed to `WithFilters` are combined with each other and with temporal and spatial filters using logical
AND.

```go
package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/tilebox/tilebox-go/datasets/v1"
	"github.com/tilebox/tilebox-go/query"
)

func main() {
	ctx := context.Background()
	client := datasets.NewClient()
	datasetID := uuid.MustParse("019c0123-4567-7890-abcd-ef0123456789")
	start := time.Date(2026, time.July, 1, 0, 0, 0, 0, time.UTC)
	end := start.Add(24 * time.Hour)

	datapoints, err := datasets.Collect(client.Datapoints.Query(ctx,
		datasetID,
		datasets.WithTemporalExtent(query.NewTimeInterval(start, end)),
		datasets.WithFilters(
			query.Field("eo_cloud_cover").LessThan(20.0),
			query.Or(
				query.Field("quality").GreaterThanOrEqual(80),
				query.Field("quality").IsNull(),
			),
		),
	))
	if err != nil {
		slog.ErrorContext(ctx, "Failed to query datapoints", slog.Any("error", err))
		return
	}

	slog.InfoContext(ctx, "Found datapoints", slog.Int("count", len(datapoints)))
}
```

Comparisons with absent fields evaluate to unknown and therefore do not match. Use `IsNull`, as above, when absent values
should be included explicitly. `NotEqual` also excludes absent values.

### Writing a Task

Here we define a simple task that logs "Hello World!":

```go
package helloworld

import (
	"context"
	"log/slog"

	"github.com/tilebox/tilebox-go/workflows/v1"
)

type HelloTask struct {
	Name string // You can add any fields you need to the task struct.
}

// The Execute method isn't needed to submit a task but is required to run a task.
func (t *HelloTask) Execute(ctx context.Context) error {
	slog.InfoContext(ctx, "Hello World!", slog.String("Name", t.Name))
	return nil
}

// The Identifier method is optional and will be generated if not provided.
func (t *HelloTask) Identifier() workflows.TaskIdentifier {
	return workflows.NewTaskIdentifier("hello-world", "v1.0")
}
```

### Submitting a Job

Here we create a Workflows client and submit a job with a single task:

```go
package main

import (
	"context"
	"log/slog"

	"github.com/tilebox/tilebox-go/workflows/v1"
)

type HelloTask struct {
	Name string
}

func main() {
	ctx := context.Background()
	workflows.ConfigureConsoleLogging(slog.LevelInfo)
	client := workflows.NewClient()

	job, err := client.Jobs.Submit(ctx, "hello-world",
		[]workflows.Task{
			&HelloTask{
				Name: "Tilebox",
			},
		},
	)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to submit job", slog.Any("error", err))
		return
	}

	slog.InfoContext(ctx, "Job submitted", slog.String("job_id", job.ID.String()))
}
```

`workflows.NewClient()` configures Tilebox OpenTelemetry export when an API key is available. `workflows.ConfigureConsoleLogging` adds console output without replacing the Tilebox exporter, so logs are written to both destinations when both are configured.

### Running a Worker

Here we create a TaskRunner and run a worker that is capable of executing `HelloTask` tasks:

```go
package main

import (
	"context"
	"log/slog"

	"github.com/tilebox/tilebox-go/workflows/v1"
)

type HelloTask struct {
	Name string
}

// The Execute method is required to run a task.
func (t *HelloTask) Execute(ctx context.Context) error {
	slog.InfoContext(ctx, "Hello World!", slog.String("Name", t.Name))
	return nil
}

func main() {
	ctx := context.Background()
	workflows.ConfigureConsoleLogging(slog.LevelInfo)
	client := workflows.NewClient()

	runner, err := client.NewTaskRunner(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "failed to create task runner", slog.Any("error", err))
		return
	}

	err = runner.RegisterTasks(&HelloTask{})
	if err != nil {
		slog.ErrorContext(ctx, "failed to register tasks", slog.Any("error", err))
		return
	}

	if err := runner.RunForever(ctx); err != nil {
		slog.ErrorContext(ctx, "task runner stopped", slog.Any("error", err))
	}
}
```

### Querying Job Telemetry

After a job has been executed, we can query runner logs associated with it.

```go
package main

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/tilebox/tilebox-go/workflows/v1"
	"github.com/tilebox/tilebox-go/workflows/v1/job"
)

func main() {
	ctx := context.Background()
	client := workflows.NewClient()
	jobID := uuid.MustParse("019e070c-63ba-1c7e-5f1d-65be3e22d52a")

	for logRecord, err := range client.Jobs.QueryLogs(ctx, jobID,
		job.WithLimit(100),
		job.WithSortDirection(job.Ascending),
	) {
		if err != nil {
			slog.ErrorContext(ctx, "failed to query job logs", slog.Any("error", err))
			return
		}

		slog.InfoContext(ctx, "job log", slog.String("body", logRecord.Body))
	}
}
```

## License

Distributed under the MIT License (`The MIT License`).
