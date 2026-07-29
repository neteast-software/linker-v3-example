package gateway

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/neteast-software/go-module/http/gateway/declaration"
	gatewaycomponent "github.com/neteast-software/go-module/http/gateway/linker"
	linker "github.com/neteast-software/linker/v3"
)

const maximumRoutesBytes = 4 << 20

type routesSource struct {
	path string
}

// Routes 把独立 Gateway 声明文件投影到组件配置 namespace。
func Routes(path string) linker.Source {
	return routesSource{path: path}
}

func (p routesSource) Name() string {
	return "example/gateway/routes"
}

func (p routesSource) Load(ctx context.Context, _ linker.BootstrapContext) (linker.Setting, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	info, err := os.Lstat(p.path)
	if err != nil {
		return nil, fmt.Errorf("读取 Gateway 声明信息失败: %w", err)
	}
	if info.Mode()&os.ModeSymlink != 0 || !info.Mode().IsRegular() {
		return nil, fmt.Errorf("Gateway 声明必须是普通文件")
	}
	file, err := os.Open(p.path)
	if err != nil {
		return nil, fmt.Errorf("打开 Gateway 声明失败: %w", err)
	}
	defer file.Close()
	content, err := io.ReadAll(io.LimitReader(file, maximumRoutesBytes+1))
	if err != nil {
		return nil, fmt.Errorf("读取 Gateway 声明失败: %w", err)
	}
	if len(content) > maximumRoutesBytes {
		return nil, fmt.Errorf("Gateway 声明超过 4 MiB")
	}
	document, err := declaration.Decode(content)
	if err != nil {
		return nil, err
	}
	normalized, err := declaration.Encode(document)
	if err != nil {
		return nil, err
	}
	return linker.NewSetting(map[linker.Namespace][]byte{
		gatewaycomponent.RoutesNamespace: normalized,
	}), nil
}

var _ linker.Source = routesSource{}
