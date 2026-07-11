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
	"github.com/neteast-software/go-module/observe/metrics"
	metricserver "github.com/neteast-software/go-module/observe/metrics/linker/server"
	metricgrpc "github.com/neteast-software/go-module/observe/metrics/rpc/grpc"
	tracegrpc "github.com/neteast-software/go-module/observe/tracing/rpc/grpc"
	rpc "github.com/neteast-software/go-module/rpc/grpc/linker"
	cron "github.com/neteast-software/go-module/scheduler/cron"
	schedule "github.com/neteast-software/go-module/scheduler/cron/linker"
	linker "github.com/neteast-software/linker/v3"
)

func TestRecommendedSemanticAPIsCompile(t *testing.T) {
	recorder := metrics.Memory()
	observation := metricserver.Observe(recorder)
	var _ server.Observer = observation
	if server.WithObserver(observation) == nil {
		t.Fatal("server observer option 不能为空")
	}
	if rpc.WithUnaryInterceptors(tracegrpc.UnaryServer(), metricgrpc.UnaryServer(recorder)) == nil ||
		rpc.WithStreamInterceptors(tracegrpc.StreamServer(), metricgrpc.StreamServer(recorder)) == nil {
		t.Fatal("rpc interceptor option 不能为空")
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

	_ = http.RegisterIn
	_ = response.Success
	_ = response.Message
	_ = response.Data
	_ = response.Warning
	_ = response.Fail
	_ = response.Error
}
