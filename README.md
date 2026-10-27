<h1 align="center">
  <img src="https://storage.googleapis.com/tbx-web-assets-2bad228/banners/tilebox-banner.svg" alt="Tilebox Logo">
  <br>
</h1>

<div align="center">
  <a href="https://pkg.go.dev/github.com/tilebox/tilebox-go">
    <img src="https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square&color=f43f5e" alt="Go.dev reference badge"/>
  </a>
  <a href="https://github.com/tilebox/tilebox-go/blob/main/LICENSE">
    <img src="https://img.shields.io/github/license/tilebox/tilebox-go.svg?style=flat-square&color=f43f5e" alt="MIT License"/>
  </a>
  <a href="https://github.com/tilebox/tilebox-go/actions">
    <img src="https://img.shields.io/github/actions/workflow/status/tilebox/tilebox-go/main.yml?style=flat-square&color=f43f5e" alt="Build Status"/>
  </a>
  <a href="https://tilebox.com/discord">
    <img src="https://img.shields.io/badge/Discord-%235865F2.svg?style=flat-square&logo=discord&logoColor=white" alt="Join us on Discord"/>
  </a>
</div>

<p align="center">
  <a href="https://docs.tilebox.com/introduction"><b>Documentation</b></a>
  |
  <a href="https://console.tilebox.com/"><b>Console</b></a>
  |
  <a href="https://examples.tilebox.com/"><b>Example Gallery</b></a>
</p>

# Tilebox Go

