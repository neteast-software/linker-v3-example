package accesslog

import (
	"context"
	"sync"
)

// MemoryPublisher 保存测试与本地 profile 的审计记录。
type MemoryPublisher struct {
	mutex   sync.Mutex
	records []Record
}

// Memory 创建有界队列后端使用的内存 publisher。
func Memory() *MemoryPublisher {
	return &MemoryPublisher{}
}

func (p *MemoryPublisher) Publish(_ context.Context, record Record) error {
	p.mutex.Lock()
	p.records = append(p.records, record)
	p.mutex.Unlock()
	return nil
}

// Records 返回不共享底层切片的审计记录。
func (p *MemoryPublisher) Records() []Record {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	return append([]Record(nil), p.records...)
}
