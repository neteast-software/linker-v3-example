package user

import usermodel "linker-v3-example/internal/model/user"

type loginResult struct {
	Token string  `json:"token"`
	User  profile `json:"user"`
}

type profile struct {
	ID          uint64 `json:"id"`
	Username    string `json:"username"`
	Avatar      string `json:"avatar"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Role        string `json:"role"`
	Application string `json:"application,omitempty"`
}

func newLoginResult(token string, user usermodel.User) loginResult {
	return loginResult{
		Token: token,
		User:  newProfile(user),
	}
}

func newProfile(user usermodel.User) profile {
	return profile{
		ID:       user.ID,
		Username: user.Username,
		Avatar:   user.Avatar,
		Email:    user.Email,
		Phone:    user.Phone,
		Role:     user.Role,
	}
}

func newApplicationProfile(user usermodel.User, application string) profile {
	ret := newProfile(user)
	ret.Application = application
	return ret
}
