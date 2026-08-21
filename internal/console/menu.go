package console

import (
	_ "embed"
	"slices"
)

//go:embed menu.yaml
var menuDeclaration []byte

// MenuDeclaration 返回项目唯一的 Graph Console 菜单声明。
func MenuDeclaration() []byte {
	return slices.Clone(menuDeclaration)
}
