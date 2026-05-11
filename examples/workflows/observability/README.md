# Workflow Observability

Tilebox-native workflow observability is automatically configured by `workflows.NewClient()` for the standard examples.

This folder contains examples for querying workflow telemetry and for exporting telemetry to custom third-party backends:

- [query](query): Query logs and traces for an existing workflow job.
- [axiom](axiom): Export workflow traces and logs to [Axiom](https://axiom.co/).
- [opentelemetry](opentelemetry): Export workflow traces and logs to an OpenTelemetry-compatible endpoint.
