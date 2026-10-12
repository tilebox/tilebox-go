package grpc

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	"connectrpc.com/connect/v2"
	"connectrpc.com/connect/v2/connecthttp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/wrapperspb"
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

	client, _ := NewHTTPClient(server.URL)
	response, err := client.Do(request)
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

func TestRetryHTTPClientReplaysPOST(t *testing.T) {
	for _, test := range []struct {
		name     string
		status   int
		body     string
		attempts int32
	}{
		{"transient", 503, "unavailable", 3},
		{"rate limited", 429, `{"code":"unavailable"}`, 3},
		{"resource exhausted", 429, `{"code":"resource_exhausted"}`, 1},
	} {
		t.Run(test.name, func(t *testing.T) {
			var attempts atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				body, err := io.ReadAll(r.Body)
				assert.NoError(t, err)
				assert.Equal(t, "the complete POST payload", string(body))
				assert.Equal(t, "preserved", r.Header.Get("X-Test"))
				if attempts.Add(1) < 3 {
					w.WriteHeader(test.status)
					_, _ = io.WriteString(w, test.body)
				}
			}))
			t.Cleanup(server.Close)
			request, err := http.NewRequestWithContext(t.Context(), http.MethodPost, server.URL, strings.NewReader("the complete POST payload"))
			require.NoError(t, err)
			request.Header.Set("X-Test", "preserved")
			client, _ := NewHTTPClient(server.URL)
			response, err := client.Do(request)
			require.NoError(t, err)
			defer response.Body.Close()
			body, err := io.ReadAll(response.Body)
			require.NoError(t, err)
			assert.Equal(t, test.attempts, attempts.Load())
			if test.attempts == 1 {
				assert.Equal(t, test.body, string(body))
			}
		})
	}
}

func TestRetryHTTPClientReplaysConnectPayload(t *testing.T) {
	var attempts atomic.Int32
	spec := connect.Spec{Procedure: "/test.Retry/Echo"}
	rpcServer := connect.NewServer()
	rpcServer.Register(connect.Method{Spec: spec, Handler: func(_ context.Context, _ connect.Spec, stream connect.ServerStream) error {
		var request wrapperspb.StringValue
		if err := stream.Receive(&request); err != nil {
			return err
		}
		if attempts.Add(1) == 1 {
			return connect.NewError(connect.CodeUnavailable, "try again")
		}
		return stream.Send(&request)
	}})
	mux := http.NewServeMux()
	connecthttp.Mount(mux, rpcServer)
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)
	httpClient, url := NewHTTPClient(server.URL)
	client := connect.NewClient(connecthttp.NewTransport(httpClient, url, connecthttp.WithSendGzip()))
	request := wrapperspb.String(strings.Repeat("tilebox", 1024))
	var response wrapperspb.StringValue
	require.NoError(t, client.CallUnary(t.Context(), spec, request, &response))
	assert.Equal(t, request.GetValue(), response.GetValue())
	assert.Equal(t, int32(2), attempts.Load())
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

type observedBody struct {
	reads  int
	closed bool
}

func (b *observedBody) Read([]byte) (int, error) { b.reads++; return 0, io.EOF }
func (b *observedBody) Close() error             { b.closed = true; return nil }

func TestRetryHTTPClientStreamingPassesThrough(t *testing.T) {
	body := &observedBody{}
	request, err := http.NewRequestWithContext(t.Context(), http.MethodPost, "http://localhost", body)
	require.NoError(t, err)
	attempts := 0
	client := newRetryHTTPClient(&http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		attempts++
		assert.Same(t, request, r)
		assert.Zero(t, body.reads)
		_ = r.Body.Close()
		return &http.Response{StatusCode: http.StatusServiceUnavailable, Body: http.NoBody}, nil
	})})
	response, err := client.Do(request)
	require.NoError(t, err)
	defer response.Body.Close()
	assert.Equal(t, 1, attempts)
	assert.True(t, body.closed)
}

