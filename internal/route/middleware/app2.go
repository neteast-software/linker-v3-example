package middleware

import (
	applicationcore "github.com/neteast-software/go-module/application"
	apphttp "github.com/neteast-software/go-module/application/http/gin"
	applicationcomponent "github.com/neteast-software/go-module/application/linker"
	http "github.com/neteast-software/go-module/http/gin/linker"
)

func Application(scope string) http.HandlerFunc {
	return apphttp.New(nil,
		apphttp.WithScope(scope),
		apphttp.Required(true),
		apphttp.WithStoreProvider(func(c *http.Context) applicationcore.Store {
			store, ok := http.Resolve(c, applicationcomponent.StoreKey())
			if !ok {
				return nil
			}
			return store
		}),
	)
}

func CurrentApplication(c *http.Context) (applicationcore.Application, bool) {
	return apphttp.Current(c)
}
