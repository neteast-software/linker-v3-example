package main

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	yaml "github.com/neteast-software/go-module/config/yaml/linker"
	postgresql "github.com/neteast-software/go-module/db/postgresql/linker"
	linker "github.com/neteast-software/linker/v3"

	usercomponent "linker-v3-example/internal/component/user"
)

func TestPlanCommand(t *testing.T) {
	t.Setenv("LINKER_V3_EXAMPLE_NACOS_DATA_ID", "")

	var output bytes.Buffer
	if err := printPlan(&output); err != nil {
		t.Fatalf("print plan: %v", err)
	}

	var body map[string]any
	if err := json.Unmarshal(output.Bytes(), &body); err != nil {
		t.Fatalf("decode plan: %v\n%s", err, output.String())
	}
	if body["mode"] != "server" {
		t.Fatalf("unexpected plan mode: %#v", body["mode"])
	}
	components, ok := body["components"].([]any)
	if !ok || len(components) == 0 {
		t.Fatalf("plan missing components: %#v", body)
	}
	assets, ok := body["assets"].([]any)
	if !ok || len(assets) == 0 {
		t.Fatalf("plan missing assets: %#v", body)
	}
	if !jsonPlanHasAsset(assets, "rpc/grpc/server", "127.0.0.1:9900") {
		t.Fatalf("plan missing grpc server asset: %#v", assets)
	}
	if !jsonPlanHasAsset(assets, "rpc/grpc/client", "rpc/client/tts") {
		t.Fatalf("plan missing grpc client asset: %#v", assets)
	}
	if !jsonPlanHasAsset(assets, "observe/metrics", "prometheus") {
		t.Fatalf("plan missing metrics asset: %#v", assets)
	}
	if !jsonPlanHasAsset(assets, "observe/tracing", "http+grpc") {
		t.Fatalf("plan missing tracing asset: %#v", assets)
	}
}

func TestExampleConfigDoesNotCarryCredentials(t *testing.T) {
	setting, err := yaml.File("config/app.example.yaml").Load(context.Background(), linker.BootstrapContext{})
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	for namespace, field := range map[linker.Namespace]string{
		postgresql.Namespace:    "password",
		usercomponent.Namespace: "token_key",
	} {
		content, ok := setting.Lookup(namespace)
		if !ok {
			t.Fatalf("namespace %s missing", namespace)
		}
		var value map[string]any
		if err = json.Unmarshal(content, &value); err != nil {
			t.Fatalf("decode %s: %v", namespace, err)
		}
		if _, exists := value[field]; exists {
			t.Fatalf("config %s contains credential field %s", namespace, field)
		}
	}
}

func TestPlanCommandArg(t *testing.T) {
	if !isPlanCommand([]string{"linker-v3-example", "--plan"}) {
		t.Fatal("expected --plan command")
	}
	if isPlanCommand([]string{"linker-v3-example"}) {
		t.Fatal("unexpected plan command")
	}
}

func jsonPlanHasAsset(assets []any, kind string, name string) bool {
	for _, item := range assets {
		asset, ok := item.(map[string]any)
		if !ok {
			continue
		}
		if asset["kind"] == kind && asset["name"] == name {
			return true
		}
	}
	return false
}
