package example

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	stdhttp "net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	applicationcomponent "github.com/neteast-software/go-module/application/linker"
	audit "github.com/neteast-software/go-module/audit/operate/linker"
	postgresqlcore "github.com/neteast-software/go-module/db/postgresql"
	postgresql "github.com/neteast-software/go-module/db/postgresql/linker"
	http "github.com/neteast-software/go-module/http/gin/linker"
	mq "github.com/neteast-software/go-module/mq/consumer/linker"
	"github.com/neteast-software/go-module/observe/tracing"
	rpccore "github.com/neteast-software/go-module/rpc/grpc"
	rpc "github.com/neteast-software/go-module/rpc/grpc/linker"
	schedule "github.com/neteast-software/go-module/scheduler/cron/linker"
	grpcdiscovery "github.com/neteast-software/grpc-discovery"
	linker "github.com/neteast-software/linker/v3"

	exampleapp "linker-v3-example/internal/app"
	ttsclient "linker-v3-example/internal/client/tts"
	inspectioncomponent "linker-v3-example/internal/component/inspection"
	notificationcomponent "linker-v3-example/internal/component/notification"
	observabilitycomponent "linker-v3-example/internal/component/observability"
	ttscomponent "linker-v3-example/internal/component/tts"
	usercomponent "linker-v3-example/internal/component/user"
	userconstant "linker-v3-example/internal/constant/user"
)

