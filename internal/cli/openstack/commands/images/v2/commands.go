package v2

import (
	"github.com/gophercloud/gophercloud/internal/cli/openstack/commands/images/v2/image"
	cli "gopkg.in/urfave/cli.v1"
)

// Get returns all the commands allowed for a `files` v1 request.
func Get() []cli.Command {
	return []cli.Command{
		{
			Name:        "image",
			Usage:       "Manage virtual machine images",
			Subcommands: image.Get(),
		},
	}
}
