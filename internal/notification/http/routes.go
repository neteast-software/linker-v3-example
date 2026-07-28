package notification

import (
	"context"

	http "github.com/neteast-software/go-module/http/gin/linker"
	mq "github.com/neteast-software/go-module/mq/consumer/linker"
	linker "github.com/neteast-software/linker/v3"
)

var routes = http.NewRouteSet()

const ID linker.ID = "example/notification-api"

type Component struct{}

func New() *Component {
	return &Component{}
}

func (p *Component) Identity() linker.ID {
	return ID
}

func (p *Component) Capabilities() linker.Capabilities {
	return linker.Capabilities{linker.Need(mq.ConsumersKey())}
}

func (p *Component) Assets(context.Context, linker.Runtime) ([]linker.Asset, error) {
	return http.Assets(Routes()...), nil
}

// Routes 返回通知能力拥有的路由声明。
func Routes() []http.Route {
	return routes.Routes()
}
