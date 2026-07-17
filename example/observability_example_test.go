package example

import (
	"context"
	"io"
	stdhttp "net/http"
	"strings"
	"testing"
	"time"

	http "github.com/neteast-software/go-module/http/gin/linker"
	server "github.com/neteast-software/go-module/linker/server"
	consumer "github.com/neteast-software/go-module/mq/consumer"
	mq "github.com/neteast-software/go-module/mq/consumer/linker"
	prometheus "github.com/neteast-software/go-module/observe/metrics/prometheus/linker"
	trace "github.com/neteast-software/go-module/observe/tracing"
	"github.com/neteast-software/go-module/observe/tracing/mq/consumer"
	opentelemetry "github.com/neteast-software/go-module/observe/tracing/opentelemetry/linker"
	rpccore "github.com/neteast-software/go-module/rpc/grpc"
	rpc "github.com/neteast-software/go-module/rpc/grpc/linker"
	cron "github.com/neteast-software/go-module/scheduler/cron"
	schedule "github.com/neteast-software/go-module/scheduler/cron/linker"
	linker "github.com/neteast-software/linker/v3"
	stdgrpc "google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"

	ttsclient "linker-v3-example/internal/client/tts"
	ttsrpc "linker-v3-example/internal/rpc/tts"
)

func TestFrameworkObservabilityExample(t *testing.T) {
	grpcAddr := freeLocalAddr(t)
	httpConfig := http.DefaultConfig()
	httpConfig.Addr = "127.0.0.1:0"
	metricComponent := prometheus.New(prometheus.WithConfig(prometheus.Config{
		Enabled: true, Namespace: "linker_v3_example", ConstLabels: map[string]string{"service": "observability-example"},
	}))
	traceComponent := opentelemetry.New(opentelemetry.WithConfig(opentelemetry.InMemory("linker-v3-example")))
	mqTrace := make(chan trace.Trace, 1)
	item := consumer.New("trace", consumer.HandlerFunc(func(ctx context.Context, _ consumer.Message) error {
		current, _ := trace.FromContext(ctx)
		mqTrace <- current
		return nil
	}), consumer.WithTopic("trace.message"))
	cronTrace := make(chan trace.Trace, 1)
	job := cron.NewJob("trace.job", "@every 10ms", cron.HandlerFunc(func(ctx context.Context) error {
		current, _ := trace.FromContext(ctx)
		select {
		case cronTrace <- current:
		default:
		}
		return nil
	}))
	settings := map[linker.Namespace]any{
		rpc.Namespace: rpccore.ServerConfig{Addr: grpcAddr},
		linker.Namespace(ttsclient.ID): rpccore.ClientConfig{
			Target: grpcAddr, Insecure: true, Timeout: time.Second,
		},
	}
	app := server.New(
		server.WithoutStartupLog(),
		server.WithoutEvent(),
		server.WithoutNotice(),
		server.WithoutAudit(),
		server.WithHTTP(httpConfig),
		server.Config(businessConfigSource(t, settings)),
		server.WithMetrics(metricComponent),
		server.WithHTTPRoutes(
			http.POST("trace/grpc", traceGRPC),
			http.POST("trace/mq", traceMQ),
		),
		server.WithComponents(
			traceComponent,
			mq.New(mq.WithConsumers(item)),
			schedule.New(
				schedule.WithStore(cron.NewMemoryStore()),
				schedule.WithJobs(job),
			),
			rpc.New(
				rpc.WithRegisters(func(server *stdgrpc.Server) { ttsrpc.Register(server, traceTTS{}) }),
			),
			ttsclient.Provider(),
		),
	)
	if err := app.Start(context.Background()); err != nil {
		t.Fatalf("start: %v", err)
	}
	t.Cleanup(func() { _ = app.Stop(context.Background()) })
	httpServer, err := http.RequireServer(app)
	if err != nil {
		t.Fatalf("http server: %v", err)
	}
	baseURL := "http://" + httpServer.Addr()

	grpcResponse := traceRequest(t, baseURL+"/trace/grpc", exampleTraceID, exampleSpanID)
	if grpcResponse.Header.Get(trace.HeaderTraceID) != exampleTraceID {
		t.Fatalf("gRPC response trace = %q", grpcResponse.Header.Get(trace.HeaderTraceID))
	}
	_ = grpcResponse.Body.Close()
	mqResponse := traceRequest(t, baseURL+"/trace/mq", httpMQTraceID, httpMQSpanID)
	if mqResponse.Header.Get(trace.HeaderTraceID) != httpMQTraceID {
		t.Fatalf("MQ response trace = %q", mqResponse.Header.Get(trace.HeaderTraceID))
	}
	_ = mqResponse.Body.Close()

	select {
	case trace := <-mqTrace:
		if trace.TraceID != httpMQTraceID {
			t.Fatalf("MQ handler trace = %#v", trace)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("等待 MQ trace 超时")
	}
	var scheduled trace.Trace
	select {
	case scheduled = <-cronTrace:
	case <-time.After(2 * time.Second):
		t.Fatal("等待 cron trace 超时")
	}

	provider, err := opentelemetry.Require(app)
	if err != nil || provider.Memory() == nil {
		t.Fatalf("memory provider = %v, %v", provider, err)
	}
	waitFor(t, time.Second, func() bool { return len(provider.Memory().Spans()) >= 6 })
	spans := provider.Memory().Spans()
	httpGRPC := requireTraceSpan(t, spans, "HTTP POST /trace/grpc", exampleTraceID, exampleSpanID)
	grpcClient := requireTraceSpan(t, spans, "gRPC "+ttsrpc.FullMethodTranscribe, exampleTraceID, httpGRPC.SpanID)
	_ = requireTraceSpan(t, spans, "gRPC "+ttsrpc.FullMethodTranscribe, exampleTraceID, grpcClient.SpanID)
	httpMQ := requireTraceSpan(t, spans, "HTTP POST /trace/mq", httpMQTraceID, httpMQSpanID)
	_ = requireTraceSpan(t, spans, "MQ trace.message", httpMQTraceID, httpMQ.SpanID)
	_ = requireTraceSpan(t, spans, "cron trace.job", scheduled.TraceID, "")

	metricsText := getRaw(t, baseURL+"/metrics")
	for _, name := range []string{
		"linker_v3_example_http_requests_total",
		"linker_v3_example_grpc_server_requests_total",
		"linker_v3_example_grpc_client_requests_total",
		"linker_v3_example_mq_consumer_messages_total",
		"linker_v3_example_scheduler_cron_runs_total",
		"linker_v3_example_linker_component_lifecycle_total",
	} {
		if !strings.Contains(metricsText, name) {
			t.Fatalf("metric %s missing:\n%s", name, metricsText)
		}
	}
}

func traceGRPC(c *http.Context) {
	runtime, ok := http.Runtime(c)
	if !ok {
		c.Status(stdhttp.StatusServiceUnavailable)
		return
	}
	client, err := ttsclient.Require(runtime)
	if err != nil {
		c.Status(stdhttp.StatusServiceUnavailable)
		return
	}
	result, err := client.Transcribe(c.Request.Context(), "trace")
	if err != nil {
		c.Status(stdhttp.StatusBadGateway)
		return
	}
	c.JSON(stdhttp.StatusOK, map[string]string{"result": result})
}

func traceMQ(c *http.Context) {
	runtime, ok := http.Runtime(c)
	if !ok {
		c.Status(stdhttp.StatusServiceUnavailable)
		return
	}
	executor, err := mq.Require(runtime, "trace")
	if err != nil {
		c.Status(stdhttp.StatusServiceUnavailable)
		return
	}
	message := tracing.InjectMessage(c.Request.Context(), consumer.Message{Topic: "trace.message", Body: []byte("trace")})
	if err = executor.Submit(c.Request.Context(), message); err != nil {
		c.Status(stdhttp.StatusServiceUnavailable)
		return
	}
	c.Status(stdhttp.StatusAccepted)
}

func traceRequest(t *testing.T, url string, traceID string, spanID string) *stdhttp.Response {
	t.Helper()
	req, err := stdhttp.NewRequest(stdhttp.MethodPost, url, nil)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	req.Header.Set(trace.HeaderTraceID, traceID)
	req.Header.Set(trace.HeaderSpanID, spanID)
	resp, err := stdhttp.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("do request: %v", err)
	}
	if resp.StatusCode >= stdhttp.StatusBadRequest {
		body, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		t.Fatalf("status=%d body=%s", resp.StatusCode, body)
	}
	return resp
}

type traceTTS struct{}

func (traceTTS) Transcribe(_ context.Context, req *wrapperspb.StringValue) (*wrapperspb.StringValue, error) {
	return wrapperspb.String("trace:" + req.Value), nil
}
