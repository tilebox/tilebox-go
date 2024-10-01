# Examples

This repository contains examples for using the Tilebox Go library.

## Environment variables

To run the examples, `TILEBOX_API_KEY` environment variable needs to be set to your Tilebox API key.
Head over to [Tilebox Console](https://console.tilebox.com/settings/api-keys) if you don't have one already.

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

- [Hello World](helloworld): How to submit a task and run it.
- [Sample Workflow](sampleworkflow): A example of how to submit a task and run a workflow using protobuf messages and with OpenTelemetry setup.

## Datasets examples

- [Load Data](load): How to load data from a collection.
