package vendorauth

import "errors"

var (
	ErrToken       = errors.New("开放平台 token 不合法")
	ErrScope       = errors.New("开放平台 token 缺少访问范围")
	ErrCertificate = errors.New("开放平台 token 与客户端证书不匹配")
)
