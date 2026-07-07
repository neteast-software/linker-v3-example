# linker-v3-example

`linker-v3-example` 是 linker v3 的演示业务系统，用来验证 framework、HTTP route、ACL resource、PostgreSQL 生命周期和业务组件声明方式。

## 业务能力

- `POST /system/login`：后台管理员登录。
- `GET /system/profile`：后台管理员信息。
- `POST /user/login`：前台用户登录。
- `GET /api/profile`：返回前台用户信息，包括用户名、头像、邮箱和手机号。
- `GET /api/v1/app2/user/:id/profile`：多层 route 示例，实际访问形如 `/api/v1/app2/user/3/profile`。
- `GET /api/v1/app2/inspection/tasks`：巡检任务列表，演示 application data scope、分页查询和响应白名单。
- `GET /api/v1/app2/notification/events`：SSE 事件入口，演示长连接 route 的局部声明。
- `GET /metrics`：Prometheus scrape 入口，演示 observability 组件、HTTP 指标 middleware 和低基数 label。
- `GET /api/v1/app2/graph/orders`：graph/naive viewer 示例。
- `GET /api/v1/app2/graph/orders/form`：graph/naive form 示例。
- `GET /api/v1/app2/graph/refresh`：graph/naive behavior 示例。
- `example.tts.TTS/Transcribe`：gRPC service，演示 RPC register、typed client provider 和表资产。

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
- `internal/constant/user`：业务错误。
- `internal/component/user`：组件 identity、linker 组件生命周期、表资产和 service capability 挂载。
- `internal/route/inspection`、`internal/model/inspection`、`internal/service/inspection`、`internal/component/inspection`：接近真实业务的列表接口结构，route 负责 HTTP 参数和输出，service/store 负责批量查询和数据范围。
- `internal/component/notification`、`internal/service/notification`、`internal/route/notification`：MQ consumer、cron job、SSE route 和 provider mock 的长生命周期组合。
- `internal/component/observability`、`internal/service/observability`、`internal/route/observability`：Prometheus recorder capability、`/metrics` route 和 HTTP 请求指标 middleware。
- `internal/rpc/tts`、`internal/client/tts`、`internal/component/tts`：gRPC server/client 的声明、注册和 capability provider。

组件 identity 必须由组件 package 自己声明，例如 `component/user.ID`。需要依赖该组件时引用这个符号，不把组件 ID 放到 `constant` 或其他公共包里代管。

推荐边界：

- `component/user` 不聚合所有 HTTP handler，也不替 route 维护完整 API 树。
- `component/user` 通过 blank import 纳入 `route/user`，表示启用该组件时才编译这些 route。
- service capability key 由 `service/user` 声明，route 通过 `http.Require(c, user.ServiceKey())` 获取能力，避免 route 反向依赖 component。
- 跨 route middleware 应集中在明确位置；单个 API 只声明 middleware 影响面。

默认演示数据只作为本地 example fixture：

- 后台账号：`admin`
- 前台手机号：`18558755877`
- 默认密码：`linfunlinfun`

## 运行

默认连接 pi2 PostgreSQL 的局域网地址 `192.168.3.13:5432`，账号为 `neteast`，数据库名为 `linker_v3_example`。数据库密码不写入默认配置，必须通过 `LINKER_V3_EXAMPLE_PG_PASSWORD` 显式提供。

查看 linker 装配计划不需要连接数据库：

```bash
go run . --plan
```

输出会包含 mode、components、dependencies、capabilities 和 route/table/RPC/job 等 assets。缺少 `LINKER_V3_EXAMPLE_PG_PASSWORD` 时，`--plan` 只使用本地占位值构建计划，不会启动 PostgreSQL component。

```bash
go run .
```

本地需要覆盖数据库地址时：

```bash
LINKER_V3_EXAMPLE_PG_HOST=127.0.0.1 LINKER_V3_EXAMPLE_PG_PASSWORD=... go run .
```

配置源推荐顺序是 `local seed -> registry final -> env override`。`example/server_yaml_test.go` 用 `registryMockSource` 演示 Nacos 类注册中心 source 如何读取本地 seed，再由环境变量覆盖最终配置。

## Example

测试文件集中在 `example/` 目录。真实 PostgreSQL example 会尝试连接 `192.168.3.13`，如果当前环境无法连接，会跳过该集成用例。

```bash
go test ./...
```

推荐先看：

- `main.go`：保持极薄，只分发 server 启动和 `--plan`。
- `internal/app/app.go`：集中装配 framework、组件、配置源和 adapter。
- `internal/route/graph/*_api.go`：一个 API 一个文件，route/resource/middleware 和 handler 放在同一个入口重心内。
- `example/graph_example_test.go`：验证 graph route、Plan asset 和 renderer capability。
