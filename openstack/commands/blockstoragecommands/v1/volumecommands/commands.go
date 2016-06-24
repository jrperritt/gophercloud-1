package volumecommands

import "github.com/codegangsta/cli"

var commandPrefix = "block-storage volume"

//var serviceClientType = "block-storage"

type VolumeV1Command struct{}

func (_ VolumeV1Command) ServiceClientType() string {
	return "block-storage"
}

// Get returns all the commands allowed for a `block-storage volumes` request.
func Get() []cli.Command {
	return []cli.Command{
		list,
		get,
		update,
		create,
		remove,
	}
}
