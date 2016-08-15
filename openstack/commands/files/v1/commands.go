package v1

import (
	"gopkg.in/urfave/cli.v1"
	"github.com/gophercloud/cli/openstack/commands/files/v1/container"
)

// Get returns all the commands allowed for a `files` v2 request.
func Get() []cli.Command {
	return []cli.Command{
		{
			Name:        "container",
			Usage:       "Object storage containers",
			Subcommands: container.Get(),
		},
	}
}
