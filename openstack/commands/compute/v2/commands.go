package v2

import (
	"gopkg.in/urfave/cli.v1"
	"github.com/gophercloud/cli/openstack/commands/compute/v2/flavor"
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
		{
			Name:        "flavor",
			Usage:       "Server flavors (options for operating system and architecture)",
			Subcommands: flavor.Get(),
		},
	}
}
