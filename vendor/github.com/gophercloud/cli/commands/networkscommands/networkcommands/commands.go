package networkcommands

import "github.com/gophercloud/cli/vendor/github.com/codegangsta/cli"

var commandPrefix = "networks network"
var serviceClientType = "network"

// Get returns all the commands allowed for a `networks network` request.
func Get() []cli.Command {
	return []cli.Command{
		create,
		get,
		remove,
		list,
		update,
	}
}
