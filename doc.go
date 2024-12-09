// Package tilebox is the Go client for Tilebox.
//
// Usage:
//
//	import "github.com/tilebox/tilebox-go/datasets/v1" // When using Datasets
//	import "github.com/tilebox/tilebox-go/workflows/v1" // When using Workflows
//	import "github.com/tilebox/tilebox-go/grpc" // When using gRPC helpers
//	import "github.com/tilebox/tilebox-go/observability" // When using observability helpers
//
// To construct a client:
//
//	client := datasets.NewClient(
//		datasets.WithAPIKey(os.Getenv("TILEBOX_API_KEY")),
//	)
//
// List all datasets:
//
//	datasets, err := client.Datasets(ctx)
//
// For examples on how to use the library, see the [examples] directory.
//
// [examples]: https://github.com/tilebox/tilebox-go/tree/main/examples
package tilebox
