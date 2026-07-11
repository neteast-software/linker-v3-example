# linker-v3-example

`linker-v3-example` 是 linker v3 的演示业务系统，用来验证 framework、HTTP route、ACL resource、PostgreSQL 生命周期和业务组件声明方式。

## 最短运行

本地 YAML server：

```bash
APP_DB_POSTGRESQL__PASSWORD='...' \
APP_EXAMPLE_USER__TOKEN_KEY='至少 32 个字符' \
go run .
```

在本地 seed 后叠加 Nacos 完整配置，再由 env 最终覆盖：

```bash
LINKER_V3_EXAMPLE_NACOS_DATA_ID='app.yaml' \
LINKER_V3_EXAMPLE_NACOS_HOST='nacos.example.internal' \
APP_DB_POSTGRESQL__PASSWORD='...' \
APP_EXAMPLE_USER__TOKEN_KEY='至少 32 个字符' \
go run .
```

构建 bin 并只查看装配计划，不连接外部资源：

```bash
go build -o ./bin/linker-v3-example .
./bin/linker-v3-example --plan
```

配置来源始终按声明顺序覆盖：`local YAML -> optional Nacos final -> explicit env override`。完整变量、动态配置与集成测试说明见后文。

## 业务能力

- `POST /system/login`：后台管理员登录。
- `GET /system/profile`：后台管理员信息。
- `POST /user/login`：前台用户登录。
- `GET /api/profile`：返回前台用户信息，包括用户名、头像、邮箱和手机号。
- `GET /api/v1/app2/user/:id/profile`：多层 route 示例，实际访问形如 `/api/v1/app2/user/3/profile`。
- `GET /api/v1/app2/inspection/tasks`：巡检任务列表，演示 application data scope、分页查询和响应白名单。
- `GET /api/v1/app2/notification/events`：SSE 事件入口，演示长连接 route 的局部声明。
- `POST /api/v1/app2/notification/send`：HTTP 到 MQ mock 的 trace 贯穿示例。
- `POST /api/v1/app2/tts/transcribe`：HTTP 到 typed gRPC client 的 trace 贯穿示例。
- `GET /metrics`：Prometheus scrape 入口，演示 observability 组件、HTTP/gRPC/MQ/cron 指标、低基数 label 和 Grafana dashboard。
- `GET /api/v1/app2/graph/orders`：graph/naive viewer 示例。
- `GET /api/v1/app2/graph/orders/form`：graph/naive form 示例。
- `GET /api/v1/app2/graph/refresh`：graph/naive behavior 示例。
- `example.tts.TTS/Transcribe`：gRPC service，演示 RPC register、typed client provider、trace 传播和表资产。
- `example/http_client_example_test.go`：出站 HTTP client 示例，演示 `http/client/linker` capability、Plan asset、credential、trace hook 和业务 typed client。

登录链路使用 modules 的边界：

- `user/account`：账号来源和密码凭据，后台账号、手机号账号都映射到内部 user。
- `token`：HMAC token 签发、验证、过期和 claims。
- `acl/session`：token 存活和撤销状态。
- `acl.Resource`：每个 route 自己声明权限资源和 scope。
- `audit/operate`：登录这类敏感操作记录 actor、client、request、resource 和成功状态。

HTTP API 按文件自声明 route：`system_login.go`、`user_login.go`、`profile_api.go` 各自通过 `init()` 注册自己的入口，常用形式是 `http.RegisterIn("api", http.GET("profile", profileAPI).Resource("http.front.user.profile", acl.Scope(...)))`。组件声明自己的 identity、表资产、生命周期和 service capability，API 通过 `http.Require(c, user.ServiceKey())` 从 linker runtime 解析能力，不在业务侧维护全局 runtime 容器。

多层 route 可以把稳定前缀放进 `RegisterIn`：

```go
func init() {
    http.RegisterIn("api/v1/app2",
        http.GET("user/:id/profile", app2Profile).Resource(
            "http.app2.user.profile",
            acl.Scope("app2", 1, "应用二用户资料"),
        ),
    )
}
```

需要控制 middleware 影响面时，再用 `Group` 表达局部层级：

```go
func init() {
    http.RegisterIn("api/v1/app2",
        http.Group("user/:id",
            http.Use(requireApp2),
            http.GET("profile", app2Profile),
            http.GET("settings", app2Settings).With(requireOwner),
        ),
    )
}
```

业务代码按职责域 package 组织：

