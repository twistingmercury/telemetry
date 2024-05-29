package metrics_test

import (
	"github.com/twistingmercury/telemetry/metrics"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func TestInitializePanics(t *testing.T) {
	defer metrics.Reset()
	assert.Panics(t, func() { metrics.Initialize("", "unit", "test") })

	metrics.Reset()
	assert.Panics(t, func() { metrics.Initialize("1023", "unit", "test") })

	metrics.Reset()
	assert.Panics(t, func() { metrics.Initialize("49152", "unit", "test") })

	metrics.Reset()
	assert.Panics(t, func() { metrics.Initialize("1234", "", "test") })

	metrics.Reset()
	assert.Panics(t, func() { metrics.Initialize("1234", "unit", "") })

	metrics.Reset()
	assert.Panics(t, func() { metrics.RegisterCustomMetrics() })

	metrics.Reset()
	assert.Panics(t, func() { metrics.Publish() })
}

func TestInitalize(t *testing.T) {
	defer metrics.Reset()
	metrics.Initialize("1024", "unit", "test")
	assert.Equal(t, []string{"path", "http_method", "status_code"}, metrics.MetricApiLabels())
}

func TestPublish(t *testing.T) {
	defer metrics.Reset()
	metrics.Initialize("1024", "unit", "test")
	assert.NotPanics(t, func() { metrics.Publish() })
}

func TestGinMiddlewareNames(t *testing.T) {
	defer metrics.Reset()
	expected := []string{
		"test_concurrent_calls",
		"test_total_calls",
		"test_call_duration"}
	metrics.Initialize("1024", "unit", "test")
	metrics.Publish()
	assert.Equal(t, expected, metrics.MetricNames())
}

func TestRegisterCustomMetrics(t *testing.T) {
	defer metrics.Reset()
	metrics.Initialize("1024", "unit", "test")
	metrics.Publish()
	metrics.RegisterCustomMetrics(customMetrics()...)
}

// Metrics returns a slice of prometheus.Collector that can be registered
func customMetrics() (c []prometheus.Collector) {
	labels := []string{"function", "cmd_type", "cmd", "result"}
	ctr := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "custom",
		Name:      `counter_1`,
		Help:      "The total count of calls to the func data.DoBusinessLogicStuff",
	}, labels)

	dur := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "custom",
		Name:      `duration_milliseconds`,
		Help:      "Duration the func data.DoDatabaseStuff took to execute successfully",
		Buckets:   prometheus.ExponentialBuckets(0.1, 1.5, 5),
	}, labels)

	c = []prometheus.Collector{ctr, dur}
	return
}
