package example

import (
	"context"
	"io"
	stdhttp "net/http"
	"testing"

	http "github.com/neteast-software/go-module/http/gin/linker"
	linker "github.com/neteast-software/linker/v3"
)

type multiListenerRoutes struct {
	assets []linker.Asset
}

func (p *multiListenerRoutes) Identity() linker.ID {
	return "example/multi-listener-routes"
}

func (p *multiListenerRoutes) Assets(context.Context, linker.Runtime) ([]linker.Asset, error) {
	return append([]linker.Asset(nil), p.assets...), nil
}

func TestMultipleHTTPListenerExample(t *testing.T) {
	public := http.Named("public")
	admin := http.Named("admin")
	publicConfig := http.DefaultConfig()
	publicConfig.Addr = "127.0.0.1:0"
	adminConfig := http.DefaultConfig()
	adminConfig.Addr = "127.0.0.1:0"

	routes := &multiListenerRoutes{assets: append(
		public.Routes(http.GET("status", func(c *http.Context) {
			c.String(stdhttp.StatusOK, "public")
		})),
		admin.Routes(http.Group("admin",
			http.Use(func(c *http.Context) {
				c.Header("X-Listener", "admin")
				c.Next()
			}),
			http.GET("status", func(c *http.Context) {
				c.String(stdhttp.StatusOK, "admin")
			}),
		))...,
	)}
	app := linker.New(
		linker.WithMode(linker.Server),
		linker.WithComponents(
			http.New(public, http.WithConfig(publicConfig)),
			http.New(admin, http.WithConfig(adminConfig)),
			routes,
		),
	)
	if err := app.Start(context.Background()); err != nil {
		t.Fatalf("启动多 HTTP listener example: %v", err)
	}
	t.Cleanup(func() {
		if err := app.Stop(context.Background()); err != nil {
			t.Errorf("关闭多 HTTP listener example: %v", err)
		}
	})

	publicServer, err := public.RequireServer(app)
	if err != nil {
		t.Fatalf("获取 public listener: %v", err)
	}
	adminServer, err := admin.RequireServer(app)
	if err != nil {
		t.Fatalf("获取 admin listener: %v", err)
	}
	assertListenerResponse(t, publicServer, "/status", stdhttp.StatusOK, "public", "")
	assertListenerResponse(t, publicServer, "/admin/status", stdhttp.StatusNotFound, "", "")
	assertListenerResponse(t, adminServer, "/admin/status", stdhttp.StatusOK, "admin", "admin")
	assertListenerResponse(t, adminServer, "/status", stdhttp.StatusNotFound, "", "")

	plan := app.Plan()
	if !planHasComponent(plan, public.ID()) || !planHasComponent(plan, admin.ID()) ||
		!planHasTargetedRoute(plan, public.ID(), "/status") ||
		!planHasTargetedRoute(plan, admin.ID(), "/admin/status") {
		t.Fatalf("多 HTTP listener Plan 不完整: %#v", plan)
	}
}

func assertListenerResponse(t *testing.T, server *http.Server, path string, status int, body string, listener string) {
	t.Helper()
	resp, err := stdhttp.Get("http://" + server.Addr() + path)
	if err != nil {
		t.Fatalf("GET %s: %v", path, err)
	}
	defer resp.Body.Close()
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("读取 %s: %v", path, err)
	}
	if resp.StatusCode != status || body != "" && string(content) != body {
		t.Fatalf("GET %s status=%d body=%q，期望 status=%d body=%q", path, resp.StatusCode, content, status, body)
	}
	if resp.Header.Get("X-Listener") != listener {
		t.Fatalf("GET %s listener header=%q，期望 %q", path, resp.Header.Get("X-Listener"), listener)
	}
}

func planHasTargetedRoute(plan linker.Plan, target linker.ID, path string) bool {
	for _, asset := range plan.Assets {
		if asset.Kind == linker.AssetRoute && asset.Target == target && asset.Detail["path"] == path {
			return true
		}
	}
	return false
}
