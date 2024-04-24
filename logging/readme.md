# Logging Package

The Logging package is a Go library that provides a simple and flexible logging system built on top of the Zerolog library. It allows you to log messages with different severity levels, add custom key-value pairs to log messages, and integrate with OpenTelemetry tracing to include trace and span information in log messages.

## Features

- Six log levels: debug, info, warn, error, fatal, and panic
- Structured logging with custom key-value pairs
- Integration with OpenTelemetry tracing to include trace and span information
- Customizable output writer
- Initialization with common attribs for consistent logging across the application

## Installation

To install the Logging package, use the following command:

```
go get github.com/twistingmercury/telemetry
```

## Usage

### Initialization

Before using the Logging package, you need to initialize it with the desired log level, common attribs, and output writer:

```go
import "github.com/twistingmercury/telemetry/logging"

err := logging.Initialize(zerolog.DebugLevel, attribs, os.Stdout)
if err != nil {
    // Handle initialization error
}
```

- `zerolog.DebugLevel` specifies the minimum log level to be logged. You can use `zerolog.InfoLevel`, `zerolog.WarnLevel`, `zerolog.ErrorLevel`, `zerolog.FatalLevel`, or `zerolog.PanicLevel` based on your needs.
- `attribs` is an instance of `common.Attributes` that contains common attribs to be included in every log message.
- `os.Stdout` is the output writer where the log messages will be written. You can use any `io.Writer` implementation, such as a file or a network connection.

### Logging Messages

The Logging package provides functions to log messages with different severity levels:

```go
logging.Debug("Debug message", logging.KeyValue{Key: "key", Value: "value"})
logging.Info("Info message", logging.KeyValue{Key: "key", Value: 123})
logging.Warn("Warn message")
logging.Error(err, "Error message")
logging.Fatal(err, "Fatal message")
logging.Panic(err, "Panic message")
```

Each function takes a message string and optional `logging.KeyValue` pairs to add custom key-value pairs to the log message.

### Logging with Context

If you have a span context from OpenTelemetry tracing, you can include it in the log messages using the `WithContext` variants of the logging functions:

```go
// _, span := tracing.StartSpan(...)
// spanCtx := span.SpanContext()

logging.DebugWithContext(&spanCtx, "debug message")
logging.InfoWithContext(&spanCtx, "info message")
logging.WarnWithContext(&spanCtx, "warn message")
logging.ErrorWithContext(&spanCtx, err, "error message")
logging.FatalWithContext(&spanCtx, err, "fatal message")
logging.PanicWithContext(&spanCtx, err, "panic message")
```

The `spanCtx` parameter is a pointer to `trace.SpanContext` that contains the trace and span information to be included in the log message.

### Custom Key-Value Pairs

You can add custom key-value pairs to log messages using the `logging.KeyValue` struct:

```go
logging.Info(
    "Info message", 
    logging.KeyValue{Key: "key1", Value: "value1"}, 
    logging.KeyValue{Key: "key2", Value: 123},
)
```

```go
logging.InfoWithContext(
    &spanCtx, 
    "Info message", 
    logging.KeyValue{Key: "key1", Value: "value1"}, 
    logging.KeyValue{Key: "key2", Value: 123},
)
```

The `Key` field is a string representing the key, and the `Value` field can be of any type that Zerolog supports.

## Contributing

Contributions to the Logging package are welcome! If you find any issues or have suggestions for improvements, please open an issue or submit a pull request on the GitHub repository.

## License

The Logging package is open-source and released under the [MIT License](../LICENSE).