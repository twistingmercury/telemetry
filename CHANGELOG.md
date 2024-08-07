# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.0.1] - 2024-07-10

### Fixed
- Had to update module path to `../v2` since it was a major verion update.

### Changed
- Retracted v2.0.0

## [2.0.0] - 2024-07-10
### Updated
- Updated logging functions to accept `context.Context` which is used for tracing data extraction.
- Marked logging `..WithContext` functions as deprecated.

### Breaking Changes
- The function `logging.Debug` now accepts the `context.Context`
- The function `logging.Info` now accepts the `context.Context`
- The function `logging.Warn` now accepts the `context.Context`
- The function `logging.Error` now accepts the `context.Context`
- The function `logging.Fatal` now accepts the `context.Context`
- The function `logging.Panic` now accepts the `context.Context`

## [1.0.3] - 2024-06-10
### Added
- Added a new function `metrics.Shutdown` to explicitly shutdown the metrics http server.
- Updated `metrics.Publish` to use an [http.Server] instead of [gin.Engine] `Run`.

## [1.0.2] - 2024-06-03
### Added
- minor test fixes

## [1.0.1] - 2024-05-30
### Added
- Added a new function `metrics.Registry` to expose the `prometheus.Registry` for use directly by the `github.com/twistingmercury/telemetry/middleware` package.

### Updated
- Updated examples to make them more comprehensive.
- Clarified examples and instructions in README files.

## [1.0.0] - 2024-05-30

### Breaking Changes
- Changed to using Prometheus as means of creating metrics.
- The function `logging.Initialize` now accepts the args `serviceName, serviceVersion, environment string` instead of type `attributes.Attributes`.
- The function `metrics.Initialize` no longer accepts a `port` argument. Use `metrics.InitializeWithPort` to specify the port.
- The function `metrics.Initialize` now accepts the args `namespace, serviceName` instead of type `attributes.Attributes`.
- The function `tracing.Initialize` now accepts the args `serviceName, serviceVersion, environment string` instead of type `attributes.Attributes`.
- The function `tracing.Initialize` No longer accepts a `sampleRate` argumenmt. Use `tracing.InitializeWithSampleRate` to specify the tracing sample rate.
- The function `tracing.StartSpan` has been removed. Use `tracing.Start` instead.

### Added
- Added a new function `metrics.InitializeWithPort` to accept a port value for publishing metrics.
- Added a new function `tracing.InitializeWithSampleRate` to accept a float64 value for setting the trace sample rate.
- Added a new function `tracing.Start` to replace the deprecated `StartSpan`.

### Changed
- Updated `metrics.InitializeWithSampleRate` to return an error instead of panicking.
- Updated `metrics.Initialize` to default to port 9090.
- Updated `tracing.Initialize` to accept the args `serviceName, serviceVersion, environ string` in place of the 

### Removed
- Removed the Middleware code base to a separate repository.
- removed `telemetry.attributes` package.

## [0.9.2] - 2024-04-24

### Fixed
- Resolved an issue in the middleware when a 2xx status code other than 200 is treated as an error for a trace span.

## [0.9.1] - 2024-04-24

### Added
- Versioned middleware

## [0.9.0] - 2024-04-24

### Added
- Initial release of the project.
- Logging based on zerolog
- OTel Distributed Tracing
- OTel metrics

[2.0.1]: https://github.com/twistingmercury/telemetry/compare/v2.0.0...v2.0.1
[2.0.0]: https://github.com/twistingmercury/telemetry/compare/v1.0.3...v2.0.0
[1.0.3]: https://github.com/twistingmercury/telemetry/compare/v1.0.2...v1.0.3
[1.0.2]: https://github.com/twistingmercury/telemetry/compare/v1.0.1...v1.0.2
[1.0.1]: https://github.com/twistingmercury/telemetry/compare/v1.0.0...v1.0.1
[1.0.0]: https://github.com/twistingmercury/telemetry/compare/v0.9.2...v1.0.0
[0.9.2]: https://github.com/twistingmercury/telemetry/compare/v0.9.1...v0.9.2
[0.9.1]: https://github.com/twistingmercury/telemetry/compare/v0.9.0...v0.9.1
[0.9.0]: https://github.com/twistingmercury/telemetry/releases/tag/v0.9.0
