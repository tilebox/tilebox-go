package workflows

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/tilebox/tilebox-go/internal/grpc"
)

func (c storageLocationClient) NotifyObjectCreated(ctx context.Context, subscriptionID uuid.UUID, objectPath string, eventTime time.Time) error {
	payload := struct {
		ObjectPath string     `json:"object_path"`
		EventTime  *time.Time `json:"event_time,omitempty"`
	}{ObjectPath: objectPath}
	if !eventTime.IsZero() {
		payload.EventTime = &eventTime
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("encode storage notification: %w", err)
	}
	endpoint := strings.TrimRight(c.baseURL, "/") + "/v1/storage-notifications/tilebox-cli/" + subscriptionID.String()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create storage notification request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}
	if c.clientMetadata != "" {
		req.Header.Set(grpc.ClientHeader, c.clientMetadata)
	}
	response, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("send storage notification (delivery may have succeeded): %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusNoContent {
		body, err := io.ReadAll(io.LimitReader(response.Body, 4096))
		if err != nil {
			return fmt.Errorf("storage notification returned HTTP %d: read response: %w", response.StatusCode, err)
		}
		return fmt.Errorf("storage notification returned HTTP %d: %s", response.StatusCode, strings.TrimSpace(string(body)))
	}
	return nil
}
