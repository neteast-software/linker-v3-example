package permission

import (
	"context"

	linker "github.com/neteast-software/linker/v3"

	permission "linker-v3-example/internal/permission"
	_ "linker-v3-example/internal/permission/http"
)

const ID linker.ID = "example/permission"

type Component struct {
	service *permission.Service
}

func New() *Component {
	return &Component{service: permission.New()}
}

func (p *Component) Identity() linker.ID {
	return ID
}

func (p *Component) OnMounted(_ context.Context, runtime linker.Runtime) error {
	return linker.Provide(runtime, permission.ServiceKey(), p.service)
}
