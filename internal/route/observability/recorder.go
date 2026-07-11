package observability

import (
	http "github.com/neteast-software/go-module/http/gin/linker"
	metricsprom "github.com/neteast-software/go-module/observe/metrics/prometheus"

	observabilityservice "linker-v3-example/internal/service/observability"
)

func prometheusRecorder(c *http.Context) (*metricsprom.Recorder, bool) {
	recorder, ok := http.Resolve(c, observabilityservice.PrometheusRecorderKey())
	return recorder, ok && recorder != nil
}
