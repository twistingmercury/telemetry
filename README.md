[![CodeQL](https://github.com/twistingmercury/telemetry/actions/workflows/codeql.yml/badge.svg)](https://github.com/twistingmercury/telemetry/actions/workflows/codeql.yml)
[![Run Tests](https://github.com/twistingmercury/telemetry/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/twistingmercury/telemetry/actions/workflows/go.yml)
[![codecov](https://codecov.io/github/twistingmercury/telemetry/graph/badge.svg?token=U6C4TE88OP)](https://codecov.io/github/twistingmercury/telemetry)
# Telemetry Package

This package was created to help me reduce a bunch of repetitive tasks in creating a Go application. All of the apps and services need logging, distributed tracing, and metrics.

## Features

- Logging: The package provides a logging system built on top of the [zerolog](https://pkg.go.dev/github.com/rs/zerolog)   library. It supports various log levels (debug, info, warn, error, fatal, panic) and allows adding custom key-value pairs to log messages. The logging system also integrates with OpenTelemetry tracing to include trace and span information in log messages.

- Metrics: The package utilizes [Prometheus](https://pkg.go.dev/github.com/prometheus/client_golang/prometheus) to collect metrics. It provides a simple way to initialize a metrics collector, which can be used to record metrics throughout the application.

- Tracing: The package integrates with OpenTelemetry tracing to collect and export trace data. It provides functions to initialize a tracer provider, extract trace context from incoming requests, and start new spans for outgoing requests or internal operations. The package allows configuring the batching duration for the tracing batch processor.

- Common Attributes: The package defines a common set of attribs that can be used across logging, metrics, and tracing. These attribs include service name, service version, namespace, and env. The package provides an `Attributes` interface and a default implementation for convenience.

## Installation

To install the Telemetry package, use the following command:

```
go get github.com/twistingmercury/telemetry
```

## Usage

### Common Attributes

To create a set of common attribs, use the `attribs.New` or `attribs.NewWithBatchDuration` function:

```go
import "github.com/twistingmercury/telemetry/attribs"

attribs := attribs.New("namespace", "service", "1.0.0", "production",
    attribute.String("custom_key", "custom_value"))

attribsWithBatchingDuration := attribs.NewWithBatchDuration("namespace", "service", "1.0.0", "production", 5*time.Second,
    attribute.String("custom_key", "custom_value"))
```

The common attribs can be passed to the initialization functions of logging, metrics, and tracing systems.

### Metrics

To use the metrics system, first initialize it with an exporter and common attribs:

```go
import "github.com/twistingmercury/telemetry/metrics"

err := metrics.InitializeWithPort(exporter, attribs)
if err != nil {
    // Handle initialization error
}
```

Then, you can create and record metrics using the meter instance:

```go
meter := metrics.Meter()
counter := meter.NewInt64Counter("my_counter")
counter.Add(ctx, 1, attribute.String("key", "value"))
```

### Tracing

To use the tracing system, first initialize it with an exporter, sampling rate, and common attribs:

```go
import "github.com/twistingmercury/telemetry/tracing"

err := tracing.InitializeWithPort(exporter, 1.0, attribs)
if err != nil {
    // Handle initialization error
}
```

Then, you can create and manage spans using the tracer instance:

```go
ctx := tracing.ExtractContext(ctx, carrier)
ctx, span := tracing.StartSpan(ctx, "my_span", oteltrace.SpanKindServer)
defer span.End()
```

## Contributing

Contributions to the Telemetry package are welcome! If you find any issues or have suggestions for improvements, please open an issue or submit a pull request on the GitHub repository.

## License

The Telemetry package is open-source and released under the [MIT License](LICENSE).