# Linker v3 项目骨架

这份文档描述 Linker v3 业务项目当前验证过的推荐结构。它约束依赖方向和语义入口，不要求所有项目拥有相同目录。

## 入口

```text
main.go
source.go
internal/app/app.go
internal/<capability>/
```

- `main.go` 只负责 CLI/server/bin 分发、中文错误表达和退出码。
- `source.go` 只读取配置文件位置和注册中心 bootstrap 参数，按声明顺序返回 Source。
- `internal/app` 只装配 framework 和能力适配层，不解析业务配置、不直接调用业务对象。
- `internal/<capability>` 是业务能力的语义入口，例如 `user`、`order`、`inspection`。
- 不建立全局 `component`、`route`、`model`、`service`、`constant`、`adapter` 或 `client` 技术分层。

`app` 只看见能力的装配入口：

```go
users := user.NewComponent()

return server.New(
    server.Config(sources...),
    server.WithComponents(
        postgresql.New(),
        users,
        order.New(),
        console.New(users.Service()),
    ),
)
```

业务调用方无需理解组件内部如何持久化、提供 route 或注册 capability。

## 能力

一个能力先形成自己的业务闭环，再按需增加技术适配层：

```text
internal/user/
├── user.go
├── account.go
├── auth.go
├── store.go
├── http/
│   ├── profile.go
│   └── user_login.go
└── linker/
    ├── component.go
    └── config.go
```

依赖方向固定为：

```text
app -> capability/linker -> capability
                     └── blank import capability/http
capability/http -> capability
capability -> autonomous modules
```

能力根 package 不反向依赖 `app`、`http` 或 `linker`。末级目录只有在它本身就是调用能力时才适合作为 package 名；`http`、`linker` 只是技术路径，因此其中的 package 仍使用 `user`、`order` 等业务名。

一个文件保持一个重心。对象、流程或适配器形成稳定中心后单独成文件；很轻的聚合关系可以同文件保留。不要为了对象化制造没有状态、策略、生命周期、反馈或审计闭环的包装。

## Linker 适配

需要 framework 管理生命周期时，在能力内部增加：

```text
internal/<capability>/linker/component.go
```

组件适配层可以承载：

- `const ID linker.ID` 和 `Identity()`。
- `Dependencies()`、`Assets()`、`Init()`、`OnMounted()`、`Start()`、`Stop()`。
- `Bootstrap()` 和 `Configs()`。
- 能力的 table、RPC、consumer、job 等 Asset。
- capability provide。

约束：

- 组件 identity 由组件自己声明，依赖方引用该符号，不重复字符串。
- `*Component` 方法接收者统一使用 `p`。
- 组件不承载 `http.Context` handler、response 或 route tree。
- 有 HTTP 入口时，由组件 blank import `internal/<capability>/http`，使 route 服从组件编译边界。
- 硬依赖使用 `linker.RequireComponent`；软启动次序由组件自己的 `Dependencies()` 声明。
- `server.WithComponents` 按业务顺序声明。framework 遇到未满足依赖会让位、回扫并计算最终拓扑。

稳定周期任务优先使用自治 `worker/periodic`，再通过它的 Linker 适配层进入生命周期。需要日历表达、持久化或补跑时使用 `scheduler/cron`；单次并行处理使用 `worker/parallel`。

## HTTP

HTTP 是能力内部的协议适配：

```text
internal/<capability>/http/<action>.go
internal/<capability>/http/query.go
internal/<capability>/http/response.go
internal/access/<policy>.go
```

一个稳定 API 一个业务文件，并在自己的 `init()` 中注册：

```go
func init() {
    http.RegisterIn("api/v1/app2",
        http.GET("user/:id/profile", profile).Resource(
            "http.app2.user.profile",
            acl.Scope("app2", 1, "应用二用户资料"),
        ),
    )
}
```

route 文件集中表达 path、method、resource、scope、middleware 影响面、请求结构、响应结构和该入口独有的短流程。跨入口复用的能力回到能力根 package。

要求：

- 不维护中心 route tree。
- 不让 Linker 组件代替 HTTP 文件声明入口。
- 文件和 handler 不重复 `_api` 或 `API` 后缀。
- payload、param、response 默认贴近入口，不强制拆 DTO/VO。
- 跨 route 的访问策略统一放在 `internal/access`，具体 route 只声明影响面。
- route 通过 `http.Require(c, user.ServiceKey())` 获取能力，不接触全局 runtime 容器。

## 持久化对象

model 是可持久化资源映射的概念，不是必须建立的全局目录。对象直接位于能力 package：

```text
internal/user/user.go
internal/user/account.go
internal/inspection/task.go
```

