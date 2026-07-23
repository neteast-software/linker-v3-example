package user

import user "linker-v3-example/internal/user"

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

func newLoginResult(token string, current user.User) loginResult {
	return loginResult{
		Token: token,
		User:  newProfile(current),
	}
}

func newProfile(current user.User) profile {
	return profile{
		ID:       current.ID,
		Username: current.Username,
		Avatar:   current.Avatar,
		Email:    current.Email,
		Phone:    current.Phone,
		Role:     current.Role,
	}
}

func newApplicationProfile(current user.User, application string) profile {
	ret := newProfile(current)
	ret.Application = application
	return ret
}
