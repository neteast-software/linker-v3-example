package permission

import (
	"github.com/neteast-software/go-module/acl"
	http "github.com/neteast-software/go-module/http/gin/linker"
	"github.com/neteast-software/go-module/http/gin/response"

	permissionconstant "linker-v3-example/internal/constant/permission"
	routemiddleware "linker-v3-example/internal/route/middleware"
)

func init() {
	http.RegisterIn("api/v1/app2",
		http.GET("permission/role/:id/resource", list).
			With(
				routemiddleware.Console(),
				routemiddleware.Application("app2"),
			).
			Resource(permissionconstant.Read, acl.Scope("app2", 2, "角色权限读取", acl.Read)),
	)
}

func list(c *http.Context) {
	service, ok := require(c)
	if !ok {
		return
	}
	response.Data(c, map[string]any{
		"selected": c.Param("id"),
		"assigned": service.List(c.Param("id")),
	})
}
