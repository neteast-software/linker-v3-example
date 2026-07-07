package observability

import (
	"github.com/neteast-software/go-module/observe/metrics"
	metricsprom "github.com/neteast-software/go-module/observe/metrics/prometheus"
	linker "github.com/neteast-software/linker/v3"
)

const MetricRecorderID linker.ID = "example/observability/metrics"
const PrometheusRecorderID linker.ID = "example/observability/prometheus"

func MetricRecorderKey() linker.CapabilityKey[metrics.Recorder] {
	return linker.NewCapabilityKey[metrics.Recorder](MetricRecorderID)
}

func PrometheusRecorderKey() linker.CapabilityKey[*metricsprom.Recorder] {
	return linker.NewCapabilityKey[*metricsprom.Recorder](PrometheusRecorderID)
}
