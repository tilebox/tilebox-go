# AGENTS.md

## Commands
- Build: `go build ./...`
- Test all: `go test ./...`
- Test single: `go test -run TestName ./path/to/package`
- Lint: `golangci-lint run ./...` (or `golangci-lint run --fix ./...` when explicitly fixing lint)
- Generate protobuf/Connect code: `go generate ./...` (runs `buf generate` via `generate.go`)

## Repo Architecture
- `workflows/v1`: Workflows public client API (job submission/query, cluster lookup, task runner, task/subtask options).
- `datasets/v1`: Datasets public client API (datasets, collections, datapoint query/ingest/delete, schema field helpers).
- `query`: Reusable temporal/spatial query extent/filter builders used by dataset/workflow query APIs.
- `observability`: Shared tracing span wrappers and OpenTelemetry logger/tracer setup helpers.
- `internal/grpc`: Internal transport utilities (retrying HTTP client, auth interceptor, record/replay round trippers for tests).
- `internal/span`: Trace-context helper utilities.
- `protogen`: Generated protobuf, gRPC, and Connect stubs used by public clients.
- `apis`: Local proto sources in this repo (currently example protos; API service protos mostly come from external module input).
- `examples`: Runnable usage examples for workflows, datasets, OpenTelemetry, and protobuf-backed tasks.

## Design Patterns And Paradigms
- Public packages expose high-level clients (`NewClient`) with focused sub-clients (`Jobs`, `Datasets`, `Datapoints`, etc.) and hide transport/service details behind interfaces.
- Configuration follows the functional options pattern (`WithURL`, `WithAPIKey`, `WithHTTPClient`, tracing/logging options).
- RPC-facing logic is separated into service abstractions (`*Service` interfaces + concrete implementations) to keep client methods testable and composable.
- Errors are wrapped with context using `%w` and traced with `observability.WithSpan` / `WithSpanResult`.
- Query APIs favor lazy iteration (`iter.Seq2`) for paginated datapoint reads; convenience collectors exist when eager loading is preferred.
- Workflows task submission and execution are strongly typed around task identifiers/versioning, with default identifier derivation via reflection when not explicitly provided.
- Tests use table-driven style and property-based tests (`pgregory.net/rapid`) where useful.
- Integration-style tests rely on deterministic RPC record/replay fixtures in `*/testdata/recordings`.

## Protobuf And Codegen Rules
- Treat `protogen/**` as generated output: do not hand-edit generated files.
- `go generate ./...` regenerates code through Buf using `buf.gen.yaml`.
- `buf.gen.yaml` includes inputs from both local `apis` and external `buf.build/tilebox/api`.
- If a required proto/schema/service change belongs to the external `api` repository (source of `buf.build/tilebox/api`), stop and ask the developer to point you to that `api` repo so the change can be made there first.
- After proto changes (local or external), regenerate and then run at least `go test ./...`.

## Code Style
- Use `stretchr/testify` for assertions and `pgregory.net/rapid` for property-based tests.
- Prefer table-driven tests for multi-case behavior.
- Prefer descriptive variable names over heavy abbreviation.
- Keep imports grouped as stdlib, external, internal (lint-enforced).
- Avoid `init` functions; prefer explicit constructors.
- Keep comments focused on intent/constraints, not restating obvious code.

## Typical Development Flow
1. Make code changes in public packages (`datasets/v1`, `workflows/v1`, etc.) and supporting internals.
2. If protobuf contracts changed, update proto source in the owning repo and run `go generate ./...`.
3. Run `go test ./...`.
4. Run `golangci-lint run ./...`.
