package accesslog

import "context"

// Publisher 只负责把结构化审计记录交给外部设施。
type Publisher interface {
	Publish(context.Context, Record) error
}

// PublisherFunc 把函数适配为 Publisher。
type PublisherFunc func(context.Context, Record) error

func (p PublisherFunc) Publish(ctx context.Context, record Record) error {
	if p == nil {
		return nil
	}
	return p(ctx, record)
}

// Discard 创建默认 profile 使用的无外部 IO publisher。
func Discard() Publisher {
	return PublisherFunc(nil)
}
