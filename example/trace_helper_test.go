package example

import (
	"testing"

	telemetry "github.com/neteast-software/go-module/observe/tracing/opentelemetry"
)

const (
	exampleTraceID = "0123456789abcdef0123456789abcdef"
	exampleSpanID  = "0123456789abcdef"
	httpMQTraceID  = "fedcba9876543210fedcba9876543210"
	httpMQSpanID   = "fedcba9876543210"
	sseTraceID     = "11111111111111111111111111111111"
	sseSpanID      = "1111111111111111"
)

func requireTraceSpan(t *testing.T, spans []telemetry.Snapshot, name string, traceID string, parentSpanID string) telemetry.Snapshot {
	t.Helper()
	for _, span := range spans {
		if span.Name != name || span.TraceID != traceID {
			continue
		}
		if parentSpanID != "" && span.ParentSpanID != parentSpanID {
			continue
		}
		return span
	}
	t.Fatalf("span name=%q trace=%q parent=%q missing: %#v", name, traceID, parentSpanID, spans)
	return telemetry.Snapshot{}
}
