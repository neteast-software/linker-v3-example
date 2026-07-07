package graph

import (
	graphlinker "github.com/neteast-software/go-module/graph/naive/linker"
	linker "github.com/neteast-software/linker/v3"

	_ "linker-v3-example/internal/route/graph" // route 声明随组件进入编译
)

const ID linker.ID = graphlinker.ID

func NewComponent() linker.Component {
	return graphlinker.New()
}
