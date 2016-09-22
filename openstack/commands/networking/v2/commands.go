package v2

import (
	"github.com/gophercloud/cli/openstack/commands/networking/v2/network"
	"gopkg.in/urfave/cli.v1"
)

// Get returns all the commands allowed for a `networking` v2 request.
func Get() []cli.Command {
	return []cli.Command{
		{
			Name:        "network",
			Usage:       "Software-defined networks",
			Subcommands: network.Get(),
		},
	}
}
