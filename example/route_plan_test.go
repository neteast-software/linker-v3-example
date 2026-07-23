package example

import (
	"strings"
	"testing"

	http "github.com/neteast-software/go-module/http/gin/linker"

	_ "linker-v3-example/internal/user/http"
)

func TestStaticRoutePlanShowsApp2ResourceAndMiddleware(t *testing.T) {
	plan := findRoutePlan(http.RegisteredRoutePlans(), "GET", "/api/v1/app2/user/:id/profile")
	if plan == nil {
		t.Fatalf("missing app2 profile route plan")
	}
	if plan.Resource == nil || plan.Resource.Tag != "http.app2.user.profile" || plan.Resource.Application != "app2" {
		t.Fatalf("unexpected resource: %#v", plan.Resource)
	}
	if !hasRouteHandlerName(plan.Middlewares, "application/http/gin.New") {
		t.Fatalf("missing app2 middleware: %#v", plan.Middlewares)
	}
}

func findRoutePlan(plans []http.RoutePlan, method string, path string) *http.RoutePlan {
	for i := range plans {
		if plans[i].Method == method && plans[i].Path == path {
			return &plans[i]
		}
	}
	return nil
}

func hasRouteHandlerName(names []string, want string) bool {
	for _, name := range names {
		if strings.Contains(name, want) {
			return true
		}
	}
	return false
}
