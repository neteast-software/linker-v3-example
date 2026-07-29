package captcha

import (
	"context"
	"sync"
	"time"
)

// MemoryStore 是默认 profile 的一次性 challenge store。
type MemoryStore struct {
	mutex      sync.Mutex
	challenges map[string]Challenge
	now        func() time.Time
}

// Memory 创建不访问外部基础设施的 challenge store。
func Memory() *MemoryStore {
	return &MemoryStore{challenges: make(map[string]Challenge), now: time.Now}
}

func (p *MemoryStore) Save(_ context.Context, challenge Challenge) error {
	if !p.now().Before(challenge.ExpiresAt) {
		return ErrExpired
	}
	p.mutex.Lock()
	p.challenges[challenge.ID] = challenge
	p.mutex.Unlock()
	return nil
}

func (p *MemoryStore) Take(_ context.Context, id string) (Challenge, error) {
	p.mutex.Lock()
	challenge, ok := p.challenges[id]
	delete(p.challenges, id)
	p.mutex.Unlock()
	if !ok {
		return Challenge{}, ErrNotFound
	}
	if !p.now().Before(challenge.ExpiresAt) {
		return Challenge{}, ErrExpired
	}
	return challenge, nil
}
