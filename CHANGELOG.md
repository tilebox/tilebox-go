# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.17.0] - 2026-10-27

### Added

- `datasets`: Define enum fields with `field.Enum(name, values...).Queryable()` and filter them using public string values with `query.Field(...).Equal` or `NotEqual`. Schema updates may append enum values while preserving existing numbers.
- `datasets`: Resolve public enum values to protobuf numbers with `DatapointDescriptor.EnumNumber` for dynamic ingestion. Dynamic decoding returns public strings for scalar and repeated enum fields; absent fields remain absent, while zero and unknown enum numbers return errors.
- `workflows`: Submit workflow-scoped root tasks with `job.WithWorkflowID`. Omission remains unscoped; an explicit zero UUID is rejected.
- `workflows`: Preserve automation task scope through `AutomationTaskPrototype.WorkflowID`. Nil means unscoped and clears scope on replacement updates.
- `workflows`: Filter jobs with `job.WithWorkflowReferences` and read optional task attribution through `JobTask.WorkflowReference`.
- `workflows`: Expose workflow `ID`, `ProjectSlug`, and presence-aware `LegacyTaskAliases`. Set the creation policy with `WithLegacyTaskAliases`; omission preserves the server's legacy-compatible default.

### Changed

- **Breaking:** Rename the protobuf-enum field helper from `field.Enum(name, protobufValue)` to `field.ProtoEnum(name, protobufValue)`. `field.Enum` now accepts public string values.
- **Breaking:** Remove generated `FieldQueryValue.EnumName` accessors and builder field. Use `StringValue` with a declared public value for custom enum comparisons.
- **Breaking:** Rename `Workflow.Slug` to `Reference`, matching remote workflow selectors. `ProjectSlug` is the portable project name, not a remote selector. Workflow RPC requests use `WorkflowReference`.
- **Breaking:** `WorkflowService.CreateWorkflow` takes an additional `*bool` alias-policy argument. Pass nil to preserve omission.
- **Breaking:** Generated release task definitions and worker discovery use `UnboundTaskIdentifier` and `UnboundTaskIdentifiers`. The high-level `TaskIdentifier` interface and `ReleaseContent.Tasks` remain unchanged and unbound.
- Update Connect-Go, gRPC, and other dependencies.

### Fixed

- `datasets`: Select the datapoint message using the schema's `type_url`, independently of file and declaration order. Descriptors without `type_url` retain the previous first-message fallback.

## [0.16.0] - 2026-10-13

### Added

- `workflows`: Add `Automations.Create`, `Update`, and `Delete` for cron and storage-event automations. Update replaces the complete definition, including triggers; delete optionally cancels jobs created by the automation.
- Configure retries with `client.WithRetryPolicy`. Defaults remain four retries (five attempts), exponential jitter starting at 10–20 ms, and a 10-second maximum wait. Transient network failures and HTTP 429/5xx are retried, except Connect `resource_exhausted` responses. `client.RetryPolicy{}` disables retries; non-replayable request bodies are not retried.
- Unix socket connections now use the same configurable retry policy as HTTP connections.
- Customize TLS, proxies, or connection pools with `client.NewHTTPTransport` and `client.WithHTTPTransport` while retaining SDK retries. `client.WithHTTPClient` replaces all HTTP behavior, including retries; `client.WithTransport` accepts a native Connect transport and bypasses HTTP configuration entirely.
- Use `client.NewRetryHTTPClient(httpClient, policy)` to apply SDK retries to a custom HTTP client.
- `workflows`: Workflow results include `Summary` and optional, read-only `Presentation` metadata (image URL, alt text, imagery types, use cases, and dataset labels). Set the summary on create with `workflows.WithSummary` or on update with `workflow.WithSummary`. Create omits the Markdown description unless `workflows.WithDescription` is supplied. Update leaves omitted fields unchanged; an empty summary or description clears that field.

### Changed

