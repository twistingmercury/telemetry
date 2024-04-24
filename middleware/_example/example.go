package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/twistingmercury/telemetry/attributes"
	"github.com/twistingmercury/telemetry/logging"
	"github.com/twistingmercury/telemetry/metrics"
	gmw "github.com/twistingmercury/telemetry/middleware/gin"
	"github.com/twistingmercury/telemetry/tracing"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"log"
	"net/http"
	"os"
)

func main() {
	// Create an instance of attributes
	// In production you would want to use environment variables, or some other configuration mechanism.
	attribs := attributes.New("my_namespace", "my_service", "1.0.0", "production")

	// Initialize logging so that we can see what's happening
	err := logging.Initialize(zerolog.DebugLevel, attribs, os.Stdout)
	if err != nil {
		log.Panicf("failed to initialize logging: %v", err)
	}

	// Create metrics exporter
	// In production you would want to use a gRPC exporter
	metricsExporter, err := stdoutmetric.New(stdoutmetric.WithPrettyPrint())
	if err != nil {
		logging.Fatal(err, "failed to create metrics exporter")
	}

	// Create trace exporter
	// In production you would want to use a gRPC exporter
	traceExporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		logging.Fatal(err, "failed to create trace exporter")
	}

	// Initialize metrics and tracing
	err = metrics.Initialize(metricsExporter, attribs)
	if err != nil {
		logging.Fatal(err, "failed to initialize metrics")
	}

	err = tracing.Initialize(traceExporter, 1.0, attribs)
	if err != nil {
		logging.Fatal(err, "failed to initialize tracing")
	}

	// use gin.ReleaseMode in production!!!
	gin.SetMode(gin.DebugMode)

	// Create a new Gin router
	router := gin.New()

	// Register the TelemetryMiddleware
	router.Use(gin.Recovery(), gmw.TelemetryMiddleware(attribs))

	// Define a simple route
	router.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, World!")
	})

	// Start the server
	log.Println("Server running on http://localhost:8080")
	if err := router.Run(":8080"); err != nil {
		logging.Fatal(err, "failed to start server")
	}
}
