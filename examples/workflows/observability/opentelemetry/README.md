# OpenTelemetry Observability

How to set up tracing and logging for workflows using OpenTelemetry.

It uses [OpenTelemetry](https://opentelemetry.io/) to add tracing and logging and the standard [slog](https://pkg.go.dev/log/slog) package for logging.

- [main.go](main.go) sets up tracing and logging, submits a job, starts a task runner to execute it, and shows how to use logs and traces within a task.

Run the example:

```bash
go run ./examples/workflows/observability/opentelemetry
```
