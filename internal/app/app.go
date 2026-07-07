package app

import (
	"encoding/json"

	applicationcore "github.com/neteast-software/go-module/application"
	applicationcomponent "github.com/neteast-software/go-module/application/linker"
	postgresql "github.com/neteast-software/go-module/db/postgresql/linker"
	server "github.com/neteast-software/go-module/linker/server"
	mq "github.com/neteast-software/go-module/mq/consumer/linker"
	rpccore "github.com/neteast-software/go-module/rpc/grpc"
	rpc "github.com/neteast-software/go-module/rpc/grpc/linker"
	cron "github.com/neteast-software/go-module/scheduler/cron/linker"
	linker "github.com/neteast-software/linker/v3"
	stdgrpc "google.golang.org/grpc"

	ttsclient "linker-v3-example/internal/client/tts"
	inspectioncomponent "linker-v3-example/internal/component/inspection"
	notificationcomponent "linker-v3-example/internal/component/notification"
	observabilitycomponent "linker-v3-example/internal/component/observability"
	ttscomponent "linker-v3-example/internal/component/tts"
	usercomponent "linker-v3-example/internal/component/user"
	"linker-v3-example/internal/config"
)

func New(config config.Config) (*server.App, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}
	postgresqlConfig, err := json.Marshal(config.PostgreSQL)
	if err != nil {
		return nil, err
	}
	notification := notificationcomponent.NewComponent()
	return server.New(
		server.WithMode(linker.Server),
		server.WithShutdownTimeout(config.ShutdownTimeout),
		server.WithHTTP(config.HTTP),
		server.WithComponents(
			postgresql.New(),
			applicationcomponent.New(
				applicationcomponent.WithApplications(
					applicationcore.Application{ID: "front", Scope: "front", Name: "前台应用", Status: applicationcore.StatusEnabled},
					applicationcore.Application{ID: "console", Scope: "console", Name: "后台应用", Status: applicationcore.StatusEnabled},
					applicationcore.Application{ID: "app2", Scope: "app2", Name: "应用二", Host: "app2.local", Status: applicationcore.StatusEnabled},
				),
			),
			observabilitycomponent.NewComponent(),
			inspectioncomponent.NewComponent(),
			usercomponent.NewComponent(),
			ttscomponent.NewComponent(),
			notification,
			mq.New(
				mq.WithAfter(notificationcomponent.ID),
				mq.WithConsumers(notification.Consumer()),
			),
			cron.New(
				cron.WithAfter(notificationcomponent.ID),
				cron.WithJobs(notification.Job()),
			),
			rpc.New(
				rpc.WithConfig(config.GRPC),
				rpc.WithAfter(ttscomponent.ID),
				rpc.WithServerOptions(stdgrpc.ChainUnaryInterceptor(rpccore.UnaryServerMeta())),
			),
			ttsclient.Provider(config.TTSClient),
		),
		server.WithMapSetting(map[linker.Namespace][]byte{
			linker.Namespace(postgresql.ID): postgresqlConfig,
		}),
	), nil
}