func TestBusinessSystemExampleWithPostgreSQL(t *testing.T) {
	httpConfig := http.DefaultConfig()
	httpConfig.Addr = "127.0.0.1:0"
	grpcConfig := rpccore.ServerConfig{Addr: freeLocalAddr(t)}
	discovery := &exampleDiscovery{addr: grpcConfig.Addr}
	rpccore.ConfigureDiscovery("example", discovery)
	postgresqlConfig := prepareExampleDatabase(t)
	ttsConfig := rpccore.ClientConfig{
		Discovery: rpccore.ClientDiscoveryConfig{
			Scheme:  "example",
			Service: "tts",
			Group:   "DEFAULT_GROUP",
			Metadata: map[string]string{
				"version": "v1",
			},
		},
		Timeout: time.Second,
		Metadata: map[string]string{
			"scope": "front",
		},
	}
	app := exampleapp.New(businessConfigSource(t, map[linker.Namespace]any{
		http.Namespace:                 httpConfig,
		rpc.Namespace:                  grpcConfig,
		linker.Namespace(ttsclient.ID): ttsConfig,
		postgresql.Namespace:           postgresqlConfig,
		usercomponent.Namespace:        usercomponent.Config{TokenKey: strings.Repeat("a", 64)},
	}))
	plan := preparedPlan(t, app)
	if !planHasComponent(plan, postgresql.ID) ||
		!planHasComponent(plan, applicationcomponent.ID) ||
		!planHasComponent(plan, audit.ID) ||
		!planHasComponent(plan, inspectioncomponent.ID) ||
		!planHasComponent(plan, notificationcomponent.ID) ||
		!planHasComponent(plan, observabilitycomponent.ID) ||
		!planHasComponent(plan, usercomponent.ID) ||
		!planHasComponent(plan, ttscomponent.ID) ||
		!planHasComponent(plan, mq.ID) ||
		!planHasComponent(plan, schedule.ID) ||
		!planHasComponent(plan, rpc.ID) ||
		!planHasComponent(plan, ttsclient.ID) ||
		!planHasComponent(plan, http.ID) {
		t.Fatalf("plan missing component: %#v", plan.Components)
	}
	if !planHasAsset(plan, "application", "app2") {
		t.Fatalf("plan missing application asset: %#v", plan.Assets)
	}
	if !planHasAsset(plan, "mq/consumer", "notification") {
		t.Fatalf("plan missing notification consumer asset: %#v", plan.Assets)
	}
	if !planHasAsset(plan, "scheduler/cron/job", "notification.health") {
		t.Fatalf("plan missing notification cron asset: %#v", plan.Assets)
	}
	if !planHasAsset(plan, "rpc/grpc/server", grpcConfig.Addr) {
		t.Fatalf("plan missing grpc server asset: %#v", plan.Assets)
	}
	if !planHasAsset(plan, "rpc/grpc/client", ttsclient.ID.String()) {
		t.Fatalf("plan missing tts grpc client asset: %#v", plan.Assets)
	}
	if !planHasAsset(plan, "observe/metrics", "prometheus") || !planHasAsset(plan, "observe/tracing", "http+grpc") {
		t.Fatalf("plan missing observability asset: %#v", plan.Assets)
	}
	if !planHasRouteAsset(plan, "GET", "/api/v1/app2/notification/events", "http.app2.notification.events") {
		t.Fatalf("plan missing notification SSE route asset: %#v", plan.Assets)
	}
	if !planHasRouteAsset(plan, "POST", "/api/v1/app2/notification/send", "http.app2.notification.send") {
		t.Fatalf("plan missing notification send route asset: %#v", plan.Assets)
	}
	if !planHasRouteAsset(plan, "POST", "/api/v1/app2/tts/transcribe", "http.app2.tts.transcribe") {
		t.Fatalf("plan missing tts transcribe route asset: %#v", plan.Assets)
	}
	postgresqlOrder := planOrder(plan, postgresql.ID)
	applicationOrder := planOrder(plan, applicationcomponent.ID)
	auditOrder := planOrder(plan, audit.ID)
	inspectionOrder := planOrder(plan, inspectioncomponent.ID)
	notificationOrder := planOrder(plan, notificationcomponent.ID)
	userOrder := planOrder(plan, usercomponent.ID)
	ttsOrder := planOrder(plan, ttscomponent.ID)
	mqOrder := planOrder(plan, mq.ID)
	cronOrder := planOrder(plan, schedule.ID)
	rpcOrder := planOrder(plan, rpc.ID)
	ttsClientOrder := planOrder(plan, ttsclient.ID)
	httpOrder := planOrder(plan, http.ID)
	if postgresqlOrder >= auditOrder ||
		postgresqlOrder >= applicationOrder ||
		postgresqlOrder >= inspectionOrder ||
		postgresqlOrder >= userOrder ||
		postgresqlOrder >= ttsOrder ||
		applicationOrder >= httpOrder ||
		inspectionOrder >= httpOrder ||
		notificationOrder >= mqOrder ||
		notificationOrder >= cronOrder ||
		ttsOrder >= rpcOrder ||
		rpcOrder >= ttsClientOrder ||
		auditOrder >= httpOrder ||
		userOrder >= httpOrder ||
		ttsClientOrder >= httpOrder {
		t.Fatalf("unexpected startup order: postgresql=%d application=%d audit=%d inspection=%d notification=%d user=%d tts=%d mq=%d cron=%d rpc=%d ttsClient=%d http=%d plan=%#v", postgresqlOrder, applicationOrder, auditOrder, inspectionOrder, notificationOrder, userOrder, ttsOrder, mqOrder, cronOrder, rpcOrder, ttsClientOrder, httpOrder, plan.Components)
	}
	if err := app.Start(context.Background()); err != nil {
		t.Fatalf("start: %v", err)
	}
	t.Cleanup(func() {
		if err := app.Stop(context.Background()); err != nil {
			t.Fatalf("stop: %v", err)
		}
	})

	httpServer, err := linker.RequireCapability(app, linker.NewCapabilityKey[*http.Server](http.ID))
	if err != nil {
		t.Fatalf("http capability: %v", err)
	}
	baseURL := "http://" + httpServer.Addr()

	admin := postJSON(t, baseURL+"/system/login", map[string]string{
		"username": "admin",
		"password": userconstant.ExampleLoginPassword,
	}, "")
	adminData := responseData(t, admin)
	adminToken, _ := adminData["token"].(string)
	if adminToken == "" {
		t.Fatalf("missing admin token: %#v", admin)
	}
	adminUser := responseMap(t, adminData["user"])
	if adminUser["username"] != "admin" || adminUser["role"] != "admin" {
		t.Fatalf("unexpected admin user: %#v", adminUser)
	}
	adminProfile := getJSON(t, baseURL+"/system/profile", adminToken)
	adminProfileData := responseData(t, adminProfile)
	if adminProfileData["username"] != "admin" || adminProfileData["role"] != "admin" {
		t.Fatalf("unexpected admin profile: %#v", adminProfileData)
	}

	front := postJSON(t, baseURL+"/user/login", map[string]string{
		"phone":    userconstant.ExamplePhone,
		"password": userconstant.ExampleLoginPassword,
	}, "")
	frontData := responseData(t, front)
	token, _ := frontData["token"].(string)
	if token == "" {
		t.Fatalf("missing user token: %#v", front)
	}

	profile := getJSON(t, baseURL+"/api/profile", token)
	profileData := responseData(t, profile)
	if profileData["username"] != "example_user" ||
		profileData["phone"] != userconstant.ExamplePhone ||
		profileData["email"] != "example.user@neteast.cn" ||
		profileData["avatar"] == "" {
		t.Fatalf("unexpected profile: %#v", profileData)
	}

	userID, ok := profileData["id"].(float64)
	if !ok || userID == 0 {
		t.Fatalf("unexpected profile id: %#v", profileData)
	}
	app2Profile := getJSON(t, fmt.Sprintf("%s/api/v1/app2/user/%d/profile", baseURL, int(userID)), "")
	app2ProfileData := responseData(t, app2Profile)
	if app2ProfileData["username"] != "example_user" ||
		app2ProfileData["phone"] != userconstant.ExamplePhone ||
		app2ProfileData["application"] != "app2" {
		t.Fatalf("unexpected app2 profile: %#v", app2ProfileData)
	}

	tasks := getJSON(t, baseURL+"/api/v1/app2/inspection/tasks?page=1&pageSize=10&status=open", token)
	taskData := responseData(t, tasks)
	items, ok := taskData["items"].([]any)
	if !ok || len(items) != 1 {
		t.Fatalf("unexpected task items: %#v", taskData)
	}
	task, ok := items[0].(map[string]any)
	if !ok || task["application_scope"] != "app2" || task["status"] != "open" || task["owner_id"] != userID {
		t.Fatalf("unexpected scoped task: %#v", task)
	}

	tts, err := ttsclient.Require(app)
	if err != nil {
		t.Fatalf("tts client capability: %v", err)
	}
	ttsResult, err := tts.Transcribe(context.Background(), "hello")
	if err != nil {
		t.Fatalf("tts transcribe: %v", err)
	}
	if ttsResult != "tts:hello:front" {
		t.Fatalf("unexpected tts result: %q", ttsResult)
	}
	httpTTS, traceHeaders := postJSONHeaders(t, baseURL+"/api/v1/app2/tts/transcribe", map[string]string{
		"text": "hello-http",
	}, "", map[string]string{
		tracing.HeaderTraceID: "trace-http-grpc",
	})
	if traceHeaders.Get(tracing.HeaderTraceID) != "trace-http-grpc" ||
		traceHeaders.Get(tracing.HeaderRequestID) == "" {
		t.Fatalf("trace headers missing from http tts response: %#v", traceHeaders)
	}
	httpTTSData := responseData(t, httpTTS)
	if httpTTSData["result"] != "tts:hello-http:front" {
		t.Fatalf("unexpected http tts result: %#v", httpTTSData)
	}
	discovery.assertLookup(t, "tts", map[string]string{"version": "v1"})
	cancelCtx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err = tts.Transcribe(cancelCtx, "cancel"); err == nil {
		t.Fatalf("expected canceled transcribe error")
	}

	db, err := postgresql.Require(app)
	if err != nil {
		t.Fatalf("db capability: %v", err)
	}
	var ttsCount int64
	if err = db.Table("tts_conversion").
		Where("text = ? AND result = ? AND scope = ?", "hello", "tts:hello:front", "front").
		Count(&ttsCount).Error; err != nil {
		t.Fatalf("count tts request: %v", err)
	}
	if ttsCount != 1 {
		t.Fatalf("tts records = %d, want 1", ttsCount)
	}
	var auditCount int64
	if err = db.Table("operate").
		Where("resource IN ? AND successful = ?", []string{"http.console.auth.login", "http.front.auth.login"}, true).
		Count(&auditCount).Error; err != nil {
		t.Fatalf("count audit: %v", err)
	}
	if auditCount < 2 {
		t.Fatalf("audit records = %d, want at least 2", auditCount)
	}
}

