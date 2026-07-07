package inspection

import (
	"context"

	postgresql "github.com/neteast-software/go-module/db/postgresql/linker"
	linker "github.com/neteast-software/linker/v3"
	"gorm.io/gorm"

	inspectionmodel "linker-v3-example/internal/model/inspection"
	_ "linker-v3-example/internal/route/inspection" // route 声明随组件进入编译
	inspectionservice "linker-v3-example/internal/service/inspection"
)

const ID linker.ID = "example/inspection"

type Component struct {
	store   inspectionservice.Store
	service inspectionservice.Service
}

func NewComponent() *Component {
	return &Component{}
}

func (p *Component) Identity() linker.ID {
	return ID
}

func (p *Component) Dependencies() []linker.Dependency {
	return []linker.Dependency{linker.RequireID(postgresql.ID)}
}

func (p *Component) Assets(context.Context, linker.Runtime) ([]linker.Asset, error) {
	return []linker.Asset{
		postgresql.Table(&inspectionmodel.Task{}, postgresql.Comment("演示巡检任务")),
		postgresql.Table(&inspectionmodel.Archive{}, postgresql.Comment("外部巡检归档"), postgresql.External()),
	}, nil
}

func (p *Component) Init(ctx context.Context, runtime linker.Runtime) error {
	db, err := linker.RequireCapability(runtime, linker.NewCapabilityKey[*gorm.DB](postgresql.ID))
	if err != nil {
		return err
	}
	p.store = inspectionservice.NewStore(db)
	p.service = inspectionservice.NewService(p.store)
	return inspectionservice.Seed(ctx, p.store)
}

func (p *Component) OnMounted(_ context.Context, runtime linker.Runtime) error {
	return linker.Provide(runtime, inspectionservice.ServiceKey(), p.service)
}
