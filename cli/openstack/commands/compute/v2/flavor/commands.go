package flavor

import (
	"github.com/gophercloud/gophercloud/cli/lib/traits"
	"gopkg.in/urfave/cli.v1"
)

var commandPrefix = "compute flavor"

type FlavorCommand struct {
	traits.Commandable
	traits.ComputeV2able
}

// Get returns all the commands allowed for a `compute flavor` v2 request.
func Get() []cli.Command {
	return []cli.Command{
		list,
		get,
	}
}
