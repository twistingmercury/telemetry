[![CodeQL](https://github.com/twistingmercury/telemetry/actions/workflows/codeql.yml/badge.svg)](https://github.com/twistingmercury/telemetry/actions/workflows/codeql.yml)
[![Run Tests](https://github.com/twistingmercury/telemetry/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/twistingmercury/telemetry/actions/workflows/go.yml)
[![codecov](https://codecov.io/github/twistingmercury/telemetry/graph/badge.svg?token=U6C4TE88OP)](https://codecov.io/github/twistingmercury/telemetry)
# Telemetry Package

This package was created to help me reduce a bunch of repetitive tasks in creating a Go application. All of the apps and services need logging, distributed tracing, and metrics.

## Features

- Logging: The package provides a logging system built on top of the [zerolog](https://pkg.go.dev/github.com/rs/zerolog)   library. It supports various log levels (debug, info, warn, error, fatal, panic) and allows adding custom key-value pairs to log messages. The logging system also integrates with OpenTelemetry tracing to include trace and span information in log messages.

- Metrics: The package utilizes [Prometheus](https://pkg.go.dev/github.com/prometheus/client_golang/prometheus) to collect metrics. It provides a simple way to initialize a metrics collector, which can be used to record metrics throughout the application.

- Tracing: The package integrates with OpenTelemetry tracing to collect and export trace data. It provides functions to initialize a tracer provider, extract trace context from incoming requests, and start new spans for outgoing requests or internal operations. The package allows configuring the batching duration for the tracing batch processor.

## Installation

To install the Telemetry package, use the following command:

```
go get github.com/twistingmercury/telemetry
```

## Usage and examples

A make file exists in the [_example](./_example/Makefile) directory where by you can run the examples.

### Logging
* [Details](./logging/README.md)
* [Example](./_example/logging/main.go)

### Metrics
* [Details](./metrics/README.md)
* [Example](./_example/metrics/main.go)

### Tracing
* [Details](./tracing/README.md)
* [Example](./_example/metrics/main.go)

A complete example of using all three at once can be found here: [Complete Example](./_example/complete/main.go)

## Contributing


Contributions to the Telemetry package are welcome! If you find any issues or have suggestions for improvements, please open an issue or submit a pull request on the GitHub repository.

## License

The Telemetry package is open-source and released under the [MIT License](LICENSE).