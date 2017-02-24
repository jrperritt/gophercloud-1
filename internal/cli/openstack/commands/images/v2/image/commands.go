package image

import (
	"github.com/gophercloud/gophercloud/internal/cli/lib/traits"
	cli "gopkg.in/urfave/cli.v1"
)

var commandPrefix = "images image"

type ImageV2Command struct {
	traits.Commandable
	traits.ImagesV2able
}

// Get returns all the commands allowed for a `files container` v1 request.
func Get() []cli.Command {
	return []cli.Command{
		list,
		//get,
		//update,
		create,
		remove,
		upload,
		//download,
	}
}
