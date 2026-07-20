package example

import (
	"context"
	"io"
	stdhttp "net/http"
	"strings"
	"sync"
	"testing"
	"time"

	http "github.com/neteast-software/go-module/http/gin/linker"
	server "github.com/neteast-software/go-module/linker/server"
)

func TestServerHTTPProductionBoundaryExample(t *testing.T) {
	config := http.DefaultConfig()
	config.Addr = "127.0.0.1:0"
	config.MaxBodyBytes = 4
	config.Proxy.TrustedProxies = []string{"127.0.0.1"}
	entered := make(chan struct{})
	release := make(chan struct{})
	var releaseOnce sync.Once
	releaseRequest := func() { releaseOnce.Do(func() { close(release) }) }
	app := server.New(
		server.WithShutdownTimeout(2*time.Second),
		server.WithHTTP(config),
		server.WithHTTPRoot(func(next stdhttp.Handler) stdhttp.Handler {
			return stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
				w.Header().Set("X-HTTP-Root", "active")
				next.ServeHTTP(w, r)
			})
		}),
		server.WithHTTPRoutes(
			http.GET("client-ip", func(c *http.Context) {
				c.String(stdhttp.StatusOK, c.ClientIP())
			}),
			http.POST("upload", func(c *http.Context) {
				c.Status(stdhttp.StatusNoContent)
			}).ReadBodyWithin(time.Minute),
			http.GET("slow", func(c *http.Context) {
				close(entered)
				<-release
				c.String(stdhttp.StatusOK, "done")
			}),
		),
	)
	if err := app.Start(context.Background()); err != nil {
		t.Fatalf("start: %v", err)
	}
	t.Cleanup(func() {
		releaseRequest()
		_ = app.Stop(context.Background())
	})
	httpServer, err := http.RequireServer(app)
	if err != nil {
		t.Fatalf("http server: %v", err)
	}

	for _, path := range []string{"/live", "/ready", "/startup"} {
		response, err := stdhttp.Get("http://" + httpServer.Addr() + path)
		if err != nil {
			t.Fatalf("GET %s: %v", path, err)
		}
		_ = response.Body.Close()
		if response.StatusCode != stdhttp.StatusOK {
			t.Fatalf("GET %s status = %d", path, response.StatusCode)
		}
		if response.Header.Get("X-HTTP-Root") != "active" {
			t.Fatalf("GET %s 未经过 HTTP root", path)
		}
	}

	request, _ := stdhttp.NewRequest(stdhttp.MethodGet, "http://"+httpServer.Addr()+"/client-ip", nil)
	request.Header.Set("X-Forwarded-For", "203.0.113.8")
	response, err := stdhttp.DefaultClient.Do(request)
	if err != nil {
		t.Fatalf("client ip: %v", err)
	}
	body, _ := io.ReadAll(response.Body)
	_ = response.Body.Close()
	if string(body) != "203.0.113.8" {
		t.Fatalf("client ip = %q", body)
	}

	response, err = stdhttp.Post("http://"+httpServer.Addr()+"/upload", "text/plain", strings.NewReader("12345"))
	if err != nil {
		t.Fatalf("upload: %v", err)
	}
	_ = response.Body.Close()
	if response.StatusCode != stdhttp.StatusRequestEntityTooLarge {
		t.Fatalf("upload status = %d", response.StatusCode)
	}

	requestDone := make(chan error, 1)
	go func() {
		response, err := stdhttp.Get("http://" + httpServer.Addr() + "/slow")
		if err == nil {
			defer response.Body.Close()
			body, readErr := io.ReadAll(response.Body)
			if readErr != nil {
				err = readErr
			} else if string(body) != "done" {
				err = io.ErrUnexpectedEOF
			}
		}
		requestDone <- err
	}()
	select {
	case <-entered:
	case <-time.After(time.Second):
		t.Fatal("slow request did not enter handler")
	}
	stopDone := make(chan error, 1)
	go func() { stopDone <- app.Stop(context.Background()) }()
	select {
	case err := <-stopDone:
		t.Fatalf("stop returned before request drained: %v", err)
	case <-time.After(20 * time.Millisecond):
	}
	releaseRequest()
	if err := <-requestDone; err != nil {
		t.Fatalf("slow request: %v", err)
	}
	if err := <-stopDone; err != nil {
		t.Fatalf("stop: %v", err)
	}
}
