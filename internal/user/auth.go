package user

import (
	"context"

	"github.com/neteast-software/go-module/token"
)

// Auth 是 user 组件对认证、会话和用户资料能力的稳定边界。
type Auth interface {
	AdminLogin(context.Context, string, string) (User, string, error)
	Current(context.Context, string, string) (User, token.Claims, error)
	Refresh(context.Context, string, string) (User, token.Token, error)
	Revoke(context.Context, string, string) error
	ProfileByID(context.Context, uint64) (User, error)
}
