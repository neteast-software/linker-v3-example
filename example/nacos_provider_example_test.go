package example

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	nacoscore "github.com/neteast-software/go-module/registry/nacos"
	nacos "github.com/neteast-software/go-module/registry/nacos/linker"
	linker "github.com/neteast-software/linker/v3"
)

const nacosProviderNamespace linker.Namespace = "example/provider"

type nacosProviderConfig struct {
	Value string `json:"value"`
}

type nacosProviderProbe struct {
	config nacosProviderConfig
}

func (p *nacosProviderProbe) Identity() linker.ID {
	return "example/nacos-provider"
}

func (p *nacosProviderProbe) Bootstrap(_ context.Context, boot linker.BootstrapContext) error {
	content, ok := boot.Seed.Lookup(nacosProviderNamespace)
	if !ok {
		return fmt.Errorf("Nacos 配置缺少 %s", nacosProviderNamespace)
	}
	if err := json.Unmarshal(content, &p.config); err != nil {
		return fmt.Errorf("解析 Nacos example 配置: %w", err)
	}
	return nil
}

func TestOptionalNacosProviderExample(t *testing.T) {
	config, ok := nacosExampleConfig(t)
	if !ok {
		t.Skip("未设置 LINKER_V3_EXAMPLE_NACOS_HOST，跳过真实 Nacos example")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := nacoscore.NewClient(config)
	if err != nil {
		t.Fatalf("创建 Nacos example client: %v", err)
	}
	t.Cleanup(func() {
		if err := client.Close(); err != nil {
			t.Errorf("关闭 Nacos example client: %v", err)
		}
	})

	dataID := fmt.Sprintf("linker-v3-example-%d.yaml", time.Now().UnixNano())
	group := config.Group
	content := "example/provider:\n  value: ready\n"
	if err = client.Config().Publish(ctx, dataID, content, nacoscore.Group(group), nacoscore.Type("yaml")); err != nil {
		t.Fatalf("发布 Nacos example 配置: %v", err)
	}
	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cleanupCancel()
		if err := client.Config().Delete(cleanupCtx, dataID, nacoscore.Group(group)); err != nil {
			t.Errorf("清理 Nacos example data id: %v", err)
		}
	})

	probe := &nacosProviderProbe{}
	app := linker.New(
		linker.WithMode(linker.Bin),
		linker.WithSource(nacos.Config(dataID, nacos.Group(group), nacos.Bootstrap(config))),
		linker.WithComponents(probe),
	)
	if err = app.Start(ctx); err != nil {
		t.Fatalf("从 Nacos 配置启动 example: %v", err)
	}
	t.Cleanup(func() {
		stopCtx, stopCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer stopCancel()
		if err := app.Stop(stopCtx); err != nil {
			t.Errorf("关闭 Nacos example: %v", err)
		}
	})
	if probe.config.Value != "ready" {
		t.Fatalf("Nacos example value = %q，期望 ready", probe.config.Value)
	}
	if !planHasComponent(app.Plan(), probe.Identity()) {
		t.Fatalf("Nacos 配置 owner 未进入 Plan: %#v", app.Plan())
	}
}

func nacosExampleConfig(t *testing.T) (nacoscore.Config, bool) {
	t.Helper()
	host := strings.TrimSpace(os.Getenv("LINKER_V3_EXAMPLE_NACOS_HOST"))
	if host == "" {
		return nacoscore.Config{}, false
	}
	config := nacoscore.DefaultConfig()
	config.Host = host
	config.Username = strings.TrimSpace(os.Getenv("LINKER_V3_EXAMPLE_NACOS_USERNAME"))
	config.Password = os.Getenv("LINKER_V3_EXAMPLE_NACOS_PASSWORD")
	config.NamespaceID = strings.TrimSpace(os.Getenv("LINKER_V3_EXAMPLE_NACOS_NAMESPACE"))
	if group := strings.TrimSpace(os.Getenv("LINKER_V3_EXAMPLE_NACOS_GROUP")); group != "" {
		config.Group = group
	}
	if value := strings.TrimSpace(os.Getenv("LINKER_V3_EXAMPLE_NACOS_PORT")); value != "" {
		port, err := strconv.ParseUint(value, 10, 64)
		if err != nil || port == 0 || port > 65535 {
			t.Fatalf("LINKER_V3_EXAMPLE_NACOS_PORT 必须是 1 到 65535 的端口")
		}
		config.Port = port
	}
	return config, true
}
