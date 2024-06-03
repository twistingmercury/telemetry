package data

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/twistingmercury/telemetry/metrics"
)

var (
	tCtr  *prometheus.CounterVec
	dHist *prometheus.HistogramVec
)

const (
	apiLabel = "api"
	pckLabel = "pkg"
	fncLabel = "func"
	errLabel = "isError"
)

// Metrics returns the metrics that are defined for the data package.
func Metrics() (c []prometheus.Collector) {
	labels := []string{apiLabel, pckLabel, fncLabel, errLabel}
	tCtr = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: metrics.Namespace(),
		Name:      "data_total_calls",
		Help:      "The total count of calls to the funcs in the data package"},
		labels)

	dHist = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: metrics.Namespace(),
		Name:      "data_call_duration",
		Help:      "Duration each func call within the data package",
		Buckets:   prometheus.ExponentialBuckets(0.1, 1.5, 5),
	}, labels)

	c = []prometheus.Collector{tCtr, dHist}
	return
}

// DoDatabaseStuff simulates a database call.
func DoDatabaseStuff() (err error) {
	s := time.Now()
	defer func() {
		duration := float64(time.Since(s))
		incMetrics("DoDatabaseStuff", duration, err)
	}()

	src := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(src)

	minSleep := 10
	maxSleep := 100

	// simulate some random latency
	time.Sleep(time.Duration(rnd.Intn(maxSleep-minSleep)+minSleep) * time.Millisecond)

	// simulate a random error...
	if rnd.Intn(24)%7 == 0 {
		err = fmt.Errorf("random simulated error")
		return
	}

	return
}

// incMetrics is a helper function to increment the metrics for the data package.
func incMetrics(fName string, d float64, err error) {
	isErr := err != nil
	tCtr.WithLabelValues("examples", "data", fName, strconv.FormatBool(isErr)).Inc()
	dHist.WithLabelValues("examples", "data", fName, strconv.FormatBool(isErr)).Observe(d)
}
