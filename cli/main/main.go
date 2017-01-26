package main

import (
	"fmt"
	"os"

	"github.com/gophercloud/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/cli/version"

	"gopkg.in/urfave/cli.v1"
)

func main() {
	cli.HelpPrinter = printHelp
	cli.AppHelpTemplate = appHelpTemplate
	cli.CommandHelpTemplate = commandHelpTemplate
	cli.SubcommandHelpTemplate = subcommandHelpTemplate

	app := cli.NewApp()
	app.Name = "stack"
	app.Flags = globalFlags
	app.Commands = commands
	app.Version = fmt.Sprintf("%v version %v\n   commit: %v\n", app.Name, version.Version, version.Commit)
	app.Usage = "Command-line interface to manage OpenStack resources"
	app.HideVersion = true
	app.EnableBashCompletion = true
	app.BashComplete = globalComplete
	app.CommandNotFound = commandNotFound
	app.Run(os.Args)
}

func globalComplete(ctx *cli.Context) {
	for _, cmd := range ctx.App.Commands {
		fmt.Println(cmd.Name)
	}
	util.CompleteFlags(globalFlags)
}