Go library for [Tilebox](https://tilebox.com), a lightweight space data management and orchestration software - on ground and in orbit.

## Installation

Run the following command to add the library to your project:

```bash
go get github.com/tilebox/tilebox-go
```

For Tilebox datasets type generation, you will need to install [tilebox-generate](https://github.com/tilebox/tilebox-generate) command-line tool.

## Examples

For examples on how to use the library, see the [examples](examples) directory.

## Usage

### Configuring RPC clients

Accounts, Datasets, and Workflows constructors require a `*connect.Client`. Create it with `client.New()`;
the defaults are `https://api.tilebox.com`, `TILEBOX_API_KEY`, SDK metadata, and binary Connect protobuf.
HTTPS negotiates HTTP/2 with HTTP/1.1 fallback. Plain HTTP and Unix sockets use HTTP/1.1 by default.
Reuse one RPC client across domain clients to share its connection pool and interceptor chain.

A CLI can configure TLS, a proxy, and retries once for all three APIs without implementing its own retry loop:

```go
import (
	"crypto/tls"
	"net/http"
	"net/url"
	"time"

	"github.com/tilebox/tilebox-go/accounts/v1"
	"github.com/tilebox/tilebox-go/client"
	"github.com/tilebox/tilebox-go/datasets/v1"
	"github.com/tilebox/tilebox-go/workflows/v1"
)

transport := client.NewHTTPTransport()
transport.TLSClientConfig = &tls.Config{MinVersion: tls.VersionTLS12}
proxyURL, err := url.Parse("http://localhost:8080")
if err != nil {
	return err
}
transport.Proxy = http.ProxyURL(proxyURL)
rpcClient := client.New(
	client.WithHTTPTransport(transport),
	client.WithRetryPolicy(client.RetryPolicy{
		MaxRetries: 2,
		MinWait: 100 * time.Millisecond,
		MaxWait: 5 * time.Second,
	}),
)
datasetsClient := datasets.NewClient(rpcClient)
workflowsClient := workflows.NewClient(rpcClient)
accountsClient := accounts.NewClient(rpcClient)
```

`NewHTTPTransport` creates an independent pool without changing `http.DefaultTransport`. Configure it before use.
`client.WithHTTPTransport` uses the supplied transport unchanged while preserving SDK retries. For Unix sockets,
a custom transport must provide its own socket dialer.

Responses have no size limit by default. To bound memory use, set a per-message limit,
measured after decompression:

```go
import (
	"connectrpc.com/connect/v2/connecthttp"
	"github.com/tilebox/tilebox-go/client"
	"github.com/tilebox/tilebox-go/datasets/v1"
)

rpcClient := client.New(client.WithTransportOptions(connecthttp.WithReadMaxBytes(64 << 20)))
datasetsClient := datasets.NewClient(rpcClient)
```

URL, API key, metadata, HTTP client, native transport, interceptor, and retry options now belong only to `client`.
`client.WithTransportOptions` accepts `connecthttp.Option` values for protocol, compression, and message limits.
`client.WithInterceptors` accepts native `connect.ClientInterceptor` functions in outermost-first order. Both replace the
former `WithConnectClientOptions` API. The SDK applies its configured API key and client metadata after custom
interceptors. An empty API key leaves caller-provided authentication intact.

For per-call headers, use `connect.NewClientContext(ctx)` and its `CallInfo.RequestHeader()`. Read response headers
and trailers from the same `CallInfo` after the call. Use a separate client context for each concurrent call.

The default HTTP client retries replayable requests on transient transport failures, HTTP 429 (except Connect
`resource_exhausted`), and HTTP 5xx, with five total attempts. Retries reuse Connect's request body without a second
payload buffer. Streaming request bodies are not retried. A lost response can still mean
the server processed the request; retries do not guarantee exactly-once execution. Use call deadlines to bound retries.
The default policy permits four retries, starts with `MinWait` of 20 ms, and caps exponential equal-jitter waits at
10 seconds. Unix sockets use the same retry policy as HTTP. `client.WithRetryPolicy(client.RetryPolicy{})` disables
retries. Default clients do not follow redirects. `client.WithHTTPClient` replaces the entire HTTP client, including
retries, redirects, and timeouts; `WithHTTPTransport` and `WithRetryPolicy` then have no effect.
`client.WithTransport` accepts a native `connect.Transport` and bypasses all HTTP setup, including URL, retry policy,
and transport options. Authentication, metadata, and client interceptors still apply.

Domain clients retain `WithDisableTracing`; `WithDisableLogging` has been removed. Constructors do not install
telemetry exporters. To export to Tilebox, call `workflows.ConfigureTileboxTelemetry(ctx, apiURL, apiKey)` before
constructing domain clients, which capture the global tracer provider. Console logging remains a separate choice.

Storage notifications use the same RPC client as other Workflows operations:

```go
workflowsClient := workflows.NewClient(rpcClient)
err := workflowsClient.StorageLocations.NotifyObjectCreated(ctx, subscriptionID, "path/to/object", eventTime)
```

Notifications share the configured transport, authentication, metadata, interceptors, and retry policy, including
support for in-process transports. A zero `eventTime` omits the timestamp. The server must implement the
`StorageLocationService.NotifyObjectCreated` RPC; the SDK no longer calls the JSON notification endpoint.
There is no separate notification timeout. Use a context deadline to bound the whole call, including retries.
Retries may create duplicate jobs if a response is lost.

### Filtering Dataset Queries

Fields marked queryable in a dataset schema can be filtered with fluent Boolean, string, and numeric expressions. Multiple
expressions passed to `WithFilters` are combined with each other and with temporal and spatial filters using logical
AND.

```go
package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/tilebox/tilebox-go/client"
	"github.com/tilebox/tilebox-go/datasets/v1"
	"github.com/tilebox/tilebox-go/query"
)

func main() {
	ctx := context.Background()
	client := datasets.NewClient(client.New())
	datasetID := uuid.MustParse("019c0123-4567-7890-abcd-ef0123456789")
	start := time.Date(2026, time.July, 1, 0, 0, 0, 0, time.UTC)
	end := start.Add(24 * time.Hour)

	datapoints, err := datasets.Collect(client.Datapoints.Query(ctx,
		datasetID,
		datasets.WithTemporalExtent(query.NewTimeInterval(start, end)),
		datasets.WithFilters(
			query.Field("eo_cloud_cover").LessThan(20.0),
			query.Field("granule_name").Equal("S2A_GRANULE"),
			query.Or(
				query.Field("quality").GreaterThanOrEqual(80),
				query.Field("quality").IsNull(),
			),
		),
	))
	if err != nil {
		slog.ErrorContext(ctx, "Failed to query datapoints", slog.Any("error", err))
		return
	}

	slog.InfoContext(ctx, "Found datapoints", slog.Int("count", len(datapoints)))
}
```

Comparisons with absent fields evaluate to unknown and therefore do not match. Use `IsNull`, as above, when absent values
should be included explicitly. `NotEqual` also excludes absent values.

### Enum Fields

Use `field.Enum` to define a finite set of public string values:

```go
fields := []datasets.Field{
	field.Enum("processing_level", "L1", "L2-SP").Queryable(),
}
```

Pass these fields to `Datasets.Create`, `Update`, or `CreateOrUpdate`. Values must be nonempty and unique.
Updates may append values, but cannot remove, rename, or reorder them. The server assigns protobuf symbols and numbers;
list positions are not enum numbers. Queryable enums do not count toward the queryable string field limit.
Use `field.ProtoEnum` for supported built-in protobuf enum types. This replaces the previous `field.Enum(name, protobufValue)` API.

Filter custom enums with case-sensitive public strings. Only `Equal` and `NotEqual` are supported; the server rejects
unknown values, numeric literals, and ordering comparisons. Combine equality comparisons with `query.Or` to match
several values:

```go
filter := query.Or(
	query.Field("processing_level").Equal("L1"),
	query.Field("processing_level").Equal("L2-SP"),
)
```

`datasets.UnmarshalDatapoint` returns custom enums as their public strings, including repeated enum values.
Absent fields remain absent. Zero and unknown enum numbers have no public value and produce a decoding error.
Built-in enum decoding is unchanged.

For ingestion with `dynamicpb`, resolve the public value using the returned dataset schema:

```go
descriptor, err := datasets.NewDatapointDescriptor(dataset)
if err != nil {
	return err
}
number, err := descriptor.EnumNumber("processing_level", "L2-SP")
if err != nil {
	return err
}
datapoint := dynamicpb.NewMessage(descriptor.MessageDescriptor)
datapoint.Set(descriptor.MessageDescriptor.Fields().ByName("processing_level"), protoreflect.ValueOfEnum(number))
datapoint.Set(descriptor.MessageDescriptor.Fields().ByName("time"), protoreflect.ValueOfMessage(timestamppb.Now().ProtoReflect()))
datapoints := []*dynamicpb.Message{datapoint}
_, err = client.Datapoints.Ingest(ctx, collection.ID, &datapoints, false)
```

This uses `dynamicpb`, `protoreflect`, and `timestamppb` from `google.golang.org/protobuf`. Set any other required datapoint
fields before ingestion. Generated protobuf messages remain supported by the same `Ingest` method.

### Writing a Task

Here we define a simple task that logs "Hello World!":

```go
package helloworld

import (
	"context"
	"log/slog"

	"github.com/tilebox/tilebox-go/workflows/v1"
)

type HelloTask struct {
	Name string // You can add any fields you need to the task struct.
}

// The Execute method isn't needed to submit a task but is required to run a task.
func (t *HelloTask) Execute(ctx context.Context) error {
	slog.InfoContext(ctx, "Hello World!", slog.String("Name", t.Name))
	return nil
}

// The Identifier method is optional and will be generated if not provided.
func (t *HelloTask) Identifier() workflows.TaskIdentifier {
	return workflows.NewTaskIdentifier("hello-world", "v1.0")
}
```

### Submitting a Job

Here we create a Workflows client and submit a job with a single task:

```go
package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/tilebox/tilebox-go/client"
	"github.com/tilebox/tilebox-go/workflows/v1"
)

type HelloTask struct {
	Name string
}

func main() {
	ctx := context.Background()
	workflows.ConfigureConsoleLogging(slog.LevelInfo)
	workflows.ConfigureTileboxTelemetry(ctx, "https://api.tilebox.com", os.Getenv("TILEBOX_API_KEY"))
	client := workflows.NewClient(client.New())

	job, err := client.Jobs.Submit(ctx, "hello-world",
		[]workflows.Task{
			&HelloTask{
				Name: "Tilebox",
			},
		},
	)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to submit job", slog.Any("error", err))
		return
	}

	slog.InfoContext(ctx, "Job submitted", slog.String("job_id", job.ID.String()))
}
```

`workflows.ConfigureConsoleLogging` adds console output without replacing the Tilebox exporter, so logs are written to both destinations when both are configured.

### Running a Worker

Here we create a TaskRunner and run a worker that is capable of executing `HelloTask` tasks:

```go
package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/tilebox/tilebox-go/client"
	"github.com/tilebox/tilebox-go/workflows/v1"
)

type HelloTask struct {
	Name string
}

// The Execute method is required to run a task.
func (t *HelloTask) Execute(ctx context.Context) error {
	slog.InfoContext(ctx, "Hello World!", slog.String("Name", t.Name))
	return nil
}

func main() {
	ctx := context.Background()
	workflows.ConfigureConsoleLogging(slog.LevelInfo)
	workflows.ConfigureTileboxTelemetry(ctx, "https://api.tilebox.com", os.Getenv("TILEBOX_API_KEY"))
	client := workflows.NewClient(client.New())

	runner, err := client.NewTaskRunner(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "failed to create task runner", slog.Any("error", err))
		return
	}

	err = runner.RegisterTasks(&HelloTask{})
	if err != nil {
		slog.ErrorContext(ctx, "failed to register tasks", slog.Any("error", err))
		return
	}

	if err := runner.RunForever(ctx); err != nil {
		slog.ErrorContext(ctx, "task runner stopped", slog.Any("error", err))
	}
}
```

### Querying Job Telemetry

After a job has been executed, we can query runner logs associated with it.

```go
package main

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/tilebox/tilebox-go/client"
	"github.com/tilebox/tilebox-go/workflows/v1"
	"github.com/tilebox/tilebox-go/workflows/v1/job"
)

func main() {
	ctx := context.Background()
	client := workflows.NewClient(client.New())
	jobID := uuid.MustParse("019e070c-63ba-1c7e-5f1d-65be3e22d52a")

	for logRecord, err := range client.Jobs.QueryLogs(ctx, jobID,
		job.WithLimit(100),
		job.WithSortDirection(job.Ascending),
	) {
		if err != nil {
			slog.ErrorContext(ctx, "failed to query job logs", slog.Any("error", err))
			return
		}

		slog.InfoContext(ctx, "job log", slog.String("body", logRecord.Body))
	}
}
```

## License

Distributed under the MIT License (`The MIT License`).
