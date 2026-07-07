package user

import (
	"github.com/neteast-software/go-module/acl"
	http "github.com/neteast-software/go-module/http/gin/linker"
	"github.com/neteast-software/go-module/http/gin/response"
)

func init() {
	http.RegisterIn("system",
		http.POST("login", systemLogin).Resource(
			"http.console.auth.login",
			acl.Scope("console", 0, "后台登录"),
		),
	)
}

func systemLogin(c *http.Context) {
	req, ok := bindLoginRequest(c)
	if !ok {
		return
	}
	svc, ok := service(c)
	if !ok {
		return
	}
	currentUser, token, err := svc.AdminLogin(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		response.Warning(c, "%s", err.Error())
		return
	}
	setUserOperator(c, currentUser)
	response.Data(c, newLoginResult(token, currentUser))
}
