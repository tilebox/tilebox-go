package main

import (
	"connectrpc.com/connect"
	"context"
	"fmt"
	"github.com/integrii/flaggy"
	slogotel "github.com/remychantenay/slog-otel"
	"github.com/tilebox/tilebox-go/task-runner/sample_workflow"
	"github.com/tilebox/tilebox-go/task-runner/services"
	"github.com/tilebox/tilebox-go/util/grpc"
	"github.com/tilebox/tilebox-go/workflows-service/protogen/go/workflows/v1/workflowsv1connect"
	"go.opentelemetry.io/otel/sdk/trace"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strings"
)

var serviceName = "task-runner"
var version = "dev"

func main() {
	ctx := context.Background()

	// Parse CLI args
	config, err := parseCliArgs()
	if err != nil {
		flaggy.ShowHelpAndExit(err.Error())
	}

	// Setup slog
	var logHandler slog.Handler
	if config.axiomLogsDataset != "" {
		axiomHandler, err := services.AxiomLogHandler(config.axiomLogsDataset, config.axiomApiKey, slog.LevelDebug)
		if err != nil {
			slog.Error("failed to set up axiom log handler", "error", err.Error())
		} else {
			defer axiomHandler.Close()
		}
	}
	if logHandler == nil { // if axiomLogsDataset is empty or if the axiom log handler failed to set up with an error
		logHandler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	}
	// wrap the log handler with the Otel handler and set it as the default logger
	slog.SetDefault(slog.New(slogotel.OtelHandler{Next: logHandler}))

	// Setup OpenTelemetry tracing
	traceExporters := make([]trace.SpanExporter, 0)
	if config.axiomTracesDataset != "" {
		axiomTraceExporter, err := services.AxiomTraceExporter(ctx, config.axiomTracesDataset, config.axiomApiKey)
		if err != nil {
			slog.Error("failed to set up axiom trace exporter", "error", err)
		}
		traceExporters = append(traceExporters, axiomTraceExporter)
	}

	shutDownFunc, err := services.SetupOtelTracing(serviceName, version, traceExporters...)
	if err != nil {
		slog.Error("failed to set up opentelemetry tracing", "error", err)
	} else {
		defer shutDownFunc(ctx)
	}

	taskClient := clientFromConfig(config.serverSocket, config.authToken)
	runner := services.NewTaskRunner(taskClient)

	err = runner.RegisterTasks(
		&sample_workflow.SampleTask{},
		&sample_workflow.SpawnWorkflowTreeTask{},
	)
	if err != nil {
		slog.Error("failed to register task", "error", err.Error())
		return // exit the program if we can't register one of the tasks
	}

	runner.Run(ctx)
}

func isHttp(serverSocket string) bool {
	return strings.HasPrefix(serverSocket, "http://") || strings.HasPrefix(serverSocket, "https://")
}

type Config struct {
	serverSocket       string
	authToken          string
	axiomApiKey        string
	axiomTracesDataset string
	axiomLogsDataset   string
}

func parseCliArgs() (Config, error) {
	serverSocket := "https://api.tilebox.com"
	authToken := os.Getenv("TILEBOX_API_KEY")
	axiomApiKey := os.Getenv("AXIOM_API_KEY")
	axiomTracesDataset := os.Getenv("AXIOM_TRACES_DATASET")
	axiomLogsDataset := os.Getenv("AXIOM_LOGS_DATASET")

	flaggy.SetName("task-runner")
	flaggy.AddPositionalValue(&serverSocket, "workflows-service", 1, false, "The workflows service to connect to. Either a HTTP(s) URL or a unix socket. Defaults to https://api.tilebox.com")
	flaggy.String(&authToken, "t", "token", "The auth token to use for the task server. Defaults to the value of the TILEBOX_API_KEY environment variable.")
	flaggy.String(&axiomApiKey, "", "axiom-api-key", "The axiom API key to use. Defaults to the value of the AXIOM_API_KEY environment variable.")
	flaggy.String(&axiomTracesDataset, "", "axiom-traces-dataset", "The axiom dataset to use for ingesting traces. Defaults to the value of the AXIOM_TRACES_DATASET environment variable.")
	flaggy.String(&axiomLogsDataset, "", "axiom-logs-dataset", "The axiom dataset to use for ingesting logs. Defaults to the value of the AXIOM_LOGS_DATASET environment variable.")
	flaggy.Parse()

	if isHttp(serverSocket) && authToken == "" {
		return Config{}, fmt.Errorf("no auth token provided. Please provide an auth token via the -t flag or the TILEBOX_API_KEY environment variable")
	}

	return Config{
		serverSocket:       serverSocket,
		authToken:          authToken,
		axiomApiKey:        axiomApiKey,
		axiomTracesDataset: axiomTracesDataset,
		axiomLogsDataset:   axiomLogsDataset,
	}, nil
}

// clientFromConfig returns a task service client based on the CLI flags
// The client is either an HTTP client or a unix socket client, depending on the serverSocket parameter.
// In case of an HTTP client, also an auth token is required.
// This function is deliberately kept here directly in main, and a similar one in the submitter/main.go file
// to keep those as examples as standalone as possible.
// That's also the reason we don't use our config package here.
func clientFromConfig(serverSocket, authToken string) workflowsv1connect.TaskServiceClient {
	if isHttp(serverSocket) {
		slog.Info("Starting task runner connected to HTTP server", "workflows-service", serverSocket)
		return workflowsv1connect.NewTaskServiceClient(
			services.RetryHttpClient(), serverSocket, connect.WithInterceptors(
				grpc.NewAddAuthTokenInterceptor(func() string {
					return authToken
				})),
		)
	}

	// connect to a unix socket
	dial := func(ctx context.Context, proto, addr string) (net.Conn, error) {
		return net.Dial("unix", serverSocket)
	}
	transport := &http.Transport{DialContext: dial}
	client := &http.Client{Transport: transport}
	slog.Info("Starting task runner connected to unix socket", "workflows-service-socket", serverSocket)
	return workflowsv1connect.NewTaskServiceClient(client, "http://workflows-service") // URL here is irrelevant
}
