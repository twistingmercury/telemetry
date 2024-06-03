package logging_test

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/twistingmercury/telemetry/logging"

	"go.opentelemetry.io/otel/trace"
)

const (
	serviceName    = "test-service"
	serviceVersion = "1.0.0"
	environment    = "unit-test"
)

func TestInitialize(t *testing.T) {
	var buf bytes.Buffer
	err := logging.Initialize(zerolog.DebugLevel, &buf, serviceName, serviceVersion, environment)
	assert.NoError(t, err, "Initialize should not return an error")
}

func TestInitializeWithNilWriter(t *testing.T) {
	err := logging.Initialize(zerolog.DebugLevel, nil, serviceName, serviceVersion, environment)
	assert.Error(t, err, "Initialize should return an error when writer is nil")
}

func TestLoggingWithSpanContext(t *testing.T) {
	defer func() {
		logging.SetExitFunc(os.Exit)
	}()
	logging.SetExitFunc(func(int) {})
	// Create a buffer to capture the log output
	var buf bytes.Buffer

	err := logging.Initialize(zerolog.DebugLevel, &buf, serviceName, serviceVersion, environment)
	assert.NoError(t, err, "Initialize should not return an error")

	traceID := trace.TraceID{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10}
	spanID := trace.SpanID{0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18}
	spanCtx := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: traceID,
		SpanID:  spanID,
	})

	logging.DebugWithContext(&spanCtx, "Debug message", logging.KeyValue{Key: "key1", Value: "value1"})
	logging.InfoWithContext(&spanCtx, "Info message", logging.KeyValue{Key: "key2", Value: 123})
	logging.WarnWithContext(&spanCtx, "Warn message", logging.KeyValue{Key: "key3", Value: true})
	logging.ErrorWithContext(&spanCtx, errors.New("test error"), "Error message", logging.KeyValue{Key: "key4", Value: 3.14})
	logging.FatalWithContext(&spanCtx, errors.New("test error"), "Fatal message", logging.KeyValue{Key: "key5", Value: 6.28})

	pfunc := func() {
		logging.PanicWithContext(&spanCtx, errors.New("test panic"), "Panic message", logging.KeyValue{Key: "key6", Value: []string{"a", "b", "c"}})
	}
	assert.Panics(t, pfunc, "Panic should be called")

	logOutput := buf.String()
	assert.Contains(t, logOutput, "Debug message")
	assert.Contains(t, logOutput, "Info message")
	assert.Contains(t, logOutput, "Warn message")
	assert.Contains(t, logOutput, "Error message")
	assert.Contains(t, logOutput, "Fatal message")
	assert.Contains(t, logOutput, "Panic message")
	assert.Contains(t, logOutput, fmt.Sprintf(`"%s":"%s"`, logging.TraceIDAttr, traceID.String()))
	assert.Contains(t, logOutput, fmt.Sprintf(`"%s":"%s"`, logging.SpanIDAttr, spanID.String()))
	assert.Contains(t, logOutput, `"key1":"value1"`)
	assert.Contains(t, logOutput, `"key2":123`)
	assert.Contains(t, logOutput, `"key3":true`)
	assert.Contains(t, logOutput, `"key4":3.14`)
	assert.Contains(t, logOutput, `"key5":6.28`)
	assert.Contains(t, logOutput, `"key6":["a","b","c"]`)
}

