package console

import (
	"context"
	"time"

	"github.com/neteast-software/go-module/graph/console/login"
	"github.com/neteast-software/go-module/token"

	usermodel "linker-v3-example/internal/model/user"
)

func (p *Provider) Current(ctx context.Context, raw string) (login.Session, error) {
	user, claims, err := p.user.Current(ctx, raw, "console")
	if err != nil {
		return login.Session{}, providerError(err)
	}
	return session(user, raw, claims), nil
}

func (p *Provider) Refresh(ctx context.Context, raw string) (login.Session, error) {
	user, issued, err := p.user.Refresh(ctx, raw, "console")
	if err != nil {
		return login.Session{}, providerError(err)
	}
	return session(user, issued.Raw, issued.Claims), nil
}

func (p *Provider) Revoke(ctx context.Context, raw string) error {
	return providerError(p.user.Revoke(ctx, raw, "console"))
}

func session(user usermodel.User, raw string, claims token.Claims) login.Session {
	return login.Session{
		Subject: claims.Subject,
		Token: login.Token{
			Access:    raw,
			ExpiresAt: time.Unix(claims.ExpiresAt, 0),
		},
		Profile: profile(user),
	}
}
