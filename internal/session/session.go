package session

import (
	"context"
	"time"
)

// Session 是服务端可撤销的登录状态。
type Session struct {
	ID        string    `json:"id"`
	UserID    uint64    `json:"user_id"`
	Username  string    `json:"username"`
	Platform  string    `json:"platform"`
	Source    string    `json:"source"`
	Scope     string    `json:"scope"`
	ExpiresAt time.Time `json:"expires_at"`
}

// Store 管理服务端 session 的保存、读取和撤销。
type Store interface {
	Save(context.Context, Session) error
	Lookup(context.Context, string) (Session, error)
	Revoke(context.Context, string) error
}

func (p Session) validate(now time.Time) error {
	if p.ID == "" || p.UserID == 0 || p.Username == "" || p.Platform == "" || p.ExpiresAt.IsZero() {
		return ErrNotFound
	}
	if !now.Before(p.ExpiresAt) {
		return ErrExpired
	}
	return nil
}

func (p Session) matches(claims *Claims) bool {
	return claims != nil &&
		p.ID == claims.ID &&
		p.UserID == claims.UserID &&
		p.Username == claims.Username &&
		p.Platform == claims.Platform &&
		p.Source == claims.Source &&
		p.Scope == claims.Scope
}
