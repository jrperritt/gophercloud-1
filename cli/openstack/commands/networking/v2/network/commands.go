package network

import (
	"github.com/gophercloud/gophercloud/cli/lib/traits"
	"gopkg.in/urfave/cli.v1"
)

var commandPrefix = "networking newtork"

type NetworkV2Command struct {
	traits.Commandable
	traits.NetworkingV2able
	name string
}

// Get returns all the commands allowed for a `networking network` v2 request.
func Get() []cli.Command {
	return []cli.Command{
		list,
		get,
		update,
		create,
		remove,
	}
}
