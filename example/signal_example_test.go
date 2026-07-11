//go:build !windows

package example

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	http "github.com/neteast-software/go-module/http/gin"
	server "github.com/neteast-software/go-module/linker/server"
	linker "github.com/neteast-software/linker/v3"
)

const (
	signalStartupID linker.ID = "example/signal-startup"
	signalRuntimeID linker.ID = "example/signal-runtime"
)

type signalStartupComponent struct {
	started string
	closed  string
}

func (p *signalStartupComponent) Identity() linker.ID {
	return signalStartupID
}

func (p *signalStartupComponent) Start(ctx context.Context) error {
	if err := os.WriteFile(p.started, []byte("starting"), 0o600); err != nil {
		return err
	}
	<-ctx.Done()
	return ctx.Err()
}

func (p *signalStartupComponent) Stop(context.Context) error {
	return nil
}

func (p *signalStartupComponent) Close(context.Context) error {
	return os.WriteFile(p.closed, []byte("closed"), 0o600)
}

type signalRuntimeComponent struct {
	closed string
}

func (p *signalRuntimeComponent) Identity() linker.ID {
	return signalRuntimeID
}

func (p *signalRuntimeComponent) Stop(context.Context) error {
	return nil
}

func (p *signalRuntimeComponent) Close(context.Context) error {
	return os.WriteFile(p.closed, []byte("closed"), 0o600)
}

func TestServerReceivesSIGTERMDuringStartup(t *testing.T) {
	runSignalProcess(t, "startup")
}

func TestServerReceivesSIGTERMDuringRuntime(t *testing.T) {
	runSignalProcess(t, "runtime")
}

func runSignalProcess(t *testing.T, phase string) {
	t.Helper()
	dir := t.TempDir()
	ready := filepath.Join(dir, "ready")
	closed := filepath.Join(dir, "closed")
	var output bytes.Buffer
	cmd := exec.Command(os.Args[0], "-test.run=^TestServerSignalHelper$")
	cmd.Env = append(os.Environ(),
		"LINKER_EXAMPLE_SIGNAL_HELPER=1",
		"LINKER_EXAMPLE_SIGNAL_PHASE="+phase,
		"LINKER_EXAMPLE_SIGNAL_READY="+ready,
		"LINKER_EXAMPLE_SIGNAL_CLOSED="+closed,
	)
	cmd.Stdout = &output
	cmd.Stderr = &output
	if err := cmd.Start(); err != nil {
		t.Fatalf("启动 %s signal helper: %v", phase, err)
	}
	if err := waitSignalFile(ready, 5*time.Second); err != nil {
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
		t.Fatalf("等待 %s server: %v\n%s", phase, err, output.String())
	}
	if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
		t.Fatalf("发送 %s SIGTERM: %v", phase, err)
	}
	if err := cmd.Wait(); err != nil {
		t.Fatalf("%s server 退出: %v\n%s", phase, err, output.String())
	}
	if err := waitSignalFile(closed, time.Second); err != nil {
		t.Fatalf("等待 %s graceful close: %v\n%s", phase, err, output.String())
	}
}

func TestServerSignalHelper(t *testing.T) {
	if os.Getenv("LINKER_EXAMPLE_SIGNAL_HELPER") != "1" {
		return
	}
	phase := os.Getenv("LINKER_EXAMPLE_SIGNAL_PHASE")
	ready := os.Getenv("LINKER_EXAMPLE_SIGNAL_READY")
	closed := os.Getenv("LINKER_EXAMPLE_SIGNAL_CLOSED")
	httpConfig := http.DefaultConfig()
	httpConfig.Addr = "127.0.0.1:0"
	httpConfig.BasePath = "/api"
	var component linker.Component
	options := []server.Option{server.WithHTTP(httpConfig)}
	switch phase {
	case "startup":
		component = &signalStartupComponent{started: ready, closed: closed}
	case "runtime":
		component = &signalRuntimeComponent{closed: closed}
		options = append(options, server.WithLifecycleObserver(func(_ context.Context, event linker.LifecycleEvent) {
			if event.Stage == linker.LifecycleStarted && event.Err == nil {
				_ = os.WriteFile(ready, []byte("running"), 0o600)
			}
		}))
	default:
		t.Fatalf("未知 signal phase: %s", phase)
	}
	options = append(options, server.WithComponents(component))
	err := server.New(options...).Run(context.Background())
	if phase == "startup" {
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("startup signal result: %v", err)
		}
		return
	}
	if err != nil {
		t.Fatalf("runtime signal result: %v", err)
	}
}

func waitSignalFile(path string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for {
		if _, err := os.Stat(path); err == nil {
			return nil
		} else if !errors.Is(err, os.ErrNotExist) {
			return err
		}
		if time.Now().After(deadline) {
			return context.DeadlineExceeded
		}
		time.Sleep(10 * time.Millisecond)
	}
}
