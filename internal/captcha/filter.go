package captcha

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/neteast-software/go-module/http/gateway"
)

const FactoryID = "captcha"

// Factory 把登录入口的验证码声明编译为一次性校验 Filter。
type Factory struct {
	store Store
	salt  string
	now   func() time.Time
}

// Filter 创建验证码 factory；salt 应来自 secret provider。
func Filter(store Store, salt string) *Factory {
	return &Factory{store: store, salt: salt, now: time.Now}
}

func (p *Factory) Identity() string {
	return FactoryID
}

func (p *Factory) Build(config map[string]any) (gateway.Filter, error) {
	if len(config) != 0 {
		return nil, fmt.Errorf("验证码 Filter 不接受 route 配置")
	}
	return gateway.Before(FactoryID, func(ctx context.Context, request *http.Request) (*http.Response, error) {
		id := request.Header.Get("X-Captcha-Challenge")
		answer := request.Header.Get("X-Captcha-Answer")
		request.Header.Del("X-Captcha-Challenge")
		request.Header.Del("X-Captcha-Answer")
		if strings.TrimSpace(id) != id || id == "" || len(id) > 128 ||
			answer == "" || len(answer) > 256 {
			return problem(http.StatusBadRequest, "captcha-required", "请完成验证码后再提交登录请求"), nil
		}
		if p == nil || p.store == nil || len(p.salt) < 16 {
			return problem(http.StatusServiceUnavailable, "captcha-unavailable", "验证码服务暂时不可用，请稍后重试"), nil
		}
		challenge, err := p.store.Take(ctx, id)
		if err != nil || challenge.verify(answer, p.salt, p.now()) != nil {
			return problem(http.StatusBadRequest, "captcha-invalid", "验证码不正确或已经失效"), nil
		}
		return nil, nil
	}), nil
}
