package user

import "errors"

var ErrLogin = errors.New("账号或密码错误")
var ErrUnavailable = errors.New("用户服务尚未就绪")

var errDefaultUserDuplicate = errors.New("默认用户重复")
var errDefaultAccountDuplicate = errors.New("默认账号重复")
