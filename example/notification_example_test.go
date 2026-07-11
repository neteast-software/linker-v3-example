package example

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	stdhttp "net/http"
	"strings"
	"testing"
	"time"

	auditcore "github.com/neteast-software/go-module/audit/operate"
	eventcore "github.com/neteast-software/go-module/fault/event"
	http "github.com/neteast-software/go-module/http/gin/linker"
	server "github.com/neteast-software/go-module/linker/server"
	consumer "github.com/neteast-software/go-module/mq/consumer"
	mq "github.com/neteast-software/go-module/mq/consumer/linker"
	"github.com/neteast-software/go-module/observe/metrics"
	"github.com/neteast-software/go-module/observe/tracing"
	traceconsumer "github.com/neteast-software/go-module/observe/tracing/mq/consumer"
	cron "github.com/neteast-software/go-module/scheduler/cron"
	schedule "github.com/neteast-software/go-module/scheduler/cron/linker"
	linker "github.com/neteast-software/linker/v3"

	notificationcomponent "linker-v3-example/internal/component/notification"
	observabilitycomponent "linker-v3-example/internal/component/observability"
	notificationservice "linker-v3-example/internal/service/notification"
)

func TestNotificationLifecycleExample(t *testing.T) {
	eventRecorder := eventcore.NewMemoryRecorder()
	auditRecorder := auditcore.NewMemoryRecorder()
	metricRecorder := metrics.Memory()
	notification := notificationcomponent.NewComponent(
		notificationcomponent.WithCronSpec("@every 1s"),
		notificationcomponent.WithMetricRecorder(metricRecorder),
		notificationcomponent.WithMetricLabels(metrics.Label("service", "example")),
	)
	httpConfig := http.DefaultConfig()
	httpConfig.Addr = "127.0.0.1:0"
	app := server.New(
		server.WithShutdownTimeout(2*time.Second),
		server.WithHTTP(httpConfig),
		server.WithEventRecorder(eventRecorder),
		server.WithAuditRecorder(auditRecorder),
		server.WithoutNotice(),
		server.WithComponents(
			observabilitycomponent.NewComponent(),
			notification,
			mq.New(),
			schedule.New(
				schedule.WithStore(cron.NewMemoryStore()),
			),
		),
	)

	plan := preparedPlan(t, app)
	if !planHasAsset(plan, "mq/consumer", "notification") {
		t.Fatalf("missing consumer asset plan: %#v", plan.Assets)
	}
	if !planHasAsset(plan, "scheduler/cron/job", "notification.health") {
		t.Fatalf("missing cron asset plan: %#v", plan.Assets)
	}
	if !planHasRouteAsset(plan, "GET", "/api/v1/app2/notification/events", "http.app2.notification.events") {
		t.Fatalf("missing SSE route asset plan: %#v", plan.Assets)
	}
	if !planHasRouteAsset(plan, "POST", "/api/v1/app2/notification/send", "http.app2.notification.send") {
		t.Fatalf("missing notification send route asset plan: %#v", plan.Assets)
	}

	if err := app.Start(context.Background()); err != nil {
		t.Fatalf("start: %v", err)
	}
	t.Cleanup(func() {
		if err := app.Stop(context.Background()); err != nil {
			t.Fatalf("stop: %v", err)
		}
	})

	executor, err := mq.Require(app, "notification")
	if err != nil {
		t.Fatalf("consumer capability: %v", err)
	}
	mqCtx, _ := tracing.Ensure(context.Background(), tracing.WithTraceID("trace-notification-message"))
	message := traceconsumer.InjectMessage(mqCtx, consumer.Message{
		Topic: "notification.message",
		Key:   "n1",
		Body:  []byte("hello notification"),
	})
	if err = executor.Submit(context.Background(), message); err != nil {
		t.Fatalf("submit: %v", err)
	}
	waitFor(t, 2*time.Second, func() bool {
		return hasProviderMessage(notification.Provider().Messages(), "mq") &&
			hasProviderTrace(notification.Provider().Messages(), "mq", "trace-notification-message") &&
			len(auditRecorder.Records()) > 0 &&
			len(eventRecorder.Events()) > 0
	})
	assertMetricLabel(t, metricRecorder, "mq_consumer_messages_total", "status", "handled")
	assertMetricLabel(t, metricRecorder, "mq_consumer_messages_total", "service", "example")

	httpServer, err := http.RequireServer(app)
	if err != nil {
		t.Fatalf("http capability: %v", err)
	}
	sendBody, err := json.Marshal(map[string]string{"key": "http", "body": "hello from http"})
	if err != nil {
		t.Fatalf("marshal send body: %v", err)
	}
	sendReq, err := stdhttp.NewRequest(
		stdhttp.MethodPost,
		"http://"+httpServer.Addr()+"/api/v1/app2/notification/send",
		bytes.NewReader(sendBody),
	)
	if err != nil {
		t.Fatalf("new send request: %v", err)
	}
	sendReq.Header.Set("Content-Type", "application/json")
	sendReq.Header.Set(tracing.HeaderTraceID, "trace-http-mq")
	sendResp, err := stdhttp.DefaultClient.Do(sendReq)
	if err != nil {
		t.Fatalf("post notification send: %v", err)
	}
	sendPayload, err := io.ReadAll(sendResp.Body)
	_ = sendResp.Body.Close()
	if err != nil {
		t.Fatalf("read send response: %v", err)
	}
	if sendResp.StatusCode != stdhttp.StatusOK ||
		sendResp.Header.Get(tracing.HeaderTraceID) != "trace-http-mq" ||
		!strings.Contains(string(sendPayload), "queued") {
		t.Fatalf("unexpected send response: status=%d header=%#v body=%q", sendResp.StatusCode, sendResp.Header, sendPayload)
	}
	waitFor(t, 2*time.Second, func() bool {
		return hasProviderTrace(notification.Provider().Messages(), "mq", "trace-http-mq")
	})
	waitFor(t, time.Second, func() bool {
		return hasAuditRecord(
			auditRecorder.Records(),
			stdhttp.MethodPost,
			"/api/v1/app2/notification/send",
			"http.app2.notification.send",
		)
	})

	sseReq, err := stdhttp.NewRequest(stdhttp.MethodGet, "http://"+httpServer.Addr()+"/api/v1/app2/notification/events", nil)
	if err != nil {
		t.Fatalf("new sse request: %v", err)
	}
	sseReq.Header.Set(tracing.HeaderTraceID, "trace-sse")
	resp, err := stdhttp.DefaultClient.Do(sseReq)
	if err != nil {
		t.Fatalf("get sse: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read sse: %v", err)
	}
	if resp.StatusCode != stdhttp.StatusOK ||
		!strings.Contains(resp.Header.Get("Content-Type"), "text/event-stream") ||
		resp.Header.Get(tracing.HeaderTraceID) != "trace-sse" ||
		!strings.Contains(string(body), "ready") {
		t.Fatalf("unexpected sse response: status=%d content-type=%q body=%q", resp.StatusCode, resp.Header.Get("Content-Type"), body)
	}

	waitFor(t, 2500*time.Millisecond, func() bool {
		return hasProviderMessage(notification.Provider().Messages(), "cron") &&
			hasProviderNonEmptyTrace(notification.Provider().Messages(), "cron")
	})
	assertMetricLabel(t, metricRecorder, "scheduler_cron_runs_total", "status", "success")
}

