package console

import (
	"slices"

	"github.com/neteast-software/go-module/acl"
	graphconsole "github.com/neteast-software/go-module/graph/console"
	"github.com/neteast-software/go-module/graph/console/protocol"
	linker "github.com/neteast-software/linker/v3"

	consoleadapter "linker-v3-example/internal/adapter/console"
	ordercomponent "linker-v3-example/internal/component/order"
	permissioncomponent "linker-v3-example/internal/component/permission"
	usercomponent "linker-v3-example/internal/component/user"
	consoleconstant "linker-v3-example/internal/constant/console"
	orderconstant "linker-v3-example/internal/constant/order"
	permissionconstant "linker-v3-example/internal/constant/permission"
	consoledoc "linker-v3-example/internal/page/console"
	"linker-v3-example/internal/page/dashboard"
	orderpage "linker-v3-example/internal/page/order"
	permissionpage "linker-v3-example/internal/page/permission"
	userservice "linker-v3-example/internal/service/user"
)

func New(user userservice.Auth, options ...graphconsole.Option) *graphconsole.Component {
	provider := consoleadapter.New(user)
	defaults := []graphconsole.Option{
		graphconsole.WithDependencies(
			linker.RequireComponent(usercomponent.ID),
			linker.RequireComponent(ordercomponent.ID),
			linker.RequireComponent(permissioncomponent.ID),
		),
		graphconsole.WithEntry(consoledoc.Entry()),
		graphconsole.WithMenu(consoledoc.Menu()),
		graphconsole.WithPages(map[string]protocol.Object{
			"dashboard":                dashboard.Page(),
			"order.list":               orderpage.List(),
			"order.form":               orderpage.Form(),
			"permission.role-resource": permissionpage.Relation(),
		}),
		graphconsole.WithResources(
			acl.NewResource(consoleconstant.Dashboard, acl.Scope("console", 0, "后台工作台", acl.Read)),
			acl.NewResource(orderconstant.List, acl.Scope("console", 1, "后台订单列表", acl.Read)),
			acl.NewResource(orderconstant.Update, acl.Scope("app2", 2, "应用二订单维护", acl.Read|acl.Update)),
			acl.NewResource(permissionconstant.Manage, acl.Scope("console", 3, "角色权限配置", acl.Read|acl.Update)),
		),
		graphconsole.WithLogin(provider),
		graphconsole.WithSession(provider),
		graphconsole.WithProfile(provider),
		graphconsole.WithAccess(provider),
		graphconsole.WithNotice(provider),
	}
	return graphconsole.New(slices.Concat(defaults, options)...)
}
