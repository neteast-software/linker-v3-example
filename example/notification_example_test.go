package example

import (
	"context"
	"io"
	stdhttp "net/http"
	"strings"
	"testing"
	"time"

	auditcore "github.com/neteast-software/go-module/audit/operate"
	eventcore "github.com/neteast-software/go-module/fault/event"
	http "github.com/neteast-software/go-module/http/gin/linker"
	server "github.com/neteast-software/go-module/linker/server"
	mq "github.com/neteast-software/go-module/mq/consumer/linker"
	cron "github.com/neteast-software/go-module/scheduler/cron/linker"
	linker "github.com/neteast-software/linker/v3"

	notificationcomponent "linker-v3-example/internal/component/notification"
	notificationservice "linker-v3-example/internal/service/notification"
)

func TestNotificationLifecycleExample(t *testing.T) {
	eventRecorder := eventcore.NewMemoryRecorder()
	auditRecorder := auditcore.NewMemoryRecorder()
	notification := notificationcomponent.NewComponent(notificationcomponent.WithCronSpec("@every 1s"))
	app := server.New(
		server.WithMode(linker.Server),
		server.WithShutdownTimeout(2*time.Second),
		server.WithHTTP(http.Config{Addr: "127.0.0.1:0"}),
		server.WithEventRecorder(eventRecorder),
		server.WithAuditRecorder(auditRecorder),
		server.WithoutNotice(),
		server.WithComponents(
			notification,
			mq.New(
				mq.WithAfter(notificationcomponent.ID),
				mq.WithConsumers(notification.Consumer()),
			),
			cron.New(
				cron.WithAfter(notificationcomponent.ID),
				cron.WithStore(cron.NewMemoryStore()),
				cron.WithJobs(notification.Job()),
			),
		),
	)

	plan := app.Plan()
	if !planHasAsset(plan, "mq/consumer", "notification") {
		t.Fatalf("missing consumer asset plan: %#v", plan.Assets)
	}
	if !planHasAsset(plan, "scheduler/cron/job", "notification.health") {
		t.Fatalf("missing cron asset plan: %#v", plan.Assets)
	}
	if !planHasRouteAsset(plan, "GET", "/api/v1/app2/notification/events", "http.app2.notification.events") {
		t.Fatalf("missing SSE route asset plan: %#v", plan.Assets)
	}

	if err := app.Start(context.Background()); err != nil {
		t.Fatalf("start: %v", err)
	}
	t.Cleanup(func() {
		if err := app.Stop(context.Background()); err != nil {
			t.Fatalf("stop: %v", err)
		}
	})

	consumer, err := linker.RequireCapability(app.App(), mq.ConsumerKey("notification"))
	if err != nil {
		t.Fatalf("consumer capability: %v", err)
	}
	if err = consumer.Submit(context.Background(), mq.Message{
		Topic: "notification.message",
		Key:   "n1",
		Body:  []byte("hello notification"),
	}); err != nil {
		t.Fatalf("submit: %v", err)
	}
	waitFor(t, 2*time.Second, func() bool {
		return hasProviderMessage(notification.Provider().Messages(), "mq") &&
			len(auditRecorder.Records()) > 0 &&
			len(eventRecorder.Events()) > 0
	})

	httpServer, err := linker.RequireCapability(app.App(), linker.NewCapabilityKey[*http.Server](http.ID))
	if err != nil {
		t.Fatalf("http capability: %v", err)
	}
	resp, err := stdhttp.Get("http://" + httpServer.Addr() + "/api/v1/app2/notification/events")
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
		!strings.Contains(string(body), "ready") {
		t.Fatalf("unexpected sse response: status=%d content-type=%q body=%q", resp.StatusCode, resp.Header.Get("Content-Type"), body)
	}

	waitFor(t, 2500*time.Millisecond, func() bool {
		return hasProviderMessage(notification.Provider().Messages(), "cron")
	})
}

func planHasAsset(plan server.Plan, kind string, name string) bool {
	for _, asset := range plan.Assets {
		if asset.Kind == kind && asset.Name == name {
			return true
		}
	}
	return false
}

func planHasRouteAsset(plan server.Plan, method string, path string, resource string) bool {
	for _, asset := range plan.Assets {
		if asset.Kind == "http/route" &&
			asset.Method == method &&
			asset.Path == path &&
			asset.Resource == resource {
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
