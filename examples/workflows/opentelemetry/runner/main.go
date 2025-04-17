package main

import (
	"context"
	"go.opentelemetry.io/otel"
	"log/slog"
	"os"

	"github.com/tilebox/tilebox-go/examples/workflows/axiom"
	"github.com/tilebox/tilebox-go/observability"
	"github.com/tilebox/tilebox-go/workflows/v1"
)

// specify a service name and version to identify the instrumenting application in traces and logs
var service = &observability.Service{Name: "task-runner", Version: "dev"}

func main() {
	ctx := context.Background()

	tileboxAPIKey := os.Getenv("TILEBOX_API_KEY")
	endpoint := "http://localhost:4318"
	headers := map[string]string{
		"Authorization": "Bearer <ENDPOINT_AUTH>",
	}

	// Setup an OpenTelemetry log handler, exporting logs to an OTLP compatible log endpoint
	otelHandler, shutdownLogger, err := logger.NewOtelHandler(ctx, service,
		logger.WithEndpointURL(endpoint),
		logger.WithHeaders(headers),
		logger.WithLevel(slog.LevelInfo), // export logs at info level and above as OTEL logs
	)
	defer shutdownLogger(ctx)
	if err != nil {
		slog.Error("failed to set up otel log handler", slog.Any("error", err))
		return
	}
	tileboxLogger := logger.New( // initialize a slog.Logger and set it as the default
		otelHandler, // export logs to the OTEL handler
		logger.NewConsoleHandler( // and additionally export WARN and ERROR logs to stdout
			logger.WithLevel(slog.LevelWarn),
		),
	)
	slog.SetDefault(tileboxLogger) // all slog.Log calls will be forwarded to the tilebox logger

	// Setup an OpenTelemetry trace span processor, exporting traces and spans to an OTLP compatible trace endpoint
	tileboxTracerProvider, shutdown, err := tracer.NewOtelProvider(ctx, service,
		tracer.WithEndpointURL(endpoint),
		tracer.WithHeaders(headers),
	)

	/* Or alternatively / under the hood this does:
	processor1 := tracer.NewOtelSpanProcesser(ctx,
		tracer.WithEndpointURL(endpoint),
		tracer.WithHeaders(headers),
	)
	processor2 := tracer.NewAxiomSpanProcesser(ctx,
		tracer.WithDataset("tilebox-traces-dev"),
	)

	tracer.NewOtelProviderWithProcessors(ctx, service, processor1, processor2)
	*/

	defer shutdown()
	if err != nil {
		slog.Error("failed to set up otel span processor", slog.Any("error", err))
		return
	}
	otel.SetTracerProvider(tileboxTracerProvider) // set the tilebox tracer provider as the global OTEL tracer provider

	client := workflows.NewClient(workflows.WithAPIKey(tileboxAPIKey))

	cluster, err := client.Clusters.Get(ctx, "testing-4qgCk4qHH85qR7")
	if err != nil {
		slog.Error("failed to get cluster", slog.Any("error", err))
		return
	}

	taskRunner, err := client.NewTaskRunner(cluster)
	if err != nil {
		slog.Error("failed to create task runner", slog.Any("error", err))
		return
	}

	err = taskRunner.RegisterTasks(&axiom.MyAxiomTask{})
	if err != nil {
		slog.Error("failed to register tasks", slog.Any("error", err))
		return
	}

	taskRunner.RunForever(ctx)
}
