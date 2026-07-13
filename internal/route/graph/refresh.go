package graph

import (
	"github.com/neteast-software/go-module/acl"
	"github.com/neteast-software/go-module/graph/naive/behavior"
	graphhttp "github.com/neteast-software/go-module/graph/naive/http/gin"
	http "github.com/neteast-software/go-module/http/gin/linker"

	routemiddleware "linker-v3-example/internal/route/middleware"
)

const refreshResource = "http.app2.graph.refresh"

func init() {
	http.RegisterIn("api/v1/app2",
		http.Group("graph",
			http.Use(routemiddleware.Application("app2")),
			http.GET("refresh", refresh).Resource(
				refreshResource,
				acl.Scope("app2", 1, "应用二 graph 行为"),
			),
		),
	)
}

func refresh(c *http.Context) {
	graphhttp.Behavior(c, behavior.Refresh(behavior.WithMessage("已刷新"), behavior.WithWait(0.2)))
}
