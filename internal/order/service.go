package order

import (
	"fmt"
	"slices"
	"sync"
)

type Service struct {
	mu     sync.RWMutex
	orders []Order
}

func New() *Service {
	return &Service{orders: []Order{
		{ID: 1, Number: "NO-20260716-001", Status: "open", Amount: 12800},
		{ID: 2, Number: "NO-20260716-002", Status: "closed", Amount: 6800},
	}}
}

func (p *Service) List() []Order {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return slices.Clone(p.orders)
}

func (p *Service) Save(value Order) (Order, error) {
	if value.ID == 0 {
		return Order{}, fmt.Errorf("订单 ID 不能为空")
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	for index := range p.orders {
		if p.orders[index].ID != value.ID {
			continue
		}
		p.orders[index] = value
		return value, nil
	}
	return Order{}, fmt.Errorf("订单 %d 不存在", value.ID)
}
