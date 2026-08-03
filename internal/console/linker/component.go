package console

import (
	"slices"

	"github.com/neteast-software/go-module/acl"
	graphconsole "github.com/neteast-software/go-module/graph/console/linker"
	"github.com/neteast-software/go-module/graph/console/protocol"

	console "linker-v3-example/internal/console"
	"linker-v3-example/internal/console/dashboard"
	orderresource "linker-v3-example/internal/order"
	orderpage "linker-v3-example/internal/order/console"
	permissionresource "linker-v3-example/internal/permission"
	permissionpage "linker-v3-example/internal/permission/console"
	user "linker-v3-example/internal/user"
)

func New(options ...graphconsole.Option) *graphconsole.Component {
	provider := console.New()
	defaults := []graphconsole.Option{
		graphconsole.ConfigureFrom(user.AuthKey(), provider.Configure),
		graphconsole.WithEntry(console.Entry()),
		graphconsole.WithMenu(console.Menu()),
		graphconsole.WithPages(map[string]protocol.Object{
			"dashboard":                dashboard.Page(),
			"order.list":               orderpage.List(),
			"order.form":               orderpage.Form(),
			"permission.role-resource": permissionpage.Relation(),
		}),
		graphconsole.WithResources(
			acl.NewResource(console.Dashboard, acl.Scope("console", 0, "后台工作台", acl.Read)),
			acl.NewResource(orderresource.List, acl.Scope("console", 1, "后台订单列表", acl.Read)),
			acl.NewResource(orderresource.Update, acl.Scope("app2", 2, "应用二订单维护", acl.Read|acl.Update)),
			acl.NewResource(permissionresource.Manage, acl.Scope("console", 3, "角色权限配置", acl.Read|acl.Update)),
		),
		graphconsole.WithProvider(provider),
	}
	return graphconsole.New(slices.Concat(defaults, options)...)
}
