package grpc

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShouldRetryTransientRequestRetriesTransportErrors(t *testing.T) {
	tests := map[string]error{
		"EOF":            io.EOF,
		"unexpected EOF": io.ErrUnexpectedEOF,
		"timeout": &url.Error{
			Op:  "Post",
			URL: "https://api.tilebox.dev",
			Err: syscall.ETIMEDOUT,
		},
		"connection refused": &url.Error{
			Op:  "Post",
			URL: "https://api.tilebox.dev",
			Err: syscall.ECONNREFUSED,
		},
		"connection reset": &url.Error{
			Op:  "Post",
			URL: "https://api.tilebox.dev",
			Err: syscall.ECONNRESET,
		},
	}

	for name, transportError := range tests {
		t.Run(name, func(t *testing.T) {
			shouldRetry, err := shouldRetryTransientRequest(context.Background(), nil, transportError)

			assert.True(t, shouldRetry)
			require.NoError(t, err)
		})
	}
}

func TestShouldRetryTransientRequestStopsWhenContextIsCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	shouldRetry, err := shouldRetryTransientRequest(ctx, nil, io.EOF)

	assert.False(t, shouldRetry)
	require.ErrorIs(t, err, context.Canceled)
}

func TestRetryHTTPClientMakesFiveTotalAttempts(t *testing.T) {
	var attempts atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		attempts.Add(1)
		response.WriteHeader(http.StatusServiceUnavailable)
	}))
	t.Cleanup(server.Close)

	request, err := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL, nil)
	require.NoError(t, err)

	response, err := RetryHTTPClient().Do(request)
	if response != nil {
		t.Cleanup(func() { require.NoError(t, response.Body.Close()) })
	}

	require.Error(t, err)
	assert.Equal(t, int32(5), attempts.Load())
}

func TestExponentialJitterBackoffIsBounded(t *testing.T) {
	const (
		minimum = 20 * time.Millisecond
		maximum = 10 * time.Second
	)
	tests := []struct {
		attempt    int
		lowerBound time.Duration
		upperBound time.Duration
	}{
		{attempt: 0, lowerBound: 10 * time.Millisecond, upperBound: 20 * time.Millisecond},
		{attempt: 1, lowerBound: 20 * time.Millisecond, upperBound: 40 * time.Millisecond},
		{attempt: 2, lowerBound: 40 * time.Millisecond, upperBound: 80 * time.Millisecond},
		{attempt: 20, lowerBound: 5 * time.Second, upperBound: maximum},
	}

	for _, test := range tests {
		for range 100 {
			delay := exponentialJitterBackoff(minimum, maximum, test.attempt, nil)
			assert.GreaterOrEqual(t, delay, test.lowerBound)
			assert.LessOrEqual(t, delay, test.upperBound)
		}
	}
}
