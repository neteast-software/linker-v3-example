package example

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	stdhttp "net/http"
	"strconv"
	"testing"
	"testing/fstest"
	"time"

	applicationcore "github.com/neteast-software/go-module/application"
	applicationcomponent "github.com/neteast-software/go-module/application/linker"
	graphconsole "github.com/neteast-software/go-module/graph/console"
	http "github.com/neteast-software/go-module/http/gin/linker"
	server "github.com/neteast-software/go-module/linker/server"
	"github.com/neteast-software/go-module/token"
	linker "github.com/neteast-software/linker/v3"

	console "linker-v3-example/internal/console/linker"
	order "linker-v3-example/internal/order/linker"
	permission "linker-v3-example/internal/permission/linker"
	userdata "linker-v3-example/internal/user"
	user "linker-v3-example/internal/user/linker"
)

type graphUser struct {
	user userdata.User
}

func (p *graphUser) AdminLogin(_ context.Context, username string, password string) (userdata.User, string, error) {
	if username != "admin" || password != "demo-password" {
		return userdata.User{}, "", context.Canceled
	}
	return p.user, "graph-token", nil
}

func (p *graphUser) Current(_ context.Context, raw string, scope string) (userdata.User, token.Claims, error) {
	if raw != "graph-token" || scope != "console" {
		return userdata.User{}, token.Claims{}, context.Canceled
	}
	return p.user, graphClaims(p.user.ID), nil
}

func (p *graphUser) Refresh(_ context.Context, raw string, scope string) (userdata.User, token.Token, error) {
	if raw != "graph-token" || scope != "console" {
		return userdata.User{}, token.Token{}, context.Canceled
	}
	return p.user, token.Token{Raw: raw, Claims: graphClaims(p.user.ID)}, nil
}

func (p *graphUser) Revoke(_ context.Context, raw string, scope string) error {
	if raw != "graph-token" || scope != "console" {
		return context.Canceled
	}
	return nil
}

func (p *graphUser) ProfileByID(_ context.Context, id uint64) (userdata.User, error) {
	if id != p.user.ID {
		return userdata.User{}, context.Canceled
	}
	return p.user, nil
}

type graphUserComponent struct {
	user *graphUser
}

func (p graphUserComponent) Identity() linker.ID {
	return user.ID
}

func (p graphUserComponent) OnMounted(_ context.Context, runtime linker.Runtime) error {
	return linker.Provide(runtime, userdata.AuthKey(), userdata.Auth(p.user))
}

