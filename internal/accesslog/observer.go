package accesslog

import (
	"context"
	"errors"
	"net/http"
	"sync"

	"github.com/neteast-software/go-module/http/gateway"
)

// Start 实现 Gateway Observer。
func (p *Log) Start(
	ctx context.Context,
	start gateway.RequestStart,
) (context.Context, gateway.Observation) {
	if p == nil {
		return ctx, nil
	}
	return ctx, &observation{owner: p, route: start.Route, method: start.Method}
}

type observation struct {
	owner  *Log
	route  string
	method string
	once   sync.Once
}

func (p *observation) Finish(result gateway.RequestResult) {
	if p == nil || p.owner == nil {
		return
	}
	p.once.Do(func() {
		p.owner.enqueueMutex.RLock()
		defer p.owner.enqueueMutex.RUnlock()
		if p.owner.closing.Load() {
			p.owner.dropped.Add(1)
			return
		}
		record := Record{
			EventID:        p.owner.nextEventID(),
			SendTime:       p.owner.now(),
			Topic:          "operateLog",
			RequestURL:     p.route,
			RequestMethod:  p.method,
			Success:        result.Status >= http.StatusOK && result.Status < http.StatusBadRequest,
			SimpleErrorMsg: result.Failure,
			ExecutionTime:  max(result.Duration.Milliseconds(), 0),
		}
		select {
		case p.owner.queue <- record:
		default:
			p.owner.dropped.Add(1)
			p.owner.setIssue(errors.New("accesslog 有界队列已满，记录被丢弃"))
		}
	})
}

var (
	_ gateway.Observer    = (*Log)(nil)
	_ gateway.Observation = (*observation)(nil)
)
