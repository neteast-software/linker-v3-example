package console

import (
	"context"
	"time"

	"github.com/neteast-software/go-module/graph/console/login"
	"github.com/neteast-software/go-module/token"

	user "linker-v3-example/internal/user"
)

func (p *Provider) Current(ctx context.Context, raw string) (login.Session, error) {
	auth, err := p.auth()
	if err != nil {
		return login.Session{}, providerError(err)
	}
	user, claims, err := auth.Current(ctx, raw, "console")
	if err != nil {
		return login.Session{}, providerError(err)
	}
	return session(user, raw, claims), nil
}

func (p *Provider) Refresh(ctx context.Context, raw string) (login.Session, error) {
	auth, err := p.auth()
	if err != nil {
		return login.Session{}, providerError(err)
	}
	user, issued, err := auth.Refresh(ctx, raw, "console")
	if err != nil {
		return login.Session{}, providerError(err)
	}
	return session(user, issued.Raw, issued.Claims), nil
}

func (p *Provider) Revoke(ctx context.Context, raw string) error {
	auth, err := p.auth()
	if err != nil {
		return providerError(err)
	}
	return providerError(auth.Revoke(ctx, raw, "console"))
}

func session(current user.User, raw string, claims token.Claims) login.Session {
	return login.Session{
		Subject: claims.Subject,
		Token: login.Token{
			Access:    raw,
			ExpiresAt: time.Unix(claims.ExpiresAt, 0),
		},
		Profile: profile(current),
	}
}
