package observability

import (
	"os"
	"runtime"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.30.0"
)

const (
	defaultServiceName      = "tilebox-go"
	defaultServiceNamespace = "tilebox.workflows"
	defaultServiceVersion   = "dev"
)

// processInstanceID is a process-unique identifier that distinguishes different
// instances of the same service running on the same host.
var processInstanceID = uuid.NewString()

// NewResource creates the default OpenTelemetry resource used for Tilebox workflow observability.
func NewResource(otelService *Service) *resource.Resource {
	serviceName := defaultServiceName
	serviceVersion := defaultServiceVersion
	if otelService != nil {
		if otelService.Name != "" {
			serviceName = otelService.Name
		}
		if otelService.Version != "" {
			serviceVersion = otelService.Version
		}
	}

	attrs := []attribute.KeyValue{
		semconv.ServiceNamespaceKey.String(defaultServiceNamespace),
		semconv.ServiceNameKey.String(serviceName),
		semconv.ServiceInstanceIDKey.String(processInstanceID),
		semconv.ServiceVersionKey.String(serviceVersion),
		semconv.ProcessPIDKey.Int(os.Getpid()),
		semconv.HostArchKey.String(runtime.GOARCH),
		semconv.OSTypeKey.String(runtime.GOOS),
	}

	if hostname, err := os.Hostname(); err == nil && hostname != "" {
		attrs = append(attrs, semconv.HostNameKey.String(hostname))
	}

	return resource.NewWithAttributes(semconv.SchemaURL, attrs...)
}
