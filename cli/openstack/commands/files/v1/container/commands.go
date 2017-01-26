package container

import (
	"github.com/gophercloud/cli/lib/traits"
	"gopkg.in/urfave/cli.v1"
)

var commandPrefix = "files container"

type ContainerV1Command struct {
	traits.Commandable
	traits.FilesV1able
	name string
}

// Get returns all the commands allowed for a `files container` v1 request.
func Get() []cli.Command {
	return []cli.Command{
		list,
		get,
		update,
		create,
		remove,
		empty,
	}
}
