# MapReduce Workflow

Example workflow that fans out many map tasks and then submits reduce tasks that depend on the map tasks.

- [main.go](main.go) contains the task definitions, submits a job, and then starts a task runner to execute it.

Run the example:

```bash
go run ./examples/workflows/mapreduce
```