- **Breaking:** Upgraded to Connect-Go v2. Consumers inspecting RPC errors must import `connectrpc.com/connect/v2`. Generated RPC methods accept and return protobuf messages directly; use `connect.NewClientContext` to set request headers and read response headers or trailers.
- **Breaking:** Accounts, Datasets, and Workflows constructors now require a shared `*connect.Client`. For example, replace `datasets.NewClient(datasets.WithAPIKey(key))` with `datasets.NewClient(client.New(client.WithAPIKey(key)))`. Reuse the same `client.New(...)` result across domains to share connections and interceptors. The default endpoint and `TILEBOX_API_KEY` lookup are unchanged.
- **Breaking:** Move `WithURL`, `WithAPIKey`, `WithClientMetadata`, and `WithHTTPClient` options from domain packages to `client`. Replace `WithConnectClientOptions` with `client.WithTransportOptions(...connecthttp.Option)` and `client.WithInterceptors(...connect.ClientInterceptor)`.
- **Breaking:** Client construction no longer enables Tilebox trace/log export. Call `workflows.ConfigureTileboxTelemetry(ctx, apiURL, apiKey)` before constructing domain clients to enable it. Existing global OpenTelemetry providers still supply client tracing. `WithDisableTracing` remains available; `WithDisableLogging` is removed.
- **Breaking:** Default HTTP clients no longer follow redirects. Configure the final API URL with `client.WithURL`; a redirect now fails the RPC instead of changing its destination or request method. Custom clients supplied through `client.WithHTTPClient` retain their own redirect policy.

## [0.15.1] - 2026-10-09

### Changed

- `workflows`: Rename `StorageTypeS3` to `StorageTypeAWSS3` and `StorageTypeFS` to `StorageTypeLocal`. Storage type strings are now `aws-s3`, `gcs`, `azure-blob`, and `local`.
- `workflows`: Replace `FilesystemReference` and `StorageLocationReference.Filesystem` with `LocalReference` and `StorageLocationReference.Local`, matching the regenerated API protos.

### Added

- `workflows`: Support `StorageSubscriptionTypeTileboxCLI` without provider delivery configuration.
- `workflows`: Add `StorageLocations.NotifyObjectCreated` to send file-ready notifications to Tilebox CLI subscriptions without automatic retries in the default transport.

## [0.15.0] - 2026-10-05

### Added

- `workflows`: Added `Client.StorageLocations` with location create/get/list/update/delete, notification subscription create/get/list/delete, and paginated subscription-event listing. Supports S3, GCS, Azure Blob, and filesystem references.
- `protogen`: Added the server-streaming `WorkerService.WatchLogs` RPC for structured runtime logs.

### Changed

- `workflows`: Moved `Automations.GetStorageLocation` and `Automations.ListStorageLocations` to `StorageLocations.Get` and `StorageLocations.List`, using the new storage-location service.
- `workflows`: Replaced the flat `StorageLocation.Location` and `StorageLocation.Type` fields with `Name` and `Reference`. Read the type from `Reference.Type` and provider coordinates from `Reference.AWSS3Bucket`, `GCSBucket`, `AzureBlob`, or `Filesystem`. Automation storage triggers use the same updated model.
- Updated OpenTelemetry and other dependencies.

## [0.14.0] - 2026-09-23

### Added

- `workflows`: Added paginated owned and public workflow listing with `ListPage`, `ListPublic`, and `ListPublicPage`.

### Changed

- `workflows`: `List` and `ListPublic` return lazy auto-paginating iterators with cursor and total-limit options, consistent with other paginated queries. Use `workflows.Collect` to collect results into a slice.
- `workflows`: `UndeployRelease` now accepts cluster slugs followed by options instead of a positional release ID. Omit `workflow.WithReleaseID` to remove all deployed releases on the selected clusters, or provide it to remove a specific release.

## [0.13.0] - 2026-09-04

### Added

- Added a structured `Tilebox-Client` header with automatically detected SDK, runtime, OS, execution environment, invoker, and cloud metadata. SDK wrappers can replace the metadata with `WithClientMetadata` and `client.NewMetadata`.
- `datasets`: Added `GEOMETRY_CONTAINS_FILTER` spatial filter mode to datapoint queries.

### Removed

- `datasets`: Removed the legacy client metadata from dataset list requests.

### Changed

- `workflows`: Change `Client.NewPollingTaskRunner` to accept a resolved `*Cluster`, executor, and logger directly; it no longer fetches a cluster or accepts task-runner options.
- `workflows`: Log recoverable task polling and result-reporting failures as warnings while they await retry.

### Fixed

- Retry recoverable HTTP transport failures, including EOFs, timeouts, and connection resets, with five total attempts and bounded exponential jitter.
- Allow polling task runners to finish graceful shutdown after reporting an active task's pending result.

## [0.12.0] - 2026-09-01

### Added

- `workflows`: Added `client.Jobs.ListTasks` and `client.Jobs.ListTasksPage` for listing paginated task-tree levels, selecting a parent task, and prefetching child pages.
- `workflows`: Added log severity filtering with `job.WithLogSeverityGroups(...)` and task-specific log and span filtering with `job.WithTaskID(...)`.

