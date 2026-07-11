package app

import (
	"time"

	applicationcore "github.com/neteast-software/go-module/application"
	applicationcomponent "github.com/neteast-software/go-module/application/linker"
	postgresql "github.com/neteast-software/go-module/db/postgresql/linker"
	server "github.com/neteast-software/go-module/linker/server"
	mq "github.com/neteast-software/go-module/mq/consumer/linker"
	prometheus "github.com/neteast-software/go-module/observe/metrics/prometheus/linker"
	opentelemetry "github.com/neteast-software/go-module/observe/tracing/opentelemetry/linker"
	rpc "github.com/neteast-software/go-module/rpc/grpc/linker"
	schedule "github.com/neteast-software/go-module/scheduler/cron/linker"
	linker "github.com/neteast-software/linker/v3"

	ttsclient "linker-v3-example/internal/client/tts"
	graphcomponent "linker-v3-example/internal/component/graph"
	inspectioncomponent "linker-v3-example/internal/component/inspection"
	notificationcomponent "linker-v3-example/internal/component/notification"
	ttscomponent "linker-v3-example/internal/component/tts"
	usercomponent "linker-v3-example/internal/component/user"
)

func New(sources ...linker.Source) *linker.App {
	return server.New(
		server.Config(sources...),
		server.WithShutdownTimeout(3*time.Second),
		server.WithMetrics(prometheus.New()),
		server.WithComponents(
			opentelemetry.New(),
			postgresql.New(),
			applicationcomponent.New(
				applicationcomponent.WithApplications(
					applicationcore.Application{ID: "front", Scope: "front", Name: "前台应用", Status: applicationcore.StatusEnabled},
					applicationcore.Application{ID: "console", Scope: "console", Name: "后台应用", Status: applicationcore.StatusEnabled},
					applicationcore.Application{ID: "app2", Scope: "app2", Name: "应用二", Host: "app2.local", Status: applicationcore.StatusEnabled},
				),
			),
			graphcomponent.NewComponent(),
			inspectioncomponent.NewComponent(),
			usercomponent.NewComponent(),
			ttscomponent.NewComponent(),
			notificationcomponent.NewComponent(),
			mq.New(mq.WithTracing(), mq.WithMetrics()),
			schedule.New(schedule.WithTracing(), schedule.WithMetrics()),
			rpc.New(
				rpc.WithTracing(),
				rpc.WithMetrics(),
			),
			ttsclient.Provider(),
		),
	)
}
