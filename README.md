# Tilebox Go Library

This repository contains the Go library for Tilebox.

## Getting Started

### Installation

Run the following command to add the library to your project:

```bash
go get github.com/tilebox/tilebox-go
```

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

Here we create a JobService and submit a job with a single task:

```go
package main

import (
	"context"
	"log/slog"
	"os"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/tilebox/tilebox-go/grpc"
	"github.com/tilebox/tilebox-go/protogen/go/workflows/v1/workflowsv1connect"
	"github.com/tilebox/tilebox-go/workflows/v1"
)

const serverURL = "https://api.tilebox.com"

type HelloTask struct {
	Name string
}

func main() {
	ctx := context.Background()

	jobsClient := clientFromConfig(serverURL, os.Getenv("TILEBOX_API_KEY"))
	jobs := workflows.NewJobService(jobsClient)

	job, err := jobs.Submit(ctx, "hello-world", workflows.DefaultClusterSlug,
		&HelloTask{
			Name: "Tilebox",
		},
	)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to submit job", "error", err)
		return
	}

	slog.InfoContext(ctx, "Job submitted", "job_id", uuid.Must(uuid.FromBytes(job.GetId().GetUuid())))
}

func clientFromConfig(serverURL, authToken string) workflowsv1connect.JobServiceClient {
	return workflowsv1connect.NewJobServiceClient(
		grpc.RetryHTTPClient(), serverURL, connect.WithInterceptors(
			grpc.NewAddAuthTokenInterceptor(func() string {
				return authToken
			})),
	)
}

```

### Running a Worker

Here we create a TaskRunner and run a worker that is capable of executing `HelloTask` tasks:

```go
package main

import (
	"context"
	"log/slog"
	"os"

	"connectrpc.com/connect"
	"github.com/tilebox/tilebox-go/grpc"
	"github.com/tilebox/tilebox-go/protogen/go/workflows/v1/workflowsv1connect"
	"github.com/tilebox/tilebox-go/workflows/v1"
)

const serverURL = "https://api.tilebox.com"

type HelloTask struct {
	Name string
}

// The Execute method is required to run a task.
func (t *HelloTask) Execute(context.Context) error {
	slog.Info("Hello World!", "Name", t.Name)
	return nil
}

func main() {
	taskClient := clientFromConfig(serverURL, os.Getenv("TILEBOX_API_KEY"))
	runner := workflows.NewTaskRunner(taskClient)

	err := runner.RegisterTasks(
		&HelloTask{},
	)
	if err != nil {
		slog.Error("failed to register task", "error", err)
		return
	}

	runner.Run(context.Background())
}

func clientFromConfig(serverURL, authToken string) workflowsv1connect.TaskServiceClient {
	return workflowsv1connect.NewTaskServiceClient(
		grpc.RetryHTTPClient(), serverURL, connect.WithInterceptors(
			grpc.NewAddAuthTokenInterceptor(func() string {
				return authToken
			})),
	)
}
```
