package tracing_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twistingmercury/telemetry/tracing"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
	"testing"
)

const (
	serviceName    = "test-service"
	serviceVersion = "1.0.0"
	environment    = "unit-test"
)

type mockExporter struct {
	sdktrace.SpanExporter
}

func TestInitialize(t *testing.T) {
	exporter := &mockExporter{}

	err := tracing.InitializeWithSampleRate(exporter, 1.0, serviceName, serviceVersion, environment)
	require.NoError(t, err)

	tracer := tracing.Tracer()
	assert.NotNil(t, tracer)

	err = tracing.InitializeWithSampleRate(exporter, 0.0, serviceName, serviceVersion, environment)
	require.Error(t, err)

}

func TestInitializeWithNilExporter(t *testing.T) {
	err := tracing.Initialize(nil, serviceName, serviceVersion, environment)
	assert.Error(t, err, "InitializeWithSampleRate should return an error when exporter is nil")
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

func TestExtractSpan(t *testing.T) {
	ctx := context.Background()

	extractedCtx := tracing.ExtractSpan(ctx)

	assert.NotNil(t, extractedCtx, "extracted context should not be nil")
	assert.Equal(t, false, extractedCtx.SpanContext().IsValid(), "extracted context should not be valid")
}

func TestStart(t *testing.T) {
	// Create a mock exporter and attributes
	exporter := new(mockExporter)

	// InitializeWithSampleRate the tracing package
	err := tracing.Initialize(exporter, serviceName, serviceVersion, environment)
	require.NoError(t, err, "InitializeWithSampleRate should not return an error")

	// Create a test context
	ctx := context.Background()

	// StartSpan a span using the StartSpan function
	spanCtx, span := tracing.Start(ctx, "test-span", oteltrace.SpanKindServer)
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
