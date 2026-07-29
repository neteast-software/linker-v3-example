package vendorauth

import (
	"fmt"
	"strings"
)

// Policy 是一个固定 endpoint 对应的开放平台访问范围。
type Policy struct {
	name               string
	scope              string
	requireCertificate bool
}

// Scope 声明 endpoint policy 名称、OAuth scope 和证书绑定要求。
func Scope(name, scope string, requireCertificate bool) Policy {
	return Policy{name: name, scope: scope, requireCertificate: requireCertificate}
}

func (p Policy) validate() error {
	if strings.TrimSpace(p.name) != p.name || p.name == "" {
		return fmt.Errorf("开放平台 policy 必须声明名称")
	}
	if strings.TrimSpace(p.scope) != p.scope || p.scope == "" {
		return fmt.Errorf("开放平台 policy %q 必须声明 scope", p.name)
	}
	return nil
}
