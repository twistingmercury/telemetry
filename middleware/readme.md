# Gin Middleware Package

The Gin Middleware package is a Go library that provides middleware for instrumenting and tracing incoming requests in a Gin web application. It integrates with the OpenTelemetry framework to collect and export telemetry data, including metrics and traces.

## Features

- Middleware for instrumenting incoming requests in a Gin web application
- Integration with OpenTelemetry for collecting and exporting metrics and traces
- Automatic generation of request duration and count metrics
- Parsing of user agent and request headers for detailed telemetry data
- Customizable span naming and attribute generation
- Correlation of logs with trace and span information

## Installation

To install the Gin Middleware package, use the following command:

```
go get github.com/twistingmercury/telemetry/middleware
```

## Usage

To use the Gin Middleware package, you need to create an instance of `attributes.Attributes` and initialize the metrics and tracing packages. Then, you can register the `Telemetry` with your Gin router.

Here's an example of how to use the package:

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/twistingmercury/telemetry/attributes"
    "github.com/twistingmercury/telemetry/middleware"
    "github.com/twistingmercury/telemetry/metrics"
    "github.com/twistingmercury/telemetry/tracing"
)

func main() {
    // Create an instance of attributes
    attribs := attributes.New("namespace", "service", "1.0.0", "production")

    // Initialize metrics and tracing
    metricsExporter := // Create a metrics exporter
    traceExporter := // Create a trace exporter
    err := metrics.Initialize(metricsExporter, attribs)
    if err != nil {
        // Handle initialization error
    }
    err = tracing.Initialize(traceExporter, 1.0, attribs)
    if err != nil {
        // Handle initialization error
    }

    // Create a new Gin router
    router := gin.New()

    // Register the Telemetry
    router.Use(middleware.Telemetry(attribs))

    // Define your routes and handlers
    // ...

    // Start the server
    router.Run(":8080")
}
```

In this example, we create an instance of `attributes.Attributes` with the desired namespace, service name, version, and environment. We then initialize the metrics and tracing packages with their respective exporters and the attributes.

Next, we create a new Gin router and register the `Telemetry` using `router.Use(gin.Telemetry(attribs))`.

After that, you can define your routes and handlers as usual, and the middleware will automatically instrument and trace the incoming requests.

## Configuration

The `attributes.Attributes` instance allows you to configure various aspects of the middleware, such as the namespace, service name, version, and environment. You can also provide additional custom attributes using the `attribute.KeyValue` pairs.

The middleware uses the provided attributes to generate metric names, set span attributes, and include relevant information in the logs.

## Telemetry Data

The Gin Middleware package generates the following telemetry data:

- Request duration metric: Measures the duration of each incoming request in milliseconds.
- Request count metric: Counts the number of incoming requests.
- Request trace: Creates a trace for each incoming request, including span information.
- Detailed request information: Parses the user agent and request headers to include additional attributes in the telemetry data.

The generated telemetry data can be exported using compatible exporters, such as Prometheus, Jaeger, or OTLP, depending on your setup.

## Contributing

Contributions to the Gin Middleware package are welcome! If you find any issues or have suggestions for improvements, please open an issue or submit a pull request on the GitHub repository.

## License

The Gin Middleware package is open-source and released under the [MIT License](../LICENSE).