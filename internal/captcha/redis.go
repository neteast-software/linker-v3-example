package captcha

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	redis "github.com/neteast-software/go-module/cache/redis"
)

// RedisStore 使用短锁保证一次性 challenge 不被并发重复消费。
type RedisStore struct {
	client *redis.Client
	prefix string
	now    func() time.Time
}

// Redis 创建不取得 client 生命周期的 challenge store。
func Redis(client *redis.Client, prefix string) *RedisStore {
	return &RedisStore{client: client, prefix: strings.TrimSpace(prefix), now: time.Now}
}

func (p *RedisStore) Save(ctx context.Context, challenge Challenge) error {
	if p == nil || p.client == nil || p.prefix == "" {
		return fmt.Errorf("Redis captcha store 未配置")
	}
	if !p.now().Before(challenge.ExpiresAt) {
		return ErrExpired
	}
	content, err := json.Marshal(challenge)
	if err != nil {
		return fmt.Errorf("验证码 challenge 编码失败: %w", err)
	}
	return p.client.Set(ctx, p.key(challenge.ID), content, challenge.ExpiresAt.Sub(p.now()))
}

func (p *RedisStore) Take(ctx context.Context, id string) (Challenge, error) {
	if p == nil || p.client == nil || p.prefix == "" {
		return Challenge{}, fmt.Errorf("Redis captcha store 未配置")
	}
	lock, err := p.client.Lock(ctx, p.lockKey(id), 3*time.Second)
	if err != nil {
		return Challenge{}, err
	}
	defer func() {
		_ = lock.Release(context.WithoutCancel(ctx))
	}()
	content, err := p.client.Get(ctx, p.key(id))
	if err != nil {
		if err == redis.ErrNotFound {
			return Challenge{}, ErrNotFound
		}
		return Challenge{}, err
	}
	if err = p.client.Delete(ctx, p.key(id)); err != nil {
		return Challenge{}, err
	}
	var challenge Challenge
	if err = json.Unmarshal([]byte(content), &challenge); err != nil {
		return Challenge{}, fmt.Errorf("验证码 challenge 解码失败: %w", err)
	}
	if !p.now().Before(challenge.ExpiresAt) {
		return Challenge{}, ErrExpired
	}
	return challenge, nil
}

func (p *RedisStore) key(id string) string {
	return p.prefix + ":" + id
}

func (p *RedisStore) lockKey(id string) string {
	return p.prefix + ":lock:" + id
}
