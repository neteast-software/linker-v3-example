package dashboard

import (
	"github.com/neteast-software/go-module/graph/console/behavior"
	"github.com/neteast-software/go-module/graph/console/chart"
	"github.com/neteast-software/go-module/graph/console/layout"
	"github.com/neteast-software/go-module/graph/console/permission"
	"github.com/neteast-software/go-module/graph/console/protocol"
	"github.com/neteast-software/go-module/graph/console/viewer"

	console "linker-v3-example/internal/console"
	order "linker-v3-example/internal/order"
)

func Page() *layout.Layout {
	summary := viewer.New[metric]("运行概览").
		Identity("dashboard.summary").
		Resource(console.Dashboard).
		Columns(
			viewer.Col("name", "指标"),
			viewer.Col("value", "数量"),
		).
		WithData(
			metric{Name: "运行组件", Value: 12},
			metric{Name: "业务路由", Value: 18},
			metric{Name: "待处理通知", Value: 1},
		)
	trend := chart.New("请求趋势", chart.Line,
		chart.Data(
			chart.Dim("time", "时间", chart.Time),
			chart.Dim("value", "请求数", chart.Number),
		).Rows(
			chart.Row{"time": "08:00", "value": 42},
			chart.Row{"time": "10:00", "value": 78},
			chart.Row{"time": "12:00", "value": 66},
		),
	).
		Identity("dashboard.request-trend").
		Resource(console.Dashboard).
		Encode(chart.XY("time", "value")).
		Axes(chart.CategoryAxis("时间"), chart.ValueAxis("请求数")).
		Interact(
			chart.On(chart.Click, behavior.Redirect(protocol.Native("order.list"))).
				Protect(order.List, permission.Read),
		)
	value := layout.Flowing(
		"Linker v3 工作台",
		layout.Place("summary", summary).Min(320),
		layout.Place("trend", trend).Min(520),
	).Identity("dashboard").AtMost(2)
	value.Node.WithResource(console.Dashboard)
	return value
}
