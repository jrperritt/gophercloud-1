package containercommands

import "github.com/gophercloud/cli/vendor/github.com/codegangsta/cli"

var commandPrefix = "files container"
var serviceClientType = "object-store"

// Get returns all the commands allowed for a `files container` request.
func Get() []cli.Command {
	return []cli.Command{
		create,
		list,
		get,
		remove,
		update,
		empty,
		setMetadata,
		updateMetadata,
		getMetadata,
		deleteMetadata,
	}
}
