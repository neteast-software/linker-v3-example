package order

import "testing"

func TestServiceSaveKeepsListIsolated(t *testing.T) {
	service := New()
	items := service.List()
	items[0].Number = "changed-outside"
	if service.List()[0].Number == "changed-outside" {
		t.Fatal("List 不应暴露内部切片")
	}

	saved, err := service.Save(Order{
		ID: 1, Number: "NO-20260716-009", Status: "closed", Amount: 16800,
	})
	if err != nil {
		t.Fatalf("save: %v", err)
	}
	if saved.Number != "NO-20260716-009" || service.List()[0] != saved {
		t.Fatalf("unexpected saved order: %#v", saved)
	}
}

func TestServiceRejectsUnknownOrder(t *testing.T) {
	if _, err := New().Save(Order{ID: 99}); err == nil {
		t.Fatal("不存在的订单不能被静默创建")
	}
}
