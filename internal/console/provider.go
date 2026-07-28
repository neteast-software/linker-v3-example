package console

import (
	"sync"

	user "linker-v3-example/internal/user"
)

type Provider struct {
	user user.Auth
	mu   sync.RWMutex
	read map[string]struct{}
}

func New() *Provider {
	return &Provider{
		read: make(map[string]struct{}),
	}
}

// Configure 装配 Console 运行所需的用户认证能力。
func (p *Provider) Configure(auth user.Auth) error {
	if auth == nil {
		return user.ErrUnavailable
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.user = auth
	return nil
}

func (p *Provider) auth() (user.Auth, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if p.user == nil {
		return nil, user.ErrUnavailable
	}
	return p.user, nil
}
