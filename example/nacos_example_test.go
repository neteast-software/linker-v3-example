package example

import (
	"context"
	"fmt"
	"io"
	stdhttp "net/http"
	"os"
	"path/filepath"
	"sync"
	"testing"

	yaml "github.com/neteast-software/go-module/config/yaml/linker"
	http "github.com/neteast-software/go-module/http/gin/linker"
	server "github.com/neteast-software/go-module/linker/server"
	registrynacos "github.com/neteast-software/go-module/registry/nacos/linker"
	service "github.com/neteast-software/go-module/registry/service"
	serviceregistry "github.com/neteast-software/go-module/registry/service/nacos/linker"
	rpc "github.com/neteast-software/go-module/rpc/grpc/linker"
	linker "github.com/neteast-software/linker/v3"
	nacoskit "github.com/neteast-software/nacos-kit"
)

type fakeNacosNaming struct{}

type recordingNacosNaming struct {
	mu           sync.Mutex
	registered   []nacoskit.Instance
	deregistered []nacoskit.Instance
}

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

func (p *recordingNacosNaming) Register(_ context.Context, instance nacoskit.Instance) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.registered = append(p.registered, instance)
	return nil
}

func (p *recordingNacosNaming) Discover(context.Context, string, ...nacoskit.DiscoverOption) ([]nacoskit.Instance, error) {
	return nil, nil
}

func (p *recordingNacosNaming) DiscoverHealthy(context.Context, string, ...nacoskit.DiscoverOption) ([]nacoskit.Instance, error) {
	return nil, nil
}

func (p *recordingNacosNaming) Offline(context.Context, nacoskit.Instance) error {
	return nil
}

func (p *recordingNacosNaming) Deregister(_ context.Context, instance nacoskit.Instance) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.deregistered = append(p.deregistered, instance)
	return nil
}

func (p *recordingNacosNaming) Close() error {
	return nil
}

func (p *recordingNacosNaming) counts() (int, int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return len(p.registered), len(p.deregistered)
}

type fakeNacosComponent struct {
	naming nacoskit.Naming
}

func (p fakeNacosComponent) Identity() linker.ID {
	return registrynacos.ID
}

func (p fakeNacosComponent) OnMounted(_ context.Context, runtime linker.Runtime) error {
	return linker.Provide(runtime, registrynacos.NamingKey(), p.naming)
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

	source := registrynacos.Config(
		"linker-example.yaml",
		registrynacos.Group("LINKER"),
		registrynacos.Fetch(func(_ context.Context, request registrynacos.Request) ([]byte, error) {
			if request.Client.Host != "seed.nacos.local" || request.Client.Port != 8848 {
				return nil, fmt.Errorf("nacos seed = %+v", request.Client)
			}
			return []byte("http:\n  addr: 127.0.0.1:0\n  base_path: nacos\n  health:\n    readiness:\n      enabled: true\n      path: ready\n"), nil
		}),
	)

	app := server.New(
		server.Config(yaml.File(file), source),
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

	httpServer, err := http.RequireServer(app)
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

func TestCoreBinPlansNacosServiceRegistry(t *testing.T) {
	app := linker.New(
		linker.WithMode(linker.Bin),
		linker.WithComponents(
			fakeNacosComponent{naming: fakeNacosNaming{}},
			serviceregistry.New(),
		),
	)

	plan := app.Plan()
	if planOrder(plan, registrynacos.ID) >= planOrder(plan, "registry/service") {
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
	assertPlanCapability(t, plan, "registry/service")
}

func TestOneServiceRegistryProvidesHTTPAndGRPC(t *testing.T) {
	naming := &recordingNacosNaming{}
	httpConfig := http.DefaultConfig()
	httpConfig.Addr = "127.0.0.1:0"
	httpConfig.Registry = service.Config{Policy: service.FailFast, Service: "example-http"}
	rpcConfig := rpc.ServerConfig{
		Addr:     "127.0.0.1:0",
		Registry: service.Config{Policy: service.FailFast, Service: "example-rpc"},
	}
	app := linker.New(
		linker.WithMode(linker.Server),
		linker.WithComponents(
			fakeNacosComponent{naming: naming},
			serviceregistry.New(),
			http.New(http.WithConfig(httpConfig)),
			rpc.New(rpc.WithConfig(rpcConfig)),
		),
	)
	if err := app.Start(context.Background()); err != nil {
		t.Fatalf("start: %v", err)
	}
	registered, _ := naming.counts()
	if registered != 2 {
		t.Fatalf("registered = %d", registered)
	}
	if err := app.Stop(context.Background()); err != nil {
		t.Fatalf("stop: %v", err)
	}
	_, deregistered := naming.counts()
	if deregistered != 2 {
		t.Fatalf("deregistered = %d", deregistered)
	}
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
