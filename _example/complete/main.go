package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/twistingmercury/telemetry/logging"
	"github.com/twistingmercury/telemetry/metrics"
	"github.com/twistingmercury/telemetry/tracing"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	otelmetric "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"log"
	"os"
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
	// 1. initialize logging
	//    from this point forward use the logging package to log messages
	if err := logging.Initialize(zerolog.DebugLevel, os.Stdout, serviceName, version, environment); err != nil {
		log.Panicf("failed to initialize old_elemetry: %s", err)
	}

	// 2. initialize metrics
	if err := metrics.Initialize("complete", "example"); err != nil {
		log.Panicf("failed to initialize metrics: %s", err)
	}

	// 3. initialize tracing
	tex, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		// fail fast!
		logging.Panic(err, "failed to create trace exporter")
	}

	if err := tracing.Initialize(tex, serviceName, version, environment); err != nil {
		logging.Panic(err, "failed to initialize tracing")
	}

	stopChan := make(chan struct{})
	defer close(stopChan)
	echo(stopChan)
}

func echo(stopChan chan struct{}) {
	// Make sure to start the span before logging so that the span context is available. This is important
	// for the trace to be able to corrolate the logs to the trace.
	_, span := tracing.Start(context.Background(), "echo", trace.SpanKindServer)
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
