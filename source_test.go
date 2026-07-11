package main

import (
	"strings"
	"testing"
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
