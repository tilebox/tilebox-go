package axiom

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/tilebox/tilebox-go/workflows/v1"
)

type MyAxiomTask struct{}

func (t *MyAxiomTask) Execute(ctx context.Context) error {
	// Start an openTelemetry tracing span that will be exported to Axiom
	result, err := workflows.WithTaskSpanResult(ctx, "Expensive Compute", func(ctx context.Context) (int, error) {
		return 6 * 7, nil
	})
	if err != nil {
		return fmt.Errorf("failed to compute: %w", err)
	}

	// Log the result to Axiom
	slog.InfoContext(ctx, "Computed result", slog.Int("result", result))
	return nil
}
