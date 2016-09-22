package instance

import (
	"github.com/gophercloud/cli/openstack/commands"
	"gopkg.in/urfave/cli.v1"
)

var (
	commandPrefix = "compute server"
)

type ServerV2Command struct {
	commands.Command
}

func (c *ServerV2Command) ServiceType() string {
	return "compute"
}

// Get returns all the commands allowed for a `compute server` v2 request.
func Get() []cli.Command {
	return []cli.Command{
		list,
		get,
		update,
		create,
		remove,
		resize,
		rebuild,
		reboot,
		deleteMetadata,
	}
}
