package access

import (
	"net/http"
	"testing"
)

func TestIdentityProjectionClearsUntrustedValues(t *testing.T) {
	header := http.Header{
		HeaderUserID:   {"999"},
		HeaderUsername: {"伪造用户"},
		"Scope":        {"external"},
	}
	ProjectIdentity(header, Identity{
		UserID:   7,
		Username: "verified",
		Platform: "console",
		Source:   "session",
		Scope:    "equipment.read",
	})
	if header.Get(HeaderUserID) != "7" || header.Get(HeaderUsername) != "verified" {
		t.Fatalf("identity = %#v", header)
	}
	if header.Get("Scope") != "" {
		t.Fatalf("legacy Scope was retained: %#v", header)
	}
}
