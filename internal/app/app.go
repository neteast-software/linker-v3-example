package app

import (
	applicationcore "github.com/neteast-software/go-module/application"
	applicationcomponent "github.com/neteast-software/go-module/application/linker"
	postgresql "github.com/neteast-software/go-module/db/postgresql/linker"
	server "github.com/neteast-software/go-module/linker/server"
	mq "github.com/neteast-software/go-module/mq/consumer/linker"
	"github.com/neteast-software/go-module/observe/metrics"
	metricserver "github.com/neteast-software/go-module/observe/metrics/linker/server"
	metricgrpc "github.com/neteast-software/go-module/observe/metrics/rpc/grpc"
	tracegrpc "github.com/neteast-software/go-module/observe/tracing/rpc/grpc"
	rpccore "github.com/neteast-software/go-module/rpc/grpc"
	rpc "github.com/neteast-software/go-module/rpc/grpc/linker"
	cron "github.com/neteast-software/go-module/scheduler/cron/linker"
	linker "github.com/neteast-software/linker/v3"
	stdgrpc "google.golang.org/grpc"

	ttsclient "linker-v3-example/internal/client/tts"
	graphcomponent "linker-v3-example/internal/component/graph"
	inspectioncomponent "linker-v3-example/internal/component/inspection"
	notificationcomponent "linker-v3-example/internal/component/notification"
	observabilitycomponent "linker-v3-example/internal/component/observability"
	ttscomponent "linker-v3-example/internal/component/tts"
	usercomponent "linker-v3-example/internal/component/user"
	"linker-v3-example/internal/config"
)

func New(config config.Config) (*linker.App, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}
	observability := observabilitycomponent.NewComponent()
	metricLabels := []metrics.LabelValue{metrics.Label("service", "linker-v3-example")}
	notification := notificationcomponent.NewComponent(
		notificationcomponent.WithMetricRecorder(observability.Recorder()),
		notificationcomponent.WithMetricLabels(metricLabels...),
	)
	return server.New(
		server.WithShutdownTimeout(config.ShutdownTimeout),
		server.WithHTTP(config.HTTP),
		server.WithLifecycleObserver(metricserver.Observer(
			observability.Recorder(),
			metricserver.WithConstLabels(metricLabels...),
		)),
		server.WithComponents(
			postgresql.New(postgresql.WithConfig(config.PostgreSQL)),
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
			usercomponent.NewComponent([]byte(config.TokenKey)),
			ttscomponent.NewComponent(),
			notification,
			mq.New(),
			cron.New(),
			rpc.New(
				rpc.WithConfig(config.GRPC),
				rpc.WithServerOptions(stdgrpc.ChainUnaryInterceptor(
					tracegrpc.UnaryServer(),
					metricgrpc.UnaryServer(observability.Recorder(), metricgrpc.WithConstLabels(metricLabels...)),
					rpccore.UnaryServerMeta(),
				)),
			),
			ttsclient.Provider(
				config.TTSClient,
				ttsclient.WithMetricRecorder(observability.Recorder()),
				ttsclient.WithMetricLabels(metricLabels...),
			),
		),
	), nil
}
