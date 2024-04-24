# Tracing Package

The Tracing package is a Go library that provides a simple and convenient way to instrument your application with distributed tracing using OpenTelemetry. It allows you to initialize a tracer provider and tracer instance, which can be used to create and manage spans throughout your application.

## Features

- Easy initialization of the OpenTelemetry tracer provider and tracer instance
- Integration with common attribs for consistent span labeling
- Support for creating and managing spans
- Exporting traces to external systems using OpenTelemetry exporters
- Configurable batching duration for the tracing batch processor
- Sampling of traces based on a specified sampling rate

## Installation

To install the Tracing package, use the following command:

```
go get github.com/twistingmercury/telemetry/tracing
```

## Usage

### Initialization

Before using the Tracing package, you need to initialize it with an exporter, sampling rate, and common attribs:

```go
import (
    "github.com/twistingmercury/telemetry/attribs"
    "github.com/twistingmercury/telemetry/tracing"
)

exporter := // Create an OpenTelemetry exporter
sampleRate := 1.0 // Set the desired sampling rate
attribs := attribs.New("namespace", "service", "1.0.0", "production")

err := tracing.Initialize(exporter, sampleRate, attribs)
if err != nil {
    // Handle initialization error
}
```

- `exporter` is an instance of an OpenTelemetry exporter that will be used to export the collected traces. You can use any compatible exporter, such as Jaeger, Zipkin, or OTLP.
- `sampleRate` is a float value between 0 and 1 that determines the probability of a trace being sampled. A value of 1.0 means that all traces will be sampled, while a value of 0.5 means that approximately 50% of traces will be sampled.
- `attribs` is an instance of `attribs.Attributes` that contains common attribs to be included in every span.

### Creating and Managing Spans

After initialization, you can create and manage spans using the tracer instance:

```go
ctx := context.Background()
ctx = tracing.ExtractContext(ctx, carrier)

ctx, span := tracing.StartSpan(ctx, "my_span", oteltrace.SpanKindServer)
defer span.End()

// Perform some operations
```

- `tracing.ExtractContext` extracts the trace context from the provided carrier and returns a new context with the extracted trace information.
- `tracing.StartSpan` starts a new span with the given name and span kind, and returns a new context with the span attached and the created span.
- `span.End()` ends the span when the operation is complete.

You can add additional attribs to the span using the `oteltrace.WithAttributes` option when starting the span.

### Accessing the Tracer

The Tracing package provides a function to access the initialized tracer:

```go
tracer := tracing.Tracer()
```

- `tracing.Tracer()` returns the initialized tracer instance.

This function can be used to access the tracer from different parts of your application.

## Configuration

### Exporter

The Tracing package requires an OpenTelemetry exporter to be provided during initialization. You can configure the exporter based on your specific requirements, such as the export endpoint, protocol, and authentication.

Refer to the documentation of the specific OpenTelemetry exporter you are using for more details on configuring the exporter.

### Batching Duration

The Tracing package allows you to configure the batching duration for the tracing batch processor. The batching duration determines the maximum amount of time that spans are buffered before being exported.

To set the batching duration, use the `attribs.NewWithBatchDuration` function when creating the common attribs:

```go
attribs := attribs.NewWithBatchDuration("namespace", "service", "1.0.0", "production", 5*time.Second)
```

If the batching duration is not provided or set to 0, a default value of 5,000 milliseconds will be used.

### Sampling Rate

The Tracing package allows you to specify a sampling rate to control the percentage of traces that are sampled and exported. The sampling rate is a float value between 0 and 1, where 1.0 means that all traces will be sampled, and 0.5 means that approximately 50% of traces will be sampled.

The sampling rate is set during the initialization of the Tracing package:

```go
sampleRate := 1.0 // Set the desired sampling rate
err := tracing.Initialize(exporter, sampleRate, attribs)
```

## Contributing

Contributions to the Tracing package are welcome! If you find any issues or have suggestions for improvements, please open an issue or submit a pull request on the GitHub repository.

## License

The Tracing package is open-source and released under the [MIT License](../LICENSE).