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
		http.GET("permission/role/:id/resource", list).
			With(
				access.Console(),
				access.Application("app2"),
			).
			Resource(permission.Read, acl.Scope("app2", 2, "角色权限读取", acl.Read)),
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
