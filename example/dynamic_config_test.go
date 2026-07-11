package example

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	env "github.com/neteast-software/go-module/config/env/linker"
	yaml "github.com/neteast-software/go-module/config/yaml/linker"
	postgresql "github.com/neteast-software/go-module/db/postgresql/linker"
	httpclient "github.com/neteast-software/go-module/http/client/linker"
	server "github.com/neteast-software/go-module/linker/server"
	nacos "github.com/neteast-software/go-module/registry/nacos/linker"
	linker "github.com/neteast-software/linker/v3"
)

type fakeNacosConfig struct {
	initial []byte
	updates chan []byte
	started chan struct{}
}

func newFakeNacosConfig(initial []byte) *fakeNacosConfig {
	return &fakeNacosConfig{
		initial: append([]byte(nil), initial...),
		updates: make(chan []byte, 4),
		started: make(chan struct{}, 1),
	}
}

func (p *fakeNacosConfig) fetch(context.Context, nacos.Request) ([]byte, error) {
	return append([]byte(nil), p.initial...), nil
}

func (p *fakeNacosConfig) listen(ctx context.Context, _ nacos.Request, publish func([]byte) error) error {
	p.started <- struct{}{}
	for {
		select {
		case <-ctx.Done():
			return nil
		case content := <-p.updates:
			if err := publish(content); err != nil {
				return err
			}
		}
	}
}

type postgresqlConfigProbe struct{}

func (p *postgresqlConfigProbe) Identity() linker.ID {
	return postgresql.ID
}

func (p *postgresqlConfigProbe) Configs() []linker.Config {
	return postgresql.New().Configs()
}

func TestFakeNacosWatchAppliesLiveAndMarksRestart(t *testing.T) {
	local := filepath.Join(t.TempDir(), "app.yaml")
	if err := os.WriteFile(local, []byte(`
http/client:
  base_url: https://local.example
  timeout: 1s
db/postgresql:
  host: local-db
  port: 5432
  user: example
  db_name: example
`), 0o600); err != nil {
		t.Fatalf("write local config: %v", err)
	}
	remote := newFakeNacosConfig([]byte(`
http/client:
  base_url: https://remote-old.example
  timeout: 2s
db/postgresql:
  host: remote-old-db
  port: 5432
  user: example
  db_name: example
`))
	t.Setenv("DYNAMIC_HTTP_CLIENT__TIMEOUT", "3s")
	events := make(chan linker.ConfigEvent, 8)
	client := httpclient.New()
	app := server.New(
		server.Config(
			yaml.File(local),
			nacos.Config(
				"app.yaml",
				nacos.Bootstrap(nacos.ClientConfig{Host: "127.0.0.1", Port: 8848}),
				nacos.Fetch(remote.fetch),
				nacos.Listen(remote.listen),
			),
			env.Prefix("DYNAMIC_"),
		),
		server.WithConfigDebounce(0),
		server.WithConfigObserver(func(_ context.Context, event linker.ConfigEvent) { events <- event }),
		server.WithoutHTTP(),
		server.WithoutEvent(),
		server.WithoutNotice(),
		server.WithoutAudit(),
		server.WithoutStartupLog(),
		server.WithComponents(client, &postgresqlConfigProbe{}),
	)
	if err := app.Start(context.Background()); err != nil {
		t.Fatalf("start: %v", err)
	}
	t.Cleanup(func() { _ = app.Stop(context.Background()) })
	awaitDynamicValue(t, remote.started)
	if config := client.Client().Config(); config.BaseURL != "https://remote-old.example" || config.Timeout != 3*time.Second {
		t.Fatalf("initial effective config = %#v", config)
	}

	remote.updates <- []byte(`
http/client:
  base_url: https://remote-new.example
  timeout: 9s
db/postgresql:
  host: remote-new-db
  port: 5432
  user: example
  db_name: example
`)
	event := awaitDynamicValue(t, events)
	if event.Status != linker.ConfigRestartRequired || !event.Plan.RestartRequired {
		t.Fatalf("event = %#v", event)
	}
	if config := client.Client().Config(); config.BaseURL != "https://remote-new.example" || config.Timeout != 3*time.Second {
		t.Fatalf("live config = %#v", config)
	}
	database, ok := app.Setting(postgresql.Namespace)
	if !ok || string(database) != `{"db_name":"example","host":"remote-old-db","port":5432,"user":"example"}` {
		t.Fatalf("active database config = %s, %v", database, ok)
	}
	assertConfigMode(t, event.Plan, httpclient.Namespace, httpclient.ID, linker.ConfigLive, false)
	assertConfigMode(t, event.Plan, postgresql.Namespace, postgresql.ID, linker.ConfigRestart, true)

	remote.updates <- []byte("http/client: [")
	if rejected := awaitDynamicValue(t, events); rejected.Status != linker.ConfigRejected {
		t.Fatalf("invalid YAML event = %#v", rejected)
	}
	remote.updates <- []byte(`
http/client:
  base_url: https://remote-recovered.example
db/postgresql:
  host: remote-new-db
  port: 5432
  user: example
  db_name: example
`)
	if recovered := awaitDynamicValue(t, events); recovered.Status != linker.ConfigRestartRequired {
		t.Fatalf("recovery event = %#v", recovered)
	}
	if client.Client().Config().BaseURL != "https://remote-recovered.example" {
		t.Fatalf("recovered client config = %#v", client.Client().Config())
	}
}

func assertConfigMode(t *testing.T, plan linker.ConfigPlan, namespace linker.Namespace, owner linker.ID, mode linker.ConfigMode, changed bool) {
	t.Helper()
	for _, value := range plan.Namespaces {
		if value.Namespace == namespace {
			if value.Owner != owner || value.Mode != mode || value.Changed != changed {
				t.Fatalf("namespace plan = %#v", value)
			}
			return
		}
	}
	t.Fatalf("namespace %s missing from plan", namespace)
}

func awaitDynamicValue[T any](t *testing.T, values <-chan T) T {
	t.Helper()
	select {
	case value := <-values:
		return value
	case <-time.After(2 * time.Second):
		t.Fatal("等待 example 动态配置超时")
		var zero T
		return zero
	}
}
