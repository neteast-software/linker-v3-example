package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io"

	linker "github.com/neteast-software/linker/v3"

	"linker-v3-example/internal/app"
	user "linker-v3-example/internal/user/linker"
)

func printPlan(output io.Writer) error {
	secret, err := planSecretSource()
	if err != nil {
		return err
	}
	sources, err := configSources(secret)
	if err != nil {
		return err
	}
	serverApp := app.New(sources...)
	if err = serverApp.Prepare(context.Background()); err != nil {
		return err
	}
	encoder := json.NewEncoder(output)
	encoder.SetIndent("", "  ")
	return encoder.Encode(serverApp.Plan())
}

func planSecretSource() (linker.Source, error) {
	value := make([]byte, 32)
	if _, err := rand.Read(value); err != nil {
		return nil, err
	}
	content, err := json.Marshal(user.Config{TokenKey: hex.EncodeToString(value)})
	if err != nil {
		return nil, err
	}
	return linker.MapSource{
		Label: "config/plan-secret",
		Setting: linker.NewSetting(map[linker.Namespace][]byte{
			user.Namespace: content,
		}),
	}, nil
}
