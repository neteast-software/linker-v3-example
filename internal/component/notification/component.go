package notification

import (
	"context"

	auditcore "github.com/neteast-software/go-module/audit/operate"
	audit "github.com/neteast-software/go-module/audit/operate/linker"
	eventcore "github.com/neteast-software/go-module/fault/event"
	event "github.com/neteast-software/go-module/fault/event/linker"
	consumer "github.com/neteast-software/go-module/mq/consumer"
	mq "github.com/neteast-software/go-module/mq/consumer/linker"
	cron "github.com/neteast-software/go-module/scheduler/cron"
	schedule "github.com/neteast-software/go-module/scheduler/cron/linker"
	linker "github.com/neteast-software/linker/v3"

	_ "linker-v3-example/internal/route/notification" // route 声明随组件进入编译
	service "linker-v3-example/internal/service/notification"
)

const ID linker.ID = "example/notification"

type Component struct {
	provider *service.Provider
	consumer *consumer.Consumer
	job      cron.Job
	audit    auditcore.Recorder
	event    eventcore.Recorder
	cronSpec string
}

type Option func(*Component)

func NewComponent(opts ...Option) *Component {
	p := &Component{
		provider: service.NewProvider(),
		audit:    auditcore.NoopRecorder{},
		event:    eventcore.NoopRecorder{},
		cronSpec: "@every 10m",
	}
	for _, opt := range opts {
		if opt != nil {
			opt(p)
		}
	}
	p.consumer = consumer.New("notification", consumer.HandlerFunc(p.handleMessage),
		consumer.WithDesc("通知消息消费"),
		consumer.WithTopic("notification.message"),
		consumer.WithBuffer(16),
		consumer.WithBackpressure(consumer.BackpressureReject),
	)
	p.job = cron.NewJob("notification.health", p.cronSpec, cron.HandlerFunc(p.runJob),
		cron.WithDesc("通知服务健康采样"),
	)
	return p
}

func WithProvider(provider *service.Provider) Option {
	return func(p *Component) {
		if provider != nil {
			p.provider = provider
		}
	}
}

func WithCronSpec(spec string) Option {
	return func(p *Component) {
		if spec != "" {
			p.cronSpec = spec
		}
	}
}

func (p *Component) Identity() linker.ID {
	return ID
}

func (p *Component) Dependencies() []linker.Dependency {
	return []linker.Dependency{
		linker.StartAfter(audit.ID),
		linker.StartAfter(event.ID),
	}
}

func (p *Component) Assets(context.Context, linker.Runtime) ([]linker.Asset, error) {
	return append(
		mq.Consumers(p.consumer),
		schedule.Jobs(p.job)...,
	), nil
}

func (p *Component) Init(_ context.Context, runtime linker.Runtime) error {
	if recorder, ok := audit.Resolve(runtime); ok && recorder != nil {
		p.audit = recorder
	}
	if recorder, ok := event.Resolve(runtime); ok && recorder != nil {
		p.event = recorder
	}
	return nil
}

func (p *Component) OnMounted(_ context.Context, runtime linker.Runtime) error {
	return linker.Provide(runtime, service.ProviderKey(), p.provider)
}

func (p *Component) Provider() *service.Provider {
	return p.provider
}
