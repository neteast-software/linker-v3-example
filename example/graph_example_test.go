package example

import (
	"context"
	"testing"
	"time"

	applicationcore "github.com/neteast-software/go-module/application"
	applicationcomponent "github.com/neteast-software/go-module/application/linker"
	graphlinker "github.com/neteast-software/go-module/graph/naive/linker"
	http "github.com/neteast-software/go-module/http/gin/linker"
	server "github.com/neteast-software/go-module/linker/server"
	linker "github.com/neteast-software/linker/v3"

	graphcomponent "linker-v3-example/internal/component/graph"
)

func TestGraphNaiveExample(t *testing.T) {
	app := server.New(
		server.WithMode(linker.Server),
		server.WithShutdownTimeout(3*time.Second),
		server.WithMapSetting(map[linker.Namespace][]byte{
			linker.Namespace(http.ID): []byte(`{"addr":"127.0.0.1:0"}`),
		}),
		server.WithComponents(
			applicationcomponent.New(
				applicationcomponent.WithDBTarget(""),
				applicationcomponent.WithApplications(applicationcore.Application{
					ID:     "app2",
					Scope:  "app2",
					Name:   "应用二",
					Status: applicationcore.StatusEnabled,
				}),
			),
			graphcomponent.NewComponent(),
		),
	)

	plan := app.Plan()
	if !planHasComponent(plan, graphcomponent.ID) ||
		!planHasRouteAsset(plan, "GET", "/api/v1/app2/graph/orders", "http.app2.graph.orders") ||
		!planHasRouteAsset(plan, "GET", "/api/v1/app2/graph/orders/form", "http.app2.graph.order.form") ||
		!planHasRouteAsset(plan, "GET", "/api/v1/app2/graph/refresh", "http.app2.graph.refresh") {
		t.Fatalf("plan missing graph example assets: components=%#v assets=%#v", plan.Components, plan.Assets)
	}

	if err := app.Start(context.Background()); err != nil {
		t.Fatalf("start: %v", err)
	}
	t.Cleanup(func() {
		if err := app.Stop(context.Background()); err != nil {
			t.Fatalf("stop: %v", err)
		}
	})

	if _, err := linker.RequireCapability(app.App(), graphlinker.RendererKey()); err != nil {
		t.Fatalf("graph renderer capability: %v", err)
	}

	httpServer, err := linker.RequireCapability(app.App(), linker.NewCapabilityKey[*http.Server](http.ID))
	if err != nil {
		t.Fatalf("http capability: %v", err)
	}
	baseURL := "http://" + httpServer.Addr()

	viewerPayload := getJSON(t, baseURL+"/api/v1/app2/graph/orders", "")
	viewerData := responseData(t, viewerPayload)
	if viewerData["type"] != "viewer" {
		t.Fatalf("unexpected viewer payload: %#v", viewerData)
	}
	viewerGraph := responseMap(t, viewerData["graph"])
	if viewerGraph["name"] != "订单" {
		t.Fatalf("unexpected viewer graph: %#v", viewerGraph)
	}

	formPayload := getJSON(t, baseURL+"/api/v1/app2/graph/orders/form", "")
	formData := responseData(t, formPayload)
	if formData["type"] != "popForm" {
		t.Fatalf("unexpected form payload: %#v", formData)
	}

	behaviorPayload := getJSON(t, baseURL+"/api/v1/app2/graph/refresh", "")
	behaviorData := responseData(t, behaviorPayload)
	behavior := responseMap(t, behaviorData["behavior"])
	if behavior["type"] != "refresh" {
		t.Fatalf("unexpected behavior payload: %#v", behaviorData)
	}
}
