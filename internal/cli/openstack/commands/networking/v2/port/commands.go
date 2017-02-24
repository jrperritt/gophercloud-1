package port

import (
	"github.com/gophercloud/gophercloud/internal/cli/lib/traits"
	"gopkg.in/urfave/cli.v1"
)

var CommandPrefix = "networking port"

type PortV2Command struct {
	traits.Commandable
	traits.NetworkingV2able
	name string
}

// Get returns all the commands allowed for a `networking port` v2 request.
func Get() []cli.Command {
	return []cli.Command{
		list,
		get,
		update,
		create,
		remove,
	}
}
