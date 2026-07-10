package user

import (
	"context"
	"strconv"
	"time"

	session "github.com/neteast-software/go-module/acl/session"
	"github.com/neteast-software/go-module/token"
	userconstant "linker-v3-example/internal/constant/user"
	usermodel "linker-v3-example/internal/model/user"
)

type Service struct {
	store    Store
	signer   *token.Signer
	sessions *session.Session
}

func NewService(store Store, signer *token.Signer, sessions *session.Session) Service {
	return Service{store: store, signer: signer, sessions: sessions}
}

func (s Service) AdminLogin(ctx context.Context, username string, password string) (usermodel.User, string, error) {
	user, account, err := s.store.ByAdmin(ctx, username)
	if err != nil {
		return usermodel.User{}, "", err
	}
	return s.login(ctx, user, account, password, "console")
}

func (s Service) UserLogin(ctx context.Context, phone string, password string) (usermodel.User, string, error) {
	user, account, err := s.store.ByPhone(ctx, phone)
	if err != nil {
		return usermodel.User{}, "", err
	}
	return s.login(ctx, user, account, password, "front")
}

func (s Service) Profile(ctx context.Context, raw string, scope string) (usermodel.User, error) {
	claims, err := s.signer.Verify(raw)
	if err != nil {
		return usermodel.User{}, userconstant.ErrLogin
	}
	if claims.Scope != scope {
		return usermodel.User{}, userconstant.ErrLogin
	}
	alive, err := s.sessions.Alive(ctx, claims)
	if err != nil {
		return usermodel.User{}, err
	}
	if !alive {
		return usermodel.User{}, userconstant.ErrLogin
	}
	userID, err := strconv.ParseUint(claims.Subject, 10, 64)
	if err != nil {
		return usermodel.User{}, userconstant.ErrLogin
	}
	return s.store.ByID(ctx, userID)
}

func (s Service) ProfileByID(ctx context.Context, id uint64) (usermodel.User, error) {
	return s.store.ByID(ctx, id)
}

func (s Service) login(ctx context.Context, user usermodel.User, account usermodel.Account, password string, scope string) (usermodel.User, string, error) {
	ok, err := verifyPassword(password, account.Salt, account.SecretHash)
	if err != nil {
		return usermodel.User{}, "", err
	}
	if !ok {
		return usermodel.User{}, "", userconstant.ErrLogin
	}
	issued, err := s.signer.Issue(strconv.FormatUint(user.ID, 10), time.Hour,
		token.WithScope(scope),
	)
	if err != nil {
		return usermodel.User{}, "", err
	}
	if err = s.sessions.Keep(ctx, issued.Claims); err != nil {
		return usermodel.User{}, "", err
	}
	return user, issued.Raw, nil
}
