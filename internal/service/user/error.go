package user

import "errors"

var ErrLogin = errors.New("账号或密码错误")
var ErrUnavailable = errors.New("用户服务尚未就绪")
