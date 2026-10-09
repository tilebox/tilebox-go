package workflows

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNotifyObjectCreated(t *testing.T) {
	id := uuid.New()
	for _, withTime := range []bool{false, true} {
		t.Run(map[bool]string{false: "without timestamp", true: "with timestamp"}[withTime], func(t *testing.T) {
			var observed time.Time
			expected := `{"object_path":"incoming/\u96ea %20.tif"}`
			if withTime {
				observed = time.Date(2026, time.October, 8, 14, 2, 3, 456000000, time.UTC)
				expected = `{"object_path":"incoming/\u96ea %20.tif","event_time":"2026-10-08T14:02:03.456Z"}`
			}
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPost, r.Method)
				assert.Equal(t, "/prefix/v1/storage-notifications/tilebox-cli/"+id.String(), r.URL.Path)
				assert.Equal(t, "Bearer test-key", r.Header.Get("Authorization"))
				assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
				assert.NotEmpty(t, r.Header.Get("Tilebox-Client"))
				body, err := io.ReadAll(r.Body)
				assert.NoError(t, err)
				assert.JSONEq(t, expected, string(body))
				w.WriteHeader(http.StatusNoContent)
			}))
			t.Cleanup(server.Close)
			client := NewClient(WithURL(server.URL+"/prefix/"), WithAPIKey("test-key"), WithDisableTracing(), WithDisableLogging())
			require.NoError(t, client.StorageLocations.NotifyObjectCreated(t.Context(), id, "incoming/\u96ea %20.tif", observed))
		})
	}
}

func TestNotifyObjectCreatedDoesNotRetryOrRedirect(t *testing.T) {
	for _, status := range []int{http.StatusBadRequest, http.StatusUnauthorized, http.StatusNotFound, http.StatusTooManyRequests, http.StatusServiceUnavailable, http.StatusTemporaryRedirect} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			var calls atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
				assert.Empty(t, r.Header.Get("Authorization"))
				w.Header().Set("Location", "/redirected")
				http.Error(w, "notification rejected", status)
			}))
			t.Cleanup(server.Close)
			client := NewClient(WithURL(server.URL), WithAPIKey(""), WithDisableTracing(), WithDisableLogging())
			err := client.StorageLocations.NotifyObjectCreated(t.Context(), uuid.New(), "file.tif", time.Time{})
			require.ErrorContains(t, err, "notification rejected")
			assert.EqualValues(t, 1, calls.Load())
		})
	}
}

func TestNotifyObjectCreatedCancellation(t *testing.T) {
	started := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		_, _ = io.Copy(io.Discard, r.Body)
		close(started)
		<-r.Context().Done()
	}))
	t.Cleanup(server.Close)
	client := NewClient(WithURL(server.URL), WithAPIKey(""), WithHTTPClient(server.Client()), WithDisableTracing(), WithDisableLogging())
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	done := make(chan error, 1)
	go func() { done <- client.StorageLocations.NotifyObjectCreated(ctx, uuid.New(), "file.tif", time.Time{}) }()
	select {
	case <-started:
	case <-time.After(5 * time.Second):
		t.Fatal("request did not start")
	}
	cancel()
	select {
	case err := <-done:
		require.ErrorIs(t, err, context.Canceled)
		require.ErrorContains(t, err, "delivery may have succeeded")
	case <-time.After(5 * time.Second):
		t.Fatal("request did not cancel")
	}
}
