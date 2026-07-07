package observability

import (
	http "github.com/neteast-software/go-module/http/gin/linker"
	"github.com/neteast-software/go-module/observe/metrics"
	metricsprom "github.com/neteast-software/go-module/observe/metrics/prometheus"

	observabilityservice "linker-v3-example/internal/service/observability"
)

func metricRecorder(c *http.Context) (metrics.Recorder, bool) {
	recorder, ok := http.Resolve(c, observabilityservice.MetricRecorderKey())
	return recorder, ok && recorder != nil
}

func prometheusRecorder(c *http.Context) (*metricsprom.Recorder, bool) {
	recorder, ok := http.Resolve(c, observabilityservice.PrometheusRecorderKey())
	return recorder, ok && recorder != nil
}
