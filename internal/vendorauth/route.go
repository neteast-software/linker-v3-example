package vendorauth

import "github.com/neteast-software/go-module/http/gateway/declaration"

// Protect 为一条开放平台 route 声明固定 endpoint policy。
func Protect(route declaration.Route, policy string) declaration.Route {
	route.Filters = append(route.Filters, declaration.Filter{
		External: &declaration.External{
			Factory: FactoryID,
			Config:  map[string]any{"policy": policy},
		},
	})
	return route
}
