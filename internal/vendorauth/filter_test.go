package vendorauth

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"net/http"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/neteast-software/go-module/http/gateway"

	"linker-v3-example/internal/access"
)

func TestVendorFilterVerifiesFixedClaimsAndProjectsIdentity(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now().UTC().Truncate(time.Second)
	leaf := &x509.Certificate{Raw: []byte("example-client-certificate")}
	digest := sha256.Sum256(leaf.Raw)
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "example-issuer",
			Subject:   "vendor-one",
			Audience:  jwt.ClaimStrings{"example-api"},
			ExpiresAt: jwt.NewNumericDate(now.Add(5 * time.Minute)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        "token-one",
		},
		ClientID: "vendor-one",
		Scope:    "equipment.read equipment.status",
		Confirmation: map[string]string{
			"x5t#S256": base64.RawURLEncoding.EncodeToString(digest[:]),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	raw, err := token.SignedString(key)
	if err != nil {
		t.Fatal(err)
	}

	verifier := RS256(&key.PublicKey, "example-issuer", "example-api")
	verifier.now = func() time.Time { return now.Add(time.Second) }
	factory := Filter(verifier, Scope("equipment-read", "equipment.read", true))
	filter, err := factory.Build(map[string]any{"policy": "equipment-read"})
	if err != nil {
		t.Fatal(err)
	}
	before := filter.(gateway.BeforeFilter)
	request, err := http.NewRequest(http.MethodGet, "http://gateway.local/vendor/equipment", nil)
	if err != nil {
		t.Fatal(err)
	}
	request.TLS = &tls.ConnectionState{PeerCertificates: []*x509.Certificate{leaf}}
	request.Header.Set("Authorization", "Bearer "+raw)
	request.Header.Set(access.HeaderUsername, "伪造用户")

	response, err := before.Before(context.Background(), request)
	if err != nil || response != nil {
		t.Fatalf("before = (%v, %v)", response, err)
	}
	if request.Header.Get("Authorization") != "" {
		t.Fatal("开放平台 token 被转发到 upstream")
	}
	if request.Header.Get(access.HeaderUsername) != "vendor-one" ||
		request.Header.Get(access.HeaderScope) != "equipment.read" ||
		request.Header.Get(access.HeaderSource) != "vendor-auth" {
		t.Fatalf("identity = %#v", request.Header)
	}
}

func TestVendorFilterRejectsUnknownPolicyAndUntrustedIdentity(t *testing.T) {
	factory := Filter(nil, Scope("equipment-read", "equipment.read", false))
	if _, err := factory.Build(map[string]any{"policy": "unknown"}); err == nil {
		t.Fatal("unknown policy was accepted")
	}
	filter, err := factory.Build(map[string]any{"policy": "equipment-read"})
	if err != nil {
		t.Fatal(err)
	}
	request, _ := http.NewRequest(http.MethodGet, "http://gateway.local/vendor/equipment", nil)
	request.Header.Set(access.HeaderUserID, "999")
	response, err := filter.(gateway.BeforeFilter).Before(context.Background(), request)
	if err != nil || response.StatusCode != http.StatusUnauthorized {
		t.Fatalf("response = %#v err=%v", response, err)
	}
	if request.Header.Get(access.HeaderUserID) != "" {
		t.Fatal("untrusted identity was retained")
	}
}
