package app

import (
	"time"

	audit "github.com/neteast-software/go-module/audit/operate"
	event "github.com/neteast-software/go-module/fault/event"
	notice "github.com/neteast-software/go-module/fault/notice"
	gatewaycore "github.com/neteast-software/go-module/http/gateway"
	"github.com/neteast-software/go-module/http/gateway/declaration"
	gatewaycomponent "github.com/neteast-software/go-module/http/gateway/linker"
	httpcore "github.com/neteast-software/go-module/http/gin"
	server "github.com/neteast-software/go-module/linker/server"
	gatewaymetrics "github.com/neteast-software/go-module/observe/metrics/http/gateway"
	prometheus "github.com/neteast-software/go-module/observe/metrics/prometheus/linker"
	gatewaytracing "github.com/neteast-software/go-module/observe/tracing/http/gateway"
	opentelemetry "github.com/neteast-software/go-module/observe/tracing/opentelemetry/linker"
	discoverynacos "github.com/neteast-software/go-module/registry/discovery/nacos/linker"
	nacos "github.com/neteast-software/go-module/registry/nacos/linker"
	linker "github.com/neteast-software/linker/v3"

	accesslog "linker-v3-example/internal/accesslog"
	accesslogcomponent "linker-v3-example/internal/accesslog/linker"
	captcha "linker-v3-example/internal/captcha"
	examplegateway "linker-v3-example/internal/gateway"
	session "linker-v3-example/internal/session"
	vendorauth "linker-v3-example/internal/vendorauth"
)

const localUpstream = "http://127.0.0.1:8800"
const localManagementAddress = "127.0.0.1:8820"

// Gateway 创建默认本地静态代理工作背景。
func Gateway(sources ...linker.Source) (*linker.App, error) {
	return gatewayApp(localDocument(), nil, sources...)
}

// GatewayNacos 创建显式使用 Nacos 服务发现的工作背景。
func GatewayNacos(sources ...linker.Source) (*linker.App, error) {
	return gatewayApp(
		nacosDocument(),
		[]linker.Component{
			nacos.New(),
			discoverynacos.New(discoverynacos.Poll(2*time.Second, 0.2)),
		},
		sources...,
	)
}

func gatewayApp(
	document declaration.Document,
	providers []linker.Component,
	sources ...linker.Source,
) (*linker.App, error) {
	metrics := prometheus.New()
	tracing := opentelemetry.New()
	log, err := accesslog.New(accesslog.Discard(), 256)
	if err != nil {
		return nil, err
	}
	sessions := session.Memory()
	sessionFilter := session.Filter(session.Legacy(nil, "example-login", "example-api"), sessions)
	vendorFilter := vendorauth.Filter(
		vendorauth.RS256(nil, "example-vendor", "example-api"),
		vendorauth.Scope("equipment-read", "equipment.read", true),
	)
	captchaFilter := captcha.Filter(captcha.Memory(), "")
	config := gatewaycomponent.DefaultConfig()
	config.Address = "127.0.0.1:8810"
	config.Health = gatewaycomponent.StandardHealth()
	if len(providers) > 0 {
		config.Discovery = "nacos"
	}
	component := gatewaycomponent.New(
		document,
		gatewaycomponent.WithConfig(config),
		gatewaycomponent.WithFilters(vendorFilter, sessionFilter, session.Upload(nil, sessions), captchaFilter),
		gatewaycomponent.WithGateway(gatewaycore.Observers(
			gatewaymetrics.Observe(metrics),
			gatewaytracing.Observe(tracing),
			log,
		)),
	)
	components := append([]linker.Component{tracing, accesslogcomponent.New(log)}, providers...)
	return server.New(
		server.Config(sources...),
		server.WithShutdownTimeout(10*time.Second),
		server.WithEventRecorder(event.NewMemoryRecorder()),
		server.WithNoticeSenders(notice.NewMemorySender()),
		server.WithAuditRecorder(audit.NewMemoryRecorder()),
		server.WithMetrics(metrics),
		server.WithComponents(components...),
		server.WithGateway(component),
		server.WithHTTP(gatewayManagementConfig()),
	), nil
}

func gatewayManagementConfig() httpcore.Config {
	config := httpcore.DefaultConfig()
	config.Addr = localManagementAddress
	return config
}

func localDocument() declaration.Document {
	return examplegateway.Document(
		examplegateway.Strip(examplegateway.URL("local/public", "/example/**", localUpstream), 1),
		examplegateway.Strip(
			vendorauth.Protect(
				examplegateway.URL("local/vendor-equipment", "/vendor/equipment/**", localUpstream),
				"equipment-read",
			),
			2,
		),
		examplegateway.Strip(
			session.Protect(
				examplegateway.URL("local/profile", "/session/profile/**", localUpstream),
				"profile.read",
				"console",
			),
			2,
		),
		examplegateway.Strip(
			captcha.Login(examplegateway.URL("local/login", "/login/**", localUpstream)),
			1,
		),
		examplegateway.Strip(
			session.Socket(examplegateway.URL("local/socket", "/socket/**", localUpstream), "/socket/"),
			1,
		),
		examplegateway.Strip(
			session.UploadRoute(examplegateway.URL("local/upload", "/file/upload/file", localUpstream)),
			2,
		),
		examplegateway.Strip(
			session.Protect(examplegateway.URL("local/file", "/file/**", localUpstream), "file.read", ""),
			1,
		),
	)
}

func nacosDocument() declaration.Document {
	return examplegateway.Document(
		examplegateway.StripService(
			examplegateway.Service("nacos/example", "/example-http/**", "example-http"),
		),
	)
}
