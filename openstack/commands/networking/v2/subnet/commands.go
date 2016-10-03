package subnet

import (
	"github.com/gophercloud/cli/openstack/commands"
	"gopkg.in/urfave/cli.v1"
)

var commandPrefix = "networking subnet"

type SubnetV2Command struct {
	commands.Command
	name string
}

func (_ SubnetV2Command) ServiceType() string {
	return "networking"
}

// Get returns all the commands allowed for a `networking subnet` v2 request.
func Get() []cli.Command {
	return []cli.Command{
		list,
		get,
		update,
		create,
		remove,
	}
}
