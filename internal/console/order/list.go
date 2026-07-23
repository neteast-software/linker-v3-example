package order

import (
	"github.com/neteast-software/go-module/graph/console/behavior"
	"github.com/neteast-software/go-module/graph/console/button"
	"github.com/neteast-software/go-module/graph/console/field"
	"github.com/neteast-software/go-module/graph/console/filter"
	"github.com/neteast-software/go-module/graph/console/permission"
	"github.com/neteast-software/go-module/graph/console/protocol"
	"github.com/neteast-software/go-module/graph/console/viewer"

	order "linker-v3-example/internal/order"
)

func List() *viewer.Viewer[order.Order] {
	return viewer.New[order.Order]("订单列表").
		Identity("order.list").
		Resource(order.List).
		Columns(
			viewer.Col("id", "ID").WithWidth(90).Sort(),
			viewer.Col("number", "订单号").Sort(),
			viewer.Col("status", "状态"),
			viewer.Col("amount", "金额"),
		).
		Query(filter.New(
			filter.Field(field.Text("number", "订单号").Hint("输入订单号"), filter.Where("number", filter.Like)),
			filter.Field(field.Select("status", "状态").Items(
				protocol.Choice("进行中", "open"),
				protocol.Choice("已关闭", "closed"),
			), filter.Where("status", filter.Equal)),
		).InLine()).
		Buttons(
			button.Open(
				"维护订单",
				protocol.Native("order.form"),
				button.Resource(order.Update, permission.Update),
				button.WithStatus(button.Primary),
			),
		).
		Rows(
			button.Run(
				"编辑",
				behavior.Redirect(protocol.Native("order.form")).From("id"),
				button.Resource(order.Update, permission.Update),
				button.WithKeys("id"),
			),
		).
		Paging(1, 30, 2).
		WithData(
			order.Order{ID: 1, Number: "NO-20260716-001", Status: "open", Amount: 12800},
			order.Order{ID: 2, Number: "NO-20260716-002", Status: "closed", Amount: 6800},
		)
}
