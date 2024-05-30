package main

import (
	"bufio"
	"context"
	"example/metrics/data"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/twistingmercury/telemetry/logging"
	"github.com/twistingmercury/telemetry/metrics"
	"github.com/twistingmercury/telemetry/tracing"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/trace"
	"log"
	"os"
	"time"
)

const (
	namespace   = "example"
	serviceName = "scooby"
	version     = "0.0.1"
	environment = "local"
)

func main() {
	// 1. initialize logging
	//    from this point forward use the logging package to log messages
	if err := logging.Initialize(zerolog.DebugLevel, os.Stdout, serviceName, version, environment); err != nil {
		log.Fatalf("failed to initialize old_elemetry: %s", err)
	}

	// 3. initialize tracing
	tex, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		// fail fast!
		logging.Fatal(err, "failed to create trace exporter")
	}

	if err := tracing.Initialize(tex, serviceName, version, environment); err != nil {
		logging.Fatal(err, "failed to initialize tracing")
	}

	// 2. initialize metrics
	if err := metrics.Initialize(namespace, serviceName); err != nil {
		log.Fatalf("failed to initialize metrics: %s", err)
	}

	metrics.RegisterMetrics(data.Metrics()...)

	metrics.Publish()

	for i := 0; i < 5; i++ {
		_ = data.DoDatabaseStuff()
		time.Sleep(250 * time.Millisecond)
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

	span.SetStatus(codes.Ok, "OK")
	span.End()

	<-stopChan // this is just to make it wait so that the telemetry is sent to stdout
}
