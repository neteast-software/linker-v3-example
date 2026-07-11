package example

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	yaml "github.com/neteast-software/go-module/config/yaml/linker"
	nacos "github.com/neteast-software/go-module/registry/nacos/linker"
	linker "github.com/neteast-software/linker/v3"
)

func TestLocalAndNacosRootYAMLProduceEquivalentSetting(t *testing.T) {
	content := []byte(`
http:
  addr: 127.0.0.1:8080
  read_timeout: 2s
  health:
    enabled: true
    path: ready
cache/redis:
  addr: 127.0.0.1:6379
  db: 0
`)
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
}
