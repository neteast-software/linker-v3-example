package example

import (
	"context"
	"encoding/json"
	"io"
	stdhttp "net/http"
	"strings"
	"testing"
	"time"

	eventcore "github.com/neteast-software/go-module/fault/event"
	"github.com/neteast-software/go-module/http/gin/gateway"
	http "github.com/neteast-software/go-module/http/gin/linker"
	"github.com/neteast-software/go-module/http/gin/param"
	"github.com/neteast-software/go-module/http/gin/response"
	"github.com/neteast-software/go-module/license"
	licensehttp "github.com/neteast-software/go-module/license/http/gin"
	server "github.com/neteast-software/go-module/linker/server"
	"github.com/neteast-software/go-module/observe/metrics"
	metricserver "github.com/neteast-software/go-module/observe/metrics/linker/server"
	"github.com/neteast-software/go-module/observe/tracing"
	linker "github.com/neteast-software/linker/v3"

	observabilitycomponent "linker-v3-example/internal/component/observability"
	observabilityservice "linker-v3-example/internal/service/observability"
)

func TestLinkerV3HTTPGinExample(t *testing.T) {
	recorder := eventcore.NewMemoryRecorder()
	httpConfig := http.DefaultConfig()
	httpConfig.Addr = "127.0.0.1:0"
	now := time.Date(2026, 7, 8, 0, 0, 0, 0, time.Local)
	gate := license.NewGate(license.New(
		license.WithExpireAt(now.Add(-time.Second)),
		license.WithClock(func() time.Time { return now }),
	))
	app := server.New(
		server.WithShutdownTimeout(3*time.Second),
		server.WithEventRecorder(recorder),
		server.WithHTTP(httpConfig),
		server.WithHTTPRoutes(
			http.GET("pong", func(c *http.Context) {
				c.String(stdhttp.StatusOK, "pong")
			}),
			http.GET("items/:id", func(c *http.Context) {
				id, err := param.ID(c)
				if err != nil {
					response.Warning(c, "%s", err.Error())
					return
				}
				url, err := gateway.CurrentURL(c, "items/"+c.Param("id"))
				if err != nil {
					response.Warning(c, "%s", err.Error())
					return
				}
				response.Data(c, map[string]any{"id": id, "url": url})
			}),
			http.GET("licensed", func(c *http.Context) {
				response.Success(c)
			}).With(licensehttp.Gate(gate)),
		),
	)

	if plan := app.Plan(); len(plan.Components) < 3 {
		t.Fatalf("plan components = %#v", plan.Components)
	}

	if err := app.Start(context.Background()); err != nil {
		t.Fatalf("start: %v", err)
	}
	t.Cleanup(func() {
		if err := app.Stop(context.Background()); err != nil {
			t.Fatalf("stop: %v", err)
		}
	})

	runtimePlan := app.Plan()
	if !runtimePlan.Started || !runtimePlan.Ready {
		t.Fatalf("runtime plan = %#v", runtimePlan)
	}
	httpRuntimePlan := requireCoreComponentPlan(t, runtimePlan, http.ID)
	if !httpRuntimePlan.Active || !httpRuntimePlan.Initialized || !httpRuntimePlan.Started ||
		!httpRuntimePlan.Enabled || !httpRuntimePlan.Effective || httpRuntimePlan.Degraded {
		t.Fatalf("http runtime plan = %#v", httpRuntimePlan)
	}
	if !corePlanHasCapability(runtimePlan, http.ID, http.ID) {
		t.Fatalf("runtime plan missing http capability: %#v", runtimePlan.Capabilities)
	}

	httpServer, err := linker.RequireCapability(app, linker.NewCapabilityKey[*http.Server](http.ID))
	if err != nil {
		t.Fatalf("server capability: %v", err)
	}

	resp, err := stdhttp.Get("http://" + httpServer.Addr() + "/pong")
	if err != nil {
		t.Fatalf("get ping: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	if resp.StatusCode != stdhttp.StatusOK || string(body) != "pong" {
		t.Fatalf("unexpected response: status=%d body=%q", resp.StatusCode, body)
	}

	resp, err = stdhttp.Get("http://" + httpServer.Addr() + "/licensed")
	if err != nil {
		t.Fatalf("get licensed: %v", err)
	}
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read licensed body: %v", err)
	}
	var licensed map[string]any
	if err := json.Unmarshal(body, &licensed); err != nil {
		t.Fatalf("decode licensed body: %v body=%q", err, body)
	}
	if licensed["code"].(float64) != response.CodeFailure || licensed["msg"] != "授权已到期" {
		t.Fatalf("unexpected licensed payload: %#v", licensed)
	}

	req, err := stdhttp.NewRequest(stdhttp.MethodGet, "http://"+httpServer.Addr()+"/items/7", nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("X-Forwarded-Host", "gateway.example")
	req.Header.Set("X-Forwarded-Proto", "https")
	req.Header.Set("X-Forwarded-Prefix", "/api")
	resp, err = stdhttp.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("get item: %v", err)
	}
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read item body: %v", err)
	}
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("decode item body: %v body=%q", err, body)
	}
	data, ok := payload["data"].(map[string]any)
	if !ok || data["url"] != "https://gateway.example/api/items/7" || data["id"].(float64) != 7 {
		t.Fatalf("unexpected item payload: %#v", payload)
	}
}

