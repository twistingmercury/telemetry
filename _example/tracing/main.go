package main

import (
	"bufio"
	"context"
	"log"
	"os"
	"time"

	"github.com/twistingmercury/telemetry/tracing"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	oteltrace "go.opentelemetry.io/otel/trace"
)

const (
	serviceName = "scooby"
	version     = "0.0.1"
	environment = "local"
)

func main() {
	// Create common attributes

	// Create an OpenTelemetry stdout exporter
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		log.Fatal("Failed to create stdout exporter:", err)
	}

	// InitializeWithSampleRate the Tracing package with the exporter, sampling rate, and common attributes
	err = tracing.Initialize(exporter, serviceName, version, environment)
	if err != nil {
		log.Fatal("failed to initialize tracing package:", err)
	}

	// Create a context
	ctx := context.Background()

	// StartSpan a new span
	ctx, span := tracing.Start(ctx, "main", oteltrace.SpanKindServer,
		attribute.String("operation", "example"))
	defer span.End()

	// Simulate some work
	log.Println("Doing some work...")
	time.Sleep(1 * time.Second)

	// Start a child span
	ctx, childSpan := tracing.Start(ctx, "child_operation", oteltrace.SpanKindServer,
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
	log.Print("Give it a few moments to process the traces. press ctrl+c to exit> ")
	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n')
}
