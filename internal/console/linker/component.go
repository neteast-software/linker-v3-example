package console

import (
	"slices"

	"github.com/neteast-software/go-module/acl"
	graphconsole "github.com/neteast-software/go-module/graph/console/linker"

	console "linker-v3-example/internal/console"
	"linker-v3-example/internal/console/dashboard"
	user "linker-v3-example/internal/user"
)

func New(options ...graphconsole.Option) *graphconsole.Component {
	provider := console.New()
	defaults := []graphconsole.Option{
		graphconsole.ConfigureFrom(user.AuthKey(), provider.Configure),
		graphconsole.WithEntry(console.Entry()),
		graphconsole.WithMenu(console.Menu()),
		graphconsole.WithPage("dashboard", dashboard.Page()),
		graphconsole.WithResources(
			acl.NewResource(console.Dashboard, acl.Scope("console", 0, "后台工作台", acl.Read)),
		),
		graphconsole.WithProvider(provider),
	}
	return graphconsole.New(slices.Concat(defaults, options)...)
}
