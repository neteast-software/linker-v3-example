package vendorauth

import (
	"crypto/rsa"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"linker-v3-example/internal/access"
)

const maximumTokenLifetime = 10 * time.Minute

// Verifier 固定开放平台 token 的算法、签发方、受众和时效约束。
type Verifier struct {
	key      *rsa.PublicKey
	issuer   string
	audience string
	now      func() time.Time
}

// RS256 创建只接受 RS256 的开放平台 token verifier。
func RS256(key *rsa.PublicKey, issuer, audience string) *Verifier {
	return &Verifier{
		key:      key,
		issuer:   strings.TrimSpace(issuer),
		audience: strings.TrimSpace(audience),
		now:      time.Now,
	}
}

func (p *Verifier) Verify(request *http.Request, raw string, policy Policy) (access.Identity, error) {
	if p == nil || p.key == nil || p.issuer == "" || p.audience == "" {
		return access.Identity{}, fmt.Errorf("%w: verifier 未配置", ErrToken)
	}
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(
		raw,
		claims,
		func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodRS256 || token.Header["typ"] != "JWT" {
				return nil, ErrToken
			}
			return p.key, nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Alg()}),
		jwt.WithIssuer(p.issuer),
		jwt.WithAudience(p.audience),
		jwt.WithExpirationRequired(),
		jwt.WithNotBeforeRequired(),
		jwt.WithIssuedAt(),
		jwt.WithTimeFunc(p.now),
	)
	if err != nil || token == nil || !token.Valid {
		return access.Identity{}, errors.Join(ErrToken, err)
	}
	if err = validateClaims(claims); err != nil {
		return access.Identity{}, err
	}
	if !containsScope(claims.Scope, policy.scope) {
		return access.Identity{}, ErrScope
	}
	if err = confirmCertificate(request, claims.Confirmation, policy.requireCertificate); err != nil {
		return access.Identity{}, err
	}
	return access.Identity{
		Username: claims.ClientID,
		Platform: "vendor",
		Source:   "vendor-auth",
		Scope:    policy.scope,
	}, nil
}

func validateClaims(claims *Claims) error {
	if claims == nil || claims.ID == "" || claims.Subject == "" || claims.ClientID == "" ||
		claims.IssuedAt == nil || claims.ExpiresAt == nil || claims.NotBefore == nil {
		return ErrToken
	}
	if claims.ExpiresAt.Time.Sub(claims.IssuedAt.Time) <= 0 ||
		claims.ExpiresAt.Time.Sub(claims.IssuedAt.Time) > maximumTokenLifetime ||
		claims.NotBefore.Time.Before(claims.IssuedAt.Time) ||
		!claims.NotBefore.Time.Before(claims.ExpiresAt.Time) {
		return ErrToken
	}
	return nil
}

func containsScope(raw, expected string) bool {
	for _, scope := range strings.Fields(raw) {
		if subtle.ConstantTimeCompare([]byte(scope), []byte(expected)) == 1 {
			return true
		}
	}
	return false
}

func confirmCertificate(request *http.Request, confirmation map[string]string, required bool) error {
	expected := confirmation["x5t#S256"]
	if expected == "" {
		if required {
			return ErrCertificate
		}
		return nil
	}
	if request == nil || request.TLS == nil || len(request.TLS.PeerCertificates) == 0 {
		return ErrCertificate
	}
	digest := sha256.Sum256(request.TLS.PeerCertificates[0].Raw)
	actual := base64.RawURLEncoding.EncodeToString(digest[:])
	if subtle.ConstantTimeCompare([]byte(actual), []byte(expected)) != 1 {
		return ErrCertificate
	}
	return nil
}
