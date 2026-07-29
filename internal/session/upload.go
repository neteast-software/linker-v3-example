package session

import (
	"context"
	"fmt"
	"net/http"

	"github.com/neteast-software/go-module/http/gateway"

	"linker-v3-example/internal/access"
)

const UploadFactoryID = "upload-token"

// UploadFactory 隔离历史文件上传凭证，不扩张普通 session Filter。
type UploadFactory struct {
	session *Factory
}

// Upload 创建旧文件上传 token 的命名兼容策略。
func Upload(verifier *Verifier, store Store) *UploadFactory {
	return &UploadFactory{session: Filter(verifier, store)}
}

func (p *UploadFactory) Identity() string {
	return UploadFactoryID
}

func (p *UploadFactory) Build(config map[string]any) (gateway.Filter, error) {
	if len(config) != 0 {
		return nil, fmt.Errorf("旧文件上传 token 策略不接受配置")
	}
	return gateway.Before(UploadFactoryID, func(ctx context.Context, request *http.Request) (*http.Response, error) {
		access.ClearIdentity(request.Header)
		if request.Method != http.MethodPost {
			return problem(http.StatusMethodNotAllowed, "upload-method-invalid", "文件上传入口只接受 POST 请求"), nil
		}
		raw := request.Header.Get("X-Upload-Token")
		request.Header.Del("X-Upload-Token")
		identity, err := p.session.authenticate(ctx, raw, policy{})
		if err != nil {
			return problem(http.StatusUnauthorized, "upload-token-invalid", "文件上传凭证无效或已经过期"), nil
		}
		access.ProjectIdentity(request.Header, identity)
		return nil, nil
	}), nil
}
