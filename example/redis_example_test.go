package example

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	rediscore "github.com/neteast-software/go-module/cache/redis"
	redis "github.com/neteast-software/go-module/cache/redis/linker"
	linker "github.com/neteast-software/linker/v3"
)

func TestOptionalRedisProviderExample(t *testing.T) {
	addr := strings.TrimSpace(os.Getenv("LINKER_V3_EXAMPLE_REDIS_ADDR"))
	if addr == "" {
		t.Skip("未设置 LINKER_V3_EXAMPLE_REDIS_ADDR，跳过真实 Redis example")
	}

	config := rediscore.DefaultConfig()
	config.Addr = addr
	config.Username = strings.TrimSpace(os.Getenv("LINKER_V3_EXAMPLE_REDIS_USERNAME"))
	config.Password = os.Getenv("LINKER_V3_EXAMPLE_REDIS_PASSWORD")
	config.DB = redisExampleDatabase(t)
	config.DialTimeout = 2 * time.Second
	config.ReadTimeout = 2 * time.Second
	config.WriteTimeout = 2 * time.Second

	app := linker.New(
		linker.WithMode(linker.Bin),
		linker.WithComponents(redis.New(redis.WithConfig(config))),
	)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := app.Start(ctx); err != nil {
		t.Fatalf("启动 Redis example: %v", err)
	}
	t.Cleanup(func() {
		stopCtx, stopCancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer stopCancel()
		if err := app.Stop(stopCtx); err != nil {
			t.Errorf("关闭 Redis example: %v", err)
		}
	})

	client, err := redis.Require(app)
	if err != nil {
		t.Fatalf("获取 Redis capability: %v", err)
	}
	key := fmt.Sprintf("linker-v3-example:%d", time.Now().UnixNano())
	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cleanupCancel()
		if err := client.Delete(cleanupCtx, key); err != nil {
			t.Errorf("清理 Redis example key: %v", err)
		}
	})
	if err = client.Set(ctx, key, "ready", time.Minute); err != nil {
		t.Fatalf("写入 Redis example: %v", err)
	}
	value, err := client.Get(ctx, key)
	if err != nil {
		t.Fatalf("读取 Redis example: %v", err)
	}
	if value != "ready" {
		t.Fatalf("Redis example value = %q，期望 ready", value)
	}

	plan := app.Plan()
	if !planHasComponent(plan, redis.ID) || !planHasAsset(plan, string(redis.AssetConnection), addr) {
		t.Fatalf("Redis component 或连接资产未进入 Plan: %#v", plan)
	}
}

func redisExampleDatabase(t *testing.T) int {
	t.Helper()
	value := strings.TrimSpace(os.Getenv("LINKER_V3_EXAMPLE_REDIS_DB"))
	if value == "" {
		return 0
	}
	database, err := strconv.Atoi(value)
	if err != nil || database < 0 {
		t.Fatalf("LINKER_V3_EXAMPLE_REDIS_DB 必须是非负整数")
	}
	return database
}
