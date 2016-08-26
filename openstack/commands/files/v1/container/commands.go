package container

import (
	"github.com/gophercloud/cli/openstack"
	"gopkg.in/urfave/cli.v1"
)

var commandPrefix = "files container"

type ContainerV1Command struct {
	openstack.CommandUtil
	name string
}

func (_ ContainerV1Command) ServiceClientType() string {
	return "files"
}

// Get returns all the commands allowed for a `servers instance` v2 request.
func Get() []cli.Command {
	return []cli.Command{
		list,
		get,
		update,
		create,
		remove,
	}
}
