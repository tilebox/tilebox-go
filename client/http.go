package client

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"math/rand/v2"
	"net"
	"net/http"
	"time"

	"connectrpc.com/connect/v2/connecthttp"
	"github.com/hashicorp/go-retryablehttp"
)

// shouldRetryTransientRequest retries recoverable transport failures, including EOFs, timeouts, and connection resets,
// as well as transient HTTP statuses. retryablehttp's standard policy excludes permanent transport configuration and
// certificate errors.
func shouldRetryTransientRequest(ctx context.Context, resp *http.Response, err error) (bool, error) {
	// do not retry on context.Canceled or context.DeadlineExceeded
	if ctx.Err() != nil {
		return false, ctx.Err()
	}

	if err != nil {
		shouldRetry, policyErr := retryablehttp.ErrorPropagatedRetryPolicy(ctx, resp, err)
		if shouldRetry {
			slog.InfoContext(ctx, "HTTP client retry", slog.Any("error", err))
		}
		return shouldRetry, policyErr
	}

	if resp != nil {
		// special handling of 429 errors from connect that are actually resource exhausted errors
		// https://connectrpc.com/docs/protocol#error-codes
		if resp.StatusCode == http.StatusTooManyRequests {
			body, err := io.ReadAll(resp.Body)
			_ = resp.Body.Close()
			if err != nil {
				return false, err
			}
			resp.Body = io.NopCloser(bytes.NewReader(body)) // reset body for potential future reads

			var connectErr struct {
				Code string `json:"code"`
			}
			if json.Unmarshal(body, &connectErr) == nil && connectErr.Code == "resource_exhausted" {
				return false, nil // don't retry on resource exhausted errors
			}
		}

		if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500 {
			slog.InfoContext(ctx, "HTTP client retry",
				slog.String("status", resp.Status),
				slog.Int("status_code", resp.StatusCode),
				slog.String("protocol", resp.Proto),
			)
			return true, nil
		}
	}
	return false, err
}

// RetryPolicy controls retries of transient network errors and HTTP 429/5xx
// responses, except Connect resource_exhausted errors. It does not retry RPC
// errors carried in gRPC trailers. Waits use exponential equal jitter, starting
// between MinWait/2 and MinWait, doubling each retry, and capped at MaxWait.
type RetryPolicy struct {
	MaxRetries int
	MinWait    time.Duration
	MaxWait    time.Duration
}

// DefaultRetryPolicy permits four retries (five attempts total).
func DefaultRetryPolicy() RetryPolicy {
	return RetryPolicy{MaxRetries: 4, MinWait: 20 * time.Millisecond, MaxWait: 10 * time.Second}
}

// NewHTTPTransport creates an independent connection pool with HTTP/2 over TLS
// and HTTP/1.1 fallback. It never reads or modifies http.DefaultTransport.
// Configure it before use; a transport can be shared by multiple clients.
func NewHTTPTransport() *http.Transport {
	return &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           (&net.Dialer{Timeout: 30 * time.Second, KeepAlive: 30 * time.Second}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   10,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: time.Second,
	}
}

func newHTTPClient(socketPath string, transport http.RoundTripper, policy RetryPolicy) connecthttp.HTTPClient {
	var owned *http.Transport
	if transport == nil {
		owned = NewHTTPTransport()
		transport = owned
	}
	httpClient := &http.Client{Transport: transport, CheckRedirect: noRedirects}
	if socketPath != "" && owned != nil {
		dialer := &net.Dialer{Timeout: 30 * time.Second, KeepAlive: 30 * time.Second}
		owned.Proxy = nil
		owned.ForceAttemptHTTP2 = false
		owned.DialContext = func(ctx context.Context, _, _ string) (net.Conn, error) {
			return dialer.DialContext(ctx, "unix", socketPath)
		}
	}
	return NewRetryHTTPClient(httpClient, policy)
}

func noRedirects(_ *http.Request, _ []*http.Request) error { return http.ErrUseLastResponse }

type retryHTTPClient struct {
	client *retryablehttp.Client
}

func (client *retryHTTPClient) Do(request *http.Request) (*http.Response, error) {
	if request.Body != nil && request.Body != http.NoBody && request.GetBody == nil {
		return client.client.HTTPClient.Do(request)
	}
	if request.Body != nil {
		defer request.Body.Close()
	}
	var body any
	if request.GetBody != nil {
		body = retryablehttp.ReaderFunc(func() (io.Reader, error) { return request.GetBody() })
	}
	retryRequest, err := retryablehttp.NewRequestWithContext(request.Context(), request.Method, request.URL.String(), body)
	if err != nil {
		return nil, err
	}
	// Preserve headers, host, content length and all other caller metadata; the
	// retryable request retains only its factory and rewinds before each attempt.
	retryRequest.Request = request.Clone(request.Context())
	return client.client.Do(retryRequest)
}

// NewRetryHTTPClient applies the SDK's replay-aware retry policy to an HTTP client.
// It preserves the supplied client's transport, per-attempt timeout, and redirect
// policy. Use a request context deadline to bound all attempts and backoff together.
// A zero policy returns the supplied client unchanged. Invalid retry limits panic.
func NewRetryHTTPClient(httpClient *http.Client, policy RetryPolicy) connecthttp.HTTPClient {
	if policy.MaxRetries < 0 || (policy.MaxRetries > 0 && (policy.MinWait <= 0 || policy.MaxWait < policy.MinWait)) {
		panic("client: retries require MaxRetries >= 0 and 0 < MinWait <= MaxWait")
	}
	if policy.MaxRetries == 0 {
		return httpClient
	}
	retryClient := retryablehttp.NewClient()
	retryClient.HTTPClient = httpClient
	retryClient.Logger = nil
	retryClient.RetryWaitMin = policy.MinWait
	retryClient.RetryWaitMax = policy.MaxWait
	retryClient.RetryMax = policy.MaxRetries
	retryClient.Backoff = exponentialJitterBackoff
	retryClient.CheckRetry = shouldRetryTransientRequest

	return &retryHTTPClient{client: retryClient}
}

func exponentialJitterBackoff(minimum, maximum time.Duration, attempt int, _ *http.Response) time.Duration {
	upperBound := minimum
	for range attempt {
		if upperBound >= maximum/2 {
			upperBound = maximum
			break
		}
		upperBound *= 2
	}

	lowerBound := upperBound / 2
	return lowerBound + rand.N(upperBound-lowerBound+1)
}
