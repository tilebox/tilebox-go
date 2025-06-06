<h1 align="center">
  <img src="https://storage.googleapis.com/tbx-web-assets-2bad228/banners/tilebox-banner.svg" alt="Tilebox Logo">
  <br>
</h1>

<div align="center">
  <a href="https://pkg.go.dev/github.com/tilebox/tilebox-go">
    <img src="https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square&color=f43f5e" alt="PyPi Latest Release badge"/>
  </a>
  <a href="https://github.com/tilebox/tilebox-python/blob/main/LICENSE">
    <img src="https://img.shields.io/github/license/tilebox/tilebox-go.svg?style=flat-square&color=f43f5e" alt="MIT License"/>
  </a>
  <a href="https://github.com/tilebox/tilebox-python/actions">
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

Go library for [Tilebox](https://tilebox.com) - a framework for space data management and workflow orchestration.

## Installation

Run the following command to add the library to your project:

```bash
go get github.com/tilebox/tilebox-go
```

For Tilebox datasets type generation, you will need to install [tilebox-generate](https://github.com/tilebox/tilebox-generate) command-line tool.

## Examples

For examples on how to use the library, see the [examples](examples) directory.

## Usage

### Writing a Task

Here we define a simple task that prints "Hello World!" to the console:

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
func (t *HelloTask) Execute(context.Context) error {
	slog.Info("Hello World!", "Name", t.Name)
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
	client := workflows.NewClient()

	job, err := client.Jobs.Submit(ctx, "hello-world",
		[]workflows.Task{
			&HelloTask{
				Name: "Tilebox",
			},
		},
	)
	if err != nil {
		slog.Error("Failed to submit job", slog.Any("error", err))
		return
	}

	slog.Info("Job submitted", slog.String("job_id", job.ID.String()))
}
```

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
func (t *HelloTask) Execute(context.Context) error {
	slog.Info("Hello World!", "Name", t.Name)
	return nil
}

func main() {
	ctx := context.Background()
	client := workflows.NewClient()
	
	runner, err := client.NewTaskRunner(ctx)
	if err != nil {
		slog.Error("failed to create task runner", slog.Any("error", err))
		return
	}

	err = runner.RegisterTasks(&HelloTask{})
	if err != nil {
		slog.Error("failed to register tasks", slog.Any("error", err))
		return
	}

	runner.RunForever(ctx)
}
```


## License

Distributed under the MIT License (`The MIT License`).
