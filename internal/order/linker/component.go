package order

import (
	"context"

	"github.com/neteast-software/go-module/acl"
	graphconsole "github.com/neteast-software/go-module/graph/console/linker"
	http "github.com/neteast-software/go-module/http/gin/linker"
	linker "github.com/neteast-software/linker/v3"

	order "linker-v3-example/internal/order"
	orderconsole "linker-v3-example/internal/order/console"
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
	assets := http.Assets(orderhttp.Routes()...)
	return append(assets,
		graphconsole.PageAsset("order.list", orderconsole.List()),
		graphconsole.PageAsset("order.form", orderconsole.Form()),
		graphconsole.ResourceAsset(acl.NewResource(
			order.List,
			acl.Scope("console", 1, "后台订单列表", acl.Read),
		)),
		graphconsole.ResourceAsset(acl.NewResource(
			order.Update,
			acl.Scope("app2", 2, "应用二订单维护", acl.Read|acl.Update),
		)),
	), nil
}
