package main

import (
	"context"
	"fmt"
	"os"

	"linker-v3-example/internal/app"
	"linker-v3-example/internal/config"
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
	if err := run(ctx); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "linker-v3-example 启动失败：%v\n", err)
		os.Exit(1)
	}
}

func isPlanCommand(args []string) bool {
	return len(args) > 1 && args[1] == "--plan"
}

func run(ctx context.Context) error {
	cfg, err := config.Load(ctx)
	if err != nil {
		return fmt.Errorf("配置加载失败: %w", err)
	}
	serverApp, err := app.New(cfg)
	if err != nil {
		return fmt.Errorf("配置错误: %w", err)
	}
	if err = serverApp.Run(ctx); err != nil {
		return fmt.Errorf("运行失败: %w", err)
	}
	return nil
}
