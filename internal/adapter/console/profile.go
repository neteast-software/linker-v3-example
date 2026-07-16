package console

import (
	"context"
	"strconv"

	"github.com/neteast-software/go-module/graph/console/login"

	usermodel "linker-v3-example/internal/model/user"
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

func profile(user usermodel.User) login.Profile {
	return login.Profile{
		ID:       strconv.FormatUint(user.ID, 10),
		Username: user.Username,
		Name:     user.Username,
		Avatar:   user.Avatar,
		Email:    user.Email,
		Phone:    user.Phone,
	}
}
