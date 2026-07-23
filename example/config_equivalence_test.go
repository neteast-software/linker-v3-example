package example

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	configcore "github.com/neteast-software/go-module/config"
	yaml "github.com/neteast-software/go-module/config/yaml/linker"
	postgresql "github.com/neteast-software/go-module/db/postgresql"
	http "github.com/neteast-software/go-module/http/gin"
	nacos "github.com/neteast-software/go-module/registry/nacos/linker"
	rpc "github.com/neteast-software/go-module/rpc/grpc"
	linker "github.com/neteast-software/linker/v3"

	user "linker-v3-example/internal/user/linker"
)

type effectiveConfig struct {
	http       http.Config
	postgresql postgresql.Config
	rpc        rpc.ServerConfig
	user       user.Config
}

func TestLocalAndNacosRootYAMLProduceEquivalentSetting(t *testing.T) {
	content := fmt.Appendf(nil, `
http:
  addr: 127.0.0.1:8080
  read_timeout: 2s
  health:
    readiness:
      enabled: true
      path: ready
cache/redis:
  addr: 127.0.0.1:6379
  db: 0
db/postgresql:
  host: 127.0.0.1
  port: 5432
  user: example
  db_name: example
rpc/grpc:
  addr: 127.0.0.1:9900
example/user:
  token_key: %s
`, strings.Repeat("a", 32))
	path := filepath.Join(t.TempDir(), "app.yaml")
	if err := os.WriteFile(path, content, 0o600); err != nil {
		t.Fatalf("write yaml: %v", err)
	}

	local, err := yaml.File(path).Load(context.Background(), linker.BootstrapContext{})
	if err != nil {
		t.Fatalf("load local: %v", err)
	}
	remote, err := nacos.Config(
		"app.yaml",
		nacos.Bootstrap(nacos.ClientConfig{Host: "127.0.0.1", Port: 8848}),
		nacos.Fetch(func(context.Context, nacos.Request) ([]byte, error) { return content, nil }),
	).Load(context.Background(), linker.BootstrapContext{})
	if err != nil {
		t.Fatalf("load nacos: %v", err)
	}
	if !reflect.DeepEqual(local.Namespaces(), remote.Namespaces()) {
		t.Fatalf("namespaces: local=%v remote=%v", local.Namespaces(), remote.Namespaces())
	}
	for _, namespace := range local.Namespaces() {
		localValue, _ := local.Lookup(namespace)
		remoteValue, _ := remote.Lookup(namespace)
		if !reflect.DeepEqual(localValue, remoteValue) {
			t.Fatalf("namespace %s: local=%s remote=%s", namespace, localValue, remoteValue)
		}
	}
	localConfig := projectEffectiveConfig(t, local)
	remoteConfig := projectEffectiveConfig(t, remote)
	if !reflect.DeepEqual(localConfig, remoteConfig) {
		t.Fatalf("effective config: local=%#v remote=%#v", localConfig, remoteConfig)
	}
}

func projectEffectiveConfig(t *testing.T, setting linker.Setting) effectiveConfig {
	t.Helper()
	ret := effectiveConfig{
		http:       http.DefaultConfig(),
		postgresql: postgresql.DefaultConfig(),
	}
	decodeNamespace(t, setting, "http", &ret.http)
	decodeNamespace(t, setting, "db/postgresql", &ret.postgresql)
	decodeNamespace(t, setting, "rpc/grpc", &ret.rpc)
	decodeNamespace(t, setting, user.Namespace, &ret.user)
	ret.http = ret.http.Normalize()
	ret.postgresql = ret.postgresql.Normalize()
	ret.rpc = ret.rpc.Normalize()
	return ret
}

func decodeNamespace(t *testing.T, setting linker.Setting, namespace linker.Namespace, target any) {
	t.Helper()
	content, ok := setting.Lookup(namespace)
	if !ok {
		t.Fatalf("namespace %s missing", namespace)
	}
	if err := configcore.Decode(content, target); err != nil {
		t.Fatalf("decode %s: %v", namespace, err)
	}
}
