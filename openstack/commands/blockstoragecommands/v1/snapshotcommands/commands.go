package snapshotcommands

import (
	"github.com/codegangsta/cli"
	"github.com/gophercloud/cli/openstack"
)

var commandPrefix = "block-storage snapshots"
var serviceClientType = "blockstorage"

// Get returns all the commands allowed for a `block-storage snapshots` request.
func Get() []cli.Command {
	return []cli.Command{
		//list(),
		openstack.NewCommand(new(commandGet), flagsGet, serviceClientType),
		openstack.NewCommand(new(commandCreate), flagsCreate, serviceClientType),
		openstack.NewCommand(new(commandDelete), flagsDelete, serviceClientType),
	}
}
