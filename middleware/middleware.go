package middleware

import (
	"errors"
	"fmt"
	gonic "github.com/gin-gonic/gin"
	"github.com/mileusna/useragent"
	"github.com/twistingmercury/telemetry/attributes"
	"github.com/twistingmercury/telemetry/logging"
	"github.com/twistingmercury/telemetry/metrics"
	"github.com/twistingmercury/telemetry/tracing"
	"go.opentelemetry.io/otel/attribute"
	otelCodes "go.opentelemetry.io/otel/codes"
	otelmetric "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	oteltrace "go.opentelemetry.io/otel/trace"
	"strings"
	"sync"
	"time"
)

const (
	method = "method"
	status = "status"
	path   = "path"
)

var (
	meterDuration otelmetric.Float64Histogram
	meterCounter  otelmetric.Int64Counter
	once          sync.Once
	tAttribs      attributes.Attributes
)

// Telemetry returns middleware that will instrument and trace incoming requests.
func Telemetry(attribs attributes.Attributes) gonic.HandlerFunc {
	once.Do(func() {
		tAttribs = attribs
		var err error
		meterDuration, err = metrics.Meter().Float64Histogram(
			fmt.Sprintf("%s.%s.request.duration", tAttribs.Namespace(), tAttribs.ServiceName()),
			otelmetric.WithDescription("Measures the duration of inbound call."),
			otelmetric.WithUnit("ms"))
		if err != nil {
			panic(err)
		}

		meterCounter, err = metrics.Meter().Int64Counter(
			fmt.Sprintf("%s.%s.request.count", tAttribs.Namespace(), tAttribs.ServiceName()),
			otelmetric.WithDescription("Measures the number of inbound calls."),
			otelmetric.WithUnit("{count}"))
		if err != nil {
			panic(err)
		}
	})
	return func(c *gonic.Context) {
		savedCtx := c.Request.Context()
		defer func() {
			c.Request = c.Request.WithContext(savedCtx)
		}()

		spanName := fmt.Sprintf("%s: %s", c.Request.Method, c.Request.URL.Path)
		parentCtx := tracing.ExtractContext(savedCtx, propagation.HeaderCarrier(c.Request.Header))
		opts := attribs.All()
		opts = append(opts, semconv.HTTPRoute(spanName))
		spanctx, span := tracing.StartSpan(parentCtx, spanName, oteltrace.SpanKindServer, opts...)
		defer span.End()

		before := time.Now()
		c.Next()
		elapsedTime := float64(time.Since(before)) / float64(time.Millisecond)

		logRequest(span.SpanContext(), c, elapsedTime)
		code, desc := SpanStatus(c.Writer.Status())
		span.SetStatus(code, desc)

		metricAttrs := []attribute.KeyValue{
			attribute.String(method, c.Request.Method),
			attribute.String(path, c.Request.URL.Path),
			attribute.Int(status, c.Writer.Status())}

		meterCounter.Add(spanctx, 1, otelmetric.WithAttributes(metricAttrs...))
		meterDuration.Record(spanctx, elapsedTime, otelmetric.WithAttributes(metricAttrs...))
	}
}

// SpanStatus returns the OpenTelemetry status code as defined in
// go.opentelemetry.io/old_elemetry/codes and a brief description for a given HTTP status code.
func SpanStatus(status int) (code otelCodes.Code, desc string) {
	switch status {
	case 200:
		code = otelCodes.Ok
		desc = "OK"
	case 400:
		code = otelCodes.Ok
		desc = "Bad Request"
	case 401:
		code = otelCodes.Ok
		desc = "Unauthorized"
	case 403:
		code = otelCodes.Ok
		desc = "Forbidden"
	case 404:
		code = otelCodes.Ok
		desc = "Not Found"
	case 405:
		code = otelCodes.Ok
		desc = "Method Not Allowed"
	case 500:
		code = otelCodes.Error
		desc = "Internal Server Error"
	case 502:
		code = otelCodes.Error
		desc = "Bad Gateway"
	case 503:
		code = otelCodes.Error
		desc = "Service Unavailable"
	default:
		code = otelCodes.Error
		desc = "Unknown"
	}
	return
}

