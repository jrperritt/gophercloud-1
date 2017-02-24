package account

import (
	"github.com/gophercloud/gophercloud/internal/cli/lib/traits"
	"gopkg.in/urfave/cli.v1"
)

var commandPrefix = "files account"

type AccountV1Command struct {
	traits.Commandable
	traits.FilesV1able
	name string
}

// Get returns all the commands allowed for a `files account` v1 request.
func Get() []cli.Command {
	return []cli.Command{
		get,
		update,
	}
}
