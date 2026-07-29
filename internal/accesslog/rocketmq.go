package accesslog

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	rocketmq "github.com/neteast-software/go-module/mq/rocketmq"
)

type rocketMQPublisher struct {
	client *rocketmq.Client
	topic  string
}

// RocketMQ 创建不取得 client 生命周期的 operateLog publisher。
func RocketMQ(client *rocketmq.Client, topic string) Publisher {
	return &rocketMQPublisher{client: client, topic: strings.TrimSpace(topic)}
}

func (p *rocketMQPublisher) Publish(ctx context.Context, record Record) error {
	if p == nil || p.client == nil || p.topic == "" {
		return fmt.Errorf("RocketMQ accesslog publisher 未配置")
	}
	content, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("operateLog 编码失败: %w", err)
	}
	_, err = p.client.SendSync(ctx, p.topic, content, rocketmq.WithTag("operateLog"))
	return err
}
