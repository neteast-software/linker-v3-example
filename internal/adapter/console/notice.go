package console

import (
	"context"
	"time"

	"github.com/neteast-software/go-module/graph/console/notification"
	"github.com/neteast-software/go-module/graph/console/protocol"
)

func (p *Provider) List(_ context.Context, _ string) ([]notification.Notice, error) {
	target := protocol.Native("order.list")
	p.mu.RLock()
	defer p.mu.RUnlock()
	_, firstRead := p.read["order-delayed"]
	_, secondRead := p.read["system-ready"]
	return []notification.Notice{
		{
			ID:        "order-delayed",
			Level:     notification.Warning,
			Title:     "订单处理超时",
			Content:   "NO-20260716-001 已超过演示处理时限",
			Read:      firstRead,
			CreatedAt: time.Date(2026, 7, 16, 16, 30, 0, 0, time.Local),
			Target:    &target,
		},
		{
			ID:        "system-ready",
			Level:     notification.Success,
			Title:     "Graph Console 已就绪",
			Content:   "登录、权限、页面和业务操作链路已经装配",
			Read:      secondRead,
			CreatedAt: time.Date(2026, 7, 16, 15, 0, 0, 0, time.Local),
		},
	}, nil
}

func (p *Provider) Read(_ context.Context, _ string, id string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.read[id] = struct{}{}
	return nil
}
