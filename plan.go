package main

import (
	"encoding/json"
	"io"

	"linker-v3-example/internal/app"
	"linker-v3-example/internal/config"
)

const planOnlyPostgreSQLPassword = "linker-v3-example-plan-only"

func printPlan(output io.Writer) error {
	cfg := config.FromEnv()
	if cfg.PostgreSQL.Password == "" {
		cfg.PostgreSQL.Password = planOnlyPostgreSQLPassword
	}
	serverApp, err := app.New(cfg)
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(output)
	encoder.SetIndent("", "  ")
	return encoder.Encode(serverApp.Plan())
}
