package opentelemetry

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/tilebox/tilebox-go/workflows/v1"
)

type MyOpenTelemetryTask struct{}

func (t *MyOpenTelemetryTask) Execute(ctx context.Context) error {
	// Start an openTelemetry tracing span that will be exported to the OpenTelemetry endpoint
	result, err := workflows.WithTaskSpanResult(ctx, "Expensive Compute", func(ctx context.Context) (int, error) {
		return 6 * 7, nil
	})
	if err != nil {
		return fmt.Errorf("failed to compute: %w", err)
	}

	// Log the result to the OpenTelemetry endpoint
	slog.InfoContext(ctx, "Computed result", slog.Int("result", result))
	return nil
}
