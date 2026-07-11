# linker v3 项目脚手架

这份文档描述 linker v3 推荐的业务项目骨架。它不是强制目录模板，而是当前 example 验证过的默认起点；真实项目可以按业务重心增减目录，但不建议打散这些边界。

## 项目入口

推荐入口保持极薄：

```text
main.go
source.go
internal/app/app.go
internal/component/<domain>/config.go
```

- `main.go` 只负责 CLI/server/bin 入口分发、中文错误表达和退出码。
- `source.go` 只读取文件位置、注册中心 bootstrap 等最小启动参数，并按顺序声明 Source。
- `internal/app` 负责装配 linker framework、components 和 adapter，只接收 Source，不拆分全局业务 Config。
- 每个 component 或自治 module 在自己的 package 定义、解码和校验 typed config；业务配置不回收到中心 `internal/config`。
- 不在 `main.go` 里拼组件细节，不在业务 route 里读取全局配置。

## 组件

组件只对自己的业务生命周期负责：

```text
internal/component/<domain>/component.go
```

组件文件建议只承载：

- `const ID linker.ID = "example/<domain>"`。
- `func (p *Component) Identity() linker.ID`。
- `Dependencies()`、`Assets()`、`Init()`、`OnMounted()`、`Start()`、`Stop()`。
- `Bootstrap()` 和 `Configs()`，由配置对象自己声明 `Live` 或 `Restart`。
- capability provide 和 component 自己拥有的 table、route、rpc、consumer、job 等资产声明。

约束：

- 组件 identity 必须由组件自己声明，其他组件依赖这个导出符号，不重复手写字符串。
- `*Component` 方法接收者统一使用 `p`，`c` 留给 context 或框架上下文。
- component 不集中塞所有 route handler。
- 如果组件带 HTTP route，用 blank import 让 route 随组件进入编译：`_ "project/internal/route/<domain>"`。
- 组件通过 `Assets()` 声明自己拥有的 table、route、rpc、consumer、job 等资产；framework 和 adapter 负责收集、装配和管理生命周期。
- 组件硬依赖写 `linker.RequireComponent(postgresql.ID)`；软性的启动让位在组件自己的 `Dependencies()` 内写 core 的 `linker.StartAfter(...)`，业务装配层不改写组件依赖。
- `server.WithComponents(...)` 的默认心智模型是声明顺序启动；遇到依赖未就绪的组件会先让位，framework 会回扫并生成最终启动顺序。

反例：

```text
internal/component/msgx/component.go
```

如果一个 component 文件同时出现 `http.Context` handler、`response.*`、route tree 聚合、鉴权 helper 和请求解析，它已经退化成 controller 大文件。正确方向是把 HTTP 入口拆回 `internal/route/<domain>/*_api.go`，把跨 route middleware 放到 `internal/route/middleware` 或 domain middleware 文件，把复用流程留给 `internal/service/<domain>`。

## Route

一个 API 一个 route 文件：

```text
internal/route/<domain>/<action>_api.go
internal/route/<domain>/query.go
internal/route/<domain>/response.go
internal/route/middleware/<name>.go
```

推荐写法：

```go
func init() {
    http.RegisterIn("api/v1/app2",
        http.GET("user/:id/profile", profileAPI).Resource(
            "http.app2.user.profile",
            acl.Scope("app2", 1, "应用二用户资料"),
        ),
    )
}
```

route 的职责：

- 声明路径、method、resource、scope 和 middleware 影响面。
- 解析 HTTP param/query/body。
- 调用 service 或本 route 私有流程。
- 表达 response。

route 不做的事：

- 不维护中心 route tree。
- 不读全局 DB。
- 不把所有 API handler 放到同一个文件。
- 不让 component 替 route 声明 method、path、resource 和 middleware 影响面。
- 不因为 Java 习惯强制拆 DTO/VO；payload、param、response 默认贴近 route，跨入口复用时再沉淀为 service 或稳定对象。

## Model

model 默认表示可持久化资源映射，大多数时候就是 DB model：

```text
internal/model/<domain>/<object>.go
```

推荐：

