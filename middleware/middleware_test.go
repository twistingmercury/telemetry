package middleware_test

import (
	"bytes"
	"context"
	"encoding/json"
	gonic "github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/twistingmercury/telemetry/attributes"
	"github.com/twistingmercury/telemetry/logging"
	"github.com/twistingmercury/telemetry/metrics"
	"github.com/twistingmercury/telemetry/middleware"
	"github.com/twistingmercury/telemetry/tracing"
	"go.opentelemetry.io/otel/sdk/metric"
	oteltrace "go.opentelemetry.io/otel/trace"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
)

const (
	namespace   = "unit"
	serviceName = "test"
	version     = "1.0.0"
	environment = "test"
)

var (
	tex         *stdouttrace.Exporter
	mex         metric.Exporter
	attribs     attributes.Attributes
	logBuffer   *bytes.Buffer
	traceBuffer *bytes.Buffer
	meterBuffer *bytes.Buffer
)

func initializeTests(t *testing.T) {
	var err error
	traceBuffer = &bytes.Buffer{}
	tex, err = stdouttrace.New(stdouttrace.WithWriter(traceBuffer), stdouttrace.WithPrettyPrint())
	require.NoError(t, err)

	meterBuffer = &bytes.Buffer{}
	mex, err = stdoutmetric.New(stdoutmetric.WithEncoder(json.NewEncoder(meterBuffer)))
	require.NoError(t, err)

	attribs = attributes.New(namespace, serviceName, version, environment,
		[]attribute.KeyValue{attribute.String("foo", "bar")}...)

	err = tracing.Initialize(tex, 1.0, attribs)
	require.NoError(t, err)

	err = metrics.Initialize(mex, attribs)
	require.NoError(t, err)

	logBuffer = &bytes.Buffer{}
	err = logging.Initialize(zerolog.DebugLevel, attribs, logBuffer)
	require.NoError(t, err)
}

func resetTests() {
	tex.Shutdown(context.TODO())
	tex = nil
	mex.Shutdown(context.TODO())
	mex = nil
	logBuffer.Reset()
	traceBuffer.Reset()
	meterBuffer.Reset()
}

func TestCommonAttributes(t *testing.T) {
	attribs := attributes.New(namespace, serviceName, version, environment,
		[]attribute.KeyValue{attribute.String("foo", "bar")}...)

	require.Equal(t, 5, len(attribs.All()))
	require.Equal(t, attribs.ServiceName(), serviceName)
	require.Equal(t, attribs.ServiceVersion(), version)
	require.Equal(t, attribs.Namespace(), namespace)
	require.Equal(t, attribs.Environment(), environment)
	require.Equal(t, attribs.Get("foo"), "bar")
}

func TestGinOTelMiddleware(t *testing.T) {
	initializeTests(t)
	defer resetTests()

	gonic.SetMode(gonic.TestMode)
	r := gonic.New()
	r.Use(middleware.Telemetry(attribs))
	r.GET("/test", func(c *gonic.Context) {
		c.Status(http.StatusOK)
	})

	req, err := http.NewRequest(http.MethodGet, "/test", nil)
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	checkLog(t)
}

func checkLog(t *testing.T) {
	logText := logBuffer.String()
	require.NotEmptyf(t, logText, "no logs found")

	var logEntry map[string]any

	err := json.Unmarshal([]byte(logText), &logEntry)
	require.NoError(t, err)

	require.Contains(t, logEntry, "otel.trace_id")
	require.Len(t, logEntry["otel.trace_id"], 32)
	require.NotEqual(t, logEntry["otel.trace_id"], oteltrace.TraceID{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})

	require.Contains(t, logEntry, "otel.span_id")
	require.Len(t, logEntry["otel.span_id"], 16)
	require.NotEqual(t, logEntry["otel.trace_id"], oteltrace.SpanID{0, 0, 0, 0, 0, 0, 0, 0})
	require.Equal(t, logEntry["level"], "info")
	require.Equal(t, logEntry["service"], serviceName)
	require.Equal(t, logEntry["version"], version)
	require.Equal(t, logEntry["env"], environment)
}

