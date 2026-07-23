package console

import (
	"context"
	"strconv"

	"github.com/neteast-software/go-module/acl"

	user "linker-v3-example/internal/user"
)

func (p *Provider) Access(ctx context.Context, subject string) (acl.Access, error) {
	id, err := strconv.ParseUint(subject, 10, 64)
	if err != nil {
		return acl.Access{}, err
	}
	user, err := p.user.ProfileByID(ctx, id)
	if err != nil {
		return acl.Access{}, err
	}
	return Access(user, subject), nil
}

// Access 把业务 user 编译为 ACL 热路径使用的访问能力投影。
func Access(current user.User, subject string) acl.Access {
	access := acl.Access{
		ID:      subject,
		Level:   1,
		Actions: acl.Read,
		Groups: []acl.ResourceGroup{
			acl.Role("console-reader",
				acl.NewResource("console.*"),
				acl.NewResource("graph.console.*"),
			),
		},
	}
	if current.Role == "admin" {
		access.Level = 9
		access.Actions = acl.AllActions
		access.Groups = []acl.ResourceGroup{
			acl.Role("console-admin", acl.NewResource("*")),
		}
	}
	return access
}
