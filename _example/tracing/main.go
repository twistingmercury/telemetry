package main

import (
	"context"
	"log"
	"time"

	"github.com/twistingmercury/telemetry/attributes"
	"github.com/twistingmercury/telemetry/tracing"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	oteltrace "go.opentelemetry.io/otel/trace"
)

func main() {
	// Create common attributes
	attribs := attributes.NewWithBatchDuration(
		"namespace",
		"service",
		"1.0.0",
		"production",
		10*time.Millisecond, // set very short batching duration for demonstration purposes
		attribute.String("custom_key", "custom_value"))

	// Create an OpenTelemetry stdout exporter
	exporter, err := stdouttrace.New()
	if err != nil {
		log.Fatal("Failed to create stdout exporter:", err)
	}

	// Initialize the Tracing package with the exporter, sampling rate, and common attributes
	err = tracing.Initialize(exporter, 1.0, attribs)
	if err != nil {
		log.Fatal("Failed to initialize Tracing package:", err)
	}

	// Create a context
	ctx := context.Background()

	// Start a new span
	ctx, span := tracing.StartSpan(ctx, "main", oteltrace.SpanKindServer,
		attribute.String("operation", "example"))
	defer span.End()

	// Simulate some work
	log.Println("Doing some work...")
	time.Sleep(1 * time.Second)

	// Start a child span
	ctx, childSpan := tracing.StartSpan(ctx, "child_operation", oteltrace.SpanKindServer,
		attribute.String("child_key", "child_value"))
	defer childSpan.End()

	// Simulate some work in the child span
	log.Println("Doing some work in the child span...")
	time.Sleep(500 * time.Millisecond)

	// Add an event to the child span
	childSpan.AddEvent("child_event", oteltrace.WithAttributes(
		attribute.String("event_key", "event_value"),
	))

	// End the child span
	childSpan.End()

	// Add an event to the main span
	span.AddEvent("main_event", oteltrace.WithAttributes(
		attribute.String("event_key", "event_value"),
	))

	// End the main span
	span.End()

	log.Println("Tracing example completed. press ctrl+c to exit.")

	// Sleep for a while to allow the exporter to process the spans
	time.Sleep(1 * time.Second)
}
