package main

import (
	"errors"
	"os"

	"github.com/rs/zerolog"
	"github.com/twistingmercury/telemetry/attributes"
	"github.com/twistingmercury/telemetry/logging"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func main() {
	// Create attributes attributes
	attribs := attributes.New("namespace", "service", "1.0.0", "production",
		attribute.String("custom_key", "custom_value"))

	// Initialize the logging system
	err := logging.Initialize(zerolog.DebugLevel, attribs, os.Stdout)
	if err != nil {
		panic(err)
	}

	// Log messages without context
	logging.Debug("Debug message", logging.KeyValue{Key: "key", Value: "value"})
	logging.Info("Info message", logging.KeyValue{Key: "key", Value: 123})
	logging.Warn("Warn message")

	err = errors.New("something bad happened, whatever it was")
	logging.Error(err, "Error message")

	// mock a span context; under normal circumstances, this would be passed in from the active trace.Span
	spanCtx := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: trace.TraceID{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10},
		SpanID:  trace.SpanID{0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18},
	})

	// Log messages with context
	kv1 := logging.KeyValue{Key: "key", Value: "value"}
	kv2 := logging.KeyValue{Key: "key", Value: 123}
	logging.DebugWithContext(&spanCtx, "Debug message with context", kv1, kv2)
	logging.InfoWithContext(&spanCtx, "Info message with context")
	logging.WarnWithContext(&spanCtx, "Warn message with context")
	logging.ErrorWithContext(&spanCtx, err, "Error message with context")

	// Log fatal message
	//logging.Fatal(nil, "Fatal message")

	// Log panic message
	//logging.Panic(nil, "Panic message")
}
