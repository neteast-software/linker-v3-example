package session

import (
	"context"
	"crypto/subtle"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/neteast-software/go-module/http/gateway"

	"linker-v3-example/internal/access"
)

const FactoryID = "session"

// Factory 将 route 的 session policy 编译为身份投影 Filter。
type Factory struct {
	verifier *Verifier
	store    Store
	now      func() time.Time
}

// Filter 创建普通 HTTP 与 WebSocket session factory。
func Filter(verifier *Verifier, store Store) *Factory {
	return &Factory{verifier: verifier, store: store, now: time.Now}
}

func (p *Factory) Identity() string {
	return FactoryID
}

func (p *Factory) Build(config map[string]any) (gateway.Filter, error) {
	policy, err := parsePolicy(config)
	if err != nil {
		return nil, err
	}
	return gateway.Before(FactoryID, func(ctx context.Context, request *http.Request) (*http.Response, error) {
		access.ClearIdentity(request.Header)
		if policy.websocketPath != "" {
			if !strings.HasPrefix(request.URL.Path, policy.websocketPath) ||
				!headerContains(request.Header, "Connection", "upgrade") ||
				!strings.EqualFold(request.Header.Get("Upgrade"), "websocket") {
				return problem(http.StatusUpgradeRequired, "websocket-upgrade-required", "该入口只接受 WebSocket Upgrade 请求"), nil
			}
		}
		raw, ok := strings.CutPrefix(request.Header.Get("Authorization"), "Bearer ")
		if !ok || strings.TrimSpace(raw) != raw || raw == "" {
			return problem(http.StatusUnauthorized, "session-token-required", "请先登录后再访问"), nil
		}
		identity, err := p.authenticate(ctx, raw, policy)
		if err != nil {
			return problem(http.StatusUnauthorized, "session-token-invalid", "登录状态无效或已经过期"), nil
		}
		request.Header.Del("Authorization")
		access.ProjectIdentity(request.Header, identity)
		return nil, nil
	}), nil
}

type policy struct {
	scope         string
	platform      string
	websocketPath string
}

func parsePolicy(config map[string]any) (policy, error) {
	if len(config) > 3 {
		return policy{}, fmt.Errorf("session policy 只接受 scope、platform 和 websocket_path")
	}
	var result policy
	for key, value := range config {
		raw, ok := value.(string)
		if !ok || strings.TrimSpace(raw) != raw {
			return policy{}, fmt.Errorf("session policy %s 必须是无首尾空白的字符串", key)
		}
		switch key {
		case "scope":
			result.scope = raw
		case "platform":
			result.platform = raw
		case "websocket_path":
			if raw == "" || !strings.HasPrefix(raw, "/") {
				return policy{}, fmt.Errorf("websocket_path 必须是绝对路径前缀")
			}
			result.websocketPath = raw
		default:
			return policy{}, fmt.Errorf("session policy 不支持 %q", key)
		}
	}
	return result, nil
}

func (p *Factory) authenticate(ctx context.Context, raw string, policy policy) (access.Identity, error) {
	if p == nil || p.verifier == nil || p.store == nil {
		return access.Identity{}, ErrToken
	}
	claims, err := p.verifier.Verify(raw)
	if err != nil {
		return access.Identity{}, err
	}
	current, err := p.store.Lookup(ctx, claims.ID)
	if err != nil {
		return access.Identity{}, err
	}
	if !current.matches(claims) || current.validate(p.now()) != nil {
		return access.Identity{}, ErrToken
	}
	if policy.platform != "" &&
		subtle.ConstantTimeCompare([]byte(policy.platform), []byte(current.Platform)) != 1 {
		return access.Identity{}, ErrToken
	}
	if policy.scope != "" && !scopeContains(current.Scope, policy.scope) {
		return access.Identity{}, ErrToken
	}
	scope := current.Scope
	if policy.scope != "" {
		scope = policy.scope
	}
	return access.Identity{
		UserID:   current.UserID,
		Username: current.Username,
		Platform: current.Platform,
		Source:   current.Source,
		Scope:    scope,
	}, nil
}

func scopeContains(raw, expected string) bool {
	for _, scope := range strings.Fields(raw) {
		if subtle.ConstantTimeCompare([]byte(scope), []byte(expected)) == 1 {
			return true
		}
	}
	return false
}

func headerContains(header http.Header, name, expected string) bool {
	for _, value := range header.Values(name) {
		for part := range strings.SplitSeq(value, ",") {
			if strings.EqualFold(strings.TrimSpace(part), expected) {
				return true
			}
		}
	}
	return false
}
