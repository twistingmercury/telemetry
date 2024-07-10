package logging

import (
	"context"
	"errors"
	"go.opentelemetry.io/otel/trace"
	"io"
	"os"

	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

const (
	TraceIDAttr = "otel.trace_id"
	SpanIDAttr  = "otel.span_id"
)

var (
	logger   zerolog.Logger
	exitFunc = os.Exit
)

type KeyValue struct {
	Key   string
	Value any
}

func toMap(values ...KeyValue) map[string]any {
	m := make(map[string]any, len(values))
	for _, v := range values {
		m[v.Key] = v.Value
	}
	return m
}

// Initialize initializes the logging system.
// It returns a logger that can be used to log messages, though it is not required.
func Initialize(level zerolog.Level, writer io.Writer, serviceName, serviceVersion, environment string) (err error) {

	if writer == nil {
		return errors.New("writer is required")
	}

	zerolog.SetGlobalLevel(level)
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	logger = zerolog.New(writer).
		With().
		Timestamp().
		Str("service", serviceName).
		Str("version", serviceVersion).
		Str("environment", environment).
		Logger()

	return
}

// DebugWithContext logs a debug message and adds the trace id and span id found in the ctx.
// The args are key value pairs and are optional.
//
// Deprecated: use Debug func instead. It extracts tracing data from the context automatically.
func DebugWithContext(spanCtx *trace.SpanContext, message string, args ...KeyValue) {
	tInf := traceInfo(spanCtx)
	margs := MergeMaps(toMap(args...), tInf)
	logger.Debug().
		Fields(margs).
		Msg(message)
}

// InfoWithContext logs an info message and adds the trace id and span id found in the ctx.
// The args are key value pairs and are optional.
//
// Deprecated: use Info func instead. It extracts tracing data from the context automatically.
func InfoWithContext(spanCtx *trace.SpanContext, message string, args ...KeyValue) {
	tInf := traceInfo(spanCtx)
	margs := MergeMaps(toMap(args...), tInf)
	logger.Info().
		Fields(margs).
		Msg(message)
}

// WarnWithContext logs a warning message and adds the trace id and span id found in the ctx.
// The args are key value pairs and are optional.
//
// Deprecated: use Warn func instead. It extracts tracing data from the context automatically.
func WarnWithContext(spanCtx *trace.SpanContext, message string, args ...KeyValue) {
	tInf := traceInfo(spanCtx)
	margs := MergeMaps(toMap(args...), tInf)
	logger.Warn().
		Fields(margs).
		Msg(message)
}

// ErrorWithContext logs an error message and adds the trace id and span id found in the ctx.
//
// Deprecated: use Error func instead. It extracts tracing data from the context automatically.
func ErrorWithContext(spanCtx *trace.SpanContext, err error, message string, args ...KeyValue) {
	tInf := traceInfo(spanCtx)
	margs := MergeMaps(toMap(args...), tInf)
	logger.Error().
		Fields(margs).
		Err(err).
		Str("is-fatal", "false").
		Msg(message)
}

// FatalWithContext logs a fatal message and adds the trace id and span id found in the ctx.
//
// Deprecated: use Fatal func instead. It extracts tracing data from the context automatically.
func FatalWithContext(spanCtx *trace.SpanContext, err error, message string, args ...KeyValue) {
	tInf := traceInfo(spanCtx)
	margs := MergeMaps(toMap(args...), tInf)
	logger.Error().
		Fields(margs).
		Err(err).
		Str("is-fatal", "true").
		Msg(message)
	exitFunc(1)
}

func PanicWithContext(spanCtx *trace.SpanContext, err error, message string, args ...KeyValue) {
	tInf := traceInfo(spanCtx)
	margs := MergeMaps(toMap(args...), tInf)
	logger.Panic().
		Fields(margs).
		Err(err).
		Msg(message)
}

// traceInfo returns the trace id and span id found in the ctx.
func traceInfo(spanCtx *trace.SpanContext) (tMap map[string]any) {
	tMap = make(map[string]any)
	if spanCtx == nil {
		return
	}

	if !spanCtx.IsValid() || !spanCtx.IsSampled() {
		return
	}

	tMap[TraceIDAttr] = spanCtx.TraceID().String()
	tMap[SpanIDAttr] = spanCtx.SpanID().String()
	return
}

// MergeMaps takes any two maps and combines them.
func MergeMaps(m1 map[string]any, m2 map[string]any) map[string]any {
	merged := make(map[string]any)
	for k, v := range m1 {
		merged[k] = v
	}
	for key, value := range m2 {
		merged[key] = value
	}
	return merged
}

// Debug logs a debug message. Tracing data (if present) is automatically retrieved from the [context.Context].
func Debug(ctx context.Context, message string, args ...KeyValue) {
	fields := MergeMaps(toMap(args...), getTracingAttributes(ctx))

	logger.Debug().
		Fields(fields).
		Msg(message)
}

// Info logs an info message. Tracing data (if present) is automatically retrieved from the [context.Context].
func Info(ctx context.Context, message string, args ...KeyValue) {
	fields := MergeMaps(toMap(args...), getTracingAttributes(ctx))

	logger.Info().
		Fields(fields).
		Msg(message)
}

// Warn logs a warning message. Tracing data (if present) is automatically retrieved from the [context.Context].
func Warn(ctx context.Context, message string, args ...KeyValue) {
	fields := MergeMaps(toMap(args...), getTracingAttributes(ctx))

	logger.Warn().
		Fields(fields).
		Msg(message)
}

// Error logs an error message. Tracing data (if present) is automatically retrieved from the [context.Context].
func Error(ctx context.Context, err error, message string, args ...KeyValue) {
	fields := MergeMaps(toMap(args...), getTracingAttributes(ctx))

	logger.Error().
		Fields(fields).
		Err(err).
		Str("is-fatal", "false").
		Msg(message)
}

// Fatal logs a fatal message. Tracing data (if present) is automatically retrieved from the [context.Context].
func Fatal(ctx context.Context, err error, message string, args ...KeyValue) {
	fields := MergeMaps(toMap(args...), getTracingAttributes(ctx))

	logger.Error().
		Fields(fields).
		Err(err).
		Str("is-fatal", "true").
		Msg(message)
	exitFunc(1)
}

// Panic logs a panic
func Panic(ctx context.Context, err error, message string, args ...KeyValue) {
	fields := MergeMaps(toMap(args...), getTracingAttributes(ctx))

	logger.Panic().
		Fields(fields).
		Err(err).
		Msg(message)
}

// getTracingAttributes retrieves tracing data from [context.Context]
func getTracingAttributes(ctx context.Context) map[string]any {
	span := trace.SpanFromContext(ctx)
	spanCtx := span.SpanContext()

	return traceInfo(&spanCtx)
}