func planHasAsset(plan linker.Plan, kind string, name string) bool {
	for _, asset := range plan.Assets {
		if string(asset.Kind) == kind && asset.Name == name {
			return true
		}
	}
	return false
}

func planHasRouteAsset(plan linker.Plan, method string, path string, resource string) bool {
	for _, asset := range plan.Assets {
		if asset.Kind == linker.AssetRoute &&
			asset.Detail["method"] == method &&
			asset.Detail["path"] == path &&
			asset.Detail["resource"] == resource {
			return true
		}
	}
	return false
}

func hasProviderMessage(messages []notificationservice.Message, source string) bool {
	for _, message := range messages {
		if message.Source == source {
			return true
		}
	}
	return false
}

func hasProviderTrace(messages []notificationservice.Message, source string, traceID string) bool {
	for _, message := range messages {
		if message.Source == source && message.TraceID == traceID {
			return true
		}
	}
	return false
}

func hasProviderNonEmptyTrace(messages []notificationservice.Message, source string) bool {
	for _, message := range messages {
		if message.Source == source && message.TraceID != "" {
			return true
		}
	}
	return false
}

func hasAuditRecord(records []auditcore.Record, method string, path string, resource string) bool {
	for _, record := range records {
		if record.Method == method && record.Path == path && record.Resource == resource {
			return true
		}
	}
	return false
}

func waitFor(t *testing.T, timeout time.Duration, ok func() bool) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if ok() {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("condition not met within %s", timeout)
}
