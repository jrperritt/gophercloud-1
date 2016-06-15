package v1

import (
	"github.com/codegangsta/cli"
	"github.com/gophercloud/cli/openstack/commands/blockstoragecommands/v1/volumecommands"
)

// Get returns all the commands allowed for a `block-storage` request.
func Get() []cli.Command {
	return []cli.Command{
		{
			Name:        "volume",
			Usage:       "Block level volumes to add storage capacity to your servers.",
			Subcommands: volumecommands.Get(),
		},
	}
}
