package example

import (
	"context"
	stdhttp "net/http"
	"sync"
	"testing"
	"time"

	http "github.com/neteast-software/go-module/http/gin/linker"
	server "github.com/neteast-software/go-module/linker/server"
	linker "github.com/neteast-software/linker/v3"
)

type startupExampleGate struct {
	entered chan struct{}
	release chan struct{}
}

func (p *startupExampleGate) Identity() linker.ID {
	return "example/startup-gate"
}

func (p *startupExampleGate) Start(ctx context.Context) error {
	close(p.entered)
	select {
	case <-p.release:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func TestStartupEndpointReflectsFrameworkState(t *testing.T) {
	config := http.DefaultConfig()
	config.Addr = "127.0.0.1:0"
	httpComponent := http.New(http.WithConfig(config))
	gate := &startupExampleGate{entered: make(chan struct{}), release: make(chan struct{})}
	var releaseOnce sync.Once
	release := func() { releaseOnce.Do(func() { close(gate.release) }) }
	app := server.New(
		server.WithStartupTimeout(3*time.Second),
		server.WithComponents(httpComponent, gate),
	)
	started := make(chan error, 1)
	go func() { started <- app.Start(context.Background()) }()
	t.Cleanup(func() {
		release()
		_ = app.Stop(context.Background())
	})
	select {
	case <-gate.entered:
	case <-time.After(time.Second):
		t.Fatal("startup gate did not start")
	}
	httpServer, err := http.RequireServer(app)
	if err != nil {
		t.Fatalf("http server during startup: %v", err)
	}
	assertExampleHealth(t, httpServer, "/live", stdhttp.StatusOK)
	assertExampleHealth(t, httpServer, "/ready", stdhttp.StatusServiceUnavailable)
	assertExampleHealth(t, httpServer, "/startup", stdhttp.StatusServiceUnavailable)

	release()
	if err := <-started; err != nil {
		t.Fatalf("start: %v", err)
	}
	assertExampleHealth(t, httpServer, "/ready", stdhttp.StatusOK)
	assertExampleHealth(t, httpServer, "/startup", stdhttp.StatusOK)
}

func assertExampleHealth(t *testing.T, server *http.Server, path string, want int) {
	t.Helper()
	response, err := stdhttp.Get("http://" + server.Addr() + path)
	if err != nil {
		t.Fatalf("GET %s: %v", path, err)
	}
	defer response.Body.Close()
	if response.StatusCode != want {
		t.Fatalf("GET %s status = %d, want %d", path, response.StatusCode, want)
	}
}
