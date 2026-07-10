package observability

import (
	"context"

	"github.com/neteast-software/go-module/observe/metrics"
	metricsprom "github.com/neteast-software/go-module/observe/metrics/prometheus"
	linker "github.com/neteast-software/linker/v3"

	_ "linker-v3-example/internal/route/observability" // route 声明随组件进入编译
	observabilityservice "linker-v3-example/internal/service/observability"
)

const ID linker.ID = "example/observability"

type Component struct {
	recorder *metricsprom.Recorder
}

type Option func(*Component)

func NewComponent(opts ...Option) *Component {
	p := &Component{
		recorder: metricsprom.New(metricsprom.WithNamespace("linker_v3_example")),
	}
	for _, opt := range opts {
		if opt != nil {
			opt(p)
		}
	}
	return p
}

func WithRecorder(recorder *metricsprom.Recorder) Option {
	return func(p *Component) {
		if recorder != nil {
			p.recorder = recorder
		}
	}
}

func (p *Component) Identity() linker.ID {
	return ID
}

func (p *Component) Recorder() metrics.Recorder {
	if p == nil || p.recorder == nil {
		return metrics.Noop()
	}
	return p.recorder
}

func (p *Component) Assets(context.Context, linker.Runtime) ([]linker.Asset, error) {
	return []linker.Asset{
		linker.NewAsset("observe/metrics", ID, p.recorder).Describe(
			"prometheus",
			map[string]string{
				"endpoint":  "/metrics",
				"namespace": "linker_v3_example",
				"recorder":  "prometheus",
			},
		),
		linker.NewAsset("observe/tracing", ID, "http+grpc").Describe(
			"http+grpc",
			map[string]string{
				"headers":   "X-Trace-ID,X-Span-ID,X-Request-ID",
				"sampling":  "always",
				"transport": "http,grpc",
			},
		),
	}, nil
}

func (p *Component) OnMounted(_ context.Context, runtime linker.Runtime) error {
	if err := linker.Provide(runtime, observabilityservice.MetricRecorderKey(), metrics.Recorder(p.recorder)); err != nil {
		return err
	}
	return linker.Provide(runtime, observabilityservice.PrometheusRecorderKey(), p.recorder)
}