### Changed

- `workflows`: Removed `Job.TaskSummaries`; job tasks are now retrieved through the paginated task-tree API.

## [0.11.1] - 2026-07-27

### Fixed

- `datasets`: Resolve registered protobuf dependencies, including STAC types, when constructing dynamic datapoint descriptors.

## [0.11.0] - 2026-07-23

### Added

- `accounts`: Added account details, active plan, and usage report clients.
- `datasets`: Added source JSON pointers, queryable metadata, JSON Schema references, semantic roles, and well-known protobuf message and enum fields to dataset creation and updates, including generated STAC types.
- `datasets`: Added fluent Boolean, string, and numeric expressions for filtering datapoints by custom queryable fields.

### Changed

- Updated the minimum supported Go version to 1.26.

## [0.10.0] - 2026-07-01

### Added

- `workflows`: Added `client.Clusters.Update`.
- `workflows`: Added `client.Workflows.Update`.
- `workflows`: Added cluster descriptions.
- `workflows`: Added deployed cluster details on workflow releases.

## [0.9.0] - 2026-06-25

### Added

- `workflows`: Added `client.Workflows.Delete` and `client.Workflows.UnpublishRelease`.
- `workflows`: Added `job.WithClusterSlugs(...)` to filter jobs by cluster slug (any task in that cluster).

### Fixed

- `workflows`: Fixed `PollingTaskRunner` retry handling so request-shaped `NextTask` and `TaskFailed` RPC errors do not cause endless retry loops.

## [0.8.0] - 2026-06-05

### Added

- `workflows`: Added `TaskExecutor`, `PollingTaskRunner`, and `client.NewPollingTaskRunner` for custom task execution backends that need to use the Tilebox task polling protocol directly.

### Changed

- `workflows`: Changed `TaskRunner.RunForever` and `TaskRunner.RunAll` to return errors from the polling loop.

## [0.7.1] - 2026-05-22

### Fixed

- `workflows`: Fixed Windows builds by using platform-specific task runner shutdown signals.

## [0.7.0] - 2026-05-22

### Added

- `datasets`: Added `client.Datasets.Create` and `client.Datasets.Update` with optional summary and markdown description options.
- `datasets`: Added datapoint query pagination and dynamic datapoint decoding helpers for converting raw protobuf datapoints to maps.
- `examples`: Added a dynamic dataset query example that decodes datapoints without generated Go types.
- `workflows`: Added `client.Automations` to list and inspect automations and storage locations.

## [0.6.0] - 2026-05-20

### Added

- `workflows`: Added cursor and limit query options plus `QueryPage`, `QueryLogsPage`, and `QuerySpansPage` for manual pagination of job, log, and span queries.
- `workflows`: Added `job.WithTaskStates(...)` to filter job queries by task state.
- `workflows`: Added string and JSON representations for job and task states.

### Changed

- `workflows`: Changed job, log, and span queries to support cursor-based pagination while keeping existing sequence APIs auto-paginated.
- `workflows`: Moved telemetry query options to the `workflows/v1/job` package so job and telemetry queries share `job.WithLimit`, `job.WithCursor`, and `job.WithSortDirection`.

## [0.5.0] - 2026-05-11

### Added

- `workflows`: Added `client.Jobs.QueryLogs` and `client.Jobs.QuerySpans` to query logs and trace spans for a job, plus a telemetry query example.
- `workflows`: Added `ConfigureConsoleLogging` to enable console log output that composes with the Tilebox OpenTelemetry log exporter.
- `workflows`: Added `WithSpan` and `WithSpanResult` helpers to start spans from the current task execution context without manually passing a tracer.

### Changed

- `workflows`: Correlate context-aware task logs with traces by adding `trace_id`, `span_id`, and `task_id` attributes and recording log messages as span events.
- `examples`: Updated workflow and dataset examples to use context-aware `slog` methods.

## [0.4.0] - 2026-03-06

### Added

- `datasets`: Added `WithCollections(...)` and `WithCollectionIDs(...)` query options for datapoint queries.

### Changed

- `datasets`: Updated `client.Datapoints.GetInto`, `client.Datapoints.Query`, and `client.Datapoints.QueryInto` to take a `datasetID` as the primary identifier, with collection filtering now configured via query options.

## [0.3.2] - 2026-02-26

### Added

