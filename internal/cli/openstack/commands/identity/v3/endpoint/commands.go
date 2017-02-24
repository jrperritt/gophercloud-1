package endpoint

import (
	"github.com/gophercloud/gophercloud/internal/cli/lib/traits"
	cli "gopkg.in/urfave/cli.v1"
)

var commandPrefix = "identity endpoint"

type EndpointV3Command struct {
	traits.Commandable
	traits.IdentityV3able
}

// Get returns all the commands allowed for an `idenity endpoint` v3 request.
func Get() []cli.Command {
	return []cli.Command{
		list,
	}
}
