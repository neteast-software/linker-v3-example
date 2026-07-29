package session

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/neteast-software/go-module/http/gateway"

	"linker-v3-example/internal/access"
)

func TestSessionFilterUsesRevocableStateAndProjectsIdentity(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	key := []byte("example-session-key-with-at-least-32-bytes")
	verifier := Legacy(key, "example-login", "example-api")
	verifier.now = func() time.Time { return now }
	store := Memory()
	store.now = func() time.Time { return now }
	current := Session{
		ID:        "session-one",
		UserID:    7,
		Username:  "linfun",
		Platform:  "console",
		Source:    "password",
		Scope:     "profile.read order.read",
		ExpiresAt: now.Add(5 * time.Minute),
	}
	if err := store.Save(context.Background(), current); err != nil {
		t.Fatal(err)
	}
	raw := signSessionToken(t, key, now, current)
	factory := Filter(verifier, store)
	filter, err := factory.Build(map[string]any{"platform": "console", "scope": "profile.read"})
	if err != nil {
		t.Fatal(err)
	}
	request, _ := http.NewRequest(http.MethodGet, "http://gateway.local/profile", nil)
	request.Header.Set("Authorization", "Bearer "+raw)
	request.Header.Set(access.HeaderUserID, "999")

	response, err := filter.(gateway.BeforeFilter).Before(context.Background(), request)
	if err != nil || response != nil {
		t.Fatalf("before = (%v, %v)", response, err)
	}
	if request.Header.Get(access.HeaderUserID) != "7" ||
		request.Header.Get(access.HeaderScope) != "profile.read" ||
		request.Header.Get("Authorization") != "" {
		t.Fatalf("identity = %#v", request.Header)
	}

	if err = store.Revoke(context.Background(), current.ID); err != nil {
		t.Fatal(err)
	}
	request, _ = http.NewRequest(http.MethodGet, "http://gateway.local/profile", nil)
	request.Header.Set("Authorization", "Bearer "+raw)
	response, err = filter.(gateway.BeforeFilter).Before(context.Background(), request)
	if err != nil || response.StatusCode != http.StatusUnauthorized {
		t.Fatalf("revoked response = %#v err=%v", response, err)
	}
}

func TestWebSocketPolicyRequiresPathAndUpgrade(t *testing.T) {
	factory := Filter(nil, Memory())
	filter, err := factory.Build(map[string]any{"websocket_path": "/socket/"})
	if err != nil {
		t.Fatal(err)
	}
	request, _ := http.NewRequest(http.MethodGet, "http://gateway.local/socket/events", nil)
	response, err := filter.(gateway.BeforeFilter).Before(context.Background(), request)
	if err != nil || response.StatusCode != http.StatusUpgradeRequired {
		t.Fatalf("ordinary response = %#v err=%v", response, err)
	}
	request, _ = http.NewRequest(http.MethodGet, "http://gateway.local/ordinary", nil)
	request.Header.Set("Connection", "Upgrade")
	request.Header.Set("Upgrade", "websocket")
	response, err = filter.(gateway.BeforeFilter).Before(context.Background(), request)
	if err != nil || response.StatusCode != http.StatusUpgradeRequired {
		t.Fatalf("wrong path response = %#v err=%v", response, err)
	}
}

func TestUploadCompatibilityPolicyUsesDedicatedHeader(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	key := []byte("example-session-key-with-at-least-32-bytes")
	verifier := Legacy(key, "example-login", "example-api")
	verifier.now = func() time.Time { return now }
	store := Memory()
	store.now = func() time.Time { return now }
	current := Session{
		ID: "upload-one", UserID: 8, Username: "uploader", Platform: "front",
		Source: "phone", Scope: "file.upload", ExpiresAt: now.Add(time.Minute),
	}
	if err := store.Save(context.Background(), current); err != nil {
		t.Fatal(err)
	}
	raw := signSessionToken(t, key, now, current)
	filter, err := Upload(verifier, store).Build(nil)
	if err != nil {
		t.Fatal(err)
	}
	request, _ := http.NewRequest(http.MethodPost, "http://gateway.local/file/upload/file", nil)
	request.Header.Set("X-Upload-Token", raw)
	response, err := filter.(gateway.BeforeFilter).Before(context.Background(), request)
	if err != nil || response != nil {
		t.Fatalf("upload = (%v, %v)", response, err)
	}
	if request.Header.Get("X-Upload-Token") != "" ||
		request.Header.Get(access.HeaderScope) != "file.upload" {
		t.Fatalf("upload identity = %#v", request.Header)
	}
}

func signSessionToken(t *testing.T, key []byte, now time.Time, current Session) string {
	t.Helper()
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "example-login",
			Subject:   current.Username,
			Audience:  jwt.ClaimStrings{"example-api"},
			ExpiresAt: jwt.NewNumericDate(current.ExpiresAt),
			IssuedAt:  jwt.NewNumericDate(now.Add(-time.Second)),
			ID:        current.ID,
		},
		UserID: current.UserID, Username: current.Username, Platform: current.Platform,
		Source: current.Source, Scope: current.Scope,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	raw, err := token.SignedString(key)
	if err != nil {
		t.Fatal(err)
	}
	return raw
}
