package console

import (
	"context"
	"strconv"

	"github.com/neteast-software/go-module/acl"

	usermodel "linker-v3-example/internal/model/user"
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
func Access(user usermodel.User, subject string) acl.Access {
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
	if user.Role == "admin" {
		access.Level = 9
		access.Actions = acl.AllActions
		access.Groups = []acl.ResourceGroup{
			acl.Role("console-admin", acl.NewResource("*")),
		}
	}
	return access
}
