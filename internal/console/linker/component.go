package console

import (
	"slices"

	"github.com/neteast-software/go-module/acl"
	graphconsole "github.com/neteast-software/go-module/graph/console"
	"github.com/neteast-software/go-module/graph/console/protocol"
	linker "github.com/neteast-software/linker/v3"

	console "linker-v3-example/internal/console"
	"linker-v3-example/internal/console/dashboard"
	orderpage "linker-v3-example/internal/console/order"
	permissionpage "linker-v3-example/internal/console/permission"
	orderresource "linker-v3-example/internal/order"
	order "linker-v3-example/internal/order/linker"
	permissionresource "linker-v3-example/internal/permission"
	permission "linker-v3-example/internal/permission/linker"
	userauth "linker-v3-example/internal/user"
	user "linker-v3-example/internal/user/linker"
)

func New(auth userauth.Auth, options ...graphconsole.Option) *graphconsole.Component {
	provider := console.New(auth)
	defaults := []graphconsole.Option{
		graphconsole.WithDependencies(
			linker.RequireComponent(user.ID),
			linker.RequireComponent(order.ID),
			linker.RequireComponent(permission.ID),
		),
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
		graphconsole.WithLogin(provider),
		graphconsole.WithSession(provider),
		graphconsole.WithProfile(provider),
		graphconsole.WithAccess(provider),
		graphconsole.WithNotice(provider),
	}
	return graphconsole.New(slices.Concat(defaults, options)...)
}
