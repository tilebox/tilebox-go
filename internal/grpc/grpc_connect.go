package grpc

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"math/rand/v2"
	"net/http"
	"time"

	"connectrpc.com/connect"
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

func RetryHTTPClient() connect.HTTPClient {
	retryClient := retryablehttp.NewClient()
	retryClient.Logger = nil
	retryClient.RetryWaitMin = 20 * time.Millisecond
	retryClient.RetryWaitMax = 10 * time.Second
	retryClient.RetryMax = 4 // Five total attempts: the initial request plus four retries.
	retryClient.Backoff = exponentialJitterBackoff
	retryClient.CheckRetry = shouldRetryTransientRequest

	return retryClient.StandardClient()
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
