package filescommands

import (
	"github.com/gophercloud/cli/src/commands/filescommands/accountcommands"
	"github.com/gophercloud/cli/src/commands/filescommands/containercommands"
	"github.com/gophercloud/cli/src/commands/filescommands/largeobjectcommands"
	"github.com/gophercloud/cli/src/commands/filescommands/objectcommands"
	"github.com/gophercloud/cli/vendor/github.com/codegangsta/cli"
)

// Get returns all the commands allowed for a `files` request.
func Get() []cli.Command {
	return []cli.Command{
		{
			Name:        "account",
			Usage:       "Storage for you account metadata.",
			Subcommands: accountcommands.Get(),
		},
		{
			Name:        "container",
			Usage:       "Storage compartments for your objects/files.",
			Subcommands: containercommands.Get(),
		},
		{
			Name:        "object",
			Usage:       "Data storage for objects/files/media.",
			Subcommands: objectcommands.Get(),
		},
		{
			Name:        "large-object",
			Usage:       "Data storage for large objects/files/media.",
			Subcommands: largeobjectcommands.Get(),
		},
	}
}
