package v3

import (
	"github.com/gophercloud/gophercloud/internal/cli/openstack/commands/identity/v3/endpoint"
	"github.com/gophercloud/gophercloud/internal/cli/openstack/commands/identity/v3/service"
	cli "gopkg.in/urfave/cli.v1"
)

// Get returns all the commands allowed for an `identity` v3 request.
func Get() []cli.Command {
	return []cli.Command{
		{
			Name:        "service",
			Usage:       "Available OpenStack services",
			Subcommands: service.Get(),
		},
		{
			Name:        "endpoint",
			Usage:       "Available service endpoints",
			Subcommands: endpoint.Get(),
		},
	}
}
