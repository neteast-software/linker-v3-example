package session

import "github.com/neteast-software/go-module/http/gateway/declaration"

// Protect 为普通业务 route 声明 session 影响面。
func Protect(route declaration.Route, scope, platform string) declaration.Route {
	config := make(map[string]any, 2)
	if scope != "" {
		config["scope"] = scope
	}
	if platform != "" {
		config["platform"] = platform
	}
	route.Filters = append(route.Filters, declaration.Filter{
		External: &declaration.External{Factory: FactoryID, Config: config},
	})
	return route
}

// Socket 为 WebSocket route 同时声明身份与 Upgrade 条件。
func Socket(route declaration.Route, pathPrefix string) declaration.Route {
	route.Filters = append(route.Filters, declaration.Filter{
		External: &declaration.External{
			Factory: FactoryID,
			Config:  map[string]any{"websocket_path": pathPrefix},
		},
	})
	return route
}

// UploadRoute 为旧文件上传入口声明隔离的 token policy。
func UploadRoute(route declaration.Route) declaration.Route {
	route.Filters = append(route.Filters, declaration.Filter{
		External: &declaration.External{Factory: UploadFactoryID},
	})
	return route
}
