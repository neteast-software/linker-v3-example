package app

import (
	"time"

	applicationcore "github.com/neteast-software/go-module/application"
	applicationcomponent "github.com/neteast-software/go-module/application/linker"
	postgresql "github.com/neteast-software/go-module/db/postgresql/linker"
	server "github.com/neteast-software/go-module/linker/server"
	mq "github.com/neteast-software/go-module/mq/consumer/linker"
	"github.com/neteast-software/go-module/observe/metrics"
	metricserver "github.com/neteast-software/go-module/observe/metrics/linker/server"
	metricgrpc "github.com/neteast-software/go-module/observe/metrics/rpc/grpc"
	tracegrpc "github.com/neteast-software/go-module/observe/tracing/rpc/grpc"
	rpc "github.com/neteast-software/go-module/rpc/grpc/linker"
	schedule "github.com/neteast-software/go-module/scheduler/cron/linker"
	linker "github.com/neteast-software/linker/v3"

	ttsclient "linker-v3-example/internal/client/tts"
	graphcomponent "linker-v3-example/internal/component/graph"
	inspectioncomponent "linker-v3-example/internal/component/inspection"
	notificationcomponent "linker-v3-example/internal/component/notification"
	observabilitycomponent "linker-v3-example/internal/component/observability"
	ttscomponent "linker-v3-example/internal/component/tts"
	usercomponent "linker-v3-example/internal/component/user"
)

func New(sources ...linker.Source) *linker.App {
	observability := observabilitycomponent.NewComponent()
	metricLabels := []metrics.LabelValue{metrics.Label("service", "linker-v3-example")}
	notification := notificationcomponent.NewComponent(
		notificationcomponent.WithMetricRecorder(observability.Recorder()),
		notificationcomponent.WithMetricLabels(metricLabels...),
	)
	return server.New(
		server.Config(sources...),
		server.WithShutdownTimeout(3*time.Second),
		server.WithObserver(metricserver.Observe(
			observability.Recorder(),
			metricserver.WithConstLabels(metricLabels...),
		)),
		server.WithComponents(
			postgresql.New(),
			applicationcomponent.New(
				applicationcomponent.WithApplications(
					applicationcore.Application{ID: "front", Scope: "front", Name: "前台应用", Status: applicationcore.StatusEnabled},
					applicationcore.Application{ID: "console", Scope: "console", Name: "后台应用", Status: applicationcore.StatusEnabled},
					applicationcore.Application{ID: "app2", Scope: "app2", Name: "应用二", Host: "app2.local", Status: applicationcore.StatusEnabled},
				),
			),
			graphcomponent.NewComponent(),
			observability,
			inspectioncomponent.NewComponent(),
			usercomponent.NewComponent(),
			ttscomponent.NewComponent(),
			notification,
			mq.New(),
			schedule.New(),
			rpc.New(
				rpc.WithUnaryInterceptors(
					tracegrpc.UnaryServer(),
					metricgrpc.UnaryServer(observability.Recorder(), metricgrpc.WithConstLabels(metricLabels...)),
				),
				rpc.WithStreamInterceptors(
					tracegrpc.StreamServer(),
					metricgrpc.StreamServer(observability.Recorder(), metricgrpc.WithConstLabels(metricLabels...)),
				),
			),
			ttsclient.Provider(
				ttsclient.WithMetricRecorder(observability.Recorder()),
				ttsclient.WithMetricLabels(metricLabels...),
			),
		),
	)
}
