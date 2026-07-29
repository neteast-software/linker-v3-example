package session

import "github.com/golang-jwt/jwt/v5"

// Claims 是历史业务 token 中仍需明确接管的身份声明。
type Claims struct {
	jwt.RegisteredClaims
	UserID   uint64 `json:"user_id"`
	Username string `json:"username"`
	Platform string `json:"platform"`
	Source   string `json:"source"`
	Scope    string `json:"scope"`
}