func requireCoreComponentPlan(t *testing.T, plan linker.Plan, id linker.ID) linker.ComponentPlan {
	t.Helper()
	for _, component := range plan.Components {
		if component.ID == id {
			return component
		}
	}
	t.Fatalf("core component %s not found in %#v", id, plan.Components)
	return linker.ComponentPlan{}
}

func corePlanHasCapability(plan linker.Plan, id linker.ID, owner linker.ID) bool {
	for _, capability := range plan.Capabilities {
		if capability.ID == id && capability.Owner == owner {
			return true
		}
	}
	return false
}

func TestLinkerV3PrometheusMetricsExample(t *testing.T) {
	observability := observabilitycomponent.NewComponent()
	httpConfig := http.DefaultConfig()
	httpConfig.Addr = "127.0.0.1:0"
	app := server.New(
		server.WithShutdownTimeout(3*time.Second),
		server.WithHTTP(httpConfig),
		server.WithObserver(metricserver.Observe(
			observability.Recorder(),
			metricserver.WithConstLabels(metrics.Label("service", "linker-v3-example")),
		)),
		server.WithComponents(
			observability,
		),
		server.WithHTTPRoutes(
			http.GET("pong", func(c *http.Context) {
				c.String(stdhttp.StatusOK, "pong")
			}),
		),
	)

	plan := preparedPlan(t, app)
	if !planHasComponent(plan, observabilitycomponent.ID) ||
		!planHasRouteAsset(plan, "GET", "/metrics", "http.observe.metrics") ||
		!planHasAsset(plan, "observe/metrics", "prometheus") ||
		!planHasAsset(plan, "observe/tracing", "http+grpc") {
		t.Fatalf("plan missing observability assets: components=%#v assets=%#v", plan.Components, plan.Assets)
	}

	if err := app.Start(context.Background()); err != nil {
		t.Fatalf("start: %v", err)
	}
	t.Cleanup(func() {
		if err := app.Stop(context.Background()); err != nil {
			t.Fatalf("stop: %v", err)
		}
	})
	runtimePlan := app.Plan()
	observabilityPlan := requireCoreComponentPlan(t, runtimePlan, observabilitycomponent.ID)
	if !observabilityPlan.Enabled || !observabilityPlan.Effective || observabilityPlan.Degraded {
		t.Fatalf("observability runtime plan = %#v", observabilityPlan)
	}
	if !corePlanHasCapability(runtimePlan, observabilityservice.MetricRecorderID, observabilitycomponent.ID) {
		t.Fatalf("runtime plan missing capability owner: %#v", runtimePlan.Capabilities)
	}

	httpServer, err := linker.RequireCapability(app, linker.NewCapabilityKey[*http.Server](http.ID))
	if err != nil {
		t.Fatalf("server capability: %v", err)
	}
	req, err := stdhttp.NewRequest(stdhttp.MethodGet, "http://"+httpServer.Addr()+"/pong", nil)
	if err != nil {
		t.Fatalf("new pong request: %v", err)
	}
	req.Header.Set(tracing.HeaderTraceID, "trace-example")
	resp, err := stdhttp.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("get pong: %v", err)
	}
	if resp.Header.Get(tracing.HeaderTraceID) != "trace-example" || resp.Header.Get(tracing.HeaderRequestID) == "" {
		t.Fatalf("trace headers missing: %#v", resp.Header)
	}
	_ = resp.Body.Close()

	resp, err = stdhttp.Get("http://" + httpServer.Addr() + "/metrics")
	if err != nil {
		t.Fatalf("get metrics: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read metrics: %v", err)
	}
	text := string(body)
	if !strings.Contains(text, "linker_v3_example_http_requests_total") ||
		!strings.Contains(text, "linker_v3_example_linker_component_state") ||
		!strings.Contains(text, `component="example/observability"`) ||
		!strings.Contains(text, `route="/pong"`) ||
		!strings.Contains(text, `status="200"`) {
		t.Fatalf("unexpected metrics body:\n%s", text)
	}
}
