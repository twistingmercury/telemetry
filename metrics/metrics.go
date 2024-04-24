package metrics

import (
	"context"
	"github.com/pkg/errors"
	"github.com/twistingmercury/telemetry/attributes"
	"go.opentelemetry.io/otel"
	otelmetric "go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"time"
)

const IntervalDefault = 60 * time.Second

var (
	meter         otelmetric.Meter
	meterProvider otelmetric.MeterProvider
)

func Meter() otelmetric.Meter {
	return meter
}

func MeterProvider() otelmetric.MeterProvider {
	return meterProvider
}

// Initialize initializes the OpenTelemetry middleware.
func Initialize(exporter sdkmetric.Exporter, attribs attributes.Attributes) (err error) {
	if exporter == nil {
		return errors.Errorf("metric exporter is required")
	}

	if attribs == nil {
		return errors.New("attributes are required")
	}

	res, err := resource.New(context.Background(), resource.WithAttributes(attribs.All()...))
	if err != nil {
		return
	}

	res, err = resource.New(context.Background(), resource.WithAttributes(attribs.All()...))
	if err != nil {
		return
	}

	batchDuration := attribs.BatchingDuration()
	if batchDuration == 0 {
		batchDuration = IntervalDefault
	}

	reader := sdkmetric.NewPeriodicReader(exporter, sdkmetric.WithInterval(batchDuration))
	meterProvider = sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader), sdkmetric.WithResource(res))
	otel.SetMeterProvider(meterProvider)
	mp := otel.GetMeterProvider()
	meter = mp.Meter(attribs.ServiceName(), otelmetric.WithInstrumentationVersion(attribs.ServiceVersion()))

	return
}
