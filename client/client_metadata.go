// Package client provides metadata describing clients that call Tilebox APIs.
package client

import (
	"os"
	"runtime"
	"runtime/debug"
	"strings"
)

const (
	clientHeaderMaxSize = 2 * 1024
	clientFieldMaxSize  = 256
	modulePath          = "github.com/tilebox/tilebox-go"
)

// Metadata describes the client and environment making a Tilebox API request.
// The metadata is sent for analytics only and must not contain secrets.
type Metadata struct {
	Name                 string
	Version              string
	Runtime              string
	RuntimeVersion       string
	OS                   string
	OSVersion            string
	Arch                 string
	ExecutionEnvironment string
	Invoker              string
	InvokerVersion       string
	CloudProvider        string
	CloudPlatform        string
	CloudRegion          string
}

// DefaultMetadata detects metadata for the Tilebox Go SDK from local process information.
// Detection does not make network requests or start subprocesses.
func DefaultMetadata() Metadata {
	metadata := NewMetadata("go", clientVersion())
	metadata.Runtime = "go"
	metadata.RuntimeVersion = strings.TrimPrefix(runtime.Version(), "go")
	return metadata
}

// NewMetadata detects environment metadata for a client with the supplied identity.
// It is useful for wrappers such as the Tilebox CLI, which should identify themselves while
// retaining detected OS, execution environment, invoker, and cloud information.
func NewMetadata(name, version string) Metadata {
	metadata := Metadata{
		Name:      name,
		Version:   version,
		OS:        runtime.GOOS,
		OSVersion: osVersion(),
		Arch:      runtime.GOARCH,
	}
	metadata.ExecutionEnvironment = executionEnvironment()
	metadata.Invoker, metadata.InvokerVersion = invoker()
	metadata.CloudProvider, metadata.CloudPlatform, metadata.CloudRegion = cloudEnvironment()
	return metadata
}

// HeaderValue serializes metadata as an RFC 9651 Structured Fields dictionary.
// Empty and invalid values are omitted.
func (m Metadata) HeaderValue() string {
	fields := [...]struct {
		name  string
		value string
	}{
		{"name", m.Name},
		{"version", m.Version},
		{"runtime", m.Runtime},
		{"runtime-version", m.RuntimeVersion},
		{"os", m.OS},
		{"os-version", m.OSVersion},
		{"arch", m.Arch},
		{"execution-environment", m.ExecutionEnvironment},
		{"invoker", m.Invoker},
		{"invoker-version", m.InvokerVersion},
		{"cloud-provider", m.CloudProvider},
		{"cloud-platform", m.CloudPlatform},
		{"cloud-region", m.CloudRegion},
	}

	var header strings.Builder
	header.Grow(256)
	for _, field := range fields {
		value, ok := structuredString(field.value)
		if !ok {
			continue
		}
		additionalSize := len(field.name) + len(value) + 3
		if header.Len() != 0 {
			additionalSize += 2
		}
		if header.Len()+additionalSize > clientHeaderMaxSize {
			continue
		}
		if header.Len() != 0 {
			header.WriteString(", ")
		}
		header.WriteString(field.name)
		header.WriteString("=\"")
		header.WriteString(value)
		header.WriteByte('"')
	}
	return header.String()
}

func structuredString(value string) (string, bool) {
	if value == "" || len(value) > clientFieldMaxSize {
		return "", false
	}
	var escaped strings.Builder
	for i := range len(value) {
		character := value[i]
		if character < 0x20 || character > 0x7e {
			return "", false
		}
		if character == '"' || character == '\\' {
			escaped.WriteByte('\\')
		}
		escaped.WriteByte(character)
	}
	return escaped.String(), true
}

func clientVersion() string {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return "dev"
	}
	if buildInfo.Main.Path == modulePath && buildInfo.Main.Version != "(devel)" {
		return buildInfo.Main.Version
	}
	for _, dependency := range buildInfo.Deps {
		if dependency.Path == modulePath {
			return dependency.Version
		}
	}
	return "dev"
}

