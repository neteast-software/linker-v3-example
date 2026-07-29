package vendorauth

import "github.com/golang-jwt/jwt/v5"

// Claims 是开放平台 RS256 token 的固定声明。
type Claims struct {
	jwt.RegisteredClaims
	ClientID     string            `json:"client_id"`
	Scope        string            `json:"scope"`
	Confirmation map[string]string `json:"cnf"`
}
