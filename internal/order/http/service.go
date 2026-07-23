package order

import (
	http "github.com/neteast-software/go-module/http/gin/linker"
	"github.com/neteast-software/go-module/http/gin/response"

	order "linker-v3-example/internal/order"
)

func require(c *http.Context) (*order.Service, bool) {
	value, err := http.Require(c, order.ServiceKey())
	if err != nil {
		response.Error(c, err, "订单服务暂时不可用")
		return nil, false
	}
	return value, true
}
