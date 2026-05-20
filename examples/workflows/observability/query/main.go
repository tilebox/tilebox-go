package main

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/tilebox/tilebox-go/workflows/v1"
	"github.com/tilebox/tilebox-go/workflows/v1/job"
)

func main() {
	ctx := context.Background()
	workflows.ConfigureConsoleLogging(slog.LevelInfo)
	client := workflows.NewClient()
	jobID := uuid.MustParse("019e070c-63ba-1c7e-5f1d-65be3e22d52a")

	fmt.Println("Log messages:")                             //nolint:forbidigo // This example intentionally prints query results to stdout.
	fmt.Printf("%-8s  %-30s  %s\n", "LEVEL", "TIME", "BODY") //nolint:forbidigo // This example intentionally prints query results to stdout.
	for logRecord, err := range client.Jobs.QueryLogs(ctx, jobID, job.WithSortDirection(job.Ascending)) {
		if err != nil {
			slog.ErrorContext(ctx, "failed to query job logs", slog.Any("error", err))
			return
		}

		fmt.Printf("%-8s  %-30s  %s\n", //nolint:forbidigo // This example intentionally prints query results to stdout.
			logRecord.SeverityText,
			logRecord.Time.Format(time.RFC3339Nano),
			logRecord.Body,
		)
	}

	fmt.Println()                                                     //nolint:forbidigo // This example intentionally prints query results to stdout.
	fmt.Println("Trace spans:")                                       //nolint:forbidigo // This example intentionally prints query results to stdout.
	fmt.Printf("%-30s  %-40s  %s\n", "TIME", "SPAN NAME", "DURATION") //nolint:forbidigo // This example intentionally prints query results to stdout.
	for span, err := range client.Jobs.QuerySpans(ctx, jobID) {
		if err != nil {
			slog.ErrorContext(ctx, "failed to query job spans", slog.Any("error", err))
			return
		}

		fmt.Printf("%-30s  %-40s  %s\n", //nolint:forbidigo // This example intentionally prints query results to stdout.
			span.StartTime.Format(time.RFC3339Nano),
			span.Name,
			span.Duration(),
		)
	}
}
