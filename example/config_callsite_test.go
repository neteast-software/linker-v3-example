package example

import (
	"testing"

	env "github.com/neteast-software/go-module/config/env/linker"
	yaml "github.com/neteast-software/go-module/config/yaml/linker"
	server "github.com/neteast-software/go-module/linker/server"
	nacos "github.com/neteast-software/go-module/registry/nacos/linker"
)

func TestConfigurationCallSiteCompiles(t *testing.T) {
	app := server.New(
		server.Config(
			yaml.File("config/app.yaml"),
			nacos.Config("app.yaml", nacos.Group("LINKER")),
			env.Prefix("APP_"),
		),
	)
	if app == nil {
		t.Fatal("server app is nil")
	}
}