- `internal/route/user`：HTTP 入口和 route 声明。
- `internal/model/user`：数据表模型，包含 user 主体表和 account 凭据表。
- `internal/service/user`：登录、资料读取、token/session 和存储流程；service capability key 由 service 自己声明。
- `internal/constant/user`：业务错误和明确的 example fixture。
- `internal/component/user`：组件 identity、linker 组件生命周期、表资产和 service capability 挂载。

record-level 权限建议放在具体业务 store 的查询入口处完成。`internal/service/inspection` 用 `TaskAccess` 把 `acl.Access`、`acl.Resource` 和 `RecordRange` 组合在一起：route 只提供当前 application 和 actor，store 在一次查询里同时应用 application scope、业务 filter 和 owner range，避免为了权限判断额外做 N+1 查询或维护 RBAC 关系表。
- `internal/route/inspection`、`internal/model/inspection`、`internal/service/inspection`、`internal/component/inspection`：接近真实业务的列表接口结构，route 负责 HTTP 参数和输出，service/store 负责批量查询和数据范围。
- `internal/model/inspection/archive.go`：外部维护表资产示例，只改业务 model 和 component asset，使用 `postgresql.External()` 避免启动期迁移。
- `internal/component/notification`、`internal/service/notification`、`internal/route/notification`：MQ consumer、cron job、SSE route、trace/metrics wrapper 和 provider mock 的长生命周期组合。
- `internal/component/observability`、`internal/service/observability`、`internal/route/observability`：Prometheus recorder capability、`/metrics` route、middleware 影响面和 Plan 里的 metrics/tracing asset；HTTP 指标 middleware 实现统一位于 `internal/route/middleware`。
- `license/http/gin`：示例只在需要保护的入口显式挂 `licensehttp.Gate(gate)`；license 不进入 core，也不默认挂到 server framework。
- `internal/rpc/tts`、`internal/client/tts`、`internal/component/tts`：gRPC server/client 的声明、注册、trace/metrics interceptor 和 capability provider。
- `internal/client/directory`：出站 HTTP typed client 示例，第三方用户目录 API 的业务语义在这里承载，通用 HTTP 执行者来自 `http/client` capability。

组件 identity 必须由组件 package 自己声明，例如 `component/user.ID`。需要依赖该组件时引用这个符号，不把组件 ID 放到 `constant` 或其他公共包里代管。

推荐边界：

- `component/user` 不作为 HTTP controller，不出现成片的 `http.Context` handler、`response.*` 和 route tree 聚合。
- `component/user` 不聚合所有 HTTP handler，也不替 route 维护完整 API 树。
- `component/user` 通过 blank import 纳入 `route/user`，表示启用该组件时才编译这些 route。
- service capability key 由 `service/user` 声明，route 通过 `http.Require(c, user.ServiceKey())` 获取能力，避免 route 反向依赖 component。
- 跨 route middleware 应集中在明确位置；单个 API 只声明 middleware 影响面。

默认演示数据只作为本地 example fixture：

- 后台账号：`admin`
- 前台手机号：`18558755877`
- 默认密码：`linfunlinfun`

## 运行

默认读取 `config/app.example.yaml`。文件只声明本地 listener、RPC 和 PostgreSQL 示例目标，不包含数据库密码、Nacos 凭据或 token key。启动前至少需要注入：

```bash
export APP_DB_POSTGRESQL__PASSWORD='...'
export APP_EXAMPLE_USER__TOKEN_KEY='至少 32 个字符'
```

查看 linker 装配计划不需要连接数据库：

```bash
go run . --plan
```

输出会包含 mode、components、dependencies、capabilities、Config mode/revision 和 application、route、gRPC、MQ consumer、cron job、metrics、tracing 等 assets。`--plan` 只在进程内生成一次临时 token key 用于 Bootstrap，不向仓库、Plan 或日志写入凭据，也不会启动 PostgreSQL component。

```bash
go run .
```

本地需要覆盖数据库地址时：

```bash
APP_DB_POSTGRESQL__HOST=127.0.0.1 \
APP_DB_POSTGRESQL__PASSWORD=... \
APP_EXAMPLE_USER__TOKEN_KEY=... \
go run .
```

项目只有一条配置主路径：`local YAML -> optional Nacos final -> explicit env override`。每个 component 在 Bootstrap 解码和校验自己的 namespace，`internal/app` 不再接收或拆分一份全局 typed Config：

