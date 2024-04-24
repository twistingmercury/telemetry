package tracing

import (
	"context"
	"github.com/pkg/errors"
	"github.com/twistingmercury/telemetry/attributes"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
	"time"
)

const IntervalDefault = 5000 * time.Millisecond

var (
	tracer      oteltrace.Tracer
	propagator  propagation.TextMapPropagator
	commonAttrs []attribute.KeyValue
)

func Tracer() oteltrace.Tracer {
	return tracer
}

// Initialize initializes the OpenTelemetry tracing.
func Initialize(exporter sdktrace.SpanExporter, sampleRate float64, attribs attributes.Attributes) (err error) {
	if exporter == nil {
		return errors.New("trace exporter is required")
	}

	if attribs == nil {
		return errors.New("attributes are required")
	}

	res, err := resource.New(context.Background(), resource.WithAttributes(attribs.All()...))
	if err != nil {
		return
	}

	batchDuration := attribs.BatchingDuration()
	if batchDuration == 0 {
		batchDuration = IntervalDefault
	}

	bsp := sdktrace.NewBatchSpanProcessor(exporter, sdktrace.WithBatchTimeout(batchDuration))
	traceProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(sampleRate)),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)

	commonAttrs = attribs.All()

	otel.SetTracerProvider(traceProvider)
	propagator = propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	otel.SetTextMapPropagator(propagator)

	tp := otel.GetTracerProvider()

	tracer = tp.Tracer(attribs.ServiceName(), oteltrace.WithInstrumentationVersion(attribs.ServiceVersion()))

	return
}

func ExtractContext(ctx context.Context, carrier propagation.TextMapCarrier) context.Context {
	return propagator.Extract(ctx, carrier)
}

func StartSpan(ctx context.Context, name string, kind oteltrace.SpanKind, attribs ...attribute.KeyValue) (spanCtx context.Context, span oteltrace.Span) {
	commonAttrs = append(commonAttrs, attribs...)

	spanCtx, span = tracer.Start(
		ctx,
		name,
		oteltrace.WithSpanKind(kind),
		oteltrace.WithAttributes(commonAttrs...))
	return
}
