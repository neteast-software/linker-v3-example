package example

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
	"github.com/neteast-software/go-module/acl"
	httpgin "github.com/neteast-software/go-module/http/gin"
	"github.com/neteast-software/go-module/http/gin/middleware"
	oauthcore "github.com/neteast-software/go-module/security/oauth"
	oauthhttp "github.com/neteast-software/go-module/security/oauth/http/gin"
	oauthjwt "github.com/neteast-software/go-module/security/oauth/jwt"
)

const oauthIssuer = "https://identity.example"
const oauthAudience = "linker-v3-example"
const oauthScope = "profile.read"

var oauthKey = []byte("linker-v3-example-oauth-key-32-bytes")

func TestOAuthJWTMiddlewareExample(t *testing.T) {
	verifier := newOAuthVerifier(t, oauthjwt.Static(oauthKey))
	resource := acl.NewResource("http.example.oauth.profile", acl.Scope("example", 1, "OAuth 用户资料", acl.Read))
	access := acl.NewStaticProvider(acl.NewGrant("user-1",
		acl.GrantLevel(1),
		acl.GrantActions(acl.Read),
		acl.GrantGroups(acl.Role("profile-reader", resource)),
	))
	unavailable := newOAuthVerifier(t, oauthjwt.KeyFunc(func(context.Context, oauthjwt.Header) (any, error) {
		return nil, errors.New("key provider unavailable")
	}))

	server := httpgin.NewServer(httpgin.Config{Addr: "127.0.0.1:0"},
		httpgin.Group("oauth",
			httpgin.GET("profile", oauthProfile).With(oauthhttp.Auth(verifier)),
			httpgin.GET("acl-profile", oauthProfile).
				With(middleware.ACL(oauthhttp.New(verifier), access)).
				WithResource(resource),
			httpgin.GET("unavailable", oauthProfile).With(oauthhttp.Auth(unavailable)),
		),
	)
	valid := signOAuthToken(t, "user-1", oauthScope)

	assertOAuthResponse(t, server.Engine(), "/oauth/profile", "", http.StatusUnauthorized)
	assertOAuthResponse(t, server.Engine(), "/oauth/profile", signOAuthToken(t, "user-1", "profile.write"), http.StatusForbidden)
	assertOAuthResponse(t, server.Engine(), "/oauth/profile", valid, http.StatusOK)
	assertOAuthResponse(t, server.Engine(), "/oauth/acl-profile", valid, http.StatusOK)
	assertOAuthResponse(t, server.Engine(), "/oauth/acl-profile", signOAuthToken(t, "user-2", oauthScope), http.StatusForbidden)
	assertOAuthResponse(t, server.Engine(), "/oauth/unavailable", valid, http.StatusServiceUnavailable)
}

func newOAuthVerifier(t *testing.T, key oauthjwt.Key) *oauthcore.Verifier {
	t.Helper()
	provider, err := oauthjwt.New(key, oauthjwt.Algorithm("HS256"))
	if err != nil {
		t.Fatalf("创建 OAuth JWT provider: %v", err)
	}
	verifier, err := oauthcore.New(provider,
		oauthcore.Issuer(oauthIssuer),
		oauthcore.Audience(oauthAudience),
		oauthcore.Scope(oauthScope),
		oauthcore.RequireExpiry(),
	)
	if err != nil {
		t.Fatalf("创建 OAuth verifier: %v", err)
	}
	return verifier
}

func signOAuthToken(t *testing.T, subject string, scope string) string {
	t.Helper()
	now := time.Now().UTC()
	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, jwtlib.MapClaims{
		"iss":   oauthIssuer,
		"sub":   subject,
		"aud":   []string{oauthAudience},
		"scope": scope,
		"iat":   now.Unix(),
		"exp":   now.Add(time.Minute).Unix(),
	})
	value, err := token.SignedString(oauthKey)
	if err != nil {
		t.Fatalf("签发 OAuth 测试 token: %v", err)
	}
	return value
}

func oauthProfile(c *httpgin.Context) {
	claims, ok := oauthhttp.Claims(c)
	if !ok {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, map[string]string{
		"id":     claims.Identity(),
		"issuer": claims.Issuer,
	})
}

func assertOAuthResponse(t *testing.T, handler http.Handler, path string, token string, status int) {
	t.Helper()
	request := httptest.NewRequest(http.MethodGet, path, nil)
	if token != "" {
		request.Header.Set("Authorization", "Bearer "+token)
	}
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != status {
		t.Fatalf("%s status=%d want=%d body=%s", path, response.Code, status, response.Body.String())
	}
	if status != http.StatusOK {
		return
	}
	var body map[string]string
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("解析 OAuth 响应: %v", err)
	}
	if body["id"] != "user-1" || body["issuer"] != oauthIssuer {
		t.Fatalf("OAuth claims 投影错误: %#v", body)
	}
}
