package accesslog

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

const publishTimeout = 2 * time.Second

// Log 是 Gateway Observer 与异步审计投递 worker 的共同 owner。
type Log struct {
	publisher Publisher
	queue     chan Record
	stop      chan struct{}
	done      chan struct{}
	now       func() time.Time
	prefix    string

	started   atomic.Bool
	closing   atomic.Bool
	sequence  atomic.Uint64
	published atomic.Uint64
	failed    atomic.Uint64
	dropped   atomic.Uint64

	enqueueMutex sync.RWMutex
	issueMutex   sync.RWMutex
	issue        error
}

// New 创建使用有界队列的 accesslog；capacity 必须大于零。
func New(publisher Publisher, capacity int) (*Log, error) {
	if capacity <= 0 {
		return nil, fmt.Errorf("accesslog queue capacity 必须大于零")
	}
	if publisher == nil {
		publisher = Discard()
	}
	value := make([]byte, 8)
	if _, err := rand.Read(value); err != nil {
		return nil, fmt.Errorf("accesslog event prefix 生成失败: %w", err)
	}
	return &Log{
		publisher: publisher,
		queue:     make(chan Record, capacity),
		stop:      make(chan struct{}),
		done:      make(chan struct{}),
		now:       time.Now,
		prefix:    hex.EncodeToString(value),
	}, nil
}

func (p *Log) Identity() string {
	return "accesslog"
}

// Open 启动唯一异步投递 worker。
func (p *Log) Open() error {
	if p == nil {
		return fmt.Errorf("accesslog 不能为空")
	}
	if p.closing.Load() {
		return fmt.Errorf("accesslog 已进入关闭流程")
	}
	if p.started.CompareAndSwap(false, true) {
		go p.run()
	}
	return nil
}

func (p *Log) run() {
	defer close(p.done)
	for {
		select {
		case record := <-p.queue:
			p.publish(record)
		case <-p.stop:
			for {
				select {
				case record := <-p.queue:
					p.publish(record)
				default:
					return
				}
			}
		}
	}
}

func (p *Log) publish(record Record) {
	ctx, cancel := context.WithTimeout(context.Background(), publishTimeout)
	defer cancel()
	if err := p.publisher.Publish(ctx, record); err != nil {
		p.failed.Add(1)
		p.setIssue(fmt.Errorf("accesslog 最近一次投递失败: %w", err))
		return
	}
	p.published.Add(1)
	p.setIssue(nil)
}

// Close 停止接收新记录并排空队列。
func (p *Log) Close(ctx context.Context) error {
	if p == nil {
		return nil
	}
	if !p.started.Load() {
		if err := p.Open(); err != nil {
			return err
		}
	}
	p.enqueueMutex.Lock()
	if p.closing.CompareAndSwap(false, true) {
		close(p.stop)
	}
	p.enqueueMutex.Unlock()
	select {
	case <-p.done:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("accesslog 队列排空未完成: %w", ctx.Err())
	}
}

// Health 表达最近一次异步投递结果。
func (p *Log) Health() error {
	if p == nil {
		return errors.New("accesslog 未初始化")
	}
	p.issueMutex.RLock()
	defer p.issueMutex.RUnlock()
	return p.issue
}

// Stats 返回当前投递计数和有界队列状态。
func (p *Log) Stats() Stats {
	if p == nil {
		return Stats{}
	}
	return Stats{
		Published: p.published.Load(),
		Failed:    p.failed.Load(),
		Dropped:   p.dropped.Load(),
		Queued:    len(p.queue),
		Capacity:  cap(p.queue),
	}
}

func (p *Log) setIssue(err error) {
	p.issueMutex.Lock()
	p.issue = err
	p.issueMutex.Unlock()
}

func (p *Log) nextEventID() string {
	return fmt.Sprintf("%s-%016x", p.prefix, p.sequence.Add(1))
}