- 一个对象一个文件，例如 `user.go`、`account.go`、`conversion.go`。
- 默认有 `id`，优先复用约定好的 DB head。
- model 可以承载表名、基础格式化、转换、映射、批处理等资源自身能力。
- 数据库约束尽量保持 GORM 兼容，不优先依赖数据库自定义函数、触发器或外键。

不推荐：

- 把 request payload、query param、response view 默认放进 `model`。
- 用一个 `model.go` 装多个中心不同的表对象。
- 让 model 反向依赖 route、component 或 framework。

## Service

service 只承载稳定复用的业务流程：

```text
internal/service/<domain>/service.go
internal/service/<domain>/store.go
internal/service/<domain>/service_key.go
```

推荐：

- service capability key 由 service package 自己声明。
- route 通过 `http.Require(c, domain.ServiceKey())` 获取能力。
- store/query 把 application scope、actor、resource、record range 和业务 filter 尽量合并到一次查询里。
- route 已经能自然承载的入口过程，不必过早拆进 service。

## Client

出站调用优先通过 typed client 表达业务语义：

```text
internal/client/<domain>/client.go
internal/client/<domain>/<object>.go
```

推荐：

- `go-module/http/client` 承载通用 HTTP 执行能力：timeout、retry、credential、错误映射、trace/log hook 和可测试 transport。
- 业务 typed client 承载第三方 API 的业务对象和方法，例如 `directory.New(api).Badge(ctx, userID)`。
- 需要 linker 生命周期识别时，用 `http/client/linker` 暴露 capability 和 Plan asset。
- route/service 不直接散落 `http.NewRequestWithContext`、`Do`、凭据拼接和原始响应解析。
- 真实外部 provider 测试默认使用 mock server；访问真实服务必须通过显式环境变量开启。

## Config

推荐配置加载顺序：

```text
local YAML -> registry final -> env override
```

要求：

- 示例 YAML 不写真实密码、token、secret。
- 敏感字段通过 env 或部署系统注入。
- Source 只负责完整配置层；component 在 Bootstrap 从 effective Setting 解码并校验自己的 namespace。
- `Live` 配置先准备整批 immutable snapshot，成功后只影响新操作；`Restart` 配置只更新 desired 并标记服务重启。
- 未声明 namespace 默认 `Restart`，但标准 component 应显式声明 owner 和 mode。
- registry source 只提供配置来源，不把业务逻辑塞进配置层。

## Observability

推荐默认具备：

- `GET /metrics` Prometheus scrape 样板。
- HTTP/gRPC/MQ/cron/component runtime 指标。
- trace/request id 在 HTTP、gRPC、MQ、cron、SSE 中贯穿。
- Grafana dashboard JSON 或指标说明。
- audit 表达业务操作事实，fault 表达错误事件，metrics 表达运行指标，trace 表达请求链路。

边界：

- trace id 不作为 metrics label。
- metrics label 必须低基数、低敏感。
- 日志或 OpenTelemetry attribute 通过 `observe/tracing/attribute` 做安全投影，不直接输出 payload、token、手机号、邮箱或错误正文。

## Example Test

推荐测试目录：

```text
example/<domain>_example_test.go
```

测试应覆盖：

- `go run . --plan` 能解释组件、依赖、capability 和 asset。
- 新增 API 不需要维护中心 route tree。
- app 层不 import route 包，component 通过 blank import 控制 route 是否进入编译。
- component 文件不承载 HTTP controller 职责。
- 组件缺少依赖时在 Init 阶段失败。
- graceful stop、context cancel、provider failure、DB store failure 有反馈。
- 真实外部依赖不可用时跳过集成测试，不阻断普通 `go test ./...`。

## 最小新增 API 流程

1. 新增 `internal/route/<domain>/<action>_api.go`，在 `init()` 里声明 route/resource/scope。
2. 如果需要持久化对象，新增 `internal/model/<domain>/<object>.go`。
3. 如果流程会被复用，新增或扩展 `internal/service/<domain>`。
4. 如果这是新业务域，新增 `internal/component/<domain>/component.go` 并声明 identity、依赖、资产和 capability。
5. 在 `internal/app` 的 `server.WithComponents(...)` 挂载组件。
6. 在 `example/<domain>_example_test.go` 增加 route、Plan、业务行为或失败场景测试。
