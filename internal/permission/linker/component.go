package permission

import (
	"context"

	"github.com/neteast-software/go-module/acl"
	graphconsole "github.com/neteast-software/go-module/graph/console/linker"
	http "github.com/neteast-software/go-module/http/gin/linker"
	linker "github.com/neteast-software/linker/v3"

	permission "linker-v3-example/internal/permission"
	permissionconsole "linker-v3-example/internal/permission/console"
	permissionhttp "linker-v3-example/internal/permission/http"
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

func (p *Component) Capabilities() linker.Capabilities {
	return linker.Capabilities{
		linker.Offer(permission.ServiceKey(), func() *permission.Service {
			return p.service
		}),
	}
}

func (p *Component) Assets(context.Context, linker.Runtime) ([]linker.Asset, error) {
	assets := http.Assets(permissionhttp.Routes()...)
	return append(assets,
		graphconsole.PageAsset("permission.role-resource", permissionconsole.Relation()),
		graphconsole.ResourceAsset(acl.NewResource(
			permission.Manage,
			acl.Scope("console", 3, "角色权限配置", acl.Read|acl.Update),
		)),
	), nil
}
