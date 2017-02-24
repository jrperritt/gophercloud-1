package service

import (
	"github.com/gophercloud/gophercloud/internal/cli/lib/traits"
	cli "gopkg.in/urfave/cli.v1"
)

var commandPrefix = "identity service"

type ServiceV3Command struct {
	traits.Commandable
	traits.IdentityV3able
}

// Get returns all the commands allowed for an `idenity service` v3 request.
func Get() []cli.Command {
	return []cli.Command{
		list,
	}
}
