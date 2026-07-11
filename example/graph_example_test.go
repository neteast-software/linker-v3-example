package example

import (
	"context"
	"io"
	stdhttp "net/http"
	"strings"
	"testing"
	"time"

	applicationcore "github.com/neteast-software/go-module/application"
	applicationcomponent "github.com/neteast-software/go-module/application/linker"
	graphlinker "github.com/neteast-software/go-module/graph/naive/linker"
	http "github.com/neteast-software/go-module/http/gin/linker"
	server "github.com/neteast-software/go-module/linker/server"
	"github.com/neteast-software/go-module/observe/tracing"

	graphcomponent "linker-v3-example/internal/component/graph"
	observabilitycomponent "linker-v3-example/internal/component/observability"
)

func TestGraphNaiveExample(t *testing.T) {
	httpConfig := http.DefaultConfig()
	httpConfig.Addr = "127.0.0.1:0"
	app := server.New(
		server.WithShutdownTimeout(3*time.Second),
		server.WithHTTP(httpConfig),
		server.WithComponents(
			observabilitycomponent.NewComponent(),
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

	plan := preparedPlan(t, app)
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

	if _, err := graphlinker.Require(app); err != nil {
		t.Fatalf("graph renderer capability: %v", err)
	}

	httpServer, err := http.RequireServer(app)
	if err != nil {
		t.Fatalf("http capability: %v", err)
	}
	baseURL := "http://" + httpServer.Addr()

	viewerReq, err := stdhttp.NewRequest(stdhttp.MethodGet, baseURL+"/api/v1/app2/graph/orders", nil)
	if err != nil {
		t.Fatalf("new graph request: %v", err)
	}
	viewerReq.Header.Set(tracing.HeaderTraceID, "trace-graph")
	viewerPayload, viewerHeaders := doJSONHeaders(t, viewerReq)
	if viewerHeaders.Get(tracing.HeaderTraceID) != "trace-graph" || viewerHeaders.Get(tracing.HeaderRequestID) == "" {
		t.Fatalf("trace headers missing from graph response: %#v", viewerHeaders)
	}
	viewerData := responseData(t, viewerPayload)
	if viewerData["type"] != "viewer" {
		t.Fatalf("unexpected viewer payload: %#v", viewerData)
	}
	viewerGraph := responseMap(t, viewerData["graph"])
	if viewerGraph["name"] != "订单" || viewerGraph["resource"] != "http.app2.graph.orders" {
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

	metricsPayload := getRaw(t, baseURL+"/metrics")
	if !strings.Contains(metricsPayload, `route="/api/v1/app2/graph/orders"`) ||
		!strings.Contains(metricsPayload, `status="200"`) {
		t.Fatalf("graph metrics missing:\n%s", metricsPayload)
	}
}

func getRaw(t *testing.T, url string) string {
	t.Helper()
	resp, err := stdhttp.Get(url)
	if err != nil {
		t.Fatalf("get raw: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read raw: %v", err)
	}
	return string(body)
}
