package grpc

import (
	"connectrpc.com/connect"
	"context"
	"errors"
	"github.com/hashicorp/go-retryablehttp"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// RetryOnStatusUnavailable provides a retry policy for retrying requests if the server is unavailable.
func RetryOnStatusUnavailable(ctx context.Context, resp *http.Response, err error) (bool, error) {
	// do not retry on context.Canceled or context.DeadlineExceeded
	if ctx.Err() != nil {
		return false, ctx.Err()
	}

	if err != nil {
		var v *url.Error
		if errors.As(err, &v) {
			// Retry if the error was due to a connection refused.
			if strings.Contains(v.Error(), "connect: connection refused") {
				slog.InfoContext(ctx, "Auth client retry", "error", v.Error())
				return true, v
			}
		}
	}

	if resp != nil {
		if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500 {
			slog.InfoContext(ctx, "Auth client retry",
				"status", resp.Status,
				"status_code", resp.StatusCode,
				"protocol", resp.Proto,
			)
			return true, nil
		}
	}
	return false, err
}

func RetryHttpClient() connect.HTTPClient {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryWaitMin = 20 * time.Millisecond
	retryClient.RetryWaitMax = 5 * time.Second
	retryClient.RetryMax = 5
	retryClient.Backoff = retryablehttp.DefaultBackoff // exponential backoff
	retryClient.CheckRetry = RetryOnStatusUnavailable

	return retryClient.StandardClient()
}
