package observability

import (
	http "github.com/neteast-software/go-module/http/gin/linker"

	routemiddleware "linker-v3-example/internal/route/middleware"
)

func init() {
	http.Register(http.Use(routemiddleware.Metrics()))
}
