package serverscommands

import (
	"github.com/gophercloud/cli/commands/serverscommands/flavorcommands"
	"github.com/gophercloud/cli/commands/serverscommands/imagecommands"
	"github.com/gophercloud/cli/commands/serverscommands/instancecommands"
	"github.com/gophercloud/cli/commands/serverscommands/keypaircommands"
	"github.com/gophercloud/cli/commands/serverscommands/volumeattachmentcommands"
	"github.com/gophercloud/cli/vendor/github.com/codegangsta/cli"
)

// Get returns all the commands allowed for a `servers` request.
func Get() []cli.Command {
	return []cli.Command{
		{
			Name:        "instance",
			Usage:       "Virtual and bare metal servers.",
			Subcommands: instancecommands.Get(),
		},
		{
			Name:        "image",
			Usage:       "Base operating system layout for a server.",
			Subcommands: imagecommands.Get(),
		},
		{
			Name:        "flavor",
			Usage:       "Resource allocations for servers.",
			Subcommands: flavorcommands.Get(),
		},
		{
			Name:        "keypair",
			Usage:       "SSH keypairs for accessing servers.",
			Subcommands: keypaircommands.Get(),
		},
		{
			Name:        "volume-attachment",
			Usage:       "Volumes attached to servers.",
			Subcommands: volumeattachmentcommands.Get(),
		},
	}
}
