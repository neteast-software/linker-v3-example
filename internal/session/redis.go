package session

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	redis "github.com/neteast-software/go-module/cache/redis"
)

// RedisStore 把 session 状态映射到现有 Redis component 的 client。
type RedisStore struct {
	client *redis.Client
	prefix string
	now    func() time.Time
}

// Redis 创建不取得 client 生命周期的 session store。
func Redis(client *redis.Client, prefix string) *RedisStore {
	return &RedisStore{client: client, prefix: strings.TrimSpace(prefix), now: time.Now}
}

func (p *RedisStore) Save(ctx context.Context, current Session) error {
	if p == nil || p.client == nil || p.prefix == "" {
		return fmt.Errorf("Redis session store 未配置")
	}
	now := p.now()
	if err := current.validate(now); err != nil {
		return err
	}
	content, err := json.Marshal(current)
	if err != nil {
		return fmt.Errorf("session 编码失败: %w", err)
	}
	return p.client.Set(ctx, p.key(current.ID), content, current.ExpiresAt.Sub(now))
}

func (p *RedisStore) Lookup(ctx context.Context, id string) (Session, error) {
	if p == nil || p.client == nil || p.prefix == "" {
		return Session{}, fmt.Errorf("Redis session store 未配置")
	}
	content, err := p.client.Get(ctx, p.key(id))
	if err != nil {
		if err == redis.ErrNotFound {
			return Session{}, ErrNotFound
		}
		return Session{}, err
	}
	var current Session
	if err = json.Unmarshal([]byte(content), &current); err != nil {
		return Session{}, fmt.Errorf("session 解码失败: %w", err)
	}
	if err = current.validate(p.now()); err != nil {
		return Session{}, err
	}
	return current, nil
}

func (p *RedisStore) Revoke(ctx context.Context, id string) error {
	if p == nil || p.client == nil || p.prefix == "" {
		return fmt.Errorf("Redis session store 未配置")
	}
	return p.client.Delete(ctx, p.key(id))
}

func (p *RedisStore) key(id string) string {
	return p.prefix + ":" + id
}
