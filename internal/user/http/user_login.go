package user

import (
	"github.com/neteast-software/go-module/acl"
	http "github.com/neteast-software/go-module/http/gin/linker"
	"github.com/neteast-software/go-module/http/gin/response"
)

func init() {
	http.RegisterIn("user",
		http.POST("login", userLogin).Resource(
			"http.front.auth.login",
			acl.Scope("front", 0, "用户登录"),
		),
	)
}

func userLogin(c *http.Context) {
	req, ok := bindLoginRequest(c)
	if !ok {
		return
	}
	svc, ok := service(c)
	if !ok {
		return
	}
	currentUser, token, err := svc.UserLogin(c.Request.Context(), req.Phone, req.Password)
	if err != nil {
		response.Warning(c, "%s", err.Error())
		return
	}
	setUserOperator(c, currentUser)
	response.Data(c, newLoginResult(token, currentUser))
}
