package session

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Verifier 固定历史业务 token 的算法和签发边界。
type Verifier struct {
	key      []byte
	issuer   string
	audience string
	now      func() time.Time
}

// Legacy 创建只接受 HS256 的历史业务 token verifier。
func Legacy(key []byte, issuer, audience string) *Verifier {
	return &Verifier{
		key:      append([]byte(nil), key...),
		issuer:   strings.TrimSpace(issuer),
		audience: strings.TrimSpace(audience),
		now:      time.Now,
	}
}

func (p *Verifier) Verify(raw string) (*Claims, error) {
	if p == nil || len(p.key) < 32 || p.issuer == "" || p.audience == "" {
		return nil, ErrToken
	}
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(
		raw,
		claims,
		func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodHS256 || token.Header["typ"] != "JWT" {
				return nil, ErrToken
			}
			return append([]byte(nil), p.key...), nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuer(p.issuer),
		jwt.WithAudience(p.audience),
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt(),
		jwt.WithTimeFunc(p.now),
	)
	if err != nil || token == nil || !token.Valid {
		return nil, errors.Join(ErrToken, err)
	}
	if claims.ID == "" || claims.UserID == 0 || claims.Username == "" ||
		claims.Platform == "" || claims.IssuedAt == nil || claims.ExpiresAt == nil {
		return nil, ErrToken
	}
	return claims, nil
}
