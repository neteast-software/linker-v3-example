package permission

import (
	"slices"
	"testing"
)

func TestServiceAppliesRelationDelta(t *testing.T) {
	service := New()
	service.Assign("2", "http.app2.order.update", "console.order.list")
	if got := service.List("2"); !slices.Equal(got, []string{
		"console.order.list",
		"http.app2.order.update",
	}) {
		t.Fatalf("unexpected assigned resources: %#v", got)
	}

	service.Remove("2", "console.order.list")
	if got := service.List("2"); !slices.Equal(got, []string{"http.app2.order.update"}) {
		t.Fatalf("unexpected remaining resources: %#v", got)
	}
}
