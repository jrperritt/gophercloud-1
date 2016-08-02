package v2

import (
	"github.com/codegangsta/cli"
	"github.com/gophercloud/cli/openstack/commands/compute/v2/instance"
)

// Get returns all the commands allowed for a `servers` v2 request.
func Get() []cli.Command {
	return []cli.Command{
		{
			Name:        "instance",
			Usage:       "Virtual and bare metal servers.",
			Subcommands: instance.Get(),
		},
	}
}
