# Axiom Observability

How to set up tracing and logging for workflows using [Axiom](https://axiom.co/) observability platform.

It uses [OpenTelemetry](https://opentelemetry.io/) to add tracing and the standard [slog](https://pkg.go.dev/log/slog) package for logging.

- [main.go](main.go) sets up tracing and logging, submits a job, starts a task runner to execute it, and shows how to use logs and traces within a task.

To run the example, you need to set the additional following environment variables:

- `AXIOM_API_KEY` to your Axiom API key.
- `AXIOM_TRACES_DATASET` to the dataset where traces should be stored.
- `AXIOM_LOGS_DATASET` to the dataset where logs should be stored.

Run the example:

```bash
go run ./examples/workflows/observability/axiom
```
