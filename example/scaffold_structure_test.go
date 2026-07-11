package example

import (
	"go/parser"
	"go/token"
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

func TestScaffoldRouteAPIFilesSelfRegister(t *testing.T) {
	files, err := filepath.Glob("../internal/route/*/*_api.go")
	if err != nil {
		t.Fatalf("scan route api files: %v", err)
	}
	if len(files) == 0 {
		t.Fatalf("missing route api files")
	}
	for _, file := range files {
		data := readFile(t, file)
		if !strings.Contains(data, ".RegisterIn(") && !strings.Contains(data, ".Register(") {
			t.Fatalf("%s should declare its route with RegisterIn or Register", file)
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
