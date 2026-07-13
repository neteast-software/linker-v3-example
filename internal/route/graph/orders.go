package graph

import (
	"github.com/neteast-software/go-module/acl"
	"github.com/neteast-software/go-module/graph/naive/behavior"
	"github.com/neteast-software/go-module/graph/naive/button"
	"github.com/neteast-software/go-module/graph/naive/field"
	"github.com/neteast-software/go-module/graph/naive/filter"
	graphhttp "github.com/neteast-software/go-module/graph/naive/http/gin"
	"github.com/neteast-software/go-module/graph/naive/types"
	"github.com/neteast-software/go-module/graph/naive/viewer"
	http "github.com/neteast-software/go-module/http/gin/linker"

	routemiddleware "linker-v3-example/internal/route/middleware"
)

const ordersResource = "http.app2.graph.orders"

func init() {
	http.RegisterIn("api/v1/app2",
		http.Group("graph",
			http.Use(routemiddleware.Application("app2")),
			http.GET("orders", orders).Resource(
				ordersResource,
				acl.Scope("app2", 1, "应用二订单 graph 列表"),
			),
		),
	)
}

func orders(c *http.Context) {
	rows := sampleOrders()
	graph := viewer.New[order]("订单").
		WithResource(ordersResource).
		AddColumn(
			viewer.NewColumn("id", "ID").SortableColumn().WithWidth(80),
			viewer.NewColumn("number", "订单号").SortableColumn(),
			viewer.NewColumn("status", "状态", types.ColumnText).
				WithDynamic("status", map[viewer.FieldValue]viewer.Condition{
					"open":   {Type: viewer.Success},
					"closed": {Type: viewer.Default},
				}),
			viewer.NewColumn("amount", "金额", types.ColumnText),
		).
		WithFilter(filter.New().
			Add(field.Text("number", "订单号").Placeholder("输入订单号"), filter.Where("number", filter.Like)).
			Add(field.Select("status", "状态").Items(
				types.NewItem("进行中", "open"),
				types.NewItem("已关闭", "closed"),
			), filter.Where("status", filter.Eq)),
		).
		AddButton(button.Refresh("刷新", button.WithResource(refreshResource), button.WithBehavior(behavior.Refresh(behavior.WithMessage("已刷新"))))).
		AddRowButton(button.Redirect("编辑", "/api/v1/app2/graph/orders/form", button.WithResource(orderFormResource), button.WithKeys("id"))).
		WithPage(1, len(rows), int64(len(rows))).
		WithData(rows...)

	graphhttp.Typed(c, graph)
}
