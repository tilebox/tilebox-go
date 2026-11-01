package workflows

import (
	"io"
	"log/slog"
	"net/http"
	"testing"
	"testing/synctest"

	"github.com/stretchr/testify/require"
	"github.com/tilebox/tilebox-go/client"
)

func TestNewPollingTaskRunnerUsesResolvedCluster(t *testing.T) {
	service := &mockMinimalTaskService{}
	client := &Client{taskService: service}
	executor := &failedResponseExecutor{}
	logger := slog.New(slog.DiscardHandler)
	cluster := &Cluster{Slug: "resolved-cluster"}

	pollingRunner, err := client.NewPollingTaskRunner(cluster, executor, logger)

	require.NoError(t, err)
	require.Equal(t, "resolved-cluster", pollingRunner.clusterSlug)
	require.Same(t, service, pollingRunner.service)
	require.Same(t, executor, pollingRunner.executor)
	require.Same(t, logger, pollingRunner.logger)
}

func TestNewPollingTaskRunnerRequiresResolvedCluster(t *testing.T) {
	client := &Client{}

	_, err := client.NewPollingTaskRunner(nil, &failedResponseExecutor{}, nil)
	require.EqualError(t, err, "cluster is required")

	_, err = client.NewPollingTaskRunner(&Cluster{}, &failedResponseExecutor{}, nil)
	require.EqualError(t, err, "cluster slug is required")
}

type hostedClaimTransport func(*http.Request) (*http.Response, error)

func (transport hostedClaimTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	return transport(request)
}

func TestHostedClientDoesNotHideLostClaimsBehindHTTPRetries(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		lifecycle := NewHostedLifecycle()
		attempts := 0
		transport := hostedClaimTransport(func(*http.Request) (*http.Response, error) {
			attempts++
			lifecycle.Drain()
			return nil, io.EOF
		})
		workflowClient := NewClient(client.New(client.WithURL("http://localhost"), client.WithHTTPTransport(transport)))
		poller, err := workflowClient.NewPollingTaskRunner(&Cluster{Slug: "default"}, &hostedTestExecutor{}, slog.New(slog.DiscardHandler), WithHostedLifecycle(lifecycle))
		require.NoError(t, err)
		lifecycle.MarkReady()
		require.True(t, lifecycle.StartIntake())
		require.NoError(t, poller.RunForever(t.Context()))
		require.Equal(t, 1, attempts, "NextTask has no deduplication ID; SDK HTTP retries must not hide lost claims")
		require.True(t, lifecycle.Snapshot().GetClaimUncertain())
		require.False(t, lifecycle.Snapshot().GetSafeToStop())
	})
}
