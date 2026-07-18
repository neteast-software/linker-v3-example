package example

import (
	"context"
	"testing"
	"time"

	server "github.com/neteast-software/go-module/linker/server"
	periodic "github.com/neteast-software/go-module/worker/periodic"
	managed "github.com/neteast-software/go-module/worker/periodic/linker"
)

func TestManagedPeriodicWorkerExample(t *testing.T) {
	runs := make(chan struct{}, 1)
	worker := periodic.New(
		"profile-refresh",
		periodic.HandlerFunc(func(context.Context) error {
			runs <- struct{}{}
			return nil
		}),
		periodic.WithDescription("刷新用户资料缓存"),
		periodic.Every(time.Hour),
		periodic.WithTimeout(30*time.Second),
		periodic.WithBackoff(periodic.Exponential(time.Second, time.Minute)),
		periodic.UnhealthyAfter(3),
	)
	app := server.New(
		server.WithoutStartupLog(),
		server.WithoutHTTP(),
		server.WithoutEvent(),
		server.WithoutNotice(),
		server.WithoutAudit(),
		server.WithComponents(managed.New(managed.WithWorkers(worker))),
	)
	if err := app.Start(context.Background()); err != nil {
		t.Fatalf("启动周期 Worker：%v", err)
	}
	t.Cleanup(func() { _ = app.Stop(context.Background()) })

	select {
	case <-runs:
	case <-time.After(time.Second):
		t.Fatal("周期 Worker 未执行首次任务")
	}
	resolved, err := managed.Require(app, "profile-refresh")
	if err != nil || resolved != worker {
		t.Fatalf("解析周期 Worker capability：%#v，%v", resolved, err)
	}
	if !planHasAsset(app.Plan(), "worker/periodic", "profile-refresh") {
		t.Fatalf("Plan 缺少周期 Worker Asset：%#v", app.Plan().Assets)
	}
	if err := app.Stop(context.Background()); err != nil {
		t.Fatalf("停止周期 Worker：%v", err)
	}
}
