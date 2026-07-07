package observability

import (
	"strconv"
	"time"

	http "github.com/neteast-software/go-module/http/gin/linker"
	"github.com/neteast-software/go-module/observe/metrics"
)

var httpRequests = metrics.Counter(
	"http_requests_total",
	metrics.WithHelp("HTTP 请求数"),
	metrics.WithConstLabels(metrics.Label("service", "linker-v3-example")),
)

var httpRequestSeconds = metrics.Histogram(
	"http_request_seconds",
	metrics.WithHelp("HTTP 请求耗时"),
	metrics.WithUnit("seconds"),
	metrics.WithConstLabels(metrics.Label("service", "linker-v3-example")),
	metrics.WithBuckets(0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5),
)

func init() {
	http.Register(http.Use(HTTPMetrics()))
}

func HTTPMetrics() http.HandlerFunc {
	return func(c *http.Context) {
		started := time.Now()
		c.Next()
		recorder, ok := metricRecorder(c)
		if !ok {
			return
		}
		labels := requestLabels(c)
		ctx := c.Request.Context()
		_ = httpRequests.Add(ctx, recorder, 1, labels...)
		_ = httpRequestSeconds.Observe(ctx, recorder, time.Since(started).Seconds(), labels...)
	}
}

func requestLabels(c *http.Context) []metrics.LabelValue {
	route := c.FullPath()
	if route == "" {
		route = "unmatched"
	}
	return []metrics.LabelValue{
		metrics.Label("method", c.Request.Method),
		metrics.Label("route", route),
		metrics.Label("status", strconv.Itoa(c.Writer.Status())),
	}
}
