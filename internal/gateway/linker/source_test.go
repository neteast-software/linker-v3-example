package gateway

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	gatewaycomponent "github.com/neteast-software/go-module/http/gateway/linker"
	linker "github.com/neteast-software/linker/v3"
)

func TestRoutesSourceProjectsStrictDocument(t *testing.T) {
	path := filepath.Join(t.TempDir(), "routes.yaml")
	content := []byte(`schema: linker.gateway.v1
plan:
  routes:
    - id: local
      match:
        paths: [/api/**]
      upstream:
        url: http://127.0.0.1:8800
`)
	if err := os.WriteFile(path, content, 0o600); err != nil {
		t.Fatal(err)
	}
	setting, err := Routes(path).Load(context.Background(), linker.BootstrapContext{})
	if err != nil {
		t.Fatal(err)
	}
	projected, ok := setting.Lookup(gatewaycomponent.RoutesNamespace)
	if !ok || len(projected) == 0 {
		t.Fatalf("projected = %q", projected)
	}
}

func TestRoutesSourceRejectsSymlink(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "target.yaml")
	link := filepath.Join(dir, "routes.yaml")
	if err := os.WriteFile(target, []byte("schema: linker.gateway.v1"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, link); err != nil {
		t.Fatal(err)
	}
	if _, err := Routes(link).Load(context.Background(), linker.BootstrapContext{}); err == nil {
		t.Fatal("symlink was accepted")
	}
}
