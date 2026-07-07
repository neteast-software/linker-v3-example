package notification

import (
	"context"

	mq "github.com/neteast-software/go-module/mq/consumer/linker"
)

func (p *Component) handleMessage(ctx context.Context, message mq.Message) error {
	body := string(message.Body)
	if body == "" {
		body = message.Key
	}
	if err := p.provider.Send(ctx, "mq", body); err != nil {
		return err
	}
	return p.recordFeedback(ctx, feedback{
		auditAction:   "通知消息消费",
		auditResource: "mq.notification.message",
		eventKind:     "notification.consumer",
		eventMessage:  "通知消息已消费",
		context:       map[string]any{"key": message.Key, "topic": message.Topic},
	})
}

func (p *Component) runJob(ctx context.Context) error {
	if err := p.provider.Send(ctx, "cron", "notification-health"); err != nil {
		return err
	}
	return p.recordFeedback(ctx, feedback{
		auditAction:   "通知健康采样",
		auditResource: "cron.notification.health",
		eventKind:     "notification.cron",
		eventMessage:  "通知健康采样完成",
		context:       map[string]any{"job": "notification.health"},
	})
}
