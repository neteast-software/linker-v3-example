package main

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	yaml "github.com/neteast-software/go-module/config/yaml/linker"
	postgresql "github.com/neteast-software/go-module/db/postgresql/linker"
	linker "github.com/neteast-software/linker/v3"

	"linker-v3-example/internal/app"
	user "linker-v3-example/internal/user/linker"
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
	if !jsonPlanHasAsset(assets, "rpc/grpc/client", "example/tts-client") {
		t.Fatalf("plan missing grpc client asset: %#v", assets)
	}
	if !jsonPlanHasAsset(assets, "observe/metrics", "prometheus") {
		t.Fatalf("plan missing metrics asset: %#v", assets)
	}
	if !jsonPlanHasAsset(assets, "observe/tracing", "linker-v3-example") {
		t.Fatalf("plan missing tracing asset: %#v", assets)
	}
}

func TestExampleConfigDoesNotCarryCredentials(t *testing.T) {
	setting, err := yaml.File("config/app.example.yaml").Load(context.Background(), linker.BootstrapContext{})
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	for namespace, fields := range map[linker.Namespace][]string{
		postgresql.Namespace: {"password"},
		user.Namespace:       {"token_key", "seed_password"},
	} {
		content, ok := setting.Lookup(namespace)
		if !ok {
			t.Fatalf("namespace %s missing", namespace)
		}
		var value map[string]any
		if err = json.Unmarshal(content, &value); err != nil {
			t.Fatalf("decode %s: %v", namespace, err)
		}
		for _, field := range fields {
			if _, exists := value[field]; exists {
				t.Fatalf("config %s contains credential field %s", namespace, field)
			}
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

func TestGatewayCommandArg(t *testing.T) {
	if nacos, ok := gatewayCommand([]string{"linker-v3-example", "--gateway"}); !ok || nacos {
		t.Fatalf("--gateway = (%v, %v)", nacos, ok)
	}
	if nacos, ok := gatewayCommand([]string{"linker-v3-example", "--gateway-nacos"}); !ok || !nacos {
		t.Fatalf("--gateway-nacos = (%v, %v)", nacos, ok)
	}
	if _, ok := gatewayCommand([]string{"linker-v3-example"}); ok {
		t.Fatal("empty command was recognized as Gateway")
	}
}

func TestGatewayNacosProfileCanPrepareWithoutCredentialsOrNetwork(t *testing.T) {
	t.Setenv("LINKER_V3_EXAMPLE_NACOS_DATA_ID", "")
	sources, err := gatewaySources(true)
	if err != nil {
		t.Fatal(err)
	}
	gatewayApp, err := app.GatewayNacos(sources...)
	if err != nil {
		t.Fatal(err)
	}
	if err = gatewayApp.Prepare(context.Background()); err != nil {
		t.Fatal(err)
	}
	required := map[linker.ID]bool{
		"registry/nacos":           false,
		"registry/discovery/nacos": false,
		"http/gateway":             false,
	}
	for _, component := range gatewayApp.Plan().Components {
		if _, ok := required[component.ID]; ok {
			required[component.ID] = true
		}
	}
	for id, found := range required {
		if !found {
			t.Fatalf("Gateway Nacos plan missing %s", id)
		}
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