func TestGraphConsoleExample(t *testing.T) {
	httpConfig := http.DefaultConfig()
	httpConfig.Addr = "127.0.0.1:0"
	current := &graphUser{user: userdata.User{
		Username: "admin",
		Avatar:   "https://static.neteast.cn/avatar/admin.png",
		Email:    "admin@neteast.cn",
		Phone:    "18558755877",
		Role:     "admin",
	}}
	current.user.ID = 1
	app := server.New(
		server.WithShutdownTimeout(3*time.Second),
		server.WithHTTP(httpConfig),
		server.WithComponents(
			applicationcomponent.New(
				applicationcomponent.WithDBTarget(""),
				applicationcomponent.WithApplications(applicationcore.Application{
					ID: "app2", Scope: "app2", Name: "应用二", Status: applicationcore.StatusEnabled,
				}),
			),
			graphUserComponent{user: current},
			order.New(),
			permission.New(),
			console.New(current, graphconsole.WithStatic(fstest.MapFS{
				"index.html":    {Data: []byte("<!doctype html><title>Graph Console</title>")},
				"assets/app.js": {Data: []byte("console.log('graph-console')")},
			})),
		),
	)

	plan := preparedPlan(t, app)
	for _, id := range []linker.ID{
		graphconsole.ID,
		order.ID,
		permission.ID,
		user.ID,
	} {
		if !planHasComponent(plan, id) {
			t.Fatalf("plan missing component %s: %#v", id, plan.Components)
		}
	}
	for _, route := range []struct {
		method   string
		path     string
		resource string
	}{
		{"GET", "/console/entry", ""},
		{"POST", "/console/login", ""},
		{"GET", "/console/menu", "graph.console.menu"},
		{"GET", "/console/page/:id", "graph.console.page"},
		{"PUT", "/api/v1/app2/order", "http.app2.order.update"},
		{"GET", "/api/v1/app2/permission/role/:id/resource", "http.app2.permission.role-resource.read"},
	} {
		if !planHasRouteAsset(plan, route.method, route.path, route.resource) {
			t.Fatalf("plan missing route %#v: %#v", route, plan.Assets)
		}
	}

	if err := app.Start(context.Background()); err != nil {
		t.Fatalf("start: %v", err)
	}
	t.Cleanup(func() {
		if err := app.Stop(context.Background()); err != nil {
			t.Fatalf("stop: %v", err)
		}
	})
	if _, ok := graphconsole.Resolve(app); !ok {
		t.Fatal("Graph Console service capability 未发布")
	}

	httpServer, err := http.RequireServer(app)
	if err != nil {
		t.Fatalf("http capability: %v", err)
	}
	baseURL := "http://" + httpServer.Addr()
	if body := getRaw(t, baseURL+"/"); body != "<!doctype html><title>Graph Console</title>" {
		t.Fatalf("unexpected static index: %s", body)
	}
	if body := getRaw(t, baseURL+"/assets/app.js"); body != "console.log('graph-console')" {
		t.Fatalf("unexpected static asset: %s", body)
	}

	entry := getJSON(t, baseURL+"/console/entry", "")
	if responseData(t, entry)["type"] != "config" {
		t.Fatalf("unexpected entry: %#v", entry)
	}
	loginPayload := postJSON(t, baseURL+"/console/login", map[string]any{
		"method": "password",
		"values": map[string]any{"username": "admin", "password": "demo-password"},
	}, "")
	loginData := responseData(t, loginPayload)
	access, _ := responseMap(t, loginData["token"])["access"].(string)
	if access == "" {
		t.Fatalf("missing console token: %#v", loginData)
	}

	menu := getJSON(t, baseURL+"/console/menu", access)
	if responseData(t, menu)["type"] != "menu" {
		t.Fatalf("unexpected menu: %#v", menu)
	}
	for _, pageCase := range []struct {
		identity string
		kind     string
	}{
		{"dashboard", "layout"},
		{"order.list", "viewer"},
		{"order.form", "form"},
		{"permission.role-resource", "multilist"},
	} {
		page := getJSON(t, baseURL+"/console/page/"+pageCase.identity, access)
		if responseData(t, page)["type"] != pageCase.kind {
			t.Fatalf("page %s = %#v", pageCase.identity, page)
		}
	}

	status, unauthorized := putJSONStatus(t, baseURL+"/api/v1/app2/order", map[string]any{
		"id": 1, "number": "NO-20260716-009", "status": "closed", "amount": 16800,
	}, "")
	if status != stdhttp.StatusUnauthorized || unauthorized["code"] != float64(stdhttp.StatusUnauthorized) {
		t.Fatalf("业务 API 未执行最终认证: status=%d payload=%#v", status, unauthorized)
	}
	invalidStatus, invalid := putJSONStatus(t, baseURL+"/api/v1/app2/order", map[string]any{
		"id": 1, "number": "", "status": "open", "amount": 100,
	}, access)
	if invalidStatus != stdhttp.StatusOK || invalid["code"] == float64(0) {
		t.Fatalf("服务端验证未拒绝无效订单: status=%d payload=%#v", invalidStatus, invalid)
	}
	valid := putJSON(t, baseURL+"/api/v1/app2/order", map[string]any{
		"id": 1, "number": "NO-20260716-009", "status": "closed", "amount": 16800,
	}, access)
	if responseData(t, valid)["number"] != "NO-20260716-009" {
		t.Fatalf("unexpected saved order: %#v", valid)
	}
}

func graphClaims(id uint64) token.Claims {
	now := time.Now()
	return token.Claims{
		Subject:   strconv.FormatUint(id, 10),
		Scope:     "console",
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(time.Hour).Unix(),
		Nonce:     "graph-example",
	}
}

func putJSON(t *testing.T, url string, value any, access string) map[string]any {
	t.Helper()
	status, payload := putJSONStatus(t, url, value, access)
	if status != stdhttp.StatusOK {
		t.Fatalf("unexpected status=%d payload=%#v", status, payload)
	}
	if code, _ := payload["code"].(float64); code != 0 {
		t.Fatalf("business response failed: %#v", payload)
	}
	return payload
}

func putJSONStatus(t *testing.T, url string, value any, access string) (int, map[string]any) {
	t.Helper()
	data, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}
	request, err := stdhttp.NewRequest(stdhttp.MethodPut, url, bytes.NewReader(data))
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	request.Header.Set("Content-Type", "application/json")
	if access != "" {
		request.Header.Set("Authorization", "Bearer "+access)
	}
	result, err := stdhttp.DefaultClient.Do(request)
	if err != nil {
		t.Fatalf("http request: %v", err)
	}
	defer result.Body.Close()
	body, err := io.ReadAll(result.Body)
	if err != nil {
		t.Fatalf("read response: %v", err)
	}
	var payload map[string]any
	if err = json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("decode response: %v body=%s", err, body)
	}
	return result.StatusCode, payload
}
