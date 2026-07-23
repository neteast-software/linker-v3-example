# Graph Console 运行方式

Graph Console 的 Go 协议和 Component 位于 `go-module/graph/console`，浏览器实现位于
独立的 `graph-console` 仓库。本项目只演示业务装配，不复制协议或前端源码。

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
