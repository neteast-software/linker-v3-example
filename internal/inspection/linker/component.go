package inspection

import (
	"context"

	postgresql "github.com/neteast-software/go-module/db/postgresql/linker"
	linker "github.com/neteast-software/linker/v3"

	inspection "linker-v3-example/internal/inspection"
	_ "linker-v3-example/internal/inspection/http" // route 声明随组件进入编译
)

const ID linker.ID = "example/inspection"

type Component struct {
	store   inspection.Store
	service inspection.Service
}

func NewComponent() *Component {
	return &Component{}
}

func (p *Component) Identity() linker.ID {
	return ID
}

func (p *Component) Dependencies() []linker.Dependency {
	return []linker.Dependency{linker.RequireComponent(postgresql.ID)}
}

func (p *Component) Assets(context.Context, linker.Runtime) ([]linker.Asset, error) {
	return []linker.Asset{
		postgresql.Table(&inspection.Task{}, postgresql.Comment("演示巡检任务")),
		postgresql.Table(&inspection.Archive{}, postgresql.Comment("外部巡检归档"), postgresql.External()),
	}, nil
}

func (p *Component) Init(ctx context.Context, runtime linker.Runtime) error {
	db, err := postgresql.Require(runtime)
	if err != nil {
		return err
	}
	p.store = inspection.NewStore(db)
	p.service = inspection.NewService(p.store)
	return inspection.Seed(ctx, p.store)
}

func (p *Component) OnMounted(_ context.Context, runtime linker.Runtime) error {
	return linker.Provide(runtime, inspection.ServiceKey(), p.service)
}
