package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/twistingmercury/telemetry/metrics"
	"github.com/twistingmercury/telemetry/tracing"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	otelmetric "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"log"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/twistingmercury/telemetry/attributes"
	"github.com/twistingmercury/telemetry/logging"

	"go.opentelemetry.io/otel/attribute"
)

var (
	meterCounter otelmetric.Int64Counter
)

const (
	namespace   = "example"
	serviceName = "monex"
	version     = "0.0.1"
	environment = "dev"
)

func main() {
	otherAttribs := []attribute.KeyValue{
		attribute.String("build-date", time.Now().Format(time.RFC3339)),
		attribute.String("commit", "ca45gf32"),
		attribute.Int("my-attribute", 123),
	}

	// 1. Create the commonly used attributes for logging, metrics, and tracing
	opt := attributes.NewWithBatchDuration(
		namespace,
		serviceName,
		version,
		environment,
		100*time.Millisecond, // set very short batching duration for demonstration purposes
		otherAttribs...)

	// 2. initialize logging functionality
	//    from this point forward use the logging package to log messages
	err := logging.Initialize(zerolog.DebugLevel, opt, os.Stdout)
	if err != nil {
		log.Panicf("failed to initialize old_elemetry: %s", err)
	}

	metrics.Initialize("9090", "complete", "example")

	// 4. Initialize the tracing functionality
	tex, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		// fail fast!
		logging.Panic(err, "failed to create trace exporter")
	}

	err = tracing.Initialize(tex, 1.0, opt)
	if err != nil {
		// fail fast!
		logging.Panic(err, "failed to initialize tracing")
	}

	stopChan := make(chan struct{})
	defer close(stopChan)
	echo(stopChan)
}

func echo(stopChan chan struct{}) {
	// Make sure to start the span before logging so that the span context is available. This is important
	// for the trace to be able to corrolate the logs to the trace.
	_, span := tracing.StartSpan(context.Background(), "echo", trace.SpanKindServer)
	sctx := span.SpanContext()
	logging.InfoWithContext(&sctx, "echoing")

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("press enter to continue example> ")
	_, _ = reader.ReadString('\n')

	meterCounter.Add(context.Background(), 1)
	span.SetStatus(codes.Ok, "OK")
	span.End()

	<-stopChan // this is just to make it wait so that the telemetry is sent to stdout
}
