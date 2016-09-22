package flavor

import (
	"github.com/gophercloud/cli/openstack/commands"
	"gopkg.in/urfave/cli.v1"
)

var commandPrefix = "compute flavor"

type FlavorV2Command struct {
	commands.Command
}

func (_ FlavorV2Command) ServiceType() string {
	return "compute"
}

// Get returns all the commands allowed for a `compute flavor` v2 request.
func Get() []cli.Command {
	return []cli.Command{
		list,
		get,
	}
}
