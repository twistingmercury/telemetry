package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/twistingmercury/telemetry/attributes"
	"github.com/twistingmercury/telemetry/metrics"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	otelmetric "go.opentelemetry.io/otel/metric"
)

func main() {
	// Create attributes attributes
	attribs := attributes.NewWithBatchDuration(
		"namespace",
		"service",
		"1.0.0",
		"production",
		1*time.Second, // set very short batching duration for demonstration purposes
		attribute.String("custom_key", "custom_value"))

	// Create an OpenTelemetry stdout exporter
	exporter, err := stdoutmetric.New()
	if err != nil {
		log.Fatal("failed to create stdout exporter:", err)
	}

	// Initialize the Metrics package with the exporter and common attributes
	err = metrics.Initialize(exporter, attribs)
	if err != nil {
		log.Fatal("failed to initialize metrics package:", err)
	}

	// Get the meter instance
	meter := metrics.Meter()

	// Create a new int64 counter metric
	counter, err := meter.Int64Counter(
		"example_counter",
		otelmetric.WithDescription("An example counter metric"))
	if err != nil {
		log.Fatal("failed to create counter metric: ", err)
	}

	// Create a new float64 histogram metric
	histogram, err := meter.Float64Histogram(
		"example_histogram",
		otelmetric.WithDescription("An example histogram metric"))
	if err != nil {
		log.Fatal("failed to create histogram metric: ", err)
	}
	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Simulate some work and record metrics
	for i := 0; i < 5; i++ {
		// Simulate some work
		time.Sleep(1 * time.Second)

		// Record a value for the counter metric
		counter.Add(ctx, 1)

		// Record a value for the histogram metric
		value := rand.Float64() * 100
		histogram.Record(ctx, value)

		log.Printf("Recorded metrics - Counter: 1, Histogram: %.2f", value)
	}
}
