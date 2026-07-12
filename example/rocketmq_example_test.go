package example

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	rocketmqcore "github.com/neteast-software/go-module/mq/rocketmq"
	rocketmq "github.com/neteast-software/go-module/mq/rocketmq/linker"
	linker "github.com/neteast-software/linker/v3"
)

func TestOptionalRocketMQProviderExample(t *testing.T) {
	endpoint := strings.TrimSpace(os.Getenv("LINKER_V3_EXAMPLE_ROCKETMQ_ENDPOINT"))
	if endpoint == "" {
		t.Skip("未设置 LINKER_V3_EXAMPLE_ROCKETMQ_ENDPOINT，跳过真实 RocketMQ example")
	}
	topic := requireRocketMQEnv(t, "LINKER_V3_EXAMPLE_ROCKETMQ_TOPIC")
	group := requireRocketMQEnv(t, "LINKER_V3_EXAMPLE_ROCKETMQ_CONSUMER_GROUP")

	config := rocketmqcore.DefaultConfig()
	config.Endpoint = endpoint
	config.NameSpace = strings.TrimSpace(os.Getenv("LINKER_V3_EXAMPLE_ROCKETMQ_NAMESPACE"))
	config.AccessKey = strings.TrimSpace(os.Getenv("LINKER_V3_EXAMPLE_ROCKETMQ_ACCESS_KEY"))
	config.AccessSecret = os.Getenv("LINKER_V3_EXAMPLE_ROCKETMQ_ACCESS_SECRET")
	config.SecurityToken = os.Getenv("LINKER_V3_EXAMPLE_ROCKETMQ_SECURITY_TOKEN")
	config.EnableTLS = rocketMQTLS(t)
	config.Producer.Topics = []string{topic}
	config.Consumer.Group = group
	config.Consumer.Subscriptions = []rocketmqcore.Subscription{rocketmqcore.Subscribe(topic)}
	config.Consumer.ThreadCount = 1
	config.Consumer.MaxCacheMessageCount = 8
	// PushConsumer 会等当前 long poll 结束；framework 的总关闭窗口必须更大。
	config.ShutdownTimeout = 45 * time.Second

	messageKey := fmt.Sprintf("linker-v3-example-%d", time.Now().UnixNano())
	received := make(chan rocketmqcore.Message, 1)
	component := rocketmq.New(
		rocketmq.WithConfig(config),
		rocketmq.Consume(func(_ context.Context, message *rocketmqcore.Message) error {
			if message != nil && containsRocketMQKey(message.Keys, messageKey) {
				select {
				case received <- message.Clone():
				default:
				}
			}
			return nil
		}),
	)
	app := linker.New(
		linker.WithMode(linker.Bin),
		linker.WithShutdownTimeout(50*time.Second),
		linker.WithComponents(component),
	)

	startCtx, startCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer startCancel()
	if err := app.Start(startCtx); err != nil {
		t.Fatalf("启动 RocketMQ example: %v", err)
	}
	t.Cleanup(func() {
		stopCtx, stopCancel := context.WithTimeout(context.Background(), 50*time.Second)
		defer stopCancel()
		if err := app.Stop(stopCtx); err != nil {
			t.Errorf("关闭 RocketMQ example: %v", err)
		}
	})

	client, err := rocketmq.Require(app)
	if err != nil {
		t.Fatalf("获取 RocketMQ producer capability: %v", err)
	}
	if _, err = rocketmq.RequireConsumer(app); err != nil {
		t.Fatalf("获取 RocketMQ consumer capability: %v", err)
	}
	sendCtx, sendCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer sendCancel()
	if _, err = client.SendSync(
		sendCtx,
		topic,
		[]byte("linker-v3-example"),
		rocketmqcore.WithKeys(messageKey),
		rocketmqcore.WithProperty("source", "linker-v3-example"),
	); err != nil {
		t.Fatalf("发送 RocketMQ 消息: %v", err)
	}

	select {
	case message := <-received:
		if string(message.Body) != "linker-v3-example" || message.Properties["source"] != "linker-v3-example" {
			t.Fatalf("RocketMQ 消息内容不完整: %#v", message)
		}
	case <-time.After(45 * time.Second):
		t.Fatal("RocketMQ consumer 未在 45 秒内收到测试消息")
	}

	plan := app.Plan()
	if !planHasComponent(plan, rocketmq.ID) ||
		!planHasAsset(plan, string(rocketmq.AssetProducer), topic) ||
		!planHasAsset(plan, string(rocketmq.AssetConsumer), group) {
		t.Fatalf("RocketMQ component 或资产未进入 Plan: %#v", plan)
	}
}

func requireRocketMQEnv(t *testing.T, name string) string {
	t.Helper()
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		t.Fatalf("启用 RocketMQ example 时必须设置 %s", name)
	}
	return value
}

func rocketMQTLS(t *testing.T) bool {
	t.Helper()
	value := strings.TrimSpace(os.Getenv("LINKER_V3_EXAMPLE_ROCKETMQ_TLS"))
	if value == "" {
		return false
	}
	enabled, err := strconv.ParseBool(value)
	if err != nil {
		t.Fatalf("LINKER_V3_EXAMPLE_ROCKETMQ_TLS 必须是布尔值")
	}
	return enabled
}

func containsRocketMQKey(keys []string, target string) bool {
	for _, key := range keys {
		if key == target {
			return true
		}
	}
	return false
}
