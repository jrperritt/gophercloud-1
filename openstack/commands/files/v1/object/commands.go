package object

import (
	"github.com/gophercloud/cli/openstack/commands"
	"gopkg.in/urfave/cli.v1"
)

var commandPrefix = "files object"

type ObjectV1Command struct {
	commands.Command
	container string
	name      string
}

func (_ ObjectV1Command) ServiceType() string {
	return "files"
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
