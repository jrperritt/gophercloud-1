package main

import (
	"fmt"

	"github.com/gophercloud/cli/openstack/commands/compute"
	"github.com/gophercloud/cli/openstack/commands/files"
	"github.com/gophercloud/cli/openstack/commands/loadbalancing"
	"github.com/gophercloud/cli/openstack/commands/networking"
	"github.com/gophercloud/cli/setup"
	"gopkg.in/urfave/cli.v1"
)

var commands = []cli.Command{
	{
		Name:   "configure",
		Usage:  "Interactively create a config file for authentication.",
		Action: configure,
	},
	{
		Name: "init",
		Usage: "Enable tab for command completion." +
			"\n\tFor Linux and OS X, creates the `rack` man page and sets up" +
			"\n\tcommand completion for the Bash shell. Run `man ./rack.1` to" +
			"\n\tview the generated man page." +
			"\n\tFor Windows, creates a `posh_autocomplete.ps1` file in the" +
			"\n\t`$HOME/.rack` directory. You must run the file to set up" +
			"\n\tcommand completion.",
		Action: func(c *cli.Context) error {
			setup.Init(c)
			//man()
			return nil
		},
	},
	{
		Name:  "version",
		Usage: "Print the version of this binary.",
		Action: func(c *cli.Context) error {
			fmt.Fprintf(c.App.Writer, c.App.Version)
			return nil
		},
	},
	{
		Name:        "compute",
		Usage:       "Operations on cloud servers, both virtual and bare metal.",
		Subcommands: compute.Get(),
	},
	{
		Name:        "files",
		Usage:       "Object storage for files and media.",
		Subcommands: files.Get(),
	},
	{
		Name:        "networking",
		Usage:       "Software-defined networks, subnets, ports.",
		Subcommands: networking.Get(),
	},
	{
		Name:        "load-balancing",
		Usage:       "Software-defined LBs.",
		Subcommands: loadbalancing.Get(),
	},
}
