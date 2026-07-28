package user

import "context"

// Actor 只负责把认证信息解析为当前业务操作者。
type Actor interface {
	ActorID(context.Context, string, string) (uint64, error)
}

// ActorID 返回认证信息对应的内部用户标识。
func (p *Service) ActorID(ctx context.Context, raw string, scope string) (uint64, error) {
	current, _, err := p.Current(ctx, raw, scope)
	if err != nil {
		return 0, err
	}
	return current.ID, nil
}
