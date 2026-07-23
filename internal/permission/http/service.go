package permission

import (
	http "github.com/neteast-software/go-module/http/gin/linker"
	"github.com/neteast-software/go-module/http/gin/response"

	permission "linker-v3-example/internal/permission"
)

func require(c *http.Context) (*permission.Service, bool) {
	value, err := http.Require(c, permission.ServiceKey())
	if err != nil {
		response.Error(c, err, "权限配置服务暂时不可用")
		return nil, false
	}
	return value, true
}
