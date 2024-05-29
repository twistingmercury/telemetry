package metrics

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	initErrMsg = "metrics must be initialized before registering metrics"
)

var (
	apiName         string
	mPort           string
	nspace          string
	registry        *prometheus.Registry
	isInit          bool
	pubOnce         = &sync.Once{}
	server          *http.Server
	mNames          []string
	totalCalls      *prometheus.CounterVec
	concurrentCalls *prometheus.GaugeVec
	callDuration    *prometheus.HistogramVec
)

// MetricNames returns the names of the metrics associated with the Collector.
func MetricNames() []string {
	return mNames
}

// Namespace returns the Namespace for the metrics of the API.
func Namespace() string {
	return nspace
}

func ServiceName() string {
	return apiName
}

// Initialize initializes metrics system so it can TestRegisterFuncs metrics.
// This must be called before any metrics are registered.
func Initialize(port string, namespace, serviceName string) {
	if len(port) == 0 {
		panic("port for metrics must be specified")
	}
	if len(namespace) == 0 {
		panic("namespace for metrics must be specified")
	}
	if len(serviceName) == 0 {
		panic("serviceName for metrics must be specified")
	}

	p, err := strconv.Atoi(port)
	if err != nil || p < 1024 || p > 49151 {
		panic(fmt.Sprintf("invalid port value: `%s`; a valid port is a number between 1024 and 49151", port))
	}

	mPort = port
	nspace = namespace
	apiName = serviceName

	newApiMetrics()

	isInit = true
}

func MetricApiLabels() []string {
	return []string{"path", "http_method", "status_code"}
}

// newApiMetrics creates a new metrics object 'p' is the path of the API and 'm' is the HTTP method of the API.
func newApiMetrics() {
	registry = prometheus.NewRegistry()

	concurentCallsName := normalize(fmt.Sprintf("%s_concurrent_calls", ServiceName()))
	concurrentCalls = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: Namespace(),
		Name:      concurentCallsName,
		Help:      "the count of concurrent calls to the APIs, grouped by API name, path, and response code"},
		[]string{"path", "http_method"})

	totalCallsName := normalize(fmt.Sprintf("%s_total_calls", ServiceName()))
	totalCalls = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: Namespace(),
		Name:      totalCallsName,
		Help:      "The count of all call to the API, grouped by API name, path, and response code"},
		MetricApiLabels())

	callDurationName := normalize(fmt.Sprintf("%s_call_duration", ServiceName()))
	callDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: Namespace(),
		Name:      callDurationName,
		Help:      "The duration in milliseconds calls to the API, grouped by API name, path, and response code",
		Buckets:   prometheus.ExponentialBuckets(0.1, 1.5, 5)},
		MetricApiLabels())

	mNames = []string{concurentCallsName, totalCallsName, callDurationName}

	registry.MustRegister(concurrentCalls, totalCalls, callDuration)
	log.Debug().Msg("newApiMetrics invoked")
}

// Publish exposes the metrics for scraping.
func Publish() {
	pubOnce.Do(func() {
		if !isInit {
			panic(initErrMsg)
		}
		go func() {
			gin.SetMode(gin.ReleaseMode)
			router := gin.New()
			router.Use(gin.Recovery())
			promHandler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
			router.GET("/metrics", gin.WrapH(promHandler))
			router.GET("/metrics/names", func(c *gin.Context) {
				c.JSON(http.StatusOK, mNames)
			})
			server = &http.Server{
				Addr:    fmt.Sprintf(":%s", mPort),
				Handler: router.Handler(),
			}

			if err := server.ListenAndServe(); err != nil {
				log.Error().Err(err).Msg("metrics endpoint failed with error")
			}
		}()
		log.Info().Msg("metrics endpoint started")
	})
}

// RegisterCustomMetrics allows one to add a custom metric to the registry. This will panic if Initialize has not
// been called first. This is useful for adding metrics that are not API related. You can add Gauge, Counter, and
// Histogram metrics that you have defined.
func RegisterCustomMetrics(cMetrics ...prometheus.Collector) {
	if !isInit {
		panic(initErrMsg)
	}
	registry.MustRegister(cMetrics...)
}

func normalize(name string) string {
	r := regexp.MustCompile(`\s+`)
	name = r.ReplaceAllString(name, "_")
	r = regexp.MustCompile(`[./:_-]`)
	return r.ReplaceAllString(strings.ToLower(name), "_")
}
