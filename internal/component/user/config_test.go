package user

import (
	"context"
	"strings"
	"testing"

	linker "github.com/neteast-software/linker/v3"
)

func TestComponentOwnsTokenConfig(t *testing.T) {
	component := NewComponent()
	declarations := component.Configs()
	if len(declarations) != 1 || declarations[0].Namespace() != Namespace || declarations[0].Mode() != linker.ConfigRestart {
		t.Fatalf("configs = %#v", declarations)
	}
	if err := component.Bootstrap(context.Background(), linker.BootstrapContext{Seed: linker.NewSetting(map[linker.Namespace][]byte{
		Namespace: []byte(`{"token_key":"short"}`),
	})}); err == nil || !strings.Contains(err.Error(), "至少需要 32") {
		t.Fatalf("short token err = %v", err)
	}
	if err := component.Bootstrap(context.Background(), linker.BootstrapContext{Seed: linker.NewSetting(map[linker.Namespace][]byte{
		Namespace: []byte(`{"token_key":"0123456789abcdef0123456789abcdef"}`),
	})}); err != nil {
		t.Fatalf("bootstrap: %v", err)
	}
}