func businessConfigSource(t *testing.T, values map[linker.Namespace]any) linker.Source {
	t.Helper()
	setting := make(map[linker.Namespace][]byte, len(values))
	for namespace, value := range values {
		content, err := json.Marshal(value)
		if err != nil {
			t.Fatalf("marshal config %s: %v", namespace, err)
		}
		setting[namespace] = content
	}
	return linker.MapSource{Label: "config/test", Setting: linker.NewSetting(setting)}
}

func prepareExampleDatabase(t *testing.T) postgresqlcore.Config {
	t.Helper()
	config := postgresqlcore.Config{
		Host:     getenv("LINKER_V3_EXAMPLE_PG_HOST", "127.0.0.1"),
		Port:     getenvInt("LINKER_V3_EXAMPLE_PG_PORT", 5432),
		User:     getenv("LINKER_V3_EXAMPLE_PG_USER", "neteast"),
		Password: os.Getenv("LINKER_V3_EXAMPLE_PG_PASSWORD"),
		DBName:   "postgres",
		SSLMode:  "disable",
	}
	if config.Password == "" {
		t.Skip("未设置 LINKER_V3_EXAMPLE_PG_PASSWORD，跳过真实数据库 example")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	adminDB, err := postgresqlcore.Open(ctx, config)
	if err != nil {
		t.Skipf("pi2 PostgreSQL 暂不可用，跳过真实数据库 example: %v", err)
	}
	defer func() {
		_ = postgresqlcore.Close(adminDB)
	}()

	dbName := fmt.Sprintf("linker_v3_example_%d", time.Now().UnixNano())
	if err = adminDB.WithContext(ctx).Exec("CREATE DATABASE " + quotePGIdent(dbName)).Error; err != nil {
		t.Skipf("无法创建 example 数据库，跳过真实数据库 example: %v", err)
	}
	adminConfig := config
	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cleanupCancel()
		cleanupDB, cleanupErr := postgresqlcore.Open(cleanupCtx, adminConfig)
		if cleanupErr != nil {
			t.Logf("清理 example 数据库失败: %v", cleanupErr)
			return
		}
		defer func() {
			_ = postgresqlcore.Close(cleanupDB)
		}()
		if dropErr := cleanupDB.WithContext(cleanupCtx).Exec("DROP DATABASE IF EXISTS " + quotePGIdent(dbName)).Error; dropErr != nil {
			t.Logf("删除 example 数据库失败: %v", dropErr)
		}
	})
	config.DBName = dbName
	return config
}

