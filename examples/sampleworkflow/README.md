# Example Workflow

An example of how to use Tilebox Workflows to implement a workflow in Go.

- `submitter/main.go` submits a workflow to Tilebox Workflows.
- `runner/main.go` starts a task runner.

## Protobuf

This example uses a protobuf message to define `SpawnWorkflowTreeTask` task.

The task definition is defined in [workflow.proto](../../apis/examples/v1/workflow.proto) and the task implementation is in [example_workflow.go](example_workflow.go).

The generated code of `workflow.proto` is in [workflow.pb.go](../../protogen/go/examples/v1/workflow.pb.go) and can be (re-)generated using the following command:

```bash
go generate ./...
```

## OpenTelemetry

This example uses OpenTelemetry to add tracing and logging.

To run the example, you need to set the additional following environment variables:

- `AXIOM_API_KEY`to your Axiom API key.
- `AXIOM_TRACES_DATASET` to the dataset where traces should be stored.
- `AXIOM_LOGS_DATASET` to the dataset where logs should be stored.
