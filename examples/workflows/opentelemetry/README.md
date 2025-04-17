# OpenTelemetry Observability

How to set up tracing and logging for workflows using OpenTelemetry.

It uses [OpenTelemetry](https://opentelemetry.io/) to add tracing and logging and the standard [slog](https://pkg.go.dev/log/slog) package for logging.

- [submitter/main.go](submitter/main.go) sets up tracing and logging and submits a job.
- [runner/main.go](runner/main.go) sets up tracing and logging for a runner.
- [task.go](task.go) shows how to use logs and traces within a task.