当前新建 server 的关系型数据库推荐 PostgreSQL，通过 `db/postgresql` 接入。普通业务资产嵌入 `db/gorm/model.Head`，使用 `uint64` 自增 ID、创建时间和更新时间；UUID 不作为默认主键。没有独立身份和生命周期的纯关系记录可以使用复合主键或唯一复合索引。

推荐：

- 一个对象一个文件。
- `user` 已表达业务归属后，字段和方法使用 `Email`、`Phone`、`Validate` 等相对语义。
- 表名使用单数业务名，不增加 `sys_`、`tbl_`、项目名或 `model_` prefix。
- 对象可以自然生长格式化、转换、映射和批处理等自治能力。
- 查询优先使用 GORM 链式 API；只有明确超出其表达范围时才使用底层 SQL。
- 约束保持 GORM 可迁移，不优先依赖数据库函数、触发器或外键。
- 数据范围、application、actor 和业务 filter 尽量合并成一次查询，避免 N+1。

棕地 SQL-only 表可以通过 PostgreSQL adapter 的 `ExternalTable` 或 `ReadOnlyTable` 声明，不为门禁制造空 model。sqlc 只作为历史项目兼容路径，新项目统一采用 GORM。

## 可复用流程

可复用流程同样属于能力根 package：

```text
internal/user/auth.go
internal/user/store.go
internal/user/service.go
```

`service.go` 是可选文件，不是强制层。route 已经能自然闭环的入口过程留在 route；跨入口复用、具有稳定策略或生命周期的流程再进入能力根。

能力 key 由能力根 package 自己声明，适配层只负责 `Provide`，调用方通过语义入口解析。避免 `UserServiceImpl`、`IUserRepository`、`UserModel` 等重复业务归属或技术路径的名称。

## 出站调用

通用连接、重试、凭据、错误映射和观测由 go-module 提供；业务项目只保留 typed client 的业务语义：

```text
internal/directory/client.go
internal/tts/client/client.go
```

- `directory.New(api).Badge(ctx, userID)` 这类调用应直接表达业务。
- 当 `client` 本身是稳定能力时可以作为 package；当它只是某能力的技术适配路径时，package 仍使用该能力名。
- 需要 Linker 管理连接生命周期时，通过 module 自带的 Linker adapter 提供 capability 和 Plan Asset。
- 不在 route 或业务流程中散落 `http.NewRequestWithContext`、原始 `Do`、凭据拼接和响应解码。
- 真实 provider 测试必须通过显式环境变量开启，普通测试使用本地 mock。

## 配置

默认配置链：

```text
local YAML -> registry final -> env override
```

- 后声明 Source 覆盖前者，启动日志解释覆盖来源。
- YAML 不写真实密码、token 或 secret。
- component 在 `Bootstrap` 中解码并校验自己的 namespace。
- `Live` 先完成整批 decode/validate，再原子发布给新操作。
- `Restart` 只更新 desired 配置并标记服务需要重启。
- registry 只提供配置来源，不承载业务逻辑。
- 不建立中心 `internal/config` 拆分所有组件配置。

## 观测

server 默认应具备 Prometheus metrics、OpenTelemetry trace、业务 audit 和故障反馈：

- metrics 使用低基数、低敏感 label，不把 trace ID 作为 label。
- trace/request ID 在 HTTP、gRPC、MQ、cron 和 SSE 中传播。
- audit 表达操作事实，fault 表达异常事件，metrics 表达运行状态，trace 表达调用链。
- attribute 使用安全投影，不直接输出 payload、token、手机号、邮箱或错误正文。
- app 使用 `server.WithMetrics(prometheus.New())` 等 framework 语义入口，不手写 SDK interceptor chain。

## 验证

测试集中在 `example/`，至少覆盖：

- `go run . --plan` 能解释组件、依赖、capability 和 Asset。
- 新 API 无需修改中心 route tree。
- `app` 只导入能力适配层。
- 能力根不反向依赖 `app/http/linker`。
- Linker 适配层不承载 HTTP controller。
- 依赖缺失、provider failure、DB failure、context cancel 和 graceful stop 有明确反馈。
- 外部依赖未显式配置时跳过集成测试，不阻断 `go test ./...`。

新增 API 的最短路径：

1. 在 `internal/<capability>/http/<action>.go` 自注册 route。
2. 需要持久化对象时，在能力根新增对象文件。
3. 流程需要复用时，在能力根形成稳定方法或对象。
4. 新能力需要生命周期时，增加 `internal/<capability>/linker/component.go`。
5. 在 `internal/app` 挂载该能力的 Linker 适配入口。
6. 在 `example/` 验证 Plan、业务行为和失败路径。
