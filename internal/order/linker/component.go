package order

import (
	"context"

	linker "github.com/neteast-software/linker/v3"

	order "linker-v3-example/internal/order"
	_ "linker-v3-example/internal/order/http"
)

const ID linker.ID = "example/order"

type Component struct {
	service *order.Service
}

func New() *Component {
	return &Component{service: order.New()}
}

func (p *Component) Identity() linker.ID {
	return ID
}

func (p *Component) OnMounted(_ context.Context, runtime linker.Runtime) error {
	return linker.Provide(runtime, order.ServiceKey(), p.service)
}
