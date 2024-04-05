# Example Workflow

An example of how to use the task runner library (see [../services/README.md](../services/README.md)) to implement
a workflow in go.

The actual tasks are implemented in `example_workflow.go`. `submitter/main.go` can be used to submit a workflow
to the processing service and `runner/main.go` can be used to start a task runner.
