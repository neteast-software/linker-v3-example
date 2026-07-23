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

func New(auth user.Auth) *Provider {
	return &Provider{
		user: auth,
		read: make(map[string]struct{}),
	}
}
