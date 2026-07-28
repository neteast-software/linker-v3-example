package permission

import http "github.com/neteast-software/go-module/http/gin/linker"

var routes = http.NewRouteSet()

// Routes 返回权限能力拥有的路由声明。
func Routes() []http.Route {
	return routes.Routes()
}
