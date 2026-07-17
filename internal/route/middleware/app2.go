package middleware

import (
	applicationcore "github.com/neteast-software/go-module/application"
	"github.com/neteast-software/go-module/application/http/gin"
	applicationcomponent "github.com/neteast-software/go-module/application/linker"
	http "github.com/neteast-software/go-module/http/gin/linker"
)

func Application(scope string) http.HandlerFunc {
	return application.New(nil,
		application.WithScope(scope),
		application.Required(true),
		application.WithStoreProvider(func(c *http.Context) applicationcore.Store {
			store, ok := http.Resolve(c, applicationcomponent.StoreKey())
			if !ok {
				return nil
			}
			return store
		}),
	)
}

func CurrentApplication(c *http.Context) (applicationcore.Application, bool) {
	return application.Current(c)
}