const (
	Http            = "http"
	Https           = "https"
	HttpMethod      = "http.request.method"
	HttpPath        = "http.request.path"
	HttpRemoteAddr  = "http.request.remoteAddr"
	HttpRequestHost = "http.request.host"
	HttpStatus      = "http.response.status"
	HttpLatency     = "http.response.latency"
	TLSVersion      = "http.tls.version"
	HttpScheme      = "http.scheme"

	//QueryString = "http.request.queryString"
)

func logRequest(spanCtx oteltrace.SpanContext, c *gonic.Context, elapsedTime float64) {
	defer func() {
		if r := recover(); r != nil {
			logging.Error(errors.New("panic in logging middleware"),
				"panic in logging middleware", logging.KeyValue{Key: "panic", Value: r})
		}
	}()

	status := c.Writer.Status()
	args := map[string]any{
		HttpMethod:          c.Request.Method,
		HttpPath:            c.Request.URL.Path,
		HttpRemoteAddr:      c.Request.RemoteAddr,
		HttpStatus:          status,
		HttpLatency:         fmt.Sprintf("%fms", elapsedTime),
		logging.TraceIDAttr: spanCtx.TraceID().String(),
		logging.SpanIDAttr:  spanCtx.SpanID().String(),
	}

	scheme := Http
	if c.Request.TLS != nil {
		scheme = Https
		args[TLSVersion] = c.Request.TLS.Version
	}

	args[HttpScheme] = scheme
	args[HttpRequestHost] = c.Request.Host

	/* !!! this could log sensitive data. leaving out for now. !!!
	if rQuery := c.Request.URL.RawQuery; len(rQuery) > 0 {
		args[QueryString] = rQuery
	}
	*/

	hd := ParseHeaders(c.Request.Header)
	args = logging.MergeMaps(args, hd)
	ua := ParseUserAgent(c.Request.UserAgent())
	args = logging.MergeMaps(args, ua)

	logAttribs := logging.FromMap(args)
	if status > 499 || c.Errors.Last() != nil {
		errs := strings.Join(c.Errors.Errors(), ";")
		logging.Error(errors.New(errs), "request failed", logAttribs...)
		return
	}

	logging.Info("request ", logAttribs...)
}

// ParseHeaders parses the headers and returns a map of attribs.
func ParseHeaders(headers map[string][]string) (args map[string]any) {
	args = make(map[string]any)
	for k, v := range headers {
		args[strings.ToLower("http."+k)] = strings.ToLower(strings.Join(v, ", "))
	}
	return
}

const (
	UserAgentOS             = "http.user_agent.os"
	UserAgentOSVersion      = "http.user_agent.os_version"
	UserAgentDevice         = "http.user_agent.device"
	UserAgentBrowser        = "http.user_agent.browser"
	UserAgentBrowserVersion = "http.user_agent.browser_version"
	BrowserChrome           = "chrome"
	BrowserSafari           = "safari"
	BrowserFirefox          = "firefox"
	BrowserOpera            = "opera"
	BrowserIE               = "ie"
	BrowserEdge             = "edge"
	BrowserTrident          = "Trident"
	DeviceMobile            = "mobile"
	DeviceDesktop           = "desktop"
	DeviceBot               = "bot"
)

// ParseUserAgent parses the user agent string and returns a map of attribs.
func ParseUserAgent(rawUserAgent string) (args map[string]any) {
	if len(rawUserAgent) == 0 {
		return //no-op
	}

	args = make(map[string]any)
	ua := useragent.Parse(rawUserAgent)

	args[UserAgentOS] = ua.OS
	args[UserAgentOSVersion] = ua.OSVersion

	var device string
	switch {
	case ua.Mobile || ua.Tablet:
		device = DeviceMobile
	case ua.Desktop:
		device = DeviceDesktop
	case ua.Bot:
		device = DeviceBot
	}

	args[UserAgentDevice] = device

	var browser string
	if ua.Mobile || ua.Tablet || ua.Desktop {
		switch {
		case ua.IsChrome():
			browser = BrowserChrome
		case ua.IsSafari():
			browser = BrowserSafari
		case ua.IsFirefox():
			browser = BrowserFirefox
		case ua.IsOpera():
			browser = BrowserOpera
		case ua.IsInternetExplorer() || strings.Contains(rawUserAgent, BrowserTrident):
			browser = BrowserIE
		case ua.IsEdge():
			browser = BrowserEdge
		}

		args[UserAgentBrowser] = browser
		args[UserAgentBrowserVersion] = ua.Version
	}
	return
}
