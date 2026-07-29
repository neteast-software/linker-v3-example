package vendorauth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/neteast-software/go-module/http/gateway"

	"linker-v3-example/internal/access"
)

const FactoryID = "vendor-auth"

// Factory 将 endpoint policy 编译为开放平台鉴权 Filter。
type Factory struct {
	verifier *Verifier
	policies map[string]Policy
	err      error
}

// Filter 创建开放平台鉴权 factory；policy 只属于 vendorauth，不进入 Gateway core。
func Filter(verifier *Verifier, policies ...Policy) *Factory {
	p := &Factory{verifier: verifier, policies: make(map[string]Policy, len(policies))}
	for _, policy := range policies {
		if err := policy.validate(); err != nil {
			p.err = err
			continue
		}
		if _, exists := p.policies[policy.name]; exists {
			p.err = fmt.Errorf("开放平台 policy %q 重复", policy.name)
			continue
		}
		p.policies[policy.name] = policy
	}
	return p
}

func (p *Factory) Identity() string {
	return FactoryID
}

func (p *Factory) Build(config map[string]any) (gateway.Filter, error) {
	if p == nil {
		return nil, fmt.Errorf("开放平台鉴权 factory 不能为空")
	}
	if p.err != nil {
		return nil, p.err
	}
	if len(config) != 1 {
		return nil, fmt.Errorf("开放平台鉴权只接受 policy 配置")
	}
	name, ok := config["policy"].(string)
	if !ok {
		return nil, fmt.Errorf("开放平台鉴权 policy 必须是字符串")
	}
	policy, ok := p.policies[name]
	if !ok {
		return nil, fmt.Errorf("开放平台鉴权 policy %q 未声明", name)
	}
	return gateway.Before(FactoryID+"/"+name, func(_ context.Context, request *http.Request) (*http.Response, error) {
		access.ClearIdentity(request.Header)
		raw, ok := strings.CutPrefix(request.Header.Get("Authorization"), "Bearer ")
		if !ok || strings.TrimSpace(raw) != raw || raw == "" {
			return problem(http.StatusUnauthorized, "vendor-token-required", "请提供开放平台访问凭证"), nil
		}
		identity, err := p.verifier.Verify(request, raw, policy)
		if err != nil {
			return problem(http.StatusUnauthorized, "vendor-token-invalid", "开放平台访问凭证未通过校验"), nil
		}
		request.Header.Del("Authorization")
		access.ProjectIdentity(request.Header, identity)
		return nil, nil
	}), nil
}
