package main

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	yaml "github.com/neteast-software/go-module/config/yaml/linker"
	opentelemetry "github.com/neteast-software/go-module/observe/tracing/opentelemetry/linker"
	linker "github.com/neteast-software/linker/v3"
)

func TestConfigSourcesRejectInvalidNacosPort(t *testing.T) {
	t.Setenv("LINKER_V3_EXAMPLE_NACOS_DATA_ID", "app.yaml")
	t.Setenv("LINKER_V3_EXAMPLE_NACOS_PORT", "invalid")

	_, err := configSources()
	if err == nil || !strings.Contains(err.Error(), "NACOS_PORT") {
		t.Fatalf("err = %v", err)
	}
}

func TestConfigSourcesKeepDeclaredOrder(t *testing.T) {
	t.Setenv("LINKER_V3_EXAMPLE_NACOS_DATA_ID", "app.yaml")
	t.Setenv("LINKER_V3_EXAMPLE_NACOS_HOST", "127.0.0.1")

	sources, err := configSources()
	if err != nil {
		t.Fatalf("config sources: %v", err)
	}
	if len(sources) != 3 {
		t.Fatalf("sources = %d", len(sources))
	}
	want := []string{"config/yaml", "registry/nacos/config", "config/env"}
	for index, source := range sources {
		if source.Name() != want[index] {
			t.Fatalf("source %d = %s", index, source.Name())
		}
	}
}

func TestTracingDefaultsToMemoryAndEnvCanEnableOTLP(t *testing.T) {
	setting, err := yaml.File("config/app.example.yaml").Load(context.Background(), linker.BootstrapContext{})
	if err != nil {
		t.Fatalf("load yaml: %v", err)
	}
	content, ok := setting.Lookup(opentelemetry.Namespace)
	if !ok {
		t.Fatal("observe/tracing config missing")
	}
	var local opentelemetry.Config
	if err = json.Unmarshal(content, &local); err != nil {
		t.Fatalf("decode local tracing: %v", err)
	}
	if local.Mode != opentelemetry.ModeMemory || local.Service != "linker-v3-example" {
		t.Fatalf("local tracing = %#v", local)
	}

	t.Setenv("APP_OBSERVE_TRACING__MODE", "otlp")
	t.Setenv("APP_OBSERVE_TRACING__ENDPOINT", "127.0.0.1:4317")
	t.Setenv("APP_OBSERVE_TRACING__PROTOCOL", "grpc")
	t.Setenv("APP_OBSERVE_TRACING__INSECURE", "true")
	sources, err := configSources()
	if err != nil {
		t.Fatalf("sources: %v", err)
	}
	override, err := sources[len(sources)-1].Load(context.Background(), linker.BootstrapContext{})
	if err != nil {
		t.Fatalf("load env: %v", err)
	}
	content, ok = override.Lookup(opentelemetry.Namespace)
	if !ok {
		t.Fatal("OTLP env override missing")
	}
	var remote opentelemetry.Config
	if err = json.Unmarshal(content, &remote); err != nil {
		t.Fatalf("decode env tracing: %v", err)
	}
	if remote.Mode != opentelemetry.ModeOTLP || remote.Endpoint != "127.0.0.1:4317" || !remote.Insecure {
		t.Fatalf("remote tracing = %#v", remote)
	}
}
