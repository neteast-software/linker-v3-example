package main

import (
	"context"
	"fmt"
	"os"

	linker "github.com/neteast-software/linker/v3"

	"linker-v3-example/internal/app"
)

func main() {
	ctx := context.Background()
	if isPlanCommand(os.Args) {
		if err := printPlan(os.Stdout); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "linker-v3-example 计划输出失败：%v\n", err)
			os.Exit(1)
		}
		return
	}
	if nacosProfile, ok := gatewayCommand(os.Args); ok {
		if err := runGateway(ctx, nacosProfile); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "linker-v3-example Gateway 启动失败：%v\n", err)
			os.Exit(1)
		}
		return
	}
	if err := run(ctx); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "linker-v3-example 启动失败：%v\n", err)
		os.Exit(1)
	}
}

func gatewayCommand(args []string) (bool, bool) {
	if len(args) < 2 {
		return false, false
	}
	switch args[1] {
	case "--gateway":
		return false, true
	case "--gateway-nacos":
		return true, true
	default:
		return false, false
	}
}

func isPlanCommand(args []string) bool {
	return len(args) > 1 && args[1] == "--plan"
}

func run(ctx context.Context) error {
	sources, err := configSources()
	if err != nil {
		return fmt.Errorf("配置来源错误: %w", err)
	}
	serverApp := app.New(sources...)
	if err := serverApp.Run(ctx); err != nil {
		return fmt.Errorf("运行失败: %w", err)
	}
	return nil
}

func runGateway(ctx context.Context, nacosProfile bool) error {
	sources, err := gatewaySources(nacosProfile)
	if err != nil {
		return fmt.Errorf("Gateway 配置来源错误: %w", err)
	}
	var gatewayApp *linker.App
	if nacosProfile {
		gatewayApp, err = app.GatewayNacos(sources...)
	} else {
		gatewayApp, err = app.Gateway(sources...)
	}
	if err != nil {
		return fmt.Errorf("Gateway 装配失败: %w", err)
	}
	if err = gatewayApp.Run(ctx); err != nil {
		return fmt.Errorf("Gateway 运行失败: %w", err)
	}
	return nil
}
