package object

import (
	"github.com/gophercloud/cli/openstack"
	"gopkg.in/urfave/cli.v1"
)

var commandPrefix = "files object"

type ObjectV1Command struct {
	openstack.CommandUtil
	container string
	name      string
}

func (_ ObjectV1Command) ServiceClientType() string {
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
