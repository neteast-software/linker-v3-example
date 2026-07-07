package config

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadAppliesYAMLRegistryAndEnvInOrder(t *testing.T) {
	t.Setenv(EnvConfig, "")
	dir := t.TempDir()
	file := filepath.Join(dir, "app.yaml")
	if err := os.WriteFile(file, []byte(`
linker-v3-example:
  http:
    addr: 127.0.0.1:8801
  postgresql:
    host: yaml-host
    user: yaml-user
`), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	t.Setenv("LINKER_V3_EXAMPLE_PG_HOST", "env-host")

	registry := SourceFunc(func(_ context.Context, seed Config) (Config, error) {
		if seed.PostgreSQL.Host != "yaml-host" || seed.PostgreSQL.User != "yaml-user" {
			t.Fatalf("unexpected seed: %+v", seed.PostgreSQL)
		}
		seed.PostgreSQL.Host = "registry-host"
		seed.PostgreSQL.DBName = "registry_db"
		return seed, nil
	})

	cfg, err := Load(context.Background(), WithFiles(file), WithSource(registry))
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if cfg.HTTP.Addr != "127.0.0.1:8801" {
		t.Fatalf("http addr = %s", cfg.HTTP.Addr)
	}
	if cfg.PostgreSQL.Host != "env-host" {
		t.Fatalf("pg host = %s", cfg.PostgreSQL.Host)
	}
	if cfg.PostgreSQL.User != "yaml-user" || cfg.PostgreSQL.DBName != "registry_db" {
		t.Fatalf("pg config = %+v", cfg.PostgreSQL)
	}
}

func TestLoadReadsConfigPathFromEnv(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "app.yaml")
	if err := os.WriteFile(file, []byte(`
linker-v3-example:
  grpc:
    addr: 127.0.0.1:9901
`), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	t.Setenv(EnvConfig, file)

	cfg, err := Load(context.Background(), WithoutEnv())
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if cfg.GRPC.Addr != "127.0.0.1:9901" {
		t.Fatalf("grpc addr = %s", cfg.GRPC.Addr)
	}
}

func TestLoadLocalConfigErrorIsActionable(t *testing.T) {
	t.Setenv(EnvConfig, "")
	_, err := Load(context.Background(), WithFiles("missing.yaml"), WithoutEnv())
	if err == nil {
		t.Fatal("expected config error")
	}
	if !strings.Contains(err.Error(), "读取本地配置失败") {
		t.Fatalf("unexpected error: %v", err)
	}
}
