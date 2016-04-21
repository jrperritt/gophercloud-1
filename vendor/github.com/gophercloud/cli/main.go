package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/gophercloud/cli/commands/blockstoragecommands"
	"github.com/gophercloud/cli/commands/filescommands"
	"github.com/gophercloud/cli/commands/networkscommands"
	"github.com/gophercloud/cli/commands/orchestrationcommands"
	"github.com/gophercloud/cli/commands/serverscommands"
	"github.com/gophercloud/cli/lib"
	"github.com/gophercloud/cli/setup"
	"github.com/gophercloud/cli/version"

	"github.com/gophercloud/cli/vendor/github.com/codegangsta/cli"
)

func init() {
	lib.Cloud = Context{
		lib.Context{},
	}
}

func main() {
	cli.HelpPrinter = printHelp
	cli.AppHelpTemplate = appHelpTemplate
	cli.CommandHelpTemplate = commandHelpTemplate
	cli.SubcommandHelpTemplate = subcommandHelpTemplate
	app := cli.NewApp()
	app.Name = "stack"
	app.Version = fmt.Sprintf("%v version %v\n   commit: %v\n", app.Name, version.Version, version.Commit)
	app.Usage = Usage()
	app.HideVersion = true
	app.EnableBashCompletion = true
	app.Commands = Cmds()
	app.CommandNotFound = commandNotFound
	app.Run(os.Args)
}

// Usage returns, you guessed it, the usage information
func Usage() string {
	return "Command-line interface to manage resources"
}

// Desc returns, you guessed it, the description
func Desc() string {
	return `A CLI that manages authentication, configures a local setup, and provides workflows for operations on resources.`
}

// Cmds returns a list of commands supported by the tool
func Cmds() []cli.Command {
	return []cli.Command{
		{
			Name:   "configure",
			Usage:  "Interactively create a config file for authentication.",
			Action: configure,
		},
		{
			Name: "init",
			Usage: strings.Join([]string{"Enable tab for command completion.",
				"\tFor Linux and OS X, creates the `rack` man page and sets up",
				"\tcommand completion for the Bash shell. Run `man ./rack.1` to",
				"\tview the generated man page.",
				"\tFor Windows, creates a `posh_autocomplete.ps1` file in the",
				"\t`$HOME/.rack` directory. You must run the file to set up",
				"\tcommand completion."}, "\n"),
			Action: func(c *cli.Context) {
				setup.Init(c)
				man()
			},
		},
		{
			Name:  "version",
			Usage: "Print the version of this binary.",
			Action: func(c *cli.Context) {
				fmt.Fprintf(c.App.Writer, "%v version %v\ncommit: %v\n", c.App.Name, version.Version, version.Commit)
			},
		},
		{
			Name:        "servers",
			Usage:       "Operations on cloud servers, both virtual and bare metal.",
			Subcommands: serverscommands.Get(),
		},
		{
			Name:        "files",
			Usage:       "Object storage for files and media.",
			Subcommands: filescommands.Get(),
		},
		{
			Name:        "networks",
			Usage:       "Software-defined networking.",
			Subcommands: networkscommands.Get(),
		},
		{
			Name: "block-storage",
			Usage: "Block-level storage, exposed as volumes to mount to host servers." +
				"\n\tWork with volumes and their associated snapshots.",
			Subcommands: blockstoragecommands.Get(),
		},
		{
			Name:        "orchestration",
			Usage:       "Use a template language to orchestrate cloud services.",
			Subcommands: orchestrationcommands.Get(),
		},
	}
}
