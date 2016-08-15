package instance

import "gopkg.in/urfave/cli.v1"

var commandPrefix = "servers instance"

type InstanceV2Command struct{}

func (_ InstanceV2Command) ServiceClientType() string {
	return "compute"
}

// Get returns all the commands allowed for a `servers instance` v2 request.
func Get() []cli.Command {
	return []cli.Command{
		list,
		get,
		update,
		create,
		remove,
		resize,
		rebuild,
		reboot,
		deleteMetadata,
	}
}
