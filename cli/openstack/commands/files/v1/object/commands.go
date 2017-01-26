package object

import (
	"github.com/gophercloud/cli/lib/traits"
	"gopkg.in/urfave/cli.v1"
)

var commandPrefix = "files object"

type ObjectV1Command struct {
	traits.Commandable
	traits.ComputeV2able
	container string
	name      string
}

// Get returns all the commands allowed for a `files object` v1 request.
func Get() []cli.Command {
	return []cli.Command{
		list,
		get,
		upload,
		download,
		remove,
	}
}
