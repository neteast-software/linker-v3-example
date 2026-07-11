package example

import (
	"context"
	"io"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	notice "github.com/neteast-software/go-module/fault/notice"
	server "github.com/neteast-software/go-module/linker/server"
	prometheus "github.com/neteast-software/go-module/observe/metrics/prometheus/linker"
	linker "github.com/neteast-software/linker/v3"
)

func TestFaultRecoveryNoticeMetricsExample(t *testing.T) {
	sender := notice.NewMemorySender()
	reporter := &noticeReporter{}
	metricComponent := prometheus.New(prometheus.WithConfig(prometheus.Config{
		Enabled: true, Namespace: "linker_v3_example", ConstLabels: map[string]string{"service": "fault-example"},
	}))
	app := server.New(
		server.WithoutStartupLog(),
		server.WithoutHTTP(),
		server.WithoutAudit(),
		server.WithoutEvent(),
		server.WithMetrics(metricComponent),
		server.WithNoticeSenders(sender),
		server.WithComponents(reporter),
	)
	if err := app.Start(context.Background()); err != nil {
		t.Fatalf("start: %v", err)
	}
	t.Cleanup(func() { _ = app.Stop(context.Background()) })

	if err := linker.Report(context.Background(), reporter.runtime, linker.Detected("provider_unavailable", nil)); err != nil {
		t.Fatalf("detected: %v", err)
	}
	if err := linker.Report(context.Background(), reporter.runtime, linker.Recovering("provider_unavailable", 1, "1s 后重试", nil)); err != nil {
		t.Fatalf("recovering: %v", err)
	}
	if err := linker.Report(context.Background(), reporter.runtime, linker.Recovered("provider_unavailable")); err != nil {
		t.Fatalf("recovered: %v", err)
	}
	waitFor(t, 2*time.Second, func() bool { return len(sender.Notices()) == 3 })
	waitFor(t, 2*time.Second, func() bool {
		return strings.Contains(scrapeMetrics(t, metricComponent), "linker_v3_example_linker_fault_notice_results_total")
	})

	text := scrapeMetrics(t, metricComponent)
	for _, value := range []string{
		"linker_v3_example_linker_fault_transition_total",
		"linker_v3_example_linker_fault_recovery_seconds",
		"linker_v3_example_linker_component_fault_active",
		"linker_v3_example_linker_fault_notice_results_total",
		`state="recovering"`,
		`state="recovered"`,
		`status="sent"`,
	} {
		if !strings.Contains(text, value) {
			t.Fatalf("metrics missing %q:\n%s", value, text)
		}
	}
}

func scrapeMetrics(t *testing.T, component *prometheus.Component) string {
	t.Helper()
	response := httptest.NewRecorder()
	component.Handler().ServeHTTP(response, httptest.NewRequest("GET", "/metrics", nil))
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("read metrics: %v", err)
	}
	return string(body)
}
