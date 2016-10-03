package port

import (
	"github.com/gophercloud/cli/openstack/commands"
	"gopkg.in/urfave/cli.v1"
)

var CommandPrefix = "networking port"

type PortV2Command struct {
	commands.Command
	name string
}

func (_ PortV2Command) ServiceType() string {
	return "networking"
}

// Get returns all the commands allowed for a `networking port` v2 request.
func Get() []cli.Command {
	return []cli.Command{
		list,
		get,
		update,
		create,
		remove,
	}
}
