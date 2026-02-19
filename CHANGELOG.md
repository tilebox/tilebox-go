# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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

[Unreleased]: https://github.com/tilebox/tilebox-go/compare/v0.3.0...HEAD
[0.2.1]: https://github.com/tilebox/tilebox-go/compare/v0.2.0...v0.3.0
[0.2.1]: https://github.com/tilebox/tilebox-go/compare/v0.2.0...v0.2.1
[0.2.0]: https://github.com/tilebox/tilebox-go/compare/v0.1.1...v0.2.0
[0.1.1]: https://github.com/tilebox/tilebox-go/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/tilebox/tilebox-go/releases/tag/v0.1.0
