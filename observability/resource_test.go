package observability

import (
	"os"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.30.0"
)

func Test_NewResource(t *testing.T) {
	res := NewResource(&Service{Name: "test-service", Version: "v1.2.3"})

	attrs := make(map[attribute.Key]attribute.Value)
	for _, attr := range res.Attributes() {
		attrs[attr.Key] = attr.Value
	}

	assert.Equal(t, "tilebox.workflows", attrs[semconv.ServiceNamespaceKey].AsString())
	assert.Equal(t, "test-service", attrs[semconv.ServiceNameKey].AsString())
	assert.Equal(t, processInstanceID, attrs[semconv.ServiceInstanceIDKey].AsString())
	assert.Equal(t, "v1.2.3", attrs[semconv.ServiceVersionKey].AsString())
	assert.Equal(t, int64(os.Getpid()), attrs[semconv.ProcessPIDKey].AsInt64())
	assert.Equal(t, runtime.GOARCH, attrs[semconv.HostArchKey].AsString())
	assert.Equal(t, runtime.GOOS, attrs[semconv.OSTypeKey].AsString())

	if hostname, err := os.Hostname(); err == nil && hostname != "" {
		assert.Equal(t, hostname, attrs[semconv.HostNameKey].AsString())
	}
}

func Test_NewResource_Defaults(t *testing.T) {
	res := NewResource(nil)

	attrs := make(map[attribute.Key]attribute.Value)
	for _, attr := range res.Attributes() {
		attrs[attr.Key] = attr.Value
	}

	assert.Equal(t, defaultServiceName, attrs[semconv.ServiceNameKey].AsString())
	assert.Equal(t, defaultServiceVersion, attrs[semconv.ServiceVersionKey].AsString())
}