func TestSpanStatus(t *testing.T) {
	testCases := []struct {
		name     string
		status   int
		expected struct {
			code codes.Code
			desc string
		}
	}{
		{
			name:   "OK",
			status: http.StatusOK,
			expected: struct {
				code codes.Code
				desc string
			}{code: codes.Ok, desc: "OK"},
		},
		{
			name:   "Bad Request",
			status: http.StatusBadRequest,
			expected: struct {
				code codes.Code
				desc string
			}{code: codes.Ok, desc: "Bad Request"},
		},
		{
			name:   "Unauthorized",
			status: http.StatusUnauthorized,
			expected: struct {
				code codes.Code
				desc string
			}{code: codes.Ok, desc: "Unauthorized"},
		},
		{
			name:   "Forbidden",
			status: http.StatusForbidden,
			expected: struct {
				code codes.Code
				desc string
			}{code: codes.Ok, desc: "Forbidden"},
		},
		{
			name:   "Not Found",
			status: http.StatusNotFound,
			expected: struct {
				code codes.Code
				desc string
			}{code: codes.Ok, desc: "Not Found"},
		},
		{
			name:   "Method Not Allowed",
			status: http.StatusMethodNotAllowed,
			expected: struct {
				code codes.Code
				desc string
			}{code: codes.Ok, desc: "Method Not Allowed"},
		},
		{
			name:   "Internal Server Error",
			status: http.StatusInternalServerError,
			expected: struct {
				code codes.Code
				desc string
			}{code: codes.Error, desc: "Internal Server Error"},
		},
		{
			name:   "Bad Gateway",
			status: http.StatusBadGateway,
			expected: struct {
				code codes.Code
				desc string
			}{code: codes.Error, desc: "Bad Gateway"},
		},
		{
			name:   "Service Unavailable",
			status: http.StatusServiceUnavailable,
			expected: struct {
				code codes.Code
				desc string
			}{code: codes.Error, desc: "Service Unavailable"},
		},
		{
			name:   "Unknown",
			status: http.StatusTeapot,
			expected: struct {
				code codes.Code
				desc string
			}{code: codes.Error, desc: "Unknown"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			code, desc := middleware.SpanStatus(tc.status)
			assert.Equal(t, tc.expected.code, code)
			assert.Equal(t, tc.expected.desc, desc)
		})
	}
}

type testUserAgent struct {
	ua      string
	uaType  string
	browser string
}

const nilValue = "<nil>"

var testUserAgents = []testUserAgent{
	{"Mozilla/5.0 (Linux; Android 7.0; SM-T827R4 Build/NRD90M) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.116 Safari/537.36", "mobile", middleware.BrowserChrome},
	{"Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)", "bot", nilValue},
	{"Mozilla/5.0 (compatible; bingbot/2.0; +http://www.bing.com/bingbot.htm)", "bot", nilValue},
	{"Mozilla/5.0 (compatible; Yahoo! Slurp; http://help.yahoo.com/help/us/ysearch/slurp)", "bot", nilValue},
	{"Mozilla/5.0 (compatible; Baiduspider/2.0; +http://www.baidu.com/search/spider.html)", "bot", nilValue},
	{"Mozilla/5.0 (compatible; YandexBot/3.0; +http://yandex.com/bots)", "bot", nilValue},
	{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.135 Safari/537.36 Edge/12.246", "desktop", middleware.BrowserEdge},
	{"Mozilla/5.0 (iPhone13,2; U; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/602.1.50 (KHTML, like Gecko) Version/10.0 Mobile/15E148 Safari/602.1", "mobile", middleware.BrowserSafari},
	{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_2) AppleWebKit/601.3.9 (KHTML, like Gecko) Version/9.0.2 Safari/601.3.9", "desktop", middleware.BrowserSafari},
	{"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:15.0) Gecko/20100101 Firefox/15.0.1", "desktop", middleware.BrowserFirefox},
	{"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36 OPR/102.0.0.0", "desktop", middleware.BrowserOpera},
	{"Mozilla/5.0 (Windows NT 10.0; Trident/7.0; rv:11.0) like Gecko", "desktop", middleware.BrowserIE},
	{"", "desktop", middleware.BrowserIE},
}

func TestParseUserAgent(t *testing.T) {
	for _, tua := range testUserAgents {
		kvps := make(map[any]any)
		raw := middleware.ParseUserAgent(tua.ua)
		for k, v := range raw {
			kvps[k] = v
		}

		if len(tua.ua) > 0 {
			actual := kvps[middleware.UserAgentDevice]
			assert.Equal(t, tua.uaType, actual)
		}
	}
}

func TestParseHeaders(t *testing.T) {
	testCases := []struct {
		name           string
		headers        map[string][]string
		expectedResult map[string]any
	}{
		{
			name: "Single header",
			headers: map[string][]string{
				"Content-Type": {"application/json"},
			},
			expectedResult: map[string]any{
				"http.content-type": "application/json",
			},
		},
		{
			name: "Multiple headers",
			headers: map[string][]string{
				"Content-Type":  {"application/json"},
				"Authorization": {"Bearer token123"},
				"Accept":        {"application/json", "text/plain"},
			},
			expectedResult: map[string]any{
				"http.content-type":  "application/json",
				"http.authorization": "bearer token123",
				"http.accept":        "application/json, text/plain",
			},
		},
		{
			name:           "Empty headers",
			headers:        map[string][]string{},
			expectedResult: map[string]any{},
		},
		{
			name: "Header with multiple values",
			headers: map[string][]string{
				"Accept": {"application/json", "text/plain", "application/xml"},
			},
			expectedResult: map[string]any{
				"http.accept": "application/json, text/plain, application/xml",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := middleware.ParseHeaders(tc.headers)
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}
