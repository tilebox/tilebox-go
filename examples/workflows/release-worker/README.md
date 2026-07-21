# Go Release Worker

This is a complete entry point for a precompiled Go workflow release. It registers tasks and serves the execution-only
`WorkerService` expected by `tilebox runner start`.

Unlike a direct `TaskRunner`, this program does **not** create a client, select a cluster, poll the task queue, extend
leases, or report results. The parent Tilebox runner owns those responsibilities and launches this binary with a private
Unix socket in `TILEBOX_WORKER_ADDRESS`. It lists the registered tasks before initialization, then supplies the API and
release context used while executing tasks.

The binary is intended to be built by a workflow release `[build.go]` configuration rather than started directly:

```toml
[build.go]
package = "./cmd/worker"
```

See [main.go](main.go) for task registration, progress-compatible context usage, and subtask submission.
