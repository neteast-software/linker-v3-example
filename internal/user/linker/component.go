package user

import (
	"context"
	"fmt"
	"time"

	session "github.com/neteast-software/go-module/acl/session"
	postgresql "github.com/neteast-software/go-module/db/postgresql/linker"
	"github.com/neteast-software/go-module/token"
	linker "github.com/neteast-software/linker/v3"

	user "linker-v3-example/internal/user"
	_ "linker-v3-example/internal/user/http" // route 声明随组件进入编译
)

// ID 是组件自治声明的稳定身份。
// 其他组件应依赖这个符号，不要重复书写字符串。
const ID linker.ID = "example/user"

type Component struct {
	store   user.Store
	service *user.Service
	config  Config
}

func NewComponent() *Component {
	return &Component{service: user.New()}
}

func (p *Component) Configs() []linker.Config {
	return []linker.Config{linker.Restart(Namespace)}
}

func (p *Component) Bootstrap(_ context.Context, boot linker.BootstrapContext) error {
	content, ok := boot.Seed.Lookup(Namespace)
	if !ok {
		return fmt.Errorf("缺少 %s 配置", Namespace)
	}
	config, err := decodeConfig(content)
	if err != nil {
		return err
	}
	p.config = config
	return nil
}

func (p *Component) Identity() linker.ID {
	return ID
}

func (p *Component) Dependencies() []linker.Dependency {
	return []linker.Dependency{linker.RequireComponent(postgresql.ID)}
}

func (p *Component) Assets(context.Context, linker.Runtime) ([]linker.Asset, error) {
	return []linker.Asset{
		postgresql.Table(&user.User{}, postgresql.Comment("演示用户")),
		postgresql.Table(&user.Account{}, postgresql.Comment("演示用户账号")),
	}, nil
}

func (p *Component) Init(ctx context.Context, runtime linker.Runtime) error {
	db, err := postgresql.Require(runtime)
	if err != nil {
		return err
	}
	p.store = user.NewStore(db)
	p.service.Configure(
		p.store,
		token.NewHMAC([]byte(p.config.TokenKey)),
		session.New(session.NewMemoryStore(time.Now)),
	)
	password := p.config.SeedPassword
	p.config.SeedPassword = ""
	if password == "" {
		return nil
	}
	return user.Seed(ctx, p.store, password)
}

func (p *Component) OnMounted(_ context.Context, runtime linker.Runtime) error {
	if err := linker.Provide(runtime, user.AuthKey(), user.Auth(p.service)); err != nil {
		return err
	}
	return linker.Provide(runtime, user.ServiceKey(), p.service)
}

func (p *Component) Service() *user.Service {
	return p.service
}
