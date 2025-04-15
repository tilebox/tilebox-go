# Axiom Observability

How to set up tracing and logging for workflows using [Axiom](https://axiom.co/) observability platform.

It uses [OpenTelemetry](https://opentelemetry.io/) to add tracing and the standard [slog](https://pkg.go.dev/log/slog) package for logging.

- [submitter/main.go](submitter/main.go) sets up tracing and logging and submits a job.
- [runner/main.go](runner/main.go) sets up tracing and logging for a runner.
- [task.go](task.go) shows how to use logs and traces within a task.

To run the example, you need to set the additional following environment variables:

- `AXIOM_API_KEY` to your Axiom API key.
- `AXIOM_TRACES_DATASET` to the dataset where traces should be stored.
- `AXIOM_LOGS_DATASET` to the dataset where logs should be stored.
