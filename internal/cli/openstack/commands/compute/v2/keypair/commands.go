package keypair

import (
	"github.com/gophercloud/gophercloud/internal/cli/lib/traits"
	"gopkg.in/urfave/cli.v1"
)

var commandPrefix = "compute keypair"

type KeypairV2Command struct {
	traits.Commandable
	traits.ComputeV2able
}

// Get returns all the commands allowed for a `compute keypair` v2 request.
func Get() []cli.Command {
	return []cli.Command{
		generate,
		upload,
		list,
		get,
		remove,
	}
}
