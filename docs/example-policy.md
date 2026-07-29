# example 定位

`linker-v3-example` 是 linker v3 的演示项目和集成验证项目。

它负责展示推荐用法和验证框架闭环，不负责成为真实业务模板仓库，也不负责承载所有业务复杂度。

## 使用边界

- 可以用来观察 linker v3 的 server framework、component、route、model、service、config、observability 和 lifecycle 组合方式。
- 可以作为新能力进入 linker v3 前后的最小运行样板。
- 不把真实业务的全部流程搬进 example。
- 不把 example 的业务命名、账号、表结构当成公司统一业务模型。
- 不在 example 里隐藏外部依赖；真实依赖必须可跳过或由环境变量显式开启。

## 新能力进入规则

每个对业务开发者可见的新能力，都先判断是否需要 example：

- 只影响纯工具函数，且 README 和单测已经足够时，可以不进 example。
- 影响 component lifecycle、route、config、observability、DB、RPC、MQ、cron、周期 Worker、SSE、ACL、license、outbox 或 graph 这类框架使用体验时，应补一个最小示例。
- 示例只覆盖最小闭环，不追求真实业务完整度。
- 示例必须能解释 Plan、请求路径、失败反馈或资源边界中的至少一个。

## 依赖模式

example 只有两种可解释状态：

- `released`：默认主路径只依赖远端可解析版本，没有本地 `replace`。
- `source-first`：开发分支为尚未发布的明确候选提供集成证据；`replace` 只指向同级
  go-module checkout，并在对应 tag 和远端安装验证成立后删除。

`replace` 不是 Linker 的安装方式，也不是业务项目的推荐依赖方式。不得用它覆盖已经发布且可组合的
版本，不得使用 `v0.0.0` 冒充候选版本，也不得把开发者本机绝对路径提交进仓库。

当前默认分支属于 `released` 状态：P18 provider 使用协调发布的 `v0.4.0`，MQ 与 Cron
依赖闭包使用修正版 `v0.4.1`，原有 19 条本地 `replace` 已全部删除。每轮 provider 发布后，
仍须使用隔离缓存、`GOWORK=off` 和无 replace 的 `go mod download`、测试、race、vet、QADC
与漏洞扫描复验默认主路径。

## 外部依赖

默认优先使用 mock、memory、fake 或本地可控依赖：

- Nacos 使用 fake getter 演示 source 和 registry adapter。
- MQ 使用 mock consumer/provider 演示生命周期和 trace。
- tracing 默认使用有界 memory exporter；只有显式 OTLP 配置才访问 collector。
- gRPC typed client 使用本地可控地址和测试 discovery。
- 出站 HTTP typed client 使用本地 mock server 演示 credential、trace hook、错误反馈和 Plan asset。
- 外部通知、地图、短信等 provider 不在默认测试中访问真实服务。

真实 PostgreSQL、Nacos、Redis、MQ 或外部 provider 必须通过环境变量显式开启。无法连接真实依赖时，集成测试应跳过，不阻断普通 `go test ./...`。

## 测试拆分

测试按场景拆分在 `example/` 目录：

- `http_example_test.go`：业务 HTTP、response、license 和 Prometheus。
- `gateway_example_test.go`：Gateway 代理、声明热更新、管理探针、指标、trace 和 graceful stop。
- `business_system_test.go`：接近真实业务系统的组合启动、登录、DB、gRPC、MQ、cron。
- `graph_example_test.go`：Graph Console Component、登录/session、ACL、页面、静态挂载和业务 API。
- `notification_example_test.go`：MQ、cron、SSE 和 HTTP -> MQ trace。
- `rocketmq_example_test.go`：显式环境下的 RocketMQ adapter、capability、Plan、发送接收和 graceful stop。
- `grpc_example_test.go`：gRPC metadata 和 trace。
- `outbox_example_test.go`：delivery 和 dead-letter。
- `reliability_example_test.go`：依赖缺失、Stop timeout 和失败反馈。
- `signal_example_test.go`：真实 server 子进程的 startup/runtime SIGTERM 和 graceful close。
- `health_example_test.go`：三类健康端点映射 framework 生命周期。
- `production_http_example_test.go`：HTTP body/proxy/负载中 graceful 边界。
- `Caddyfile`：部署层 TLS 终止样板，linker 不管理证书。
- `nacos_example_test.go`：YAML seed、Nacos source 和 registry adapter。
- `periodic_worker_example_test.go`：固定周期 Worker 的自治声明、Plan Asset、capability 与 graceful stop。

如果单个测试文件开始承载多个不相干流程，优先按场景拆分，不把所有能力塞回一个大测试。

## 与 linker 仓库的关系

当前 example 保持独立仓库和独立 `go.mod`。

如果未来纳入 linker 仓库，推荐作为 git submodule：

```text
examples/linker-v3-example
```

linker 仓库只记录 submodule pointer 和说明，不复制 example 源码，不让 example 变成 linker core 的一部分。