func executionEnvironment() string {
	switch {
	case os.Getenv("GITHUB_ACTIONS") == "true":
		return "github-actions"
	case os.Getenv("GITLAB_CI") == "true":
		return "gitlab-ci"
	case envSet("BUILDKITE"):
		return "buildkite"
	case envSet("CIRCLECI"):
		return "circleci"
	case envSet("JENKINS_URL"):
		return "jenkins"
	case envSet("TEAMCITY_VERSION"):
		return "teamcity"
	case envSet("TF_BUILD"):
		return "azure-pipelines"
	case envSet("K_SERVICE", "CLOUD_RUN_JOB"):
		return "google-cloud-run"
	case strings.HasPrefix(os.Getenv("AWS_EXECUTION_ENV"), "AWS_Lambda_"):
		return "aws-lambda"
	case envSet("FUNCTIONS_WORKER_RUNTIME"):
		return "azure-functions"
	case envSet("KUBERNETES_SERVICE_HOST"):
		return "kubernetes"
	case isTerminal():
		return "terminal"
	default:
		return ""
	}
}

func invoker() (string, string) {
	switch {
	case os.Getenv("AGENT") == "amp":
		return "amp", ""
	case envSet("COPILOT_AGENT_SESSION_ID"):
		return "github-copilot", ""
	case os.Getenv("OPENCODE") == "1":
		return "opencode", ""
	case envSet("CLAUDECODE"):
		return "claude-code", ""
	case envSet("CURSOR_AGENT"):
		return "cursor", ""
	case envSet("CODEX_SESSION_ID", "CODEX_THREAD_ID"):
		return "codex", os.Getenv("CODEX_VERSION")
	case envSet("GEMINI_CLI"):
		return "gemini-cli", ""
	default:
		return "", ""
	}
}

func cloudEnvironment() (string, string, string) {
	var provider, platform, region string
	switch {
	case envSet("AWS_EXECUTION_ENV", "AWS_REGION", "AWS_DEFAULT_REGION"):
		provider = "aws"
		region = firstEnvironmentValue("AWS_REGION", "AWS_DEFAULT_REGION")
		executionEnvironment := os.Getenv("AWS_EXECUTION_ENV")
		switch {
		case strings.HasPrefix(executionEnvironment, "AWS_Lambda_"):
			platform = "aws_lambda"
		case strings.HasPrefix(executionEnvironment, "AWS_ECS_"):
			platform = "aws_ecs"
		}
	case envSet("K_SERVICE", "CLOUD_RUN_JOB", "GAE_ENV"):
		provider = "gcp"
		region = firstEnvironmentValue("GOOGLE_CLOUD_REGION", "CLOUD_RUN_REGION", "FUNCTION_REGION")
		if envSet("K_SERVICE", "CLOUD_RUN_JOB") {
			platform = "gcp_cloud_run"
		} else if envSet("GAE_ENV") {
			platform = "gcp_app_engine"
		}
	case envSet("WEBSITE_INSTANCE_ID", "FUNCTIONS_WORKER_RUNTIME", "REGION_NAME"):
		provider = "azure"
		region = os.Getenv("REGION_NAME")
		if envSet("FUNCTIONS_WORKER_RUNTIME") {
			platform = "azure_functions"
		} else if envSet("WEBSITE_INSTANCE_ID") {
			platform = "azure_app_service"
		}
	}
	return provider, platform, region
}

func isTerminal() bool {
	info, err := os.Stdin.Stat()
	return err == nil && info.Mode()&os.ModeCharDevice != 0
}

func envSet(names ...string) bool {
	for _, name := range names {
		if os.Getenv(name) != "" {
			return true
		}
	}
	return false
}

func firstEnvironmentValue(names ...string) string {
	for _, name := range names {
		if value := os.Getenv(name); value != "" {
			return value
		}
	}
	return ""
}
