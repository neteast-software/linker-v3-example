# 运行与联调

这篇文档承接 example 的运行参数和外部依赖边界。推荐结构和编码范式见 `scaffold.md`。

## 本地配置

默认配置链为：

```text
local YAML -> optional Nacos final -> explicit env override
```

后声明 Source 覆盖前者。每个 component 在 Bootstrap 解码并校验自己的 namespace；
`internal/app` 不拆分中心 typed Config。

启动前至少注入：

```bash
export APP_DB_POSTGRESQL__PASSWORD='...'
export APP_EXAMPLE_USER__TOKEN_KEY='至少 32 个字符'
go run .
```

覆盖配置文件位置：

```bash
LINKER_V3_EXAMPLE_CONFIG=config/app.example.yaml \
APP_DB_POSTGRESQL__PASSWORD='...' \
APP_EXAMPLE_USER__TOKEN_KEY='至少 32 个字符' \
go run .
```

加入 Nacos 完整快照：

```bash
LINKER_V3_EXAMPLE_NACOS_DATA_ID='app.yaml' \
LINKER_V3_EXAMPLE_NACOS_HOST='nacos.example.internal' \
APP_DB_POSTGRESQL__PASSWORD='...' \
APP_EXAMPLE_USER__TOKEN_KEY='至少 32 个字符' \
go run .
```

Nacos bootstrap endpoint 和凭据使用
`LINKER_V3_EXAMPLE_NACOS_HOST/PORT/USERNAME/PASSWORD`，不进入业务 namespace。
`Live` 配置先整批 decode/validate，再原子发布给新操作；`Restart` 只更新 desired Setting 并
标记整个服务需要重启。component 不实现 Nacos reload，也不会被 framework 隐式重启。

## 演示账号

只有显式注入 `APP_EXAMPLE_USER__SEED_PASSWORD` 才创建本地演示数据：

- 后台账号：`admin`
- 前台手机号：`18558755877`
- 登录密码：启动时注入的值；每个账号使用独立随机盐。

终端中无回显读取初始密码：

```bash
read -rsp '演示账号初始密码: ' APP_EXAMPLE_USER__SEED_PASSWORD && echo
export APP_EXAMPLE_USER__SEED_PASSWORD
```

## 装配计划

`--plan` 不连接 PostgreSQL，也不启动外部 provider：

```bash
go run . --plan
```

输出包括 mode、components、dependencies、capabilities、Config mode/revision，以及
application、route、gRPC、MQ consumer、cron job、metrics 和 tracing 等 Asset。plan 模式只在
进程内生成临时 token key，不向仓库、Plan 或日志写入凭据。

## 外部依赖

普通测试使用 fake、memory 或本地可控依赖。真实 provider 必须用环境变量显式开启，不可用时
测试明确 skip：

- PostgreSQL：`LINKER_V3_EXAMPLE_PG_PASSWORD`
- Redis：`LINKER_V3_EXAMPLE_REDIS_ADDR`
- Nacos：`LINKER_V3_EXAMPLE_NACOS_HOST`
- RocketMQ：`LINKER_V3_EXAMPLE_ROCKETMQ_ENDPOINT`

RocketMQ 的 topic/group 是部署资产，测试只引用预建名称。部署主机可用 go-module 的受限辅助
工具创建不存在的普通资产：

```bash
ssh rocketmq-host 'bash -s -- topic ensure linker-v3-example' \
  < ../modules/support/rocketmq-admin/rocketmq-admin.sh
ssh rocketmq-host 'bash -s -- group ensure linker-v3-example-consumer' \
  < ../modules/support/rocketmq-admin/rocketmq-admin.sh
```

真实测试还需设置
`LINKER_V3_EXAMPLE_ROCKETMQ_TOPIC` 和
`LINKER_V3_EXAMPLE_ROCKETMQ_CONSUMER_GROUP`。工具不覆盖已有配置，也不依赖 broker 的自动创建
开关。PushConsumer 的 component shutdown timeout 为 `45s`，framework 为 `50s`，外层必须大于
内层。

Apache RocketMQ Go SDK 5.1.4 在真实 PushConsumer settings/metrics 并发路径存在上游 race
报告；普通 example 仍执行完整 `go test -race ./...`，真实 provider 路径保留独立可信边界和推荐
用法说明，不复制私有 SDK fork。

## 脚手架

`scaffold/grpc.yaml` 是 `linker generate grpc` 的最小输入。安装 Linker CLI 后验证生成成果：

```bash
scripts/verify-grpc-scaffold.sh
```

联调源码 CLI 时显式传入可执行文件：

```bash
LINKER_BIN=../linker-v3/cli/linker scripts/verify-grpc-scaffold.sh
```

部署层 TLS 由 `docs/Caddyfile` 终止。Prometheus、OTel Collector 和 Grafana 配置见
`observability.md`。
