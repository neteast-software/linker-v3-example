package user

import http "github.com/neteast-software/go-module/http/gin/linker"

var routes = http.NewRouteSet()

// Routes 返回用户能力拥有的路由声明。
func Routes() []http.Route {
	return routes.Routes()
}
