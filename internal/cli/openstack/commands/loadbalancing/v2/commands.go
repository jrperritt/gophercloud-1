package v2

import (
	"github.com/gophercloud/gophercloud/internal/cli/openstack/commands/loadbalancing/v2/loadbalancer"
	"gopkg.in/urfave/cli.v1"
)

// Get returns all the commands allowed for a `load-balancer` v2 request.
func Get() []cli.Command {
	return []cli.Command{
		{
			Name:        "load-balancer",
			Usage:       "LB operations",
			Subcommands: loadbalancer.Get(),
		},
	}
}
