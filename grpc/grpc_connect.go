package grpc // import "github.com/tilebox/tilebox-go/grpc"

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"connectrpc.com/connect"
	"github.com/hashicorp/go-retryablehttp"
)

// retryOnStatusUnavailable provides a retry policy for retrying requests if the server is unavailable.
func retryOnStatusUnavailable(ctx context.Context, resp *http.Response, err error) (bool, error) {
	// do not retry on context.Canceled or context.DeadlineExceeded
	if ctx.Err() != nil {
		return false, ctx.Err()
	}

	if err != nil {
		var v *url.Error
		if errors.As(err, &v) {
			// Retry if the error was due to a connection refused.
			if strings.Contains(v.Error(), "connect: connection refused") {
				slog.InfoContext(ctx, "Auth client retry", slog.Any("error", v))
				return true, v
			}
		}
	}

	if resp != nil {
		if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500 {
			slog.InfoContext(ctx, "Auth client retry",
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
	retryClient.RetryMax = 5
	retryClient.Backoff = retryablehttp.LinearJitterBackoff
	retryClient.CheckRetry = retryOnStatusUnavailable

	return retryClient.StandardClient()
}
