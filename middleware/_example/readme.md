# Example: Gin Middleware

This example demonstrates how to use the Gin Middleware package to instrument and trace incoming requests in a Gin web application. The middleware integrates with the OpenTelemetry framework through the [github.com/twistingmercury/telemetry](https://github.com/twistingmercury/telemetry/blob/main/readme.md) to collect and export telemetry data, including metrics and traces.

## Prerequisites

Before running the example, ensure that you have the following prerequisites installed:

- Go programming language (version 1.21 or higher)
- Gin web framework (`go get -u github.com/gin-gonic/gin`)
- Telemetry packages:
    - `go get github.com/twistingmercury/telemetry/`
    - `go get github.com/twistingmercury/telemetry/middleware/gin`
- OpenTelemetry dependencies:
    - `go get go.opentelemetry.io/otel`
    - `go get go.opentelemetry.io/otel/exporters/stdout/stdouttrace`
    - `go get go.opentelemetry.io/otel/exporters/stdout/stdoutmetric`
    - `go get go.opentelemetry.io/otel/sdk/metric`
    - `go get go.opentelemetry.io/otel/sdk/resource`

## Usage

From a terminal, follow these steps to run the example:

1. Run the example:

   ```
   go run example.go
   ```

   The server will start running on `http://localhost:8080`.

2. Access the `/hello` route in your web browser or using a tool like cURL:

   ```
   curl http://localhost:8080/hello
   ```

   You should see the response "Hello, World!".

3. Observe the console output to see the collected metrics and traces.

## Configuration

The example uses default configurations for the following:

- Namespace: "my_namespace"
- Service name: "my_service"
- Service version: "1.0.0"
- Environment: "production"

You can modify these values by updating the `attribs` variable in the `main()` function:

```go
attribs := attributes.New("my_namespace", "my_service", "1.0.0", "production")
```

## Exporters

The example uses stdout exporters for metrics and traces, which print the collected telemetry data to the console. You can replace these exporters with other compatible exporters, such as Prometheus, Jaeger, or OTLP, by modifying the initialization code in the `main()` function.

