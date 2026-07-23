package console

import (
	"context"
	"strconv"

	"github.com/neteast-software/go-module/graph/console/login"

	user "linker-v3-example/internal/user"
)

func (p *Provider) Profile(ctx context.Context, subject string) (login.Profile, error) {
	id, err := strconv.ParseUint(subject, 10, 64)
	if err != nil {
		return login.Profile{}, err
	}
	user, err := p.user.ProfileByID(ctx, id)
	if err != nil {
		return login.Profile{}, err
	}
	return profile(user), nil
}

func profile(current user.User) login.Profile {
	return login.Profile{
		ID:       strconv.FormatUint(current.ID, 10),
		Username: current.Username,
		Name:     current.Username,
		Avatar:   current.Avatar,
		Email:    current.Email,
		Phone:    current.Phone,
	}
}