```bash
LINKER_V3_EXAMPLE_CONFIG=config/app.example.yaml \
APP_DB_POSTGRESQL__PASSWORD=... \
APP_EXAMPLE_USER__TOKEN_KEY=... \
go run .
```

设置 `LINKER_V3_EXAMPLE_NACOS_DATA_ID` 后会在 local 与 env 之间加入 Nacos Source；bootstrap endpoint 和凭据使用 `LINKER_V3_EXAMPLE_NACOS_HOST/PORT/USERNAME/PASSWORD`，不进入业务 namespace。`example/dynamic_config_test.go` 使用真实 Nacos adapter 的 fake fetch/listen，演示 HTTP client Live snapshot、后置 env 覆盖、非法 YAML 恢复和 PostgreSQL Restart 标记。

`Live` namespace 先完成整批 typed decode/validate，再原子发布到新操作；`Restart` namespace 只进入 desired Setting 并标记整个服务需要重启。业务 component 不实现 Nacos reload，也不会被 framework 隐式重启。

## Example

测试文件集中在 `example/` 目录。真实 PostgreSQL example 只在设置 `LINKER_V3_EXAMPLE_PG_PASSWORD` 后运行，host 默认使用 `127.0.0.1` 且可通过同前缀测试变量覆盖；当前环境不可用时会明确 skip。`signal_example_test.go` 会启动真实 server 子进程，分别在 startup 和 running 阶段发送 SIGTERM，验证反向关闭和退出码，不依赖外部 provider。`production_http_example_test.go` 展示 body limit、trusted proxy、health endpoint 和负载中 graceful stop；`docs/Caddyfile` 是部署层终止 TLS 的最小样板。

```bash
go test ./...
```

Prometheus 可抓取 `GET /metrics`，Grafana 示例面板在 `docs/grafana-dashboard.json`。当前 dashboard 对齐 HTTP、gRPC、MQ consumer、cron、linker runtime Plan 和动态配置指标，包括 active/desired revision、`restart_required` 与固定状态事件；namespace 和错误文本不进入 metrics label。

推荐先看：

- `docs/scaffold.md`：推荐项目骨架，说明 main、app、component、route、model、service、config、observability 和 example test 的边界。
- `docs/example-policy.md`：说明 example 的定位、外部依赖、测试拆分和未来 submodule 边界。
- `main.go`：保持极薄，只分发 server 启动和 `--plan`。
- `source.go`：只读取配置文件位置和 Nacos bootstrap 参数，按顺序声明 Source。
- `internal/app/app.go`：集中装配 framework、组件和 adapter，不解析业务配置。
- `internal/route/graph/*_api.go`：一个 API 一个文件，route/resource/middleware 和 handler 放在同一个入口重心内。
- `example/graph_example_test.go`：验证 graph route、Plan asset 和 renderer capability。
- `example/nacos_example_test.go`：验证 YAML seed、Nacos source、HTTP/gRPC registry adapter 和 Plan 里的依赖/capability 表达。
- `example/dynamic_config_test.go`：验证 Nacos 完整快照、Live/Restart、desired/active、env 后置覆盖和可恢复拒绝。
- `example/reliability_example_test.go`：验证 DB capability 缺失会在组件初始化期失败，以及 Stop timeout 会返回可判断的 `context.DeadlineExceeded`。
- `example/signal_example_test.go`：验证真实 server 在启动期和运行期收到 SIGTERM 后都执行 graceful close；运行期正常退出，启动期返回可判断的取消原因。
- `example/health_example_test.go`：验证 liveness、readiness、startup 与 framework starting/running 状态一致。
- `example/production_http_example_test.go`：验证 body limit、trusted proxy 和在途请求 graceful drain。
- `docs/Caddyfile`：由 Caddy 终止 TLS，再转发到只监听 loopback 的 linker HTTP server。
- `example/grpc_example_test.go`：验证 gRPC metadata 和 trace id 通过 interceptor 传播。
- `example/http_client_example_test.go`：验证出站 HTTP client linker adapter、typed client、credential、trace hook 和 Plan asset。
- `example/notification_example_test.go`：验证 MQ/cron/SSE lifecycle，并覆盖 HTTP -> MQ mock 的 trace id 贯穿。
- `example/feishu_notify_example_test.go`：验证 component 上报 Fault 后由 server 自动发现飞书 Sender，并形成 detected/recovered 通知闭环。
- `example/business_system_test.go`：验证完整业务系统，并覆盖 HTTP -> gRPC typed client 的 trace id 贯穿。
