package v2

import (
	"github.com/gophercloud/cli/openstack/commands/compute/v2/flavor"
	"github.com/gophercloud/cli/openstack/commands/compute/v2/keypair"
	"github.com/gophercloud/cli/openstack/commands/compute/v2/server"
	"gopkg.in/urfave/cli.v1"
)

// Get returns all the commands allowed for a `compute` v2 request.
func Get() []cli.Command {
	return []cli.Command{
		{
			Name:        "server",
			Usage:       "Virtual and bare metal servers.",
			Subcommands: server.Get(),
		},
		{
			Name:        "flavor",
			Usage:       "Server flavors (options for operating system and architecture)",
			Subcommands: flavor.Get(),
		},
		{
			Name:        "keypair",
			Usage:       "SSH keypairs for accessing servers.",
			Subcommands: keypair.Get(),
		},
	}
}
