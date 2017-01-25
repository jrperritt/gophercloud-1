package securitygroup

import (
	"github.com/gophercloud/cli/lib/traits"
	"gopkg.in/urfave/cli.v1"
)

var commandPrefix = "networking security-group"

type SecurityGroupV2Command struct {
	traits.Commandable
	traits.NetworkingV2able
	name string
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
