package securitygrouprule

import (
	"github.com/gophercloud/gophercloud/internal/cli/lib/traits"
	"gopkg.in/urfave/cli.v1"
)

var commandPrefix = "networking security-group-rule"

type SecurityGroupRuleV2Command struct {
	traits.Commandable
	traits.NetworkingV2able
}

// Get returns all the commands allowed for a `networking security-group-rule` v2 request.
func Get() []cli.Command {
	return []cli.Command{
		list,
		get,
		create,
		remove,
	}
}
