package user

import (
	"context"
	"fmt"
	"time"

	session "github.com/neteast-software/go-module/acl/session"
	postgresql "github.com/neteast-software/go-module/db/postgresql/linker"
	"github.com/neteast-software/go-module/token"
	linker "github.com/neteast-software/linker/v3"

	usermodel "linker-v3-example/internal/model/user"
	_ "linker-v3-example/internal/route/user" // route 声明随组件进入编译
	userservice "linker-v3-example/internal/service/user"
)

// ID is this component's stable identity.
// Other components should depend on this symbol instead of repeating strings.
const ID linker.ID = "example/user"

type Component struct {
	store   userservice.Store
	service userservice.Service
	config  Config
}

func NewComponent() *Component {
	return &Component{}
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
		postgresql.Table(&usermodel.User{}, postgresql.Comment("演示用户")),
		postgresql.Table(&usermodel.Account{}, postgresql.Comment("演示用户账号")),
	}, nil
}

func (p *Component) Init(ctx context.Context, runtime linker.Runtime) error {
	db, err := postgresql.Require(runtime)
	if err != nil {
		return err
	}
	p.store = userservice.NewStore(db)
	p.service = userservice.NewService(
		p.store,
		token.NewHMAC([]byte(p.config.TokenKey)),
		session.New(session.NewMemoryStore(time.Now)),
	)
	return userservice.Seed(ctx, p.store)
}

func (p *Component) OnMounted(_ context.Context, runtime linker.Runtime) error {
	return linker.Provide(runtime, userservice.ServiceKey(), p.service)
}
