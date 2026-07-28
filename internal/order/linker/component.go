package order

import (
	"context"

	http "github.com/neteast-software/go-module/http/gin/linker"
	linker "github.com/neteast-software/linker/v3"

	order "linker-v3-example/internal/order"
	orderhttp "linker-v3-example/internal/order/http"
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

func (p *Component) Capabilities() linker.Capabilities {
	return linker.Capabilities{
		linker.Offer(order.ServiceKey(), func() *order.Service {
			return p.service
		}),
	}
}

func (p *Component) Assets(context.Context, linker.Runtime) ([]linker.Asset, error) {
	return http.Assets(orderhttp.Routes()...), nil
}
