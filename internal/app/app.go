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

	console "linker-v3-example/internal/console/linker"
	inspection "linker-v3-example/internal/inspection/linker"
	notification "linker-v3-example/internal/notification/linker"
	order "linker-v3-example/internal/order/linker"
	permission "linker-v3-example/internal/permission/linker"
	ttsclient "linker-v3-example/internal/tts/client"
	tts "linker-v3-example/internal/tts/linker"
	user "linker-v3-example/internal/user/linker"
)

func New(sources ...linker.Source) *linker.App {
	users := user.NewComponent()
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
			inspection.NewComponent(),
			users,
			order.New(),
			permission.New(),
			console.New(users.Service()),
			tts.NewComponent(),
			notification.NewComponent(),
			mq.New(),
			schedule.New(),
			rpc.New(),
			ttsclient.Provider(),
		),
	)
}
