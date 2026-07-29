package gateway

import (
	"time"

	"github.com/neteast-software/go-module/http/gateway/declaration"
)

// Document 创建一个完整、有序的 Linker Gateway 声明。
func Document(routes ...declaration.Route) declaration.Document {
	return declaration.Document{
		Schema: declaration.Schema,
		Plan: declaration.Plan{
			Routes: append([]declaration.Route(nil), routes...),
		},
	}
}

// URL 声明固定 HTTP origin 路由。
func URL(id, path, origin string) declaration.Route {
	return declaration.Route{
		ID:       id,
		Match:    declaration.Match{Paths: []string{path}},
		Upstream: declaration.Upstream{URL: origin},
		Timeout:  declaration.NewTimeout(5 * time.Second),
	}
}

// Service 声明通过中立 discovery capability 解析的路由。
func Service(id, path, service string) declaration.Route {
	return declaration.Route{
		ID:       id,
		Match:    declaration.Match{Paths: []string{path}},
		Upstream: declaration.Upstream{Service: service},
		Timeout:  declaration.NewTimeout(5 * time.Second),
	}
}

// Strip 在当前路由末尾追加显式路径前缀删除。
func Strip(route declaration.Route, count int) declaration.Route {
	route.Filters = append(route.Filters, declaration.Filter{StripPrefix: &count})
	return route
}

// StripService 删除 `/service/**` 路由中的 service identity 层。
func StripService(route declaration.Route) declaration.Route {
	route.Filters = append(route.Filters, declaration.Filter{ServicePrefix: true})
	return route
}
