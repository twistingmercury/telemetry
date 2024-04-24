package metrics_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twistingmercury/telemetry/attributes"
	"github.com/twistingmercury/telemetry/metrics"
	"go.opentelemetry.io/otel/attribute"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"testing"
)

type mockExporter struct {
	sdkmetric.Exporter
}

func testAttributes() attributes.Attributes {
	return attributes.New(
		"test_namespace",
		"test_service",
		"1.0.0",
		"test_env",
		attribute.String("key1", "value1"),
		attribute.Int("key2", 123),
	)
}

func TestInitialize(t *testing.T) {
	exporter := &mockExporter{}
	attribs := attributes.New("namespace", "service", "1.0.0", "production")

	err := metrics.Initialize(exporter, attribs)
	require.NoError(t, err)

	meter := metrics.Meter()
	assert.NotNil(t, meter)

	meterProvider := metrics.MeterProvider()
	assert.NotNil(t, meterProvider)
}

func TestInitializeWithNilExporter(t *testing.T) {
	attribs := testAttributes()
	err := metrics.Initialize(nil, attribs)
	assert.Error(t, err, "Initialize should return an error when exporter is nil")
}

func TestInitializeWithNilAttributes(t *testing.T) {
	exporter := new(mockExporter)
	err := metrics.Initialize(exporter, nil)
	assert.Error(t, err, "Initialize should return an error when attributes is nil")
}
