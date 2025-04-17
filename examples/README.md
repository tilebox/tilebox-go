# Examples

This repository contains examples for using the Tilebox Go library.

## Environment variables

To run the examples, `TILEBOX_API_KEY` environment variable needs to be set to your Tilebox API key.
Head over to [Tilebox Console](https://console.tilebox.com/account/api-keys) if you don't have one already.

### Running the examples

Each example can be run using the following command:

```bash
go run ./<example-folder>
```

For example, to run the `workflows/helloworld/submitter` example:
x
```bash
go run ./workflows/helloworld/submitter
```

## Workflows examples

- [Hello world](workflows/helloworld): How to submit a task and run it.
- [Protobuf tasks](workflows/protobuf-task): How to use Protobuf tasks.
- [Axiom Observability](workflows/axiom): How to set up tracing and logging for workflows using [Axiom](https://axiom.co/) observability platform.
- [OpenTelemetry Observability](workflows/opentelemetry): How to set up tracing and logging for workflows using [OpenTelemetry](https://opentelemetry.io/).

## Datasets examples

- [Querying data](datasets/query): How to query datapoints from dataset collections.
- [Ingest and delete data](datasets/ingest): How to create a collection, ingest datapoints, and then delete them.
