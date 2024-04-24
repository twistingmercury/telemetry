package tracing_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twistingmercury/telemetry/attributes"
	"github.com/twistingmercury/telemetry/tracing"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
	"testing"
)

type mockExporter struct {
	sdktrace.SpanExporter
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
	attribs := testAttributes()

	err := tracing.Initialize(exporter, 1.0, attribs)
	require.NoError(t, err)

	tracer := tracing.Tracer()
	assert.NotNil(t, tracer)
}

func TestInitializeWithNilExporter(t *testing.T) {
	attribs := testAttributes()
	err := tracing.Initialize(nil, 1.0, attribs)
	assert.Error(t, err, "Initialize should return an error when exporter is nil")
}

func TestInitializeWithNilAttributes(t *testing.T) {
	exporter := new(mockExporter)
	err := tracing.Initialize(exporter, 1.0, nil)
	assert.Error(t, err, "Initialize should return an error when attributes is nil")
}

func TestExtractContext(t *testing.T) {
	// Create a test context and carrier
	ctx := context.Background()
	carrier := propagation.HeaderCarrier{}

	// Extract the context using the ExtractContext function
	extractedCtx := tracing.ExtractContext(ctx, carrier)

	// Assert that the extracted context is not nil
	assert.NotNil(t, extractedCtx, "extracted context should not be nil")
}
func TestStartSpan(t *testing.T) {
	// Create a mock exporter and attributes
	exporter := new(mockExporter)
	attribs := testAttributes()

	// Initialize the tracing package
	err := tracing.Initialize(exporter, 1.0, attribs)
	require.NoError(t, err, "Initialize should not return an error")

	// Create a test context
	ctx := context.Background()

	// Start a span using the StartSpan function
	spanCtx, span := tracing.StartSpan(ctx, "test-span", oteltrace.SpanKindServer)
	defer span.End()

	// Assert that the span context and span are not nil
	assert.NotNil(t, spanCtx, "span context should not be nil")
	assert.NotNil(t, span, "span should not be nil")

	// Assert that the span context contains a valid trace ID and span ID
	assert.NotEqual(t, oteltrace.TraceID{}, span.SpanContext().TraceID(), "span context should have a non-zero trace ID")
	assert.NotEqual(t, oteltrace.SpanID{}, span.SpanContext().SpanID(), "span context should have a non-zero span ID")

	// Assert that the trace ID and span ID are not empty (all zeros)
	assert.NotEqual(t, oteltrace.TraceID{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, span.SpanContext().TraceID(), "trace ID should not be empty")
	assert.NotEqual(t, oteltrace.SpanID{0, 0, 0, 0, 0, 0, 0, 0}, span.SpanContext().SpanID(), "span ID should not be empty")
}
