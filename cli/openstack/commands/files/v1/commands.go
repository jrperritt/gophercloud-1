package v1

import (
	"github.com/gophercloud/cli/openstack/commands/files/v1/container"
	"github.com/gophercloud/cli/openstack/commands/files/v1/object"
	"gopkg.in/urfave/cli.v1"
)

// Get returns all the commands allowed for a `files` v1 request.
func Get() []cli.Command {
	return []cli.Command{
		/*
			{
				Name:        "account",
				Usage:       "Object storage account information",
				Subcommands: account.Get(),
			},
		*/
		{
			Name:        "container",
			Usage:       "Object storage containers",
			Subcommands: container.Get(),
		},
		{
			Name:        "object",
			Usage:       "Object storage files",
			Subcommands: object.Get(),
		},
	}
}
