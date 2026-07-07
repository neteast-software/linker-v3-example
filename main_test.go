package main

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestPlanCommand(t *testing.T) {
	t.Setenv("LINKER_V3_EXAMPLE_PG_PASSWORD", "")

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
}

func TestPlanCommandArg(t *testing.T) {
	if !isPlanCommand([]string{"linker-v3-example", "--plan"}) {
		t.Fatal("expected --plan command")
	}
	if isPlanCommand([]string{"linker-v3-example"}) {
		t.Fatal("unexpected plan command")
	}
}
