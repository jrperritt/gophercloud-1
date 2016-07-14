package v1

import (
	"github.com/codegangsta/cli"
	"github.com/gophercloud/cli/openstack/commands/blockstorage/v1/volume"
)

// Get returns all the commands allowed for a `block-storage` v1 request.
func Get() []cli.Command {
	return []cli.Command{
		{
			Name:        "volume",
			Usage:       "Block level volumes to add storage capacity to your servers.",
			Subcommands: volume.Get(),
		},
	}
}
