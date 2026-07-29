# linker-v3-example

`linker-v3-example` 是 Linker v3 的演示业务系统和集成验证仓库。它用一套可运行代码展示
server framework、自治能力、HTTP route、ACL、PostgreSQL、RPC、MQ、动态配置、可观测性和
graceful shutdown 如何闭环，不把演示业务模型当成公司统一模型。

当前工具链基线为 Go `1.26.5`，已发布 framework 基线为 Linker `v3.6.0`。默认分支只使用
远端可追溯版本，不依赖本地 `replace`；候选契约的源码联调不得冒充正式消费闭环。精确依赖以
`go.mod` 为准。

## 最短路径

只查看组件、依赖、capability、配置和 Asset 装配计划，不连接外部资源：

```bash
go run . --plan
```

启动本地 YAML server：

```bash
APP_DB_POSTGRESQL__PASSWORD='...' \
APP_EXAMPLE_USER__TOKEN_KEY='至少 32 个字符' \
go run .
```

默认读取 `config/app.example.yaml`。仓库不提供数据库密码、token key、Nacos 凭据或演示账号
默认密码；完整启动、Nacos 覆盖、seed 和外部依赖变量见
[`docs/run.md`](docs/run.md)。

## 覆盖范围

- framework：`New / Use / Run`、依赖拓扑、typed capability、Asset、Plan 和反向关闭漏斗。
- HTTP：能力局部 `RouteSet`、文件级 API 声明、middleware 影响面、ACL Resource 和统一 response。
- 数据：PostgreSQL 生命周期、GORM 对象、`model.Head`、数据范围和棕地表显式边界。
- 服务：typed gRPC client/server、出站 HTTP client、Redis、Nacos、RocketMQ、cron、Worker 和 SSE。
- 工作背景：health、Prometheus、OpenTelemetry、audit、fault、notice、license 和 outbox。
- Graph Console：登录、session、菜单、权限、viewer、form、multilist、chart、theme 和 layout。

新建 server 的关系型数据库默认推荐 PostgreSQL。数据库只作为可替换仓库，业务关系、规则、
权限和流程在所属 Go 能力中自治；新能力不使用外键、自建函数、存储过程、触发器或数据库扩展
承载业务语义。

## 能力入口

```text
main.go                         # 极薄的 server/plan 分发
source.go                       # 按顺序声明配置 Source
internal/app/app.go             # framework 与能力适配层装配
internal/<capability>/          # 业务能力语义入口
internal/<capability>/http/     # 能力的 HTTP 适配
internal/<capability>/linker/   # 能力的生命周期适配
example/                        # 集成与失败路径验证
```

能力按业务纵向组织，不建立全局 `route`、`model`、`service`、`component` 或 `client` 横向层。
package 是 Go 最天然的语义入口；末级目录只有在自身就是调用能力时才适合作为 package 名。

## 阅读地图

- [`docs/scaffold.md`](docs/scaffold.md)：推荐目录、依赖方向、组件、route、持久化对象、配置和验证。
- [`docs/run.md`](docs/run.md)：本地运行、Nacos 覆盖、演示账号、外部依赖和生成式脚手架验证。
- [`docs/example-policy.md`](docs/example-policy.md)：example 边界、source-first 状态和测试拆分。
- [`docs/observability.md`](docs/observability.md)：Prometheus、OpenTelemetry、Grafana 和故障通知。
- [`docs/graph-console.md`](docs/graph-console.md)：Graph Console fixture、静态挂载和前端联调。

现代 Go 能力从 Linker
[`GO.md`](https://github.com/neteast-software/linker/blob/v3/GO.md) 进入；Linker 的底层概念和
框架能力地图从其根 `README.md` 进入。阅读路径保持为 `README -> topic -> code`，不再从顶层
一次性展开所有接口和测试文件。
