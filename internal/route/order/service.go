package order

import (
	http "github.com/neteast-software/go-module/http/gin/linker"
	"github.com/neteast-software/go-module/http/gin/response"

	service "linker-v3-example/internal/service/order"
)

func require(c *http.Context) (*service.Service, bool) {
	value, err := http.Require(c, service.ServiceKey())
	if err != nil {
		response.Error(c, err, "订单服务暂时不可用")
		return nil, false
	}
	return value, true
}
