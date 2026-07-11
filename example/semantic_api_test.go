package example

import (
	"context"
	"testing"

	postgresql "github.com/neteast-software/go-module/db/postgresql/linker"
	event "github.com/neteast-software/go-module/fault/event/linker"
	notice "github.com/neteast-software/go-module/fault/notice/linker"
	http "github.com/neteast-software/go-module/http/gin/linker"
	"github.com/neteast-software/go-module/http/gin/response"
	server "github.com/neteast-software/go-module/linker/server"
	consumer "github.com/neteast-software/go-module/mq/consumer"
	mq "github.com/neteast-software/go-module/mq/consumer/linker"
	prometheus "github.com/neteast-software/go-module/observe/metrics/prometheus/linker"
	opentelemetry "github.com/neteast-software/go-module/observe/tracing/opentelemetry/linker"
	rpc "github.com/neteast-software/go-module/rpc/grpc/linker"
	cron "github.com/neteast-software/go-module/scheduler/cron"
	schedule "github.com/neteast-software/go-module/scheduler/cron/linker"
	linker "github.com/neteast-software/linker/v3"
)

func TestRecommendedSemanticAPIsCompile(t *testing.T) {
	if server.WithMetrics(prometheus.New()) == nil {
		t.Fatal("server metrics option 不能为空")
	}
	if rpc.WithTracing() == nil || rpc.WithMetrics() == nil ||
		rpc.WithClientTracing[any]() == nil || rpc.WithClientMetrics[any]() == nil {
		t.Fatal("rpc observability option 不能为空")
	}
	if mq.WithTracing() == nil || mq.WithMetrics() == nil ||
		schedule.WithTracing() == nil || schedule.WithMetrics() == nil {
		t.Fatal("后台入口 observability option 不能为空")
	}

	item := consumer.New("notice", consumer.HandlerFunc(func(context.Context, consumer.Message) error { return nil }))
	job := cron.NewJob("sync", "@every 1m", cron.HandlerFunc(func(context.Context) error { return nil }))
	if len(mq.Consumers(item)) != 1 || len(schedule.Jobs(job)) != 1 {
		t.Fatal("业务对象应由 linker adapter 转成资产")
	}

	var runtime linker.Runtime
	_, _ = postgresql.Resolve(runtime)
	_, _ = http.ResolveServer(runtime)
	_, _ = event.Resolve(runtime)
	_, _ = notice.ResolveDispatcher(runtime)
	_, _ = mq.Resolve(runtime, "notice")
	_, _ = schedule.Resolve(runtime)
	_, _ = prometheus.Resolve(runtime)
	_, _ = opentelemetry.Resolve(runtime)

	_ = http.RegisterIn
	_ = response.Success
	_ = response.Message
	_ = response.Data
	_ = response.Warning
	_ = response.Fail
	_ = response.Error
}
