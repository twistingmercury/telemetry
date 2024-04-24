package attributes

import (
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"time"
)

// Attributes is an interface for common attributes that are shared between telemetry and middleware packages.
type Attributes interface {
	All() []attribute.KeyValue
	ServiceVersion() string
	ServiceName() string
	Namespace() string
	Environment() string
	BatchingDuration() time.Duration
	Get(key string) string
}

type attributes struct {
	namespace      attribute.KeyValue
	serviceName    attribute.KeyValue
	serviceVersion attribute.KeyValue
	environment    attribute.KeyValue
	batchDuration  time.Duration
	keyValues      []attribute.KeyValue
}

// All returns all the attributes as a slice of attribute.KeyValue.
func (a attributes) All() []attribute.KeyValue {
	return append([]attribute.KeyValue{
		a.namespace,
		a.serviceName,
		a.serviceVersion,
		a.environment,
	}, a.keyValues...)
}

// ServiceVersion returns the service version.
func (a attributes) ServiceVersion() string {
	return a.serviceVersion.Value.AsString()
}

// ServiceName returns the service name.
func (a attributes) ServiceName() string {
	return a.serviceName.Value.AsString()
}

// Namespace returns the namespace.
func (a attributes) Namespace() string {
	return a.namespace.Value.AsString()
}

// Environment returns the environment.
func (a attributes) Environment() string {
	return a.environment.Value.AsString()
}

// Get returns the value of the attribute key.
func (a attributes) Get(key string) string {
	for _, kv := range a.keyValues {
		if kv.Key == attribute.Key(key) {
			return kv.Value.AsString()
		}
	}
	return ""
}

func (a attributes) BatchingDuration() time.Duration {
	return a.batchDuration
}

func New(namespace, serviceName, version, environment string, attribs ...attribute.KeyValue) Attributes {
	return NewWithBatchDuration(namespace, serviceName, version, environment, 0, attribs...)
}

func NewWithBatchDuration(namespace, serviceName, version, environment string, batchDruation time.Duration, attribs ...attribute.KeyValue) Attributes {
	return &attributes{
		namespace:      semconv.ServiceNamespaceKey.String(namespace),
		serviceName:    semconv.ServiceNameKey.String(serviceName),
		serviceVersion: semconv.ServiceVersionKey.String(version),
		environment:    semconv.DeploymentEnvironmentKey.String(environment),
		batchDuration:  batchDruation,
		keyValues:      attribs,
	}
}
