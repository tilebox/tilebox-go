# How to use Protobuf tasks

This example uses a protobuf message to define `SpawnWorkflowTreeTask` task.

The task definition is defined in [workflow.proto](../../../apis/examples/v1/workflow.proto) and the task implementation is in [example_workflow.go](tasks.go).

The generated code of `workflow.proto` is in [workflow.pb.go](../../../protogen/go/examples/v1/workflow.pb.go) and can be (re-)generated using the following command:

```bash
go generate ./...
```

- [submitter/main.go](submitter/main.go) submits a workflow to Tilebox Workflows with a protobuf root task.
- [runner/main.go](runner/main.go) starts a task runner that can execute both protobuf and regular struct tasks.
- [tasks.go](tasks.go) shows the protobuf task implementation and how to submit protobuf sub-tasks.
