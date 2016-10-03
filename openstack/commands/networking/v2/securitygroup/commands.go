package securitygroup

import (
	"github.com/gophercloud/cli/openstack/commands"
	"gopkg.in/urfave/cli.v1"
)

var commandPrefix = "networking security-group"

type SecurityGroupV2Command struct {
	commands.Command
	name string
}

func (_ SecurityGroupV2Command) ServiceType() string {
	return "networking"
}

// Get returns all the commands allowed for a `networking security-group` v2 request.
func Get() []cli.Command {
	return []cli.Command{
		list,
		get,
		create,
		remove,
	}
}
