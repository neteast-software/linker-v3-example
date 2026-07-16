package console

import (
	"context"
	"errors"
	"fmt"

	graphconsole "github.com/neteast-software/go-module/graph/console"
	"github.com/neteast-software/go-module/graph/console/login"

	userservice "linker-v3-example/internal/service/user"
)

func (p *Provider) Login(ctx context.Context, credential login.Credential) (login.Session, error) {
	if credential.Method != login.Password {
		return login.Session{}, fmt.Errorf("当前示例只支持账号密码登录")
	}
	username, ok := text(credential.Values["username"])
	if !ok {
		return login.Session{}, fmt.Errorf("用户名不能为空")
	}
	password, ok := text(credential.Values["password"])
	if !ok {
		return login.Session{}, fmt.Errorf("密码不能为空")
	}
	_, raw, err := p.user.AdminLogin(ctx, username, password)
	if err != nil {
		return login.Session{}, err
	}
	return p.Current(ctx, raw)
}

func text(value any) (string, bool) {
	text, ok := value.(string)
	return text, ok && text != ""
}

func providerError(err error) error {
	if errors.Is(err, userservice.ErrUnavailable) {
		return fmt.Errorf("%w: %v", graphconsole.ErrSessionUnavailable, err)
	}
	return err
}
