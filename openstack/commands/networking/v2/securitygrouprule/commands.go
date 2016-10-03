package securitygrouprule

import (
	"github.com/gophercloud/cli/openstack/commands"
	"gopkg.in/urfave/cli.v1"
)

var commandPrefix = "networking security-group-rule"

type SecurityGroupRuleV2Command struct {
	commands.Command
}

func (_ SecurityGroupRuleV2Command) ServiceType() string {
	return "networking"
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
