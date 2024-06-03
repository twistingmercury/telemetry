package metrics

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
)

const defaultPort = "9090"

var (
	apiName  string
	mPort    string
	nspace   string
	registry *prometheus.Registry
	server   *http.Server
)

// Namespace returns the Namespace for the metrics of the API.
func Namespace() string {
	return nspace
}

// ServiceName returns the name of the service used to create the metrics for the API.
func ServiceName() string {
	return apiName
}

// Port returns the port being used to publish metrics. The default is 9090.
func Port() string {
	return mPort
}

// Registry returns the internal [prometheus.Registry] so it can be used directly if required.
func Registry() *prometheus.Registry {
	return registry
}

// Initialize initializes metrics system on the default port 9090.
func Initialize(namespace, serviceName string) error {
	return InitializeWithPort(defaultPort, namespace, serviceName)
}

// InitializeWithPort initializes metrics with a specific port to publish metrics on.
// This must be called before any metrics are registered.
func InitializeWithPort(port string, namespace, serviceName string) error {
	if len(port) == 0 {
		return errors.New("port for metrics must be specified")
	}
	if len(namespace) == 0 {
		return errors.New("namespace for metrics must be specified")
	}
	if len(serviceName) == 0 {
		return errors.New("serviceName for metrics must be specified")
	}

	p, err := strconv.Atoi(port)
	if err != nil || p < 1024 || p > 49151 {
		return errors.New(fmt.Sprintf("invalid port value: `%s`; a valid port is a number between 1024 and 49151", port))
	}

	registry = prometheus.NewRegistry()

	mPort = port
	nspace = namespace
	apiName = serviceName
	return nil
}

// Publish exposes the metrics for scraping.
func Publish() {
	go func() {
		gin.SetMode(gin.ReleaseMode)
		router := gin.New()
		router.Use(gin.Recovery())
		promHandler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
		router.GET("/metrics", gin.WrapH(promHandler))
		server = &http.Server{
			Addr:    fmt.Sprintf(":%s", mPort),
			Handler: router.Handler(),
		}

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Panic().Err(err).Msg("metrics endpoint failed with error")
		}
	}()
	log.Info().Msg("metrics endpoint started")
}

// RegisterMetrics is used to add one to or more metrics (collectors) to the registry.
func RegisterMetrics(cMetrics ...prometheus.Collector) {
	registry.MustRegister(cMetrics...)
}
