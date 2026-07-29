package captcha

import "github.com/neteast-software/go-module/http/gateway/declaration"

// Login 为登录 route 声明一次性验证码影响面。
func Login(route declaration.Route) declaration.Route {
	route.Filters = append(route.Filters, declaration.Filter{
		External: &declaration.External{Factory: FactoryID},
	})
	return route
}
