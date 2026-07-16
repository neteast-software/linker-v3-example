package order

import (
	"github.com/neteast-software/go-module/acl"
	http "github.com/neteast-software/go-module/http/gin/linker"
	"github.com/neteast-software/go-module/http/gin/response"

	orderconstant "linker-v3-example/internal/constant/order"
	routemiddleware "linker-v3-example/internal/route/middleware"
)

func init() {
	http.RegisterIn("api/v1/app2",
		http.GET("order", list).
			With(
				routemiddleware.Console(),
				routemiddleware.Application("app2"),
			).
			Resource(orderconstant.List, acl.Scope("console", 1, "后台订单列表", acl.Read)),
	)
}

func list(c *http.Context) {
	service, ok := require(c)
	if !ok {
		return
	}
	response.Data(c, service.List())
}
