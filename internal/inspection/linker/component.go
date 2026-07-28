package inspection

import (
	"context"

	postgresql "github.com/neteast-software/go-module/db/postgresql/linker"
	http "github.com/neteast-software/go-module/http/gin/linker"
	linker "github.com/neteast-software/linker/v3"

	inspection "linker-v3-example/internal/inspection"
	inspectionhttp "linker-v3-example/internal/inspection/http"
	user "linker-v3-example/internal/user"
)

const ID linker.ID = "example/inspection"

type Component struct {
	store   inspection.Store
	service inspection.Service
}

func New() *Component {
	return &Component{}
}

func (p *Component) Identity() linker.ID {
	return ID
}

func (p *Component) Capabilities() linker.Capabilities {
	return linker.Capabilities{
		linker.Need(postgresql.Key()),
		linker.Need(user.ActorKey()),
		linker.Offer(inspection.ServiceKey(), func() inspection.Service {
			return p.service
		}),
	}
}

func (p *Component) Assets(context.Context, linker.Runtime) ([]linker.Asset, error) {
	assets := []linker.Asset{
		postgresql.Table(&inspection.Task{}, postgresql.Comment("演示巡检任务")),
		postgresql.Table(&inspection.Archive{}, postgresql.Comment("外部巡检归档"), postgresql.External()),
	}
	return append(assets, http.Assets(inspectionhttp.Routes()...)...), nil
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
