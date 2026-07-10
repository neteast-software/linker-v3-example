package example

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	stdhttp "net/http"
	"os"
	"path/filepath"
	"testing"

	yaml "github.com/neteast-software/go-module/config/yaml/linker"
	http "github.com/neteast-software/go-module/http/gin/linker"
	"github.com/neteast-software/go-module/http/gin/response"
	server "github.com/neteast-software/go-module/linker/server"
	linker "github.com/neteast-software/linker/v3"
)

type registryMockSource struct{}

func (registryMockSource) Load(ctx context.Context, boot linker.BootstrapContext) (linker.Setting, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	seed, ok := boot.Seed.Lookup("registry/nacos")
	if !ok || string(seed) != `{"addr":"seed"}` {
		return nil, fmt.Errorf("registry seed = %q, %v", seed, ok)
	}
	return linker.NewSetting(map[linker.Namespace][]byte{
		"cache/redis": []byte(`{"addr":"registry"}`),
		"notice/mock": []byte(`{"enabled":true}`),
	}), nil
}

func TestServerFrameworkLoadsYAMLSource(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "app.yaml")
	if err := os.WriteFile(file, []byte(`
http:
  addr: 127.0.0.1:0
  basePath: api
  recovery: true
  health:
    path: health
`), 0o600); err != nil {
		t.Fatalf("write yaml: %v", err)
	}

	app := server.New(
		server.WithMode(linker.Server),
		server.WithSource(yaml.NewSource(file)),
		server.WithHTTPRoutes(http.GET("hello", func(c *http.Context) {
			response.Data(c, map[string]string{"name": "yaml"})
		})),
	)
	if err := app.Start(context.Background()); err != nil {
		t.Fatalf("start: %v", err)
	}
	t.Cleanup(func() {
		_ = app.Stop(context.Background())
	})

	httpServer, err := linker.RequireCapability(app.App(), linker.NewCapabilityKey[*http.Server](http.ID))
	if err != nil {
		t.Fatalf("http capability: %v", err)
	}

	resp, err := stdhttp.Get("http://" + httpServer.Addr() + "/api/health")
	if err != nil {
		t.Fatalf("get health: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read health body: %v", err)
	}
	if resp.StatusCode != stdhttp.StatusOK || string(body) != "ok" {
		t.Fatalf("unexpected health response: status=%d body=%q", resp.StatusCode, body)
	}

	resp, err = stdhttp.Get("http://" + httpServer.Addr() + "/api/hello")
	if err != nil {
		t.Fatalf("get hello: %v", err)
	}
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read hello body: %v", err)
	}
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("decode hello body: %v body=%q", err, body)
	}
	data, ok := payload["data"].(map[string]any)
	if !ok || data["name"] != "yaml" {
		t.Fatalf("unexpected hello payload: %#v", payload)
	}
}

func TestServerFrameworkAppliesEnvOverrideAfterYAML(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "app.yaml")
	if err := os.WriteFile(file, []byte(`
http:
  addr: 127.0.0.1:0
  basePath: yaml-api
  recovery: true
  health:
    path: yaml-health
`), 0o600); err != nil {
		t.Fatalf("write yaml: %v", err)
	}
	t.Setenv("LINKER_HTTP", `{"addr":"127.0.0.1:0","basePath":"env-api","recovery":true,"health":{"enabled":true,"path":"ready"}}`)

	app := server.New(
		server.WithMode(linker.Server),
		server.WithSource(yaml.NewSource(file), linker.EnvSource{Prefix: "LINKER_"}),
		server.WithHTTPRoutes(http.GET("hello", func(c *http.Context) {
			response.Data(c, map[string]string{"name": "env"})
		})),
	)
	if err := app.Start(context.Background()); err != nil {
		t.Fatalf("start: %v", err)
	}
	t.Cleanup(func() {
		_ = app.Stop(context.Background())
	})

	httpServer, err := linker.RequireCapability(app.App(), linker.NewCapabilityKey[*http.Server](http.ID))
	if err != nil {
		t.Fatalf("http capability: %v", err)
	}
	resp, err := stdhttp.Get("http://" + httpServer.Addr() + "/env-api/ready")
	if err != nil {
		t.Fatalf("get ready: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != stdhttp.StatusOK {
		t.Fatalf("unexpected health status=%d", resp.StatusCode)
	}
}

func TestServerFrameworkLoadsRegistryMockAfterLocalSeed(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "app.yaml")
	if err := os.WriteFile(file, []byte(`
registry/nacos:
  addr: seed
cache/redis:
  addr: local
`), 0o600); err != nil {
		t.Fatalf("write yaml: %v", err)
	}
	t.Setenv("LINKER_CACHE__REDIS", `{"addr":"env"}`)

	app := server.New(
		server.WithMode(linker.Bin),
		server.WithoutHTTP(),
		server.WithSource(yaml.NewSource(file), registryMockSource{}, linker.EnvSource{Prefix: "LINKER_"}),
	)
	if err := app.Start(context.Background()); err != nil {
		t.Fatalf("start: %v", err)
	}
	t.Cleanup(func() {
		_ = app.Close(context.Background())
	})

	cache, ok := app.App().Setting("cache/redis")
	if !ok || string(cache) != `{"addr":"env"}` {
		t.Fatalf("env override setting = %q, %v", cache, ok)
	}
	notice, ok := app.App().Setting("notice/mock")
	if !ok || string(notice) != `{"enabled":true}` {
		t.Fatalf("registry mock setting = %q, %v", notice, ok)
	}
}
