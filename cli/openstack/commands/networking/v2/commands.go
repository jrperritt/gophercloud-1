package v2

import (
	"github.com/gophercloud/gophercloud/cli/openstack/commands/networking/v2/network"
	"github.com/gophercloud/gophercloud/cli/openstack/commands/networking/v2/port"
	"github.com/gophercloud/gophercloud/cli/openstack/commands/networking/v2/securitygroup"
	"github.com/gophercloud/gophercloud/cli/openstack/commands/networking/v2/securitygrouprule"
	"github.com/gophercloud/gophercloud/cli/openstack/commands/networking/v2/subnet"
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
		{
			Name:        "subnet",
			Usage:       "Logically-partitioned address spaces",
			Subcommands: subnet.Get(),
		},
		{
			Name:        "port",
			Usage:       "Logically-defined ports",
			Subcommands: port.Get(),
		},
		{
			Name:        "security-group",
			Usage:       "Permissioned network access via groups",
			Subcommands: securitygroup.Get(),
		},
		{
			Name:        "security-group-rules",
			Usage:       "Rules that define security groups",
			Subcommands: securitygrouprule.Get(),
		},
	}
}
