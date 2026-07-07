package main

import (
	"context"
	"fmt"
	"os"

	"linker-v3-example/internal/app"
	"linker-v3-example/internal/config"
)

func main() {
	if err := run(context.Background()); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "linker-v3-example 启动失败：%v\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	serverApp, err := app.New(config.FromEnv())
	if err != nil {
		return fmt.Errorf("配置错误: %w", err)
	}
	if err = serverApp.Run(ctx); err != nil {
		return fmt.Errorf("运行失败: %w", err)
	}
	return nil
}
