# Metrics Package

The Metrics package is a Go library that provides a simple and convenient way to collect and export application metrics using OpenTelemetry. It allows you to initialize a meter provider and meter instance, which can be used to create and record metrics throughout your application.

## Features

- Easy initialization of the OpenTelemetry meter provider and meter instance
- Integration with common attribs for consistent metric labeling
- Support for various metric types, such as counters, gauges, and histograms
- Exporting metrics to external systems using OpenTelemetry exporters
- Configurable batching duration for the metrics periodic reader

## Installation

To install the Metrics package, use the following command:

```
go get github.com/twistingmercury/telemetry/metrics
```

## Usage

### Initialization

Before using the Metrics package, you need to initialize it with an exporter and common attribs:

```go
import (
    "github.com/twistingmercury/telemetry/attribs"
    "github.com/twistingmercury/telemetry/metrics"
)

exporter := // Create an OpenTelemetry exporter
attribs := attribs.New("namespace", "service", "1.0.0", "production")

err := metrics.Initialize(exporter, attribs)
if err != nil {
    // Handle initialization error
}
```

- `exporter` is an instance of an OpenTelemetry exporter that will be used to export the collected metrics. You can use any compatible exporter, such as Prometheus, Jaeger, or OTLP.
- `attribs` is an instance of `attribs.Attributes` that contains common attribs to be included in every metric.

### Creating and Recording Metrics

After initialization, you can create and record metrics using the meter instance:

```go
meter := metrics.Meter()
counter := meter.NewInt64Counter("my_counter")
counter.Add(ctx, 1, attribute.String("key", "value"))
```

- `meter.NewInt64Counter` creates a new int64 counter metric with the given name.
- `counter.Add` records a value of 1 for the counter metric, with an additional attribute key-value pair.

The Metrics package supports various metric types, such as counters, gauges, and histograms. Refer to the OpenTelemetry Go SDK documentation for more details on creating and recording different types of metrics.

### Accessing the Meter and Meter Provider

The Metrics package provides functions to access the initialized meter and meter provider:

```go
meter := metrics.Meter()
meterProvider := metrics.MeterProvider()
```

- `metrics.Meter()` returns the initialized meter instance.
- `metrics.MeterProvider()` returns the initialized meter provider instance.

These functions can be used to access the meter and meter provider from different parts of your application.

## Configuration

### Exporter

The Metrics package requires an OpenTelemetry exporter to be provided during initialization. You can configure the exporter based on your specific requirements, such as the export interval, endpoint, and authentication.

Refer to the documentation of the specific OpenTelemetry exporter you are using for more details on configuring the exporter.

### Batching Duration

The Metrics package allows you to configure the batching duration for the metrics periodic reader. The batching duration determines the frequency at which the collected metrics are exported.

To set the batching duration, use the `attribs.NewWithBatchDuration` function when creating the common attribs:

```go
attribs := attribs.NewWithBatchDuration("namespace", "service", "1.0.0", "production", 5*time.Second)
```

If the batching duration is not provided or set to 0, a default value of 60 seconds will be used.

## Contributing

Contributions to the Metrics package are welcome! If you find any issues or have suggestions for improvements, please open an issue or submit a pull request on the GitHub repository.

## License

The Metrics package is open-source and released under the [MIT License](LICENSE).