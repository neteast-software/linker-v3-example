package captcha

import "context"

// Store 保存并原子消费一次性 challenge。
type Store interface {
	Save(context.Context, Challenge) error
	Take(context.Context, string) (Challenge, error)
}
