# How to use Protobuf tasks

This example uses a protobuf message to define `SpawnWorkflowTreeTask` task.

The task definition is defined in [workflow.proto](../../../apis/examples/v1/workflow.proto) and the task implementation is in [main.go](main.go).

The generated code of `workflow.proto` is in [workflow.pb.go](../../../protogen/examples/v1/workflow.pb.go) and can be (re-)generated using the following command:

```bash
go generate ./...
```

- [main.go](main.go) submits a workflow to Tilebox Workflows, starts a task runner to execute it, and shows the protobuf task implementation and how to submit protobuf sub-tasks.

Run the example:

```bash
go run ./examples/workflows/protobuf-task
```
