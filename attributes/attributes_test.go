package attributes

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

func TestAttributes_All(t *testing.T) {
	// Create an instance of Attributes
	attribs := New("test_namespace", "test_service", "1.0.0", "test_env",
		attribute.String("key1", "value1"),
		attribute.Int("key2", 123),
	)

	// Get all Attributes
	allAttribs := attribs.All()

	// Assert the expected number of Attributes
	assert.Len(t, allAttribs, 6)

	// Assert the presence of specific Attributes
	assert.Contains(t, allAttribs, semconv.ServiceNamespaceKey.String("test_namespace"))
	assert.Contains(t, allAttribs, semconv.ServiceNameKey.String("test_service"))
	assert.Contains(t, allAttribs, semconv.ServiceVersionKey.String("1.0.0"))
	assert.Contains(t, allAttribs, semconv.DeploymentEnvironmentKey.String("test_env"))
	assert.Contains(t, allAttribs, attribute.String("key1", "value1"))
	assert.Contains(t, allAttribs, attribute.Int("key2", 123))
}

func TestAttributes_ServiceVersion(t *testing.T) {
	// Create an instance of Attributes
	attribs := New("test_namespace", "test_service", "1.0.0", "test_env")

	// Get the service version
	version := attribs.ServiceVersion()

	// Assert the expected service version
	assert.Equal(t, "1.0.0", version)
}

func TestAttributes_ServiceName(t *testing.T) {
	// Create an instance of Attributes
	attribs := New("test_namespace", "test_service", "1.0.0", "test_env")

	// Get the service name
	name := attribs.ServiceName()

	// Assert the expected service name
	assert.Equal(t, "test_service", name)
}

func TestAttributes_Namespace(t *testing.T) {
	// Create an instance of Attributes
	attribs := New("test_namespace", "test_service", "1.0.0", "test_env")

	// Get the namespace
	namespace := attribs.Namespace()

	// Assert the expected namespace
	assert.Equal(t, "test_namespace", namespace)
}

func TestAttributes_Environment(t *testing.T) {
	// Create an instance of Attributes
	attribs := New("test_namespace", "test_service", "1.0.0", "test_env")

	// Get the environment
	env := attribs.Environment()

	// Assert the expected environment
	assert.Equal(t, "test_env", env)
}

func TestAttributes_Get(t *testing.T) {
	// Create an instance of Attributes
	attribs := New("test_namespace", "test_service", "1.0.0", "test_env",
		attribute.String("key1", "value1"),
		attribute.Int("key2", 123),
	)

	// Get the value of a specific attribute
	value := attribs.Get("key1")

	// Assert the expected value
	assert.Equal(t, "value1", value)

	// Get the value of a non-existent attribute
	value = attribs.Get("non_existent")

	// Assert that an empty string is returned for a non-existent attribute
	assert.Equal(t, "", value)
}
func TestAttributes_BatchingDuration(t *testing.T) {
	target := 10 * time.Millisecond
	// Create an instance of Attributes
	attribs := NewWithBatchDuration(
		"test_namespace",
		"test_service",
		"1.0.0",
		"test_env",
		target)

	// Get the batching duration
	duration := attribs.BatchingDuration()

	// Assert the expected batching duration
	assert.Equal(t, target, duration)

	attribs = New("test_namespace", "test_service", "1.0.0", "test_env")
	duration = attribs.BatchingDuration()
	assert.Equal(t, time.Duration(0), duration)
}