type exampleDiscovery struct {
	mu    sync.Mutex
	addr  string
	query grpcdiscovery.Query
}

func (p *exampleDiscovery) Lookup(_ context.Context, query grpcdiscovery.Query) ([]grpcdiscovery.Instance, error) {
	p.mu.Lock()
	p.query = query
	p.mu.Unlock()
	return []grpcdiscovery.Instance{
		{
			Addr: p.addr,
			Metadata: map[string]string{
				"version": "v1",
			},
		},
	}, nil
}

func (p *exampleDiscovery) Watch(context.Context, grpcdiscovery.Query, func([]grpcdiscovery.Instance, error)) (grpcdiscovery.Watcher, error) {
	return noopWatcher{}, nil
}

func (p *exampleDiscovery) assertLookup(t *testing.T, service string, metadata map[string]string) {
	t.Helper()
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.query.ServiceName != service {
		t.Fatalf("discovery service = %q", p.query.ServiceName)
	}
	for key, want := range metadata {
		if p.query.Metadata[key] != want {
			t.Fatalf("discovery metadata[%s] = %q, want %q query=%#v", key, p.query.Metadata[key], want, p.query)
		}
	}
}

type noopWatcher struct{}

func (noopWatcher) Close() error {
	return nil
}

func freeLocalAddr(t *testing.T) string {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen free addr: %v", err)
	}
	addr := listener.Addr().String()
	if err = listener.Close(); err != nil {
		t.Fatalf("close free addr listener: %v", err)
	}
	return addr
}

func postJSON(t *testing.T, url string, payload any, token string) map[string]any {
	t.Helper()
	result, _ := postJSONHeaders(t, url, payload, token, nil)
	return result
}

func postJSONHeaders(t *testing.T, url string, payload any, token string, headers map[string]string) (map[string]any, stdhttp.Header) {
	t.Helper()
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}
	req, err := stdhttp.NewRequest(stdhttp.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		t.Fatalf("new post request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	return doJSONHeaders(t, req)
}

func getJSON(t *testing.T, url string, token string) map[string]any {
	t.Helper()
	req, err := stdhttp.NewRequest(stdhttp.MethodGet, url, nil)
	if err != nil {
		t.Fatalf("new get request: %v", err)
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return doJSON(t, req)
}

func doJSON(t *testing.T, req *stdhttp.Request) map[string]any {
	t.Helper()
	payload, _ := doJSONHeaders(t, req)
	return payload
}

func doJSONHeaders(t *testing.T, req *stdhttp.Request) (map[string]any, stdhttp.Header) {
	t.Helper()
	resp, err := stdhttp.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("http request: %v", err)
	}
	defer resp.Body.Close()
	headers := resp.Header.Clone()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read response: %v", err)
	}
	if resp.StatusCode != stdhttp.StatusOK {
		t.Fatalf("unexpected status=%d body=%s", resp.StatusCode, body)
	}
	var payload map[string]any
	if err = json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("decode response: %v body=%s", err, body)
	}
	if code, _ := payload["code"].(float64); code != 0 {
		t.Fatalf("business response failed: %#v", payload)
	}
	return payload, headers
}

func responseData(t *testing.T, payload map[string]any) map[string]any {
	t.Helper()
	return responseMap(t, payload["data"])
}

func responseMap(t *testing.T, value any) map[string]any {
	t.Helper()
	data, ok := value.(map[string]any)
	if !ok {
		t.Fatalf("response data is not object: %#v", value)
	}
	return data
}

func planHasComponent(plan linker.Plan, id linker.ID) bool {
	return planOrder(plan, id) > 0
}

func planOrder(plan linker.Plan, id linker.ID) int {
	for _, component := range plan.Components {
		if component.ID == id {
			return component.Order
		}
	}
	return 0
}

func quotePGIdent(name string) string {
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}

func getenv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func getenvInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}
