package container

import (
	"github.com/gophercloud/cli/openstack/commands"
	"gopkg.in/urfave/cli.v1"
)

var commandPrefix = "files container"

type ContainerV1Command struct {
	commands.Command
	name string
}

func (_ ContainerV1Command) ServiceType() string {
	return "files"
}

// Get returns all the commands allowed for a `files container` v1 request.
func Get() []cli.Command {
	return []cli.Command{
		list,
		get,
		update,
		create,
		remove,
		empty,
	}
}
