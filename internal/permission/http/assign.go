package permission

import (
	"github.com/neteast-software/go-module/acl"
	http "github.com/neteast-software/go-module/http/gin/linker"
	"github.com/neteast-software/go-module/http/gin/response"

	access "linker-v3-example/internal/access"
	permission "linker-v3-example/internal/permission"
)

func init() {
	http.RegisterIn("api/v1/app2",
		http.POST("permission/role/:id/resource", assign).
			With(
				access.Console(),
				access.Application("app2"),
			).
			Resource(permission.Assign, acl.Scope("app2", 3, "角色权限分配", acl.Write)),
	)
}

func assign(c *http.Context) {
	var request changeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(c, err, "权限分配参数无效")
		return
	}
	service, ok := require(c)
	if !ok {
		return
	}
	service.Assign(c.Param("id"), request.Resources...)
	response.Success(c)
}
