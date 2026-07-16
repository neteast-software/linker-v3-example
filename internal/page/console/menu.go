package console

import (
	"github.com/neteast-software/go-module/graph/console/menu"
	"github.com/neteast-software/go-module/graph/console/protocol"

	consoleconstant "linker-v3-example/internal/constant/console"
	orderconstant "linker-v3-example/internal/constant/order"
	permissionconstant "linker-v3-example/internal/constant/permission"
)

func Menu() *menu.Menu {
	return menu.New(
		menu.NativePage("dashboard", "工作台", "dashboard").
			WithResource(consoleconstant.Dashboard).
			WithIcon("LayoutDashboard"),
		menu.Section("business", "业务演示",
			menu.Entry("order.list", "订单列表", protocol.Native("order.list")).
				WithResource(orderconstant.List).
				WithIcon("List"),
			menu.Entry("order.form", "订单维护", protocol.Native("order.form")).
				WithResource(orderconstant.Update).
				WithIcon("SquarePen"),
		),
		menu.Section("permission", "权限演示",
			menu.Entry("permission.role-resource", "角色权限", protocol.Native("permission.role-resource")).
				WithResource(permissionconstant.Manage).
				WithIcon("ShieldCheck"),
		),
	).Identity("console.menu")
}
