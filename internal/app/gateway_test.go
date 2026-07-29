package app

import (
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/neteast-software/go-module/http/gateway/declaration"

	captcha "linker-v3-example/internal/captcha"
	session "linker-v3-example/internal/session"
	vendorauth "linker-v3-example/internal/vendorauth"
)

func TestLocalGatewayCodeAndYAMLDeclarationsStayEquivalent(t *testing.T) {
	content, err := os.ReadFile("../../config/gateway.routes.yaml")
	if err != nil {
		t.Fatal(err)
	}
	fromFile, err := declaration.Decode(content)
	if err != nil {
		t.Fatal(err)
	}
	fromCode := localDocument()
	if !reflect.DeepEqual(fromCode, fromFile) {
		t.Fatalf("Gateway code/YAML declarations drifted:\ncode=%#v\nfile=%#v", fromCode, fromFile)
	}
	if strings.Contains(strings.ToLower(string(content)), "xss") {
		t.Fatal("XSS body rewriting must not be a default Gateway policy")
	}
	if _, err = declaration.Compile(fromCode, declaration.WithFilters(
		vendorauth.Filter(nil, vendorauth.Scope("equipment-read", "equipment.read", true)),
		session.Filter(nil, session.Memory()),
		session.Upload(nil, session.Memory()),
		captcha.Filter(captcha.Memory(), ""),
	)); err != nil {
		t.Fatalf("compile local document: %v", err)
	}
}

func TestNacosGatewayDocumentUsesNeutralServiceRoute(t *testing.T) {
	plan, err := declaration.Compile(nacosDocument())
	if err != nil {
		t.Fatal(err)
	}
	routes := plan.Routes()
	if len(routes) != 1 {
		t.Fatalf("routes = %d", len(routes))
	}
	service, ok := routes[0].Service()
	if !ok || service != "example-http" {
		t.Fatalf("service = %q ok=%v", service, ok)
	}
}
