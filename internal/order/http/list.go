package order

import (
	"github.com/neteast-software/go-module/acl"
	http "github.com/neteast-software/go-module/http/gin/linker"
	"github.com/neteast-software/go-module/http/gin/response"

	access "linker-v3-example/internal/access"
	order "linker-v3-example/internal/order"
)

func init() {
	http.RegisterIn("api/v1/app2",
		http.GET("order", list).
			With(
				access.Console(),
				access.Application("app2"),
			).
			Resource(order.List, acl.Scope("console", 1, "后台订单列表", acl.Read)),
	)
}

func list(c *http.Context) {
	service, ok := require(c)
	if !ok {
		return
	}
	response.Data(c, service.List())
}
