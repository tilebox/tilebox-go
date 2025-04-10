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

For example, to run the `helloworld/submitter` example:

```bash
go run ./helloworld/submitter
```

## Workflows examples

- [Hello world](helloworld): How to submit a task and run it.
- [Sample workflow](sampleworkflow): How to submit a task and run a workflow using protobuf messages and with OpenTelemetry setup.

## Datasets examples

- [Querying data](query): How to query datapoints from dataset collections.
- [Ingest and delete data](ingest): How to create a collection, ingest datapoints, and then delete them.
