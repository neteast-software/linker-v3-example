package middleware

import (
	"context"
	"errors"

	"github.com/neteast-software/go-module/acl"
	http "github.com/neteast-software/go-module/http/gin/linker"
	ginmiddleware "github.com/neteast-software/go-module/http/gin/middleware"

	consoleadapter "linker-v3-example/internal/adapter/console"
	usermodel "linker-v3-example/internal/model/user"
	userservice "linker-v3-example/internal/service/user"
)

type consoleUserKey struct{}

// Console 完成后台业务 API 的最终认证与 ACL 判断。
func Console() http.HandlerFunc {
	authenticator := ginmiddleware.AuthenticatorFunc(func(c *http.Context) (string, error) {
		service, err := http.Require(c, userservice.AuthKey())
		if err != nil {
			return "", ginmiddleware.ErrAuthUnavailable
		}
		raw, err := Bearer(c)
		if err != nil {
			return "", ginmiddleware.ErrUnauthenticated
		}
		user, claims, err := service.Current(c.Request.Context(), raw, "console")
		if err != nil {
			if errors.Is(err, userservice.ErrUnavailable) {
				return "", ginmiddleware.ErrAuthUnavailable
			}
			return "", ginmiddleware.ErrUnauthenticated
		}
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), consoleUserKey{}, user))
		return claims.Subject, nil
	})
	provider := acl.ProviderFunc(func(ctx context.Context, subject string) (acl.Access, error) {
		user, ok := ctx.Value(consoleUserKey{}).(usermodel.User)
		if !ok {
			return acl.Access{}, ginmiddleware.ErrACLUnavailable
		}
		return consoleadapter.Access(user, subject), nil
	})
	return ginmiddleware.ACL(authenticator, provider)
}
