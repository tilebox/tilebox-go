# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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

[Unreleased]: https://github.com/tilebox/tilebox-go/compare/v0.12.0...HEAD
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
