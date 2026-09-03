package workflows

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"
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