func TestLoggingWithNilSpanContext(t *testing.T) {
	defer func() {
		logging.SetExitFunc(os.Exit)
	}()
	logging.SetExitFunc(func(int) {})

	var buf bytes.Buffer

	err := logging.Initialize(zerolog.DebugLevel, &buf, serviceName, serviceVersion, environment)
	assert.NoError(t, err, "Initialize should not return an error")

	// Log messages with different levels using a nil span context
	logging.DebugWithContext(nil, "Debug message", logging.KeyValue{Key: "key1", Value: "value1"})
	logging.InfoWithContext(nil, "Info message", logging.KeyValue{Key: "key2", Value: 123})
	logging.WarnWithContext(nil, "Warn message", logging.KeyValue{Key: "key3", Value: true})
	logging.ErrorWithContext(nil, errors.New("test error"), "Error message", logging.KeyValue{Key: "key4", Value: 3.14})
	logging.FatalWithContext(nil, errors.New("test error"), "Fatal message", logging.KeyValue{Key: "key5", Value: 6.28})

	pfunc := func() {
		logging.PanicWithContext(nil, errors.New("test panic"), "Panic message", logging.KeyValue{Key: "key6", Value: []string{"a", "b", "c"}})
	}
	assert.Panics(t, pfunc, "Panic should be called")

	logOutput := buf.String()
	assert.Contains(t, logOutput, "Debug message")
	assert.Contains(t, logOutput, "Info message")
	assert.Contains(t, logOutput, "Warn message")
	assert.Contains(t, logOutput, "Error message")
	assert.Contains(t, logOutput, "Panic message")
	assert.NotContains(t, logOutput, logging.TraceIDAttr)
	assert.NotContains(t, logOutput, logging.SpanIDAttr)
	assert.Contains(t, logOutput, `"key1":"value1"`)
	assert.Contains(t, logOutput, `"key2":123`)
	assert.Contains(t, logOutput, `"key3":true`)
	assert.Contains(t, logOutput, `"key4":3.14`)
	assert.Contains(t, logOutput, `"key5":6.28`)
	assert.Contains(t, logOutput, `"key6":["a","b","c"]`)
}

func TestLoggingWithoutSpanContext(t *testing.T) {
	defer func() {
		logging.SetExitFunc(os.Exit)
	}()
	logging.SetExitFunc(func(int) {})

	var buf bytes.Buffer

	err := logging.Initialize(zerolog.DebugLevel, &buf, serviceName, serviceVersion, environment)
	assert.NoError(t, err, "Initialize should not return an error")

	logging.Debug("Debug message", logging.KeyValue{Key: "key1", Value: "value1"})
	logging.Info("Info message", logging.KeyValue{Key: "key2", Value: 123})
	logging.Warn("Warn message", logging.KeyValue{Key: "key3", Value: true})
	logging.Error(errors.New("test error"), "Error message", logging.KeyValue{Key: "key4", Value: 3.14})
	logging.Fatal(errors.New("test error"), "Fatal message", logging.KeyValue{Key: "key5", Value: 6.28})
	pfunc := func() {
		logging.Panic(errors.New("test panic"), "Panic message", logging.KeyValue{Key: "key6", Value: []string{"a", "b", "c"}})
	}
	assert.Panics(t, pfunc, "Panic should be called")

	logOutput := buf.String()
	assert.Contains(t, logOutput, "Debug message")
	assert.Contains(t, logOutput, "Info message")
	assert.Contains(t, logOutput, "Warn message")
	assert.Contains(t, logOutput, "Error message")
	assert.Contains(t, logOutput, "Fatal message")
	assert.Contains(t, logOutput, "Panic message")
	assert.Contains(t, logOutput, `"key1":"value1"`)
	assert.Contains(t, logOutput, `"key2":123`)
	assert.Contains(t, logOutput, `"key3":true`)
	assert.Contains(t, logOutput, `"key4":3.14`)
	assert.Contains(t, logOutput, `"key5":6.28`)
	assert.Contains(t, logOutput, `"key6":["a","b","c"]`)
}

func TestTraceInfoWithInvalidTraceID(t *testing.T) {
	spanCtx := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: trace.TraceID{}, // Empty trace ID
		SpanID:  trace.SpanID{0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18},
	})

	tMap := logging.TraceInfo(&spanCtx)

	assert.Emptyf(t, tMap, "TraceInfo should return an empty map when trace ID is invalid")
}