func TestRetryHTTPClientCancellationDuringBackoff(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var attempts atomic.Int32
	client := newRetryHTTPClient(&http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		attempts.Add(1)
		return &http.Response{StatusCode: http.StatusServiceUnavailable, Body: http.NoBody}, nil
	})}).(*retryHTTPClient)
	client.client.Backoff = func(time.Duration, time.Duration, int, *http.Response) time.Duration {
		cancel()
		return time.Hour
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	require.NoError(t, err)
	response, err := client.Do(request)
	if response != nil {
		require.NoError(t, response.Body.Close())
	}
	require.ErrorIs(t, err, context.Canceled)
	assert.Equal(t, int32(1), attempts.Load())
}

func TestHTTPClientTLSProtocols(t *testing.T) {
	for _, http2 := range []bool{true, false} {
		t.Run(map[bool]string{true: "HTTP2", false: "HTTP1 fallback"}[http2], func(t *testing.T) {
			server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if http2 {
					assert.Equal(t, 2, r.ProtoMajor)
				} else {
					assert.Equal(t, 1, r.ProtoMajor)
				}
				w.WriteHeader(http.StatusOK)
			}))
			server.EnableHTTP2 = http2
			server.StartTLS()
			t.Cleanup(server.Close)
			client, baseURL := NewHTTPClient(server.URL)
			assert.Equal(t, server.URL, baseURL)
			notifications := NonRetryHTTPClient(client, 30*time.Second)
			transport := notifications.Transport.(*http.Transport)
			transport.TLSClientConfig = server.Client().Transport.(*http.Transport).TLSClientConfig.Clone()
			t.Cleanup(transport.CloseIdleConnections)
			assert.Same(t, transport, client.(*retryHTTPClient).client.HTTPClient.Transport)
			request, err := http.NewRequestWithContext(t.Context(), http.MethodGet, baseURL, nil)
			require.NoError(t, err)
			response, err := client.Do(request)
			require.NoError(t, err)
			require.NoError(t, response.Body.Close())
		})
	}
}

func TestRetryPolicyClosesInspectedBody(t *testing.T) {
	body := &observedBody{}
	response := &http.Response{StatusCode: http.StatusTooManyRequests, Body: body}
	retry, err := shouldRetryTransientRequest(context.Background(), response, nil)
	require.NoError(t, err)
	assert.True(t, retry)
	assert.True(t, body.closed)
	require.NoError(t, response.Body.Close())
}

func TestRecordRoundTripperClosesOriginalBody(t *testing.T) {
	originalTransport := http.DefaultTransport
	t.Cleanup(func() { http.DefaultTransport = originalTransport })
	body := &observedBody{}
	http.DefaultTransport = roundTripFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusOK, Body: body}, nil
	})
	var recording bytes.Buffer
	request, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "http://localhost", nil)
	require.NoError(t, err)
	response, err := NewRecordRoundTripper(&recording).RoundTrip(request)
	require.NoError(t, err)
	assert.True(t, body.closed)
	require.NoError(t, response.Body.Close())
	replayed, err := NewReplayRoundTripper(&recording).RoundTrip(request)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, replayed.StatusCode)
	require.NoError(t, replayed.Body.Close())
}

func TestHTTPClientsDoNotFollowRedirects(t *testing.T) {
	var attempts atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts.Add(1)
		http.Redirect(w, r, "/target", http.StatusTemporaryRedirect)
	}))
	t.Cleanup(server.Close)
	client, _ := NewHTTPClient(server.URL)
	for _, current := range []connecthttp.HTTPClient{client, NonRetryHTTPClient(client, 30*time.Second)} {
		request, err := http.NewRequestWithContext(t.Context(), http.MethodGet, server.URL, nil)
		require.NoError(t, err)
		response, err := current.Do(request)
		require.NoError(t, err)
		assert.Equal(t, http.StatusTemporaryRedirect, response.StatusCode)
		require.NoError(t, response.Body.Close())
	}
	assert.Equal(t, int32(2), attempts.Load())
}
