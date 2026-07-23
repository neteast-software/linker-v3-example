package example

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestProductionDatabaseStaysRepository(t *testing.T) {
	err := filepath.WalkDir("../internal", func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() || filepath.Ext(path) != ".go" || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		source, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		file, parseErr := parser.ParseFile(token.NewFileSet(), path, source, 0)
		if parseErr != nil {
			return parseErr
		}
		assertRepositoryBoundary(t, path, file)
		return nil
	})
	if err != nil {
		t.Fatalf("扫描生产数据库边界: %v", err)
	}
}

func assertRepositoryBoundary(t *testing.T, path string, file *ast.File) {
	t.Helper()
	ast.Inspect(file, func(node ast.Node) bool {
		switch value := node.(type) {
		case *ast.Field:
			if value.Tag != nil {
				assertNoDatabaseOwnedSemantics(t, path, unquote(t, path, value.Tag.Value))
			}
		case *ast.BasicLit:
			if value.Kind == token.STRING {
				assertNoDatabaseOwnedSemantics(t, path, unquote(t, path, value.Value))
			}
		case *ast.SelectorExpr:
			switch value.Sel.Name {
			case "SetupJoinTable":
				t.Fatalf("%s 不应由业务生产代码调用 GORM %s 承载数据库侧关系", path, value.Sel.Name)
			}
		}
		return true
	})
}

func assertNoDatabaseOwnedSemantics(t *testing.T, path, value string) {
	t.Helper()
	normalized := strings.ToLower(strings.Join(strings.Fields(value), " "))
	for _, forbidden := range []string{
		"foreignkey:",
		"references:",
		"many2many:",
		"constraint:on",
		"foreign key",
		"create function",
		"create or replace function",
		"create procedure",
		"create or replace procedure",
		"create trigger",
		"create event",
		"create extension",
		"create materialized view",
		"create view",
		"create policy",
		"create rule",
		"execute function",
		"execute procedure",
		"generated always",
	} {
		if strings.Contains(normalized, forbidden) {
			t.Fatalf("%s 包含数据库侧业务语义 %q", path, forbidden)
		}
	}
}

func unquote(t *testing.T, path, value string) string {
	t.Helper()
	result, err := strconv.Unquote(value)
	if err != nil {
		t.Fatalf("解析 %s 字符串: %v", path, err)
	}
	return result
}
