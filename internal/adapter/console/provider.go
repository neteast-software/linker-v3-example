package console

import (
	"sync"

	userservice "linker-v3-example/internal/service/user"
)

type Provider struct {
	user userservice.Auth
	mu   sync.RWMutex
	read map[string]struct{}
}

func New(user userservice.Auth) *Provider {
	return &Provider{
		user: user,
		read: make(map[string]struct{}),
	}
}
