package server

import (
	"github.com/gophercloud/cli/lib/traits"
	"gopkg.in/urfave/cli.v1"
)

var (
	CommandPrefix = "compute server"
)

type ServerV2Command struct {
	traits.Commandable
	traits.ComputeV2able
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
		getMetadata,
		deleteMetadata,
	}
}
