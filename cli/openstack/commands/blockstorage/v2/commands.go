package v2

import (
	"github.com/gophercloud/gophercloud/cli/openstack/commands/blockstorage/v2/volume"
	"gopkg.in/urfave/cli.v1"
)

// Get returns all the commands allowed for a `block-storage` v2 request.
func Get() []cli.Command {
	return []cli.Command{
		{
			Name:        "volume",
			Usage:       "Logically-partitioned storage",
			Subcommands: volume.Get(),
		},
	}
}
