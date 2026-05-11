package main

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/tilebox/tilebox-go/workflows/v1"
)

func main() {
	ctx := context.Background()
	client := workflows.NewClient()
	jobID := uuid.MustParse("019e070c-63ba-1c7e-5f1d-65be3e22d52a")

	fmt.Println("Log messages:")
	fmt.Printf("%-8s  %-30s  %s\n", "LEVEL", "TIME", "BODY")
	for logRecord, err := range client.Jobs.QueryLogs(ctx, jobID, workflows.WithSortDirection(workflows.Ascending)) {
		if err != nil {
			slog.Error("failed to query job logs", slog.Any("error", err))
			return
		}

		fmt.Printf("%-8s  %-30s  %s\n",
			logRecord.SeverityText,
			logRecord.Time.Format(time.RFC3339Nano),
			logRecord.Body,
		)
	}

	fmt.Println()
	fmt.Println("Trace spans:")
	fmt.Printf("%-30s  %-40s  %s\n", "TIME", "SPAN NAME", "DURATION")
	for span, err := range client.Jobs.QuerySpans(ctx, jobID) {
		if err != nil {
			slog.Error("failed to query job spans", slog.Any("error", err))
			return
		}

		fmt.Printf("%-30s  %-40s  %s\n",
			span.StartTime.Format(time.RFC3339Nano),
			span.Name,
			span.Duration(),
		)
	}
}
