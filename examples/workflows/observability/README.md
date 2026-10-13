# Workflow Observability

The task-runner examples call `workflows.ConfigureTileboxTelemetry(ctx, apiURL, apiKey)` explicitly before
`workflows.NewClient(client.New())` captures the global tracer provider. Constructors do not install exporters.
The query example only reads existing telemetry and does not configure export.

This folder contains examples for querying workflow telemetry and for exporting telemetry to custom third-party backends:

- [query](query): Query logs and traces for an existing workflow job.
- [axiom](axiom): Export workflow traces and logs to [Axiom](https://axiom.co/).
- [opentelemetry](opentelemetry): Export workflow traces and logs to an OpenTelemetry-compatible endpoint.

The custom-backend examples configure their own providers first, then add Tilebox export. Tilebox setup composes
with existing SDK tracer providers and slog handlers. Omit that setup call to export only to the custom backend.
