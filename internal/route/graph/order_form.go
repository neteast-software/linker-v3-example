package graph

import (
	"github.com/neteast-software/go-module/acl"
	"github.com/neteast-software/go-module/graph/naive/behavior"
	"github.com/neteast-software/go-module/graph/naive/button"
	"github.com/neteast-software/go-module/graph/naive/field"
	"github.com/neteast-software/go-module/graph/naive/form"
	graphhttp "github.com/neteast-software/go-module/graph/naive/http/gin"
	"github.com/neteast-software/go-module/graph/naive/types"
	http "github.com/neteast-software/go-module/http/gin/linker"

	routemiddleware "linker-v3-example/internal/route/middleware"
)

const orderFormResource = "http.app2.graph.order.form"
const orderSaveResource = "http.app2.graph.order.save"

func init() {
	http.RegisterIn("api/v1/app2",
		http.Group("graph",
			http.Use(routemiddleware.Application("app2")),
			http.GET("orders/form", orderForm).Resource(
				orderFormResource,
				acl.Scope("app2", 2, "应用二订单 graph 表单"),
			),
		),
	)
}

func orderForm(c *http.Context) {
	data := sampleOrders()[0]
	graph := form.New("订单维护").
		WithResource(orderFormResource).
		Add(
			field.Hidden("id").Value(data.ID),
			field.Text("number", "订单号").Required("请输入订单号").Value(data.Number),
			field.Select("status", "状态").
				Required("请选择状态").
				Items(types.NewItem("进行中", "open"), types.NewItem("已关闭", "closed")).
				Value(data.Status),
			field.Number("amount", "金额").Required("请输入金额").Value(data.Amount),
		).
		AddButton(button.Command("保存", "/api/v1/app2/graph/orders/save",
			button.WithResource(orderSaveResource),
			button.WithBehavior(behavior.Refresh(behavior.WithMessage("保存成功"))),
		)).
		Submit("/api/v1/app2/graph/orders/save").
		WithPrimary("id").
		WithData(data)

	graphhttp.Typed(c, graph)
}
