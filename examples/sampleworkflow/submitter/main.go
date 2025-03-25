package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/integrii/flaggy"
	slogmulti "github.com/samber/slog-multi"
	"github.com/tilebox/tilebox-go/examples/sampleworkflow"
	"github.com/tilebox/tilebox-go/observability"
	"github.com/tilebox/tilebox-go/workflows/v1"
	"go.opentelemetry.io/otel"
)

var (
	serviceName = "task-runner"
	version     = "dev"
)

func main() {
	ctx := context.Background()

	// Parse CLI args
	config, err := parseCliArgs()
	if err != nil {
		flaggy.ShowHelpAndExit(err.Error())
		return
	}

	// Setup slog
	// a fallback logger that logs to stdout, in case axiom logger setup fails
	stdoutHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	log := slog.New(stdoutHandler)
	if config.axiomLogsDataset != "" {
		axiomLogHandler, shutdownLogger, err := observability.NewAxiomLogger(config.axiomLogsDataset, config.axiomAPIKey, slog.LevelDebug, true)
		defer shutdownLogger()
		if err != nil {
			log.Error("failed to set up axiom log handler", "error", err.Error())
		} else {
			log = slog.New(slogmulti.Fanout(axiomLogHandler, stdoutHandler))
		}
	}
	slog.SetDefault(log) // set the global slog logger also to our axiom logger

	// Setup OpenTelemetry tracer provider
	tracerProvider := otel.GetTracerProvider()
	if config.axiomTracesDataset != "" {
		axiomTracerProvider, shutdownTracer, err := observability.NewAxiomTracerProvider(ctx, config.axiomTracesDataset, config.axiomAPIKey, serviceName, version)
		defer shutdownTracer()
		if err != nil {
			log.Error("failed to set up axiom trace exporter", "error", err)
		} else {
			tracerProvider = axiomTracerProvider
		}
	}

	jobs := workflows.NewJobService(
		workflows.NewJobClient(
			workflows.WithURL(config.url),
			workflows.WithAPIKey(config.authToken),
		),
		workflows.WithJobServiceTracerProvider(tracerProvider),
	)

	job, err := jobs.Submit(ctx, "spawn-workflow-tree", "testing-4qgCk4qHH85qR7", 0,
		&sampleworkflow.SampleTask{
			Message:      "hello go runner!",
			Depth:        8,
			BranchFactor: 4,
		},
	)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to submit job", "error", err)
		return
	}

	slog.InfoContext(ctx, "Job submitted", "job_id", uuid.Must(uuid.FromBytes(job.GetId().GetUuid())))
}

type Config struct {
	url                string
	authToken          string
	axiomAPIKey        string
	axiomTracesDataset string
	axiomLogsDataset   string
}

func parseCliArgs() (Config, error) {
	url := "https://api.tilebox.com"
	authToken := os.Getenv("TILEBOX_API_KEY")
	axiomAPIKey := os.Getenv("AXIOM_API_KEY")
	axiomTracesDataset := os.Getenv("AXIOM_TRACES_DATASET")
	axiomLogsDataset := os.Getenv("AXIOM_LOGS_DATASET")

	flaggy.SetName("task-runner")
	flaggy.AddPositionalValue(&url, "workflows-service", 1, false, "The workflows service to connect to. Either a HTTP(s) URL or a unix socket. Defaults to https://api.tilebox.com")
	flaggy.String(&authToken, "t", "token", "The auth token to use for the task server. Defaults to the value of the TILEBOX_API_KEY environment variable.")
	flaggy.String(&axiomAPIKey, "", "axiom-api-key", "The axiom API key to use. Defaults to the value of the AXIOM_API_KEY environment variable.")
	flaggy.String(&axiomTracesDataset, "", "axiom-traces-dataset", "The axiom dataset to use for ingesting traces. Defaults to the value of the AXIOM_TRACES_DATASET environment variable.")
	flaggy.String(&axiomLogsDataset, "", "axiom-logs-dataset", "The axiom dataset to use for ingesting logs. Defaults to the value of the AXIOM_LOGS_DATASET environment variable.")
	flaggy.Parse()

	if strings.HasPrefix(url, "http") && authToken == "" {
		return Config{}, errors.New("no auth token provided. Please provide an auth token via the -t flag or the TILEBOX_API_KEY environment variable")
	}

	return Config{
		url:                url,
		authToken:          authToken,
		axiomAPIKey:        axiomAPIKey,
		axiomTracesDataset: axiomTracesDataset,
		axiomLogsDataset:   axiomLogsDataset,
	}, nil
}
