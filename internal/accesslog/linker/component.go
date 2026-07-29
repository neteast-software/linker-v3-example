package accesslog

import (
	"context"
	"fmt"
	"strconv"

	linker "github.com/neteast-software/linker/v3"

	accesslog "linker-v3-example/internal/accesslog"
)

// ID 是 accesslog 生命周期适配器的组件标识。
const ID linker.ID = "example/accesslog"

const AssetQueue linker.AssetKind = "audit/accesslog/queue"

// Component 只接管 Log worker 的启动、健康和排空。
type Component struct {
	log *accesslog.Log
}

// New 创建 accesslog 生命周期适配器。
func New(log *accesslog.Log) *Component {
	return &Component{log: log}
}

func (p *Component) Identity() linker.ID {
	return ID
}

func (p *Component) Start(context.Context) error {
	if p == nil || p.log == nil {
		return fmt.Errorf("accesslog component 缺少 Log")
	}
	return p.log.Open()
}

func (p *Component) Health(context.Context) error {
	if p == nil || p.log == nil {
		return fmt.Errorf("accesslog component 缺少 Log")
	}
	return p.log.Health()
}

func (p *Component) Close(ctx context.Context) error {
	if p == nil || p.log == nil {
		return nil
	}
	return p.log.Close(ctx)
}

func (p *Component) Assets(context.Context, linker.Runtime) ([]linker.Asset, error) {
	if p == nil || p.log == nil {
		return nil, nil
	}
	stats := p.log.Stats()
	return []linker.Asset{
		linker.NewAsset(AssetQueue, ID, stats.Capacity).Describe("operateLog", map[string]string{
			"capacity": strconv.Itoa(stats.Capacity),
		}),
	}, nil
}

func (p *Component) AssetPolicies() linker.AssetPolicies {
	return linker.AssetPolicies{linker.Observe(AssetQueue)}
}

var (
	_ linker.Component           = (*Component)(nil)
	_ linker.Starter             = (*Component)(nil)
	_ linker.HealthChecker       = (*Component)(nil)
	_ linker.Closer              = (*Component)(nil)
	_ linker.AssetProvider       = (*Component)(nil)
	_ linker.AssetPolicyProvider = (*Component)(nil)
)
