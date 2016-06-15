package volumecommands

import (
	"github.com/codegangsta/cli"
	"github.com/gophercloud/cli/openstack"
)

var commandPrefix = "block-storage volume"
var serviceClientType = "block-storage"

// Get returns all the commands allowed for a `block-storage volumes` request.
func Get() []cli.Command {
	return []cli.Command{
		openstack.NewCommand(new(commandList), flagsList, serviceClientType),
		openstack.NewCommand(new(commandGet), flagsGet, serviceClientType),
		openstack.NewCommand(new(commandUpdate), flagsUpdate, serviceClientType),
		openstack.NewCommand(new(commandCreate), flagsCreate, serviceClientType),
		openstack.NewCommand(new(commandDelete), flagsDelete, serviceClientType),
	}
}