- `workflows`: Added OTEL metrics for the task runner, tracking `task.executed.count`, `task.computed.count`, `task.failed.count`, `task.input.size` and `task.execution.duration`.

## [0.3.1] - 2026-02-25

### Fixed

- `workflows`: Fixed a bug where the task runner could panic when a task failed.

## [0.3.0] - 2026-02-19

### Added

- `workflows`: Added `WithOptional` option to `workflows.SubmitSubtask` to mark a subtask as optional.

### Changed

- `datasets`: Changed `client.Datasets.Create()` to `client.Datasets.CreateOrUpdate()`.

## [0.2.1] - 2025-12-05

### Fixed

- `datasets`: Fixed stack overflow in `client.Datapoints.DeleteIDs`.

## [0.2.0] - 2025-12-04

### Added

- `workflows`: Added `ExecutionStats` to the `Job` object to provide programmatic access to a job's execution
  statistics.
- `workflows`: Added query filters to the `client.Jobs.Query` method to filter jobs by multiple automation ids,
  job state, and job name.
- `workflows`: Added additional `job.State` values to indicate a job's current state and progress more accurately.
- `workflows`: Removed the restriction of `64` subtasks per task.
- `datasets`: Added `client.Datasets.Create()` method to create a new dataset.

### Changed

- `tilebox-workflows`: Switched to an updated internal `TaskSubmission` message format that allows for more efficient
  submission of a very large number of tasks.

## [0.1.1] - 2025-10-30

### Fixed

- Fixed error logging when trying to extend a task lease after a context cancellation.

## [0.1.0] - 2025-10-23

### Added

- Added support for Tilebox Datasets, including operations for datasets, collections, and datapoints.
- Added support for Tilebox Workflows, including operations for runners, jobs, tasks, and clusters.
- Added support for Tilebox Observability, including logging and tracing helpers.
- Added examples for using the library.

[Unreleased]: https://github.com/tilebox/tilebox-go/compare/v0.17.0...HEAD
[0.17.0]: https://github.com/tilebox/tilebox-go/compare/v0.16.0...v0.17.0
[0.16.0]: https://github.com/tilebox/tilebox-go/compare/v0.15.1...v0.16.0
[0.15.1]: https://github.com/tilebox/tilebox-go/compare/v0.15.0...v0.15.1
[0.15.0]: https://github.com/tilebox/tilebox-go/compare/v0.14.0...v0.15.0
[0.14.0]: https://github.com/tilebox/tilebox-go/compare/v0.13.0...v0.14.0
[0.13.0]: https://github.com/tilebox/tilebox-go/compare/v0.12.0...v0.13.0
[0.12.0]: https://github.com/tilebox/tilebox-go/compare/v0.11.1...v0.12.0
[0.11.1]: https://github.com/tilebox/tilebox-go/compare/v0.11.0...v0.11.1
[0.11.0]: https://github.com/tilebox/tilebox-go/compare/v0.10.0...v0.11.0
[0.10.0]: https://github.com/tilebox/tilebox-go/compare/v0.9.0...v0.10.0
[0.9.0]: https://github.com/tilebox/tilebox-go/compare/v0.8.0...v0.9.0
[0.8.0]: https://github.com/tilebox/tilebox-go/compare/v0.7.1...v0.8.0
[0.7.1]: https://github.com/tilebox/tilebox-go/compare/v0.7.0...v0.7.1
[0.7.0]: https://github.com/tilebox/tilebox-go/compare/v0.6.0...v0.7.0
[0.6.0]: https://github.com/tilebox/tilebox-go/compare/v0.5.0...v0.6.0
[0.5.0]: https://github.com/tilebox/tilebox-go/compare/v0.4.0...v0.5.0
[0.4.0]: https://github.com/tilebox/tilebox-go/compare/v0.3.2...v0.4.0
[0.3.2]: https://github.com/tilebox/tilebox-go/compare/v0.3.1...v0.3.2
[0.3.1]: https://github.com/tilebox/tilebox-go/compare/v0.3.0...v0.3.1
[0.3.0]: https://github.com/tilebox/tilebox-go/compare/v0.2.0...v0.3.0
[0.2.1]: https://github.com/tilebox/tilebox-go/compare/v0.2.0...v0.2.1
[0.2.0]: https://github.com/tilebox/tilebox-go/compare/v0.1.1...v0.2.0
[0.1.1]: https://github.com/tilebox/tilebox-go/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/tilebox/tilebox-go/releases/tag/v0.1.0
