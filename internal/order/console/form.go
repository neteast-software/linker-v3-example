package order

import (
	"github.com/neteast-software/go-module/graph/console/button"
	"github.com/neteast-software/go-module/graph/console/field"
	"github.com/neteast-software/go-module/graph/console/form"
	"github.com/neteast-software/go-module/graph/console/permission"
	"github.com/neteast-software/go-module/graph/console/protocol"

	order "linker-v3-example/internal/order"
)

func Form() *form.Form {
	target := protocol.Route("/api/v1/app2/order").WithMethod(protocol.PUT)
	return form.New("订单维护").
		Identity("order.form").
		Resource(order.Update).
		Describe("演示 Graph Console 表单与业务 route 的职责边界").
		Add(
			field.Hidden("id").WithValue(uint64(1)),
			field.Text("number", "订单号").Required("请输入订单号").Length(4, 64).
				WithValue("NO-20260716-001"),
			field.Select("status", "状态").Required("请选择状态").Items(
				protocol.Choice("进行中", "open"),
				protocol.Choice("已关闭", "closed"),
			).WithValue("open"),
			field.Number("amount", "金额").Range(0, 1_000_000).WithStep(100).
				WithValue(int64(12800)),
		).
		Buttons(
			button.Do(
				"保存",
				target,
				button.Resource(order.Update, permission.Update),
				button.WithStatus(button.Primary),
			),
		).
		Save(target).
		Key("id").
		Guard()
}
