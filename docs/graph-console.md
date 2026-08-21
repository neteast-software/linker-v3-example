# Graph Console 运行方式

Graph Console 的 Go 协议位于 `go-module/graph/console`，独立 Linker adapter 位于
`go-module/graph/console/linker`，浏览器实现位于独立的 `graph-console` 仓库。本项目只演示业务
装配，不复制协议或前端源码。业务 page/resource 由 `internal/<capability>/linker` 使用 typed Asset
自治声明；`internal/console/linker` 只保留 dashboard 和系统级 menu/session/provider 装配。

本项目显式消费 `graph.console/v2`：GraphPage 使用完整 PageURI，ClientPage 使用业务前端 catalog
拥有的唯一 identity，IframePage 使用经过 entry allowlist 允许的 URL。唯一系统菜单位于
`internal/console/menu.yaml`，其 schema 为 `graph.console.menu.v2`。

历史系统缺少 `protocol` 时，由业务前端 Shell 的版本选择层命中 `naive-v2`；显式
`graph.console/v2` 命中 `graph-console-v2`。这是一轮确定性选择，不是请求失败后的 fallback。
本 Go 示例始终输出显式 v2，不在后端复制历史协议或 renderer runtime。

## 1. Fixture

前端开发默认使用 canonical fixture，不需要启动 Go 服务：

```bash
cd ../graph-console
pnpm dev
```

该模式只验证 renderer 和交互，不替代真实登录、session、ACL 与业务 API smoke。

## 2. 前后端开发

先启动本项目，再让 Vite 把 `/console` 和 `/api` 转发到 Go 服务：

```bash
go run .
```

```bash
cd ../graph-console
VITE_CONSOLE_SOURCE=http \
VITE_CONSOLE_BACKEND=http://127.0.0.1:8080 \
pnpm dev
```

后端 route 是权限和服务端验证的最终边界，前端 Access 投影只用于隐藏、禁用和解释
界面能力。

## 3. 同进程静态挂载

生产构建得到 `graph-console/dist` 后，把静态目录转换为 `fs.FS`，再作为
`graphconsole.WithStatic(files)` 传给 `internal/console/linker` 的 `console.New`。Component 使用 HTTP
公共 route 挂载 `/assets` 和 SPA 入口，API 继续位于 `/console`。

`example/graph_example_test.go` 使用 `fstest.MapFS` 验证了相同边界。真实项目可以用
`embed.FS`，但前端产物所有权仍属于 graph-console，不进入 linker core。

## 4. 反向代理

前端静态产物和 Go 服务可以独立部署，由 Caddy 统一入口：

```caddyfile
console.example.internal {
    handle /console/* {
        reverse_proxy 127.0.0.1:8080
    }
    handle /api/* {
        reverse_proxy 127.0.0.1:8080
    }
    handle {
        root * /srv/graph-console
        try_files {path} /index.html
        file_server
    }
}
```

TLS、静态缓存和公网入口属于部署层；linker 与 Console Component 保持中性。
