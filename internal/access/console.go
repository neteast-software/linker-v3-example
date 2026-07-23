package access

import (
	"context"
	"errors"

	"github.com/neteast-software/go-module/acl"
	http "github.com/neteast-software/go-module/http/gin/linker"
	ginmiddleware "github.com/neteast-software/go-module/http/gin/middleware"

	console "linker-v3-example/internal/console"
	user "linker-v3-example/internal/user"
)

type consoleUserKey struct{}

// Console 完成后台业务 API 的最终认证与 ACL 判断。
func Console() http.HandlerFunc {
	authenticator := ginmiddleware.AuthenticatorFunc(func(c *http.Context) (string, error) {
		service, err := http.Require(c, user.AuthKey())
		if err != nil {
			return "", ginmiddleware.ErrAuthUnavailable
		}
		raw, err := Bearer(c)
		if err != nil {
			return "", ginmiddleware.ErrUnauthenticated
		}
		current, claims, err := service.Current(c.Request.Context(), raw, "console")
		if err != nil {
			if errors.Is(err, user.ErrUnavailable) {
				return "", ginmiddleware.ErrAuthUnavailable
			}
			return "", ginmiddleware.ErrUnauthenticated
		}
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), consoleUserKey{}, current))
		return claims.Subject, nil
	})
	provider := acl.ProviderFunc(func(ctx context.Context, subject string) (acl.Access, error) {
		current, ok := ctx.Value(consoleUserKey{}).(user.User)
		if !ok {
			return acl.Access{}, ginmiddleware.ErrACLUnavailable
		}
		return console.Access(current, subject), nil
	})
	return ginmiddleware.ACL(authenticator, provider)
}
