package client

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMetadataHeaderValue(t *testing.T) {
	metadata := Metadata{
		Name:                 "go",
		Version:              "v1.2.3",
		Runtime:              "go",
		RuntimeVersion:       "1.27.0",
		OS:                   "linux",
		OSVersion:            "6.8.0",
		Arch:                 "amd64",
		ExecutionEnvironment: "github-actions",
		Invoker:              "claude-code",
		InvokerVersion:       "1.0.0",
		CloudProvider:        "gcp",
		CloudPlatform:        "gcp_cloud_run",
		CloudRegion:          "europe-west1",
	}

	require.Equal(t, `name="go", version="v1.2.3", runtime="go", runtime-version="1.27.0", os="linux", os-version="6.8.0", arch="amd64", execution-environment="github-actions", invoker="claude-code", invoker-version="1.0.0", cloud-provider="gcp", cloud-platform="gcp_cloud_run", cloud-region="europe-west1"`, metadata.HeaderValue())
}

func TestMetadataHeaderValueOmitsInvalidValues(t *testing.T) {
	metadata := Metadata{
		Name:          "cli",
		Version:       `1.0"beta\\1`,
		CloudProvider: "invalid\nvalue",
		CloudRegion:   strings.Repeat("a", clientFieldMaxSize+1),
	}

	require.Equal(t, `name="cli", version="1.0\"beta\\\\1"`, metadata.HeaderValue())
}

func TestDefaultMetadata(t *testing.T) {
	metadata := DefaultMetadata()

	require.Equal(t, "go", metadata.Name)
	require.NotEmpty(t, metadata.Version)
	require.Equal(t, "go", metadata.Runtime)
	require.NotEmpty(t, metadata.RuntimeVersion)
	require.NotEmpty(t, metadata.OS)
	require.NotEmpty(t, metadata.Arch)
	require.NotEmpty(t, metadata.HeaderValue())
}

func TestNewMetadataForWrapper(t *testing.T) {
	metadata := NewMetadata("cli", "v1.2.3")

	require.Equal(t, "cli", metadata.Name)
	require.Equal(t, "v1.2.3", metadata.Version)
	require.Empty(t, metadata.Runtime)
	require.Empty(t, metadata.RuntimeVersion)
	require.NotEmpty(t, metadata.OS)
	require.NotEmpty(t, metadata.Arch)
}

func TestEnvironmentDetection(t *testing.T) {
	t.Setenv("CLOUD_PROVIDER", "")
	t.Setenv("CLOUD_PLATFORM", "")
	t.Setenv("CLOUD_REGION", "")
	t.Setenv("GITHUB_ACTIONS", "true")
	t.Setenv("AWS_EXECUTION_ENV", "AWS_Lambda_go1.x")
	t.Setenv("AWS_REGION", "eu-west-1")

	metadata := DefaultMetadata()

	require.Equal(t, "github-actions", metadata.ExecutionEnvironment)
	require.Equal(t, "aws", metadata.CloudProvider)
	require.Equal(t, "aws_lambda", metadata.CloudPlatform)
	require.Equal(t, "eu-west-1", metadata.CloudRegion)
}

func TestInvokerDetection(t *testing.T) {
	tests := []struct {
		name            string
		environment     map[string]string
		expectedInvoker string
		expectedVersion string
	}{
		{name: "Amp", environment: map[string]string{"AGENT": "amp"}, expectedInvoker: "amp"},
		{
			name:            "GitHub Copilot",
			environment:     map[string]string{"COPILOT_AGENT_SESSION_ID": "session-id"},
			expectedInvoker: "github-copilot",
		},
		{name: "OpenCode", environment: map[string]string{"OPENCODE": "1"}, expectedInvoker: "opencode"},
		{name: "Claude Code", environment: map[string]string{"CLAUDECODE": "1"}, expectedInvoker: "claude-code"},
		{name: "Cursor", environment: map[string]string{"CURSOR_AGENT": "1"}, expectedInvoker: "cursor"},
		{
			name:            "Codex",
			environment:     map[string]string{"CODEX_SESSION_ID": "session-id", "CODEX_VERSION": "1.2.3"},
			expectedInvoker: "codex",
			expectedVersion: "1.2.3",
		},
		{name: "Gemini CLI", environment: map[string]string{"GEMINI_CLI": "1"}, expectedInvoker: "gemini-cli"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for _, name := range []string{"AGENT", "COPILOT_AGENT_SESSION_ID", "OPENCODE", "CLAUDECODE", "CURSOR_AGENT", "CODEX_SESSION_ID", "CODEX_THREAD_ID", "CODEX_VERSION", "GEMINI_CLI"} {
				t.Setenv(name, "")
			}
			for name, value := range test.environment {
				t.Setenv(name, value)
			}

			invokerName, invokerVersion := invoker()
			require.Equal(t, test.expectedInvoker, invokerName)
			require.Equal(t, test.expectedVersion, invokerVersion)
		})
	}
}
