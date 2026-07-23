package example

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestScaffoldUsesCapabilityRoots(t *testing.T) {
	for _, name := range []string{
		"adapter",
		"client",
		"component",
		"config",
		"constant",
		"model",
		"page",
		"route",
		"rpc",
		"service",
	} {
		path := filepath.Join("../internal", name)
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Fatalf("不应恢复全局技术分层 %s: %v", path, err)
		}
	}
}

func TestScaffoldAdaptersKeepCapabilityPackageName(t *testing.T) {
	for _, adapter := range []string{"client", "http", "linker"} {
		for _, file := range adapterFiles(t, adapter) {
			parsed, err := parser.ParseFile(token.NewFileSet(), file, nil, parser.PackageClauseOnly)
			if err != nil {
				t.Fatalf("解析 %s: %v", file, err)
			}
			capability := filepath.Base(filepath.Dir(filepath.Dir(file)))
			if parsed.Name.Name != capability {
				t.Fatalf("%s 是 %s 的适配层，package 应为 %q，实际为 %q", file, capability, capability, parsed.Name.Name)
			}
		}
	}
}

func TestScaffoldLinkerAdaptersOwnCompileBoundary(t *testing.T) {
	componentFiles, err := filepath.Glob("../internal/*/linker/component.go")
	if err != nil {
		t.Fatalf("扫描 Linker 适配层: %v", err)
	}
	if len(componentFiles) == 0 {
		t.Fatal("没有找到 Linker 组件适配层")
	}
	for _, file := range componentFiles {
		for _, forbidden := range []string{"http.Context", "response.", "http.RegisterIn", "http.Routes("} {
			assertFileNotContains(t, file, forbidden)
		}

		capability := filepath.Base(filepath.Dir(filepath.Dir(file)))
		httpDir := filepath.Join("../internal", capability, "http")
		if _, err := os.Stat(httpDir); err == nil {
			want := "linker-v3-example/internal/" + capability + "/http"
			if !hasBlankImport(t, file, want) {
				t.Fatalf("%s 应 blank import %q，使 HTTP 声明服从组件编译边界", file, want)
			}
		}

		content := readFile(t, file)
		if strings.Contains(content, "type Component struct") {
			if !strings.Contains(content, "const ID linker.ID") ||
				!strings.Contains(content, "func (p *Component) Identity() linker.ID") {
				t.Fatalf("%s 的组件必须自治声明 ID 和 Identity", file)
			}
		}
	}
}

func TestScaffoldHTTPFilesSelfRegister(t *testing.T) {
	files := adapterFiles(t, "http")
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
		t.Fatal("没有找到自注册 HTTP route")
	}
}

func TestScaffoldCapabilityRootsStayBusinessFacing(t *testing.T) {
	entries, err := os.ReadDir("../internal")
	if err != nil {
		t.Fatalf("读取 internal: %v", err)
	}
	for _, entry := range entries {
		if !entry.IsDir() || entry.Name() == "app" {
			continue
		}
		files, err := filepath.Glob(filepath.Join("../internal", entry.Name(), "*.go"))
		if err != nil {
			t.Fatalf("扫描 %s: %v", entry.Name(), err)
		}
		for _, file := range files {
			content := readFile(t, file)
			for _, forbidden := range []string{"linker_v3_example_", "Model struct", "Entity struct", "DTO struct", "VO struct"} {
				if strings.Contains(content, forbidden) {
					t.Fatalf("能力根 package 应使用业务节点和相对语义，不应包含 %q: %s", forbidden, file)
				}
			}
			for _, imported := range localImports(t, file) {
				if imported == "linker-v3-example/internal/app" ||
					strings.Contains(imported, "/linker") ||
					strings.Contains(imported, "/http") {
					t.Fatalf("能力根 package 不应反向依赖装配层或适配层 %q: %s", imported, file)
				}
			}
		}
	}
}

func TestScaffoldCentralizesMiddlewareImplementation(t *testing.T) {
	for _, file := range adapterFiles(t, "http") {
		if strings.Contains(readFile(t, file), ".Next()") {
			t.Fatalf("HTTP 入口不应内嵌 middleware 实现，统一由 internal/access 承担: %s", file)
		}
	}
}

func TestScaffoldKeepsFrameworkAssemblySemantic(t *testing.T) {
	app := readFile(t, "../internal/app/app.go")
	for _, imported := range localImports(t, "../internal/app/app.go") {
		parts := strings.Split(strings.TrimPrefix(imported, "linker-v3-example/internal/"), "/")
		if len(parts) < 2 || parts[1] != "linker" && parts[1] != "client" {
			t.Fatalf("internal/app 只应装配能力适配层，不应直连业务实现 %q", imported)
		}
	}
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
	for _, file := range adapterFiles(t, "linker") {
		if strings.Contains(readFile(t, file), "NewCapabilityKey") {
			t.Fatalf("Linker 适配层应使用能力根 package 自己声明的能力入口: %s", file)
		}
	}
}

func TestScaffoldKeepsOneConfigurationPath(t *testing.T) {
	if _, err := os.Stat("../internal/config"); !os.IsNotExist(err) {
		t.Fatalf("不应建立全局 internal/config: %v", err)
	}
	app := readFile(t, "../internal/app/app.go")
	for _, forbidden := range []string{
		"postgresql.WithConfig",
		"rpc.WithConfig",
		"server.WithHTTP(",
		"internal/config",
	} {
		if strings.Contains(app, forbidden) {
			t.Fatalf("internal/app 不应包含 %q", forbidden)
		}
	}
	if !strings.Contains(app, "server.Config(sources...)") {
		t.Fatal("internal/app 应只转交一条有序 Source 配置链")
	}
	user := readFile(t, "../internal/user/linker/component.go")
	if !strings.Contains(user, "func (p *Component) Bootstrap(") || !strings.Contains(user, "func (p *Component) Configs()") {
		t.Fatal("业务组件应自治声明类型化配置引导和生命周期模式")
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

func localImports(t *testing.T, file string) []string {
	t.Helper()
	parsed, err := parser.ParseFile(token.NewFileSet(), file, nil, parser.ImportsOnly)
	if err != nil {
		t.Fatalf("解析 %s 的 imports: %v", file, err)
	}
	var imports []string
	for _, spec := range parsed.Imports {
		path, err := strconv.Unquote(spec.Path.Value)
		if err != nil {
			t.Fatalf("解析 import %s: %v", spec.Path.Value, err)
		}
		if strings.HasPrefix(path, "linker-v3-example/internal/") {
			imports = append(imports, path)
		}
	}
	return imports
}

func adapterFiles(t *testing.T, adapter string) []string {
	t.Helper()
	pattern := filepath.Join("../internal", "*", adapter, "*.go")
	files, err := filepath.Glob(pattern)
	if err != nil {
		t.Fatalf("扫描 %s 适配层: %v", adapter, err)
	}
	ret := files[:0]
	for _, file := range files {
		if !strings.HasSuffix(file, "_test.go") {
			ret = append(ret, file)
		}
	}
	return ret
}

func readFile(t *testing.T, file string) string {
	t.Helper()
	data, err := os.ReadFile(file)
	if err != nil {
		t.Fatalf("read %s: %v", file, err)
	}
	return string(data)
}
