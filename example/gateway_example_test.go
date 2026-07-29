package example

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/neteast-software/go-module/http/gateway/declaration"
	gatewaycomponent "github.com/neteast-software/go-module/http/gateway/linker"
	httpcore "github.com/neteast-software/go-module/http/gin"
	httpcomponent "github.com/neteast-software/go-module/http/gin/linker"
	prometheus "github.com/neteast-software/go-module/observe/metrics/prometheus/linker"
	opentelemetry "github.com/neteast-software/go-module/observe/tracing/opentelemetry/linker"
	linker "github.com/neteast-software/linker/v3"

	"linker-v3-example/internal/app"
	examplegateway "linker-v3-example/internal/gateway"
)

func TestGatewayRecommendedWorkingBackground(t *testing.T) {
	var firstCalls atomic.Uint64
	first := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		firstCalls.Add(1)
		writer.Header().Set("Content-Type", "text/plain")
		_, _ = io.WriteString(writer, "first:"+request.URL.Path)
	}))
	defer first.Close()

	slowStarted := make(chan struct{}, 1)
	slowRelease := make(chan struct{})
	second := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path == "/slow" {
			select {
			case slowStarted <- struct{}{}:
			default:
			}
			select {
			case <-slowRelease:
			case <-request.Context().Done():
				return
			}
		}
		writer.Header().Set("Content-Type", "text/plain")
		_, _ = io.WriteString(writer, "second:"+request.URL.Path)
	}))
	defer second.Close()

	address := availableGatewayAddress(t)
	managementAddress := availableGatewayAddress(t)
	source := &gatewayExampleSource{
		initial: gatewayExampleSetting(t, address, managementAddress, first.URL),
		updates: make(chan linker.SourceSnapshot, 1),
	}
	runtime, err := app.Gateway(source)
	if err != nil {
		t.Fatal(err)
	}
	stopped := false
	t.Cleanup(func() {
		if !stopped {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			_ = runtime.Stop(ctx)
		}
	})
	if err = runtime.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err = runtime.Health(context.Background()); err != nil {
		t.Fatalf("Gateway health: %v", err)
	}
	for _, path := range []string{"/live", "/ready", "/startup"} {
		if body := gatewayExampleGET(t, address, path); !strings.Contains(body, `"status":"ok"`) {
			t.Fatalf("%s body = %q", path, body)
		}
	}
	if body := gatewayExampleGET(t, address, "/example/profile"); body != "first:/profile" {
		t.Fatalf("initial body = %q", body)
	}

	source.updates <- linker.SourceSnapshot{
		Revision: "routes-2",
		Setting:  gatewayExampleSetting(t, address, managementAddress, second.URL),
	}
	waitFor(t, 3*time.Second, func() bool {
		return gatewayExampleGET(t, address, "/example/profile") == "second:/profile"
	})
	if firstCalls.Load() == 0 {
		t.Fatal("initial Gateway route was not exercised")
	}

	recorder, err := prometheus.Require(runtime)
	if err != nil {
		t.Fatal(err)
	}
	metrics := httptest.NewRecorder()
	recorder.Handler().ServeHTTP(metrics, httptest.NewRequest(http.MethodGet, "/metrics", nil))
	if content := metrics.Body.String(); !strings.Contains(content, "gateway_requests_total") ||
		!strings.Contains(content, "gateway_plan_reloads_total") {
		t.Fatalf("Gateway metrics missing:\n%s", content)
	}
	if content := gatewayExampleGET(t, managementAddress, "/metrics"); !strings.Contains(content, "gateway_requests_total") {
		t.Fatalf("Gateway scrape endpoint missing:\n%s", content)
	}
	provider, err := opentelemetry.Require(runtime)
	if err != nil {
		t.Fatal(err)
	}
	if provider.Memory() == nil || len(provider.Memory().Spans()) == 0 {
		t.Fatal("Gateway trace span missing")
	}

	responseDone := make(chan error, 1)
	go func() {
		response, requestErr := http.Get("http://" + address + "/example/slow")
		if requestErr == nil {
			_, requestErr = io.Copy(io.Discard, response.Body)
			requestErr = errors.Join(requestErr, response.Body.Close())
		}
		responseDone <- requestErr
	}()
	<-slowStarted
	stopDone := make(chan error, 1)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		stopDone <- runtime.Stop(ctx)
	}()
	select {
	case err = <-stopDone:
		t.Fatalf("Gateway did not wait for in-flight request: %v", err)
	case <-time.After(50 * time.Millisecond):
	}
	close(slowRelease)
	if err = <-responseDone; err != nil {
		t.Fatalf("slow request: %v", err)
	}
	if err = <-stopDone; err != nil {
		t.Fatalf("Gateway graceful shutdown: %v", err)
	}
	stopped = true
}

func gatewayExampleSetting(t *testing.T, address, managementAddress, origin string) linker.Setting {
	t.Helper()
	document := examplegateway.Document(
		examplegateway.Strip(examplegateway.URL("local/public", "/example/**", origin), 1),
	)
	routes, err := declaration.Encode(document)
	if err != nil {
		t.Fatal(err)
	}
	config := gatewaycomponent.DefaultConfig()
	config.Address = address
	config.Health = gatewaycomponent.StandardHealth()
	management := httpcore.DefaultConfig()
	management.Addr = managementAddress
	return linker.NewSetting(map[linker.Namespace][]byte{
		httpcomponent.Namespace:          mustGatewayJSON(t, management),
		gatewaycomponent.Namespace:       mustGatewayJSON(t, config),
		gatewaycomponent.RoutesNamespace: routes,
		prometheus.Namespace: mustGatewayJSON(t, prometheus.Config{
			Enabled:   true,
			Namespace: "linker_v3_example",
			ConstLabels: map[string]string{
				"service": "gateway-example",
			},
		}),
		opentelemetry.Namespace: mustGatewayJSON(t, opentelemetry.InMemory("gateway-example")),
	})
}

func mustGatewayJSON(t *testing.T, value any) []byte {
	t.Helper()
	content, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	return content
}

func availableGatewayAddress(t *testing.T) string {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	address := listener.Addr().String()
	if err = listener.Close(); err != nil {
		t.Fatal(err)
	}
	return address
}

func gatewayExampleGET(t *testing.T, address, path string) string {
	t.Helper()
	response, err := http.Get("http://" + address + path)
	if err != nil {
		return ""
	}
	defer response.Body.Close()
	content, err := io.ReadAll(response.Body)
	if err != nil || response.StatusCode != http.StatusOK {
		return ""
	}
	return string(content)
}

type gatewayExampleSource struct {
	initial linker.Setting
	updates chan linker.SourceSnapshot
}

func (p *gatewayExampleSource) Name() string {
	return "example/gateway/live"
}

func (p *gatewayExampleSource) Load(context.Context, linker.BootstrapContext) (linker.Setting, error) {
	return p.initial, nil
}

func (p *gatewayExampleSource) Watch(
	ctx context.Context,
	_ linker.BootstrapContext,
	publish linker.SourcePublish,
) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case snapshot := <-p.updates:
			if err := publish(snapshot); err != nil {
				return err
			}
		}
	}
}

var _ linker.SourceWatcher = (*gatewayExampleSource)(nil)
