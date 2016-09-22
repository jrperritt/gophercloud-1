package network

import (
	"github.com/gophercloud/cli/openstack/commands"
	"gopkg.in/urfave/cli.v1"
)

var commandPrefix = "networking newtork"

type NetworkV2Command struct {
	commands.Command
	name string
}

func (_ NetworkV2Command) ServiceType() string {
	return "networking"
}

// Get returns all the commands allowed for a `networking newtork` v2 request.
func Get() []cli.Command {
	return []cli.Command{
		list,
		//get,
		//update,
		create,
		//remove,
	}
}
