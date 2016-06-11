package v1

import (
	"github.com/codegangsta/cli"
	"github.com/gophercloud/cli/openstack/commands/blockstoragecommands/v1/snapshotcommands"
	"github.com/gophercloud/cli/openstack/commands/blockstoragecommands/v1/volumecommands"
)

// Get returns all the commands allowed for a `block-storage` request.
func Get() []cli.Command {
	return []cli.Command{
		{
			Name:        "snapshot",
			Usage:       "Copies of block storage volumes at a specific moment in time. Used for backup, restoration, and other long term storage.",
			Subcommands: snapshotcommands.Get(),
		},
		{
			Name:        "volume",
			Usage:       "Block level volumes to add storage capacity to your servers.",
			Subcommands: volumecommands.Get(),
		},
	}
}
