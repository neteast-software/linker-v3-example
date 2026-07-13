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

func TestScaffoldStructureKeepsRouteOwnershipLocal(t *testing.T) {
	assertFileNotContains(t, "../internal/app/app.go", "internal/route/")

	componentFiles, err := filepath.Glob("../internal/component/*/component.go")
	if err != nil {
		t.Fatalf("scan component files: %v", err)
	}
	if len(componentFiles) == 0 {
		t.Fatalf("missing component files")
	}
	for _, file := range componentFiles {
		assertFileNotContains(t, file, "http.Context")
		assertFileNotContains(t, file, "response.")
		assertFileNotContains(t, file, "http.RegisterIn")
		assertFileNotContains(t, file, "http.Routes(")

		domain := filepath.Base(filepath.Dir(file))
		routeDir := filepath.Join("../internal/route", domain)
		if _, err := os.Stat(routeDir); err == nil {
			want := "linker-v3-example/internal/route/" + domain
			if !hasBlankImport(t, file, want) {
				t.Fatalf("%s should blank import %q so route declarations follow the component compile boundary", file, want)
			}
		}
	}
}

func TestScaffoldRouteFilesSelfRegister(t *testing.T) {
	files := goFilesUnder(t, "../internal/route")
	routeFiles := 0
	for _, file := range files {
		parsed, err := parser.ParseFile(token.NewFileSet(), file, nil, 0)
		if err != nil {
			t.Fatalf("parse %s: %v", file, err)
		}
		inits := 0
		methods := 0
		registers := 0
		for _, decl := range parsed.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if ok && fn.Recv == nil && fn.Name.Name == "init" {
				inits++
			}
			if ok && strings.HasSuffix(fn.Name.Name, "API") {
				t.Fatalf("route 已表达 API 归属，函数名不应重复 API 后缀: %s:%s", file, fn.Name.Name)
			}
		}
		ast.Inspect(parsed, func(node ast.Node) bool {
			call, ok := node.(*ast.CallExpr)
			if !ok {
				return true
			}
			selector, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			switch selector.Sel.Name {
			case "GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD", "Any", "Match":
				methods++
			case "Register", "RegisterIn":
				registers++
			}
			return true
		})
		if methods == 0 {
			continue
		}
		routeFiles++
		if strings.HasSuffix(file, "_api.go") {
			t.Fatalf("route 路径已表达 API 归属，文件名不应重复 _api 后缀: %s", file)
		}
		if inits != 1 || registers != 1 || methods != 1 {
			t.Fatalf("%s 应保持一个 init、一个注册入口和一个 method route: init=%d register=%d method=%d", file, inits, registers, methods)
		}
	}
	if routeFiles == 0 {
		t.Fatal("missing self-registering route files")
	}
}

func TestScaffoldModelNamesStayBusinessFacing(t *testing.T) {
	for _, file := range goFilesUnder(t, "../internal/model") {
		content := readFile(t, file)
		for _, forbidden := range []string{"linker_v3_example_", "Model struct", "Entity struct", "DTO struct", "VO struct"} {
			if strings.Contains(content, forbidden) {
				t.Fatalf("model 应使用业务节点和相对语义，不应包含 %q: %s", forbidden, file)
			}
		}
	}
}

func TestScaffoldCentralizesMiddlewareImplementation(t *testing.T) {
	files := goFilesUnder(t, "../internal/route")
	for _, file := range files {
		if filepath.Base(filepath.Dir(file)) == "middleware" {
			continue
		}
		if strings.Contains(readFile(t, file), ".Next()") {
			t.Fatalf("middleware 实现必须集中在 internal/route/middleware: %s", file)
		}
	}
}

func TestScaffoldKeepsFrameworkAssemblySemantic(t *testing.T) {
	app := readFile(t, "../internal/app/app.go")
	for _, forbidden := range []string{
		"WithLifecycleObserver",
		"WithConfigObserver",
		"WithServerOptions",
		"ChainUnaryInterceptor",
		"UnaryServerMeta",
		"metricserver.Observe",
		"metricgrpc.",
		"tracegrpc.",
		"google.golang.org/grpc",
		"WithTracing()",
		"WithMetrics()",
	} {
		if strings.Contains(app, forbidden) {
			t.Fatalf("internal/app 不应手工装配底层入口 %q", forbidden)
		}
	}
	if !strings.Contains(app, "server.WithMetrics(prometheus.New())") {
		t.Fatal("internal/app 应通过 framework 语义入口装配观测能力")
	}
	for _, file := range goFilesUnder(t, "../internal/component") {
		if strings.Contains(readFile(t, file), "NewCapabilityKey") {
			t.Fatalf("业务 component 应使用能力归属 package 的语义入口: %s", file)
		}
	}
}

func TestScaffoldKeepsOneConfigurationPath(t *testing.T) {
	if _, err := os.Stat("../internal/config"); !os.IsNotExist(err) {
		t.Fatalf("global internal/config should not exist: %v", err)
	}
	app := readFile(t, "../internal/app/app.go")
	for _, forbidden := range []string{
		"postgresql.WithConfig",
		"rpc.WithConfig",
		"server.WithHTTP(",
		"internal/config",
	} {
		if strings.Contains(app, forbidden) {
			t.Fatalf("internal/app should not contain %q", forbidden)
		}
	}
	if !strings.Contains(app, "server.Config(sources...)") {
		t.Fatal("internal/app should forward the single ordered Source chain")
	}
	user := readFile(t, "../internal/component/user/component.go")
	if !strings.Contains(user, "func (p *Component) Bootstrap(") || !strings.Contains(user, "func (p *Component) Configs()") {
		t.Fatal("business component should own typed config bootstrap and lifecycle mode")
	}
}

func assertFileNotContains(t *testing.T, file string, needle string) {
	t.Helper()
	data := readFile(t, file)
	if strings.Contains(data, needle) {
		t.Fatalf("%s should not contain %q", file, needle)
	}
}

func hasBlankImport(t *testing.T, file string, importPath string) bool {
	t.Helper()
	parsed, err := parser.ParseFile(token.NewFileSet(), file, nil, parser.ImportsOnly)
	if err != nil {
		t.Fatalf("parse imports from %s: %v", file, err)
	}
	for _, spec := range parsed.Imports {
		if spec.Name == nil || spec.Name.Name != "_" {
			continue
		}
		path, err := strconv.Unquote(spec.Path.Value)
		if err != nil {
			t.Fatalf("unquote import path %s: %v", spec.Path.Value, err)
		}
		if path == importPath {
			return true
		}
	}
	return false
}

func readFile(t *testing.T, file string) string {
	t.Helper()
	data, err := os.ReadFile(file)
	if err != nil {
		t.Fatalf("read %s: %v", file, err)
	}
	return string(data)
}

func goFilesUnder(t *testing.T, root string) []string {
	t.Helper()
	var files []string
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !entry.IsDir() && strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("scan %s: %v", root, err)
	}
	return files
}
