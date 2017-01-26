package volume

import (
	"github.com/gophercloud/gophercloud/cli/lib/traits"
	"gopkg.in/urfave/cli.v1"
)

var commandPrefix = "block-storage volume"

type command struct {
	traits.Commandable
	traits.BlockStorageV2able
}

// Get returns all the commands allowed for a `block-storage volume` v2 request.
func Get() []cli.Command {
	return []cli.Command{
	//create,
	//upload,
	//list,
	//get,
	//remove,
	}
}
