package console

import (
	"github.com/neteast-software/go-module/graph/console/menu"
	"github.com/neteast-software/go-module/graph/console/protocol"

	order "linker-v3-example/internal/order"
	permission "linker-v3-example/internal/permission"
)

func Menu() *menu.Menu {
	return menu.New(
		menu.NativePage("dashboard", "工作台", "dashboard").
			WithResource(Dashboard).
			WithIcon("LayoutDashboard"),
		menu.Section("business", "业务演示",
			menu.Entry("order.list", "订单列表", protocol.Native("order.list")).
				WithResource(order.List).
				WithIcon("List"),
			menu.Entry("order.form", "订单维护", protocol.Native("order.form")).
				WithResource(order.Update).
				WithIcon("SquarePen"),
		),
		menu.Section("permission", "权限演示",
			menu.Entry("permission.role-resource", "角色权限", protocol.Native("permission.role-resource")).
				WithResource(permission.Manage).
				WithIcon("ShieldCheck"),
		),
	).Identity("console.menu")
}
