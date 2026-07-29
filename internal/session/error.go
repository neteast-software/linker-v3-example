package session

import "errors"

var (
	ErrToken    = errors.New("session token 不合法")
	ErrNotFound = errors.New("session 不存在")
	ErrExpired  = errors.New("session 已过期")
)
