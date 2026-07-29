package access

import (
	"net/http"
	"strconv"
)

const (
	HeaderUserID   = "X-Internal-User-ID"
	HeaderUsername = "X-Internal-Username"
	HeaderPlatform = "X-Internal-Platform"
	HeaderSource   = "X-Internal-Source"
	HeaderScope    = "X-Internal-Scope"
)

var identityHeaders = [...]string{
	HeaderUserID,
	HeaderUsername,
	HeaderPlatform,
	HeaderSource,
	HeaderScope,
	"Scope",
}

// Identity 是经过当前入口策略确认后才可投影到内部请求的身份。
type Identity struct {
	UserID   uint64
	Username string
	Platform string
	Source   string
	Scope    string
}

// ClearIdentity 清除调用方提交的内部身份声明。
func ClearIdentity(header http.Header) {
	for _, name := range identityHeaders {
		header.Del(name)
	}
}

// ProjectIdentity 只投影已经验证的内部身份字段。
func ProjectIdentity(header http.Header, identity Identity) {
	ClearIdentity(header)
	if identity.UserID > 0 {
		header.Set(HeaderUserID, strconv.FormatUint(identity.UserID, 10))
	}
	if identity.Username != "" {
		header.Set(HeaderUsername, identity.Username)
	}
	if identity.Platform != "" {
		header.Set(HeaderPlatform, identity.Platform)
	}
	if identity.Source != "" {
		header.Set(HeaderSource, identity.Source)
	}
	if identity.Scope != "" {
		header.Set(HeaderScope, identity.Scope)
	}
}
