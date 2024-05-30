package metrics_test

import (
	"github.com/stretchr/testify/require"
	"github.com/twistingmercury/telemetry/metrics"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

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

func TestInitializePanics(t *testing.T) {
	assert.Error(t, metrics.InitializeWithPort("", "unit", "test"))
	assert.Error(t, metrics.InitializeWithPort("1023", "unit", "test"))
	assert.Error(t, metrics.InitializeWithPort("49152", "unit", "test"))
	assert.Error(t, metrics.InitializeWithPort("1234", "", "test"))
	assert.Error(t, metrics.InitializeWithPort("1234", "unit", ""))
	assert.Error(t, metrics.Initialize("unit", ""))
	assert.Error(t, metrics.Initialize("", "test"))
}

func TestInitalize(t *testing.T) {
	err := metrics.InitializeWithPort("1024", "unit", "test")
	require.NoError(t, err)
	assert.Equal(t, "unit", metrics.Namespace())
	assert.Equal(t, "test", metrics.ServiceName())
	assert.Equal(t, "1024", metrics.Port())
}

func TestPublish(t *testing.T) {
	err := metrics.InitializeWithPort("1024", "unit", "test")
	require.NoError(t, err)

	assert.NotPanics(t, func() { metrics.Publish() })
}

func TestRegisterCustomMetrics(t *testing.T) {
	err := metrics.InitializeWithPort("1024", "unit", "test")
	require.NoError(t, err)
	metrics.RegisterMetrics(customMetrics()...)
}
