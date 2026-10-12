package grpc

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"math/rand/v2"
	"net"
	"net/http"
	"strings"
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

// NewHTTPClient returns a pooled client and the Connect base URL. Socket paths
// retain HTTP/1 and no retries for compatibility with older deployments.
func NewHTTPClient(baseURL string) (connecthttp.HTTPClient, string) {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.ForceAttemptHTTP2 = true
	transport.MaxIdleConns = 100
	transport.MaxIdleConnsPerHost = 10
	transport.IdleConnTimeout = 90 * time.Second
	client := &http.Client{Transport: transport, CheckRedirect: noRedirects}
	if !strings.HasPrefix(baseURL, "https://") && !strings.HasPrefix(baseURL, "http://") {
		transport.ForceAttemptHTTP2 = false
		transport.DialContext = func(ctx context.Context, _, _ string) (net.Conn, error) {
			return (&net.Dialer{Timeout: 30 * time.Second, KeepAlive: 30 * time.Second}).DialContext(ctx, "unix", baseURL)
		}
		return client, "http://localhost"
	}
	return newRetryHTTPClient(client), baseURL
}

// NonRetryHTTPClient shares the connection pool while giving notifications an
// independent timeout and preventing retries and redirects. The client must
// have been returned by NewHTTPClient, not a custom option.
func NonRetryHTTPClient(client connecthttp.HTTPClient, timeout time.Duration) *http.Client {
	if retryClient, ok := client.(*retryHTTPClient); ok {
		client = retryClient.client.HTTPClient
	}
	base := client.(*http.Client)
	return &http.Client{Transport: base.Transport, Timeout: timeout, CheckRedirect: noRedirects}
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

func newRetryHTTPClient(httpClient *http.Client) connecthttp.HTTPClient {
	retryClient := retryablehttp.NewClient()
	retryClient.HTTPClient = httpClient
	retryClient.Logger = nil
	retryClient.RetryWaitMin = 20 * time.Millisecond
	retryClient.RetryWaitMax = 10 * time.Second
	retryClient.RetryMax = 4 // Five total attempts: the initial request plus four retries.
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
