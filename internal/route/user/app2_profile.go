package user

import (
	"github.com/neteast-software/go-module/acl"
	http "github.com/neteast-software/go-module/http/gin/linker"
	"github.com/neteast-software/go-module/http/gin/param"
	"github.com/neteast-software/go-module/http/gin/response"

	routemiddleware "linker-v3-example/internal/route/middleware"
)

func init() {
	http.RegisterIn("api/v1/app2",
		http.Group("user/:id",
			http.Use(routemiddleware.Application("app2")),
			http.GET("profile", app2Profile).Resource(
				"http.app2.user.profile",
				acl.Scope("app2", 1, "应用二用户资料"),
			),
		),
	)
}

func app2Profile(c *http.Context) {
	id, err := param.ID(c)
	if err != nil {
		response.Warning(c, "%s", err.Error())
		return
	}
	svc, ok := service(c)
	if !ok {
		return
	}
	currentUser, err := svc.ProfileByID(c.Request.Context(), uint64(id))
	if err != nil {
		response.Warning(c, "%s", err.Error())
		return
	}
	setUserOperator(c, currentUser)
	app, _ := routemiddleware.CurrentApplication(c)
	response.Data(c, newApplicationProfile(currentUser, app.Scope))
}
