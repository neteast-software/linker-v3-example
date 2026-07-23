package order

import (
	"fmt"
	"strings"

	"github.com/neteast-software/go-module/acl"
	http "github.com/neteast-software/go-module/http/gin/linker"
	"github.com/neteast-software/go-module/http/gin/response"

	access "linker-v3-example/internal/access"
	order "linker-v3-example/internal/order"
)

type saveRequest struct {
	ID     uint64 `json:"id"`
	Number string `json:"number"`
	Status string `json:"status"`
	Amount int64  `json:"amount"`
}

func init() {
	http.RegisterIn("api/v1/app2",
		http.PUT("order", save).
			With(
				access.Console(),
				access.Application("app2"),
			).
			Resource(order.Update, acl.Scope("app2", 2, "应用二订单维护", acl.Update)),
	)
}

func save(c *http.Context) {
	var request saveRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(c, err, "订单参数无效")
		return
	}
	if err := request.Validate(); err != nil {
		response.Warning(c, "%s", err.Error())
		return
	}
	value, ok := require(c)
	if !ok {
		return
	}
	saved, err := value.Save(order.Order(request))
	if err != nil {
		response.Error(c, err, "订单保存失败")
		return
	}
	response.Data(c, saved)
}

func (p saveRequest) Validate() error {
	if p.ID == 0 {
		return fmt.Errorf("订单 ID 不能为空")
	}
	if len(strings.TrimSpace(p.Number)) < 4 {
		return fmt.Errorf("订单号至少需要 4 个字符")
	}
	if p.Status != "open" && p.Status != "closed" {
		return fmt.Errorf("订单状态必须是 open 或 closed")
	}
	if p.Amount < 0 {
		return fmt.Errorf("订单金额不能小于 0")
	}
	return nil
}
