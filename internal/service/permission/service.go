package permission

import (
	"slices"
	"sync"
)

type Service struct {
	mu        sync.RWMutex
	resources map[string][]string
}

func New() *Service {
	return &Service{resources: map[string][]string{
		"1": {"console.order.list", "http.app2.order.update"},
		"2": {"console.order.list"},
	}}
}

func (p *Service) List(role string) []string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return slices.Clone(p.resources[role])
}

func (p *Service) Assign(role string, resources ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	current := make(map[string]struct{}, len(p.resources[role])+len(resources))
	for _, resource := range p.resources[role] {
		current[resource] = struct{}{}
	}
	for _, resource := range resources {
		current[resource] = struct{}{}
	}
	p.resources[role] = sorted(current)
}

func (p *Service) Remove(role string, resources ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	current := make(map[string]struct{}, len(p.resources[role]))
	for _, resource := range p.resources[role] {
		current[resource] = struct{}{}
	}
	for _, resource := range resources {
		delete(current, resource)
	}
	p.resources[role] = sorted(current)
}

func sorted(values map[string]struct{}) []string {
	ret := make([]string, 0, len(values))
	for value := range values {
		ret = append(ret, value)
	}
	slices.Sort(ret)
	return ret
}
