package inspection

import (
	"errors"
	"testing"
)

func TestStatusOwnsDefinitionAndTextBoundary(t *testing.T) {
	status, err := ParseStatus(" OPEN ")
	if err != nil || status != Open {
		t.Fatalf("parse = %q, %v", status, err)
	}
	definition, ok := status.Definition()
	if !ok || definition.Name != "待处理" {
		t.Fatalf("definition = %#v, %t", definition, ok)
	}
	var decoded Status
	if err = decoded.UnmarshalText([]byte("done")); err != nil || decoded != Done {
		t.Fatalf("decoded = %q, %v", decoded, err)
	}
	if _, err = ParseStatus("paused"); !errors.Is(err, ErrStatusInvalid) {
		t.Fatalf("invalid parse = %v", err)
	}
}

func TestDefinitionsReturnCopyInStableOrder(t *testing.T) {
	items := Definitions()
	items[0].Name = "changed"
	next := Definitions()
	if len(next) != 2 || next[0].Status != Open || next[0].Name != "待处理" || next[1].Status != Done {
		t.Fatalf("definitions = %#v", next)
	}
}
