package user

import (
	"context"
	"strconv"
	"time"

	session "github.com/neteast-software/go-module/acl/session"
	"github.com/neteast-software/go-module/token"
)

type Service struct {
	store    Store
	signer   *token.Signer
	sessions *session.Session
}

var _ Auth = (*Service)(nil)

func New() *Service {
	return &Service{}
}

func (p *Service) Configure(store Store, signer *token.Signer, sessions *session.Session) {
	p.store = store
	p.signer = signer
	p.sessions = sessions
}

func (p *Service) AdminLogin(ctx context.Context, username string, password string) (User, string, error) {
	user, account, err := p.store.ByAdmin(ctx, username)
	if err != nil {
		return User{}, "", err
	}
	return p.login(ctx, user, account, password, "console")
}

func (p *Service) UserLogin(ctx context.Context, phone string, password string) (User, string, error) {
	user, account, err := p.store.ByPhone(ctx, phone)
	if err != nil {
		return User{}, "", err
	}
	return p.login(ctx, user, account, password, "front")
}

func (p *Service) Profile(ctx context.Context, raw string, scope string) (User, error) {
	user, _, err := p.Current(ctx, raw, scope)
	return user, err
}

func (p *Service) Current(ctx context.Context, raw string, scope string) (User, token.Claims, error) {
	claims, err := p.verify(ctx, raw, scope)
	if err != nil {
		return User{}, token.Claims{}, err
	}
	userID, err := strconv.ParseUint(claims.Subject, 10, 64)
	if err != nil {
		return User{}, token.Claims{}, ErrLogin
	}
	user, err := p.store.ByID(ctx, userID)
	return user, claims, err
}

func (p *Service) Refresh(ctx context.Context, raw string, scope string) (User, token.Token, error) {
	user, claims, err := p.Current(ctx, raw, scope)
	if err != nil {
		return User{}, token.Token{}, err
	}
	issued, err := p.signer.Issue(claims.Subject, time.Hour, token.WithScope(scope))
	if err != nil {
		return User{}, token.Token{}, err
	}
	if err = p.sessions.Refresh(ctx, claims, issued.Claims); err != nil {
		return User{}, token.Token{}, err
	}
	return user, issued, nil
}

func (p *Service) Revoke(ctx context.Context, raw string, scope string) error {
	claims, err := p.verify(ctx, raw, scope)
	if err != nil {
		return err
	}
	return p.sessions.Revoke(ctx, claims)
}

func (p *Service) verify(ctx context.Context, raw string, scope string) (token.Claims, error) {
	if p.signer == nil || p.sessions == nil {
		return token.Claims{}, ErrUnavailable
	}
	claims, err := p.signer.Verify(raw)
	if err != nil {
		return token.Claims{}, ErrLogin
	}
	if claims.Scope != scope {
		return token.Claims{}, ErrLogin
	}
	alive, err := p.sessions.Alive(ctx, claims)
	if err != nil {
		return token.Claims{}, err
	}
	if !alive {
		return token.Claims{}, ErrLogin
	}
	return claims, nil
}

func (p *Service) ProfileByID(ctx context.Context, id uint64) (User, error) {
	return p.store.ByID(ctx, id)
}

func (p *Service) login(ctx context.Context, user User, account Account, password string, scope string) (User, string, error) {
	ok, err := verifyPassword(password, account.Salt, account.SecretHash)
	if err != nil {
		return User{}, "", err
	}
	if !ok {
		return User{}, "", ErrLogin
	}
	issued, err := p.signer.Issue(strconv.FormatUint(user.ID, 10), time.Hour,
		token.WithScope(scope),
	)
	if err != nil {
		return User{}, "", err
	}
	if err = p.sessions.Keep(ctx, issued.Claims); err != nil {
		return User{}, "", err
	}
	return user, issued.Raw, nil
}
