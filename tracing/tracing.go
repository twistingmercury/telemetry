package tracing

import (
	"context"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

const (
	SampleRateDefault = 1.0
)

var (
	tracer      oteltrace.Tracer
	propagator  propagation.TextMapPropagator
	commonAttrs []attribute.KeyValue
	svcName     string
	svcVersion  string
	env         string
)

func Tracer() oteltrace.Tracer {
	return tracer
}

// Initialize nitializes the OpenTelemetry tracing with the provided values.
func Initialize(exporter sdktrace.SpanExporter, serviceName, serviceVersion, environment string) error {
	return InitializeWithSampleRate(
		exporter,
		SampleRateDefault,
		serviceName,
		serviceVersion,
		environment)
}

// InitializeWithSampleRate initializes the OpenTelemetry tracing, and sets the sample rate to the value passed by the sampleRate arg.
func InitializeWithSampleRate(exporter sdktrace.SpanExporter, sampleRate float64, serviceName, serviceVersion, environment string) (err error) {
	if exporter == nil {
		return errors.New("trace exporter is required")
	}

	if sampleRate == 0 {
		return errors.New("sample-rate must be a floating point value between 0.1 and 1.0")
	}

	svcName = serviceName
	svcVersion = serviceVersion
	env = environment

	commonAttrs = []attribute.KeyValue{
		semconv.ServiceNameKey.String(svcName),
		semconv.ServiceVersionKey.String(svcVersion),
		semconv.DeploymentEnvironmentKey.String(env),
	}

	res, err := resource.New(
		context.Background(),
		resource.WithAttributes(commonAttrs...))
	if err != nil {
		return
	}

	bsp := sdktrace.NewBatchSpanProcessor(exporter)
	traceProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(sampleRate)),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)

	otel.SetTracerProvider(traceProvider)
	propagator = propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	otel.SetTextMapPropagator(propagator)

	tp := otel.GetTracerProvider()

	tracer = tp.Tracer(serviceName, oteltrace.WithInstrumentationVersion(serviceVersion))

	return
}

// ExtractContext returns the [context.Context] for OTel tracing that may be passed
func ExtractContext(ctx context.Context, carrier propagation.TextMapCarrier) context.Context {
	return propagator.Extract(ctx, carrier)
}

func Start(ctx context.Context, name string, kind oteltrace.SpanKind, attribs ...attribute.KeyValue) (spanCtx context.Context, span oteltrace.Span) {
	commonAttrs = append(commonAttrs, attribs...)

	spanCtx, span = Tracer().Start(
		ctx,
		name,
		oteltrace.WithSpanKind(kind),
		oteltrace.WithAttributes(commonAttrs...))
	return
}
