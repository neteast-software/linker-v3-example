# 可观测性样板

example 通过 `server.WithMetrics(prometheus.New())` 装配 Prometheus recorder、framework observer 和 `GET /metrics`，通过独立 OpenTelemetry component 为 HTTP、gRPC、MQ 和 cron 提供统一 trace。

默认配置使用有界 memory exporter，不访问远端 collector：

```yaml
observe/tracing:
  mode: memory
  service: linker-v3-example
```

部署时由最后声明的 env Source 显式切换到 OTLP：

```bash
export APP_OBSERVE_TRACING__MODE=otlp
export APP_OBSERVE_TRACING__ENDPOINT=127.0.0.1:4317
export APP_OBSERVE_TRACING__PROTOCOL=grpc
export APP_OBSERVE_TRACING__INSECURE=true
go run .
```

Collector 最小配置见 `otel-collector.yaml`。Prometheus 使用 `prometheus.yaml` 抓取 `/metrics`，Grafana 可导入 `grafana-dashboard.json`。

指标边界：

- HTTP 只使用 method、route template、status。
- gRPC 只使用 method、rpc type、code。
- MQ/cron 只使用稳定 consumer/topic/job/status。
- component lifecycle 只使用 component、stage、result。
- fault 使用 component/state，notice 投递结果只使用 severity/state/status。
- trace id、request id、user id、原始 path、payload、错误正文和凭据不进入 metrics label。

`example/observability_example_test.go` 使用本地 gRPC、内存 consumer、cron、Prometheus 和 OTel memory exporter 验证 HTTP -> gRPC、HTTP -> MQ、cron 的 span 关系。`example/fault_observability_example_test.go` 使用 fake sender 验证 detected -> recovering -> recovered、notice 和 metrics 闭环。

飞书只作为 framework fault notice sender 显式接入。结构见 `feishu-notice.example.yaml`；默认
example 不加载该文件，也不内置 app identity、secret 或接收人。

## Gateway

Gateway profile 使用同一 Prometheus recorder 和 OpenTelemetry provider，通过
`gatewaymetrics.Observe(metrics)` 与 `gatewaytracing.Observe(tracing)` 接入请求漏斗。
稳定观测面包括：

- `gateway_requests_total`、请求耗时、in-flight 和稳定 upstream failure 分类。
- `gateway_plan_reloads_total` 与最后正常 Plan 保留结果。
- 以 route ID 命名的 server span 和标准 trace context 传播。
- `accesslog.New` 产生的脱敏 `operateLog` 记录、有界队列状态和投递失败反馈。

Gateway metrics 不使用原始 path、upstream 地址、用户、trace ID 或错误正文作为 label；
trace 和 accesslog 不记录 token、cookie、证书、完整 query、body 或凭据。管理探针不经过
业务 Route，也不会制造 upstream 请求。默认 scrape 地址为
`http://127.0.0.1:8820/metrics`，只绑定 loopback；部署时由受控采集链路访问，不把指标端点
暴露为公网 Gateway Route。
