# linker-v3-example

`linker-v3-example` 是 linker v3 的演示业务系统，用来验证 framework、HTTP route、ACL resource、PostgreSQL 生命周期和业务组件声明方式。

当前工具链基线为 Go `1.26.5`，framework 基线为 linker `v3.3.3`；各自治 module 按自身版本发布，精确依赖以 `go.mod` 为准。现代 Go 能力和采用边界从 linker [`GO.md`](https://github.com/neteast-software/linker/blob/v3/GO.md) 进入。`go.mod` 中剩余的本地 `replace` 只对应尚未发布的 source-ready 能力，不覆盖已经发布的 canonical API。

当前阶段新建 server 的关系型数据库推荐 PostgreSQL，并通过 `db/postgresql` 与可选 linker adapter 接入；普通业务模型嵌入 `db/gorm/model` 的 `model.Head` 接管自增 ID。这是脚手架默认选型，不让 Linker core 依赖数据库，也不削弱其他自治数据库 module。

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
- `GET /metrics`：Prometheus scrape 入口，由 `server.WithMetrics` 统一装配 HTTP/gRPC/MQ/cron/component/fault/notice 指标。
- `GET /console/entry`、`POST /console/login`：Graph Console 公开入口和登录。
- `GET /console/menu`、`GET /console/page/:id`：经过 session 与 ACL 的菜单和动态页面。
- `GET|PUT /api/v1/app2/order`：viewer/form 对应的真实业务查询、服务端验证和写入边界。
- `GET|POST|DELETE /api/v1/app2/permission/role/:id/resource`：multilist 对应的权限关系增量接口。
- `example.tts.TTS/Transcribe`：gRPC service，演示 RPC register、typed client provider、trace 传播和表资产。
- `example/http_client_example_test.go`：出站 HTTP client 示例，演示 `http/client/linker` capability、Plan asset、credential、trace hook 和业务 typed client。
- `example/oauth_example_test.go`：可选 OAuth 示例，演示 JWT provider、issuer/audience/scope 校验、Gin Bearer middleware 以及与 ACL 的组合；不进入 Linker 默认组件。
- `example/postgresql_brownfield_example_test.go`：历史项目迁移示例，演示有 model 的既有表、无 model 的 SQL-only 命名表、旧 sqlc 查询层的共池共事务兼容和集中 transition Asset；新项目仍统一推荐 GORM。
- `example/periodic_worker_example_test.go`：受管周期 Worker 示例，演示自治 Worker、Linker adapter、Plan Asset、capability 与 graceful stop。
- `example/grpc_example_test.go`：进阶 provider 示例，演示官方 interceptor 与 metrics/tracing 组合；普通业务优先使用 typed client 和 Linker register asset。

主程序和业务目录只使用能力地图 `call` 主路径。测试为了取得随机 listener 地址、发送真实请求或验证 provider 边界，可以使用明确的 `advanced_call`；棕地测试只验证 `compatibility_call`，两者都不作为新业务脚手架。

登录链路使用 modules 的边界：

- `user/account`：账号来源和密码凭据，后台账号、手机号账号都映射到内部 user。
- `token`：HMAC token 签发、验证、过期和 claims。
- `acl/session`：token 存活和撤销状态。
- `acl.Resource`：每个 route 自己声明权限资源和 scope。
- `audit/operate`：登录这类敏感操作记录 actor、client、request、resource 和成功状态。

HTTP API 按业务文件自声明 route：`system_login.go`、`user_login.go`、`profile.go` 各自通过 `init()` 注册自己的入口，常用形式是 `http.RegisterIn("api", http.GET("profile", frontProfile).Resource("http.front.user.profile", acl.Scope(...)))`。`route` 已经表达 API 归属，文件和 handler 不重复携带 `_api` / `API` 后缀；同一 package 存在多个 profile 入口时，用 `frontProfile`、`systemProfile` 这类最小业务词区分。组件声明自己的 identity、表资产、生命周期和 service capability，API 通过 `http.Require(c, user.ServiceKey())` 从 linker runtime 解析能力，不在业务侧维护全局 runtime 容器。

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
- `internal/model/user`：数据表模型，包含 `user` 主体表和 `account` 凭据表；业务名直接成为 model/table 节点，不增加项目 prefix，并统一嵌入 `model.Head`。
- `internal/service/user`：登录、资料读取、token/session 和存储流程；service capability key 由 service 自己声明。
- `internal/constant/inspection`：巡检状态等稳定业务词汇，类型自己提供校验、解析、定义集和文本边界。
- `internal/fixture/user`：只承载 example 演示数据；错误回到产生它的 service 或 route/middleware 边界，不塞进 `constant`。
- `internal/component/user`：组件 identity、linker 组件生命周期、表资产和 service capability 挂载。
- `internal/page/*`：Graph Console 页面对象；页面只声明布局、数据展示意图、target 和 Resource，不查询数据库。
- `internal/adapter/console`：把 user/token/session/ACL 业务边界适配为 Console provider，不污染 Graph Console 协议。
- `internal/component/console`：集中装配 Console Component、页面、资源和 provider，不进入 linker core。
- `internal/route/order`、`internal/service/order`：表单 target 的业务 API 与可复用流程，服务端验证始终保留。
- `internal/route/permission`、`internal/service/permission`：multilist 的关系读写接口，最终权限判断仍由后端 middleware 完成。

record-level 权限建议放在具体业务 store 的查询入口处完成。`internal/service/inspection` 用 `TaskAccess` 把 `acl.Access`、`acl.Resource` 和 `RecordRange` 组合在一起：route 只提供当前 application 和 actor，store 在一次查询里同时应用 application scope、业务 filter 和 owner range，避免为了权限判断额外做 N+1 查询或维护 RBAC 关系表。
- `internal/route/inspection`、`internal/model/inspection`、`internal/service/inspection`、`internal/component/inspection`：接近真实业务的列表接口结构，route 负责 HTTP 参数和输出，service/store 负责批量查询和数据范围。
- `internal/model/inspection/archive.go`：外部维护表资产示例，只改业务 model 和 component asset，使用 `postgresql.External()` 避免启动期迁移。
- `internal/component/notification`、`internal/service/notification`、`internal/route/notification`：MQ consumer、cron job、SSE route 和 provider mock 的长生命周期组合；观测 wrapper 由 MQ/cron adapter 统一装配。
- `worker/periodic` 与 `worker/periodic/linker`：稳定固定周期后台循环及其 framework 生命周期装配；`UnhealthyAfter(0)` 可声明只观测、不影响主服务健康的可选任务，日历表达和持久化调度仍由 `scheduler/cron` 承担。
- `observe/metrics/prometheus/linker`、`observe/tracing/opentelemetry/linker`：标准 metrics/tracing 组件；example 不维护平行 recorder capability 或手写 interceptor chain。
- `license/http/gin`：示例只在需要保护的入口显式挂 `license.Gate(gate)`；license 不进入 core，也不默认挂到 server framework。
- `internal/rpc/tts`、`internal/client/tts`、`internal/component/tts`：gRPC server/client 的声明、注册、trace/metrics interceptor 和 capability provider。
- `internal/client/directory`：出站 HTTP typed client 示例，第三方用户目录 API 的业务语义在这里承载，通用 HTTP 执行者来自 `http/client` capability。

组件 identity 必须由组件 package 自己声明，例如 `component/user.ID`。需要依赖该组件时引用这个符号，不把组件 ID 放到 `constant` 或其他公共包里代管。

推荐边界：

- `component/user` 不作为 HTTP controller，不出现成片的 `http.Context` handler、`response.*` 和 route tree 聚合。
- `component/user` 不聚合所有 HTTP handler，也不替 route 维护完整 API 树。
- `component/user` 通过 blank import 纳入 `route/user`，表示启用该组件时才编译这些 route。
- service capability key 由 `service/user` 声明，route 通过 `http.Require(c, user.ServiceKey())` 获取能力，避免 route 反向依赖 component。
- 跨 route middleware 应集中在明确位置；单个 API 只声明 middleware 影响面。

仅在显式注入 `APP_EXAMPLE_USER__SEED_PASSWORD` 时创建本地演示数据：

- 后台账号：`admin`
- 前台手机号：`18558755877`
- 登录密码：使用启动时注入的值；仓库不提供默认值，每个账号写入独立随机盐。

## 运行

默认读取 `config/app.example.yaml`。文件只声明本地 listener、RPC 和 PostgreSQL 示例目标，不包含数据库密码、Nacos 凭据或 token key。启动前至少需要注入：

```bash
export APP_DB_POSTGRESQL__PASSWORD='...'
export APP_EXAMPLE_USER__TOKEN_KEY='至少 32 个字符'
```

需要创建演示账号时，在当前终端无回显读取初始密码：

```bash
read -rsp '演示账号初始密码: ' APP_EXAMPLE_USER__SEED_PASSWORD && echo
export APP_EXAMPLE_USER__SEED_PASSWORD
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

测试文件集中在 `example/` 目录。真实 PostgreSQL example 只在设置 `LINKER_V3_EXAMPLE_PG_PASSWORD` 后运行，host 默认使用 `127.0.0.1` 且可通过同前缀测试变量覆盖；当前环境不可用时会明确 skip。真实 Redis、Nacos 和 RocketMQ 样板分别由 `LINKER_V3_EXAMPLE_REDIS_ADDR`、`LINKER_V3_EXAMPLE_NACOS_HOST`、`LINKER_V3_EXAMPLE_ROCKETMQ_ENDPOINT` 显式开启；凭据只从同前缀环境变量读取。RocketMQ 的 topic/group 是部署资产，测试只引用预建名称并用唯一 message key 识别本次消息。`signal_example_test.go` 会启动真实 server 子进程，分别在 startup 和 running 阶段发送 SIGTERM，验证反向关闭和退出码，不依赖外部 provider。`production_http_example_test.go` 展示 body limit、trusted proxy、health endpoint 和负载中 graceful stop；`docs/Caddyfile` 是部署层终止 TLS 的最小样板。

可选 OAuth 不需要 Linker component，直接在 route 影响面挂 middleware：

```bash
go test ./example -run TestOAuthJWTMiddlewareExample -count=1
```

RocketMQ 5.5 Proxy 测试前先在部署主机预建专用资产；仓库辅助工具只创建不存在的普通 topic/group，不覆盖已有配置：

```bash
ssh rocketmq-host 'bash -s -- topic ensure linker-v3-example' \
  < ../modules/support/rocketmq-admin/rocketmq-admin.sh
ssh rocketmq-host 'bash -s -- group ensure linker-v3-example-consumer' \
  < ../modules/support/rocketmq-admin/rocketmq-admin.sh

LINKER_V3_EXAMPLE_ROCKETMQ_ENDPOINT='<rocketmq-proxy>:8081' \
LINKER_V3_EXAMPLE_ROCKETMQ_TOPIC='linker-v3-example' \
LINKER_V3_EXAMPLE_ROCKETMQ_CONSUMER_GROUP='linker-v3-example-consumer' \
go test ./example -run TestOptionalRocketMQProviderExample -count=1
```

工具默认从当前 Docker 主机的 `rmqbroker` 容器执行 `mqadmin`；容器、NameServer 或集群名不同，通过工具 README 中的环境变量覆盖。不要依赖 `autoCreateTopicEnable` 或 `autoCreateSubscriptionGroup` 代替资产声明。

RocketMQ PushConsumer 停止时会等待 SDK 当前 receive request、long poll 和异步 ack。示例因此显式设置 component `shutdown_timeout=45s` 和 framework `shutdown_timeout=50s`；外层必须大于内层，避免 framework 在组件完成 graceful stop 前先终止关闭漏斗。

Apache RocketMQ Go SDK 5.1.4 在真实 PushConsumer 的 settings/metrics 并发路径存在上游 race detector 报告。`go test -race` 可用于复现和跟踪该边界，但当前不能作为这条真实 provider 示例的通过门禁；项目不复制私有 SDK fork，也不通过全局关闭 race detector 隐藏证据。未启用真实 RocketMQ 环境时，example 自身仍完整执行普通 `-race` 门禁。

```bash
ruby scripts/check-go-baseline.rb
go test ./...
```

Prometheus 可抓取 `GET /metrics`，最小 scrape、OTel Collector 和 Grafana 样板分别位于 `docs/prometheus.yaml`、`docs/otel-collector.yaml`、`docs/grafana-dashboard.json`。默认 trace 使用 memory exporter；部署时通过 `APP_OBSERVE_TRACING__*` 显式切换 OTLP，详见 `docs/observability.md`。dashboard 覆盖 HTTP、gRPC、MQ、cron、component lifecycle、动态配置、fault recovery、terminal、当前异常组件和 notice 投递结果。

推荐先看：

- `docs/scaffold.md`：推荐项目骨架，说明 main、app、component、route、model、service、config、observability 和 example test 的边界。
- `docs/example-policy.md`：说明 example 的定位、外部依赖、测试拆分和未来 submodule 边界。
- `docs/graph-console.md`：Graph Console fixture、前端开发、同进程静态挂载和反向代理四种运行方式。
- `main.go`：保持极薄，只分发 server 启动和 `--plan`。
- `source.go`：只读取配置文件位置和 Nacos bootstrap 参数，按顺序声明 Source。
- `internal/app/app.go`：集中装配 framework、组件和 adapter，不解析业务配置。
- `internal/page/*`：viewer、form、multilist、chart、theme 和 layout 的业务声明。
- `internal/route/order/*.go`、`internal/route/permission/*.go`：一个稳定 API 一个业务文件，route/resource/middleware 和 handler 放在同一个入口重心内。
- `example/graph_example_test.go`：验证 Console Component、登录/session、ACL、页面、静态挂载和业务 API。
- `example/nacos_example_test.go`：验证 YAML seed、Nacos source、HTTP/gRPC registry adapter 和 Plan 里的依赖/capability 表达。
- `example/nacos_provider_example_test.go`：在显式 Nacos 环境中发布独立 data id，经 Source 启动最小 App，并验证清理和关闭。
- `example/redis_example_test.go`：在显式 Redis 环境中验证 component、capability、读写、Plan 资产和 graceful close。
- `example/rocketmq_example_test.go`：在显式 RocketMQ 环境中验证 adapter、producer/consumer capability、消息发送接收、Plan 资产和 graceful close。
- `example/dynamic_config_test.go`：验证 Nacos 完整快照、Live/Restart、desired/active、env 后置覆盖和可恢复拒绝。
- `example/reliability_example_test.go`：验证 DB capability 缺失会在组件初始化期失败，以及 Stop timeout 会返回可判断的 `context.DeadlineExceeded`。
- `example/signal_example_test.go`：验证真实 server 在启动期和运行期收到 SIGTERM 后都执行 graceful close；运行期正常退出，启动期返回可判断的取消原因。
- `example/health_example_test.go`：验证 liveness、readiness、startup 与 framework starting/running 状态一致。
- `example/production_http_example_test.go`：验证标准 Handler root、route 请求体读取期限、body limit、trusted proxy 和在途请求 graceful drain。
- `example/multi_listener_example_test.go`：验证 public/admin 命名 listener、定向 route、middleware 影响面、capability 和 Plan 隔离。
- `docs/Caddyfile`：由 Caddy 终止 TLS，再转发到只监听 loopback 的 linker HTTP server。
- `example/grpc_example_test.go`：验证 gRPC metadata 和 trace id 通过 interceptor 传播。
- `example/http_client_example_test.go`：验证出站 HTTP client linker adapter、typed client、credential、trace hook 和 Plan asset。
- `example/notification_example_test.go`：验证 MQ/cron/SSE lifecycle，并覆盖 HTTP -> MQ mock 的 trace id 贯穿。
- `example/observability_example_test.go`：用本地可控依赖验证 HTTP -> gRPC、HTTP -> MQ、cron span parent/child 和全套指标。
- `example/fault_observability_example_test.go`：用 fake sender 验证 detected/recovering/recovered、notice 和 metrics 闭环。
- `example/feishu_notify_example_test.go`：验证 component 上报 Fault 后由 server 自动发现飞书 Sender，并形成 detected/recovered 通知闭环。
- `example/business_system_test.go`：验证完整业务系统，并覆盖 HTTP -> gRPC typed client 的 trace id 贯穿。
