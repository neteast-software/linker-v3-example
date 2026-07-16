package order

import (
	"context"

	linker "github.com/neteast-software/linker/v3"

	_ "linker-v3-example/internal/route/order"
	service "linker-v3-example/internal/service/order"
)

const ID linker.ID = "example/order"

type Component struct {
	service *service.Service
}

func New() *Component {
	return &Component{service: service.New()}
}

func (p *Component) Identity() linker.ID {
	return ID
}

func (p *Component) OnMounted(_ context.Context, runtime linker.Runtime) error {
	return linker.Provide(runtime, service.ServiceKey(), p.service)
}
