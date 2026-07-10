package example

import (
	"context"
	"fmt"
	"io"
	stdhttp "net/http"
	"os"
	"path/filepath"
	"testing"

	yaml "github.com/neteast-software/go-module/config/yaml/linker"
	http "github.com/neteast-software/go-module/http/gin/linker"
	httpregistry "github.com/neteast-software/go-module/http/gin/registry/nacos/linker"
	server "github.com/neteast-software/go-module/linker/server"
	registrynacos "github.com/neteast-software/go-module/registry/nacos/linker"
	grpcregistry "github.com/neteast-software/go-module/rpc/grpc/registry/nacos/linker"
	linker "github.com/neteast-software/linker/v3"
	nacoskit "github.com/neteast-software/nacos-kit"
)

type fakeNacosNaming struct{}

func (fakeNacosNaming) Register(context.Context, nacoskit.Instance) error {
	return nil
}

func (fakeNacosNaming) Discover(context.Context, string, ...nacoskit.DiscoverOption) ([]nacoskit.Instance, error) {
	return nil, nil
}

func (fakeNacosNaming) DiscoverHealthy(context.Context, string, ...nacoskit.DiscoverOption) ([]nacoskit.Instance, error) {
	return nil, nil
}

func (fakeNacosNaming) Offline(context.Context, nacoskit.Instance) error {
	return nil
}

func (fakeNacosNaming) Deregister(context.Context, nacoskit.Instance) error {
	return nil
}

func (fakeNacosNaming) Close() error {
	return nil
}

type fakeNacosComponent struct {
	naming nacoskit.Naming
}

func (p fakeNacosComponent) Identity() linker.ID {
	return registrynacos.ID
}

func (p fakeNacosComponent) OnMounted(_ context.Context, runtime linker.Runtime) error {
	return linker.Provide(runtime, linker.NewCapabilityKey[nacoskit.Naming]("registry/nacos/naming"), p.naming)
}

func TestServerFrameworkLoadsNacosSourceAfterLocalSeed(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "app.yaml")
	if err := os.WriteFile(file, []byte(`
registry/nacos:
  host: seed.nacos.local
  port: 8848
http:
  addr: 127.0.0.1:0
  base_path: local
`), 0o600); err != nil {
		t.Fatalf("write yaml: %v", err)
	}

	source := registrynacos.NewSource(
		registrynacos.SourceConfig{
			DataID:    "linker-example.yaml",
			Group:     "LINKER",
			Namespace: "http",
		},
		registrynacos.WithGetter(func(_ context.Context, config registrynacos.SourceConfig) ([]byte, error) {
			if config.Nacos.Host != "seed.nacos.local" || config.Nacos.Port != 8848 {
				return nil, fmt.Errorf("nacos seed = %+v", config.Nacos)
			}
			return []byte(`{"addr":"127.0.0.1:0","base_path":"nacos","health":{"enabled":true,"path":"ready"}}`), nil
		}),
	)

	app := server.New(
		server.WithSource(yaml.NewSource(file), source),
		server.WithHTTPRoutes(http.GET("hello", func(c *http.Context) {
			c.String(stdhttp.StatusOK, "nacos")
		})),
	)
	if err := app.Start(context.Background()); err != nil {
		t.Fatalf("start: %v", err)
	}
	t.Cleanup(func() {
		_ = app.Stop(context.Background())
	})

	httpServer, err := linker.RequireCapability(app, http.ServerKey())
	if err != nil {
		t.Fatalf("http capability: %v", err)
	}
	resp, err := stdhttp.Get("http://" + httpServer.Addr() + "/nacos/ready")
	if err != nil {
		t.Fatalf("get ready: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	if resp.StatusCode != stdhttp.StatusOK || string(body) != "ok" {
		t.Fatalf("ready response status=%d body=%q", resp.StatusCode, body)
	}
}

func TestCoreBinPlansNacosHTTPAndGRPCRegistries(t *testing.T) {
	app := linker.New(
		linker.WithMode(linker.Bin),
		linker.WithComponents(
			fakeNacosComponent{naming: fakeNacosNaming{}},
			httpregistry.New(),
			grpcregistry.New(),
		),
	)

	plan := app.Plan()
	if planOrder(plan, registrynacos.ID) >= planOrder(plan, httpregistry.ID) ||
		planOrder(plan, registrynacos.ID) >= planOrder(plan, grpcregistry.ID) {
		t.Fatalf("registry order = %#v", plan.Components)
	}
	if err := app.Start(context.Background()); err != nil {
		t.Fatalf("start: %v", err)
	}
	t.Cleanup(func() {
		_ = app.Stop(context.Background())
	})

	plan = app.Plan()
	assertPlanCapability(t, plan, "registry/nacos/naming")
	assertPlanCapability(t, plan, "http/registry")
	assertPlanCapability(t, plan, "rpc/grpc/registry")
}

func assertPlanCapability(t *testing.T, plan linker.Plan, id linker.ID) {
	t.Helper()
	for _, capability := range plan.Capabilities {
		if capability.ID == id {
			return
		}
	}
	t.Fatalf("plan capability %s not found in %#v", id, plan.Capabilities)
}
