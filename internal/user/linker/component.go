package user

import (
	"context"
	"fmt"
	"time"

	session "github.com/neteast-software/go-module/acl/session"
	postgresql "github.com/neteast-software/go-module/db/postgresql/linker"
	http "github.com/neteast-software/go-module/http/gin/linker"
	"github.com/neteast-software/go-module/token"
	linker "github.com/neteast-software/linker/v3"

	user "linker-v3-example/internal/user"
	userhttp "linker-v3-example/internal/user/http"
)

// ID 是组件自治声明的稳定身份。
// 其他组件应依赖这个符号，不要重复书写字符串。
const ID linker.ID = "example/user"

type Component struct {
	store   user.Store
	service *user.Service
	config  Config
}

func New() *Component {
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

func (p *Component) Capabilities() linker.Capabilities {
	return linker.Capabilities{
		linker.Need(postgresql.Key()),
		linker.Offer(user.AuthKey(), func() user.Auth {
			return p.service
		}),
		linker.Offer(user.ActorKey(), func() user.Actor {
			return p.service
		}),
		linker.Offer(user.ServiceKey(), func() *user.Service {
			return p.service
		}),
	}
}

func (p *Component) Assets(context.Context, linker.Runtime) ([]linker.Asset, error) {
	assets := []linker.Asset{
		postgresql.Table(&user.User{}, postgresql.Comment("演示用户")),
		postgresql.Table(&user.Account{}, postgresql.Comment("演示用户账号")),
	}
	return append(assets, http.Assets(userhttp.Routes()...)...), nil
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
