package tts

import (
	"context"

	http "github.com/neteast-software/go-module/http/gin/linker"
	linker "github.com/neteast-software/linker/v3"

	ttsclient "linker-v3-example/internal/tts/client/linker"
)

var routes = http.NewRouteSet()

const ID linker.ID = "example/tts-proxy"

type Component struct{}

func New() *Component {
	return &Component{}
}

func (p *Component) Identity() linker.ID {
	return ID
}

func (p *Component) Capabilities() linker.Capabilities {
	return linker.Capabilities{linker.Need(ttsclient.Key())}
}

func (p *Component) Assets(context.Context, linker.Runtime) ([]linker.Asset, error) {
	return http.Assets(Routes()...), nil
}

// Routes 返回 TTS 能力拥有的路由声明。
func Routes() []http.Route {
	return routes.Routes()
}
