package session

import (
	"context"
	"sync"
	"time"
)

// MemoryStore 是默认 profile 使用的并发安全 session store。
type MemoryStore struct {
	mutex    sync.RWMutex
	sessions map[string]Session
	now      func() time.Time
}

// Memory 创建不访问外部基础设施的 session store。
func Memory() *MemoryStore {
	return &MemoryStore{sessions: make(map[string]Session), now: time.Now}
}

func (p *MemoryStore) Save(_ context.Context, current Session) error {
	if err := current.validate(p.now()); err != nil {
		return err
	}
	p.mutex.Lock()
	p.sessions[current.ID] = current
	p.mutex.Unlock()
	return nil
}

func (p *MemoryStore) Lookup(_ context.Context, id string) (Session, error) {
	p.mutex.RLock()
	current, ok := p.sessions[id]
	p.mutex.RUnlock()
	if !ok {
		return Session{}, ErrNotFound
	}
	if err := current.validate(p.now()); err != nil {
		if err == ErrExpired {
			p.mutex.Lock()
			delete(p.sessions, id)
			p.mutex.Unlock()
		}
		return Session{}, err
	}
	return current, nil
}

func (p *MemoryStore) Revoke(_ context.Context, id string) error {
	p.mutex.Lock()
	delete(p.sessions, id)
	p.mutex.Unlock()
	return nil
}
